package keeper

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	"math/big"
	"time"

	"cosmossdk.io/log"
	sdklog "cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/maps"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
)

type (
	Keeper struct {
		cdc               codec.BinaryCodec
		storeKey          storetypes.StoreKey
		transientStoreKey storetypes.StoreKey

		// Needed for getting `UsdcAsset.AtomicResolution` (converting quote quantums to a full USDC).
		assetsKeeper types.AssetsKeeper
		// Need for getting balance of module account balance and transfering tokens.
		bankKeeper types.BankKeeper
		// Needed for getting lowest maker fee.
		feeTiersKeeper types.FeeTiersKeeper
		// Neeeded for retrieve market price of rewards token.
		pricesKeeper types.PricesKeeper

		// the addresses capable of executing a MsgUpdateParams message.
		authorities map[string]struct{}
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	transientStoreKey storetypes.StoreKey,
	assetsKeeper types.AssetsKeeper,
	bankKeeper types.BankKeeper,
	feeTiersKeeper types.FeeTiersKeeper,
	pricesKeeper types.PricesKeeper,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:               cdc,
		storeKey:          storeKey,
		transientStoreKey: transientStoreKey,
		assetsKeeper:      assetsKeeper,
		bankKeeper:        bankKeeper,
		feeTiersKeeper:    feeTiersKeeper,
		pricesKeeper:      pricesKeeper,
		authorities:       maps.ArrayToMapInterface(authorities),
	}
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(sdklog.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
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

func (k Keeper) getAllRewardSharesAndTotalWeight(ctx sdk.Context) (
	list []types.RewardShare,
	totalWeight *big.Int,
) {
	store := prefix.NewStore(ctx.KVStore(k.transientStoreKey), types.KeyPrefix(types.RewardShareKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	totalWeight = big.NewInt(0)
	for ; iterator.Valid(); iterator.Next() {
		var val types.RewardShare
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
		totalWeight.Add(
			totalWeight,
			val.Weight.BigInt(),
		)
	}
	return list, totalWeight
}

// ProcessRewardsForBlock processes rewards for all fills that happened in a block.
// The amount A of the reward token to be distributed to traders is defined as:
//
//	A = min(F, T)
//
// where:
//
//	`T` is the amount of available reward tokens in the `treasury_account`.
//	`F` = fee_multiplier * (total_positive_maker_fees +
//		                    total taker fees -
//		                    maximum possible maker rebate * total taker volume)
//	                     / reward_token_price
func (k Keeper) ProcessRewardsForBlock(
	ctx sdk.Context,
) error {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.ProcessRewardsForBlock,
		metrics.Latency,
	)

	// Get reward params.
	params := k.GetParams(ctx)

	// Calculate value of `F`.
	usdcAsset, err := k.assetsKeeper.GetAsset(ctx, lib.UsdcAssetId)
	if err != nil {
		return fmt.Errorf("failed to get USDC asset: %w", err)
	}
	rewardTokenPrice, err := k.pricesKeeper.GetMarketPrice(ctx, params.GetMarketId())
	if err != nil {
		return fmt.Errorf("failed to get market price of reward token: %w", err)
	}
	allRewardShares, totalRewardWeight := k.getAllRewardSharesAndTotalWeight(ctx)
	// Measure total reward weight.
	telemetry.SetGauge(
		float32(totalRewardWeight.Int64()),
		types.ModuleName,
		metrics.TotalRewardShareWeight,
	)
	bigRatRewardTokenAmount := clobtypes.NotionalToCoinAmount(
		totalRewardWeight,
		usdcAsset.AtomicResolution,
		params.DenomExponent,
		rewardTokenPrice,
	)
	bigRatRewardTokenAmount = lib.BigRatMulPpm(
		bigRatRewardTokenAmount,
		params.FeeMultiplierPpm,
	)
	bigIntRewardTokenAmount := lib.BigRatRound(bigRatRewardTokenAmount, false)

	// Calculate value of `T`, the reward tokens balance in the `treasury_account`.
	rewardTokenBalance := k.bankKeeper.GetBalance(
		ctx,
		authtypes.NewModuleAddress(params.TreasuryAccount),
		params.Denom,
	)

	// Get tokenToDistribute as the min(F, T).
	tokensToDistribute := lib.BigMin(rewardTokenBalance.Amount.BigInt(), bigIntRewardTokenAmount)
	// Measure distributed token amount.
	telemetry.SetGauge(
		float32(tokensToDistribute.Int64()),
		types.ModuleName,
		metrics.DistributedRewardTokens,
	)
	if tokensToDistribute.Sign() == 0 {
		// Nothing to distribute. This can happen either when there is no reward token in the treasury account,
		// or if no reward shares were recorded for this block.
		return nil
	}

	// Go through each address with reward and distribute tokens.
	for _, share := range allRewardShares {
		// Calculate `tokensToDistribute` * `share.Weight` / `totalRewardWeight`.
		rewardAmountForAddress := new(big.Int).Div(
			new(big.Int).Mul(
				tokensToDistribute,
				share.Weight.BigInt(),
			),
			totalRewardWeight,
		) // big.Div() rounds down, so sum of actual distributed tokens will not exeed `tokensToDistribute`

		if rewardAmountForAddress.Sign() == 0 {
			// Nothing to distribute to this address. This will only happen due to rounding.
			continue
		}

		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx,
			params.TreasuryAccount,
			// MustAccAddressFromBech32() panics if the address is invalid.
			// This should never happen, since the address is taken from the address field
			// of the fill object.
			sdk.MustAccAddressFromBech32(share.Address),
			[]sdk.Coin{
				{
					Denom:  params.Denom,
					Amount: sdkmath.NewIntFromBigInt(rewardAmountForAddress),
				},
			},
		); err != nil {
			panic(
				fmt.Errorf(
					"failed to send reward tokens from treasury (%s) to address %s: %w",
					params.TreasuryAccount, share.Address, err,
				),
			)
		}
	}

	// Measure treasury balance after distribution.
	remainingTreasuryBalance := k.bankKeeper.GetBalance(
		ctx,
		authtypes.NewModuleAddress(params.TreasuryAccount),
		params.Denom,
	)
	telemetry.SetGauge(
		float32(remainingTreasuryBalance.Amount.Int64()),
		types.ModuleName,
		metrics.TreasuryBalanceAfterDistribution,
	)

	return nil
}
