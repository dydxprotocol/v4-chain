package msgs

import (
	upgrade "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	marketmap "github.com/dydxprotocol/slinky/x/marketmap/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	sending "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
)

func init() {
	testEncodingCfg := encoding.GetTestEncodingCfg()
	testTxBuilder := testEncodingCfg.TxConfig.NewTxBuilder()

	_ = testTxBuilder.SetMsgs(MsgSoftwareUpgrade)
	MsgSoftwareUpgradeTxBytes, _ = testEncodingCfg.TxConfig.TxEncoder()(testTxBuilder.GetTx())

	_ = testTxBuilder.SetMsgs(MsgCancelUpgrade)
	MsgCancelUpgradeTxBytes, _ = testEncodingCfg.TxConfig.TxEncoder()(testTxBuilder.GetTx())

	_ = testTxBuilder.SetMsgs(MsgSubmitProposalWithEmptyInner)
	MsgSubmitProposalWithEmptyInnerTxBytes, _ = testEncodingCfg.TxConfig.TxEncoder()(testTxBuilder.GetTx())

	_ = testTxBuilder.SetMsgs(MsgSubmitProposalWithUnsupportedInner)
	MsgSubmitProposalWithUnsupportedInnerTxBytes, _ = testEncodingCfg.TxConfig.TxEncoder()(testTxBuilder.GetTx())

	_ = testTxBuilder.SetMsgs(MsgSubmitProposalWithAppInjectedInner)
	MsgSubmitProposalWithAppInjectedInnerTxBytes, _ = testEncodingCfg.TxConfig.TxEncoder()(testTxBuilder.GetTx())

	_ = testTxBuilder.SetMsgs(MsgSubmitProposalWithDoubleNestedInner)
	MsgSubmitProposalWithDoubleNestedInnerTxBytes, _ = testEncodingCfg.TxConfig.TxEncoder()(testTxBuilder.GetTx())

	_ = testTxBuilder.SetMsgs(&MsgExecWithUnsupportedInner)
	MsgExecWithUnsupportedInnerTxBytes, _ = testEncodingCfg.TxConfig.TxEncoder()(testTxBuilder.GetTx())

	_ = testTxBuilder.SetMsgs(&MsgExecWithAppInjectedInner)
	MsgExecWithAppInjectedInnerTxBytes, _ = testEncodingCfg.TxConfig.TxEncoder()(testTxBuilder.GetTx())

	_ = testTxBuilder.SetMsgs(&MsgExecWithDoubleNestedInner)
	MsgExecWithDoubleNestedInnerTxBytes, _ = testEncodingCfg.TxConfig.TxEncoder()(testTxBuilder.GetTx())

	_ = testTxBuilder.SetMsgs(&MsgExecWithDydxMessage)
	MsgExecWithDydxMessageTxBytes, _ = testEncodingCfg.TxConfig.TxEncoder()(testTxBuilder.GetTx())

	_ = testTxBuilder.SetMsgs(&MsgExecWithSlinkyMessage)
	MsgExecWithSlinkyMessageTxBytes, _ = testEncodingCfg.TxConfig.TxEncoder()(testTxBuilder.GetTx())

	_ = testTxBuilder.SetMsgs(MsgSubmitProposalWithUpgrade)
	MsgSubmitProposalWithUpgradeTxBytes, _ = testEncodingCfg.TxConfig.TxEncoder()(testTxBuilder.GetTx())

	_ = testTxBuilder.SetMsgs(MsgSubmitProposalWithUpgradeAndCancel)
	MsgSubmitProposalWithUpgradeAndCancelTxBytes, _ = testEncodingCfg.TxConfig.TxEncoder()(testTxBuilder.GetTx())
}

const (
	testMetadata = "test-metadata"
	testTitle    = "test-title"
	testSummary  = "test-summary"
)

