package keeper

import (
	"fmt"
	"math"
	"math/big"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

type (
	Keeper struct {
		cdc              codec.BinaryCodec
		statsKeeper      types.StatsKeeper
		vaultKeeper      types.VaultKeeper
		storeKey         storetypes.StoreKey
		authorities      map[string]struct{}
		affiliatesKeeper types.AffiliatesKeeper
		revShareKeeper   types.RevShareKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	statsKeeper types.StatsKeeper,
	affiliatesKeeper types.AffiliatesKeeper,
	storeKey storetypes.StoreKey,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:              cdc,
		statsKeeper:      statsKeeper,
		storeKey:         storeKey,
		authorities:      lib.UniqueSliceToSet(authorities),
		affiliatesKeeper: affiliatesKeeper,
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

// SetVaultKeeper sets the `VaultKeeper` reference in `FeeTiersKeeper`.
// The reference is set with an explicit method call rather than during `NewKeeper`
// due to the circular dependency `Clob` -> `Vault` -> `FeeTiers` -> `Clob`.
func (k *Keeper) SetVaultKeeper(vk types.VaultKeeper) {
	k.vaultKeeper = vk
}

func (k Keeper) getUserFeeTier(
	ctx sdk.Context,
	address string,
	feeTierOverrideIdx uint32,
) (uint32, *types.PerpetualFeeTier) {
	tiers := k.GetPerpetualFeeParams(ctx).Tiers

	// A vault is always in the highest tier.
	// Invariant: there's at least one tier.
	if k.vaultKeeper.IsVault(ctx, address) {
		highestTierIdx := uint32(len(tiers) - 1)
		return highestTierIdx, tiers[highestTierIdx]
	}

	userStats := k.statsKeeper.GetUserStats(ctx, address)
	globalStats := k.statsKeeper.GetGlobalStats(ctx)

	// Invariant: we know there is at least one tier and that the first tier has no requirements
	idx := uint32(0)

	// Find the last tier we meet all requirements for
	for i := 0; i < len(tiers); i++ {
		currTier := tiers[i]
		bigUserMakerNotional := new(big.Int).SetUint64(userStats.MakerNotional)
		bigUserTakerNotional := new(big.Int).SetUint64(userStats.TakerNotional)
		bigUserTotalNotional := new(big.Int).Add(bigUserMakerNotional, bigUserTakerNotional)
		bigGlobalNotional := new(big.Int).SetUint64(globalStats.NotionalTraded)

		bigAbsVolumeRequirement := new(big.Int).SetUint64(currTier.AbsoluteVolumeRequirement)
		bigTotalVolumeShareRequirement := lib.BigIntMulPpm(
			bigGlobalNotional,
			currTier.TotalVolumeShareRequirementPpm,
		)
		bigMakerVolumeShareRequirement := lib.BigIntMulPpm(
			bigGlobalNotional,
			currTier.MakerVolumeShareRequirementPpm,
		)

		if bigUserTotalNotional.Cmp(bigAbsVolumeRequirement) == -1 ||
			bigUserTotalNotional.Cmp(bigTotalVolumeShareRequirement) == -1 ||
			bigUserMakerNotional.Cmp(bigMakerVolumeShareRequirement) == -1 {
			break
		}
		idx = uint32(i)
	}

	maxTierIdx := uint32(len(tiers) - 1)
	if feeTierOverrideIdx > maxTierIdx {
		feeTierOverrideIdx = maxTierIdx
	}

	if idx < feeTierOverrideIdx {
		_, hasReferree := k.affiliatesKeeper.GetReferredBy(ctx, address)
		if hasReferree {
			idx = feeTierOverrideIdx
		}
	}

	return idx, tiers[idx]
}

// GetPerpetualFeePpm returns the fee PPM (parts per million) for a user.
// It checks if
// 1. there's an active fee discount for the specified CLOB pair.
// 2. user qualifies for staking-based discounts.
func (k Keeper) GetPerpetualFeePpm(
	ctx sdk.Context,
	address string,
	isTaker bool,
	feeTierOverrideIdx uint32,
	clobPairId uint32,
) int32 {
	_, userTier := k.getUserFeeTier(ctx, address, feeTierOverrideIdx)
	var baseFee int32
	if isTaker {
		baseFee = userTier.TakerFeePpm
	} else {
		baseFee = userTier.MakerFeePpm
	}

	// Get the per-market discount PPM (returns MaxChargePpm = 1,000,000 = 100% if no active fee discount)
	perMarketDiscountPpm := k.GetDiscountedPpm(ctx, clobPairId)

	// Calculate the fee after per-market discount
	// For negative fees (rebates), we also apply the discount percentage
	feeAfterMarketDiscount := int32(int64(baseFee) * int64(perMarketDiscountPpm) / int64(types.MaxChargePpm))

	// Apply staking discount if fee is positive and user qualifies
	if feeAfterMarketDiscount > 0 {
		// Validate address before getting staked amount
		_, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			// Log error but do not fail fee calculation
			k.Logger(ctx).Error(
				"Failed to validate address for staking discount",
				"address", address,
				"error", err,
			)
		} else {
			stakedAmount := k.statsKeeper.GetStakedBaseTokens(ctx, address)
			stakingDiscountPpm := k.GetStakingDiscountPpm(ctx, userTier.Name, stakedAmount)
			if stakingDiscountPpm > 0 {
				// Final fee
				// = fee * (1 - staking_discount)
				// = fee * (1_000_000 - staking_discount_ppm) / 1_000_000
				remainingFeePpm := types.MaxChargePpm - stakingDiscountPpm
				feeAfterMarketDiscount = int32(int64(feeAfterMarketDiscount) * int64(remainingFeePpm) / int64(types.MaxChargePpm))
			}
		}
	}

	return feeAfterMarketDiscount
}

// GetLowestMakerFee returns the lowest maker fee among any tiers.
func (k Keeper) GetLowestMakerFee(ctx sdk.Context) int32 {
	feeParams := k.GetPerpetualFeeParams(ctx)

	return GetLowestMakerFeeFromTiers(feeParams.Tiers)
}

func (k Keeper) GetAffiliateRefereeLowestTakerFee(ctx sdk.Context) int32 {
	feeParams := k.GetPerpetualFeeParams(ctx)

	return GetAffiliateRefereeLowestTakerFeeFromTiers(feeParams.Tiers)
}

func (k *Keeper) SetRevShareKeeper(revShareKeeper types.RevShareKeeper) {
	k.revShareKeeper = revShareKeeper
}

func GetLowestMakerFeeFromTiers(tiers []*types.PerpetualFeeTier) int32 {
	lowestMakerFee := int32(math.MaxInt32)
	for _, tier := range tiers {
		if tier.MakerFeePpm < lowestMakerFee {
			lowestMakerFee = tier.MakerFeePpm
		}
	}
	return lowestMakerFee
}

// GetAffiliateRefereeLowestTakerFeeFromTiers returns the minimum of
// - the taker fee of the tier that has the max absolute volume requirement
// - the taker fee of the referee starting fee tier
func GetAffiliateRefereeLowestTakerFeeFromTiers(tiers []*types.PerpetualFeeTier) int32 {
	takerFeePpm := int32(math.MaxInt32)
	for _, tier := range tiers {
		// assumes tiers are ordered by absolute volume requirement
		if tier.AbsoluteVolumeRequirement < revsharetypes.MaxReferee30dVolumeForAffiliateShareQuantums {
			takerFeePpm = tier.TakerFeePpm
		} else {
			break
		}
	}

	if uint32(len(tiers)) > types.RefereeStartingFeeTier {
		return min(takerFeePpm, tiers[types.RefereeStartingFeeTier].TakerFeePpm)
	}

	return takerFeePpm
}
