package keeper

import (
	"context"
	"math/big"

	errorsmod "cosmossdk.io/errors"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/api"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpdateSDAIConversionRate updates the sDAI conversion rate.
func (k msgServer) UpdateSDAIConversionRate(
	goCtx context.Context,
	msg *types.MsgUpdateSDAIConversionRate,
) (*types.MsgUpdateSDAIConversionRateResponse, error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	bigConversionRate, err := ConvertStringToBigInt(msg.ConversionRate)
	if err != nil {
		return nil, err
	}

	lastTenEvents := k.sDAIEventManager.GetLastTensDAIEventsUnordered()

	conversionRate, err := findMatchingConversionRate(lastTenEvents, bigConversionRate)
	if err != nil {
		return nil, err
	}

	err = k.setNewSDaiConversionRate(ctx, conversionRate)
	if err != nil {
		return nil, err
	}

	err = k.ProcessNewTDaiConversionRateUpdate(ctx)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateSDAIConversionRateResponse{}, nil
}

func (k msgServer) setNewSDaiConversionRate(
	ctx sdk.Context,
	conversionRate *big.Int,
) error {

	currentConversionRate, converstionRateInitialized := k.GetSDAIPrice(ctx)

	if converstionRateInitialized && conversionRate.Cmp(currentConversionRate) <= 0 {
		return errorsmod.Wrap(
			types.ErrInvalidSDAIConversionRate,
			"The suggested sDAI price must be greater than the curret one",
		)
	}

	if !converstionRateInitialized {
		k.SetAssetYieldIndex(ctx, new(big.Rat).SetInt64(0))
	}

	k.SetSDAIPrice(ctx, conversionRate)
	return nil
}

func findMatchingConversionRate(
	lastTenEvents [10]api.AddsDAIEventsRequest,
	conversionRateFromMsg *big.Int,
) (*big.Int, error) {

	for _, event := range lastTenEvents {

		if event.ConversionRate == "" {
			continue
		}

		conversionRate, err := ConvertStringToBigInt(event.ConversionRate)
		if err != nil {
			return nil, err
		}

		if conversionRateFromMsg.Cmp(conversionRate) == 0 {
			return conversionRate, nil
		}

	}

	return big.NewInt(0), errorsmod.Wrap(
		types.ErrInvalidSDAIConversionRate,
		"The suggested sDAI conversion rate is not valid",
	)
}
