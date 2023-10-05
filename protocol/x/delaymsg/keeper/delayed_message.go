package keeper

import (
	"bytes"
	"sort"

	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	generic "github.com/dydxprotocol/v4-chain/protocol/generic/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
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
	b := store.Get([]byte(types.NumDelayedMessagesKey))
	var result generic.Uint32
	k.cdc.MustUnmarshal(b, &result)
	return result.Value
}

// SetNumMessages sets the number of messages in the store.
func (k Keeper) SetNumMessages(
	ctx sdk.Context,
	numMessages uint32,
) {
	store := ctx.KVStore(k.storeKey)
	value := generic.Uint32{Value: numMessages}
	store.Set([]byte(types.NumDelayedMessagesKey), k.cdc.MustMarshal(&value))
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
	b := store.Get(lib.Uint32ToKey(id))
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
	store.Delete(lib.Uint32ToKey(id))

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
	// Unpack the message and validate it.
	// For messages that are being set from genesis state, we need to unpack the Any type to hydrate the cached value.
	msg.UnpackInterfaces(k.cdc)
	sdkMsg, err := msg.GetMessage()
	if err != nil {
		return errorsmod.Wrapf(
			types.ErrInvalidInput,
			"failed to delay message: %v",
			err,
		)
	}

	if err := k.ValidateMsg(sdkMsg); err != nil {
		return errorsmod.Wrapf(
			types.ErrInvalidInput,
			"failed to delay message: %v",
			err,
		)
	}

	if msg.BlockHeight < lib.MustConvertIntegerToUint32(ctx.BlockHeight()) {
		return errorsmod.Wrapf(
			types.ErrInvalidInput,
			"failed to delay message: block height %d is in the past",
			msg.BlockHeight,
		)
	}

	// Add message to the store.
	store := k.newDelayedMessageStore(ctx)

	// Check for duplicate message id.
	if store.Get(lib.Uint32ToKey(msg.Id)) != nil {
		return errorsmod.Wrapf(
			types.ErrInvalidInput,
			"failed to delay message: message with id %d already exists",
			msg.Id,
		)
	}

	store.Set(lib.Uint32ToKey(msg.Id), k.cdc.MustMarshal(msg))

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

// ValidateMsg validates that a message is routable, passes ValidateBasic, and has the expected signer.
func (k Keeper) ValidateMsg(msg sdk.Msg) error {
	handler := k.router.Handler(msg)
	// If the message type is not routable, return an error.
	if handler == nil {
		return errorsmod.Wrapf(
			types.ErrMsgIsUnroutable,
			sdk.MsgTypeURL(msg),
		)
	}

	if err := msg.ValidateBasic(); err != nil {
		return errorsmod.Wrapf(
			types.ErrInvalidInput,
			"message failed basic validation: %v",
			err,
		)
	}

	if err := validateSigners(msg); err != nil {
		return err
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
		BlockHeight: lib.MustConvertIntegerToUint32(blockHeight),
	}

	err = k.SetDelayedMessage(ctx, &delayedMessage)
	if err != nil {
		return 0, err
	}

	// Increment the number of messages in the store.
	k.SetNumMessages(ctx, nextId+1)

	return nextId, nil
}
