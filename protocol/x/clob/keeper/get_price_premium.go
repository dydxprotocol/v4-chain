package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

// GetPricePremiumForPerpetual returns the price premium for a perpetual market,
// according to the memclob state. If the market is not active, returns zero premium ppm.
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
		return 0, errorsmod.Wrapf(
			types.ErrInvalidClob,
			"GetPricePremiumForPerpetual: did not find clob pair with clobPairId = %d",
			clobPairId,
		)
	}

	// Zero premium if ClobPair is not active.
	if clobPair.Status != types.ClobPair_STATUS_ACTIVE {
		return 0, nil
	}

	return k.MemClob.GetPricePremium(
		ctx,
		clobPair,
		params,
	)
}
