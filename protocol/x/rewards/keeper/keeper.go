package keeper

import (
	"fmt"
	"math/big"
	"time"

	errorsmod "cosmossdk.io/errors"
	cosmoslog "cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
)

type (
	Keeper struct {
		cdc               codec.BinaryCodec
		storeKey          storetypes.StoreKey
		transientStoreKey storetypes.StoreKey

		// Needed for getting `UsdcAsset.AtomicResolution` (converting quote quantums to a full USDC).
		assetsKeeper types.AssetsKeeper
		// Need for getting balance of module account balance and transferring tokens.
		bankKeeper types.BankKeeper
		// Needed for getting lowest maker fee.
		feeTiersKeeper types.FeeTiersKeeper
		// Neeeded for retrieve market price of rewards token.
		pricesKeeper        types.PricesKeeper
		indexerEventManager indexer_manager.IndexerEventManager

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
	indexerEventManager indexer_manager.IndexerEventManager,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		transientStoreKey:   transientStoreKey,
		assetsKeeper:        assetsKeeper,
		bankKeeper:          bankKeeper,
		feeTiersKeeper:      feeTiersKeeper,
		pricesKeeper:        pricesKeeper,
		indexerEventManager: indexerEventManager,
		authorities:         lib.UniqueSliceToSet(authorities),
	}
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

func (k Keeper) GetIndexerEventManager() indexer_manager.IndexerEventManager {
	return k.indexerEventManager
}

func (k Keeper) Logger(ctx sdk.Context) cosmoslog.Logger {
	return ctx.Logger().With(cosmoslog.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
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
	store := prefix.NewStore(ctx.KVStore(k.transientStoreKey), []byte(types.RewardShareKeyPrefix))
	b := store.Get([]byte(address))

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
//	reward_share_score = total_taker_fees_paid - max_possible_taker_fee_rev_share
//   - max_possible_maker_rebate * taker_volume + total_positive_maker_fees - total_rev_shared_maker_fee
//
// Hence, for each fill, increment reward share score as follow:
//   - Let F = sum(percentages of general rev-share) (excluding taker only rev share i.e. affiliate)
//   - For maker address, positive_maker_fees * (1 - F) are added to reward share score.
//   - For taker address, (positive_taker_fees - max_possible_maker_rebate
//     					  * fill_quote_quantum - max_possible_taker_fee_rev_share) * (1 - F)
//     are added to reward share score.
// max_possible_taker_fee_rev_share is 0 when taker trailing volume is > MaxReferee30dVolumeForAffiliateShareQuantums,
// since taker_fee_share is only affiliate at the moment, and they don’t generate affiliate rev share.
// When taker volume ≤ MaxReferee30dVolumeForAffiliateShareQuantums,
// max_possible_taker_fee_rev_share = max_vip_affiliate_share * taker_fee
// regardless of if the taker has an affiliate or not.

func (k Keeper) AddRewardSharesForFill(
	ctx sdk.Context,
	fill clobtypes.FillForProcess,
	revSharesForFill revsharetypes.RevSharesForFill,
) {
	// Process reward weight for taker.
	lowestMakerFee := k.feeTiersKeeper.GetLowestMakerFee(ctx)
	maxMakerRebatePpm := lib.Min(int32(0), lowestMakerFee)

	// Calculate total net fee rev share ppm on the protocol
	totalNetFeeRevSharePpm := uint32(0)
	if value, ok := revSharesForFill.FeeSourceToRevSharePpm[revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE]; ok {
		totalNetFeeRevSharePpm = value
	}
	maxPossibleAffiliateRevShareQuoteQuantum := big.NewInt(0)

	// taker revshare is always 0 if taker rolling volume is greater than or equal
	// to Max30dTakerVolumeQuantums, so no need to reduce score by `max_possible_taker_fee_rev_share`
	if fill.MonthlyRollingTakerVolumeQuantums < revsharetypes.MaxReferee30dVolumeForAffiliateShareQuantums {
		maxPossibleAffiliateRevShareQuoteQuantum = lib.BigMulPpm(fill.TakerFeeQuoteQuantums,
			lib.BigU(affiliatetypes.AffiliatesRevSharePpmCap),
			false,
		)
	}

	// Remove the taker order router rev share from the taker fee
	// Taker order router rev share is mutually exclusive with affiliate rev share
	takerOrderRouterRevShare := big.NewInt(0)
	for _, share := range revSharesForFill.AllRevShares {
		if share.RevShareFeeSource == revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE &&
			share.RevShareType == revsharetypes.REV_SHARE_TYPE_ORDER_ROUTER {
			takerOrderRouterRevShare = share.QuoteQuantums
		}
	}

	// This is the amount of remaining fee: 1 - total_net_fee_rev_share_ppm
	totalFeeSubNetRevSharePpm := lib.OneMillion - totalNetFeeRevSharePpm

	// Calculate quote_quantums * max_maker_rebate. Result is non-positive.
	// This is the amount of quote quantums that is rebate to the maker
	makerRebateMulTakerVolume := lib.BigMulPpm(fill.FillQuoteQuantums, lib.BigI(maxMakerRebatePpm), false)

	// Remove the rebate given to the maker from the taker fee
	netTakerFee := new(big.Int).Add(
		fill.TakerFeeQuoteQuantums,
		makerRebateMulTakerVolume,
	)

	// Remove the affiliate or order router fee from the taker fee
	netTakerFee = netTakerFee.Sub(
		netTakerFee,
		maxPossibleAffiliateRevShareQuoteQuantum,
	)

	// Remove the taker order router fee from the taker fee
	netTakerFee = netTakerFee.Sub(
		netTakerFee,
		takerOrderRouterRevShare,
	)

	// Factor out the protocol fees given as rev shares
	takerWeight := lib.BigMulPpm(
		netTakerFee,
		lib.BigU(totalFeeSubNetRevSharePpm),
		false,
	)

	// Give the taker the remaining reward shares
	if takerWeight.Sign() > 0 {
		// We aren't concerned with errors here because we've already validated the weight is positive.
		if err := k.AddRewardShareToAddress(
			ctx,
			fill.TakerAddr,
			takerWeight,
		); err != nil {
			log.InfoLog(
				ctx,
				"Failed to add rewards share to address",
				err,
			)
		}
	}

	// Process reward weight for maker.
	// This is the maker fee quote quantums * (1 - total_net_fee_rev_share_ppm)
	makerOrderRouterRevShare := big.NewInt(0)
	for _, share := range revSharesForFill.AllRevShares {
		if share.RevShareFeeSource == revsharetypes.REV_SHARE_FEE_SOURCE_MAKER_FEE &&
			share.RevShareType == revsharetypes.REV_SHARE_TYPE_ORDER_ROUTER {
			makerOrderRouterRevShare = share.QuoteQuantums
		}
	}

	// Remove the maker order router rev share from the maker fee
	// Maker ORRS is 0 if there is a rebate
	netMakerFee := new(big.Int).Sub(
		fill.MakerFeeQuoteQuantums,
		makerOrderRouterRevShare,
	)
	// Factor out the protocol fees from the remaining maker fee
	makerWeight := new(big.Int).Set(lib.BigMulPpm(netMakerFee, lib.BigU(totalFeeSubNetRevSharePpm), false))
	if makerWeight.Sign() > 0 {
		// We aren't concerned with errors here because we've already validated the weight is positive.
		if err := k.AddRewardShareToAddress(
			ctx,
			fill.MakerAddr,
			makerWeight,
		); err != nil {
			log.InfoLog(
				ctx,
				"Failed to add rewards share to address",
				err,
			)
		}
	}
}

// AddRewardShareToAddress adds a reward share to an address.
// If the address has a previous reward share, increment weight.
// If not, create new reward share with given weight.
func (k Keeper) AddRewardShareToAddress(
	ctx sdk.Context,
	address string,
	weight *big.Int,
) error {
	if weight.Sign() <= 0 {
		return errorsmod.Wrapf(
			types.ErrNonpositiveWeight,
			"Invalid weight %v",
			weight.String(),
		)
	}

	// Get existing reward share. If no previous reward share, 0 weight is returned.
	rewardShare := k.GetRewardShare(ctx, address)
	newWeight := new(big.Int).Add(
		rewardShare.Weight.BigInt(),
		weight,
	)

	// Set the new reward share.
	return k.SetRewardShare(ctx, types.RewardShare{
		Address: address,
		Weight:  dtypes.NewIntFromBigInt(newWeight),
	})
}

// SetRewardShare set a reward share object under rewardShare.Address.
func (k Keeper) SetRewardShare(
	ctx sdk.Context,
	rewardShare types.RewardShare,
) error {
	if rewardShare.Weight.Sign() <= 0 {
		return errorsmod.Wrapf(
			types.ErrNonpositiveWeight,
			"Invalid weight %v",
			rewardShare.Weight.String(),
		)
	}

	store := prefix.NewStore(ctx.KVStore(k.transientStoreKey), []byte(types.RewardShareKeyPrefix))
	b := k.cdc.MustMarshal(&rewardShare)

	store.Set([]byte(rewardShare.Address), b)
	return nil
}

func (k Keeper) getAllRewardSharesAndTotalWeight(ctx sdk.Context) (
	list []types.RewardShare,
	totalWeight *big.Int,
) {
	store := prefix.NewStore(ctx.KVStore(k.transientStoreKey), []byte(types.RewardShareKeyPrefix))
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
	rewardTokenPrice, err := k.pricesKeeper.GetMarketPrice(ctx, params.GetMarketId())
	if err != nil {
		return fmt.Errorf("failed to get market price of reward token: %w", err)
	}
	allRewardShares, totalRewardWeight := k.getAllRewardSharesAndTotalWeight(ctx)
	// Measure total reward weight.
	telemetry.SetGauge(
		metrics.GetMetricValueFromBigInt(totalRewardWeight),
		types.ModuleName,
		metrics.TotalRewardShareWeight,
	)
	totalRewardWeightPpm := new(big.Int).Mul(totalRewardWeight, lib.BigU(params.FeeMultiplierPpm))
	rewardTokenAmountPpm := lib.QuoteToBaseQuantums(
		totalRewardWeightPpm,
		params.DenomExponent,
		rewardTokenPrice.Price,
		rewardTokenPrice.Exponent,
	)
	rewardTokenAmount := new(big.Int).Div(rewardTokenAmountPpm, lib.BigIntOneMillion())

	// Calculate value of `T`, the reward tokens balance in the `treasury_account`.
	rewardTokenBalance := k.bankKeeper.GetBalance(ctx, types.TreasuryModuleAddress, params.Denom)

	// Get tokenToDistribute as the min(F, T).
	tokensToDistribute := lib.BigMin(rewardTokenBalance.Amount.BigInt(), rewardTokenAmount)
	// Measure distributed token amount.
	telemetry.SetGauge(
		metrics.GetMetricValueFromBigInt(tokensToDistribute),
		types.ModuleName,
		metrics.DistributedRewardTokens,
	)
	if tokensToDistribute.Sign() == 0 {
		// Nothing to distribute. This can happen either when there is no reward token in the treasury account,
		// or if no reward shares were recorded for this block.
		return nil
	}

	rewardIndexerEvent := indexerevents.TradingRewardsEventV1{}

	// Go through each address with reward and distribute tokens.
	for _, share := range allRewardShares {
		// Calculate `tokensToDistribute` * `share.Weight` / `totalRewardWeight`.
		rewardAmountForAddress := new(big.Int).Div(
			new(big.Int).Mul(
				tokensToDistribute,
				share.Weight.BigInt(),
			),
			totalRewardWeight,
		) // big.Div() rounds down, so sum of actual distributed tokens will not exceed `tokensToDistribute`

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
			log.ErrorLogWithError(
				ctx,
				"Failed to send reward tokens from treasury account to address",
				err,
				"treasury_account",
				params.TreasuryAccount,
				"address",
				share.Address,
			)
		} else {
			rewardIndexerEvent.TradingRewards = append(rewardIndexerEvent.TradingRewards,
				&indexerevents.AddressTradingReward{
					Owner:       share.Address,
					DenomAmount: dtypes.NewIntFromBigInt(rewardAmountForAddress),
				},
			)
		}
	}

	k.indexerEventManager.AddBlockEvent(
		ctx,
		indexerevents.SubtypeTradingReward,
		indexer_manager.IndexerTendermintEvent_BLOCK_EVENT_END_BLOCK,
		indexerevents.TradingRewardVersion,
		indexer_manager.GetBytes(&rewardIndexerEvent),
	)

	// Measure treasury balance after distribution.
	remainingTreasuryBalance := k.bankKeeper.GetBalance(ctx, types.TreasuryModuleAddress, params.Denom)
	telemetry.SetGauge(
		metrics.GetMetricValueFromBigInt(remainingTreasuryBalance.Amount.BigInt()),
		types.ModuleName,
		metrics.TreasuryBalanceAfterDistribution,
	)

	return nil
}
