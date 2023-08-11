package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	pricestypes "github.com/dydxprotocol/v4/x/prices/types"
)

// GetPricePremiumParams includes the parameters used by
// `ClobKeeper.GetPricePremiumForPerpetual` and
// `MemClob.GetPricePremium` to get the price premium.
type GetPricePremiumParams struct {
	Market                      pricestypes.Market
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
		address sdk.AccAddress,
	) (
		msgAddPremiumVotes *MsgAddPremiumVotes,
	)
	PerformStatefulPremiumVotesValidation(
		ctx sdk.Context,
		msg *MsgAddPremiumVotes,
	) (
		err error,
	)
}
