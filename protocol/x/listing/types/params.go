package types

import "github.com/dydxprotocol/v4-chain/protocol/dtypes"

// DefaultParams defines the default parameters for listing vault deposits.
func DefaultParams() ListingVaultDepositParams {
	return ListingVaultDepositParams{
		NewVaultDepositAmount:  dtypes.NewIntFromUint64(10_000_000_000), // 10_000 USDC
		MainVaultDepositAmount: dtypes.NewIntFromUint64(0),
		NumBlocksToLockShares:  30 * 24 * 3600, // 30 days
	}
}

// Validate checks that the parameters have valid values.
func (p ListingVaultDepositParams) Validate() error {
	// if any of the deposit params are negative, return an error
	if p.NewVaultDepositAmount.Sign() <= 0 || p.MainVaultDepositAmount.Sign() < 0 {
		return ErrInvalidDepositAmount
	}

	// if the number of blocks to lock shares is negative, return an error
	if p.NumBlocksToLockShares <= 0 {
		return ErrInvalidNumBlocksToLockShares
	}

	return nil
}
