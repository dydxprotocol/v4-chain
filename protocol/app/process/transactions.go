package process

import (
	"slices"

	errorsmod "cosmossdk.io/errors"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	minTxsCount                   = 3
	proposedOperationsTxIndex     = 0
	updateMarketPricesTxLenOffset = -1
	addPremiumVotesTxLenOffset    = -2
	lastOtherTxLenOffset          = addPremiumVotesTxLenOffset
	firstOtherTxIndex             = proposedOperationsTxIndex + 1
)

func init() {
	txIndicesAndOffsets := []int{
		proposedOperationsTxIndex,
		addPremiumVotesTxLenOffset,
		updateMarketPricesTxLenOffset,
	}
	if minTxsCount != len(txIndicesAndOffsets) {
		panic("minTxsCount does not match expected count of Txs.")
	}
	if lib.ContainsDuplicates(txIndicesAndOffsets) {
		panic("Duplicate indices/offsets defined for Txs.")
	}
	if slices.Min[[]int](txIndicesAndOffsets) != lastOtherTxLenOffset {
		panic("lastTxLenOffset is not the lowest offset")
	}
	if slices.Max[[]int](txIndicesAndOffsets)+1 != firstOtherTxIndex {
		panic("firstOtherTxIndex is <= the maximum offset")
	}
	txIndicesForMinTxsCount := []int{
		proposedOperationsTxIndex,
		addPremiumVotesTxLenOffset + minTxsCount,
		updateMarketPricesTxLenOffset + minTxsCount,
	}
	if minTxsCount != len(txIndicesForMinTxsCount) {
		panic("minTxsCount does not match expected count of Txs.")
	}
	if lib.ContainsDuplicates(txIndicesForMinTxsCount) {
		panic("Overlapping indices and offsets defined for Txs.")
	}
	if minTxsCount != firstOtherTxIndex-lastOtherTxLenOffset {
		panic("Unexpected gap between firstOtherTxIndex and lastOtherTxLenOffset which is greater than minTxsCount")
	}
}

// ProcessProposalTxs is used as an intermediary struct to validate a proposed list of txs
// for `ProcessProposal`.
type ProcessProposalTxs struct {
	// Single msg txs.
	ProposedOperationsTx *ProposedOperationsTx
	AddPremiumVotesTx    *AddPremiumVotesTx
	UpdateMarketPricesTx *UpdateMarketPricesTx

	// Multi msgs txs.
	OtherTxs []*OtherMsgsTx
}

// DecodeProcessProposalTxs returns a new `processProposalTxs`.
func DecodeProcessProposalTxs(
	ctx sdk.Context,
	decoder sdk.TxDecoder,
	req *abci.RequestProcessProposal,
	pricesKeeper ProcessPricesKeeper,
) (*ProcessProposalTxs, error) {
	// Check len.
	numTxs := len(req.Txs)
	if numTxs < minTxsCount {
		return nil, errorsmod.Wrapf(
			ErrUnexpectedNumMsgs,
			"Expected the proposal to contain at least %d txs, but got %d",
			minTxsCount,
			numTxs,
		)
	}

	// Operations.
	operationsTx, err := DecodeProposedOperationsTx(decoder, req.Txs[proposedOperationsTxIndex])
	if err != nil {
		return nil, err
	}

	// Funding samples.
	addPremiumVotesTx, err := DecodeAddPremiumVotesTx(decoder, req.Txs[numTxs+addPremiumVotesTxLenOffset])
	if err != nil {
		return nil, err
	}

	// Price updates.
	updatePricesTx, err := DecodeUpdateMarketPricesTx(
		ctx,
		pricesKeeper,
		decoder,
		req.Txs[numTxs+updateMarketPricesTxLenOffset],
	)
	if err != nil {
		return nil, err
	}

	// Other txs.
	allOtherTxs := make([]*OtherMsgsTx, numTxs-minTxsCount)
	for i, txBytes := range req.Txs[firstOtherTxIndex : numTxs+lastOtherTxLenOffset] {
		otherTx, err := DecodeOtherMsgsTx(decoder, txBytes)
		if err != nil {
			return nil, err
		}

		allOtherTxs[i] = otherTx
	}

	return &ProcessProposalTxs{
		ProposedOperationsTx: operationsTx,
		AddPremiumVotesTx:    addPremiumVotesTx,
		UpdateMarketPricesTx: updatePricesTx,
		OtherTxs:             allOtherTxs,
	}, nil
}

// Validate performs `ValidateBasic` on the underlying msgs that are part of the txs.
// Returns nil if all are valid. Otherwise, returns error.
//
// Exception: for UpdateMarketPricesTx, we perform "in-memory stateful" validation
// to ensure that the new proposed prices are "valid" in comparison to index prices.
func (ppt *ProcessProposalTxs) Validate() error {
	// Validate single msg txs.
	singleTxs := []SingleMsgTx{
		ppt.ProposedOperationsTx,
		ppt.AddPremiumVotesTx,
		ppt.UpdateMarketPricesTx,
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
