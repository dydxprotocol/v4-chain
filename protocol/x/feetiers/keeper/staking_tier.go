package keeper

import (
	"fmt"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
)

// SetStakingTiers sets staking tiers in state
func (k Keeper) SetStakingTiers(ctx sdk.Context, stakingTiers []*types.StakingTier) error {
	// Validate that staking tiers correspond to existing fee tiers
	perpetualFeeParams := k.GetPerpetualFeeParams(ctx)
	existingTiers := make(map[string]bool)
	for _, tier := range perpetualFeeParams.Tiers {
		existingTiers[tier.Name] = true
	}

	for _, stakingTier := range stakingTiers {
		if !existingTiers[stakingTier.FeeTierName] {
			return fmt.Errorf("fee tier %s does not exist", stakingTier.FeeTierName)
		}
	}

	// Clear existing staking tiers
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, []byte(types.StakingTierKeyPrefix))

	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		prefixStore.Delete(iterator.Key())
	}

	// Set new staking tiers
	for _, tier := range stakingTiers {
		bz := k.cdc.MustMarshal(tier)
		store.Set(types.StakingTierKey(tier.FeeTierName), bz)
	}

	return nil
}

// GetStakingTier retrieves a staking tier from state, if exists
func (k Keeper) GetStakingTier(ctx sdk.Context, tierName string) (*types.StakingTier, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.StakingTierKey(tierName))
	if bz == nil {
		return nil, false
	}

	var tier types.StakingTier
	k.cdc.MustUnmarshal(bz, &tier)
	return &tier, true
}

// GetAllStakingTiers retrieves all staking tiers from state
func (k Keeper) GetAllStakingTiers(ctx sdk.Context) []*types.StakingTier {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, []byte(types.StakingTierKeyPrefix))

	var stakingTiers []*types.StakingTier
	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var tier types.StakingTier
		k.cdc.MustUnmarshal(iterator.Value(), &tier)
		stakingTiers = append(stakingTiers, &tier)
	}

	return stakingTiers
}
