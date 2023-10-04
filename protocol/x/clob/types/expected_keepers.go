package types

import (
	"math/big"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	perpetualsmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type SubaccountsKeeper interface {
	CanUpdateSubaccounts(
		ctx sdk.Context,
		updates []satypes.Update,
	) (
		success bool,
		successPerUpdate []satypes.UpdateResult,
		err error,
	)
	GetNetCollateralAndMarginRequirements(
		ctx sdk.Context,
		update satypes.Update,
	) (
		bigNetCollateral *big.Int,
		bigInitialMargin *big.Int,
		bigMaintenanceMargin *big.Int,
		err error,
	)
	GetSubaccount(
		ctx sdk.Context,
		id satypes.SubaccountId,
	) (
		val satypes.Subaccount,
	)
	GetAllSubaccount(
		ctx sdk.Context,
	) (
		list []satypes.Subaccount,
	)
	ForEachSubaccountRandomStart(
		ctx sdk.Context,
		callback func(satypes.Subaccount) (finished bool),
		rand *rand.Rand,
	)
	GetRandomSubaccount(
		ctx sdk.Context,
		rand *rand.Rand,
	) (
		satypes.Subaccount,
		error,
	)
	UpdateSubaccounts(
		ctx sdk.Context,
		updates []satypes.Update,
	) (
		success bool,
		successPerUpdate []satypes.UpdateResult,
		err error,
	)
	TransferFeesToFeeCollectorModule(
		ctx sdk.Context,
		assetId uint32,
		amount *big.Int,
	) error
	TransferInsuranceFundPayments(
		ctx sdk.Context,
		amount *big.Int,
	) error
}

type AssetsKeeper interface {
	GetAsset(ctx sdk.Context, id uint32) (val assettypes.Asset, exists bool)
}

type BlockTimeKeeper interface {
	GetPreviousBlockInfo(ctx sdk.Context) blocktimetypes.BlockInfo
}

type FeeTiersKeeper interface {
	GetPerpetualFeePpm(ctx sdk.Context, address string, isTaker bool) int32
}

type PerpetualsKeeper interface {
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
	GetPerpetual(
		ctx sdk.Context,
		id uint32,
	) (val perpetualsmoduletypes.Perpetual, err error)
	GetPerpetualAndMarketPrice(
		ctx sdk.Context,
		perpetualId uint32,
	) (perpetualsmoduletypes.Perpetual, pricestypes.MarketPrice, error)
	GetSettlementPpm(
		ctx sdk.Context,
		perpetualId uint32,
		quantums *big.Int,
		index *big.Int,
	) (
		bigNetSettlement *big.Int,
		newFundingIndex *big.Int,
		err error,
	)
	MaybeProcessNewFundingTickEpoch(ctx sdk.Context)
}

type StatsKeeper interface {
	RecordFill(ctx sdk.Context, takerAddress string, makerAddress string, notional *big.Int)
}

// AccountKeeper defines the expected account keeper used for simulations.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
}

// BankKeeper defines the expected bank keeper used for simulations.
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

type RewardsKeeper interface {
	AddRewardSharesForFill(
		ctx sdk.Context,
		takerAddress string,
		makerAddress string,
		bigFillQuoteQuantums *big.Int,
		bigTakerFeeQuoteQuantums *big.Int,
		bigMakerFeeQuoteQuantums *big.Int,
	)
}
