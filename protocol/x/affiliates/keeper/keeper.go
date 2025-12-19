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
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
)

type (
	Keeper struct {
		cdc                 codec.BinaryCodec
		storeKey            storetypes.StoreKey
		authorities         map[string]struct{}
		statsKeeper         types.StatsKeeper
		feetiersKeeper      types.FeetiersKeeper
		indexerEventManager indexer_manager.IndexerEventManager
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	authorities []string,
	statsKeeper types.StatsKeeper,
	indexerEventManager indexer_manager.IndexerEventManager,
) *Keeper {
	return &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		authorities:         lib.UniqueSliceToSet(authorities),
		statsKeeper:         statsKeeper,
		indexerEventManager: indexerEventManager,
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
	if referee == affiliateAddr {
		return errorsmod.Wrapf(
			types.ErrSelfReferral, "referee: %s, affiliate: %s",
			referee, affiliateAddr,
		)
	}
	if _, err := sdk.AccAddressFromBech32(referee); err != nil {
		return errorsmod.Wrapf(types.ErrInvalidAddress, "referee: %s", referee)
	}
	if _, err := sdk.AccAddressFromBech32(affiliateAddr); err != nil {
		return errorsmod.Wrapf(types.ErrInvalidAddress, "affiliate: %s", affiliateAddr)
	}
	if _, found := k.GetReferredBy(ctx, referee); found {
		return errorsmod.Wrapf(types.ErrAffiliateAlreadyExistsForReferee, "referee: %s, affiliate: %s",
			referee, affiliateAddr)
	}
	affiliateTiers, err := k.GetAllAffiliateTiers(ctx)
	if err != nil {
		return err
	}
	// Return error if no tiers are set.
	if len(affiliateTiers.GetTiers()) == 0 {
		return types.ErrAffiliateTiersNotSet
	}
	prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ReferredByKeyPrefix)).Set([]byte(referee), []byte(affiliateAddr))
	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeRegisterAffiliate,
		indexerevents.RegisterAffiliateEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewRegisterAffiliateEventV1(
				referee,
				affiliateAddr,
			),
		),
	)
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

// GetAllAffiliateTiers returns all affiliate tiers.
func (k Keeper) GetAllAffiliateTiers(ctx sdk.Context) (types.AffiliateTiers, error) {
	store := ctx.KVStore(k.storeKey)
	affiliateTiersBytes := store.Get([]byte(types.AffiliateTiersKey))

	var affiliateTiers types.AffiliateTiers
	if affiliateTiersBytes == nil {
		// Return empty tiers if not initialized.
		return types.AffiliateTiers{}, nil
	}
	err := k.cdc.Unmarshal(affiliateTiersBytes, &affiliateTiers)
	if err != nil {
		return affiliateTiers, err
	}

	return affiliateTiers, nil
}

func (k Keeper) GetAllAffilliateOverrides(ctx sdk.Context) (types.AffiliateOverrides, error) {
	store := ctx.KVStore(k.storeKey)
	affiliateOverridesBytes := store.Get([]byte(types.AffiliateOverridesKey))

	var affiliateOverrides types.AffiliateOverrides
	if affiliateOverridesBytes == nil {
		// Return empty overrides if not initialized.
		return types.AffiliateOverrides{}, nil
	}
	err := k.cdc.Unmarshal(affiliateOverridesBytes, &affiliateOverrides)
	if err != nil {
		return affiliateOverrides, err
	}

	return affiliateOverrides, nil
}

