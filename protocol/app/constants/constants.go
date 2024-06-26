package constants

const (
	AppName       = "dydxprotocol"
	AppDaemonName = AppName + "d"
	ServiceName   = "validator"

	// MaximumPriceSize defines the maximum size of a price in bytes. This allows
	// up to 32 bytes for the price and 1 byte for the sign (positive/negative).
	MaximumPriceSize = 33

	// NumInjectedTxs is the number of transactions that were injected into
	// the proposal but are not actual transactions. In this case, the oracle
	// info is injected into the proposal but should be ignored by the application.
	NumInjectedTxs = 1

	// where in the proposal the injected VE's are located
	DeamonInfoIndex = 0
)
