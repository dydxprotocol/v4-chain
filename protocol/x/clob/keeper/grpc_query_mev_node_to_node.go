package keeper

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) MevNodeToNodeCalculation(
	c context.Context,
	req *types.MevNodeToNodeCalculationRequest,
) (
	*types.MevNodeToNodeCalculationResponse,
	error,
) {
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	// Validate that the request is valid.
	if err := validateMevNodeToNodeRequest(req); err != nil {
		log.ErrorLogWithError(ctx, "Failed to validate MEV node to node calculation request", err,
			"mev_calculation_request", req,
		)
		return nil, err
	}

	blockProposerPnL, validatorPnL := k.InitializeCumulativePnLsFromRequest(ctx, req)

	k.CalculateSubaccountPnLForMevMatches(
		ctx,
		blockProposerPnL,
		req.BlockProposerMatches,
	)
	k.CalculateSubaccountPnLForMevMatches(
		ctx,
		validatorPnL,
		req.ValidatorMevMetrics.ValidatorMevMatches,
	)

	mevAndVolumePerClob := make(
		[]types.MevNodeToNodeCalculationResponse_MevAndVolumePerClob,
		0,
		len(blockProposerPnL),
	)
	for clobPairId, blockProposerSubaccountPnL := range blockProposerPnL {
		// Calculate MEV for the given market.
		mev, _ := blockProposerSubaccountPnL.CalculateMev(validatorPnL[clobPairId]).Float32()
		validatorVolumeQuoteQuantums := new(big.Int).Div(
			validatorPnL[clobPairId].VolumeQuoteQuantums,
			big.NewInt(2),
		)

		mevAndVolumePerClob = append(
			mevAndVolumePerClob,
			types.MevNodeToNodeCalculationResponse_MevAndVolumePerClob{
				ClobPairId: clobPairId.ToUint32(),
				Mev:        mev,
				Volume:     validatorVolumeQuoteQuantums.Uint64(),
			},
		)
	}

	return &types.MevNodeToNodeCalculationResponse{
		Results: mevAndVolumePerClob,
	}, nil
}

func (k Keeper) InitializeCumulativePnLsFromRequest(
	ctx sdk.Context,
	req *types.MevNodeToNodeCalculationRequest,
) (
	blockProposerPnL map[types.ClobPairId]*CumulativePnL,
	validatorPnL map[types.ClobPairId]*CumulativePnL,
) {
	clobMetadata := make(map[types.ClobPairId]ClobMetadata, len(req.ValidatorMevMetrics.ClobMidPrices))
	for _, clobMidPrice := range req.ValidatorMevMetrics.ClobMidPrices {
		clobPairId := types.ClobPairId(clobMidPrice.ClobPair.Id)
		clobMetadata[clobPairId] = ClobMetadata{
			ClobPair: clobMidPrice.ClobPair,
			MidPrice: types.Subticks(clobMidPrice.Subticks),
		}
	}

	blockProposerPnL, validatorPnL = k.InitializeCumulativePnLs(
		ctx,
		k.perpetualsKeeper,
		clobMetadata,
	)

	return blockProposerPnL, validatorPnL
}

// ValidateMevNodeToNodeRequest validates a MEV node to node calculation request. It returns
// an error if the request is invalid.
func validateMevNodeToNodeRequest(
	req *types.MevNodeToNodeCalculationRequest,
) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ValidatorMevMetrics == nil {
		return status.Error(codes.InvalidArgument, "missing validator MEV metrics")
	}

	seenClobPairIds := make(map[types.ClobPairId]bool)
	for _, clobMidPrice := range req.ValidatorMevMetrics.ClobMidPrices {
		clobPairId := types.ClobPairId(clobMidPrice.ClobPair.Id)
		if _, exists := seenClobPairIds[clobPairId]; exists {
			return status.Error(codes.InvalidArgument, "duplicate CLOB pair")
		}
		seenClobPairIds[clobPairId] = true
	}

	// Validate that all validator matches reference a valid CLOB pair.
	for _, validatorMevMatch := range req.ValidatorMevMetrics.ValidatorMevMatches.Matches {
		if _, exists := seenClobPairIds[types.ClobPairId(validatorMevMatch.ClobPairId)]; !exists {
			return status.Error(codes.InvalidArgument, "validator MEV match references invalid CLOB pair")
		}
	}

	// Validate that all validator liquidation matches reference a valid CLOB pair.
	for _, validatorMevMatch := range req.ValidatorMevMetrics.ValidatorMevMatches.Matches {
		if _, exists := seenClobPairIds[types.ClobPairId(validatorMevMatch.ClobPairId)]; !exists {
			return status.Error(codes.InvalidArgument, "validator MEV match references invalid CLOB pair")
		}
	}

	// Validate that all validator matches reference a valid CLOB pair.
	if err := validateMevMatchesHaveValidClob(
		req.ValidatorMevMetrics.ValidatorMevMatches,
		seenClobPairIds,
	); err != nil {
		return err
	}

	// Validate that all block proposer matches reference a valid CLOB pair.
	if err := validateMevMatchesHaveValidClob(
		req.BlockProposerMatches,
		seenClobPairIds,
	); err != nil {
		return err
	}

	return nil
}

// validateMevMatchesHaveValidClob validates that all matches reference a valid CLOB pair. It returns
// an error if any MEV matches do not reference a valid CLOB pair.
func validateMevMatchesHaveValidClob(
	mevMatches *types.ValidatorMevMatches,
	seenClobPairIds map[types.ClobPairId]bool,
) error {
	// Validate that all validator matches reference a valid CLOB pair.
	for _, validatorMevMatch := range mevMatches.Matches {
		if _, exists := seenClobPairIds[types.ClobPairId(validatorMevMatch.ClobPairId)]; !exists {
			return status.Error(codes.InvalidArgument, "validator MEV match references invalid CLOB pair")
		}
	}

	// Validate that all validator liquidation matches reference a valid CLOB pair.
	for _, validatorMevMatch := range mevMatches.Matches {
		if _, exists := seenClobPairIds[types.ClobPairId(validatorMevMatch.ClobPairId)]; !exists {
			return status.Error(codes.InvalidArgument, "validator MEV match references invalid CLOB pair")
		}
	}

	return nil
}
