package types

import (
	"math/big"

	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perpetualtypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

type PricesKeeper interface {
	CreateMarket(
		ctx sdk.Context,
		marketParam pricestypes.MarketParam,
		marketPrice pricestypes.MarketPrice,
	) (pricestypes.MarketParam, error)
	AcquireNextMarketID(ctx sdk.Context) uint32
}

type ClobKeeper interface {
	CreatePerpetualClobPair(
		ctx sdk.Context,
		clobPairId uint32,
		perpetualId uint32,
		stepSizeBaseQuantums satypes.BaseQuantums,
		quantumConversionExponent int32,
		subticksPerTick uint32,
		status clobtypes.ClobPair_Status,
	) (clobtypes.ClobPair, error)
	AcquireNextClobPairID(ctx sdk.Context) uint32
	ValidateClobPairCreation(ctx sdk.Context, clobPair *clobtypes.ClobPair) error
	CreateClobPairStructures(ctx sdk.Context, clobPair clobtypes.ClobPair) error
	SetClobPair(ctx sdk.Context, clobPair clobtypes.ClobPair)
}

type MarketMapKeeper interface {
	GetMarket(
		ctx sdk.Context,
		ticker string,
	) (marketmaptypes.Market, error)
	// Only used for testing purposes
	CreateMarket(
		ctx sdk.Context,
		market marketmaptypes.Market,
	) error
}

type PerpetualsKeeper interface {
	CreatePerpetual(
		ctx sdk.Context,
		id uint32,
		ticker string,
		marketId uint32,
		atomicResolution int32,
		defaultFundingPpm int32,
		liquidityTier uint32,
		marketType perpetualtypes.PerpetualMarketType,
	) (perpetualtypes.Perpetual, error)
	AcquireNextPerpetualID(ctx sdk.Context) uint32
	GetAllPerpetuals(ctx sdk.Context) (list []perpetualtypes.Perpetual)
}

type VaultKeeper interface {
	DepositToMegavault(
		ctx sdk.Context,
		fromSubaccount satypes.SubaccountId,
		quoteQuantums *big.Int,
	) (mintedShares *big.Int, err error)
	AllocateToVault(
		ctx sdk.Context,
		vaultId vaulttypes.VaultId,
		quantums *big.Int,
	) error
	LockShares(
		ctx sdk.Context,
		ownerAddress string,
		sharesToLock vaulttypes.NumShares,
		tilBlock uint32,
	) error
	SetVaultStatus(
		ctx sdk.Context,
		vaultId vaulttypes.VaultId,
		status vaulttypes.VaultStatus,
	) error
}
