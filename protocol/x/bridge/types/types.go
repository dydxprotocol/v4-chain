package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BridgeKeeper interface {
	// Bridge Events
	GetAcknowledgedEventInfo(ctx sdk.Context) BridgeEventInfo

	GetRecognizedEventInfo(ctx sdk.Context) BridgeEventInfo

	AcknowledgeBridges(ctx sdk.Context, bridges []BridgeEvent) error

	CompleteBridge(ctx sdk.Context, bridges BridgeEvent) error

	// Event Params
	GetEventParams(ctx sdk.Context) EventParams

	UpdateEventParams(ctx sdk.Context, params EventParams) error

	// Propose Params
	GetProposeParams(ctx sdk.Context) ProposeParams

	UpdateProposeParams(ctx sdk.Context, params ProposeParams) error

	// Safety Params
	GetSafetyParams(ctx sdk.Context) SafetyParams

	UpdateSafetyParams(ctx sdk.Context, params SafetyParams) error

	// Authority.
	HasAuthority(authority string) bool
}
