package keeper

import (
	"math/big"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

// Function to serialize market mapper revenue share params and store in the module store
func (k Keeper) SetMarketMapperRevenueShareParams(
	ctx sdk.Context,
	params types.MarketMapperRevenueShareParams,
) (err error) {
	// Validate the params
	if err := params.Validate(); err != nil {
		return err
	}

	// Store the params in the module store
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&params)
	store.Set([]byte(types.MarketMapperRevenueShareParamsKey), b)

	return nil
}

// Function to get market mapper revenue share params from the module store
func (k Keeper) GetMarketMapperRevenueShareParams(
	ctx sdk.Context,
) (params types.MarketMapperRevenueShareParams) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.MarketMapperRevenueShareParamsKey))
	k.cdc.MustUnmarshal(b, &params)
	return params
}

// Function to serialize market mapper rev share details for a market
// and store in the module store
func (k Keeper) SetMarketMapperRevShareDetails(
	ctx sdk.Context,
	marketId uint32,
	params types.MarketMapperRevShareDetails,
) {
	// Store the rev share details for provided market in module store
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.MarketMapperRevSharePrefix))
	b := k.cdc.MustMarshal(&params)
	store.Set(lib.Uint32ToKey(marketId), b)
}

// Function to retrieve marketmapper revshare details for a market from module store
func (k Keeper) GetMarketMapperRevShareDetails(
	ctx sdk.Context,
	marketId uint32,
) (params types.MarketMapperRevShareDetails, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.MarketMapperRevSharePrefix))
	b := store.Get(lib.Uint32ToKey(marketId))
	if b == nil {
		return params, types.ErrMarketMapperRevShareDetailsNotFound
	}
	k.cdc.MustUnmarshal(b, &params)
	return params, nil
}

// Function to perform all market creation actions for the revshare module
func (k Keeper) CreateNewMarketRevShare(ctx sdk.Context, marketId uint32) {
	revShareParams := k.GetMarketMapperRevenueShareParams(ctx)

	validDurationSeconds := int64(revShareParams.ValidDays * 24 * 60 * 60)

	// set the rev share details for the market
	details := types.MarketMapperRevShareDetails{
		ExpirationTs: uint64(ctx.BlockTime().Unix() + validDurationSeconds),
	}
	k.SetMarketMapperRevShareDetails(ctx, marketId, details)
}

func (k Keeper) GetMarketMapperRevenueShareForMarket(ctx sdk.Context, marketId uint32) (
	address sdk.AccAddress,
	revenueSharePpm uint32,
	err error,
) {
	// get the revenue share details for the market
	revShareDetails, err := k.GetMarketMapperRevShareDetails(ctx, marketId)
	if err != nil {
		return nil, 0, err
	}

	// check if the rev share details are expired
	if revShareDetails.ExpirationTs < uint64(ctx.BlockTime().Unix()) {
		return nil, 0, nil
	}

	// Get revenue share params
	revShareParams := k.GetMarketMapperRevenueShareParams(ctx)

	revShareAddr, err := sdk.AccAddressFromBech32(revShareParams.Address)
	if err != nil {
		return nil, 0, err
	}

	return revShareAddr, revShareParams.RevenueSharePpm, nil
}

func (k Keeper) GetUnconditionalRevShareConfigParams(ctx sdk.Context) (types.UnconditionalRevShareConfig, error) {
	store := ctx.KVStore(k.storeKey)
	unconditionalRevShareConfigBytes := store.Get(
		[]byte(types.UnconditionalRevShareConfigKey),
	)
	var unconditionalRevShareConfig types.UnconditionalRevShareConfig
	k.cdc.MustUnmarshal(unconditionalRevShareConfigBytes, &unconditionalRevShareConfig)
	return unconditionalRevShareConfig, nil
}

func (k Keeper) SetUnconditionalRevShareConfigParams(ctx sdk.Context, config types.UnconditionalRevShareConfig) {
	store := ctx.KVStore(k.storeKey)
	unconditionalRevShareConfigBytes := k.cdc.MustMarshal(&config)
	store.Set([]byte(types.UnconditionalRevShareConfigKey), unconditionalRevShareConfigBytes)
}

