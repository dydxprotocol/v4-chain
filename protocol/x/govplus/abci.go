package govplus

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/x/govplus/keeper"
)

func EndBlocker(
	ctx sdk.Context,
	keeper *keeper.Keeper,
) ([]abci.ValidatorUpdate, error) {
	fmt.Println("tian", "govplus EndBlocker")
	return keeper.BlockProposerSetUpdates(ctx)
}
