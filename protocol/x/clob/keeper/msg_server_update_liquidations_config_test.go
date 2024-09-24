package keeper_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/memclob"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
					InsuranceFundFeePpm: 5_000,
					ValidatorFeePpm:     200_000,
					LiquidityFeePpm:     800_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm: 0,
					},
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
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