// GetTakerFeeShare returns the taker fee share for an address based on the affiliate tiers.
// If the address is in the whitelist, the fee share ppm is overridden.
func (k Keeper) GetTakerFeeShare(
	ctx sdk.Context,
	address string,
	affiliateOverrides map[string]bool,
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
	_, feeSharePpm, err = k.GetTierForAffiliate(ctx, affiliateAddress, affiliateOverrides)
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
	affiliateOverrides map[string]bool,
) (
	tierLevel uint32,
	feeSharePpm uint32,
	err error) {
	affiliateTiers, err := k.GetAllAffiliateTiers(ctx)
	if err != nil {
		return 0, 0, err
	}

	tiers := affiliateTiers.GetTiers()
	// Return 0 tier if no tiers are set.
	if len(tiers) == 0 {
		return 0, 0, nil
	}
	numTiers := uint32(len(tiers))
	maxTierLevel := numTiers - 1
	currentTier := uint32(0)

	// Check whether the address is overridden, if it is then set the
	// affiliate tier to the max
	if affiliateOverrides != nil {
		if _, exists := affiliateOverrides[affiliateAddr]; exists {
			feeSharePpm = affiliateTiers.Tiers[maxTierLevel].TakerFeeSharePpm
			return uint32(maxTierLevel), feeSharePpm, nil
		}
	}

	// Get the affiliate revenue generated in the last 30d
	userStats := k.statsKeeper.GetUserStats(ctx, affiliateAddr)
	referredVolume := big.NewInt(0)
	if userStats != nil {
		referredVolume = new(big.Int).SetUint64(userStats.Affiliate_30DReferredVolumeQuoteQuantums)
	}

	for index, tier := range tiers {
		// required referred volume is strictly increasing as tiers are traversed in order.
		if referredVolume.Cmp(lib.BigU(tier.ReqReferredVolumeQuoteQuantums)) < 0 {
			break
		}
		// safe to do as tier cannot be negative
		currentTier = uint32(index)
	}

	return currentTier, tiers[currentTier].TakerFeeSharePpm, nil
}

// UpdateAffiliateTiers updates the affiliate tiers.
// Used primarily through governance.
func (k Keeper) UpdateAffiliateTiers(ctx sdk.Context, affiliateTiers types.AffiliateTiers) error {
	store := ctx.KVStore(k.storeKey)
	affiliateTiersBytes := k.cdc.MustMarshal(&affiliateTiers)
	tiers := affiliateTiers.GetTiers()
	// start at 1, since 0 is the default tier.
	for i := 1; i < len(tiers); i++ {
		// Check if the taker fee share ppm is greater than the cap.
		if tiers[i].TakerFeeSharePpm > types.AffiliatesRevSharePpmCap {
			return errorsmod.Wrapf(types.ErrRevShareSafetyViolation,
				"taker fee share ppm %d is greater than the cap %d",
				tiers[i].TakerFeeSharePpm, types.AffiliatesRevSharePpmCap)
		}
		// Check if the tiers are strictly increasing.
		if tiers[i].ReqReferredVolumeQuoteQuantums <= tiers[i-1].ReqReferredVolumeQuoteQuantums {
			return errorsmod.Wrapf(types.ErrInvalidAffiliateTiers,
				"volume must be strictly increasing")
		}
	}
	store.Set([]byte(types.AffiliateTiersKey), affiliateTiersBytes)
	return nil
}

func (k *Keeper) SetFeetiersKeeper(feetiersKeeper types.FeetiersKeeper) {
	k.feetiersKeeper = feetiersKeeper
}

func (k Keeper) GetIndexerEventManager() indexer_manager.IndexerEventManager {
	return k.indexerEventManager
}

// Deprecated: This is deprecated in favor of AffiliateOverride.
func (k Keeper) GetAffiliateWhitelistMap(ctx sdk.Context) (map[string]uint32, error) {
	affiliateWhitelist, err := k.GetAffiliateWhitelist(ctx)
	if err != nil {
		return nil, err
	}
	affiliateWhitelistMap := make(map[string]uint32)
	for _, tier := range affiliateWhitelist.GetTiers() {
		for _, address := range tier.GetAddresses() {
			affiliateWhitelistMap[address] = tier.GetTakerFeeSharePpm()
		}
	}
	return affiliateWhitelistMap, nil
}

