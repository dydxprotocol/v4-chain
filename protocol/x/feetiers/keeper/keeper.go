package keeper

import (
	"fmt"
	"math"
	"math/big"

	"github.com/cometbft/cometbft/libs/log"

	sdklog "cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
)

type (
	Keeper struct {
		cdc         codec.BinaryCodec
		statsKeeper types.StatsKeeper
		storeKey    storetypes.StoreKey
		authorities map[string]struct{}
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	statsKeeper types.StatsKeeper,
	storeKey storetypes.StoreKey,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:         cdc,
		statsKeeper: statsKeeper,
		storeKey:    storeKey,
		authorities: lib.UniqueSliceToSet(authorities),
	}
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(sdklog.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {}

func (k Keeper) getUserFeeTier(ctx sdk.Context, address string) (uint32, *types.PerpetualFeeTier) {
	userStats := k.statsKeeper.GetUserStats(ctx, address)
	globalStats := k.statsKeeper.GetGlobalStats(ctx)

	// Invariant: we know there is at least one tier and that the first tier has no requirements
	tiers := k.GetPerpetualFeeParams(ctx).Tiers
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

	return idx, tiers[idx]
}

func (k Keeper) GetPerpetualFeePpm(ctx sdk.Context, address string, isTaker bool) int32 {
	_, userTier := k.getUserFeeTier(ctx, address)
	if isTaker {
		return userTier.TakerFeePpm
	}
	return userTier.MakerFeePpm
}

// GetLowestMakerFee returns the lowest maker fee among any tiers.
func (k Keeper) GetLowestMakerFee(ctx sdk.Context) int32 {
	feeParams := k.GetPerpetualFeeParams(ctx)

	lowestMakerFee := int32(math.MaxInt32)
	for _, tier := range feeParams.Tiers {
		if tier.MakerFeePpm < lowestMakerFee {
			lowestMakerFee = tier.MakerFeePpm
		}
	}

	return lowestMakerFee
}
