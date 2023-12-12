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
	GetMarginRequirements(
		ctx sdk.Context,
		id uint32,
		bigQuantums *big.Int,
	) (
		bigInitialMarginQuoteQuantums *big.Int,
		bigMaintenanceMarginQuoteQuantums *big.Int,
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
	) (Perpetual, error)
	ModifyPerpetual(
		ctx sdk.Context,
		id uint32,
		ticker string,
		marketId uint32,
		defaultFundingPpm int32,
		liquidityTier uint32,
	) (Perpetual, error)
	SetLiquidityTier(
		ctx sdk.Context,
		id uint32,
		name string,
		initialMarginPpm uint32,
		maintenanceFractionPpm uint32,
		impactNotional uint64,
	) (
		liquidityTier LiquidityTier,
		err error,
	)
	SetParams(
		ctx sdk.Context,
		params Params,
	) error
}
