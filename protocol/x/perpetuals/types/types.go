package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// GetPricePremiumParams includes the parameters used by
// `ClobKeeper.GetPricePremiumForPerpetual` and
// `MemClob.GetPricePremium` to get the price premium.
type GetPricePremiumParams struct {
	IndexPrice                  pricestypes.MarketPrice
	BaseAtomicResolution        int32
	QuoteAtomicResolution       int32
	ImpactNotionalQuoteQuantums *big.Int
	MaxAbsPremiumVotePpm        *big.Int
}

// Interface used by ABCI calls to access the perpetuals keeper.
type PerpetualsKeeper interface {
	MaybeProcessNewFundingTickEpoch(ctx sdk.Context)
	MaybeProcessNewFundingSampleEpoch(ctx sdk.Context)
	AddPremiumVotes(ctx sdk.Context, votes []FundingPremium) error
	GetNetNotional(
		ctx sdk.Context,
		id uint32,
		bigQuantums *big.Int,
	) (
		bigNetNotionalQuoteQuantums *big.Int,
		err error,
	)
	GetNotionalInBaseQuantums(
		ctx sdk.Context,
		id uint32,
		bigQuoteQuantums *big.Int,
	) (
		bigBaseQuantums *big.Int,
		err error,
	)
	GetNetCollateral(
		ctx sdk.Context,
		id uint32,
		bigQuantums *big.Int,
	) (
		bigNetCollateralQuoteQuantums *big.Int,
		err error,
	)
	GetAddPremiumVotes(
		ctx sdk.Context,
	) (
		msgAddPremiumVotes *MsgAddPremiumVotes,
	)
	PerformStatefulPremiumVotesValidation(
		ctx sdk.Context,
		msg *MsgAddPremiumVotes,
	) (
		err error,
	)
	HasAuthority(authority string) bool
	CreatePerpetual(
		ctx sdk.Context,
		id uint32,
		ticker string,
		marketId uint32,
		atomicResolution int32,
		defaultFundingPpm int32,
		liquidityTier uint32,
		marketType PerpetualMarketType,
	) (Perpetual, error)
	ModifyPerpetual(
		ctx sdk.Context,
		id uint32,
		ticker string,
		marketId uint32,
		defaultFundingPpm int32,
		liquidityTier uint32,
	) (Perpetual, error)
	ModifyOpenInterest(
		ctx sdk.Context,
		perpetualId uint32,
		openInterestDeltaBaseQuantums *big.Int,
	) (
		err error,
	)
	SetLiquidityTier(
		ctx sdk.Context,
		id uint32,
		name string,
		initialMarginPpm uint32,
		maintenanceFractionPpm uint32,
		impactNotional uint64,
		openInterestLowerCap uint64,
		openInterestUpperCap uint64,
	) (
		liquidityTier LiquidityTier,
		err error,
	)
	SetParams(
		ctx sdk.Context,
		params Params,
	) error
	SetPerpetualMarketType(
		ctx sdk.Context,
		id uint32,
		marketType PerpetualMarketType,
	) (Perpetual, error)
	GetPerpetual(
		ctx sdk.Context,
		id uint32,
	) (Perpetual, error)
	GetAllPerpetuals(
		ctx sdk.Context,
	) []Perpetual
	GetAllLiquidityTiers(ctx sdk.Context) (list []LiquidityTier)
	ValidateAndSetPerpetual(
		ctx sdk.Context,
		perpetual Perpetual,
	) error
	SetNextPerpetualID(ctx sdk.Context, nextID uint32)
}

// OpenInterestDelta represents a (perpId, openInterestDelta) tuple.
type OpenInterestDelta struct {
	// The `Id` of the `Perpetual`.
	PerpetualId uint32
	// Delta of open interest (in base quantums).
	BaseQuantums *big.Int
}
