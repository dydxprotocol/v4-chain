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

func DefaultAccountState(address sdk.AccAddress) types.AccountState {
	return types.AccountState{
		Address: address.String(),
		TimestampNonceDetails: types.TimestampNonceDetails{
			MaxEjectedNonce: 0,
			TimestampNonces: []uint64{},
		},
	}
}

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
func (k Keeper) GetAllAccountStates(ctx sdk.Context) ([]types.AccountState, error) {
	store := ctx.KVStore(k.storeKey)

	iterator := storetypes.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	accounts := []types.AccountState{}
	for ; iterator.Valid(); iterator.Next() {
		accountState, found := k.GetAccountState(ctx, iterator.Key())
		if !found {
			return accounts, errors.New("Could not get account state for address: " + sdk.AccAddress(iterator.Key()).String())
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
		k.setAccountState(ctx, address, account)
	}

	return nil
}

func (k Keeper) InitializeAccount(ctx sdk.Context, address sdk.AccAddress) error {
	if _, found := k.GetAccountState(ctx, address); found {
		return errors.New(
			"Cannot initialize AccountState for address with existing AccountState, address: " + address.String(),
		)
	}

	k.setAccountState(ctx, address, DefaultAccountState(address))

	return nil
}

// Get the AccountState from KVStore for a given account address
func (k Keeper) GetAccountState(
	ctx sdk.Context,
	address sdk.AccAddress,
) (types.AccountState, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(address.Bytes())
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
func (k Keeper) setAccountState(
	ctx sdk.Context,
	address sdk.AccAddress,
	accountState types.AccountState,
) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&accountState)
	store.Set(address.Bytes(), bz)
}
