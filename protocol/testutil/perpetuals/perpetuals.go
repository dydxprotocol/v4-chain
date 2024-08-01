package perpetuals

import (
	"fmt"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricetypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

type PerpetualModifierOption func(cp *perptypes.Perpetual)

func WithId(id uint32) PerpetualModifierOption {
	return func(cp *perptypes.Perpetual) {
		cp.Params.Id = id
	}
}

func WithMarketId(id uint32) PerpetualModifierOption {
	return func(cp *perptypes.Perpetual) {
		cp.Params.MarketId = id
	}
}

func WithPerpetual(perp perptypes.Perpetual) PerpetualModifierOption {
	return func(cp *perptypes.Perpetual) {
		cp.Params = perp.Params
	}
}

func WithTicker(ticker string) PerpetualModifierOption {
	return func(cp *perptypes.Perpetual) {
		cp.Params.Ticker = ticker
	}
}

func WithLiquidityTier(liquidityTier uint32) PerpetualModifierOption {
	return func(cp *perptypes.Perpetual) {
		cp.Params.LiquidityTier = liquidityTier
	}
}

func WithAtomicResolution(atomicResolution int32) PerpetualModifierOption {
	return func(cp *perptypes.Perpetual) {
		cp.Params.AtomicResolution = atomicResolution
	}
}

func WithMarketType(marketType perptypes.PerpetualMarketType) PerpetualModifierOption {
	return func(cp *perptypes.Perpetual) {
		cp.Params.MarketType = marketType
	}
}

// GeneratePerpetual returns a `Perpetual` object set to default values.
// Passing in `PerpetualModifierOption` methods alters the value of the `Perpetual` returned.
// It will start with the default, valid `Perpetual` value defined within the method
// and make the requested modifications before returning the object.
//
// Example usage:
// `GeneratePerpetual(WithId(10))`
// This will start with the default `Perpetual` object defined within the method and
// return the newly-created object after overriding the values of
// `Id` to 10.
func GeneratePerpetual(optionalModifications ...PerpetualModifierOption) *perptypes.Perpetual {
	perpetual := &perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                0,
			Ticker:            "BTC-USD",
			MarketId:          0,
			AtomicResolution:  -8,
			DefaultFundingPpm: 0,
			LiquidityTier:     0,
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}

	for _, opt := range optionalModifications {
		opt(perpetual)
	}

	return perpetual
}

// MustHumanSizeToBaseQuantums converts a human-readable size to quantums.
// It uses the inverse of the exponent to convert the human size to quantums,
// since the exponent applies to the quantums to derive the human-readable size.
func MustHumanSizeToBaseQuantums(
	humanSize string,
	atomicResolution int32,
) (baseQuantums uint64) {
	ratio, ok := new(big.Rat).SetString(humanSize)
	if !ok {
		panic(fmt.Sprintf("MustHumanSizeToBaseQuantums: Failed to parse humanSize: %s", humanSize))
	}
	result := lib.BigIntMulPow10(ratio.Num(), -atomicResolution, false)
	result.Quo(result, ratio.Denom())
	if !result.IsUint64() {
		panic("MustHumanSizeToBaseQuantums: result is not a uint64")
	}
	return result.Uint64()
}

// Helper function to set up default open interest for input perpetuals.
func SetUpDefaultPerpOIsForTest(
	t testing.TB,
	ctx sdk.Context,
	k perptypes.PerpetualsKeeper,
	perps []perptypes.Perpetual,
) {
	for _, perpOI := range constants.DefaultTestPerpOIs {
		for _, perp := range perps {
			if perp.Params.Id != perpOI.PerpetualId {
				continue
			}
			// If the perpetual exists in input, set up the open interest.
			require.NoError(t,
				k.ModifyOpenInterest(
					ctx,
					perp.Params.Id,
					perpOI.BaseQuantums,
				),
			)
		}
	}
}

func CreatePerpInfo(
	id uint32,
	atomicResolution int32,
	price uint64,
	priceExponent int32,
) perptypes.PerpInfo {
	return perptypes.PerpInfo{
		Perpetual: perptypes.Perpetual{
			Params: perptypes.PerpetualParams{
				Id:               id,
				Ticker:           "test ticker",
				MarketId:         id,
				AtomicResolution: atomicResolution,
				LiquidityTier:    id,
			},
			FundingIndex: dtypes.NewInt(0),
			OpenInterest: dtypes.NewInt(0),
		},
		Price: pricetypes.MarketPrice{
			Id:       id,
			Exponent: priceExponent,
			Price:    price,
		},
		LiquidityTier: perptypes.LiquidityTier{
			Id:                     id,
			InitialMarginPpm:       100_000,
			MaintenanceFractionPpm: 500_000,
			OpenInterestLowerCap:   0,
			OpenInterestUpperCap:   0,
		},
	}
}
