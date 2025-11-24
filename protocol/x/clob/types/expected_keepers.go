package types

import (
	"context"
	"math/big"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/margin"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	perpetualsmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	stattypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
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
		risk margin.Risk,
		err error,
	)
	GetSubaccount(
		ctx sdk.Context,
		id satypes.SubaccountId,
	) (
		val satypes.Subaccount,
	)
	GetStreamSubaccountUpdate(
		ctx sdk.Context,
		id satypes.SubaccountId,
		snapshot bool,
	) (
		val satypes.StreamSubaccountUpdate,
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
	TransferInsuranceFundPayments(
		ctx sdk.Context,
		amount *big.Int,
		perpetualId uint32,
	) error
	TransferBuilderFees(
		ctx sdk.Context,
		productId uint32,
		builderFeeQuantums *big.Int,
		builderAddress string,
	) error
	GetInsuranceFundBalance(
		ctx sdk.Context,
		perpetualId uint32,
	) (
		balance *big.Int,
	)
	GetCrossInsuranceFundBalance(
		ctx sdk.Context,
	) (
		balance *big.Int,
	)
	GetCollateralPoolFromPerpetualId(
		ctx sdk.Context,
		perpetualId uint32,
	) (sdk.AccAddress, error)
	DistributeFees(
		ctx sdk.Context,
		assetId uint32,
		revSharesForFill revsharetypes.RevSharesForFill,
		fillForProcess FillForProcess,
	) error

	// Leverage methods
	SetLeverage(
		ctx sdk.Context,
		subaccountId *satypes.SubaccountId,
		leverageMap map[uint32]uint32,
	)
	GetLeverage(
		ctx sdk.Context,
		subaccountId *satypes.SubaccountId,
	) (map[uint32]uint32, bool)
	UpdateLeverage(
		ctx sdk.Context,
		subaccountId *satypes.SubaccountId,
		perpetualLeverage map[uint32]uint32,
	) error
	GetMinImfForPerpetual(
		ctx sdk.Context,
		perpetualId uint32,
	) (uint32, error)
}

type AssetsKeeper interface {
	GetAsset(ctx sdk.Context, id uint32) (val assettypes.Asset, exists bool)
}

type BlockTimeKeeper interface {
	GetPreviousBlockInfo(ctx sdk.Context) blocktimetypes.BlockInfo
}

type FeeTiersKeeper interface {
	GetPerpetualFeePpm(ctx sdk.Context, address string, isTaker bool, feeTierOverrideIdx uint32, clobPairId uint32) int32
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
	GetPerpetualAndMarketPriceAndLiquidityTier(
		ctx sdk.Context,
		perpetualId uint32,
	) (
		perpetual perpetualsmoduletypes.Perpetual,
		price pricestypes.MarketPrice,
		liquidityTier perpetualsmoduletypes.LiquidityTier,
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
	MaybeProcessNewFundingTickEpoch(ctx sdk.Context)
	GetInsuranceFundModuleAddress(ctx sdk.Context, perpetualId uint32) (sdk.AccAddress, error)
}

type PricesKeeper interface {
	GetMarketParam(ctx sdk.Context, id uint32) (param pricestypes.MarketParam, exists bool)
	GetStreamPriceUpdate(ctx sdk.Context, id uint32, snapshot bool) (val pricestypes.StreamPriceUpdate)
}

type StatsKeeper interface {
	RecordFill(
		ctx sdk.Context,
		takerAddress string,
		makerAddress string,
		notional *big.Int,
		affiliateFeeGenerated *big.Int,
		affiliateAttributions []*stattypes.AffiliateAttribution,
	)
	GetUserStats(ctx sdk.Context, address string) *stattypes.UserStats
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

type RewardsKeeper interface {
	AddRewardSharesForFill(
		ctx sdk.Context,
		fill FillForProcess,
		revSharesForFill revsharetypes.RevSharesForFill,
	)
}

type RevShareKeeper interface {
	GetAllRevShares(
		ctx sdk.Context,
		fill FillForProcess,
		affiliateOverrides map[string]bool,
		affiliateParameters affiliatetypes.AffiliateParameters,
	) (
		revsharetypes.RevSharesForFill, error,
	)
	GetOrderRouterRevShare(ctx sdk.Context, orderRouterRevShare string) (
		uint32, error)
}

type AffiliatesKeeper interface {
	GetAffiliateWhitelistMap(ctx sdk.Context) (map[string]uint32, error)
	GetAffiliateOverridesMap(ctx sdk.Context) (map[string]bool, error)
	GetAffiliateParameters(ctx sdk.Context) (affiliatetypes.AffiliateParameters, error)
	GetReferredBy(ctx sdk.Context, referee string) (referrer string, found bool)
}

type AccountPlusKeeper interface {
	MaybeValidateAuthenticators(ctx sdk.Context, tx sdk.Tx) error
}
