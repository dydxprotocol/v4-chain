package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: key,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(log.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

// Func required for test setup
func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}

// Get all account details pairs in store
func (k Keeper) GetAllAccountStates(ctx sdk.Context) []*types.AccountState {
	store := ctx.KVStore(k.storeKey)

	iterator := storetypes.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	var accounts []*types.AccountState
	for ; iterator.Valid(); iterator.Next() {
		var account types.AccountState
		k.cdc.MustUnmarshal(iterator.Value(), &account)
		accounts = append(accounts, &account)
	}

	return accounts
}

// Set genesis state
func (k Keeper) SetGenesisState(ctx sdk.Context, data types.GenesisState) error {
	store := ctx.KVStore(k.storeKey)

	for _, account := range data.Accounts {
		address, err := sdk.AccAddressFromBech32(account.Address)
		if err != nil {
			return err
		}
		k.setAccountState(store, address, *account)
	}

	return nil
}

func (k Keeper) InitializeAccount(ctx sdk.Context, address sdk.AccAddress) (types.AccountState, error) {
	store := ctx.KVStore(k.storeKey)

	if _, found := k.getAccountState(store, address); found {
		return types.AccountState{}, errors.New(
			"Cannot initialize AccountState for address with existing AccountState, address: " + address.String(),
		)
	}

	initialAccountState := types.AccountState{
		Address:               address.String(),
		TimestampNonceDetails: DeepCopyTimestampNonceDetails(InitialTimestampNonceDetails),
	}

	k.setAccountState(store, address, initialAccountState)

	return initialAccountState, nil
}

// Get the AccountState from KVStore for a given account address
func (k Keeper) getAccountState(
	store storetypes.KVStore,
	address sdk.AccAddress,
) (types.AccountState, bool) {
	bz := store.Get(address.Bytes())
	if bz == nil {
		return types.AccountState{}, false
	}

	var accountState types.AccountState
	k.cdc.MustUnmarshal(bz, &accountState)
	return accountState, true
}

// Set the AccountState into KVStore for a given account address
func (k Keeper) setAccountState(
	store storetypes.KVStore,
	address sdk.AccAddress,
	accountState types.AccountState,
) {
	bz := k.cdc.MustMarshal(&accountState)
	store.Set(address.Bytes(), bz)
}
