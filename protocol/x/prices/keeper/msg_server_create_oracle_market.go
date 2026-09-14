package keeper

import (
	"context"
	"strings"

	"github.com/cosmos/cosmos-sdk/telemetry"
	pricefeedmetrics "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	gometrics "github.com/hashicorp/go-metrics"

	errorsmod "cosmossdk.io/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func (k msgServer) CreateOracleMarket(
	goCtx context.Context,
	msg *types.MsgCreateOracleMarket,
) (
	response *types.MsgCreateOracleMarketResponse,
	err error,
) {
	// Increment the appropriate success/error counter when the function finishes.
	defer func() {
		success := metrics.Success
		if err != nil {
			success = metrics.Error
		}
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, metrics.CreateOracleMarket, success},
			1,
			[]gometrics.Label{pricefeedmetrics.GetLabelForMarketId(msg.Params.Id)},
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

	// Normalize the pair to match the permissionless listing path (x/listing
	// uppercases tickers before creation). Without this, a governance proposal
	// that passes a non-canonical case (e.g. "btc-usd") stores the pair
	// verbatim, while the pricefeed daemon keys its exchange-config lookup on
	// the exact stored string, so the market is created but never priced.
	msg.Params.Pair = strings.ToUpper(msg.Params.Pair)

	exponent, err := k.Keeper.GetExponent(ctx, msg.Params.Pair)
	if err != nil {
		return nil, err
	}

	// Use zero oracle price to create the new market.
	// Note that valid oracle price updates cannot be zero (checked in MsgUpdateMarketPrices.ValidateBasic),
	// so a zero oracle price indicates that the oracle price has never been updated.
	zeroMarketPrice := types.MarketPrice{
		Id:       msg.Params.Id,
		Exponent: exponent,
		Price:    0,
	}
	if _, err = k.Keeper.CreateMarket(ctx, msg.Params, zeroMarketPrice); err != nil {
		return nil, err
	}

	return &types.MsgCreateOracleMarketResponse{}, nil
}
