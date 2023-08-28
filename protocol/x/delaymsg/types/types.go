package types

import (
	"github.com/cometbft/cometbft/libs/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

type DelayMsgKeeper interface {
	// Delayed messages
	SetNumMessages(
		ctx sdk.Context,
		numMessages uint32,
	)

	SetDelayedMessage(
		ctx sdk.Context,
		msg *DelayedMessage,
	) (
		err error,
	)
	DelayMessageByBlocks(
		ctx sdk.Context,
		msg sdk.Msg,
		blockDelay uint32,
	) (
		id uint32,
		err error,
	)

	GetMessage(
		ctx sdk.Context,
		id uint32,
	) (
		delayedMessage DelayedMessage,
		found bool,
	)

	GetAllDelayedMessages(ctx sdk.Context) []*DelayedMessage

	DeleteMessage(
		ctx sdk.Context,
		id uint32,
	) (
		err error,
	)

	// Block message ids
	GetBlockMessageIds(
		ctx sdk.Context,
		blockHeight int64,
	) (
		blockMessageIds BlockMessageIds,
		found bool,
	)

	// DecodeMessage decodes a message from bytes.
	DecodeMessage(msgBytes []byte, msg *sdk.Msg) error

	// GetAuthorities returns the set of authorities that can send delayed messages.
	GetAuthorities() map[string]struct{}

	Logger(ctx sdk.Context) log.Logger

	Router() lib.MsgRouter
}
