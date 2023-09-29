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

	// EventParamsKey defines the key for the EventParams
	EventParamsKey = "event_params"

	// ProposeParamsKey defines the key for the ProposeParams
	ProposeParamsKey = "propose_params"

	// SafetyParamsKey defines the key for the SafetyParams
	SafetyParamsKey = "safety_params"
)