// Check two conditions to ensure rev shares are safe.
// 1. totalUnconditionalRevSharePpm + totalMarketMapperRevSharePpm < 100%
// 2. lowest_taker_fee + lowest_maker_fee >= Highest_affiliate_rev_share * lowest_taker_fee
func (k Keeper) ValidateRevShareSafety(
	ctx sdk.Context,
	unconditionalRevShareConfig types.UnconditionalRevShareConfig,
	marketMapperRevShareParams types.MarketMapperRevenueShareParams,
	lowestTakerFeePpm int32,
	lowestMakerFeePpm int32,
) bool {
	totalUnconditionalRevSharePpm := uint32(0)
	for _, recipientConfig := range unconditionalRevShareConfig.Configs {
		totalUnconditionalRevSharePpm += recipientConfig.SharePpm
	}
	totalMarketMapperRevSharePpm := marketMapperRevShareParams.RevenueSharePpm

	// return false if totalUnconditionalRevSharePpm + totalMarketMapperRevSharePpm >= 100%
	if totalUnconditionalRevSharePpm+totalMarketMapperRevSharePpm >= lib.OneMillion {
		return false
	}

	bigNetFee := new(big.Int).SetUint64(
		// Casting is safe since both variables are int32.
		uint64(lowestTakerFeePpm) + uint64(lowestMakerFeePpm),
	)

	bigLowestTakerFeePpmMulRevShareRateCap := lib.BigMulPpm(
		lib.BigI(lowestTakerFeePpm),
		lib.BigU(affiliatetypes.AffiliatesRevSharePpmCap),
		true,
	)
	// TODO(OTE-826): Update ValidateRevshareSafety formula and fix tests
	return bigNetFee.Cmp(bigLowestTakerFeePpmMulRevShareRateCap) >= 0
}

func (k Keeper) GetAllRevShares(
	ctx sdk.Context,
	fill clobtypes.FillForProcess,
	affiliatesWhitelistMap map[string]uint32,
) (types.RevSharesForFill, error) {
	revShares := []types.RevShare{}
	feeSourceToQuoteQuantums := make(map[types.RevShareFeeSource]*big.Int)
	feeSourceToRevSharePpm := make(map[types.RevShareFeeSource]uint32)
	feeSourceToQuoteQuantums[types.REV_SHARE_FEE_SOURCE_TAKER_FEE] = big.NewInt(0)
	feeSourceToRevSharePpm[types.REV_SHARE_FEE_SOURCE_TAKER_FEE] = 0
	feeSourceToQuoteQuantums[types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE] = big.NewInt(0)
	feeSourceToRevSharePpm[types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE] = 0

	totalFeesShared := big.NewInt(0)
	takerFees := fill.TakerFeeQuoteQuantums
	makerFees := fill.MakerFeeQuoteQuantums
	netFees := big.NewInt(0).Add(takerFees, makerFees)

	affiliateRevShares, affiliateFeesShared, err := k.getAffiliateRevShares(ctx, fill, affiliatesWhitelistMap)
	if err != nil {
		return types.RevSharesForFill{}, err
	}
	netFeesSubAffiliateFeesShared := new(big.Int).Sub(netFees, affiliateFeesShared)
	if netFeesSubAffiliateFeesShared.Sign() <= 0 {
		return types.RevSharesForFill{}, types.ErrAffiliateFeesSharedExceedsNetFees
	}

	unconditionalRevShares, err := k.getUnconditionalRevShares(ctx, netFeesSubAffiliateFeesShared)
	if err != nil {
		return types.RevSharesForFill{}, err
	}
	marketMapperRevShares, err := k.getMarketMapperRevShare(ctx, fill.MarketId, netFeesSubAffiliateFeesShared)
	if err != nil {
		return types.RevSharesForFill{}, err
	}

	revShares = append(revShares, affiliateRevShares...)
	revShares = append(revShares, unconditionalRevShares...)
	revShares = append(revShares, marketMapperRevShares...)

	var affiliateRevShare *types.RevShare
	if len(affiliateRevShares) > 0 {
		// There should only be one affiliate rev share per fill
		affiliateRevShare = &affiliateRevShares[0]
	}

	for _, revShare := range revShares {
		totalFeesShared.Add(totalFeesShared, revShare.QuoteQuantums)

		// Add the rev share to the total for the fee source
		feeSourceToQuoteQuantums[revShare.RevShareFeeSource].Add(
			feeSourceToQuoteQuantums[revShare.RevShareFeeSource], revShare.QuoteQuantums)

		// Add the rev share ppm to the total for the fee source
		feeSourceToRevSharePpm[revShare.RevShareFeeSource] += revShare.RevSharePpm
	}
	//check total fees shared is less than or equal to net fees
	if totalFeesShared.Cmp(netFees) > 0 {
		return types.RevSharesForFill{}, types.ErrTotalFeesSharedExceedsNetFees
	}

	return types.RevSharesForFill{
		AffiliateRevShare:        affiliateRevShare,
		FeeSourceToQuoteQuantums: feeSourceToQuoteQuantums,
		FeeSourceToRevSharePpm:   feeSourceToRevSharePpm,
		AllRevShares:             revShares,
	}, nil
}

