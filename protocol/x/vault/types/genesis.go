package types

// DefaultGenesis returns the default stats genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Validate params.
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// Validate vaults, ensuring that for each vault:
	// 1. TotalShares is non-negative.
	// 2. OwnerShares is non-negative.
	// 3. TotalShares is equal to the sum of OwnerShares.
	// 4. Owner is not empty.
	for _, vault := range gs.Vaults {
		totalShares := vault.TotalShares.NumShares.BigInt()
		if totalShares.Sign() == -1 {
			return ErrNegativeShares
		}
		for _, ownerShares := range vault.OwnerShares {
			if ownerShares.Owner == "" {
				return ErrInvalidOwner
			} else if ownerShares.Shares.NumShares.Sign() == -1 {
				return ErrNegativeShares
			}
			totalShares.Sub(totalShares, ownerShares.Shares.NumShares.BigInt())
		}
		if totalShares.Sign() != 0 {
			return ErrMismatchedTotalAndOwnerShares
		}
	}
	return nil
}
