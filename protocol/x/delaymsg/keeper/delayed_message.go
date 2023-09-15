package keeper

import (
	"bytes"
	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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
		return errorsmod.Wrapf(
			types.ErrInvalidInput,
			"failed to delete message: message with id %d not found",
			id,
		)
	}
	store := k.newDelayedMessageStore(ctx)
	store.Delete(lib.Uint32ToBytesForState(id))

	// Remove message id from block message ids.
	if err := k.deleteMessageIdFromBlock(ctx, id, delayedMsg.BlockHeight); err != nil {
		return errorsmod.Wrapf(
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
		return errorsmod.Wrapf(
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

// validateSigners validates that the message has exactly one signer, and that the signer is the delaymsg module
// address.
func validateSigners(msg sdk.Msg) error {
	signers := msg.GetSigners()
	if len(signers) != 1 {
		return errorsmod.Wrapf(
			types.ErrInvalidSigner,
			"message must have exactly one signer",
		)
	}
	moduleAddress := authtypes.NewModuleAddress(types.ModuleName)
	if !bytes.Equal(signers[0], moduleAddress) {
		return errorsmod.Wrapf(
			types.ErrInvalidSigner,
			"message signer must be delaymsg module address",
		)
	}
	return nil
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
		return 0, errorsmod.Wrapf(
			types.ErrMsgIsUnroutable,
			sdk.MsgTypeURL(msg),
		)
	}

	if err := msg.ValidateBasic(); err != nil {
		return 0, errorsmod.Wrapf(
			types.ErrInvalidInput,
			"message failed basic validation: %v",
			err,
		)
	}

	if err := validateSigners(msg); err != nil {
		return 0, err
	}

	nextId := k.GetNumMessages(ctx)
	blockHeight, err := lib.AddUint32(ctx.BlockHeight(), blockDelay)
	if err != nil {
		return 0, errorsmod.Wrapf(
			types.ErrInvalidInput,
			"failed to add block delay to current block height: %v",
			err,
		)
	}

	anyMsg, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return 0, errorsmod.Wrapf(
			types.ErrInvalidInput,
			"failed to convert message to Any: %v",
			err,
		)
	}

	delayedMessage := types.DelayedMessage{
		Id:          nextId,
		Msg:         anyMsg,
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
