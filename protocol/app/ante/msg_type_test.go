package ante_test

import (
	"fmt"
	"testing"

	errorsmod "cosmossdk.io/errors"
	"golang.org/x/exp/slices"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"

	customante "github.com/dydxprotocol/v4-chain/protocol/app/ante"
	appmsgs "github.com/dydxprotocol/v4-chain/protocol/app/msgs"
	testante "github.com/dydxprotocol/v4-chain/protocol/testutil/ante"
	testmsgs "github.com/dydxprotocol/v4-chain/protocol/testutil/msgs"

	"github.com/stretchr/testify/require"
)

type txMode string

const (
	checkTx   txMode = "CheckTx"
	reCheckTx txMode = "ReCheckTx"
	deliverTx txMode = "DeliverTx"
)

type testCase struct {
	name        string
	txMode      txMode
	msgs        []sdk.Msg
	expectedErr error
}

var (
	allTxModes = []txMode{checkTx, reCheckTx, deliverTx}

	testMsg = &testdata.TestMsg{Signers: []string{"meh"}}

	invalidReqErrCannotBeEmpty = errorsmod.Wrap(
		sdkerrors.ErrInvalidRequest,
		"msgs cannot be empty",
	)
	invalidReqErrAppInjectedMustBeOnlyMsg = errorsmod.Wrap(
		sdkerrors.ErrInvalidRequest,
		"app-injected msg must be the only msg in a tx",
	)
	invalidReqErrInternalMsg = errorsmod.Wrap(
		sdkerrors.ErrInvalidRequest,
		"internal msg cannot be submitted externally",
	)
	invalidReqErrNestedUnsupportedMsg = errorsmod.Wrap(
		sdkerrors.ErrInvalidRequest,
		fmt.Errorf("Invalid nested msg: unsupported msg type").Error(),
	)
	invalidReqErrNestedAppInjectedMsg = errorsmod.Wrap(
		sdkerrors.ErrInvalidRequest,
		fmt.Errorf("Invalid nested msg: app-injected msg type").Error(),
	)
	invalidReqErrNestedDoubleNested = errorsmod.Wrap(
		sdkerrors.ErrInvalidRequest,
		fmt.Errorf("Invalid nested msg: double-nested msg type").Error(),
	)
	invalidReqErrUnsupportedMsg = errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "unsupported msg")
)

func TestValidateMsgType_Empty(t *testing.T) {
	tests := map[string]struct {
		msgs        []sdk.Msg
		expectedErr map[txMode]error
	}{
		"NilMsg": {
			msgs: nil,

			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   invalidReqErrCannotBeEmpty,
				deliverTx: invalidReqErrCannotBeEmpty,
			},
		},
		"EmptyMsg": {
			msgs: []sdk.Msg{},

			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   invalidReqErrCannotBeEmpty,
				deliverTx: invalidReqErrCannotBeEmpty,
			},
		},
	}

	// Craft all test cases.
	allTestCases := make([]testCase, 0)
	for name, tc := range tests {
		// Run test for each tx mode.
		for _, mode := range allTxModes {
			expectedErr, exists := tc.expectedErr[mode]
			require.True(t, exists)

			testName := fmt.Sprintf("Mode:%s / %s", mode, name)
			allTestCases = append(allTestCases, testCase{testName, mode, tc.msgs, expectedErr})
		}
	}
	require.Len(t, allTestCases, len(tests)*len(allTxModes))

	// Run all test cases.
	for _, tc := range allTestCases {
		runTest(t, tc.name, tc.msgs, tc.txMode, tc.expectedErr)
	}
}

