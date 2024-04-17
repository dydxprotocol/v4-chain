package keeper

import (
	"context"

	pricefeedmetrics "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/metrics"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	gometrics "github.com/hashicorp/go-metrics"

	errorsmod "cosmossdk.io/errors"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (k msgServer) UpdateMarketParam(
	goCtx context.Context,
	msg *types.MsgUpdateMarketParam,
) (
	response *types.MsgUpdateMarketParamResponse,
	err error,
) {
	// Increment the appropriate success/error counter when the function finishes.
	defer func() {
		success := metrics.Success
		if err != nil {
			success = metrics.Error
		}
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, metrics.UpdateMarketParam, success},
			1,
			[]gometrics.Label{pricefeedmetrics.GetLabelForMarketId(msg.MarketParam.Id)},
		)
	}()

	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	if _, err = k.Keeper.ModifyMarketParam(ctx, msg.MarketParam); err != nil {
		return nil, err
	}

	return &types.MsgUpdateMarketParamResponse{}, nil
}
