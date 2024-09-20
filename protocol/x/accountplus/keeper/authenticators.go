package keeper

import (
	"fmt"
	"strconv"

	"cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

// AddAuthenticator adds an authenticator to an account, this function is used to add multiple
// authenticators such as SignatureVerifications and AllOfs
func (k Keeper) AddAuthenticator(
	ctx sdk.Context,
	account sdk.AccAddress,
	authenticatorType string,
	config []byte,
) (uint64, error) {
	impl := k.authenticatorManager.GetAuthenticatorByType(authenticatorType)
	if impl == nil {
		return 0, fmt.Errorf("authenticator type %s is not registered", authenticatorType)
	}

	// Get the next global id value for authenticators from the store
	id := k.InitializeOrGetNextAuthenticatorId(ctx)

	// Each authenticator has a custom OnAuthenticatorAdded function
	err := impl.OnAuthenticatorAdded(ctx, account, config, strconv.FormatUint(id, 10))
	if err != nil {
		return 0, errors.Wrapf(err, "`OnAuthenticatorAdded` failed on authenticator type %s", authenticatorType)
	}

	k.SetNextAuthenticatorId(ctx, id+1)

	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.AuthenticatorKeyPrefix),
	)
	authenticator := types.AccountAuthenticator{
		Id:     id,
		Type:   authenticatorType,
		Config: config,
	}
	b := k.cdc.MustMarshal(&authenticator)
	store.Set(types.KeyAccountId(account, id), b)
	return id, nil
}

// RemoveAuthenticator removes an authenticator from an account
func (k Keeper) RemoveAuthenticator(ctx sdk.Context, account sdk.AccAddress, authenticatorId uint64) error {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.AuthenticatorKeyPrefix),
	)
	key := types.KeyAccountId(account, authenticatorId)

	var existing types.AccountAuthenticator
	bz := store.Get(key)
	if bz == nil {
		return errors.Wrapf(
			types.ErrAuthenticatorNotFound,
			"RemoveAuthenticator: failed to get authenticator %d for address %s",
			authenticatorId,
			account.String(),
		)
	}
	k.cdc.MustUnmarshal(bz, &existing)

	impl := k.authenticatorManager.GetAuthenticatorByType(existing.Type)
	if impl == nil {
		return fmt.Errorf("authenticator type %s is not registered", existing.Type)
	}

	// Authenticators can prevent removal. This should be used sparingly
	err := impl.OnAuthenticatorRemoved(ctx, account, existing.Config, strconv.FormatUint(authenticatorId, 10))
	if err != nil {
		return errors.Wrapf(err, "`OnAuthenticatorRemoved` failed on authenticator type %s", existing.Type)
	}

	store.Delete(key)
	return nil
}

// InitializeOrGetNextAuthenticatorId returns the next authenticator id.
func (k Keeper) InitializeOrGetNextAuthenticatorId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)

	b := store.Get([]byte(types.AuthenticatorIdKeyPrefix))
	if b == nil {
		return 0
	}

	result := gogotypes.UInt64Value{Value: 0}
	k.cdc.MustUnmarshal(b, &result)
	return result.Value
}

// SetNextAuthenticatorId sets next authenticator id.
func (k Keeper) SetNextAuthenticatorId(ctx sdk.Context, authenticatorId uint64) {
	value := gogotypes.UInt64Value{Value: authenticatorId}
	b := k.cdc.MustMarshal(&value)

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.AuthenticatorIdKeyPrefix), b)
}