func TestValidateMsgType_AppInjectedMsg(t *testing.T) {
	tests := map[string]struct {
		msgs              []sdk.Msg
		addAppInjectedMsg bool

		expectedErr map[txMode]error
	}{
		"SingleMsg / AppInjectedMsg=true": {
			msgs:              []sdk.Msg{},
			addAppInjectedMsg: true,

			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx: errorsmod.Wrap( // Should only be included in DeliverTx
					sdkerrors.ErrInvalidRequest,
					"app-injected msg must only be included in DeliverTx",
				),
				deliverTx: nil,
			},
		},
		"SingleMsg / AppInjectedMsg=false": {
			msgs:              []sdk.Msg{testMsg},
			addAppInjectedMsg: false,

			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   nil,
				deliverTx: nil,
			},
		},
		"MultiMsgs / AppInjectedMsg=true": {
			msgs:              []sdk.Msg{testMsg},
			addAppInjectedMsg: true,

			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   invalidReqErrAppInjectedMustBeOnlyMsg,
				deliverTx: invalidReqErrAppInjectedMustBeOnlyMsg,
			},
		},
		"MultiMsgs / AppInjectedMsg=false": {
			msgs:              []sdk.Msg{testMsg, testMsg},
			addAppInjectedMsg: false,

			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   nil,
				deliverTx: nil,
			},
		},
	}

	// Craft all test cases.
	testNameFormat := "%s / Mode=%s / AppInjectedMsg=%s"
	allTestCases := make([]testCase, 0)
	numAddAppInjectedMsgTrue := 0
	for name, tc := range tests {
		if tc.addAppInjectedMsg {
			numAddAppInjectedMsgTrue += 1
		}

		// Run test for each tx mode.
		for _, mode := range allTxModes {
			expectedErr, exists := tc.expectedErr[mode]
			require.True(t, exists)

			if tc.addAppInjectedMsg {
				// Run test for each app-injected msg
				appInjectedSampleMsgs := testmsgs.GetNonNilSampleMsgs(appmsgs.AppInjectedMsgSamples)
				for _, sampleMsg := range appInjectedSampleMsgs {
					testName := fmt.Sprintf(testNameFormat, name, mode, sampleMsg.Name)
					testMsgs := append(slices.Clone(tc.msgs), sampleMsg.Msg)
					require.True(t, len(testMsgs) > 0 && len(testMsgs) <= 3)
					allTestCases = append(allTestCases, testCase{testName, mode, testMsgs, expectedErr})
				}
			} else {
				testName := fmt.Sprintf(testNameFormat, name, mode, "none")
				require.True(t, len(tc.msgs) > 0 && len(tc.msgs) <= 2)
				allTestCases = append(allTestCases, testCase{testName, mode, tc.msgs, expectedErr})
			}
		}
	}

	// NumOfTxModes * ((NumOfAddInjectedMsg=true * NumAppInjectedMsgSample) + NumAddInjectedMsg=false)
	numAppInjectedMsgSamples := len(testmsgs.GetNonNilSampleMsgs(appmsgs.AppInjectedMsgSamples))
	expectedNumTestCases := len(allTxModes) *
		((numAddAppInjectedMsgTrue * numAppInjectedMsgSamples) + (len(tests) - numAddAppInjectedMsgTrue))
	require.Len(t, allTestCases, expectedNumTestCases)

	// Run tests.
	for _, tc := range allTestCases {
		runTest(t, tc.name, tc.msgs, tc.txMode, tc.expectedErr)
	}
}

func TestValidateMsgType_InternalOnlyMsg(t *testing.T) {
	tests := map[string]struct {
		msgs           []sdk.Msg
		addInternalMsg bool
		expectedErr    map[txMode]error
	}{
		"SingleMsg / InternalMsg=true": {
			msgs:           []sdk.Msg{},
			addInternalMsg: true,

			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   invalidReqErrInternalMsg,
				deliverTx: invalidReqErrInternalMsg,
			},
		},
		"SingleMsg / InternalMsg=false": {
			msgs:           []sdk.Msg{testMsg},
			addInternalMsg: false,

			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   nil,
				deliverTx: nil,
			},
		},
		"MultiMsgs / InternalMsg=true": {
			msgs:           []sdk.Msg{testMsg},
			addInternalMsg: true,

			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   invalidReqErrInternalMsg,
				deliverTx: invalidReqErrInternalMsg,
			},
		},
		"MultiMsgs / InternalMsg=false": {
			msgs:           []sdk.Msg{testMsg, testMsg},
			addInternalMsg: false,

			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   nil,
				deliverTx: nil,
			},
		},
	}

	// Craft all test cases.
	testNameFormat := "%s / Mode=%s / InternalMsg=%s"
	allTestCases := make([]testCase, 0)
	numAddInternalMsgTrue := 0
	for name, tc := range tests {
		if tc.addInternalMsg {
			numAddInternalMsgTrue += 1
		}

		// Run test for each tx mode.
		for _, mode := range allTxModes {
			expectedErr, exists := tc.expectedErr[mode]
			require.True(t, exists)

			if tc.addInternalMsg {
				// Run test for each internal msg
				internalSampleMsgs := testmsgs.GetNonNilSampleMsgs(appmsgs.InternalMsgSamplesAll)
				for _, sampleMsg := range internalSampleMsgs {
					testName := fmt.Sprintf(testNameFormat, name, mode, sampleMsg.Name)
					testMsgs := append(slices.Clone(tc.msgs), sampleMsg.Msg)
					require.True(t, len(testMsgs) > 0 && len(testMsgs) <= 3)
					allTestCases = append(allTestCases, testCase{testName, mode, testMsgs, expectedErr})
				}
			} else {
				testName := fmt.Sprintf(testNameFormat, name, mode, "none")
				require.True(t, len(tc.msgs) > 0 && len(tc.msgs) <= 2)
				allTestCases = append(allTestCases, testCase{testName, mode, tc.msgs, expectedErr})
			}
		}
	}

	// NumOfTxModes * ((NumOfAddInternalMsg=true * NumInternalMsgSample) + NumOfAddInternalMsg=false)
	numInternalMsgSamples := len(testmsgs.GetNonNilSampleMsgs(appmsgs.InternalMsgSamplesAll))
	expectedNumTestCases := len(allTxModes) *
		((numAddInternalMsgTrue * numInternalMsgSamples) + (len(tests) - numAddInternalMsgTrue))
	require.Len(t, allTestCases, expectedNumTestCases)

	// Run all test cases.
	for _, tc := range allTestCases {
		runTest(t, tc.name, tc.msgs, tc.txMode, tc.expectedErr)
	}
}

