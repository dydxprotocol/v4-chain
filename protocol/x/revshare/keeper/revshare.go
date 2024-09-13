package keeper

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
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

// ValidateRevShareSafety roughly checks if the total rev share is valid using the formula below
// highest_affiliate_taker_share + sum(unconditional_rev_shares) + market_mapper_rev_share < 100%
// Note: this is just an estimate as affiliate rev share is based on taker fees, while
// the rest of the rev share is based on net fees.
// TODO(OTE-788): Revisit this formula to ensure accuracy.
func (k Keeper) ValidateRevShareSafety(
	affiliateTiers affiliatetypes.AffiliateTiers,
	unconditionalRevShareConfig types.UnconditionalRevShareConfig,
	marketMapperRevShareParams types.MarketMapperRevenueShareParams,
) bool {
	highestTierRevSharePpm := uint32(0)
	if len(affiliateTiers.Tiers) > 0 {
		highestTierRevSharePpm = affiliateTiers.Tiers[len(affiliateTiers.Tiers)-1].TakerFeeSharePpm
	}
	totalUnconditionalRevSharePpm := uint32(0)
	for _, recipientConfig := range unconditionalRevShareConfig.Configs {
		totalUnconditionalRevSharePpm += recipientConfig.SharePpm
	}
	totalMarketMapperRevSharePpm := marketMapperRevShareParams.RevenueSharePpm

	totalRevSharePpm := totalUnconditionalRevSharePpm + totalMarketMapperRevSharePpm + highestTierRevSharePpm
	return totalRevSharePpm < lib.OneMillion
}

func (k Keeper) GetAllRevShares(
	ctx sdk.Context,
	fill clobtypes.FillForProcess,
) ([]types.RevShare, error) {
	revShares := []types.RevShare{}
	totalFeesShared := big.NewInt(0)
	takerFees := fill.TakerFeeQuoteQuantums
	makerFees := fill.MakerFeeQuoteQuantums
	netFees := big.NewInt(0).Add(takerFees, makerFees)

	affiliateRevShares, err := k.getAffiliateRevShares(ctx, fill)
	if err != nil {
		return nil, err
	}

	unconditionalRevShares, err := k.getUnconditionalRevShares(ctx, netFees)
	if err != nil {
		return nil, err
	}

	marketMapperRevShares, err := k.getMarketMapperRevShare(ctx, fill.MarketId, netFees)
	if err != nil {
		return nil, err
	}

	revShares = append(revShares, affiliateRevShares...)
	revShares = append(revShares, unconditionalRevShares...)
	revShares = append(revShares, marketMapperRevShares...)

	for _, revShare := range revShares {
		totalFeesShared.Add(totalFeesShared, revShare.QuoteQuantums)
	}
	//check total fees shared is less than or equal to net fees
	if totalFeesShared.Cmp(netFees) > 0 {
		return nil, errorsmod.Wrap(types.ErrTotalFeesSharedExceedsNetFees, "total fees shared exceeds net fees")
	}

	return revShares, nil
}

func (k Keeper) getAffiliateRevShares(
	ctx sdk.Context,
	fill clobtypes.FillForProcess,
) ([]types.RevShare, error) {
	takerAddr := fill.TakerAddr
	takerFee := fill.TakerFeeQuoteQuantums
	if fill.MonthlyRollingTakerVolumeQuantums >= types.Max30dRefereeVolumeQuantums {
		return nil, nil
	}

	takerAffiliateAddr, feeSharePpm, exists, err := k.affiliatesKeeper.GetTakerFeeShare(ctx, takerAddr)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
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
	}, nil
}

func (k Keeper) getUnconditionalRevShares(
	ctx sdk.Context,
	netFees *big.Int,
) ([]types.RevShare, error) {
	revShares := []types.RevShare{}
	unconditionalRevShareConfig, err := k.GetUnconditionalRevShareConfigParams(ctx)
	if err != nil {
		return nil, err
	}
	for _, revShare := range unconditionalRevShareConfig.Configs {
		feeShared := lib.BigMulPpm(netFees, lib.BigU(revShare.SharePpm), false)
		revShare := types.RevShare{
			Recipient:         revShare.Address,
			RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_NET_FEE,
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
	netFees *big.Int,
) ([]types.RevShare, error) {
	revShares := []types.RevShare{}
	marketMapperRevshareAddress, revenueSharePpm, err := k.GetMarketMapperRevenueShareForMarket(ctx, marketId)
	if err != nil {
		return nil, err
	}
	if revenueSharePpm == 0 {
		return nil, nil
	}

	marketMapperRevshareAmount := lib.BigMulPpm(netFees, lib.BigU(revenueSharePpm), false)
	revShares = append(revShares, types.RevShare{
		Recipient:         marketMapperRevshareAddress.String(),
		RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_NET_FEE,
		RevShareType:      types.REV_SHARE_TYPE_MARKET_MAPPER,
		QuoteQuantums:     marketMapperRevshareAmount,
		RevSharePpm:       revenueSharePpm,
	})

	return revShares, nil
}
