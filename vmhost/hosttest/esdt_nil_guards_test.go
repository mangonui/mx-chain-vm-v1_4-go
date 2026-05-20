package hostCoretest

import (
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/esdt"
	contextmock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-v1_4-go/testcommon"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
)

func TestESDTNilGuards_FailClosedForFlatAndManagedHooks(t *testing.T) {

	testCases := []struct {
		name      string
		function  string
		esdtToken *esdt.ESDigitalToken
	}{
		{
			name:      "legacy balance nil token",
			function:  "legacyBalance",
			esdtToken: nil,
		},
		{
			name:      "legacy balance nil value",
			function:  "legacyBalance",
			esdtToken: &esdt.ESDigitalToken{},
		},
		{
			name:      "legacy token data nil token",
			function:  "legacyTokenData",
			esdtToken: nil,
		},
		{
			name:      "legacy token data nil value",
			function:  "legacyTokenData",
			esdtToken: &esdt.ESDigitalToken{},
		},
		{
			name:      "managed balance nil token",
			function:  "managedBalance",
			esdtToken: nil,
		},
		{
			name:      "managed balance nil value",
			function:  "managedBalance",
			esdtToken: &esdt.ESDigitalToken{},
		},
		{
			name:      "managed token data nil token",
			function:  "managedTokenData",
			esdtToken: nil,
		},
		{
			name:      "managed token data nil value",
			function:  "managedTokenData",
			esdtToken: &esdt.ESDigitalToken{},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {

			runESDTNilGuardFixture(t, testCase.function, testCase.esdtToken, func(verify *test.VMOutputVerifier) {
				verify.
					ExecutionFailed().
					ReturnMessage(vmhost.ErrNilESDTData.Error()).
					HasRuntimeErrors(vmhost.ErrNilESDTData.Error())
			})
		})
	}
}

func TestBigIntGetESDTExternalBalance_NilGuardReturnsWithoutPanic(t *testing.T) {

	testCases := []struct {
		name      string
		esdtToken *esdt.ESDigitalToken
	}{
		{
			name:      "nil token",
			esdtToken: nil,
		},
		{
			name:      "nil value",
			esdtToken: &esdt.ESDigitalToken{},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {

			runESDTNilGuardFixture(t, "bigIntExternalBalance", testCase.esdtToken, func(verify *test.VMOutputVerifier) {
				verify.Ok()
			})
		})
	}
}

func runESDTNilGuardFixture(
	t *testing.T,
	function string,
	esdtToken *esdt.ESDigitalToken,
	assertResults func(verify *test.VMOutputVerifier),
) {
	t.Helper()

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("esdt-nil-guards", "../../"))).
		WithSetup(func(_ vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.GetESDTTokenCalled = func(_ []byte, _ []byte, _ uint64) (*esdt.ESDigitalToken, error) {
				return esdtToken, nil
			}
		}).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(100000).
			WithFunction(function).
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			assertResults(verify)
		})
}

func TestESDTNilGuardFixtureSanity(t *testing.T) {

	runESDTNilGuardFixture(t, "legacyBalance", &esdt.ESDigitalToken{Value: big.NewInt(7)}, func(verify *test.VMOutputVerifier) {
		verify.Ok()
	})
}
