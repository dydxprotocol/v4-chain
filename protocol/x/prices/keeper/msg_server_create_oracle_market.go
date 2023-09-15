package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func (k msgServer) CreateOracleMarket(
	goCtx context.Context,
	msg *types.MsgCreateOracleMarket,
) (*types.MsgCreateOracleMarketResponse, error) {
	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Use zero oracle price to create the new market.
	// Note that valid oracle price updates cannot be zero (checked in MsgUpdateMarketPrices.ValidateBasic),
	// so a zero oracle price indicates that the oracle price has never been updated.
	zeroMarketPrice := types.MarketPrice{
		Id:       msg.Params.Id,
		Exponent: msg.Params.Exponent,
		Price:    0,
	}
	if _, err := k.Keeper.CreateMarket(ctx, msg.Params, zeroMarketPrice); err != nil {
		return nil, err
	}
	return &types.MsgCreateOracleMarketResponse{}, nil
}