// Deprecated: This is deprecated in favor of AffiliateOverride.
func (k Keeper) SetAffiliateWhitelist(ctx sdk.Context, whitelist types.AffiliateWhitelist) error {
	store := ctx.KVStore(k.storeKey)
	addressSet := make(map[string]bool)
	for _, tier := range whitelist.Tiers {
		// Check if the taker fee share ppm is greater than the cap.
		if tier.TakerFeeSharePpm > types.AffiliatesRevSharePpmCap {
			return errorsmod.Wrapf(types.ErrRevShareSafetyViolation,
				"taker fee share ppm %d is greater than the cap %d",
				tier.TakerFeeSharePpm, types.AffiliatesRevSharePpmCap)
		}
		for _, address := range tier.Addresses {
			// Check for invalid addresses.
			if _, err := sdk.AccAddressFromBech32(address); err != nil {
				return errorsmod.Wrapf(types.ErrInvalidAddress, "address to whitelist: %s", address)
			}
			// Check for duplicate addresses.
			if addressSet[address] {
				return errorsmod.Wrapf(types.ErrDuplicateAffiliateAddressForWhitelist,
					"address %s is duplicated in affiliate whitelist", address)
			}
			addressSet[address] = true
		}
	}
	affiliateWhitelistBytes := k.cdc.MustMarshal(&whitelist)
	store.Set([]byte(types.AffiliateWhitelistKey), affiliateWhitelistBytes)
	return nil
}

// DO NOT USE: This will be deprecated soon.
func (k Keeper) GetAffiliateWhitelist(ctx sdk.Context) (types.AffiliateWhitelist, error) {
	store := ctx.KVStore(k.storeKey)
	affiliateWhitelistBytes := store.Get([]byte(types.AffiliateWhitelistKey))
	if affiliateWhitelistBytes == nil {
		return types.AffiliateWhitelist{
			Tiers: []types.AffiliateWhitelist_Tier{},
		}, nil
	}
	affiliateWhitelist := types.AffiliateWhitelist{}
	err := k.cdc.Unmarshal(affiliateWhitelistBytes, &affiliateWhitelist)
	if err != nil {
		return types.AffiliateWhitelist{}, err
	}
	return affiliateWhitelist, nil
}

func (k Keeper) UpdateAffiliateParameters(
	ctx sdk.Context,
	msg *types.MsgUpdateAffiliateParameters,
) error {
	store := ctx.KVStore(k.storeKey)

	affiliateParametersBytes, err := k.cdc.Marshal(&msg.AffiliateParameters)
	if err != nil {
		return err
	}
	store.Set([]byte(types.AffiliateParametersKey), affiliateParametersBytes)

	return nil
}

func (k Keeper) GetAffiliateParameters(ctx sdk.Context) (types.AffiliateParameters, error) {
	store := ctx.KVStore(k.storeKey)
	affiliateParametersBytes := store.Get([]byte(types.AffiliateParametersKey))
	if affiliateParametersBytes == nil {
		return types.AffiliateParameters{}, nil
	}
	affiliateParameters := types.AffiliateParameters{}
	err := k.cdc.Unmarshal(affiliateParametersBytes, &affiliateParameters)
	if err != nil {
		return types.AffiliateParameters{}, err
	}
	return affiliateParameters, nil
}

func (k Keeper) SetAffiliateOverrides(ctx sdk.Context, overrides types.AffiliateOverrides) error {
	store := ctx.KVStore(k.storeKey)
	affiliateOverridesBytes, err := k.cdc.Marshal(&overrides)
	if err != nil {
		return err
	}
	store.Set([]byte(types.AffiliateOverridesKey), affiliateOverridesBytes)
	return nil
}

func (k Keeper) GetAffiliateOverrides(ctx sdk.Context) (types.AffiliateOverrides, error) {
	store := ctx.KVStore(k.storeKey)
	affiliateOverridesBytes := store.Get([]byte(types.AffiliateOverridesKey))
	if affiliateOverridesBytes == nil {
		return types.AffiliateOverrides{}, nil
	}
	affiliateOverrides := types.AffiliateOverrides{}
	err := k.cdc.Unmarshal(affiliateOverridesBytes, &affiliateOverrides)
	if err != nil {
		return types.AffiliateOverrides{}, err
	}
	return affiliateOverrides, nil
}

func (k Keeper) GetAffiliateOverridesMap(ctx sdk.Context) (map[string]bool, error) {
	affiliateOverrides, err := k.GetAffiliateOverrides(ctx)
	if err != nil {
		return nil, err
	}
	affiliateOverridesMap := make(map[string]bool)
	for _, address := range affiliateOverrides.Addresses {
		affiliateOverridesMap[address] = true
	}
	return affiliateOverridesMap, nil
}
