package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

type VaultKeeper interface {
	// Orders.
	GetVaultClobOrders(
		ctx sdk.Context,
		vaultId VaultId,
	) (orders []*clobtypes.Order, err error)
	RefreshAllVaultOrders(ctx sdk.Context)
	RefreshVaultClobOrders(
		ctx sdk.Context,
		vaultId VaultId,
	) (err error)

	// Params.
	GetDefaultQuotingParams(
		ctx sdk.Context,
	) QuotingParams
	SetDefaultQuotingParams(
		ctx sdk.Context,
		params QuotingParams,
	) error

	// Shares.
	GetTotalShares(
		ctx sdk.Context,
		vaultId VaultId,
	) (val NumShares, exists bool)
	SetTotalShares(
		ctx sdk.Context,
		vaultId VaultId,
		totalShares NumShares,
	) error
	MintShares(
		ctx sdk.Context,
		vaultId VaultId,
		owner string,
		quantumsToDeposit *big.Int,
	) error

	// Vault info.
	GetVaultEquity(
		ctx sdk.Context,
		vaultId VaultId,
	) (*big.Int, error)
}
