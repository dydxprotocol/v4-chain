package process

import (
	errorsmod "cosmossdk.io/errors"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type txtype int

const (
	ProposedOperationsTxType txtype = 1
	AcknowledgeBridgesTxType txtype = 2
	AddPremiumVotesTxType    txtype = 3
	UpdateMarketPricesTxType txtype = 4
)

const (
	MinTxsCount = 4
)

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
	if numTxs < MinTxsCount {
		return nil, errorsmod.Wrapf(
			ErrUnexpectedNumMsgs,
			"Expected the proposal to contain at least %d txs, but got %d",
			MinTxsCount,
			numTxs,
		)
	}

	txTypeToIdx, idxToTxType := GetAppInjectedMsgIdxMaps(numTxs)

	// Operations.
	orderIdx, ok := txTypeToIdx[ProposedOperationsTxType]
	if !ok {
		panic("must define ProposedOperationsTxType")
	}
	operationsTx, err := DecodeProposedOperationsTx(decoder, req.Txs[orderIdx])
	if err != nil {
		return nil, err
	}

	// Acknowledge bridges.
	acknowledgeBridgesIdx, ok := txTypeToIdx[AcknowledgeBridgesTxType]
	if !ok {
		panic("must define AcknowledgeBridgesTxType")
	}
	acknowledgeBridgesTx, err := DecodeAcknowledgeBridgesTx(
		ctx,
		bridgeKeeper,
		decoder,
		req.Txs[acknowledgeBridgesIdx],
	)
	if err != nil {
		return nil, err
	}

	// Funding samples.
	addFundingIdx, ok := txTypeToIdx[AddPremiumVotesTxType]
	if !ok {
		panic("must define AddPremiumVotesTxType")
	}
	addPremiumVotesTx, err := DecodeAddPremiumVotesTx(decoder, req.Txs[addFundingIdx])
	if err != nil {
		return nil, err
	}

	// Price updates.
	updatePricesIdx, ok := txTypeToIdx[UpdateMarketPricesTxType]
	if !ok {
		panic("must define AddPremiumVotesTxType")
	}
	updatePricesTx, err := DecodeUpdateMarketPricesTx(ctx, pricesKeeper, decoder, req.Txs[updatePricesIdx])
	if err != nil {
		return nil, err
	}

	// Other txs.
	allOtherTxs := make([]*OtherMsgsTx, numTxs-len(txTypeToIdx))
	idx := 0
	for i, txBytes := range req.Txs {
		if _, exists := idxToTxType[i]; exists { // skip, because tx is not part of "others".
			continue
		}

		otherTx, err := DecodeOtherMsgsTx(decoder, txBytes)
		if err != nil {
			return nil, err
		}

		allOtherTxs[idx] = otherTx
		idx += 1
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
