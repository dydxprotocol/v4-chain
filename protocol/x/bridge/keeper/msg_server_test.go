package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

func setupMsgServer(t *testing.T) (keeper.Keeper, types.MsgServer, context.Context) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	return k, keeper.NewMsgServerImpl(k), sdk.WrapSDKContext(ctx)
}
