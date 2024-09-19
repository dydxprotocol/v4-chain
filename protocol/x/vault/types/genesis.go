package types

// DefaultGenesis returns the default stats genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		DefaultQuotingParams: DefaultQuotingParams(),
		OperatorParams:       DefaultOperatorParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Validate params.
	if err := gs.DefaultQuotingParams.Validate(); err != nil {
		return err
	}

	// Validate shares:
	// 1. TotalShares is non-negative.
	// 2. Each OwnerShares is non-negative.
	// 3. Each Owner is non-empty.
	// 4. TotalShares is equal to the sum of OwnerShares.
	totalShares := gs.TotalShares.NumShares.BigInt()
	if totalShares.Sign() == -1 {
		return ErrNegativeShares
	}
	for _, ownerShares := range gs.OwnerShares {
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

	// Validate vaults, ensuring that for each vault:
	// 1. VaultParams are valid.
	for _, vault := range gs.Vaults {
		if err := vault.VaultParams.Validate(); err != nil {
			return err
		}
	}
	return nil
}
