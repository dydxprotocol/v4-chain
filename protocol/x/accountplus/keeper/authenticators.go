package keeper

import (
	"fmt"
	"strconv"

	"cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/cosmos/gogoproto/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/authenticator"
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

// GetAuthenticatorExtension unpacks the extension for the transaction, this is used with transactions specify
// an authenticator to use
func (k Keeper) GetAuthenticatorExtension(exts []*codectypes.Any) types.AuthenticatorTxOptions {
	var authExtension types.AuthenticatorTxOptions
	for _, ext := range exts {
		err := k.cdc.UnpackAny(ext, &authExtension)
		if err == nil {
			return authExtension
		}
	}
	return nil
}

// GetSelectedAuthenticatorData gets all authenticators from an account
// from the store, the data is  prefixed by 2|<accAddr|<keyId>
func (k Keeper) GetSelectedAuthenticatorData(
	ctx sdk.Context,
	account sdk.AccAddress,
	selectedAuthenticator int,
) (*types.AccountAuthenticator, error) {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.AuthenticatorKeyPrefix),
	)
	bz := store.Get(types.KeyAccountId(account, uint64(selectedAuthenticator)))
	if bz == nil {
		return &types.AccountAuthenticator{}, errors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("authenticator %d not found for account %s", selectedAuthenticator, account),
		)
	}
	authenticatorFromStore, err := k.unmarshalAccountAuthenticator(bz)
	if err != nil {
		return &types.AccountAuthenticator{}, err
	}

	return authenticatorFromStore, nil
}

// GetInitializedAuthenticatorForAccount returns a single initialized authenticator for the account.
// It fetches the authenticator data from the store, gets the authenticator struct from the manager,
// then calls initialize on the authenticator data
func (k Keeper) GetInitializedAuthenticatorForAccount(
	ctx sdk.Context,
	account sdk.AccAddress,
	selectedAuthenticator int,
) (authenticator.InitializedAuthenticator, error) {
	// Get the authenticator data from the store
	authenticatorFromStore, err := k.GetSelectedAuthenticatorData(ctx, account, selectedAuthenticator)
	if err != nil {
		return authenticator.InitializedAuthenticator{}, err
	}

	uninitializedAuthenticator := k.authenticatorManager.GetAuthenticatorByType(authenticatorFromStore.Type)
	if uninitializedAuthenticator == nil {
		// This should never happen, but if it does, it means that stored authenticator is not registered
		// or somehow the registered authenticator was removed / malformed
		telemetry.IncrCounter(1, metrics.MissingRegisteredAuthenticator)
		k.Logger(ctx).Error(
			"account asscoicated authenticator not registered in manager",
			"type", authenticatorFromStore.Type,
			"id", selectedAuthenticator,
		)

		return authenticator.InitializedAuthenticator{},
			errors.Wrapf(
				sdkerrors.ErrLogic,
				"authenticator id %d failed to initialize, authenticator type %s not registered in manager",
				selectedAuthenticator, authenticatorFromStore.Type,
			)
	}
	// Ensure that initialization of each authenticator works as expected
	// NOTE: Always return a concrete authenticator not a pointer, do not modify in place
	// NOTE: The authenticator manager returns a struct that is reused
	initializedAuthenticator, err := uninitializedAuthenticator.Initialize(authenticatorFromStore.Config)
	if err != nil || initializedAuthenticator == nil {
		return authenticator.InitializedAuthenticator{},
			errors.Wrapf(err,
				"authenticator %d with type %s failed to initialize",
				selectedAuthenticator, authenticatorFromStore.Type,
			)
	}

	finalAuthenticator := authenticator.InitializedAuthenticator{
		Id:            authenticatorFromStore.Id,
		Authenticator: initializedAuthenticator,
	}

	return finalAuthenticator, nil
}

// unmarshalAccountAuthenticator is used to unmarshal the AccountAuthenticator from the store
func (k Keeper) unmarshalAccountAuthenticator(bz []byte) (*types.AccountAuthenticator, error) {
	var accountAuthenticator types.AccountAuthenticator
	err := k.cdc.Unmarshal(bz, &accountAuthenticator)
	if err != nil {
		return &types.AccountAuthenticator{}, errors.Wrap(err, "failed to unmarshal account authenticator")
	}
	return &accountAuthenticator, nil
}
