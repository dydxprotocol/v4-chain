package wasmbinding

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/pkg/errors"

	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	subaccountskeeper "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	subaccountstypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type QueryPlugin struct {
	pricesKeeper      *priceskeeper.Keeper
	subaccountsKeeper *subaccountskeeper.Keeper
}

// NewQueryPlugin returns a reference to a new PriceQueryPlugin.
func NewQueryPlugin(pk *priceskeeper.Keeper) *QueryPlugin {
	return &QueryPlugin{
		pricesKeeper: pk,
	}
}

func (qp QueryPlugin) HandleMarketPriceQuery(ctx sdk.Context, queryData json.RawMessage) ([]byte, error) {
	var parsedQuery pricestypes.QueryMarketPriceRequest
	if err := json.Unmarshal(queryData, &parsedQuery); err != nil {
		return nil, errorsmod.Wrap(err, "Error parsing DydxMarketPriceQuery")
	}

	marketPrice, err := qp.pricesKeeper.GetMarketPrice(ctx, parsedQuery.Id)
	if err != nil {
		return nil, errorsmod.Wrap(err, fmt.Sprintf("Error getting price for market %d", parsedQuery.Id))
	}

	res := pricestypes.QueryMarketPriceResponse{
		MarketPrice: marketPrice,
	}

	bz, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "Error encoding MarketPrice as JSON")
	}

	return bz, nil
}

func (qp QueryPlugin) HandleSubaccountsQuery(ctx sdk.Context, queryData json.RawMessage) ([]byte, error) {
	var parsedQuery subaccountstypes.QueryGetSubaccountRequest
	if err := json.Unmarshal(queryData, &parsedQuery); err != nil {
		return nil, errorsmod.Wrap(err, "Error parsing DydxGetSubaccountQuery")
	}

	subaccount := qp.subaccountsKeeper.GetSubaccount(ctx,
		subaccountstypes.SubaccountId{
			Owner:  parsedQuery.Owner,
			Number: parsedQuery.Number,
		},
	)

	res := subaccountstypes.QuerySubaccountResponse{
		Subaccount: subaccount,
	}

	bz, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "Error encoding Subaccount as JSON")
	}

	return bz, nil
}
