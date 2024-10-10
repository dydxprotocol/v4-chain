package constants

const (
	AppName       = "dydxprotocol"
	AppDaemonName = AppName + "d"
	ServiceName   = "validator"

	// MaximumPriceSize defines the maximum size of a price in bytes. This allows
	// up to 32 bytes for the price and 1 byte for the sign (positive/negative).
	MaximumPriceSizeInBytes = 33

	// MaxSDaiConversionRateLengthCharacters is the maximum length of the sDai conversion rate in characters.
	MaxSDaiConversionRateLengthCharacters = 35

	// where in the proposal the injected VE's are located
	DaemonInfoIndex    = 0
	InjectedNonTxCount = 1
	// block structure
	// this is three because the first place in the block is for VE's
	MinTxsCount                     = 2
	MinTxsCountWithVE               = MinTxsCount + 1
	ProposedOperationsTxIndex       = 0
	ProposedOperationsTxIndexWithVE = ProposedOperationsTxIndex + 1
	AddPremiumVotesTxLenOffset      = -1
	LastOtherTxLenOffset            = AddPremiumVotesTxLenOffset
	FirstOtherTxIndex               = ProposedOperationsTxIndex + 1
)
