package types

const (
	MinPriceChangePpm_LongTail uint32 = 800

	DefaultFundingPpm = 100 // 1bps per 8 hour or 0.125bps per hour

	LiquidityTier_Isolated uint32 = 4

	LiquidityTier_LongTail uint32 = 2

	DefaultStepBaseQuantums uint64 = 1_000_000

	SubticksPerTick_LongTail uint32 = 1_000_000

	DefaultQuantumConversionExponent = -9

	ResolutionOffset = -6

	DefaultMarketsHardCap = 500
)