func (k Keeper) getAffiliateRevShares(
	ctx sdk.Context,
	fill clobtypes.FillForProcess,
	affiliatesWhitelistMap map[string]uint32,
) ([]types.RevShare, *big.Int, error) {
	takerAddr := fill.TakerAddr
	takerFee := fill.TakerFeeQuoteQuantums
	if fill.MonthlyRollingTakerVolumeQuantums >= types.MaxReferee30dVolumeForAffiliateShareQuantums {
		return nil, big.NewInt(0), nil
	}

	takerAffiliateAddr, feeSharePpm, exists, err := k.affiliatesKeeper.GetTakerFeeShare(
		ctx, takerAddr, affiliatesWhitelistMap)
	if err != nil {
		return nil, big.NewInt(0), err
	}
	if !exists {
		return nil, big.NewInt(0), nil
	}
	feesShared := lib.BigMulPpm(takerFee, lib.BigU(feeSharePpm), false)
	return []types.RevShare{
		{
			Recipient:         takerAffiliateAddr,
			RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_TAKER_FEE,
			RevShareType:      types.REV_SHARE_TYPE_AFFILIATE,
			QuoteQuantums:     feesShared,
			RevSharePpm:       feeSharePpm,
		},
	}, feesShared, nil
}

func (k Keeper) getUnconditionalRevShares(
	ctx sdk.Context,
	netFeesSubAffiliateFeesShared *big.Int,
) ([]types.RevShare, error) {
	revShares := []types.RevShare{}
	unconditionalRevShareConfig, err := k.GetUnconditionalRevShareConfigParams(ctx)
	if err != nil {
		return nil, err
	}
	for _, revShare := range unconditionalRevShareConfig.Configs {
		feeShared := lib.BigMulPpm(netFeesSubAffiliateFeesShared, lib.BigU(revShare.SharePpm), false)
		revShare := types.RevShare{
			Recipient:         revShare.Address,
			RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
			RevShareType:      types.REV_SHARE_TYPE_UNCONDITIONAL,
			QuoteQuantums:     feeShared,
			RevSharePpm:       revShare.SharePpm,
		}
		revShares = append(revShares, revShare)
	}
	return revShares, nil
}

func (k Keeper) getMarketMapperRevShare(
	ctx sdk.Context,
	marketId uint32,
	netFeesSubAffiliateFeesShared *big.Int,
) ([]types.RevShare, error) {
	revShares := []types.RevShare{}
	marketMapperRevshareAddress, revenueSharePpm, err := k.GetMarketMapperRevenueShareForMarket(ctx, marketId)
	if err != nil {
		return nil, err
	}
	if revenueSharePpm == 0 {
		return nil, nil
	}

	marketMapperRevshareAmount := lib.BigMulPpm(netFeesSubAffiliateFeesShared, lib.BigU(revenueSharePpm), false)
	revShares = append(revShares, types.RevShare{
		Recipient:         marketMapperRevshareAddress.String(),
		RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
		RevShareType:      types.REV_SHARE_TYPE_MARKET_MAPPER,
		QuoteQuantums:     marketMapperRevshareAmount,
		RevSharePpm:       revenueSharePpm,
	})

	return revShares, nil
}
