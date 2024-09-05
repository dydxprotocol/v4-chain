package keeper

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
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

// RegisterAffiliate registers an affiliate address for a referee address.
func (k Keeper) RegisterAffiliate(
	ctx sdk.Context,
	referee string,
	affiliateAddr string,
) error {
	if _, found := k.GetReferredBy(ctx, referee); found {
		return errorsmod.Wrapf(types.ErrAffiliateAlreadyExistsForReferee, "referee: %s, affiliate: %s",
			referee, affiliateAddr)
	}
	if _, err := sdk.AccAddressFromBech32(referee); err != nil {
		return errorsmod.Wrapf(types.ErrInvalidAddress, "referee: %s", referee)
	}
	if _, err := sdk.AccAddressFromBech32(affiliateAddr); err != nil {
		return errorsmod.Wrapf(types.ErrInvalidAddress, "affiliate: %s", affiliateAddr)
	}
	prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ReferredByKeyPrefix)).Set([]byte(referee), []byte(affiliateAddr))
	// TODO(OTE-696): Emit indexer event.
	return nil
}

// GetReferredBy returns the affiliate address for a referee address.
func (k Keeper) GetReferredBy(ctx sdk.Context, referee string) (string, bool) {
	referredByPrefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ReferredByKeyPrefix))
	if !referredByPrefixStore.Has([]byte(referee)) {
		return "", false
	}
	return string(referredByPrefixStore.Get([]byte(referee))), true
}

// AddReferredVolume adds the referred volume from a block to the affiliate's referred volume.
func (k Keeper) AddReferredVolume(
	ctx sdk.Context,
	affiliateAddr string,
	referredVolumeFromBlock *big.Int,
) error {
	affiliateReferredVolumePrefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ReferredVolumeKeyPrefix))
	referredVolume := big.NewInt(0)

	if affiliateReferredVolumePrefixStore.Has([]byte(affiliateAddr)) {
		prevReferredVolumeFromState := dtypes.SerializableInt{}
		if err := prevReferredVolumeFromState.Unmarshal(
			affiliateReferredVolumePrefixStore.Get([]byte(affiliateAddr)),
		); err != nil {
			// maybe change to errorsmod
			return err
		}
		referredVolume = prevReferredVolumeFromState.BigInt()
	}

	referredVolume.Add(
		referredVolume,
		referredVolumeFromBlock,
	)
	updatedReferedVolume := dtypes.NewIntFromBigInt(referredVolume)

	updatedReferredVolumeBytes, err := updatedReferedVolume.Marshal()
	if err != nil {
		return errorsmod.Wrapf(types.ErrUpdatingAffiliateReferredVolume,
			"affiliate %s, error: %s", affiliateAddr, err)
	}
	affiliateReferredVolumePrefixStore.Set([]byte(affiliateAddr), updatedReferredVolumeBytes)
	return nil
}

// GetReferredVolume returns all time referred volume for an affiliate address.
func (k Keeper) GetReferredVolume(ctx sdk.Context, affiliateAddr string) (*big.Int, error) {
	affiliateReferredVolumePrefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ReferredVolumeKeyPrefix))
	if !affiliateReferredVolumePrefixStore.Has([]byte(affiliateAddr)) {
		return big.NewInt(0), nil
	}
	var referredVolume dtypes.SerializableInt
	if err := referredVolume.Unmarshal(affiliateReferredVolumePrefixStore.Get([]byte(affiliateAddr))); err != nil {
		return big.NewInt(0), err
	}
	return referredVolume.BigInt(), nil
}

// GetAllAffiliateTiers returns all affiliate tiers.
func (k Keeper) GetAllAffiliateTiers(ctx sdk.Context) (types.AffiliateTiers, error) {
	store := ctx.KVStore(k.storeKey)
	affiliateTiersBytes := store.Get([]byte(types.AffiliateTiersKey))

	var affiliateTiers types.AffiliateTiers
	if affiliateTiersBytes == nil {
		return affiliateTiers, errorsmod.Wrapf(types.ErrAffiliateTiersNotInitialized, "affiliate tiers not initialized")
	}
	err := k.cdc.Unmarshal(affiliateTiersBytes, &affiliateTiers)
	if err != nil {
		return affiliateTiers, err
	}

	return affiliateTiers, nil
}

// GetTakerFeeShare returns the taker fee share for an address.
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

// GetTierForAffiliate returns the tier an affiliate address is qualified for.
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
	tiers := affiliateTiers.GetTiers()
	numTiers := uint32(len(tiers))
	maxTierLevel := numTiers - 1
	currentTier := uint32(0)
	referredVolume, err := k.GetReferredVolume(ctx, affiliateAddr)
	if err != nil {
		return 0, 0, err
	}

	for index, tier := range tiers {
		if referredVolume.Cmp(lib.BigU(tier.ReqReferredVolumeQuoteQuantums)) >= 0 {
			// safe to do as tier cannot be negative
			currentTier = uint32(index)
		} else {
			break
		}
	}

	if currentTier == maxTierLevel {
		return currentTier, tiers[currentTier].TakerFeeSharePpm, nil
	}

	numCoinsStaked := k.statsKeeper.GetStakedAmount(ctx, affiliateAddr)
	for i := currentTier + 1; i < numTiers; i++ {
		expMultiplier, _ := lib.BigPow10(lib.BaseDenomExponent)
		reqStakedCoins := new(big.Int).Mul(
			lib.BigU(tiers[i].ReqStakedWholeCoins),
			expMultiplier,
		)
		if numCoinsStaked.Cmp(reqStakedCoins) >= 0 {
			currentTier = i
		} else {
			break
		}
	}
	return currentTier, tiers[currentTier].TakerFeeSharePpm, nil
}

// UpdateAffiliateTiers updates the affiliate tiers.
// Used primarily through governance.
func (k Keeper) UpdateAffiliateTiers(ctx sdk.Context, affiliateTiers types.AffiliateTiers) {
	store := ctx.KVStore(k.storeKey)
	// TODO(OTE-779): Check strictly increasing volume and
	// staking requirements hold in UpdateAffiliateTiers
	store.Set([]byte(types.AffiliateTiersKey), k.cdc.MustMarshal(&affiliateTiers))
}
