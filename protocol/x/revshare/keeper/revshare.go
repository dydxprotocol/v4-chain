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

func (k Keeper) GetOrderRouterRevShare(ctx sdk.Context, orderRouterAddr string) (uint32, error) {
	if orderRouterAddr == "" {
		return 0, nil
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.OrderRouterRevSharePrefix))
	orderRouterBech32Addr, err := sdk.AccAddressFromBech32(orderRouterAddr)
	if err != nil {
		return 0, types.ErrInvalidAddress
	}

	orderRouterRevShareBytes := store.Get(
		[]byte(orderRouterBech32Addr),
	)

	if orderRouterRevShareBytes == nil {
		return 0, types.ErrOrderRouterRevShareNotFound
	}

	return lib.BytesToUint32(orderRouterRevShareBytes), nil
}

func (k Keeper) SetOrderRouterRevShare(ctx sdk.Context, orderRouterAddr string, revSharePpm uint32) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.OrderRouterRevSharePrefix))
	orderRouterBech32Addr, err := sdk.AccAddressFromBech32(orderRouterAddr)
	if err != nil {
		return types.ErrInvalidAddress
	}

	store.Set([]byte(orderRouterBech32Addr), lib.Uint32ToKey(revSharePpm))
	return nil
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
	affiliateOverrides map[string]bool,
	affiliateParameters affiliatetypes.AffiliateParameters,
) (types.RevSharesForFill, error) {
	revShares := []types.RevShare{}
	feeSourceToQuoteQuantums := make(map[types.RevShareFeeSource]*big.Int)
	feeSourceToRevSharePpm := make(map[types.RevShareFeeSource]uint32)
	feeSourceToQuoteQuantums[types.REV_SHARE_FEE_SOURCE_TAKER_FEE] = big.NewInt(0)
	feeSourceToRevSharePpm[types.REV_SHARE_FEE_SOURCE_TAKER_FEE] = 0
	feeSourceToQuoteQuantums[types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE] = big.NewInt(0)
	feeSourceToRevSharePpm[types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE] = 0
	feeSourceToQuoteQuantums[types.REV_SHARE_FEE_SOURCE_MAKER_FEE] = big.NewInt(0)
	feeSourceToRevSharePpm[types.REV_SHARE_FEE_SOURCE_MAKER_FEE] = 0

	totalFeesShared := big.NewInt(0)
	takerFees := fill.TakerFeeQuoteQuantums
	makerFees := fill.MakerFeeQuoteQuantums
	netFees := big.NewInt(0).Add(takerFees, makerFees)

	// when net fee is zero, no rev share is generated from the fill
	if netFees.Sign() == 0 {
		return types.RevSharesForFill{}, nil
	}

	affiliateRevShares, affiliateFeesShared, err := k.getAffiliateRevShares(
		ctx, fill, affiliateOverrides, affiliateParameters)
	if err != nil {
		return types.RevSharesForFill{}, err
	}

	var orderRouterRevShares []types.RevShare
	netFeesSubRevenueShare := new(big.Int).Set(netFees)
	// No affiliate fees shared, so we can generate order router rev shares
	// In the case that the taker has an affiliate fee and the maker does not, then no order router fees are generated
	// for the maker or the taker
	if len(affiliateRevShares) == 0 {
		orderRouterRevShares = k.getOrderRouterRevShares(ctx, fill, takerFees, makerFees)
		for _, revShare := range orderRouterRevShares {
			netFeesSubRevenueShare.Sub(netFeesSubRevenueShare, revShare.QuoteQuantums)
		}
	} else {
		netFeesSubRevenueShare.Sub(netFeesSubRevenueShare, affiliateFeesShared)
	}

	if netFeesSubRevenueShare.Sign() <= 0 {
		return types.RevSharesForFill{}, types.ErrAffiliateFeesSharedGreaterThanOrEqualToNetFees
	}

	unconditionalRevShares, err := k.getUnconditionalRevShares(ctx, netFeesSubRevenueShare)
	if err != nil {
		return types.RevSharesForFill{}, err
	}
	marketMapperRevShares, err := k.getMarketMapperRevShare(ctx, fill.MarketId, netFeesSubRevenueShare)
	if err != nil {
		return types.RevSharesForFill{}, err
	}

	revShares = append(revShares, affiliateRevShares...)
	if orderRouterRevShares != nil {
		revShares = append(revShares, orderRouterRevShares...)
	}
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

	// Check total fees shared is less than or equal to net fees
	if totalFeesShared.Cmp(netFees) > 0 {
		k.Logger(ctx).Error(
			"Total fees exceed net fees. Total fees: ", totalFeesShared,
			"Net fees: ", netFees,
			"Revshares generated: ", revShares)
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
	affiliateOverrides map[string]bool,
	affiliateParams affiliatetypes.AffiliateParameters,
) ([]types.RevShare, *big.Int, error) {
	takerAddr := fill.TakerAddr
	takerFee := fill.TakerFeeQuoteQuantums
	if fill.MonthlyRollingTakerVolumeQuantums >= types.MaxReferee30dVolumeForAffiliateShareQuantums ||
		takerFee.Sign() == 0 {
		return nil, big.NewInt(0), nil
	}

	userStats := k.statsKeeper.GetUserStats(ctx, takerAddr)
	if userStats != nil {
		// If the affiliate revenue generated is greater than the maximum 30d attributable volume
		// per referred user notional, then no affiliate rev share is generated
		// Disable this check if it is 0
		cap := affiliateParams.Maximum_30DAffiliateRevenuePerReferredUserQuoteQuantums
		if cap != 0 &&
			userStats.Affiliate_30DRevenueGeneratedQuantums >= cap {
			// Exceeded revenue cap, no rev share is attributed
			return []types.RevShare{}, big.NewInt(0), nil
		}
	}

	takerAffiliateAddr, feeSharePpm, exists, err := k.affiliatesKeeper.GetTakerFeeShare(
		ctx,
		takerAddr,
		affiliateOverrides,
	)
	if err != nil {
		return nil, big.NewInt(0), err
	}
	if !exists {
		return nil, big.NewInt(0), nil
	}
	feesShared := lib.BigMulPpm(takerFee, lib.BigU(feeSharePpm), false)

	// Cap the affiliate revenue share if it exceeds the maximum 30d affiliate revenue per referred user
	if userStats != nil {
		cap := affiliateParams.Maximum_30DAffiliateRevenuePerReferredUserQuoteQuantums
		if cap != 0 {
			revenueGenerated := userStats.Affiliate_30DRevenueGeneratedQuantums
			// We know revenueGenerated < cap here because of the check above
			maxFee := new(big.Int).SetUint64(cap - revenueGenerated)
			if feesShared.Cmp(maxFee) > 0 {
				feesShared = maxFee
			}
		}
	}

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

func (k Keeper) getOrderRouterRevShares(
	ctx sdk.Context,
	fill clobtypes.FillForProcess,
	takerFees *big.Int,
	makerFees *big.Int,
) []types.RevShare {
	if fill.TakerOrderRouterAddr == "" && fill.MakerOrderRouterAddr == "" {
		return []types.RevShare{}
	}

	orderRouterRevShares := []types.RevShare{}
	takerOrderRouterRevSharePpm, err := k.GetOrderRouterRevShare(ctx, fill.TakerOrderRouterAddr)
	if err != nil {
		k.Logger(ctx).Warn("order router rev share invalid for taker, ignoring ",
			"taker_order_router_addr: ", fill.TakerOrderRouterAddr,
			"taker_addr: ", fill.TakerAddr,
			"error: ", err,
		)
	} else {
		if fill.TakerOrderRouterAddr != "" {
			// Orders can have 2 rev share ids, we need to calculate each side separately
			// This is taker ppm * min(taker, taker - maker_rebate)
			takerFeesSide := lib.BigMin(takerFees, new(big.Int).Add(takerFees, makerFees))
			takerRevShare := lib.BigMulPpm(lib.BigU(takerOrderRouterRevSharePpm), takerFeesSide, false)
			orderRouterRevShares = append(orderRouterRevShares, types.RevShare{
				Recipient:         fill.TakerOrderRouterAddr,
				RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_TAKER_FEE,
				RevShareType:      types.REV_SHARE_TYPE_ORDER_ROUTER,
				QuoteQuantums:     takerRevShare,
				RevSharePpm:       takerOrderRouterRevSharePpm,
			})
		}
	}

	makerOrderRouterRevSharePpm, err := k.GetOrderRouterRevShare(ctx, fill.MakerOrderRouterAddr)
	if err != nil {
		k.Logger(ctx).Warn("order router rev share invalid for maker, ignoring ",
			"maker_order_router_addr: ", fill.MakerOrderRouterAddr,
			"maker_addr: ", fill.MakerAddr,
			"error: ", err,
		)
	} else {
		if fill.MakerOrderRouterAddr != "" {
			// maker ppm * max(0, maker)
			makerFeeSide := lib.BigMax(lib.BigI(0), makerFees)
			makerRevShare := lib.BigMulPpm(makerFeeSide,
				lib.BigU(makerOrderRouterRevSharePpm),
				false,
			)

			orderRouterRevShares = append(orderRouterRevShares, types.RevShare{
				Recipient:         fill.MakerOrderRouterAddr,
				RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_MAKER_FEE,
				RevShareType:      types.REV_SHARE_TYPE_ORDER_ROUTER,
				QuoteQuantums:     makerRevShare,
				RevSharePpm:       makerOrderRouterRevSharePpm,
			})
		}
	}

	return orderRouterRevShares
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
