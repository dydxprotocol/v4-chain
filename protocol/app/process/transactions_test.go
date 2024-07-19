package process_test

import (
	"testing"

	errorsmod "cosmossdk.io/errors"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/process"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDecodeProcessProposalTxs_Error(t *testing.T) {
	invalidTxBytes := []byte{1, 2, 3}

	// Valid operations tx.
	validOperationsTx := constants.ValidEmptyMsgProposedOperationsTxBytes

	// Valid acknowledge bridges tx.
	validAcknowledgeBridgesTx := constants.MsgAcknowledgeBridges_Ids0_1_Height0_TxBytes

	// Valid add funding tx.
	validAddFundingTx := constants.ValidMsgAddPremiumVotesTxBytes

	// Valid update price tx.
	validUpdatePriceTx := constants.ValidMsgUpdateMarketPricesTxBytes

	// Valid "other" tx.
	validSendTx := constants.Msg_Send_TxBytes

	tests := map[string]struct {
		txsBytes    [][]byte
		expectedErr error
	}{
		"Less than min num txs": {
			txsBytes: [][]byte{validOperationsTx, validAddFundingTx, validUpdatePriceTx}, // need at least 4.
			expectedErr: errorsmod.Wrapf(
				process.ErrUnexpectedNumMsgs,
				"Expected the proposal to contain at least 4 txs, but got 3",
			),
		},
		"Order tx decoding fails": {
			txsBytes: [][]byte{invalidTxBytes, validAcknowledgeBridgesTx, validAddFundingTx, validUpdatePriceTx},
			expectedErr: errorsmod.Wrapf(
				process.ErrDecodingTxBytes,
				"invalid field number: tx parse error",
			),
		},
		"Acknowledge bridges tx decoding fails": {
			txsBytes: [][]byte{validOperationsTx, invalidTxBytes, validAddFundingTx, validUpdatePriceTx},
			expectedErr: errorsmod.Wrapf(
				process.ErrDecodingTxBytes,
				"invalid field number: tx parse error",
			),
		},
		"Add funding tx decoding fails": {
			txsBytes: [][]byte{validOperationsTx, validAcknowledgeBridgesTx, invalidTxBytes, validUpdatePriceTx},
			expectedErr: errorsmod.Wrapf(
				process.ErrDecodingTxBytes,
				"invalid field number: tx parse error",
			),
		},
		"Update prices tx decoding fails": {
			txsBytes: [][]byte{validOperationsTx, validAcknowledgeBridgesTx, validAddFundingTx, invalidTxBytes},
			expectedErr: errorsmod.Wrapf(
				process.ErrDecodingTxBytes,
				"invalid field number: tx parse error",
			),
		},
		"Other txs fails: invalid bytes": {
			txsBytes: [][]byte{
				validOperationsTx,
				validSendTx,    // other tx: valid.
				invalidTxBytes, // other tx: invalid.
				validAcknowledgeBridgesTx,
				validAddFundingTx,
				validUpdatePriceTx,
			},
			expectedErr: errorsmod.Wrapf(
				process.ErrDecodingTxBytes,
				"invalid field number: tx parse error",
			),
		},
		"Other txs fails: app-injected msg": {
			txsBytes: [][]byte{
				validOperationsTx,
				validSendTx,        // other tx: valid.
				validUpdatePriceTx, // other tx: invalid due to app-injected msg.
				validAcknowledgeBridgesTx,
				validAddFundingTx,
				validUpdatePriceTx,
			},
			expectedErr: errorsmod.Wrapf(
				process.ErrUnexpectedMsgType,
				"Invalid msg type or content in OtherTxs *types.MsgUpdateMarketPrices",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			_, bridgeKeeper, _, _, _, _, _ := keepertest.BridgeKeepers(t)
			ctx, pricesKeeper, _, _, _, _, _ := keepertest.PricesKeepers(t)

			// Run.
			_, err := process.DecodeProcessProposalTxs(
				ctx,
				constants.TestEncodingCfg.TxConfig.TxDecoder(),
				&abci.RequestProcessProposal{Txs: tc.txsBytes},
				bridgeKeeper,
				process.NewDefaultUpdateMarketPriceTxDecoder(pricesKeeper, constants.TestEncodingCfg.TxConfig.TxDecoder()),
			)

			// Validate.
			require.ErrorContains(t, err, tc.expectedErr.Error())
		})
	}
}