func TestValidateMsgType_NestedMsg(t *testing.T) {
	// Test cases.
	tests := map[string]struct {
		msg         sdk.Msg
		expectedErr map[txMode]error
	}{
		"Success: empty inner msgs": {
			msg: testmsgs.MsgSubmitProposalWithEmptyInner,
			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   nil,
				deliverTx: nil,
			},
		},
		"Fails: unsupported inner msg": {
			msg: testmsgs.MsgSubmitProposalWithUnsupportedInner,
			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   invalidReqErrNestedUnsupportedMsg,
				deliverTx: invalidReqErrNestedUnsupportedMsg,
			},
		},
		"Fails: app-injected inner msg": {
			msg: testmsgs.MsgSubmitProposalWithAppInjectedInner,
			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   invalidReqErrNestedAppInjectedMsg,
				deliverTx: invalidReqErrNestedAppInjectedMsg,
			},
		},
		"Fails: double-nested inner msg": {
			msg: testmsgs.MsgSubmitProposalWithDoubleNestedInner,
			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   invalidReqErrNestedDoubleNested,
				deliverTx: invalidReqErrNestedDoubleNested,
			},
		},
		"Success: single valid inner msg": {
			msg: testmsgs.MsgSubmitProposalWithUpgrade,
			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   nil,
				deliverTx: nil,
			},
		},
		"Success: multi valid inner msgs": {
			msg: testmsgs.MsgSubmitProposalWithUpgradeAndCancel,
			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   nil,
				deliverTx: nil,
			},
		},
	}

	// Run.
	testNameFormat := "%s / %s / Mode=%s"
	for name, tc := range tests {
		for _, mode := range allTxModes {
			expectedErr, exists := tc.expectedErr[mode]
			require.True(t, exists)

			singleMsgTestName := fmt.Sprintf(testNameFormat, name, "SingleMsg", mode)
			runTest(t, singleMsgTestName, []sdk.Msg{tc.msg}, mode, expectedErr)

			multiMsgsTestName := fmt.Sprintf(testNameFormat, name, "MultiMsgs", mode)
			multiMsgs := []sdk.Msg{testMsg, tc.msg, testMsg} // add extra test msg.
			runTest(t, multiMsgsTestName, multiMsgs, mode, expectedErr)
		}
	}
}

