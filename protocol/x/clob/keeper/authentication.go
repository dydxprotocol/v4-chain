package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) IsTxAuthenticated(ctx sdk.Context, txBytes []byte) error {
	tx, err := k.txDecoder(txBytes)
	if err != nil {
		return err
	}

	if _, err := k.antehandler(ctx, tx, false); err != nil {
		return err
	}
	return nil
}
