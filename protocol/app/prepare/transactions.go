package prepare

import (
	"errors"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
)

// txProtoSize returns the number of bytes a tx contributes to a block's data size, matching how
// CometBFT accounts for it in Txs.Validate (1 + varint(len) + len per tx). Budgeting against the raw
// len(tx) undercounts by the per-tx proto framing, which can push the proposal over MaxTxBytes and
// panic the proposer.
func txProtoSize(tx []byte) uint64 {
	return uint64(cmttypes.ComputeProtoSizeForTxs([]cmttypes.Tx{tx}))
}

// txProtoSizeOrZero is txProtoSize but returns 0 for an unset (empty) slot, which contributes
// nothing to the block, rather than the ~2 bytes of framing txProtoSize would report for empty bytes.
func txProtoSizeOrZero(tx []byte) uint64 {
	if len(tx) == 0 {
		return 0
	}
	return txProtoSize(tx)
}

// PrepareProposalTxs is used as an intermediary storage for transactions when creating
// a proposal for `PrepareProposal`.
type PrepareProposalTxs struct {
	// Transactions.
	UpdateMarketPricesTx []byte
	AddPremiumVotesTx    []byte
	ProposedOperationsTx []byte
	AcknowledgeBridgesTx []byte
	OtherTxs             [][]byte

	// Bytes.
	// In general, there's no need to check for int64 overflow given that it would require
	// exabytes of memory to hit the max int64 value in bytes.
	MaxBytes  uint64
	UsedBytes uint64
}

// NewPrepareProposalTxs returns a new `PrepareProposalTxs` given the request.
func NewPrepareProposalTxs(
	req *abci.RequestPrepareProposal,
) (PrepareProposalTxs, error) {
	if req.MaxTxBytes <= 0 {
		return PrepareProposalTxs{}, errors.New("MaxTxBytes must be positive")
	}

	return PrepareProposalTxs{
		MaxBytes:  uint64(req.MaxTxBytes),
		UsedBytes: 0,
	}, nil
}

// SetUpdateMarketPricesTx sets the tx used for updating market prices.
func (t *PrepareProposalTxs) SetUpdateMarketPricesTx(tx []byte) error {
	oldBytes := txProtoSizeOrZero(t.UpdateMarketPricesTx)
	newBytes := txProtoSizeOrZero(tx)
	if err := t.UpdateUsedBytes(oldBytes, newBytes); err != nil {
		return err
	}
	t.UpdateMarketPricesTx = tx
	return nil
}

// SetAddPremiumVotesTx sets the tx used for adding premium votes.
func (t *PrepareProposalTxs) SetAddPremiumVotesTx(tx []byte) error {
	oldBytes := txProtoSizeOrZero(t.AddPremiumVotesTx)
	newBytes := txProtoSizeOrZero(tx)
	if err := t.UpdateUsedBytes(oldBytes, newBytes); err != nil {
		return err
	}
	t.AddPremiumVotesTx = tx
	return nil
}

// SetProposedOperationsTx sets the tx used for order operations.
func (t *PrepareProposalTxs) SetProposedOperationsTx(tx []byte) error {
	oldBytes := txProtoSizeOrZero(t.ProposedOperationsTx)
	newBytes := txProtoSizeOrZero(tx)
	if err := t.UpdateUsedBytes(oldBytes, newBytes); err != nil {
		return err
	}
	t.ProposedOperationsTx = tx
	return nil
}

// SetAcknowledgeBridgesTx sets the tx used for acknowledging bridges.
func (t *PrepareProposalTxs) SetAcknowledgeBridgesTx(tx []byte) error {
	oldBytes := txProtoSizeOrZero(t.AcknowledgeBridgesTx)
	newBytes := txProtoSizeOrZero(tx)
	if err := t.UpdateUsedBytes(oldBytes, newBytes); err != nil {
		return err
	}
	t.AcknowledgeBridgesTx = tx
	return nil
}

// AddOtherTxs adds txs to the "other" tx category.
func (t *PrepareProposalTxs) AddOtherTxs(allTxs [][]byte) error {
	bytesToAdd := uint64(0)
	for _, tx := range allTxs {
		if len(tx) == 0 {
			return fmt.Errorf("Cannot add zero length tx: %v", tx)
		}
		bytesToAdd += txProtoSize(tx)
	}

	if bytesToAdd == 0 { // no new txs, so return early.
		return errors.New("No txs to add.")
	}

	if err := t.UpdateUsedBytes(0, bytesToAdd); err != nil {
		return err
	}

	t.OtherTxs = append(t.OtherTxs, allTxs...)
	return nil
}

// UpdateUsedBytes updates the used bytes field. This returns an error if the num used bytes
// exceeds the max byte limit.
func (t *PrepareProposalTxs) UpdateUsedBytes(
	bytesToRemove uint64,
	bytesToAdd uint64,
) error {
	if t.UsedBytes < bytesToRemove {
		return errors.New("Result cannot be negative")
	}

	finalBytes := t.UsedBytes - bytesToRemove + bytesToAdd
	if finalBytes > t.MaxBytes {
		return fmt.Errorf("Exceeds max: max=%d, used=%d, adding=%d", t.MaxBytes, t.UsedBytes, bytesToAdd)
	}

	t.UsedBytes = finalBytes
	return nil
}

// GetAvailableBytes returns the available bytes for the proposal.
func (t *PrepareProposalTxs) GetAvailableBytes() uint64 {
	return t.MaxBytes - t.UsedBytes
}

// GetTxsInOrder returns a list of txs in an order that the `ProcessProposal` expects.
func (t *PrepareProposalTxs) GetTxsInOrder() ([][]byte, error) {
	if len(t.UpdateMarketPricesTx) == 0 {
		return nil, errors.New("UpdateMarketPricesTx must be set")
	}

	if len(t.AddPremiumVotesTx) == 0 {
		return nil, errors.New("AddPremiumVotesTx must be set")
	}

	if len(t.AcknowledgeBridgesTx) == 0 {
		return nil, errors.New("AcknowledgeBridgesTx must be set")
	}

	var txsToReturn [][]byte

	// 1. Proposed operations.
	if len(t.ProposedOperationsTx) > 0 {
		txsToReturn = append(txsToReturn, t.ProposedOperationsTx)
	}

	// 2. "Other" txs.
	if len(t.OtherTxs) > 0 {
		txsToReturn = append(txsToReturn, t.OtherTxs...)
	}

	// 3. Acknowledge bridges.
	txsToReturn = append(txsToReturn, t.AcknowledgeBridgesTx)

	// 4. Funding samples.
	// The validation for `AddPremiumVotesTx` is done at the beginning.
	txsToReturn = append(txsToReturn, t.AddPremiumVotesTx)

	// 5. Price updates.
	// The validation for `UpdateMarketPricesTx` is done at the beginning.
	txsToReturn = append(txsToReturn, t.UpdateMarketPricesTx)

	return txsToReturn, nil
}
