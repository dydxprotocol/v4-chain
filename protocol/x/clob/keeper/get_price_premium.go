package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

// GetPricePremiumForPerpetual returns the price premium for a perpetual market,
// according to the memclob state.
func (k Keeper) GetPricePremiumForPerpetual(
	ctx sdk.Context,
	perpetualId uint32,
	params perptypes.GetPricePremiumParams,
) (
	premiumPpm int32,
	err error,
) {
	clobPairId, err := k.GetClobPairIdForPerpetual(ctx, perpetualId)
	if err != nil {
		return 0, err
	}

	clobPair, found := k.GetClobPair(ctx, clobPairId)
	if !found {
		return 0, sdkerrors.Wrapf(
			types.ErrInvalidClob,
			"GetPricePremiumForPerpetual: did not find clob pair with clobPairId = %d",
			clobPairId,
		)
	}

	return k.MemClob.GetPricePremium(
		ctx,
		clobPair,
		params,
	)
}
