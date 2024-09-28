package keeper

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
		k.SetAssetYieldIndex(ctx, new(big.Rat).SetInt64(1))
	}

	k.SetSDAIPrice(ctx, conversionRate)
	return nil
}
