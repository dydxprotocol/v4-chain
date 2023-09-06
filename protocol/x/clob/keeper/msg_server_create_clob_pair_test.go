package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	perptest "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestCreateClobPair(t *testing.T) {
	testClobPair1 := *clobtest.GenerateClobPair(
		clobtest.WithId(1),
		clobtest.WithPerpetualId(1),
	)
	testPerp1 := *perptest.GeneratePerpetual(
		perptest.WithId(1),
		perptest.WithMarketId(1),
	)
	testMarket1 := *pricestest.GenerateMarketParamPrice(pricestest.WithId(1))
	testCases := map[string]struct {
		setup             func(t *testing.T, ks keepertest.ClobKeepersTestContext, manager *mocks.IndexerEventManager)
		msg               *types.MsgCreateClobPair
		expectedClobPairs []types.ClobPair
		expectedErr       string
	}{
		"Succeeds: create new clob pair": {
			setup: func(t *testing.T, ks keepertest.ClobKeepersTestContext, mockIndexerEventManager *mocks.IndexerEventManager) {
				keepertest.CreateTestPricesAndPerpetualMarkets(
					t,
					ks.Ctx,
					ks.PerpetualsKeeper,
					ks.PricesKeeper,
					[]perptypes.Perpetual{testPerp1},
					[]pricestypes.MarketParamPrice{testMarket1},
				)
				mockIndexerEventManager.On("AddTxnEvent",
					ks.Ctx,
					indexerevents.SubtypePerpetualMarket,
					indexer_manager.GetB64EncodedEventMessage(
						indexerevents.NewPerpetualMarketCreateEvent(
							testClobPair1.MustGetPerpetualId(),
							testClobPair1.GetId(),
							testPerp1.Params.Ticker,
							testPerp1.Params.MarketId,
							testClobPair1.Status,
							testClobPair1.QuantumConversionExponent,
							testPerp1.Params.AtomicResolution,
							testClobPair1.SubticksPerTick,
							testClobPair1.MinOrderBaseQuantums,
							testClobPair1.StepBaseQuantums,
							testPerp1.Params.LiquidityTier,
						),
					),
				).Return()
			},
			msg: &types.MsgCreateClobPair{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				ClobPair:  testClobPair1,
			},
			expectedClobPairs: []types.ClobPair{testClobPair1},
		},
		"Succeeds: clob pair already exists": {
			setup: func(t *testing.T, ks keepertest.ClobKeepersTestContext, mockIndexerEventManager *mocks.IndexerEventManager) {
				keepertest.CreateTestPricesAndPerpetualMarkets(
					t,
					ks.Ctx,
					ks.PerpetualsKeeper,
					ks.PricesKeeper,
					[]perptypes.Perpetual{testPerp1},
					[]pricestypes.MarketParamPrice{testMarket1},
				)
				mockIndexerEventManager.On("AddTxnEvent",
					ks.Ctx,
					indexerevents.SubtypePerpetualMarket,
					indexer_manager.GetB64EncodedEventMessage(
						indexerevents.NewPerpetualMarketCreateEvent(
							testClobPair1.MustGetPerpetualId(),
							testClobPair1.GetId(),
							testPerp1.Params.Ticker,
							testPerp1.Params.MarketId,
							testClobPair1.Status,
							testClobPair1.QuantumConversionExponent,
							testPerp1.Params.AtomicResolution,
							testClobPair1.SubticksPerTick,
							testClobPair1.MinOrderBaseQuantums,
							testClobPair1.StepBaseQuantums,
							testPerp1.Params.LiquidityTier,
						),
					),
				).Return()
				_, err := ks.ClobKeeper.CreatePerpetualClobPair(
					ks.Ctx,
					testClobPair1.Id,
					testClobPair1.MustGetPerpetualId(),
					satypes.BaseQuantums(testClobPair1.MinOrderBaseQuantums),
					satypes.BaseQuantums(testClobPair1.StepBaseQuantums),
					testClobPair1.QuantumConversionExponent,
					testClobPair1.SubticksPerTick,
					testClobPair1.Status,
				)
				require.NoError(t, err)
			},
			msg: &types.MsgCreateClobPair{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				ClobPair:  testClobPair1,
			},
			expectedClobPairs: []types.ClobPair{testClobPair1},
			expectedErr:       "ClobPair with id already exists",
		},
		"Failure: refers to non-existing perpetual": {
			setup: func(t *testing.T, ks keepertest.ClobKeepersTestContext, mockIndexerEventManager *mocks.IndexerEventManager) {
			},
			msg: &types.MsgCreateClobPair{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				ClobPair:  testClobPair1,
			},
			expectedClobPairs: nil,
			expectedErr:       "has invalid perpetual.: 1: Perpetual does not exist",
		},
		"Failure: invalid authority": {
			setup: func(t *testing.T, ks keepertest.ClobKeepersTestContext, mockIndexerEventManager *mocks.IndexerEventManager) {
			},
			msg: &types.MsgCreateClobPair{
				Authority: "invalid",
				ClobPair:  testClobPair1,
			},
			expectedClobPairs: nil,
			expectedErr:       "invalid authority invalid: expected gov account as only signer for proposal message",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)
			tc.setup(t, ks, mockIndexerEventManager)

			msgServer := keeper.NewMsgServerImpl(ks.ClobKeeper)
			wrappedCtx := sdk.WrapSDKContext(ks.Ctx)

			_, err := msgServer.CreateClobPair(wrappedCtx, tc.msg)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedClobPairs, ks.ClobKeeper.GetAllClobPairs(ks.Ctx))
		})
	}
}
