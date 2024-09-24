package perpetuals

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
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

func WithMarketType(marketType perptypes.PerpetualMarketType) PerpetualModifierOption {
	return func(cp *perptypes.Perpetual) {
		cp.Params.MarketType = marketType
	}
}

func WithIsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock(delta uint64) PerpetualModifierOption {
	return func(cp *perptypes.Perpetual) {
		cp.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock = delta
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
			DangerIndexPpm:    0,
			IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock: 0,
		},
		FundingIndex:    dtypes.ZeroInt(),
		OpenInterest:    dtypes.ZeroInt(),
		LastFundingRate: dtypes.ZeroInt(),
	}

	for _, opt := range optionalModifications {
		opt(perpetual)
	}

	return perpetual
}

func MustHumanSizeToBaseQuantums(
	humanSize string,
	atomicResolution int32,
) (baseQuantums uint64) {
	// Parse the humanSize string to a big rational
	ratValue, ok := new(big.Rat).SetString(humanSize)
	if !ok {
		panic("Failed to parse humanSize to big.Rat")
	}

	// Convert atomicResolution to int64 for calculations
	resolution := int64(atomicResolution)

	// Create a multiplier which is 10 raised to the power of the absolute atomicResolution
	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(abs(resolution)), nil)

	// Depending on the sign of atomicResolution, multiply or divide
	if atomicResolution > 0 {
		ratValue.Mul(ratValue, new(big.Rat).SetInt(multiplier))
	} else if atomicResolution < 0 {
		divisor := new(big.Rat).SetInt(multiplier)
		ratValue.Mul(ratValue, divisor)
	}

	// Convert the result to an unsigned 64-bit integer
	resultInt := ratValue.Num() // Get the numerator which now represents the whole value

	return resultInt.Uint64()
}

// Helper function to get the absolute value of an int64
func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// Helper function to set up default open interest for input perpetuals.
func SetUpDefaultPerpOIsForTest(
	t *testing.T,
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
