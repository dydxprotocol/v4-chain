package keeper

import (
	"fmt"
	"strings"

	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/authenticator"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

	authenticatorManager *authenticator.AuthenticatorManager
	authorities          map[string]struct{}
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	authenticatorManager *authenticator.AuthenticatorManager,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:                  cdc,
		storeKey:             key,
		authenticatorManager: authenticatorManager,
		authorities:          lib.UniqueSliceToSet(authorities),
	}
}

func (k Keeper) GetStoreKey() storetypes.StoreKey {
	return k.storeKey
}

func (k Keeper) GetCdc() codec.BinaryCodec {
	return k.cdc
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(log.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

// Func required for test setup
func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}

// Get all AccountStates from kvstore
func (k Keeper) GetAllAccountStates(ctx sdk.Context) ([]types.AccountState, error) {
	iterator := storetypes.KVStorePrefixIterator(
		ctx.KVStore(k.storeKey),
		[]byte(types.AccountStateKeyPrefix),
	)
	defer iterator.Close()

	accounts := []types.AccountState{}
	for ; iterator.Valid(); iterator.Next() {
		var accountState types.AccountState
		err := k.cdc.Unmarshal(iterator.Value(), &accountState)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal account state: %w", err)
		}

		accounts = append(accounts, accountState)
	}

	return accounts, nil
}

// Set genesis state
func (k Keeper) SetGenesisState(ctx sdk.Context, data types.GenesisState) error {
	for _, account := range data.Accounts {
		address, err := sdk.AccAddressFromBech32(account.Address)
		if err != nil {
			return err
		}
		k.SetAccountState(ctx, address, account)
	}

	k.SetParams(ctx, data.Params)
	k.SetNextAuthenticatorId(ctx, data.NextAuthenticatorId)

	for _, data := range data.GetAuthenticatorData() {
		address := data.GetAddress()
		for _, authenticator := range data.GetAuthenticators() {
			k.SetAuthenticator(ctx, address, authenticator.Id, authenticator)
		}
	}

	return nil
}

// GetAllAuthenticatorData is used in genesis export to export all the authenticator for all accounts
func (k Keeper) GetAllAuthenticatorData(ctx sdk.Context) ([]types.AuthenticatorData, error) {
	var accountAuthenticators []types.AuthenticatorData

	parse := func(key []byte, value []byte) error {
		var authenticator types.AccountAuthenticator
		err := k.cdc.Unmarshal(value, &authenticator)
		if err != nil {
			return err
		}

		// Key is of format `SA/A/{accountAddr}/{authenticatorId}`.
		// Extract account address from key.
		accountAddr := strings.Split(string(key), "/")[2]

		// Check if this entry is for a new address or the same as the last one processed
		if len(accountAuthenticators) == 0 ||
			accountAuthenticators[len(accountAuthenticators)-1].Address != accountAddr {
			// If it's a new address, create a new AuthenticatorData entry
			accountAuthenticators = append(accountAuthenticators, types.AuthenticatorData{
				Address:        accountAddr,
				Authenticators: []types.AccountAuthenticator{authenticator},
			})
		} else {
			// If it's the same address, append the authenticator to the last entry in the list
			lastIndex := len(accountAuthenticators) - 1
			accountAuthenticators[lastIndex].Authenticators = append(
				accountAuthenticators[lastIndex].Authenticators,
				authenticator,
			)
		}

		return nil
	}

	// Iterate over all entries in the store using a prefix iterator
	iterator := storetypes.KVStorePrefixIterator(ctx.KVStore(k.storeKey), []byte(types.AuthenticatorKeyPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		err := parse(iterator.Key(), iterator.Value())
		if err != nil {
			return nil, err
		}
	}

	return accountAuthenticators, nil
}

func AccountStateFromTimestampNonceDetails(
	address sdk.AccAddress,
	tsNonce uint64,
) types.AccountState {
	return types.AccountState{
		Address: address.String(),
		TimestampNonceDetails: types.TimestampNonceDetails{
			MaxEjectedNonce: TimestampNonceSequenceCutoff,
			TimestampNonces: []uint64{tsNonce},
		},
	}
}

// Get the AccountState from KVStore for a given account address
func (k Keeper) GetAccountState(
	ctx sdk.Context,
	address sdk.AccAddress,
) (types.AccountState, bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AccountStateKeyPrefix))
	bz := prefixStore.Get(address.Bytes())
	if bz == nil {
		return types.AccountState{}, false
	}

	var accountState types.AccountState
	k.cdc.MustUnmarshal(bz, &accountState)

	// By default empty slices are Unmarshed into nil
	if accountState.TimestampNonceDetails.TimestampNonces == nil {
		accountState.TimestampNonceDetails.TimestampNonces = make([]uint64, 0)
	}

	return accountState, true
}

// Set the AccountState into KVStore for a given account address
func (k Keeper) SetAccountState(
	ctx sdk.Context,
	address sdk.AccAddress,
	accountState types.AccountState,
) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AccountStateKeyPrefix))
	bz := k.cdc.MustMarshal(&accountState)
	prefixStore.Set(address.Bytes(), bz)
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}
