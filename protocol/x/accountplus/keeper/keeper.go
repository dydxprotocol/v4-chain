package keeper

import (
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
func (k Keeper) GetAllAccoutDetails(ctx sdk.Context) []*types.AccountState {
	store := ctx.KVStore(k.storeKey)

	iterator := storetypes.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	var accounts []*types.AccountState
	for ; iterator.Valid(); iterator.Next() {
		var details types.TimestampNonceDetails
		k.cdc.MustUnmarshal(iterator.Value(), &details)
		address := sdk.AccAddress(iterator.Key()).String()
		accounts = append(accounts, &types.AccountState{
			Address: address,
			Details: &details,
		})
	}

	return accounts
}

// Set genesis state
func (k Keeper) SetGenesisState(ctx sdk.Context, data types.GenesisState) {
	store := ctx.KVStore(k.storeKey)

	for _, account := range data.Accounts {
		address, err := sdk.AccAddressFromBech32(account.Address)
		if err != nil {
			panic(err)
		}
		k.setTimestampNonceDetails(store, address, *account.Details)
	}
}

// Set the TimestampNonceDetails into KVStore for a given account address
func (k Keeper) setTimestampNonceDetails(store storetypes.KVStore, address sdk.AccAddress, details types.TimestampNonceDetails) {
	bz := k.cdc.MustMarshal(&details)
	store.Set(address.Bytes(), bz)
}

// Get the TimestampNonceDetails from KVStore for a given account address
func (k Keeper) getTimestampNonceDetails(store storetypes.KVStore, address sdk.AccAddress) (types.TimestampNonceDetails, bool) {
	bz := store.Get(address.Bytes())
	if bz == nil {
		return types.TimestampNonceDetails{}, false
	}

	var details types.TimestampNonceDetails
	k.cdc.MustUnmarshal(bz, &details)
	return details, true
}
