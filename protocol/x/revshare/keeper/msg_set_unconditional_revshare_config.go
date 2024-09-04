package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

// TODO(OTE-775): Add methods to set and get rev share config
func (k msgServer) SetUnconditionalRevShareConfig(
	goCtx context.Context,
	msg *types.MsgUpdateUnconditionalRevShareConfig,
) (*types.MsgUpdateUnconditionalRevShareConfigResponse, error) {
	return nil, nil
}
