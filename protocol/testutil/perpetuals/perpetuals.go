package perpetuals

import (
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
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
		},
		FundingIndex: dtypes.ZeroInt(),
	}

	for _, opt := range optionalModifications {
		opt(perpetual)
	}

	return perpetual
}
