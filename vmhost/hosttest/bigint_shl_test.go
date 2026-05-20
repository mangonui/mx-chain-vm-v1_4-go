package hostCoretest

import (
	"testing"

	contextmock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-v1_4-go/testcommon"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
)

func TestBigIntShlHugeShiftFailsBeforeAllocation(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("big-int-shl", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(100000).
			WithFunction("hugeShift").
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				OutOfGas().
				ReturnMessage(vmhost.ErrNotEnoughGas.Error()).
				HasRuntimeErrors(vmhost.ErrNotEnoughGas.Error())
		})
}
