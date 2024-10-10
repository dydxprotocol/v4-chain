package process

import (
	"slices"

	errorsmod "cosmossdk.io/errors"
	constants "github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func init() {
	txIndicesAndOffsets := []int{
		constants.ProposedOperationsTxIndex,
		constants.AddPremiumVotesTxLenOffset,
	}
	if constants.MinTxsCount != len(txIndicesAndOffsets) {
		panic("minTxsCount does not match expected count of Txs.")
	}
	if lib.ContainsDuplicates(txIndicesAndOffsets) {
		panic("Duplicate indices/offsets defined for Txs.")
	}
	if slices.Min[[]int](txIndicesAndOffsets) != constants.LastOtherTxLenOffset {
		panic("lastTxLenOffset is not the lowest offset")
	}
	if slices.Max[[]int](txIndicesAndOffsets)+1 != constants.FirstOtherTxIndex {
		panic("firstOtherTxIndex is <= the maximum offset")
	}
	txIndicesForMinTxsCount := []int{
		constants.ProposedOperationsTxIndex,
		constants.AddPremiumVotesTxLenOffset + constants.MinTxsCount,
	}
	if constants.MinTxsCount != len(txIndicesForMinTxsCount) {
		panic("minTxsCount does not match expected count of Txs.")
	}
	if lib.ContainsDuplicates(txIndicesForMinTxsCount) {
		panic("Overlapping indices and offsets defined for Txs.")
	}
	if constants.MinTxsCount != constants.FirstOtherTxIndex-constants.LastOtherTxLenOffset {
		panic("Unexpected gap between firstOtherTxIndex and lastOtherTxLenOffset which is greater than minTxsCount")
	}
}

// TODO: add extInfo into this and use to decode (cleanup)
// ProcessProposalTxs is used as an intermediary struct to validate a proposed list of txs
// for `ProcessProposal`.
type ProcessProposalTxs struct {
	// Single msg txs.
	ProposedOperationsTx *ProposedOperationsTx
	AddPremiumVotesTx    *AddPremiumVotesTx

	// Multi msgs txs.
	OtherTxs []*OtherMsgsTx
}

// DecodeProcessProposalTxs returns a new `processProposalTxs`.
func DecodeProcessProposalTxs(
	decoder sdk.TxDecoder,
	req *abci.RequestProcessProposal,
	pricesKeeper ve.PreBlockExecPricesKeeper,
) (*ProcessProposalTxs, error) {
	// Check len.
	numTxs := len(req.Txs)
	if err := validateNumTxs(numTxs); err != nil {
		return nil, err
	}

	// Operations.
	operationsTx, err := DecodeProposedOperationsTx(decoder, req.Txs[constants.ProposedOperationsTxIndex])
	if err != nil {
		return nil, err
	}

	// Funding samples.
	addPremiumVotesTx, err := DecodeAddPremiumVotesTx(decoder, req.Txs[numTxs+constants.AddPremiumVotesTxLenOffset])
	if err != nil {
		return nil, err
	}
	// Other txs.
	allOtherTxs := make([]*OtherMsgsTx, numTxs-constants.MinTxsCount)
	for i, txBytes := range req.Txs[constants.FirstOtherTxIndex : numTxs+constants.LastOtherTxLenOffset] {
		otherTx, err := DecodeOtherMsgsTx(decoder, txBytes)
		if err != nil {
			return nil, err
		}

		allOtherTxs[i] = otherTx
	}
	return &ProcessProposalTxs{
		ProposedOperationsTx: operationsTx,
		AddPremiumVotesTx:    addPremiumVotesTx,
		OtherTxs:             allOtherTxs,
	}, nil
}

// Validate performs `ValidateBasic` on the underlying msgs that are part of the txs.
// Returns nil if all are valid. Otherwise, returns error.
//
// Exception: for UpdateMarketPricesTx, we perform "in-memory stateful" validation
// to ensure that the new proposed prices are "valid" in comparison to daemon prices.
func (ppt *ProcessProposalTxs) Validate() error {
	// Validate single msg txs.
	singleTxs := []SingleMsgTx{
		ppt.ProposedOperationsTx,
		ppt.AddPremiumVotesTx,
	}
	for _, smt := range singleTxs {
		if err := smt.Validate(); err != nil {
			return err
		}
	}

	// Validate multi msgs txs.
	for _, mmt := range ppt.OtherTxs {
		if err := mmt.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func validateNumTxs(numTxs int) error {
	if numTxs < constants.MinTxsCount {
		return errorsmod.Wrapf(
			ErrUnexpectedNumMsgs,
			"Expected the proposal to contain at least %d txs, but got %d",
			constants.MinTxsCount,
			numTxs,
		)
	}

	return nil
}
