package types

import (
	"context"
	"math/big"
	"math/rand"

	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	blocktimetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
	perpetualsmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SubaccountsKeeper interface {
	CanUpdateSubaccounts(
		ctx sdk.Context,
		updates []satypes.Update,
		updateType satypes.UpdateType,
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
		updateType satypes.UpdateType,
	) (
		success bool,
		successPerUpdate []satypes.UpdateResult,
		err error,
	)
	SetNegativeTncSubaccountSeenAtBlock(
		ctx sdk.Context,
		perpetualId uint32,
		blockHeight uint32,
	) error
	TransferFeesToFeeCollectorModule(
		ctx sdk.Context,
		assetId uint32,
		amount *big.Int,
		perpetualId uint32,
	) error
	TransferInsuranceFundPayments(
		ctx sdk.Context,
		amount *big.Int,
		perpetualId uint32,
	) error
	TransferLiquidityFee(
		ctx sdk.Context,
		liquidityFeeQuoteQuantums *big.Int,
		perpetualId uint32,
	) error
	TransferValidatorFee(
		ctx sdk.Context,
		validatorFeeQuoteQuantums *big.Int,
		perpetualId uint32,
	) error
	GetCollateralPoolFromPerpetualId(
		ctx sdk.Context,
		perpetualId uint32,
	) (sdk.AccAddress, error)
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
	GetAllPerpetuals(ctx sdk.Context) (list []perpetualsmoduletypes.Perpetual)
	GetAllLiquidityTiers(ctx sdk.Context) (list []perpetualsmoduletypes.LiquidityTier)
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
	GetInsuranceFundModuleAddress(ctx sdk.Context, perpetualId uint32) (sdk.AccAddress, error)
}

type PricesKeeper interface {
	GetMarketParam(ctx sdk.Context, id uint32) (param pricestypes.MarketParam, exists bool)
	GetAllMarketPrices(ctx sdk.Context) []pricestypes.MarketPrice
}

type StatsKeeper interface {
	RecordFill(ctx sdk.Context, takerAddress string, makerAddress string, notional *big.Int)
}

// AccountKeeper defines the expected account keeper used for simulations.
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

// BankKeeper defines the expected bank keeper used for simulations.
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}
