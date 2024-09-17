package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RevShareKeeper interface {
	// MarketMapperRevenueShareParams
	SetMarketMapperRevenueShareParams(
		ctx sdk.Context,
		params MarketMapperRevenueShareParams,
	) (err error)
	GetMarketMapperRevenueShareParams(
		ctx sdk.Context,
	) (params MarketMapperRevenueShareParams)

	// MarketMapperRevShareDetails
	SetMarketMapperRevShareDetails(
		ctx sdk.Context,
		marketId uint32,
		params MarketMapperRevShareDetails,
	)
	GetMarketMapperRevShareDetails(
		ctx sdk.Context,
		marketId uint32,
	) (params MarketMapperRevShareDetails, err error)
	CreateNewMarketRevShare(ctx sdk.Context, marketId uint32)
}

type RevShare struct {
	Recipient         string
	RevShareFeeSource RevShareFeeSource
	RevShareType      RevShareType
	QuoteQuantums     *big.Int
	RevSharePpm       uint32
}

type RevShareFeeSource int

const (
	REV_SHARE_FEE_SOURCE_UNSPECIFIED RevShareFeeSource = iota
	REV_SHARE_FEE_SOURCE_NET_FEE
	REV_SHARE_FEE_SOURCE_TAKER_FEE
)

type RevShareType int

const (
	REV_SHARE_TYPE_UNSPECIFIED RevShareType = iota
	REV_SHARE_TYPE_MARKET_MAPPER
	REV_SHARE_TYPE_UNCONDITIONAL
	REV_SHARE_TYPE_AFFILIATE
)

type RevSharesForFill struct {
	AffiliateRevShare        *RevShare
	FeeSourceToQuoteQuantums map[RevShareFeeSource]*big.Int
	FeeSourceToRevSharePpm   map[RevShareFeeSource]uint32
	AllRevShares             []RevShare
}
