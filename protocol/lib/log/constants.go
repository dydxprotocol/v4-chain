package log

const (
	SourceModuleKey = "source_module"
	Error           = "error"
)

// Tag keys
// Do not have anything generic in here. For example, `Status` is too vague
// and can be clarified as `OrderStatus` or `DaemonHealthStatus`.
const (
	Module              = "module"
	TxMode              = "tx_mode"
	OperationsQueue     = "operations_queue"
	Callback            = "callback"
	BlockHeight         = "block_height"
	Msg                 = "msg"
	ProposerConsAddress = "proposer_cons_address"
	Handler             = "handler"
	Tx                  = "tx"
	OrderHash           = "order_hash"
	OrderStatus         = "order_status"

	OrderSizeOptimisticallyFilledFromMatchingQuantums = "order_size_optimistically_filled_from_matching_quantums"
	NewLocalValidatorOperationsQueue                  = "new_local_validator_operations_queue"
	LocalValidatorOperationsQueue                     = "local_validator_operations_queue"
)

// Tag values
const (
	Clob      = "x/clob"
	CheckTx   = "check_tx"
	RecheckTx = "recheck_tx"
	DeliverTx = "deliver_tx"
)

// Tag values that should be camelcased (i.e function names)
const (
	AnteHandler        = "AnteHandler"
	PlaceOrder         = "PlaceOrder"
	CancelOrder        = "CancelOrder"
	ProposedOperations = "ProposedOperations"
	BeginBlocker       = "BeginBlocker"
	EndBlocker         = "EndBlocker"
	PrepareCheckState  = "PrepareCheckState"
)
