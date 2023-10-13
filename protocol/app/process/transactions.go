package process

import (
	errorsmod "cosmossdk.io/errors"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"slices"
)

const (
	minTxsCount                   = 4
	proposedOperationsTxIndex     = 0
	updateMarketPricesTxLenOffset = -1
	addPremiumVotesTxLenOffset    = -2
	acknowledgeBridgesTxLenOffset = -3
	lastOtherTxLenOffset          = acknowledgeBridgesTxLenOffset
	firstOtherTxOffset            = 1
)

func init() {
	txIndicesAndOffsets := []int{
		proposedOperationsTxIndex,
		updateMarketPricesTxLenOffset,
		addPremiumVotesTxLenOffset,
		acknowledgeBridgesTxLenOffset,
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
}

// ProcessProposalTxs is used as an intermediary struct to validate a proposed list of txs
// for `ProcessProposal`.
type ProcessProposalTxs struct {
	// Single msg txs.
	ProposedOperationsTx *ProposedOperationsTx
	AcknowledgeBridgesTx *AcknowledgeBridgesTx
	AddPremiumVotesTx    *AddPremiumVotesTx
	UpdateMarketPricesTx *UpdateMarketPricesTx

	// Multi msgs txs.
	OtherTxs []*OtherMsgsTx
}

// DecodeProcessProposalTxs returns a new `processProposalTxs`.
func DecodeProcessProposalTxs(
	ctx sdk.Context,
	decoder sdk.TxDecoder,
	req abci.RequestProcessProposal,
	bridgeKeeper ProcessBridgeKeeper,
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

	// Acknowledge bridges.
	acknowledgeBridgesTx, err := DecodeAcknowledgeBridgesTx(
		ctx,
		bridgeKeeper,
		decoder,
		req.Txs[numTxs+acknowledgeBridgesTxLenOffset],
	)
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
	for i, txBytes := range req.Txs[firstOtherTxOffset : numTxs+lastOtherTxLenOffset] {
		otherTx, err := DecodeOtherMsgsTx(decoder, txBytes)
		if err != nil {
			return nil, err
		}

		allOtherTxs[i] = otherTx
	}

	return &ProcessProposalTxs{
		ProposedOperationsTx: operationsTx,
		AcknowledgeBridgesTx: acknowledgeBridgesTx,
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
		ppt.AcknowledgeBridgesTx,
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
