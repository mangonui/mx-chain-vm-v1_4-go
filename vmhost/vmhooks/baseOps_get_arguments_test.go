package vmhooks

import (
	"encoding/binary"
	"testing"

	"github.com/multiversx/mx-chain-vm-v1_4-go/config"
	contextmock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/context"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	"github.com/stretchr/testify/require"
)

type countingRuntimeContext struct {
	*contextmock.RuntimeContextMock
	memLoadCalls         int
	memLoadMultipleCalls int
}

func (r *countingRuntimeContext) MemLoad(offset int32, length int32) ([]byte, error) {
	r.memLoadCalls++
	return r.RuntimeContextMock.MemLoad(offset, length)
}

func (r *countingRuntimeContext) MemLoadMultiple(offset int32, lengths []int32) ([][]byte, error) {
	r.memLoadMultipleCalls++
	return r.RuntimeContextMock.MemLoadMultiple(offset, lengths)
}

func newGetArgumentsHost(gasLeft uint64, dataCopyGas uint64, lengthBytes []byte, payload [][]byte) (*contextmock.VMHostMock, *countingRuntimeContext, *contextmock.MeteringContextMock) {
	runtime := &countingRuntimeContext{
		RuntimeContextMock: &contextmock.RuntimeContextMock{
			MemLoadResult:         lengthBytes,
			MemLoadMultipleResult: payload,
		},
	}
	metering := &contextmock.MeteringContextMock{
		GasCost: &config.GasCost{
			BaseOperationCost: config.BaseOperationCost{
				DataCopyPerByte: dataCopyGas,
			},
		},
		GasLeftMock: gasLeft,
	}

	host := &contextmock.VMHostMock{
		RuntimeContext:  runtime,
		MeteringContext: metering,
	}

	return host, runtime, metering
}

func encodeArgumentLengths(lengths ...int32) []byte {
	data := make([]byte, len(lengths)*4)
	for i, length := range lengths {
		binary.LittleEndian.PutUint32(data[i*4:], uint32(length))
	}
	return data
}

func TestGetArgumentsFromMemory_PreChargesBeforeLoadsAndReturnsData(t *testing.T) {
	t.Parallel()

	host, runtime, metering := newGetArgumentsHost(
		1000,
		2,
		encodeArgumentLengths(3, 5),
		[][]byte{[]byte("one"), []byte("three")},
	)

	data, actualLen, err := getArgumentsFromMemory(host, 2, 10, 20)

	require.NoError(t, err)
	require.Equal(t, [][]byte{[]byte("one"), []byte("three")}, data)
	require.Equal(t, int32(8), actualLen)
	require.Equal(t, 1, runtime.memLoadCalls)
	require.Equal(t, 1, runtime.memLoadMultipleCalls)
	require.Equal(t, uint64(968), metering.GasLeft())
}

func TestGetArgumentsFromMemory_AllowsZeroArgumentsWithNoGasLeft(t *testing.T) {
	t.Parallel()

	host, runtime, metering := newGetArgumentsHost(0, 10, nil, [][]byte{})

	data, actualLen, err := getArgumentsFromMemory(host, 0, 10, 20)

	require.NoError(t, err)
	require.Empty(t, data)
	require.Zero(t, actualLen)
	require.Equal(t, 1, runtime.memLoadCalls)
	require.Equal(t, 1, runtime.memLoadMultipleCalls)
	require.Zero(t, metering.GasLeft())
}

func TestGetArgumentsFromMemory_RejectsNegativeAndTooManyArgumentsBeforeLoading(t *testing.T) {
	t.Parallel()

	host, runtime, _ := newGetArgumentsHost(1000, 2, nil, nil)

	_, _, err := getArgumentsFromMemory(host, -1, 10, 20)
	require.ErrorContains(t, err, "negative numArguments")
	require.Zero(t, runtime.memLoadCalls)
	require.Zero(t, runtime.memLoadMultipleCalls)

	_, _, err = getArgumentsFromMemory(host, maxNumArgumentsFromMemory+1, 10, 20)
	require.ErrorContains(t, err, "exceeds maximum")
	require.Zero(t, runtime.memLoadCalls)
	require.Zero(t, runtime.memLoadMultipleCalls)
}

func TestGetArgumentsFromMemory_FailsBeforeMetadataLoadWhenMetadataGasIsUnavailable(t *testing.T) {
	t.Parallel()

	host, runtime, _ := newGetArgumentsHost(80, 10, encodeArgumentLengths(1, 1), nil)

	_, _, err := getArgumentsFromMemory(host, 2, 10, 20)

	require.ErrorIs(t, err, vmhost.ErrNotEnoughGas)
	require.Zero(t, runtime.memLoadCalls)
	require.Zero(t, runtime.memLoadMultipleCalls)
}

func TestGetArgumentsFromMemory_RejectsMalformedLengthsBeforePayloadLoad(t *testing.T) {
	t.Parallel()

	t.Run("negative length", func(t *testing.T) {
		host, runtime, _ := newGetArgumentsHost(1000, 1, encodeArgumentLengths(-1), nil)

		_, _, err := getArgumentsFromMemory(host, 1, 10, 20)

		require.ErrorContains(t, err, "negative argument length")
		require.Equal(t, 1, runtime.memLoadCalls)
		require.Zero(t, runtime.memLoadMultipleCalls)
	})

	t.Run("sum overflow", func(t *testing.T) {
		host, runtime, _ := newGetArgumentsHost(1000, 1, encodeArgumentLengths(2147483647, 1), nil)

		_, _, err := getArgumentsFromMemory(host, 2, 10, 20)

		require.ErrorContains(t, err, "total argument bytes exceeds int32 max")
		require.Equal(t, 1, runtime.memLoadCalls)
		require.Zero(t, runtime.memLoadMultipleCalls)
	})
}

func TestGetArgumentsFromMemory_FailsBeforePayloadLoadWhenPayloadGasIsUnavailable(t *testing.T) {
	t.Parallel()

	host, runtime, _ := newGetArgumentsHost(100, 10, encodeArgumentLengths(100), nil)

	_, _, err := getArgumentsFromMemory(host, 1, 10, 20)

	require.ErrorIs(t, err, vmhost.ErrNotEnoughGas)
	require.Equal(t, 1, runtime.memLoadCalls)
	require.Zero(t, runtime.memLoadMultipleCalls)
}
