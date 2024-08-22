package keeper

import (
	"fmt"
	"math/big"

	"errors"

	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
)

type (
	Keeper struct {
		cdc         codec.BinaryCodec
		storeKey    storetypes.StoreKey
		authorities map[string]struct{}
		statsKeeper types.StatsKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	authorities []string,
	statsKeeper types.StatsKeeper,
) *Keeper {
	return &Keeper{
		cdc:         cdc,
		storeKey:    storeKey,
		authorities: lib.UniqueSliceToSet(authorities),
		statsKeeper: statsKeeper,
	}
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(log.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {}

func (k Keeper) RegisterAffiliate(
	ctx sdk.Context,
	referee string,
	affiliateAddr string,
) error {
	referredByPrefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ReferredByKeyPrefix))
	if referredByPrefixStore.Has([]byte(referee)) {
		return errors.New("Referee already exists for address: " + referee)
	}
	referredByPrefixStore.Set([]byte(referee), []byte(affiliateAddr))
	return nil
}

func (k Keeper) GetReferredBy(ctx sdk.Context, referee string) (string, bool) {
	referredByPrefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ReferredByKeyPrefix))
	if !referredByPrefixStore.Has([]byte(referee)) {
		return "", false
	}
	return string(referredByPrefixStore.Get([]byte(referee))), true
}

func (k Keeper) AddReferredVolume(
	ctx sdk.Context,
	affiliateAddr string,
	referredVolumeFromBlock dtypes.SerializableInt,
) error {
	affiliateReferredVolumePrefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ReferredVolumeKeyPrefix))
	var referredVolume dtypes.SerializableInt
	if !affiliateReferredVolumePrefixStore.Has([]byte(affiliateAddr)) {
		referredVolume = dtypes.NewInt(0)
	} else {
		err := referredVolume.Unmarshal(affiliateReferredVolumePrefixStore.Get([]byte(affiliateAddr)))
		if err != nil {
			return err
		}
	}
	referredVolumeBigInt := referredVolume.BigInt()
	referredVolumeBigInt.Add(referredVolumeBigInt, referredVolumeFromBlock.BigInt())
	updatedReferedVolume := dtypes.NewIntFromBigInt(referredVolumeBigInt)

	updatedReferredVolumeBytes, err := updatedReferedVolume.Marshal()
	if err != nil {
		return errors.New("Error marshalling referred volume")
	}
	affiliateReferredVolumePrefixStore.Set([]byte(affiliateAddr), updatedReferredVolumeBytes)
	return nil
}

func (k Keeper) GetReferredVolume(ctx sdk.Context, affiliateAddr string) (dtypes.SerializableInt, bool) {
	affiliateReferredVolumePrefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ReferredVolumeKeyPrefix))
	if !affiliateReferredVolumePrefixStore.Has([]byte(affiliateAddr)) {
		return dtypes.NewInt(0), false
	}
	var referredVolume dtypes.SerializableInt
	err := referredVolume.Unmarshal(affiliateReferredVolumePrefixStore.Get([]byte(affiliateAddr)))
	if err != nil {
		return dtypes.NewInt(0), false
	}
	return referredVolume, true
}

func (k Keeper) GetAllAffiliateTiers(ctx sdk.Context) (types.AffiliateTiers, error) {
	store := ctx.KVStore(k.storeKey)
	affiliateTiersBytes := store.Get([]byte(types.AffiliateTiersKey))

	var affiliateTiers types.AffiliateTiers
	if affiliateTiersBytes == nil {
		return affiliateTiers, errors.New("Affiliate tiers not found")
	}
	err := k.cdc.Unmarshal(affiliateTiersBytes, &affiliateTiers)
	if err != nil {
		return affiliateTiers, err
	}

	return affiliateTiers, nil
}

func (k Keeper) GetTakerFeeShare(
	ctx sdk.Context,
	address string,
) (
	affiliateAddress string,
	feeSharePpm uint32,
	exists bool,
	err error,
) {
	affiliateAddress, exists = k.GetReferredBy(ctx, address)
	if !exists {
		return "", 0, false, nil
	}
	_, feeSharePpm, err = k.GetTierForAffiliate(ctx, affiliateAddress)
	if err != nil {
		return "", 0, false, err
	}
	return affiliateAddress, feeSharePpm, true, nil
}

// Assumes that the affiliate tiers are sorted by level in ascending order.
func (k Keeper) GetTierForAffiliate(
	ctx sdk.Context,
	affiliateAddr string,
) (
	tierLevel uint32,
	feeSharePpm uint32,
	err error) {
	affiliateTiers, err := k.GetAllAffiliateTiers(ctx)
	if err != nil {
		return 0, 0, err
	}
	numTiers := uint32(len(affiliateTiers.GetTiers()))
	maxTierLevel := numTiers - 1
	currentTier := uint32(0)
	referredVolume, exists := k.GetReferredVolume(ctx, affiliateAddr)

	if !exists {
		// If referred volume is not found, set it to 0.
		referredVolume = dtypes.NewInt(0)
	}

	for _, tier := range affiliateTiers.GetTiers() {
		if referredVolume.BigInt().Int64() >= int64(tier.ReqReferredVolume) {
			currentTier = tier.GetLevel()
		}
	}

	if currentTier == maxTierLevel {
		return currentTier, affiliateTiers.GetTiers()[currentTier].TakerFeeSharePpm, nil
	}
	numCoinsStaked := k.statsKeeper.GetStakedAmount(ctx, affiliateAddr)
	for _, tier := range affiliateTiers.GetTiers() {
		if numCoinsStaked.Cmp(big.NewInt(
			int64(affiliateTiers.GetTiers()[currentTier].ReqStakedWholeCoins))) >= 0 &&
			tier.GetLevel() > currentTier {
			currentTier = tier.GetLevel()
		}
	}
	return currentTier, affiliateTiers.GetTiers()[currentTier].TakerFeeSharePpm, nil
}

func (k Keeper) UpdateAffiliateTiers(ctx sdk.Context, affiliateTiers types.AffiliateTiers) error {
	numTiers := uint32(len(affiliateTiers.GetTiers()))
	for i := uint32(0); i < numTiers; i++ {
		if affiliateTiers.GetTiers()[i].GetLevel() != i {
			return errors.New("tiers are not sorted by level in ascending order")
		}
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.AffiliateTiersKey), k.cdc.MustMarshal(&affiliateTiers))
	return nil
}
