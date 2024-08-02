package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices"
)

func TestMevNodeToNodeCalculation(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	ks.MarketMapKeeper.InitGenesis(ks.Ctx, constants.MarketMap_DefaultGenesisState)
	prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
	perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)
	for testName, tc := range map[string]struct {
		request  *types.MevNodeToNodeCalculationRequest
		response *types.MevNodeToNodeCalculationResponse
		err      error
	}{
		"Nil request returns an error": {
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
		"Nil validator MEV metrics returns an error": {
			request: &types.MevNodeToNodeCalculationRequest{},
			err:     status.Error(codes.InvalidArgument, "missing validator MEV metrics"),
		},
		"Can successfully run validator MEV calculation on zero validator and block proposer matches": {
			request: &types.MevNodeToNodeCalculationRequest{
				ValidatorMevMetrics: &types.MevNodeToNodeMetrics{
					ValidatorMevMatches: &types.ValidatorMevMatches{},
					ClobMidPrices: []types.ClobMidPrice{
						{
							ClobPair: constants.ClobPair_Btc,
							Subticks: 50_000_000_000, // $50,000 / BTC
						},
						{
							ClobPair: constants.ClobPair_Eth,
							Subticks: 3_000_000_000, // $3000 / ETH
						},
					},
				},
				BlockProposerMatches: &types.ValidatorMevMatches{},
			},
			response: &types.MevNodeToNodeCalculationResponse{
				Results: []types.MevNodeToNodeCalculationResponse_MevAndVolumePerClob{
					{
						ClobPairId: constants.ClobPair_Btc.Id,
						Mev:        0,
						Volume:     0,
					},
					{
						ClobPairId: constants.ClobPair_Eth.Id,
						Mev:        0,
						Volume:     0,
					},
				},
			},
		},
		"Can successfully run validator MEV calculation on validator and block proposer matches": {
			request: &types.MevNodeToNodeCalculationRequest{
				ValidatorMevMetrics: &types.MevNodeToNodeMetrics{
					ValidatorMevMatches: &types.ValidatorMevMatches{
						Matches: []types.MEVMatch{
							{
								TakerOrderSubaccountId: &constants.Alice_Num0,
								TakerFeePpm:            0,

								MakerOrderSubaccountId: &constants.Bob_Num0,
								MakerOrderSubticks:     49_000_000_000, // $49,000 / BTC
								MakerOrderIsBuy:        true,
								MakerFeePpm:            0,

								ClobPairId: 0,
								FillAmount: 100_000_000, // 1 BTC
							},
						},
					},
					ClobMidPrices: []types.ClobMidPrice{
						{
							ClobPair: constants.ClobPair_Btc,
							Subticks: 50_000_000_000, // $50,000 / BTC
						},
					},
				},
				BlockProposerMatches: &types.ValidatorMevMatches{
					Matches: []types.MEVMatch{},
					LiquidationMatches: []types.MEVLiquidationMatch{
						{
							LiquidatedSubaccountId:          constants.Alice_Num0,
							InsuranceFundDeltaQuoteQuantums: 0, // $0 paid to insurance fund

							MakerOrderSubaccountId: constants.Carl_Num0,
							MakerOrderSubticks:     48_000_000_000, // $48,000 / BTC
							MakerOrderIsBuy:        true,
							MakerFeePpm:            0,

							ClobPairId: 0,
							FillAmount: 100_000_000, // 1 BTC
						},
					},
				},
			},
			response: &types.MevNodeToNodeCalculationResponse{
				Results: []types.MevNodeToNodeCalculationResponse_MevAndVolumePerClob{
					{
						ClobPairId: constants.ClobPair_Btc.Id,
						Mev:        2_000_000_000,  // $2,000 of MEV
						Volume:     49_000_000_000, // $49,000 of volume
					},
				},
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			response, err := ks.ClobKeeper.MevNodeToNodeCalculation(ks.Ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.ElementsMatch(t, tc.response.Results, response.Results)
			}
		})
	}
}
