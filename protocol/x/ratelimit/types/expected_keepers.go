package types

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// BankKeeper defines the expected bank keeper used for simulations.
type BankKeeper interface {
	GetSupply(ctx context.Context, denom string) sdk.Coin
}

type BlockTimeKeeper interface {
	GetTimeSinceLastBlock(ctx sdk.Context) time.Duration
}

// ICS4Wrapper defines the expected ICS4Wrapper for middleware
type ICS4Wrapper interface {
	WriteAcknowledgement(
		ctx sdk.Context,
		chanCap *capabilitytypes.Capability,
		packet ibcexported.PacketI,
		acknowledgement ibcexported.Acknowledgement,
	) error
	SendPacket(
		ctx sdk.Context,
		chanCap *capabilitytypes.Capability,
		sourcePort string,
		sourceChannel string,
		timeoutHeight clienttypes.Height,
		timeoutTimestamp uint64,
		data []byte,
	) (sequence uint64, err error)
	GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool)
}
