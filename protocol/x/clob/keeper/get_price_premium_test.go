package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	errorsmod "cosmossdk.io/errors"

	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testMemClobMethodArgs struct {
	// Input args for `MemClob.GetPricePremium`
	clobPair              types.ClobPair
	indexPrice            pricestypes.MarketPrice
	baseAtomicResolution  int32
	quoteAtomicResolution int32
	impactNotionalAmount  *big.Int
	maxAbsPremiumVotePpm  *big.Int

	// Return values from `MemClob.GetPricePremium`
	premiumPpm         int32
	getPricePremiumErr error
}

type testCase map[string]struct {
	setUpMockMemClob        func(mck *mocks.MemClob, args testMemClobMethodArgs)
	args                    testMemClobMethodArgs
	perpetualId             uint32
	perpetualIdToClobPairId map[uint32][]types.ClobPairId
	expectedErr             error
}

func TestGetPricePremiumForPerpetual(t *testing.T) {
	tests := testCase{
		"Success": {
			perpetualId: 0,
			args: testMemClobMethodArgs{
				clobPair: constants.ClobPair_Btc,
				indexPrice: pricestypes.MarketPrice{
					Price:    1_000_000_000, // $10_000
					Exponent: -5,
				},
				baseAtomicResolution:  -9,
				quoteAtomicResolution: -6,
				impactNotionalAmount:  big.NewInt(5000),
				maxAbsPremiumVotePpm:  big.NewInt(1000),
			},
			setUpMockMemClob: func(mck *mocks.MemClob, args testMemClobMethodArgs) {
				mck.On(
					"GetPricePremium",
					mock.Anything,
					args.clobPair,
					perptypes.GetPricePremiumParams{
						IndexPrice:                  args.indexPrice,
						BaseAtomicResolution:        args.baseAtomicResolution,
						QuoteAtomicResolution:       args.quoteAtomicResolution,
						ImpactNotionalQuoteQuantums: args.impactNotionalAmount,
						MaxAbsPremiumVotePpm:        args.maxAbsPremiumVotePpm,
					},
				).Return(
					args.premiumPpm,
					args.getPricePremiumErr,
				)
			},
		},
		"Failure: GetClobPairIdForPerpetual error": {
			perpetualId: 1,
			args: testMemClobMethodArgs{
				clobPair: constants.ClobPair_Btc,
			},
			setUpMockMemClob: func(mck *mocks.MemClob, args testMemClobMethodArgs) {},
			expectedErr: errors.New(
				"Perpetual ID 1 has no associated CLOB pairs: " +
					"The provided perpetual ID does not have any associated CLOB pairs",
			),
		},
		"Failure: GetPricePremium failure": {
			perpetualId: 0,
			args: testMemClobMethodArgs{
				clobPair: constants.ClobPair_Btc,
				indexPrice: pricestypes.MarketPrice{
					Price:    1_000_000_000, // $10_000
					Exponent: -5,
				},
				baseAtomicResolution:  -9,
				quoteAtomicResolution: -6,
				impactNotionalAmount:  big.NewInt(5000),
				maxAbsPremiumVotePpm:  big.NewInt(1000),
				getPricePremiumErr:    errors.New("GetPricePremium error"),
			},
			setUpMockMemClob: func(mck *mocks.MemClob, args testMemClobMethodArgs) {
				mck.On(
					"GetPricePremium",
					mock.Anything,
					args.clobPair,
					perptypes.GetPricePremiumParams{
						IndexPrice:                  args.indexPrice,
						BaseAtomicResolution:        args.baseAtomicResolution,
						QuoteAtomicResolution:       args.quoteAtomicResolution,
						ImpactNotionalQuoteQuantums: args.impactNotionalAmount,
						MaxAbsPremiumVotePpm:        args.maxAbsPremiumVotePpm,
					},
				).Return(
					args.premiumPpm,
					args.getPricePremiumErr,
				)
			},
			expectedErr: errors.New("GetPricePremium error"),
		},
		"Failure, clob pair not found": {
			perpetualId: 0,
			args: testMemClobMethodArgs{
				clobPair: constants.ClobPair_Btc, // clobPairId = 1000
			},
			setUpMockMemClob: func(mck *mocks.MemClob, args testMemClobMethodArgs) {},
			// clob pair is created with id 0, but we override in-memory datastructure with 1 to cause error
			perpetualIdToClobPairId: map[uint32][]types.ClobPairId{
				0: {1},
			},
			expectedErr: errorsmod.Wrapf(
				types.ErrInvalidClob,
				"GetPricePremiumForPerpetual: did not find clob pair with clobPairId = %d",
				1,
			),
		},
		"Success, premium is zeroed for initializing clob pair": {
			perpetualId: 0,
			args: testMemClobMethodArgs{
				clobPair: constants.ClobPair_Btc_Initializing,
			},
			setUpMockMemClob: func(mck *mocks.MemClob, args testMemClobMethodArgs) {},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything)
			memClob.On("CreateOrderbook", mock.Anything, mock.Anything, mock.Anything)

			tc.setUpMockMemClob(memClob, tc.args)

			mockIndexerEventManager := &mocks.IndexerEventManager{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, nil, mockIndexerEventManager)

			ks.MarketMapKeeper.InitGenesis(ks.Ctx, constants.MarketMap_DefaultGenesisState)
			prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			perpetualId := clobtest.MustPerpetualId(tc.args.clobPair)
			perpetual := constants.Perpetuals_DefaultGenesisState.Perpetuals[perpetualId]
			// TODO(IND-362): Refactor into helper function that takes in perpetual/clobPair args.
			// PerpetualMarketCreateEvents are emitted when initializing the genesis state, so we need to mock
			// the indexer event manager to expect these events.
			mockIndexerEventManager.On("AddTxnEvent",
				ks.Ctx,
				indexerevents.SubtypePerpetualMarket,
				indexerevents.PerpetualMarketEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewPerpetualMarketCreateEvent(
						perpetualId,
						0,
						perpetual.Params.Ticker,
						perpetual.Params.MarketId,
						tc.args.clobPair.Status,
						tc.args.clobPair.QuantumConversionExponent,
						perpetual.Params.AtomicResolution,
						tc.args.clobPair.SubticksPerTick,
						tc.args.clobPair.StepBaseQuantums,
						perpetual.Params.LiquidityTier,
						perpetual.Params.MarketType,
						perpetual.Params.DefaultFundingPpm,
					),
				),
			).Return()
			_, err := ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
				ks.Ctx,
				tc.args.clobPair.Id,
				clobtest.MustPerpetualId(tc.args.clobPair),
				satypes.BaseQuantums(tc.args.clobPair.StepBaseQuantums),
				tc.args.clobPair.QuantumConversionExponent,
				tc.args.clobPair.SubticksPerTick,
				tc.args.clobPair.Status,
			)
			require.NoError(t, err)

			// override clob keeper's PerpetualIdToClobPairId
			if tc.perpetualIdToClobPairId != nil {
				ks.ClobKeeper.PerpetualIdToClobPairId = tc.perpetualIdToClobPairId
			}

			premiumPpm, err := ks.ClobKeeper.GetPricePremiumForPerpetual(
				ks.Ctx,
				tc.perpetualId,
				perptypes.GetPricePremiumParams{
					IndexPrice:                  tc.args.indexPrice,
					BaseAtomicResolution:        tc.args.baseAtomicResolution,
					QuoteAtomicResolution:       tc.args.quoteAtomicResolution,
					ImpactNotionalQuoteQuantums: tc.args.impactNotionalAmount,
					MaxAbsPremiumVotePpm:        tc.args.maxAbsPremiumVotePpm,
				},
			)

			if tc.expectedErr != nil {
				require.ErrorContains(t,
					tc.expectedErr,
					err.Error(),
				)
				return
			}

			require.Equal(t,
				tc.args.premiumPpm,
				premiumPpm,
			)
		})
	}
}
