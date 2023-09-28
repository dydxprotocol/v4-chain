package types

// Module name and store keys.
const (
	// ModuleName defines the module name
	ModuleName = "bridge"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// State
const (
	// AcknowledgedEventInfoKey defines the key for the AcknowledgedEventInfo
	AcknowledgedEventInfoKey = "acknowledged_event_info"

	EventParamsKey   = "event_params"
	ProposeParamsKey = "propose_params"
	SafetyParamsKey  = "safety_params"
)
