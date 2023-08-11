package keeper

import (
	"fmt"
	"math/big"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4/dtypes"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/lib/metrics"
	"github.com/dydxprotocol/v4/x/rewards/types"
)

type (
	Keeper struct {
		cdc               codec.BinaryCodec
		storeKey          storetypes.StoreKey
		transientStoreKey storetypes.StoreKey

		feeTiersKeeper types.FeeTiersKeeper

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	transientStoreKey storetypes.StoreKey,
	feeTiersKeeper types.FeeTiersKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:               cdc,
		storeKey:          storeKey,
		transientStoreKey: transientStoreKey,
		feeTiersKeeper:    feeTiersKeeper,
		authority:         authority,
	}
}

// GetAuthority returns the x/rewards module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}

// Get `RewardShare` for a given address.
// If the address does not have existing reward share, return a
// `RewardShare` with 0 weight.
func (k Keeper) GetRewardShare(
	ctx sdk.Context,
	address string,
) (val types.RewardShare) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.GetRewardShare,
		metrics.Latency,
	)

	// Check state for the subaccount.
	store := prefix.NewStore(ctx.KVStore(k.transientStoreKey), types.KeyPrefix(types.RewardShareKeyPrefix))
	b := store.Get(types.RewardShareKey(address))

	// If RewardShare does not exist in state, return a default value.
	if b == nil {
		return types.RewardShare{
			Address: address,
			Weight:  dtypes.NewInt(0),
		}
	}

	// If RewardShare does exist in state, unmarshall and return the value.
	k.cdc.MustUnmarshal(b, &val)
	return val
}

// Add reward shares for the maker and taker of a fill. Intended for being called in `x/clob` when a fill is persisted.
//
// Within each block, total reward share score for an address is defined as:
//
//	reward_share_score = total_taker_fees_paid - max_possible_maker_rebate*taker_volume + total_positive_maker_fees
//
// Hence, for each fill, increment reward share score as follow:
//   - For maker address, positive maker fees are added directly.
//   - For taker address, positive taker fees are reduced by the largest possible maker rebate in x/fee-tiers multiplied
//     by quote quantums of the fill.
func (k Keeper) AddRewardSharesForFill(
	ctx sdk.Context,
	takerAddress string,
	makerAddress string,
	bigFillQuoteQuantums *big.Int,
	bigTakerFeeQuoteQuantums *big.Int,
	bigMakerFeeQuoteQuantums *big.Int,
) {
	// Process reward weight for taker.
	lowestMakerFee := k.feeTiersKeeper.GetLowestMakerFee(ctx)
	maxMakerRebatePpm := lib.Min(int32(0), lowestMakerFee)
	// Calculate quote_quantums * max_maker_rebate. Result is non-positive.
	makerRebateMulTakerVolume := lib.BigIntMulSignedPpm(bigFillQuoteQuantums, maxMakerRebatePpm)
	takerWeight := new(big.Int).Add(
		bigTakerFeeQuoteQuantums,
		makerRebateMulTakerVolume,
	)
	if takerWeight.Cmp(lib.BigInt0()) > 0 {
		k.AddRewardShareToAddress(
			ctx,
			takerAddress,
			takerWeight,
		)
	}

	// Process reward weight for maker.
	makerWeight := new(big.Int).Set(bigMakerFeeQuoteQuantums)
	if makerWeight.Cmp(lib.BigInt0()) > 0 {
		k.AddRewardShareToAddress(
			ctx,
			makerAddress,
			makerWeight,
		)
	}
}

// AddRewardShareToAddress adds a reward share to an address.
// If the address has a previous reward share, increment weight.
// If not, create new reward share with given weight.
func (k Keeper) AddRewardShareToAddress(
	ctx sdk.Context,
	address string,
	weight *big.Int,
) {
	// Get existing reward share. If no previous reward share, 0 weight is returned.
	rewardShare := k.GetRewardShare(ctx, address)
	newWeight := new(big.Int).Add(
		rewardShare.Weight.BigInt(),
		weight,
	)

	// Set the new reward share.
	k.SetRewardShare(ctx, types.RewardShare{
		Address: address,
		Weight:  dtypes.NewIntFromBigInt(newWeight),
	})
}

// SetRewardShare set a reward share object under rewardShare.Address.
func (k Keeper) SetRewardShare(
	ctx sdk.Context,
	rewardShare types.RewardShare,
) {
	store := prefix.NewStore(ctx.KVStore(k.transientStoreKey), types.KeyPrefix(types.RewardShareKeyPrefix))
	b := k.cdc.MustMarshal(&rewardShare)

	store.Set(types.RewardShareKey(
		rewardShare.Address,
	), b)
}
