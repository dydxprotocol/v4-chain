package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Subaccounts: []Subaccount{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	includedAccounts := make(map[SubaccountId]bool)
	for _, sa := range gs.Subaccounts {
		subaccountId := sa.GetId()
		if err := subaccountId.Validate(); err != nil {
			return err
		}
		if includedAccounts[*subaccountId] {
			return errorsmod.Wrapf(ErrDuplicateSubaccountIds,
				"duplicate subaccount id %+v found within genesis state", subaccountId)
		}
		includedAccounts[*subaccountId] = true

		// Validate AssetPositions.
		// TODO(DEC-582): once we support different assets, remove this validation.
		if len(sa.GetAssetPositions()) > 1 {
			return ErrMultAssetPositionsNotSupported
		}
		for i := 0; i < len(sa.GetAssetPositions()); i++ {
			assetP := sa.GetAssetPositions()[i]
			if i > 0 && assetP.AssetId <= sa.GetAssetPositions()[i-1].AssetId {
				return ErrAssetPositionsOutOfOrder
			}
			// TODO(DEC-582): once we support different assets, remove this validation.
			if assetP.AssetId != 0 {
				return ErrAssetPositionNotSupported
			}
			if assetP.GetBigQuantums().Sign() == 0 {
				return ErrAssetPositionZeroQuantum
			}
		}

		// Validate PerpetualPositions.
		for i := 0; i < len(sa.GetPerpetualPositions()); i++ {
			perpP := sa.GetPerpetualPositions()[i]
			if i > 0 && perpP.PerpetualId <= sa.GetPerpetualPositions()[i-1].PerpetualId {
				return ErrPerpPositionsOutOfOrder
			}
			if perpP.GetBigQuantums().Sign() == 0 {
				return ErrPerpPositionZeroQuantum
			}
		}
	}
	return nil
}
