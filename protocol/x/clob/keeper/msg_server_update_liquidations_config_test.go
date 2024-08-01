package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices"
	"github.com/stretchr/testify/require"
)

func TestUpdateLiquidationsConfig(t *testing.T) {
	testCases := map[string]struct {
		msg           *types.MsgUpdateLiquidationsConfig
		expectedError error
	}{
		"Succeeds": {
			msg: &types.MsgUpdateLiquidationsConfig{
				Authority:          lib.GovModuleAddress.String(),
				LiquidationsConfig: constants.LiquidationsConfig_No_Limit,
			},
		},
		"Error: invalid liquidations config": {
			msg: &types.MsgUpdateLiquidationsConfig{
				Authority: lib.GovModuleAddress.String(),
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 5_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm: 0,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
				},
			},
			expectedError: types.ErrInvalidLiquidationsConfig,
		},
		"Error: invalid authority": {
			msg: &types.MsgUpdateLiquidationsConfig{
				Authority:          "foobar",
				LiquidationsConfig: types.LiquidationsConfig{},
			},
			expectedError: govtypes.ErrInvalidSigner,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			ks.MarketMapKeeper.InitGenesis(ks.Ctx, constants.MarketMap_DefaultGenesisState)
			prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			msgServer := keeper.NewMsgServerImpl(ks.ClobKeeper)
			_, err := msgServer.UpdateLiquidationsConfig(ks.Ctx, tc.msg)

			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, ks.ClobKeeper.GetLiquidationsConfig(ks.Ctx), tc.msg.LiquidationsConfig)
			}
		})
	}
}
