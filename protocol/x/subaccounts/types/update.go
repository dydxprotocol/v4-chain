package types

import (
	"math/big"
	"time"

	errorsmod "cosmossdk.io/errors"
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

func (u UpdateResult) IsIsolatedSubaccountError() bool {
	return u == ViolatesIsolatedSubaccountConstraints
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
	Success:                               "Success",
	NewlyUndercollateralized:              "NewlyUndercollateralized",
	StillUndercollateralized:              "StillUndercollateralized",
	WithdrawalsAndTransfersBlocked:        "WithdrawalsAndTransfersBlocked",
	UpdateCausedError:                     "UpdateCausedError",
	ViolatesIsolatedSubaccountConstraints: "ViolatesIsolatedSubaccountConstraints",
}

const (
	Success UpdateResult = iota
	NewlyUndercollateralized
	StillUndercollateralized
	WithdrawalsAndTransfersBlocked
	UpdateCausedError
	ViolatesIsolatedSubaccountConstraints
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

func (u *AssetUpdate) DeepCopy() AssetUpdate {
	return AssetUpdate{
		AssetId:          u.AssetId,
		BigQuantumsDelta: new(big.Int).Set(u.BigQuantumsDelta),
	}
}

type PerpetualUpdate struct {
	// The `Id` of the `Perpetual` for which the `PerpetualPosition` is for.
	PerpetualId uint32
	// The signed change in the `Quantums` of the `PerpetualPosition`
	// represented in base quantums.
	BigQuantumsDelta *big.Int
	// The signed change in the `QuoteBalance` of the `PerpetualPosition`.
	BigQuoteBalanceDelta *big.Int
}

func (u *PerpetualUpdate) DeepCopy() PerpetualUpdate {
	r := PerpetualUpdate{
		PerpetualId:          u.PerpetualId,
		BigQuantumsDelta:     new(big.Int).Set(u.BigQuantumsDelta),
		BigQuoteBalanceDelta: new(big.Int),
	}

	if u.BigQuoteBalanceDelta != nil {
		r.BigQuoteBalanceDelta.Set(u.BigQuoteBalanceDelta)
	}
	return r
}

type UpdateType uint

const (
	UpdateTypeUnspecified UpdateType = iota
	Withdrawal
	Transfer
	Deposit
	Match
	CollatCheck
)

var updateTypeStringMap = map[UpdateType]string{
	UpdateTypeUnspecified: "UpdateTypeUnspecified",
	Withdrawal:            "Withdrawal",
	Transfer:              "Transfer",
	Deposit:               "Deposit",
	Match:                 "Match",
	CollatCheck:           "CollatCheck",
}

func (u UpdateType) String() string {
	result, exists := updateTypeStringMap[u]
	if !exists {
		return "UnexpectedUpdateTypeError"
	}

	return result
}

// WITHDRAWAL_AND_TRANSFER_DELAY_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS defines the number of
// blocks withdrawals and transfers will be blocked if a negative TNC subaccount is seen in state,
// starting from the last block a negative TNC subaccount was seen.
const WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS = 50

const WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_CHAIN_OUTAGE_DURATION = 5 * time.Minute
