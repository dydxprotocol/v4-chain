package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
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
	return tsNonce >= referenceTs-MaxTimeInPastMs && tsNonce <= referenceTs+MaxTimeInFutureMs
}

func (k Keeper) ProcessTimestampNonce(ctx sdk.Context, acc sdk.AccountI, tsNonce uint64) error {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.TimestampNonce,
		metrics.Latency,
	)

	blockTs := uint64(ctx.BlockTime().UnixMilli())
	address := acc.GetAddress()

	if !IsValidTimestampNonce(tsNonce, blockTs) {
		k.Logger(ctx).Warn(
			"timestamp nonce not within valid time window",
			"tsNonce", tsNonce,
			"blockTs", blockTs,
			"minAcceptedTs", blockTs-MaxTimeInPastMs,
			"maxAcceptedTs", blockTs+MaxTimeInFutureMs,
		)
		return fmt.Errorf("timestamp nonce %d not within valid time window", tsNonce)
	}
	accountState, found := k.GetAccountState(ctx, address)
	if !found {
		// initialize accountplus state with ts nonce details
		k.SetAccountState(ctx, address, AccountStateFromTimestampNonceDetails(address, tsNonce))
	} else {
		EjectStaleTimestampNonces(&accountState, blockTs)
		tsNonceAccepted := AttemptTimestampNonceUpdate(tsNonce, &accountState)
		if !tsNonceAccepted {
			return fmt.Errorf("timestamp nonce %d rejected", tsNonce)
		}
		k.SetAccountState(ctx, address, accountState)
	}
	return nil
}

// Inplace eject all stale timestamps.
func EjectStaleTimestampNonces(accountState *types.AccountState, referenceTs uint64) {
	tsNonceDetails := &accountState.TimestampNonceDetails
	oldestAllowedTs := referenceTs - MaxTimeInPastMs
	var newTsNonces []uint64
	for _, tsNonce := range tsNonceDetails.TimestampNonces {
		if tsNonce >= oldestAllowedTs {
			newTsNonces = append(newTsNonces, tsNonce)
		} else {
			if tsNonce > tsNonceDetails.MaxEjectedNonce {
				tsNonceDetails.MaxEjectedNonce = tsNonce
			}
		}
	}
	tsNonceDetails.TimestampNonces = newTsNonces
}

// Check if the new tsNonce should be accepted. If satisfies conditions, inplace update AccountState.
// Returns bool indicating if new tsNonce was accepted (update was made).
func AttemptTimestampNonceUpdate(
	tsNonce uint64,
	accountState *types.AccountState,
) bool {
	tsNonceDetails := &accountState.TimestampNonceDetails

	if tsNonce <= tsNonceDetails.MaxEjectedNonce {
		return false
	}

	// Must be unique
	if lib.SliceContains(tsNonceDetails.TimestampNonces, tsNonce) {
		return false
	}

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

// Check if input value is larger than smallest value in arr and return index of the min value. If minimum value has
// duplicates, will return smallest index. index = -1 empty slice.
func isLargerThanSmallestValue(value uint64, values []uint64) (bool, int) {
	if len(values) == 0 {
		return false, -1
	}

	minIndex := 0
	for i, ts := range values {
		if ts < values[minIndex] {
			minIndex = i
		}
	}

	return value > values[minIndex], minIndex
}