func TestDecodeProcessProposalTxs_Valid(t *testing.T) {
	// Valid order tx.
	validOperationsTx := constants.ValidEmptyMsgProposedOperationsTxBytes

	// Valid acknowledge bridges tx.
	validAcknowledgeBridgesTx := constants.MsgAcknowledgeBridges_Ids0_1_Height0_TxBytes

	// Valid add funding tx.
	validAddFundingTx := constants.ValidMsgAddPremiumVotesTxBytes

	// Valid update price tx.
	validUpdatePriceTx := constants.ValidMsgUpdateMarketPricesTxBytes

	// Valid "other" tx.
	validSingleMsgOtherTx := constants.Msg_Send_TxBytes

	// Valid "other" multi msgs tx.
	validMultiMsgOtherTx := constants.Msg_SendAndTransfer_TxBytes

	tests := map[string]struct {
		txsBytes [][]byte

		expectedOtherTxsNum    int
		expectedOtherTxOneMsgs []sdk.Msg
		expectedOtherTxTwoMsgs []sdk.Msg
	}{
		"Valid: no other tx": {
			txsBytes: [][]byte{
				validOperationsTx,
				validAcknowledgeBridgesTx,
				validAddFundingTx,
				validUpdatePriceTx,
			},
		},
		"Valid: single other tx": {
			txsBytes: [][]byte{
				validOperationsTx,
				validSingleMsgOtherTx,
				validAcknowledgeBridgesTx,
				validAddFundingTx,
				validUpdatePriceTx,
			},
			expectedOtherTxsNum:    1,
			expectedOtherTxOneMsgs: []sdk.Msg{constants.Msg_Send},
		},
		"Valid: mult other txs": {
			txsBytes: [][]byte{
				validOperationsTx,
				validSingleMsgOtherTx,
				validMultiMsgOtherTx,
				validAcknowledgeBridgesTx,
				validAddFundingTx,
				validUpdatePriceTx,
			},
			expectedOtherTxsNum:    2,
			expectedOtherTxOneMsgs: []sdk.Msg{constants.Msg_Send},
			expectedOtherTxTwoMsgs: []sdk.Msg{constants.Msg_Send, constants.Msg_Transfer},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, pricesKeeper, _, _, _, _, _ := keepertest.PricesKeepers(t)
			_, bridgeKeeper, _, _, _, _, _ := keepertest.BridgeKeepers(t)

			// Run.
			ppt, err := process.DecodeProcessProposalTxs(
				ctx,
				constants.TestEncodingCfg.TxConfig.TxDecoder(),
				&abci.RequestProcessProposal{Txs: tc.txsBytes},
				bridgeKeeper,
				process.NewDefaultUpdateMarketPriceTxDecoder(pricesKeeper, constants.TestEncodingCfg.TxConfig.TxDecoder()),
			)

			// Validate.
			require.NoError(t, err)
			require.NotNil(t, ppt)

			require.Equal(t, constants.ValidEmptyMsgProposedOperations, ppt.ProposedOperationsTx.GetMsg())
			require.Equal(
				t,
				constants.MsgAcknowledgeBridges_Ids0_1_Height0,
				ppt.AcknowledgeBridgesTx.GetMsg(),
			)
			require.Equal(t, constants.ValidMsgAddPremiumVotes, ppt.AddPremiumVotesTx.GetMsg())
			require.Equal(t, constants.ValidMsgUpdateMarketPrices, ppt.UpdateMarketPricesTx.GetMsg())

			require.Len(t, ppt.OtherTxs, tc.expectedOtherTxsNum)

			if tc.expectedOtherTxTwoMsgs != nil {
				require.Len(t, ppt.OtherTxs, 2)
				require.ElementsMatch(t, tc.expectedOtherTxOneMsgs, ppt.OtherTxs[0].GetMsgs())
				require.ElementsMatch(t, tc.expectedOtherTxTwoMsgs, ppt.OtherTxs[1].GetMsgs())
			} else if tc.expectedOtherTxOneMsgs != nil {
				require.Len(t, ppt.OtherTxs, 1)
				require.ElementsMatch(t, tc.expectedOtherTxOneMsgs, ppt.OtherTxs[0].GetMsgs())
			}
		})
	}
}

