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
			return errorsmod.Wrapf(types.ErrUpdatingAffiliateReferredVolume,
				"affiliate %s, error: %s", affiliateAddr, err)
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
		// Return empty tiers if not initialized.
		return types.AffiliateTiers{}, nil
	}
	err := k.cdc.Unmarshal(affiliateTiersBytes, &affiliateTiers)
	if err != nil {
		return affiliateTiers, err
	}

	return affiliateTiers, nil
}

// GetTakerFeeShare returns the taker fee share for an address based on the affiliate tiers.
// If the address is in the whitelist, the fee share ppm is overridden.
func (k Keeper) GetTakerFeeShare(
	ctx sdk.Context,
	address string,
	affiliatesWhitelistMap map[string]uint32,
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
	// Override fee share ppm if the address is in the whitelist.
	if _, exists := affiliatesWhitelistMap[affiliateAddress]; exists {
		feeSharePpm = affiliatesWhitelistMap[affiliateAddress]
		return affiliateAddress, feeSharePpm, true, nil
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
	// Return 0 tier if no tiers are set.
	if len(tiers) == 0 {
		return 0, 0, nil
	}
	numTiers := uint32(len(tiers))
	maxTierLevel := numTiers - 1
	currentTier := uint32(0)
	referredVolume, err := k.GetReferredVolume(ctx, affiliateAddr)
	if err != nil {
		return 0, 0, err
	}

	for index, tier := range tiers {
		// required referred volume is strictly increasing as tiers are traversed in order.
		if referredVolume.Cmp(lib.BigU(tier.ReqReferredVolumeQuoteQuantums)) < 0 {
			break
		}
		// safe to do as tier cannot be negative
		currentTier = uint32(index)
	}

	if currentTier == maxTierLevel {
		return currentTier, tiers[currentTier].TakerFeeSharePpm, nil
	}

	numCoinsStaked := k.statsKeeper.GetStakedAmount(ctx, affiliateAddr)
	for i := currentTier + 1; i < numTiers; i++ {
		// required staked coins is strictly increasing as tiers are traversed in order.
		expMultiplier, _ := lib.BigPow10(-lib.BaseDenomExponent)
		reqStakedCoins := new(big.Int).Mul(
			lib.BigU(tiers[i].ReqStakedWholeCoins),
			expMultiplier,
		)
		if numCoinsStaked.Cmp(reqStakedCoins) < 0 {
			break
		}
		currentTier = i
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
		if tiers[i].ReqReferredVolumeQuoteQuantums <= tiers[i-1].ReqReferredVolumeQuoteQuantums ||
			tiers[i].ReqStakedWholeCoins <= tiers[i-1].ReqStakedWholeCoins {
			return errorsmod.Wrapf(types.ErrInvalidAffiliateTiers,
				"tiers values must be strictly increasing")
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

func (k Keeper) AggregateAffiliateReferredVolumeForFills(
	ctx sdk.Context,
) error {
	blockStats := k.statsKeeper.GetBlockStats(ctx)
	referredByCache := make(map[string]string)

	for _, fill := range blockStats.Fills {
		// Process taker's referred volume
		referredByAddrTaker, cached := referredByCache[fill.Taker]
		if !cached {
			var found bool
			referredByAddrTaker, found = k.GetReferredBy(ctx, fill.Taker)
			if found {
				referredByCache[fill.Taker] = referredByAddrTaker
			}
		}
		if referredByAddrTaker != "" {
			if err := k.AddReferredVolume(ctx, referredByAddrTaker, lib.BigU(fill.Notional)); err != nil {
				return err
			}
		}

		// Process maker's referred volume
		referredByAddrMaker, cached := referredByCache[fill.Maker]
		if !cached {
			var found bool
			referredByAddrMaker, found = k.GetReferredBy(ctx, fill.Maker)
			if found {
				referredByCache[fill.Maker] = referredByAddrMaker
			}
		}
		if referredByAddrMaker != "" {
			if err := k.AddReferredVolume(ctx, referredByAddrMaker, lib.BigU(fill.Notional)); err != nil {
				return err
			}
		}
	}
	return nil
}
