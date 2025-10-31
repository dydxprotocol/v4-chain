package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	statstypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
)

// StatsKeeper defines the expected stats keeper
type StatsKeeper interface {
	GetUserStats(ctx sdk.Context, address string) *statstypes.UserStats
	GetGlobalStats(ctx sdk.Context) *statstypes.GlobalStats
	GetStakedBaseTokens(ctx sdk.Context, delegatorAddr string) *big.Int
}

// VaultKeeper defines the expected vault keeper.
type VaultKeeper interface {
	IsVault(ctx sdk.Context, address string) bool
}

// AffiliatesKeeper defines the expected affiliates keeper.
type AffiliatesKeeper interface {
	GetReferredBy(ctx sdk.Context, referee string) (string, bool)
	GetAllAffiliateTiers(ctx sdk.Context) (affiliatetypes.AffiliateTiers, error)
	GetAffiliateParameters(ctx sdk.Context) (affiliatetypes.AffiliateParameters, error)
}

// RevShareKeeper defines the expected revshare keeper.
type RevShareKeeper interface {
	GetUnconditionalRevShareConfigParams(ctx sdk.Context) (revsharetypes.UnconditionalRevShareConfig, error)
	GetMarketMapperRevenueShareParams(
		ctx sdk.Context,
	) revsharetypes.MarketMapperRevenueShareParams
	ValidateRevShareSafety(
		ctx sdk.Context,
		unconditionalRevShareConfig revsharetypes.UnconditionalRevShareConfig,
		marketMapperRevShareParams revsharetypes.MarketMapperRevenueShareParams,
		lowestTakerFee int32,
		lowestMakerFee int32,
	) bool
}
