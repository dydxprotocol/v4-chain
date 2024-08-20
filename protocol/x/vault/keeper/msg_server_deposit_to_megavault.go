package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// DepositToMegavault deposits from a subaccount to megavault.
func (k msgServer) DepositToMegavault(
	goCtx context.Context,
	msg *types.MsgDepositToMegavault,
) (*types.MsgDepositToMegavaultResponse, error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)
	quoteQuantums := msg.QuoteQuantums.BigInt()

	// Mint shares.
	err := k.MintShares(
		ctx,
		msg.SubaccountId.Owner,
		quoteQuantums,
	)
	if err != nil {
		return nil, err
	}

	// Transfer from sender subaccount to megavault.
	// Note: Transfer should take place after minting shares for
	// shares calculation to be correct.
	err = k.sendingKeeper.ProcessTransfer(
		ctx,
		&sendingtypes.Transfer{
			Sender:    *msg.SubaccountId,
			Recipient: types.MegavaultMainSubaccount,
			AssetId:   assettypes.AssetUsdc.Id,
			Amount:    quoteQuantums.Uint64(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgDepositToMegavaultResponse{}, nil
}
