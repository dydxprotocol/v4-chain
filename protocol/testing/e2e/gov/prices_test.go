package gov_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	marketmaptypes "github.com/dydxprotocol/slinky/x/marketmap/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

var (
	GENESIS_MARKET_PARAM = pricestypes.MarketParam{
		Id:                0,
		Pair:              "btc-adv4tnt",
		MinPriceChangePpm: 1_000,
	}

	MODIFIED_MARKET_PARAM = pricestypes.MarketParam{
		Id:                GENESIS_MARKET_PARAM.Id,
		Pair:              GENESIS_MARKET_PARAM.Pair,
		MinPriceChangePpm: 2_002,
	}
)

// This tests `MsgUpdateMarketParam` in `x/prices`.
func TestUpdateMarketParam(t *testing.T) {
	tests := map[string]struct {
		msg                       *pricestypes.MsgUpdateMarketParam
		expectCheckTxFails        bool
		expectSubmitProposalFails bool
		expectedProposalStatus    govtypesv1.ProposalStatus
	}{
		"Success": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority:   lib.GovModuleAddress.String(),
				MarketParam: MODIFIED_MARKET_PARAM,
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Failure: market param does not exist": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: lib.GovModuleAddress.String(),
				MarketParam: pricestypes.MarketParam{
					Id:                MODIFIED_MARKET_PARAM.Id + 1, // id does not exist
					Pair:              MODIFIED_MARKET_PARAM.Pair,
					MinPriceChangePpm: MODIFIED_MARKET_PARAM.MinPriceChangePpm,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
		},
		"Failure: new pair name does not exist in marketmap": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: lib.GovModuleAddress.String(),
				MarketParam: pricestypes.MarketParam{
					Id:                MODIFIED_MARKET_PARAM.Id,
					Pair:              "nonexistent-pair",
					MinPriceChangePpm: MODIFIED_MARKET_PARAM.MinPriceChangePpm,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
		},
		"Failure: empty pair": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: lib.GovModuleAddress.String(),
				MarketParam: pricestypes.MarketParam{
					Id:                MODIFIED_MARKET_PARAM.Id,
					Pair:              "", // invalid
					MinPriceChangePpm: MODIFIED_MARKET_PARAM.MinPriceChangePpm,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: invalid authority": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority:   constants.AliceAccAddress.String(),
				MarketParam: MODIFIED_MARKET_PARAM,
			},
			expectSubmitProposalFails: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *govtypesv1.GenesisState) {
						genesisState.Params.VotingPeriod = &testapp.TestVotingPeriod
					},
				)
				// Initialize marketmap module with genesis market.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *marketmaptypes.GenesisState) {
						markets := make(map[string]marketmaptypes.Market)
						cp, _ := slinky.MarketPairToCurrencyPair(GENESIS_MARKET_PARAM.Pair)
						ticker := marketmaptypes.Ticker{
							CurrencyPair:     cp,
							Decimals:         8,
							MinProviderCount: 3,
							Enabled:          true,
							Metadata_JSON:    "",
						}
						markets[cp.String()] = marketmaptypes.Market{
							Ticker: ticker,
							ProviderConfigs: []marketmaptypes.ProviderConfig{
								{Name: "binance_ws", OffChainTicker: "test"},
								{Name: "bybit_ws", OffChainTicker: "test"},
								{Name: "coinbase_ws", OffChainTicker: "test"},
							},
						}
						marketMap := marketmaptypes.MarketMap{
							Markets: markets,
						}
						genesisState.MarketMap = marketMap
					},
				)

				// Initialize prices module with genesis market param.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *pricestypes.GenesisState) {
						marketParamPrice := pricestest.GenerateMarketParamPrice(
							pricestest.WithId(GENESIS_MARKET_PARAM.Id),
							pricestest.WithPair(GENESIS_MARKET_PARAM.Pair),
							pricestest.WithMinPriceChangePpm(GENESIS_MARKET_PARAM.MinPriceChangePpm),
						)
						genesisState.MarketParams = []pricestypes.MarketParam{marketParamPrice.Param}
						genesisState.MarketPrices = []pricestypes.MarketPrice{marketParamPrice.Price}
					},
				)
				// Initialize perpetuals module with no perpetuals.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Perpetuals = []perptypes.Perpetual{}
					},
				)
				// Initialize clob module with no clob pairs.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.ClobPairs = []clobtypes.ClobPair{}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			initialMarketParams := tApp.App.PricesKeeper.GetAllMarketParams(ctx)

			// Submit and tally governance proposal that includes `MsgUpdateMarketParam`.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				[]sdk.Msg{tc.msg},
				testapp.TestSubmitProposalTxHeight,
				tc.expectCheckTxFails,
				tc.expectSubmitProposalFails,
				tc.expectedProposalStatus,
			)

			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				// If proposal is supposed to pass, verify that maret param is updated.
				updatedMarketParam, exists := tApp.App.PricesKeeper.GetMarketParam(ctx, tc.msg.MarketParam.Id)
				require.True(t, exists)
				require.Equal(t, tc.msg.MarketParam, updatedMarketParam)
			} else {
				// Otherwise, verify that market params are unchanged.
				require.Equal(t, initialMarketParams, tApp.App.PricesKeeper.GetAllMarketParams(ctx))
			}
		})
	}
}
