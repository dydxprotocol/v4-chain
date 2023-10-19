package ante_test

import (
	"fmt"
	"testing"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	appmsgs "github.com/dydxprotocol/v4-chain/protocol/app/msgs"
	"github.com/dydxprotocol/v4-chain/protocol/lib/ante"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	testante "github.com/dydxprotocol/v4-chain/protocol/testutil/ante"
	testmsgs "github.com/dydxprotocol/v4-chain/protocol/testutil/msgs"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testMsg = &testdata.TestMsg{Signers: []string{"meh"}}
)

func TestValidateMsgType_AppInjectedMsg(t *testing.T) {
	tests := map[string]struct {
		msgs              []sdk.Msg
		addAppInjectedMsg bool

		expectSkip bool
	}{
		"No Skip: No msg": {
			msgs:              []sdk.Msg{},
			addAppInjectedMsg: false,

			expectSkip: false,
		},
		"Yes Skip: Single msg / Has app-injected msg": {
			msgs:              []sdk.Msg{},
			addAppInjectedMsg: true,

			expectSkip: true,
		},
		"No Skip: Single msg / No app-injected msg": {
			msgs:              []sdk.Msg{testMsg},
			addAppInjectedMsg: false,

			expectSkip: false,
		},
		"No Skip: Mult msgs / Has app-injected msg": {
			msgs:              []sdk.Msg{testMsg},
			addAppInjectedMsg: true,

			expectSkip: false,
		},
		"No Skip: mult msgs / No app-injected msg": {
			msgs:              []sdk.Msg{testMsg, testMsg},
			addAppInjectedMsg: false,

			expectSkip: false,
		},
	}

	type testCase struct {
		name       string
		msgs       []sdk.Msg
		expectSkip bool
	}
	allTestCases := make([]testCase, 0)

	// Craft all test cases.
	for name, tc := range tests {
		msgs := tc.msgs
		testName := ""
		if tc.addAppInjectedMsg {
			// Run test for each app-injected msg
			appInjectedSampleMsgs := testmsgs.GetNonNilSampleMsgs(appmsgs.AppInjectedMsgSamples)
			for _, sampleMsg := range appInjectedSampleMsgs {
				testName = fmt.Sprintf("AppInjectedMsg:%s / %s", sampleMsg.Name, name)
				msgs = append(tc.msgs, sampleMsg.Msg)
				allTestCases = append(allTestCases, testCase{testName, msgs, tc.expectSkip})
			}
		} else {
			testName = fmt.Sprintf("AppInjectedMsg:%s / %s", "none", name)
			allTestCases = append(allTestCases, testCase{testName, msgs, tc.expectSkip})
		}
	}

	for _, tc := range allTestCases {
		runTest(t, tc.name, tc.msgs, tc.expectSkip)
	}
}

func runTest(t *testing.T, name string, msgs []sdk.Msg, expectSkip bool) {
	t.Run(name, func(t *testing.T) {
		suite := testante.SetupTestSuite(t, true)
		suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder()

		mockAntehandler := &mocks.AnteDecorator{}
		mockAntehandler.On("AnteHandle", suite.Ctx, mock.Anything, false, mock.Anything).
			Return(suite.Ctx, nil)

		wrappedHandler := ante.NewAppInjectedMsgAnteWrapper(mockAntehandler)
		antehandler := sdk.ChainAnteDecorators(wrappedHandler)

		require.NoError(t, suite.TxBuilder.SetMsgs(msgs...))

		// Empty private key, so tx's signature should be empty.
		privs, accNums, accSeqs := []cryptotypes.PrivKey{}, []uint64{}, []uint64{}

		tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.Ctx.ChainID())
		require.NoError(t, err)

		resultCtx, err := antehandler(suite.Ctx, tx, false)
		require.NoError(t, err)
		require.Equal(t, suite.Ctx, resultCtx)

		if expectSkip {
			mockAntehandler.AssertNotCalled(
				t,
				"AnteHandle",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			)
		} else {
			mockAntehandler.AssertCalled(
				t,
				"AnteHandle",
				suite.Ctx,
				tx,
				false,
				mock.Anything,
			)
			mockAntehandler.AssertNumberOfCalls(t, "AnteHandle", 1)
		}
	})
}
