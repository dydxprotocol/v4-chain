package constants

const (
	AppName       = "dydxprotocol"
	AppDaemonName = AppName + "d"
	ServiceName   = "validator"

	// MaximumPriceSize defines the maximum size of a price in bytes. This allows
	// up to 32 bytes for the price and 1 byte for the sign (positive/negative).
	MaximumPriceSize = 33

	// where in the proposal the injected VE's are located
	DeamonInfoIndex = 0

	// block structure
	// this is three becuase the first place in the block is for VE's
	MinTxsCount                = 3
	ExtInfoBzIndex             = 0
	ProposedOperationsTxIndex  = 1
	AddPremiumVotesTxLenOffset = -1
	LastOtherTxLenOffset       = AddPremiumVotesTxLenOffset
	FirstOtherTxIndex          = ProposedOperationsTxIndex + 1
)
