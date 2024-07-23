package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

const TimestampNonceSequenceCutoff uint64 = 1 << 40 // 2^40
const MaxTimestampNonceArrSize = 20
const MaxTimeInPastMs = 30000
const MaxTimeInFutureMs = 30000

func IsTimestampNonce(ts uint64) bool {
	return ts >= TimestampNonceSequenceCutoff
}

func IsValidTimestampNonce(tsNonce uint64, referenceTs uint64) bool {
	return tsNonce >= referenceTs-MaxTimeInPastMs && tsNonce <= referenceTs-MaxTimeInFutureMs
}

// Inplace eject all stale timestamps.
func EjectStaleTsNonces(accountState *types.AccountState, referenceTs uint64) {
	tsNonceDetails := &accountState.TimestampNonceDetails
	var newTsNonces []uint64
	for _, tsNonce := range tsNonceDetails.TimestampNonces {
		if tsNonce >= referenceTs-MaxTimeInPastMs {
			newTsNonces = append(newTsNonces, tsNonce)
		} else {
			if tsNonce > tsNonceDetails.MaxEjectedNonce {
				tsNonceDetails.MaxEjectedNonce = tsNonce
			}
		}
	}
	tsNonceDetails.TimestampNonces = newTsNonces
}

// Check if the new tsNonce should be accepted.
// Inplace update TimestampNonceDetails with the new ts nonce if necessary.
// Returns bool indicating if new tsNonce was accepted (update was made).
func ValidateAndUpdateTimestampNonceDetails(
	tsNonce uint64,
	referenceTs uint64,
	accountState types.AccountState,
) bool {
	tsNonceDetails := &accountState.TimestampNonceDetails
	if len(tsNonceDetails.TimestampNonces) < MaxTimestampNonceArrSize {
		tsNonceDetails.TimestampNonces = append(tsNonceDetails.TimestampNonces, tsNonce)
		return true
	}

	isSufficientlyLargeTsNonce, minIdx := isLargerThanSmallestValue(tsNonce, tsNonceDetails.TimestampNonces)
	if isSufficientlyLargeTsNonce {
		tsNonceDetails.MaxEjectedNonce = tsNonceDetails.TimestampNonces[minIdx]
		tsNonceDetails.TimestampNonces[minIdx] = tsNonce
		return true
	}

	return false
}

// Check if input value is larger than smallest value in arr and return index. index = -1 if false.
func isLargerThanSmallestValue(value uint64, values []uint64) (bool, int) {
	minIndex := -1

	if len(values) == 0 {
		return false, minIndex
	}

	for i, ts := range values {
		if ts < values[minIndex] {
			minIndex = i
		}
	}

	if value > values[minIndex] {
		return true, minIndex
	} else {
		return false, minIndex
	}
}
