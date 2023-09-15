package types

import (
	errorsmod "cosmossdk.io/errors"
	"math/big"
)

type UpdateResult uint

func (u UpdateResult) String() string {
	result, exists := updateResultStringMap[u]
	if !exists {
		return "UnexpectedError"
	}

	return result
}

// IsSuccess returns true if the `UpdateResult` value is `Success`.
func (u UpdateResult) IsSuccess() bool {
	return u == Success
}

// GetErrorFromUpdateResults generates a helpful error when UpdateSubaccounts or
// CanUpdateSubaccounts returns one or more failed updates.
func GetErrorFromUpdateResults(
	success bool,
	successPerUpdate []UpdateResult,
	updates []Update,
) error {
	if success {
		return nil
	}

	for index, result := range successPerUpdate {
		if !result.IsSuccess() {
			subaccountId := updates[index].SubaccountId
			return errorsmod.Wrapf(
				ErrFailedToUpdateSubaccounts,
				"Subaccount with id %v failed with UpdateResult: %v",
				subaccountId,
				result,
			)
		}
	}

	// Should not reach here since successPerUpdate must contains a failure,
	// if success = false.
	panic("UpdateSubaccounts/CanUpdateSubaccounts returns success, but UpdateResults contains failure")
}

var updateResultStringMap = map[UpdateResult]string{
	0: "Success",
	1: "NewlyUndercollateralized",
	2: "StillUndercollateralized",
	3: "UpdateCausedError",
}

const (
	Success UpdateResult = iota
	NewlyUndercollateralized
	StillUndercollateralized
	UpdateCausedError
)

// Update is used by the subaccounts keeper to allow other modules
// to specify changes to one or more `Subaccounts` (for example the
// result of a trade, transfer, etc)
type Update struct {
	// The `Id` of the `Subaccount` for which this update applies.
	SubaccountId SubaccountId
	// A list of changes to make to any `AssetPositions` in the `Subaccount`.
	AssetUpdates []AssetUpdate
	// A list of changes to make to any `PerpetualPositions` in the `Subaccount`.
	PerpetualUpdates []PerpetualUpdate
}

type AssetUpdate struct {
	// The `Id` of the `Asset` for which the `AssetPosition` is for.
	AssetId uint32
	// The signed change in the Size of the `AssetPosition`.
	BigQuantumsDelta *big.Int
}

type PerpetualUpdate struct {
	// The `Id` of the `Perpetual` for which the `PerpetualPosition` is for.
	PerpetualId uint32
	// The signed change in the `Quantums` of the `PerpetualPosition`
	// represented in base quantums.
	BigQuantumsDelta *big.Int
}
