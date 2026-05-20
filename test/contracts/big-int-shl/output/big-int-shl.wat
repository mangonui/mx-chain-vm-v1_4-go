(module
  (type (;0;) (func (param i64) (result i32)))
  (type (;1;) (func (param i32 i32 i32)))
  (type (;2;) (func))
  (import "env" "bigIntNew" (func (;0;) (type 0)))
  (import "env" "bigIntShl" (func (;1;) (type 1)))
  (func (;2;) (type 2))
  (func (;3;) (type 2)
    (local i32 i32)
    i64.const 1
    call 0
    local.set 0
    i64.const 0
    call 0
    local.set 1
    local.get 1
    local.get 0
    i32.const 2147483647
    call 1)
  (table (;0;) 1 1 funcref)
  (memory (;0;) 2)
  (global (;0;) (mut i32) (i32.const 66576))
  (export "memory" (memory 0))
  (export "init" (func 2))
  (export "hugeShift" (func 3)))