func TestValidateMsgType_UnsupportedMsg(t *testing.T) {
	tests := map[string]struct {
		msgs              []sdk.Msg
		addUnsupportedMsg bool
		expectedErr       map[txMode]error
	}{
		"SingleMsg / UnsupportedMsg=true": {
			msgs:              []sdk.Msg{},
			addUnsupportedMsg: true,

			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   invalidReqErrUnsupportedMsg,
				deliverTx: invalidReqErrUnsupportedMsg,
			},
		},
		"SingleMsg / UnsupportedMsg=false": {
			msgs:              []sdk.Msg{testMsg},
			addUnsupportedMsg: false,

			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   nil,
				deliverTx: nil,
			},
		},
		"MultiMsgs / UnsupportedMsg=true": {
			msgs:              []sdk.Msg{testMsg},
			addUnsupportedMsg: true,

			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   invalidReqErrUnsupportedMsg,
				deliverTx: invalidReqErrUnsupportedMsg,
			},
		},
		"MultiMsgs / UnsupportedMsg=false": {
			msgs:              []sdk.Msg{testMsg, testMsg},
			addUnsupportedMsg: false,

			expectedErr: map[txMode]error{
				reCheckTx: nil, // ReCheck skips the AnteHandler.
				checkTx:   nil,
				deliverTx: nil,
			},
		},
	}

	// Craft all test cases.
	testNameFormat := "%s / Mode=%s / UnsupportedMsg=%s"
	allTestCases := make([]testCase, 0)
	addUnsupportedMsgTrue := 0
	for name, tc := range tests {
		if tc.addUnsupportedMsg {
			addUnsupportedMsgTrue += 1
		}

		// Run test for each tx mode.
		for _, mode := range allTxModes {
			expectedErr, exists := tc.expectedErr[mode]
			require.True(t, exists)

			if tc.addUnsupportedMsg {
				// Run test for each unsupported msg
				unsupportedSampleMsgs := testmsgs.GetNonNilSampleMsgs(appmsgs.UnsupportedMsgSamples)
				for _, sampleMsg := range unsupportedSampleMsgs {
					testName := fmt.Sprintf(testNameFormat, name, mode, sampleMsg.Name)
					testMsgs := append(slices.Clone(tc.msgs), sampleMsg.Msg)
					require.True(t, len(testMsgs) > 0 && len(testMsgs) <= 3)
					allTestCases = append(allTestCases, testCase{testName, mode, testMsgs, expectedErr})
				}
			} else {
				testName := fmt.Sprintf(testNameFormat, name, mode, "none")
				require.True(t, len(tc.msgs) > 0 && len(tc.msgs) <= 2)
				allTestCases = append(allTestCases, testCase{testName, mode, tc.msgs, expectedErr})
			}
		}
	}

	// NumOfTxModes * ((NumOfAddUnsupportedMsgTrue=true * NumUnsupportedMsgSample) + NumOfAddUnsupportedMsg=false)
	numUnsupportedMsgSamples := len(testmsgs.GetNonNilSampleMsgs(appmsgs.UnsupportedMsgSamples))
	expectedNumTestCases := len(allTxModes) *
		((addUnsupportedMsgTrue * numUnsupportedMsgSamples) + (len(tests) - addUnsupportedMsgTrue))
	require.Len(t, allTestCases, expectedNumTestCases)

	// Run all test cases.
	for _, tc := range allTestCases {
		runTest(t, tc.name, tc.msgs, tc.txMode, tc.expectedErr)
	}
}

func getCtxWithTxMode(mode txMode, suite *testante.AnteTestSuite) sdk.Context {
	switch mode {
	case checkTx:
		return suite.Ctx.WithIsCheckTx(true).WithIsReCheckTx(false)
	case reCheckTx:
		return suite.Ctx.WithIsCheckTx(false).WithIsReCheckTx(true)
	case deliverTx:
		return suite.Ctx.WithIsCheckTx(false).WithIsReCheckTx(false)
	default:
		panic("invalid mode")
	}
}

func runTest(t *testing.T, name string, msgs []sdk.Msg, mode txMode, expectedErr error) {
	t.Run(name, func(t *testing.T) {
		// Setup.
		suite := testante.SetupTestSuite(t, true)
		suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder()
		vmt := customante.NewValidateMsgTypeDecorator()
		antehandler := sdk.ChainAnteDecorators(vmt)
		require.NoError(t, suite.TxBuilder.SetMsgs(msgs...))
		// Empty private key, so tx's signature should be empty.
		privs, accNums, accSeqs := []cryptotypes.PrivKey{}, []uint64{}, []uint64{}
		tx, err := suite.CreateTestTx(
			suite.Ctx,
			privs,
			accNums,
			accSeqs,
			suite.Ctx.ChainID(),
			signing.SignMode_SIGN_MODE_DIRECT,
		)
		require.NoError(t, err)
		suite.Ctx = getCtxWithTxMode(mode, suite)

		// Run.
		_, err = antehandler(suite.Ctx, tx, false)

		// Verify.
		if expectedErr != nil {
			require.EqualError(t, err, expectedErr.Error())
		} else {
			require.NoError(t, err)
		}
	})
}