func TestProcessProposalTxs_Validate_Error(t *testing.T) {
	encodingCfg := encoding.GetTestEncodingCfg()
	txBuilder := encodingCfg.TxConfig.NewTxBuilder()

	// Operations tx.
	validOperationsTx := constants.ValidEmptyMsgProposedOperationsTxBytes

	// Acknowledge bridges tx.
	validAcknowledgeBridgesTx := constants.MsgAcknowledgeBridges_Ids0_1_Height0_TxBytes
	validAcknowledgeBridgesMsg := constants.MsgAcknowledgeBridges_Ids0_1_Height0
	invalidAcknowledgeBridgesTx := constants.MsgAcknowledgeBridges_Id55_Height15_TxBytes

	// Add funding tx.
	validAddFundingTx := constants.ValidMsgAddPremiumVotesTxBytes
	invalidAddFundingTx := constants.InvalidMsgAddPremiumVotesTxBytes

	// Update price tx.
	validUpdatePriceTx := constants.ValidMsgUpdateMarketPricesTxBytes
	invalidUpdatePriceTx := constants.InvalidMsgUpdateMarketPricesStatelessTxBytes

	// "Other" tx.
	validSingleMsgOtherTx := constants.Msg_Send_TxBytes
	invalidSingleMsgOtherTx := constants.Msg_Transfer_Invalid_SameSenderAndRecipient_TxBytes
	_ = txBuilder.SetMsgs(constants.Msg_Send, constants.Msg_Transfer_Invalid_SameSenderAndRecipient)
	invalidMultiMsgOtherTx, _ := encodingCfg.TxConfig.TxEncoder()(txBuilder.GetTx())

	tests := map[string]struct {
		txsBytes         [][]byte
		bridgingDisabled bool
		expectedErr      error
	}{
		"AcknowledgeBridges tx validation fails as event ID is not expected": {
			txsBytes: [][]byte{
				validOperationsTx,
				invalidAcknowledgeBridgesTx,
				validAddFundingTx,
				validUpdatePriceTx,
			},
			expectedErr: bridgetypes.ErrBridgeIdNotNextToAcknowledge,
		},
		"AcknowledgeBridges tx validation fails as events are non-empty and bridging is disabled": {
			txsBytes: [][]byte{
				validOperationsTx,
				validAcknowledgeBridgesTx,
				validAddFundingTx,
				validUpdatePriceTx,
			},
			bridgingDisabled: true,
			expectedErr:      bridgetypes.ErrBridgingDisabled,
		},
		"AddFunding tx validation fails": {
			txsBytes: [][]byte{validOperationsTx, validAcknowledgeBridgesTx, invalidAddFundingTx, validUpdatePriceTx},
			expectedErr: errorsmod.Wrap(
				process.ErrMsgValidateBasic,
				"premium votes must be sorted by perpetual id in ascending order and "+
					"cannot contain duplicates: MsgAddPremiumVotes is invalid"),
		},
		"UpdatePrices tx validation fails": {
			txsBytes: [][]byte{validOperationsTx, validAcknowledgeBridgesTx, validAddFundingTx, invalidUpdatePriceTx},
			expectedErr: errorsmod.Wrap(
				process.ErrMsgValidateBasic,
				"price cannot be 0 for market id (0): Market price update is invalid: stateless.",
			),
		},
		"Other txs validation fails: single tx": {
			txsBytes: [][]byte{
				validOperationsTx,
				validSingleMsgOtherTx,
				invalidSingleMsgOtherTx,
				validAcknowledgeBridgesTx,
				validAddFundingTx,
				validUpdatePriceTx,
			},
			expectedErr: errorsmod.Wrap(process.ErrMsgValidateBasic, "Sender is the same as recipient"),
		},
		"Other txs validation fails: multi txs": {
			txsBytes: [][]byte{
				validOperationsTx,
				validSingleMsgOtherTx,
				invalidMultiMsgOtherTx,
				validAcknowledgeBridgesTx,
				validAddFundingTx,
				validUpdatePriceTx,
			},
			expectedErr: errorsmod.Wrap(process.ErrMsgValidateBasic, "Sender is the same as recipient"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, pricesKeeper, _, indexPriceCache, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)
			indexPriceCache.UpdatePrices(constants.AtTimeTSingleExchangePriceUpdate)

			mockBridgeKeeper := &mocks.ProcessBridgeKeeper{}
			mockBridgeKeeper.On("GetSafetyParams", mock.Anything).Return(
				bridgetypes.SafetyParams{
					IsDisabled:  tc.bridgingDisabled,
					DelayBlocks: 5, // dummy value, not considered by Validate.
				},
			)
			mockBridgeKeeper.On("GetAcknowledgedEventInfo", mock.Anything).Return(
				constants.AcknowledgedEventInfo_Id0_Height0,
			)
			mockBridgeKeeper.On("GetRecognizedEventInfo", mock.Anything).Return(
				constants.RecognizedEventInfo_Id2_Height0,
			)
			for _, bridgeEvent := range validAcknowledgeBridgesMsg.Events {
				mockBridgeKeeper.On("GetBridgeEventFromServer", mock.Anything, bridgeEvent.Id).Return(bridgeEvent, true).Once()
			}

			ppt, err := process.DecodeProcessProposalTxs(
				ctx,
				encodingCfg.TxConfig.TxDecoder(),
				&abci.RequestProcessProposal{Txs: tc.txsBytes},
				mockBridgeKeeper,
				process.NewDefaultUpdateMarketPriceTxDecoder(pricesKeeper, constants.TestEncodingCfg.TxConfig.TxDecoder()),
			)
			require.NoError(t, err)

			// Run.
			err = ppt.Validate()

			// Validate.
			require.ErrorContains(t, err, tc.expectedErr.Error())
		})
	}
}

