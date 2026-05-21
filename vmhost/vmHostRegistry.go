package vmhost

import (
	"sync"
	"sync/atomic"
)

// ISSUE-013: typed handle registry for VMHost references.
//
// Same pattern as DESIGN-ISSUE-013 / DESIGN-ISSUE-011: replace the
// `uintptr → unsafe.Pointer → VMHost` round-trip with a `uint64` handle
// looked up in a process-wide registry. The registry holds a strong
// reference to the VMHost, keeping it alive regardless of the
// runtimeContext's GC-eligibility.
//
// Lifecycle. The legacy v1_x VMs do NOT expose an explicit
// runtimeContext destroy method, but vmHost.Close() now explicitly asks
// the runtime context to release its handle. This keeps the registry
// bounded even in integrations that repeatedly create and close legacy
// VM hosts in one process.
//
// Concurrency. Same RWMutex pattern as wasmer2's vmHooksRegistry.
// Hook callbacks are read-heavy (Lookup); registration happens once
// per runtimeContext.
//
// Panic-on-miss. The exported GetVMHost helper PANICS with a clear
// message if a handle isn't registered. Every VMHost-using hook in
// vmhost/vmhooks/ assumes the returned interface is non-nil and
// immediately invokes a method on it; returning nil would defer the
// failure to a cryptic nil-interface dispatch elsewhere. Panicking
// here surfaces the bug clearly.

type vmHostRegistry struct {
	mu      sync.RWMutex
	nextID  uint64
	entries map[uint64]VMHost
}

var globalVMHostRegistry = &vmHostRegistry{
	entries: make(map[uint64]VMHost),
}

// Register inserts the host and returns a stable monotonic handle.
// Handles are never reused; a released handle's value will never be
// reissued, so a stale-handle lookup always returns nil (safest
// possible failure mode).
func (r *vmHostRegistry) Register(host VMHost) uint64 {
	id := atomic.AddUint64(&r.nextID, 1)
	r.mu.Lock()
	r.entries[id] = host
	r.mu.Unlock()
	return id
}

// Lookup resolves a handle. Returns nil for unregistered or
// already-released handles. Caller-side nil handling is the security
// boundary — see lookupVMHostOrPanic.
func (r *vmHostRegistry) Lookup(id uint64) VMHost {
	r.mu.RLock()
	h := r.entries[id]
	r.mu.RUnlock()
	return h
}

// Release frees a handle. Subsequent Lookup(id) returns nil. Idempotent
// (delete from a Go map on a missing key is a no-op).
func (r *vmHostRegistry) Release(id uint64) {
	r.mu.Lock()
	delete(r.entries, id)
	r.mu.Unlock()
}

// lookupVMHostOrPanic resolves a handle and panics with a clear message
// on miss. Used by GetVMHost — see panic-on-miss rationale in the
// package-level comment above.
func lookupVMHostOrPanic(handle uint64) VMHost {
	host := globalVMHostRegistry.Lookup(handle)
	if host == nil {
		panic(vmHostRegistryMissPanicMessage(handle))
	}
	return host
}

// RegisterVMHostHandle is a package-public wrapper used by
// vmhost/contexts to register a host on first SetContextData call.
// Returns a stable handle that survives the runtimeContext's lifetime.
// See lifecycle note in the package-level comment.
func RegisterVMHostHandle(host VMHost) uint64 {
	return globalVMHostRegistry.Register(host)
}

// ReleaseVMHostHandle releases a previously registered VMHost handle.
// It is idempotent; releasing an unknown or already-released handle is a no-op.
func ReleaseVMHostHandle(handle uint64) {
	globalVMHostRegistry.Release(handle)
}

func vmHostRegistryMissPanicMessage(handle uint64) string {
	return "vmhost: VMHost handle not found in registry; runtimeContext lifecycle bug or stale wasmer context (handle=" +
		uint64ToDecimalString(handle) + ")"
}

// uint64ToDecimalString avoids strconv import bloat. Allocation-free
// for the success path (only called inside the panic branch).
func uint64ToDecimalString(v uint64) string {
	if v == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	return string(buf[i:])
}
