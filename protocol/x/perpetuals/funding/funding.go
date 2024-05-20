package funding

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// GetFundingIndexDelta returns `fundingIndexDelta` which represents the change of the funding index
// given the funding rate, the time since the last funding tick, and the oracle price. The index delta
// is in parts-per-million (PPM) and is calculated as follows:
//
//	 indexDelta =
//		  fundingRatePpm *
//	   (time / realizationPeriod) *
//	   quoteQuantumsPerBaseQuantum
//
// Any multiplication is done before division to avoid precision loss.
func GetFundingIndexDelta(
	perp types.Perpetual,
	marketPrice pricestypes.MarketPrice,
	big8hrFundingRatePpm *big.Int,
	timeSinceLastFunding uint32,
) (fundingIndexDelta *big.Int) {
	// Get pro-rated funding rate adjusted by time delta.
	result := new(big.Int).SetUint64(uint64(timeSinceLastFunding))

	// Multiply by the time-delta numerator upfront.
	result.Mul(result, big8hrFundingRatePpm)

	// Multiply by the price of the asset.
	result = lib.BaseToQuoteQuantums(
		result,
		perp.Params.AtomicResolution,
		marketPrice.Price,
		marketPrice.Exponent,
	)

	// Divide by the time-delta denominator.
	// Use truncated division (towards zero) instead of Euclidean division.
	// TODO(DEC-1536): Make the 8-hour funding rate period configurable.
	result.Quo(result, big.NewInt(60*60*8))

	return result
}
