package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// DepositToMegavault deposits from a subaccount to megavault.
func (k msgServer) DepositToMegavault(
	goCtx context.Context,
	msg *types.MsgDepositToMegavault,
) (*types.MsgDepositToMegavaultResponse, error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)
	quoteQuantums := msg.QuoteQuantums.BigInt()

	mintedShares, err := k.Keeper.DepositToMegavault(
		ctx,
		*msg.SubaccountId,
		quoteQuantums,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgDepositToMegavaultResponse{
		MintedShares: types.BigIntToNumShares(mintedShares),
	}, nil
}
