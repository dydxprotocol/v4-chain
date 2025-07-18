package types

import "math/big"

type FillForProcess struct {
	TakerAddr             string
	TakerFeeQuoteQuantums *big.Int
	MakerAddr             string
	MakerFeeQuoteQuantums *big.Int
	FillQuoteQuantums     *big.Int
	ProductId             uint32
	// MonthlyRollingTakerVolumeQuantums is the total taker volume for
	// the given taker address in the last 30 days. This rolling volume
	// does not include stats of the current block being processed.
	// If there are multiple fills for the taker address in the
	// same block, this volume will not be included in the function
	// below
	MonthlyRollingTakerVolumeQuantums uint64
	MarketId                          uint32
	// Order router addresses for taker and maker
	TakerOrderRouterAddr string
	MakerOrderRouterAddr string
}
