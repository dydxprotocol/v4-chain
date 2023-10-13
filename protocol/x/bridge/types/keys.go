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
	AcknowledgedEventInfoKey = "AckEventInfo"

	// EventParamsKey defines the key for the EventParams
	EventParamsKey = "EventParams"

	// ProposeParamsKey defines the key for the ProposeParams
	ProposeParamsKey = "ProposeParams"

	// SafetyParamsKey defines the key for the SafetyParams
	SafetyParamsKey = "SafetyParams"
)
