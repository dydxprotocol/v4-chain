package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	// Input args for `MemClob.GetClobPairForPerpetual`
	perpetualId uint32

	// Return values from `MemClob.GetClobPairForPerpetual`
	clobPairId                 types.ClobPairId
	getClobPairForPerpetualErr error

	// Input args for `MemClob.GetPricePremium`
	clobPair              types.ClobPair
	marketPrice           pricestypes.MarketPrice
	baseAtomicResolution  int32
	quoteAtomicResolution int32
	impactNotionalAmount  *big.Int
	maxAbsPremiumVotePpm  *big.Int

	// Return values from `MemClob.GetPricePremium`
	premiumPpm         int32
	getPricePremiumErr error
}

type testCase map[string]struct {
	setUpMockMemClob func(mck *mocks.MemClob, args testMemClobMethodArgs)
	args             testMemClobMethodArgs
	expectedErr      error
}

func TestGetPricePremiumForPerpetual(t *testing.T) {
	tests := testCase{
		"Success": {
			args: testMemClobMethodArgs{
				perpetualId: 0,
				clobPairId:  0,
				clobPair:    constants.ClobPair_Btc,
				marketPrice: pricestypes.MarketPrice{
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
					"GetClobPairForPerpetual",
					mock.Anything,
					args.perpetualId,
				).Return(
					args.clobPairId,
					args.getClobPairForPerpetualErr,
				)
				mck.On(
					"GetPricePremium",
					mock.Anything,
					args.clobPair,
					perptypes.GetPricePremiumParams{
						MarketPrice:                 args.marketPrice,
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
		"Failure: GetClobPairForPerpetual error": {
			args: testMemClobMethodArgs{
				perpetualId:                0,
				clobPair:                   constants.ClobPair_Btc,
				getClobPairForPerpetualErr: errors.New("GetClobPairForPerpetual error"),
			},
			setUpMockMemClob: func(mck *mocks.MemClob, args testMemClobMethodArgs) {
				mck.On(
					"GetClobPairForPerpetual",
					mock.Anything,
					args.perpetualId,
				).Return(
					args.clobPairId,
					args.getClobPairForPerpetualErr,
				)
			},
			expectedErr: errors.New("GetClobPairForPerpetual error"),
		},
		"Failure: GetPricePremium failure": {
			args: testMemClobMethodArgs{
				perpetualId: 0,
				clobPairId:  0,
				clobPair:    constants.ClobPair_Btc,
				marketPrice: pricestypes.MarketPrice{
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
					"GetClobPairForPerpetual",
					mock.Anything,
					args.perpetualId,
				).Return(
					args.clobPairId,
					args.getClobPairForPerpetualErr,
				)
				mck.On(
					"GetPricePremium",
					mock.Anything,
					args.clobPair,
					perptypes.GetPricePremiumParams{
						MarketPrice:                 args.marketPrice,
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
			args: testMemClobMethodArgs{
				perpetualId: 0,
				clobPairId:  types.ClobPairId(1),
				clobPair:    constants.ClobPair_Btc, // clobPairId = 1000
			},
			setUpMockMemClob: func(mck *mocks.MemClob, args testMemClobMethodArgs) {
				mck.On(
					"GetClobPairForPerpetual",
					mock.Anything,
					args.perpetualId,
				).Return(
					args.clobPairId,
					args.getClobPairForPerpetualErr,
				)
			},
			expectedErr: sdkerrors.Wrapf(
				types.ErrInvalidClob,
				"GetPricePremiumForPerpetual: did not find clob pair with clobPairId = %d",
				1,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything)
			memClob.On("CreateOrderbook", mock.Anything, mock.Anything, mock.Anything)

			tc.setUpMockMemClob(memClob, tc.args)

			ks := keepertest.NewClobKeepersTestContext(t, memClob, nil, &mocks.IndexerEventManager{})

			prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			_, err := ks.ClobKeeper.CreatePerpetualClobPair(
				ks.Ctx,
				tc.args.clobPair.Id,
				clobtest.MustPerpetualId(tc.args.clobPair),
				satypes.BaseQuantums(tc.args.clobPair.StepBaseQuantums),
				tc.args.clobPair.QuantumConversionExponent,
				tc.args.clobPair.SubticksPerTick,
				tc.args.clobPair.Status,
				tc.args.clobPair.MakerFeePpm,
				tc.args.clobPair.TakerFeePpm,
			)
			require.NoError(t, err)

			premiumPpm, err := ks.ClobKeeper.GetPricePremiumForPerpetual(
				ks.Ctx,
				tc.args.perpetualId,
				perptypes.GetPricePremiumParams{
					MarketPrice:                 tc.args.marketPrice,
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
