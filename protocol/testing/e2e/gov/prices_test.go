package gov_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
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
		Id:                 0,
		Pair:               "btc-adv4tnt",
		Exponent:           -8,
		MinExchanges:       2,
		MinPriceChangePpm:  1_000,
		ExchangeConfigJson: "{}",
	}

	MODIFIED_MARKET_PARAM = pricestypes.MarketParam{
		Id:                 GENESIS_MARKET_PARAM.Id,
		Pair:               "eth-adv4tnt",
		Exponent:           GENESIS_MARKET_PARAM.Exponent, // exponent cannot be updated
		MinExchanges:       3,
		MinPriceChangePpm:  2_002,
		ExchangeConfigJson: `{"exchanges":[{"exchangeName":"Bitfinex","ticker":"tBTCUSD"}]}`,
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
					Id:                 MODIFIED_MARKET_PARAM.Id + 1, // id does not exist
					Pair:               MODIFIED_MARKET_PARAM.Pair,
					Exponent:           MODIFIED_MARKET_PARAM.Exponent,
					MinExchanges:       MODIFIED_MARKET_PARAM.MinExchanges,
					MinPriceChangePpm:  MODIFIED_MARKET_PARAM.MinPriceChangePpm,
					ExchangeConfigJson: MODIFIED_MARKET_PARAM.ExchangeConfigJson,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
		},
		"Failure: exponent is updated": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: lib.GovModuleAddress.String(),
				MarketParam: pricestypes.MarketParam{
					Id:                 MODIFIED_MARKET_PARAM.Id,
					Pair:               MODIFIED_MARKET_PARAM.Pair,
					Exponent:           MODIFIED_MARKET_PARAM.Exponent + 1, // update to exponent is not permitted.
					MinExchanges:       MODIFIED_MARKET_PARAM.MinExchanges,
					MinPriceChangePpm:  MODIFIED_MARKET_PARAM.MinPriceChangePpm,
					ExchangeConfigJson: MODIFIED_MARKET_PARAM.ExchangeConfigJson,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
		},
		"Failure: empty pair": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: lib.GovModuleAddress.String(),
				MarketParam: pricestypes.MarketParam{
					Id:                 MODIFIED_MARKET_PARAM.Id,
					Pair:               "", // invalid
					Exponent:           MODIFIED_MARKET_PARAM.Exponent,
					MinExchanges:       MODIFIED_MARKET_PARAM.MinExchanges,
					MinPriceChangePpm:  MODIFIED_MARKET_PARAM.MinPriceChangePpm,
					ExchangeConfigJson: MODIFIED_MARKET_PARAM.ExchangeConfigJson,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: min exchanges is 0": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: lib.GovModuleAddress.String(),
				MarketParam: pricestypes.MarketParam{
					Id:                 MODIFIED_MARKET_PARAM.Id,
					Pair:               MODIFIED_MARKET_PARAM.Pair,
					Exponent:           MODIFIED_MARKET_PARAM.Exponent,
					MinExchanges:       0, // invalid
					MinPriceChangePpm:  MODIFIED_MARKET_PARAM.MinPriceChangePpm,
					ExchangeConfigJson: MODIFIED_MARKET_PARAM.ExchangeConfigJson,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: malformed exchange config json": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: lib.GovModuleAddress.String(),
				MarketParam: pricestypes.MarketParam{
					Id:                 MODIFIED_MARKET_PARAM.Id,
					Pair:               MODIFIED_MARKET_PARAM.Pair,
					Exponent:           MODIFIED_MARKET_PARAM.Exponent,
					MinExchanges:       MODIFIED_MARKET_PARAM.MinExchanges,
					MinPriceChangePpm:  MODIFIED_MARKET_PARAM.MinPriceChangePpm,
					ExchangeConfigJson: `{{"exchanges":[{"exchangeName":"Bitfinex","ticker":"tBTCUSD"}]}`, // invalid
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
				// Initialize prices module with genesis market param.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *pricestypes.GenesisState) {
						marketParamPrice := pricestest.GenerateMarketParamPrice(
							pricestest.WithId(GENESIS_MARKET_PARAM.Id),
							pricestest.WithPair(GENESIS_MARKET_PARAM.Pair),
							pricestest.WithExponent(GENESIS_MARKET_PARAM.Exponent),
							pricestest.WithMinExchanges(GENESIS_MARKET_PARAM.MinExchanges),
							pricestest.WithMinPriceChangePpm(GENESIS_MARKET_PARAM.MinPriceChangePpm),
							pricestest.WithExchangeConfigJson(GENESIS_MARKET_PARAM.ExchangeConfigJson),
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
