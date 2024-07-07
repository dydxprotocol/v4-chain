package prepare

import (
	"errors"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
)

// PrepareProposalTxs is used as an intermediary storage for transactions when creating
// a proposal for `PrepareProposal`.
type PrepareProposalTxs struct {
	// Transactions.

	AddPremiumVotesTx    []byte
	ProposedOperationsTx []byte
	OtherTxs             [][]byte
	ExtInfoBz            []byte
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

func (t *PrepareProposalTxs) SetExtInfoBz(extInfoBz []byte) error {
	oldBytes := uint64(len(t.ExtInfoBz))
	newBytes := uint64(len(extInfoBz))
	if err := t.UpdateUsedBytes(oldBytes, newBytes); err != nil {
		return err
	}
	t.ExtInfoBz = extInfoBz
	return nil
}

// SetAddPremiumVotesTx sets the tx used for adding premium votes.
func (t *PrepareProposalTxs) SetAddPremiumVotesTx(tx []byte) error {
	oldBytes := uint64(len(t.AddPremiumVotesTx))
	newBytes := uint64(len(tx))
	if err := t.UpdateUsedBytes(oldBytes, newBytes); err != nil {
		return err
	}
	t.AddPremiumVotesTx = tx
	return nil
}

// SetProposedOperationsTx sets the tx used for order operations.
func (t *PrepareProposalTxs) SetProposedOperationsTx(tx []byte) error {
	oldBytes := uint64(len(t.ProposedOperationsTx))
	newBytes := uint64(len(tx))
	if err := t.UpdateUsedBytes(oldBytes, newBytes); err != nil {
		return err
	}
	t.ProposedOperationsTx = tx
	return nil
}

// AddOtherTxs adds txs to the "other" tx category.
func (t *PrepareProposalTxs) AddOtherTxs(allTxs [][]byte) error {
	bytesToAdd := uint64(0)
	for _, tx := range allTxs {
		txSize := uint64(len(tx))
		if txSize == 0 {
			return fmt.Errorf("Cannot add zero length tx: %v", tx)
		}
		bytesToAdd += txSize
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
func (t *PrepareProposalTxs) GetTxsInOrder(veEnabled bool) ([][]byte, error) {

	if !veEnabled && len(t.ExtInfoBz) > 0 {
		return nil, errors.New("extInfoBz must not be set; VE is disabled")
	}

	if len(t.AddPremiumVotesTx) == 0 {
		return nil, errors.New("AddPremiumVotesTx must be set")
	}

	var txsToReturn [][]byte

	// 1. ve info. it gets included even if empty
	if t.ExtInfoBz != nil {
		txsToReturn = append(txsToReturn, t.ExtInfoBz)
	} else {
		return nil, errors.New("ExtInfoBz must be set")
	}

	// 2. Proposed operations.
	if len(t.ProposedOperationsTx) > 0 {
		txsToReturn = append(txsToReturn, t.ProposedOperationsTx)
	}

	// 3. "Other" txs.
	if len(t.OtherTxs) > 0 {
		txsToReturn = append(txsToReturn, t.OtherTxs...)
	}

	// 4. Funding samples.
	// The validation for `AddPremiumVotesTx` is done at the beginning.
	txsToReturn = append(txsToReturn, t.AddPremiumVotesTx)

	return txsToReturn, nil
}
