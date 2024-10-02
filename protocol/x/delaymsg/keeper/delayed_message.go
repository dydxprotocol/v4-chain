package keeper

import (
	"bytes"
	"sort"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
)

// newDelayedMessageStore returns a prefix store for delayed messages.
func (k Keeper) newDelayedMessageStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.DelayedMessageKeyPrefix))
}

// GetNextDelayedMessageId returns the next delayed message id in the store.
func (k Keeper) GetNextDelayedMessageId(
	ctx sdk.Context,
) uint32 {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.NextDelayedMessageIdKey))
	var result gogotypes.UInt32Value
	k.cdc.MustUnmarshal(b, &result)
	return result.Value
}

// SetNextDelayedMessageId sets the next delayed message id in the store.
func (k Keeper) SetNextDelayedMessageId(
	ctx sdk.Context,
	nextId uint32,
) {
	store := ctx.KVStore(k.storeKey)
	value := gogotypes.UInt32Value{Value: nextId}
	store.Set([]byte(types.NextDelayedMessageIdKey), k.cdc.MustMarshal(&value))
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
	if err = msg.UnpackInterfaces(k.cdc); err != nil {
		return err
	}

	sdkMsg, err := msg.GetMessage()
	if err != nil {
		return errorsmod.Wrapf(
			types.ErrInvalidInput,
			"failed to delay msg: failed to get message with error '%v'",
			err,
		)
	}

	signers, _, err := k.cdc.GetMsgAnySigners(msg.Msg)
	if err != nil {
		return errorsmod.Wrapf(
			types.ErrInvalidInput,
			"failed to delay msg: failed to get signers with error '%v'",
			err,
		)
	}

	if err := k.ValidateMsg(sdkMsg, signers); err != nil {
		return errorsmod.Wrapf(
			types.ErrInvalidInput,
			"failed to delay message: failed to validate with error '%v'",
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
func validateSigners(signers [][]byte) error {
	if len(signers) != 1 {
		return errorsmod.Wrapf(
			types.ErrInvalidSigner,
			"message must have exactly one signer",
		)
	}
	if !bytes.Equal(signers[0], types.ModuleAddress) {
		return errorsmod.Wrap(
			types.ErrInvalidSigner,
			"message signer must be delaymsg module address",
		)
	}
	return nil
}

// ValidateMsg validates that a message is routable, passes ValidateBasic, and has the expected signer.
func (k Keeper) ValidateMsg(msg sdk.Msg, signers [][]byte) error {
	handler := k.router.Handler(msg)
	// If the message type is not routable, return an error.
	if handler == nil {
		return errorsmod.Wrap(
			types.ErrMsgIsUnroutable,
			sdk.MsgTypeURL(msg),
		)
	}

	if m, ok := msg.(sdk.HasValidateBasic); ok {
		if err := m.ValidateBasic(); err != nil {
			return errorsmod.Wrapf(
				types.ErrInvalidInput,
				"message failed basic validation: %v",
				err,
			)
		}
	}

	if err := validateSigners(signers); err != nil {
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

	nextId := k.GetNextDelayedMessageId(ctx)
	delayedMessage := types.DelayedMessage{
		Id:          nextId,
		Msg:         anyMsg,
		BlockHeight: lib.MustConvertIntegerToUint32(blockHeight),
	}

	err = k.SetDelayedMessage(ctx, &delayedMessage)
	if err != nil {
		return 0, err
	}

	// Increment next delayed message id in the store.
	k.SetNextDelayedMessageId(ctx, nextId+1)

	return nextId, nil
}
