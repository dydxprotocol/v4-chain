package keeper

import (
	"fmt"
	"strconv"

	"cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	gogotypes "github.com/cosmos/gogoproto/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

// MaybeValidateAuthenticators checks if the transaction has authenticators specified and if so,
// validates them. It returns an error if the authenticators are not valid or removed from state.
func (k Keeper) MaybeValidateAuthenticators(ctx sdk.Context, tx sdk.Tx) error {
	// Check if the tx had authenticator specified.
	specified, txOptions := lib.HasSelectedAuthenticatorTxExtensionSpecified(tx, k.cdc)
	if !specified {
		return nil
	}

	// The tx had authenticators specified.
	// First make sure smart account flow is enabled.
	if active := k.GetIsSmartAccountActive(ctx); !active {
		return types.ErrSmartAccountNotActive
	}

	// Make sure txn is a SigVerifiableTx and get signers from the tx.
	sigVerifiableTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return errors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	signers, err := sigVerifiableTx.GetSigners()
	if err != nil {
		return err
	}

	if len(signers) != 1 {
		return errors.Wrap(types.ErrTxnHasMultipleSigners, "only one signer is allowed")
	}

	account := sdk.AccAddress(signers[0])

	// Retrieve the selected authenticators from the extension and make sure they are valid, i.e. they
	// are registered and not removed from state.
	//
	// Note that we only verify the existence of the authenticators here without actually
	// runnning them. This is because all current authenticators are stateless and do not read/modify any states.
	selectedAuthenticators := txOptions.GetSelectedAuthenticators()
	for _, authenticatorId := range selectedAuthenticators {
		_, err := k.GetInitializedAuthenticatorForAccount(
			ctx,
			account,
			authenticatorId,
		)
		if err != nil {
			return errors.Wrapf(
				err,
				"selected authenticator (%s, %d) is not registered or removed from state",
				account.String(),
				authenticatorId,
			)
		}
	}
	return nil
}

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
	requireSigVerification, err := impl.OnAuthenticatorAdded(ctx, account, config, strconv.FormatUint(id, 10))
	if err != nil {
		return 0, errors.Wrapf(err, "`OnAuthenticatorAdded` failed on authenticator type %s", authenticatorType)
	}

	if !requireSigVerification {
		return 0, fmt.Errorf(
			"unsafe: authenticator tree does not require signature verification all possible paths, type = %v, config = %v",
			authenticatorType,
			config,
		)
	}

	k.SetNextAuthenticatorId(ctx, id+1)
	authenticator := types.AccountAuthenticator{
		Id:     id,
		Type:   authenticatorType,
		Config: config,
	}
	k.SetAuthenticator(ctx, account.String(), id, authenticator)
	return id, nil
}

func (k Keeper) SetAuthenticator(
	ctx sdk.Context,
	account string,
	authenticatorId uint64,
	authenticator types.AccountAuthenticator,
) {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.AuthenticatorKeyPrefix),
	)
	b := k.cdc.MustMarshal(&authenticator)
	store.Set(types.BuildKey(account, authenticatorId), b)
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

// GetSelectedAuthenticatorData gets a single authenticator for the account from the store.
func (k Keeper) GetSelectedAuthenticatorData(
	ctx sdk.Context,
	account sdk.AccAddress,
	selectedAuthenticator uint64,
) (*types.AccountAuthenticator, error) {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.AuthenticatorKeyPrefix),
	)
	bz := store.Get(types.KeyAccountId(account, selectedAuthenticator))
	if bz == nil {
		return &types.AccountAuthenticator{}, errors.Wrap(
			types.ErrAuthenticatorNotFound,
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
	selectedAuthenticator uint64,
) (types.InitializedAuthenticator, error) {
	// Get the authenticator data from the store
	authenticatorFromStore, err := k.GetSelectedAuthenticatorData(ctx, account, selectedAuthenticator)
	if err != nil {
		return types.InitializedAuthenticator{}, err
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
			"account", account.String(),
		)

		return types.InitializedAuthenticator{},
			errors.Wrapf(
				sdkerrors.ErrLogic,
				"authenticator id %d failed to initialize for account %s, authenticator type %s not registered in manager",
				selectedAuthenticator, account.String(), authenticatorFromStore.Type,
			)
	}
	// Ensure that initialization of each authenticator works as expected
	// NOTE: Always return a concrete authenticator not a pointer, do not modify in place
	// NOTE: The authenticator manager returns a struct that is reused
	initializedAuthenticator, err := uninitializedAuthenticator.Initialize(authenticatorFromStore.Config)
	if err != nil {
		return types.InitializedAuthenticator{},
			errors.Wrapf(
				err,
				"authenticator %d with type %s failed to initialize for account %s",
				selectedAuthenticator, authenticatorFromStore.Type, account.String(),
			)
	}
	if initializedAuthenticator == nil {
		return types.InitializedAuthenticator{},
			errors.Wrapf(
				types.ErrInitializingAuthenticator,
				"authenticator.Initialize returned nil for %d with type %s for account %s",
				selectedAuthenticator, authenticatorFromStore.Type, account.String(),
			)
	}

	finalAuthenticator := types.InitializedAuthenticator{
		Id:            authenticatorFromStore.Id,
		Authenticator: initializedAuthenticator,
	}

	return finalAuthenticator, nil
}

// GetAuthenticatorDataForAccount gets all authenticators AccAddressFromBech32 with an account
// from the store.
func (k Keeper) GetAuthenticatorDataForAccount(
	ctx sdk.Context,
	account sdk.AccAddress,
) ([]*types.AccountAuthenticator, error) {
	authenticators := make([]*types.AccountAuthenticator, 0)

	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.AuthenticatorKeyPrefix),
	)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(account.String()))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		authenticator, err := k.unmarshalAccountAuthenticator(iterator.Value())
		if err != nil {
			return nil, err
		}
		authenticators = append(authenticators, authenticator)
	}

	return authenticators, nil
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
