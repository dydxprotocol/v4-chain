package keeper

import (
	"fmt"
	"math/big"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
)

// validateStakingTiers validates that:
// - No duplicate fee tier names
// - Fee tier name is not empty
// - Staking levels are valid:
//   - Min staked tokens is a valid non-negative number
//   - Levels are in strictly increasing order of min staked tokens
//   - Discount is not more than 100%
//
// - Staking tiers correspond to existing fee tiers
func (k Keeper) validateStakingTiers(ctx sdk.Context, stakingTiers []*types.StakingTier) error {
	seenTiers := make(map[string]bool)
	for _, tier := range stakingTiers {
		// Validate fee tier name is not empty
		if tier.FeeTierName == "" {
			return fmt.Errorf("fee tier name cannot be empty")
		}
		// Check for duplicate fee tier names
		if seenTiers[tier.FeeTierName] {
			return fmt.Errorf("duplicate staking tier for fee tier: %s", tier.FeeTierName)
		}
		seenTiers[tier.FeeTierName] = true

		// Validate staking levels
		var prevMinStaked *big.Int
		for i, level := range tier.Levels {
			// Validate min staked tokens is a valid number
			minStaked := new(big.Int)
			if _, ok := minStaked.SetString(level.MinStakedBaseTokens, 10); !ok {
				return fmt.Errorf("invalid min staked tokens for tier %s level %d: %s",
					tier.FeeTierName, i, level.MinStakedBaseTokens)
			}

			// Check that min staked is non-negative
			if minStaked.Sign() < 0 {
				return fmt.Errorf("min staked tokens cannot be negative for tier %s level %d",
					tier.FeeTierName, i)
			}

			// Check that levels are in increasing order
			if prevMinStaked != nil && minStaked.Cmp(prevMinStaked) <= 0 {
				return fmt.Errorf("staking levels must be in increasing order for tier %s",
					tier.FeeTierName)
			}
			prevMinStaked = minStaked

			// Validate discount is not more than 100%
			if level.FeeDiscountPpm > 1_000_000 {
				return fmt.Errorf("fee discount cannot exceed 100%% for tier %s level %d",
					tier.FeeTierName, i)
			}
		}
	}

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

	return nil
}

// SetStakingTiers sets staking tiers in state
func (k Keeper) SetStakingTiers(ctx sdk.Context, stakingTiers []*types.StakingTier) error {
	// Validate staking tiers
	if err := k.validateStakingTiers(ctx, stakingTiers); err != nil {
		return err
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
