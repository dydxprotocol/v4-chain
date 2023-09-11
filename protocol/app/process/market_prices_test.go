package process_test

import (
	errorsmod "cosmossdk.io/errors"
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/process"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestDecodeUpdateMarketPricesTx(t *testing.T) {
	encodingCfg := encoding.GetTestEncodingCfg()
	txBuilder := encodingCfg.TxConfig.NewTxBuilder()

	// Valid.
	validMsgTxBytes := constants.ValidMsgUpdateMarketPricesTxBytes

	// Duplicate.
	_ = txBuilder.SetMsgs(constants.ValidMsgUpdateMarketPrices, constants.ValidMsgUpdateMarketPrices)
	duplicateMsgTxBytes, _ := encodingCfg.TxConfig.TxEncoder()(txBuilder.GetTx())

	// Incorrect type.
	incorrectMsgTxBytes := constants.ValidMsgAddPremiumVotesTxBytes

	tests := map[string]struct {
		txBytes []byte

		expectedErr error
		expectedMsg *types.MsgUpdateMarketPrices
	}{
		"Error: decode fails": {
			txBytes:     []byte{1, 2, 3}, // invalid bytes.
			expectedErr: errors.New("tx parse error: Decoding tx bytes failed"),
		},
		"Error: empty bytes": {
			txBytes: []byte{}, // empty returns 0 msgs.
			expectedErr: errors.New("Msg Type: types.MsgUpdateMarketPrices, " +
				"Expected 1 num of msgs, but got 0: Unexpected num of msgs"),
		},
		"Error: incorrect msg len": {
			txBytes: duplicateMsgTxBytes,
			expectedErr: errors.New("Msg Type: types.MsgUpdateMarketPrices, " +
				"Expected 1 num of msgs, but got 2: Unexpected num of msgs"),
		},
		"Error: incorrect msg type": {
			txBytes: incorrectMsgTxBytes,
			expectedErr: errors.New(
				"Expected MsgType types.MsgUpdateMarketPrices, but " +
					"got *types.MsgAddPremiumVotes: Unexpected msg type",
			),
		},
		"Valid": {
			txBytes:     validMsgTxBytes,
			expectedMsg: constants.ValidMsgUpdateMarketPrices,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, k, _, _, _, _ := keepertest.PricesKeepers(t)
			umpt, err := process.DecodeUpdateMarketPricesTx(ctx, k, encodingCfg.TxConfig.TxDecoder(), tc.txBytes)
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
				require.Nil(t, umpt)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedMsg, umpt.GetMsg())
			}
		})
	}
}

func TestUpdateMarketPricesTx_Validate(t *testing.T) {
	// Valid.
	validMsgTxBytes := constants.ValidMsgUpdateMarketPricesTxBytes

	// Invalid (stateless).
	invalidStatelessMsgTxBytes := constants.InvalidMsgUpdateMarketPricesStatelessTxBytes

	// Invalid (stateful + deterministic).
	invalidStatefulMsgTxBytes := constants.InvalidMsgUpdateMarketPricesStatefulTxBytes

	tests := map[string]struct {
		txBytes     []byte
		indexPrices []*api.MarketPriceUpdate

		expectedErr error
	}{
		"Error: Stateful + Deterministic validation fails": {
			txBytes:     invalidStatefulMsgTxBytes,
			indexPrices: constants.AtTimeTSingleExchangePriceUpdate,
			expectedErr: errorsmod.Wrap(
				types.ErrInvalidMarketPriceUpdateDeterministic,
				"market param price (99) does not exist",
			),
		},
		"Error: Stateful + NonDeterministic validation fails": {
			txBytes: validMsgTxBytes, // Msg is valid, but there's no corresponding index price.
			// Skip index price updates, so the validation fails.
			expectedErr: errorsmod.Wrapf(types.ErrIndexPriceNotAvailable, "index price for market (0) is not available"),
		},
		"Error: ValidateBasic fails": {
			txBytes:     invalidStatelessMsgTxBytes,
			indexPrices: constants.AtTimeTSingleExchangePriceUpdate,
			expectedErr: errorsmod.Wrap(
				process.ErrMsgValidateBasic,
				"price cannot be 0 for market id (0): Market price update is invalid: stateless.",
			),
		},
		"Valid: ValidateBasic passes": {
			txBytes:     validMsgTxBytes,
			indexPrices: constants.AtTimeTSingleExchangePriceUpdate,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, k, _, indexPriceCache, _, mockTimeProvider := keepertest.PricesKeepers(t)
			keepertest.CreateTestMarkets(t, ctx, k)
			indexPriceCache.UpdatePrices(tc.indexPrices)
			mockTimeProvider.On("Now").Return(constants.TimeT)
			umpt, err := process.DecodeUpdateMarketPricesTx(ctx, k, constants.TestEncodingCfg.TxConfig.TxDecoder(), tc.txBytes)
			require.NoError(t, err)

			// Run and Validate.
			err = umpt.Validate()
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdateMarketPricesTx_GetMsg(t *testing.T) {
	validMsgTxBytes := constants.ValidMsgUpdateMarketPricesTxBytes

	tests := map[string]struct {
		txWrapper   process.UpdateMarketPricesTx
		txBytes     []byte
		expectedMsg *types.MsgUpdateMarketPrices
	}{
		"Returns nil msg": {
			txWrapper: process.UpdateMarketPricesTx{},
		},
		"Returns valid msg": {
			txBytes:     validMsgTxBytes,
			expectedMsg: constants.ValidMsgUpdateMarketPrices,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var msg sdk.Msg
			if tc.txBytes != nil {
				ctx, k, _, _, _, _ := keepertest.PricesKeepers(t)
				umpt, err := process.DecodeUpdateMarketPricesTx(ctx, k, constants.TestEncodingCfg.TxConfig.TxDecoder(), tc.txBytes)
				require.NoError(t, err)
				msg = umpt.GetMsg()
			} else {
				msg = tc.txWrapper.GetMsg()
				require.Equal(t, tc.expectedMsg, msg)
			}
			require.Equal(t, tc.expectedMsg, msg)
		})
	}
}
