package keeper

import (
	"context"
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// UpdateParams updates the parameters of the vault module.
// Deprecated since v6.x in favor of UpdateDefaultQuotingParams.
func (k msgServer) UpdateParams(
	goCtx context.Context,
	msg *types.MsgUpdateParams,
) (*types.MsgUpdateParamsResponse, error) {
	return nil, fmt.Errorf("deprecated since v6.x in favor of UpdateDefaultQuotingParams")
}