func TestProcessProposalTxs_Validate_Valid(t *testing.T) {
	// Valid order tx.
	validOperationsTx := constants.ValidEmptyMsgProposedOperationsTxBytes

	// Valid acknowledge bridges tx.
	validAcknowledgeBridgesTx := constants.MsgAcknowledgeBridges_Ids0_1_Height0_TxBytes
	validAcknowledgeBridgesMsg := constants.MsgAcknowledgeBridges_Ids0_1_Height0
	emptyAcknowledgeBridgesTx := constants.MsgAcknowledgeBridges_NoEvents_TxBytes

	// Valid add funding tx.
	validAddFundingTx := constants.ValidMsgAddPremiumVotesTxBytes

	// Valid update price tx.
	validUpdatePriceTx := constants.ValidMsgUpdateMarketPricesTxBytes

	// Valid "other" tx.
	validSingleMsgOtherTx := constants.Msg_Send_TxBytes

	// Valid "other" multi msgs tx.
	validMultiMsgOtherTx := constants.Msg_SendAndTransfer_TxBytes

	tests := map[string]struct {
		txsBytes         [][]byte
		bridgingDisabled bool
	}{
		"No other txs": {
			txsBytes: [][]byte{
				validOperationsTx,
				validAcknowledgeBridgesTx,
				validAddFundingTx,
				validUpdatePriceTx,
			},
		},
		"Single other txs": {
			txsBytes: [][]byte{
				validOperationsTx,
				validSingleMsgOtherTx,
				validAcknowledgeBridgesTx,
				validAddFundingTx,
				validUpdatePriceTx,
			},
		},
		"Multi other txs": {
			txsBytes: [][]byte{
				validOperationsTx,
				validSingleMsgOtherTx,
				validMultiMsgOtherTx,
				validAcknowledgeBridgesTx,
				validAddFundingTx,
				validUpdatePriceTx,
			},
		},
		"Empty bridge events and bridging is disabled": {
			txsBytes: [][]byte{
				validOperationsTx,
				validSingleMsgOtherTx,
				emptyAcknowledgeBridgesTx,
				validAddFundingTx,
				validUpdatePriceTx,
			},
			bridgingDisabled: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, pricesKeeper, _, indexPriceCache, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)
			indexPriceCache.UpdatePrices(constants.AtTimeTSingleExchangePriceUpdate)

			mockBridgeKeeper := &mocks.ProcessBridgeKeeper{}
			mockBridgeKeeper.On("GetSafetyParams", mock.Anything).Return(
				bridgetypes.SafetyParams{
					IsDisabled:  tc.bridgingDisabled,
					DelayBlocks: 5, // dummy value, not considered by Validate.
				},
			)
			mockBridgeKeeper.On("GetAcknowledgedEventInfo", mock.Anything).Return(
				constants.AcknowledgedEventInfo_Id0_Height0,
			)
			mockBridgeKeeper.On("GetRecognizedEventInfo", mock.Anything).Return(
				constants.RecognizedEventInfo_Id2_Height0,
			)
			for _, bridgeEvent := range validAcknowledgeBridgesMsg.Events {
				mockBridgeKeeper.On("GetBridgeEventFromServer", mock.Anything, bridgeEvent.Id).Return(bridgeEvent, true).Once()
			}

			ppt, err := process.DecodeProcessProposalTxs(
				ctx,
				constants.TestEncodingCfg.TxConfig.TxDecoder(),
				&abci.RequestProcessProposal{Txs: tc.txsBytes},
				mockBridgeKeeper,
				process.NewDefaultUpdateMarketPriceTxDecoder(pricesKeeper, constants.TestEncodingCfg.TxConfig.TxDecoder()),
			)
			require.NoError(t, err)

			// Run.
			err = ppt.Validate()

			// Validate.
			require.NoError(t, err)
		})
	}
}
