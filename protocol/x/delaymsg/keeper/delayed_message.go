package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"sort"
)

// newDelayedMessageStore returns a prefix store for delayed messages.
func (k Keeper) newDelayedMessageStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.DelayedMessageKeyPrefix))
}

// GetNumMessages returns the number of messages in the store.
func (k Keeper) GetNumMessages(
	ctx sdk.Context,
) uint32 {
	store := ctx.KVStore(k.storeKey)
	var numMessagesBytes = store.Get([]byte(types.NumDelayedMessagesKeyPrefix))
	return lib.BytesToUint32(numMessagesBytes)
}

// SetNumMessages sets the number of messages in the store.
func (k Keeper) SetNumMessages(
	ctx sdk.Context,
	numMessages uint32,
) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.NumDelayedMessagesKeyPrefix), lib.Uint32ToBytes(numMessages))
}

// GetMessage returns a message from its id.
func (k Keeper) GetMessage(
	ctx sdk.Context,
	id uint32,
) (
	delayedMessage types.DelayedMessage,
	found bool,
) {
	store := k.newDelayedMessageStore(ctx)
	b := store.Get(lib.Uint32ToBytesForState(id))
	if b == nil {
		return types.DelayedMessage{}, false
	}

	delayedMessage = types.DelayedMessage{}
	k.cdc.MustUnmarshal(b, &delayedMessage)
	return delayedMessage, true
}

// GetAllDelayedMessages returns all messages in the store, sorted by id ascending.
// This is primarily used to export the genesis state.
func (k Keeper) GetAllDelayedMessages(ctx sdk.Context) []*types.DelayedMessage {
	delayedMessages := make([]*types.DelayedMessage, 0)
	store := k.newDelayedMessageStore(ctx)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delayedMessage := types.DelayedMessage{}
		k.cdc.MustUnmarshal(iterator.Value(), &delayedMessage)
		delayedMessages = append(delayedMessages, &delayedMessage)
	}

	// Sort delayed messages by id ascending. This is useful for tests and ease of reading / checking
	// genesis state exports.
	sort.Slice(delayedMessages, func(i, j int) bool {
		return delayedMessages[i].Id < delayedMessages[j].Id
	})

	return delayedMessages
}

// DeleteMessage deletes a message from the store.
func (k Keeper) DeleteMessage(
	ctx sdk.Context,
	id uint32,
) (
	err error,
) {
	delayedMsg, found := k.GetMessage(ctx, id)
	if !found {
		return sdkerrors.Wrapf(
			types.ErrInvalidInput,
			"failed to delete message: message with id %d not found",
			id,
		)
	}
	store := k.newDelayedMessageStore(ctx)
	store.Delete(lib.Uint32ToBytesForState(id))

	// Remove message id from block message ids.
	if err := k.deleteMessageIdFromBlock(ctx, id, delayedMsg.BlockHeight); err != nil {
		return sdkerrors.Wrapf(
			types.ErrInvalidInput,
			"failed to delete message: %v",
			err,
		)
	}
	return nil
}

// SetDelayedMessage delays a message to be executed at the specified block height. The delayed
// message is assigned the specified id. This method is suitable for initializing from genesis state.
func (k Keeper) SetDelayedMessage(
	ctx sdk.Context,
	msg *types.DelayedMessage,
) (
	err error,
) {
	if msg.BlockHeight < ctx.BlockHeight() {
		return sdkerrors.Wrapf(
			types.ErrInvalidInput,
			"failed to delay message: block height %d is in the past",
			msg.BlockHeight,
		)
	}

	// Add message to the store.
	store := k.newDelayedMessageStore(ctx)
	store.Set(lib.Uint32ToBytesForState(msg.Id), k.cdc.MustMarshal(msg))

	// Add message id to the list of message ids for the block.
	k.addMessageIdToBlock(ctx, msg.Id, msg.BlockHeight)
	return nil
}

// DecodeMessage decodes message bytes into a sdk.Msg. This method is added to the keeper
// to allow for mocking in msg service tests.
func (k Keeper) DecodeMessage(msgBytes []byte, msg *sdk.Msg) error {
	return k.cdc.UnmarshalInterface(msgBytes, msg)
}

// EncodeMessage encodes a sdk.Msg into bytes. This method is added to the keeper
// for ease of testing.
func (k Keeper) EncodeMessage(msg sdk.Msg) ([]byte, error) {
	return k.cdc.MarshalInterface(msg)
}

// DelayMessageByBlocks registers an sdk.Msg to be executed after blockDelay blocks.
func (k Keeper) DelayMessageByBlocks(
	ctx sdk.Context,
	msg sdk.Msg,
	blockDelay uint32,
) (
	id uint32,
	err error,
) {
	handler := k.router.Handler(msg)
	// If the message type is not routable, return an error.
	if handler == nil {
		return 0, sdkerrors.Wrapf(
			types.ErrMsgIsUnroutable,
			sdk.MsgTypeURL(msg),
		)
	}

	nextId := k.GetNumMessages(ctx)
	blockHeight, err := lib.AddUint32(ctx.BlockHeight(), blockDelay)
	if err != nil {
		return 0, sdkerrors.Wrapf(
			types.ErrInvalidInput,
			"failed to add block delay to current block height: %v",
			err,
		)
	}

	messageBytes, err := k.cdc.MarshalInterface(msg)
	if err != nil {
		return 0, sdkerrors.Wrapf(
			types.ErrInvalidInput,
			"failed to marshal message: %v. Is the message type registered?",
			err,
		)
	}

	delayedMessage := types.DelayedMessage{
		Id:          nextId,
		Msg:         messageBytes,
		BlockHeight: blockHeight,
	}

	err = k.SetDelayedMessage(ctx, &delayedMessage)
	if err != nil {
		return 0, err
	}

	// Increment the number of messages in the store.
	k.SetNumMessages(ctx, nextId+1)

	return nextId, nil
}
