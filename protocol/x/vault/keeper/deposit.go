package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// DepositToMegavault deposits from a subaccount to megavault by
// 1. Minting shares for owner address of `fromSubaccount`.
// 2. Transferring `quoteQuantums` from `fromSubaccount` to megavault subaccount 0.
func (k Keeper) DepositToMegavault(
	ctx sdk.Context,
	fromSubaccount satypes.SubaccountId,
	quoteQuantums *big.Int,
) (mintedShares *big.Int, err error) {
	// Mint shares.
	mintedShares, err = k.MintShares(
		ctx,
		fromSubaccount.Owner,
		quoteQuantums,
	)
	if err != nil {
		return nil, err
	}

	// Transfer from sender subaccount to megavault.
	// Note: Transfer should take place after minting shares for
	// shares calculation to be correct.
	err = k.sendingKeeper.ProcessTransfer(
		ctx,
		&sendingtypes.Transfer{
			Sender:    fromSubaccount,
			Recipient: types.MegavaultMainSubaccount,
			AssetId:   assettypes.AssetUsdc.Id,
			Amount:    quoteQuantums.Uint64(),
		},
	)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		types.NewDepositToMegavaultEvent(
			fromSubaccount.Owner,
			quoteQuantums.Uint64(),
			mintedShares.Uint64(),
		),
	)

	return mintedShares, nil
}

// MintShares mints shares for `owner` based on `quantumsToDeposit` by:
// 1. Increasing total shares.
// 2. Increasing owner shares for given `owner`.
func (k Keeper) MintShares(
	ctx sdk.Context,
	owner string,
	quantumsToDeposit *big.Int,
) (mintedShares *big.Int, err error) {
	// Quantums to deposit should be positive.
	if quantumsToDeposit.Sign() <= 0 {
		return nil, types.ErrInvalidQuoteQuantums
	}
	// Get existing TotalShares of the vault.
	existingTotalShares := k.GetTotalShares(ctx).NumShares.BigInt()
	// Calculate shares to mint.
	var sharesToMint *big.Int
	if existingTotalShares.Sign() <= 0 {
		// Mint `quoteQuantums` number of shares.
		sharesToMint = new(big.Int).Set(quantumsToDeposit)
	} else {
		// Get megavault equity.
		equity, err := k.GetMegavaultEquity(ctx)
		if err != nil {
			return nil, err
		}
		// Don't mint shares if equity is non-positive.
		if equity.Sign() <= 0 {
			return nil, types.ErrNonPositiveEquity
		}
		// Mint `deposit (in quote quantums) * existing shares / vault equity (in quote quantums)`
		// number of shares.
		// For example:
		// - a vault currently has 5000 shares and 4000 equity (in quote quantums)
		// - each quote quantum is worth 5000 / 4000 = 1.25 shares
		// - a deposit of 1000 quote quantums should thus be given 1000 * 1.25 = 1250 shares
		sharesToMint = new(big.Int).Set(quantumsToDeposit)
		sharesToMint = sharesToMint.Mul(sharesToMint, existingTotalShares)
		sharesToMint = sharesToMint.Quo(sharesToMint, equity)

		// Return error if `sharesToMint` is rounded down to 0.
		if sharesToMint.Sign() == 0 {
			return nil, types.ErrZeroSharesToMint
		}
	}

	// Increase total shares.
	err = k.SetTotalShares(
		ctx,
		types.BigIntToNumShares(
			existingTotalShares.Add(existingTotalShares, sharesToMint),
		),
	)
	if err != nil {
		return nil, err
	}

	// Increase owner shares.
	ownerShares, exists := k.GetOwnerShares(ctx, owner)
	if !exists {
		// Set owner shares to be sharesToMint.
		err := k.SetOwnerShares(
			ctx,
			owner,
			types.BigIntToNumShares(sharesToMint),
		)
		if err != nil {
			return nil, err
		}
	} else {
		// Increase existing owner shares by sharesToMint.
		existingOwnerShares := ownerShares.NumShares.BigInt()
		err = k.SetOwnerShares(
			ctx,
			owner,
			types.BigIntToNumShares(
				existingOwnerShares.Add(existingOwnerShares, sharesToMint),
			),
		)
		if err != nil {
			return nil, err
		}
	}

	return sharesToMint, nil
}