var (
	testProposer = constants.Bob_Num0.Owner

	// Inner msgs.
	MsgSoftwareUpgrade = &upgrade.MsgSoftwareUpgrade{
		Authority: constants.Bob_Num0.Owner,
		Plan: upgrade.Plan{
			Name:   "test-plan",
			Height: 10,
			Info:   "test-info",
		},
	}
	MsgSoftwareUpgradeTxBytes []byte

	MsgCancelUpgrade = &upgrade.MsgCancelUpgrade{
		Authority: constants.Bob_Num0.Owner,
	}
	MsgCancelUpgradeTxBytes []byte

	// Invalid MsgSubmitProposals
	MsgSubmitProposalWithEmptyInner, _ = gov.NewMsgSubmitProposal(
		[]sdk.Msg{}, nil, testProposer, testMetadata, testTitle, testSummary, false)
	MsgSubmitProposalWithEmptyInnerTxBytes []byte

	MsgSubmitProposalWithUnsupportedInner, _ = gov.NewMsgSubmitProposal(
		[]sdk.Msg{GovBetaMsgSubmitProposal}, nil, testProposer, testMetadata, testTitle, testSummary, false)
	MsgSubmitProposalWithUnsupportedInnerTxBytes []byte

	MsgSubmitProposalWithAppInjectedInner, _ = gov.NewMsgSubmitProposal(
		[]sdk.Msg{&prices.MsgUpdateMarketPrices{}}, nil, testProposer, testMetadata, testTitle, testSummary, false)
	MsgSubmitProposalWithAppInjectedInnerTxBytes []byte

	MsgSubmitProposalWithDoubleNestedInner, _ = gov.NewMsgSubmitProposal(
		[]sdk.Msg{MsgSubmitProposalWithUpgradeAndCancel}, nil, testProposer, testMetadata, testTitle, testSummary, false)
	MsgSubmitProposalWithDoubleNestedInnerTxBytes []byte

	// Invalid MsgExec
	MsgExecWithUnsupportedInner = authz.NewMsgExec(
		constants.AliceAccAddress,
		[]sdk.Msg{GovBetaMsgSubmitProposal},
	)
	MsgExecWithUnsupportedInnerTxBytes []byte

	MsgExecWithAppInjectedInner = authz.NewMsgExec(
		constants.AliceAccAddress,
		[]sdk.Msg{&prices.MsgUpdateMarketPrices{}},
	)
	MsgExecWithAppInjectedInnerTxBytes []byte

	MsgExecWithDoubleNestedInner = authz.NewMsgExec(
		constants.AliceAccAddress,
		[]sdk.Msg{MsgSubmitProposalWithUpgradeAndCancel},
	)
	MsgExecWithDoubleNestedInnerTxBytes []byte

	MsgExecWithDydxMessage = authz.NewMsgExec(
		constants.AliceAccAddress,
		[]sdk.Msg{&sending.MsgCreateTransfer{}},
	)
	MsgExecWithDydxMessageTxBytes []byte

	MsgExecWithSlinkyMessage = authz.NewMsgExec(
		constants.AliceAccAddress,
		[]sdk.Msg{&marketmap.MsgUpsertMarkets{}},
	)
	MsgExecWithSlinkyMessageTxBytes []byte

	// Valid MsgSubmitProposals
	MsgSubmitProposalWithUpgrade, _ = gov.NewMsgSubmitProposal(
		[]sdk.Msg{MsgSoftwareUpgrade}, nil, testProposer, testMetadata, testTitle, testSummary, false)
	MsgSubmitProposalWithUpgradeTxBytes []byte

	MsgSubmitProposalWithUpgradeAndCancel, _ = gov.NewMsgSubmitProposal(
		[]sdk.Msg{
			MsgSoftwareUpgrade,
			MsgCancelUpgrade,
		}, nil, testProposer, testMetadata, testTitle, testSummary, false)
	MsgSubmitProposalWithUpgradeAndCancelTxBytes []byte
)
