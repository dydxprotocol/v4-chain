package wasmbinding

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/pkg/errors"

	clobKeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	subaccountskeeper "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	subaccountstypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type QueryPlugin struct {
	pricesKeeper      *priceskeeper.Keeper
	subaccountsKeeper *subaccountskeeper.Keeper
	clobKeeper        *clobKeeper.Keeper
}

// NewQueryPlugin returns a reference to a new PriceQueryPlugin.
func NewQueryPlugin(pk *priceskeeper.Keeper, sk *subaccountskeeper.Keeper, ck *clobKeeper.Keeper) *QueryPlugin {
	return &QueryPlugin{
		pricesKeeper:      pk,
		subaccountsKeeper: sk,
		clobKeeper:        ck,
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
	parsedSubaccountId := subaccountstypes.SubaccountId(parsedQuery)
	subaccount := qp.subaccountsKeeper.GetSubaccount(ctx, parsedSubaccountId)

	res := subaccountstypes.QuerySubaccountResponse{
		Subaccount: subaccount,
	}

	bz, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "Error encoding Subaccount as JSON")
	}

	return bz, nil
}

func (qp QueryPlugin) HandlePerpetualClobDetailsQuery(ctx sdk.Context, queryData json.RawMessage) ([]byte, error) {
	var parsedQuery clobtypes.QueryGetPerpetualClobDetailsRequest
	if err := json.Unmarshal(queryData, &parsedQuery); err != nil {
		return nil, errorsmod.Wrap(err, "Error parsing DydxGetClobQuery")
	}

	perpetualClobDetails, err := qp.clobKeeper.GetPerpetualClobDetails(ctx, clobtypes.ClobPairId(parsedQuery.Id))

	if err != nil {
		return nil, errorsmod.Wrap(err, fmt.Sprintf("Error getting clob details for pair %d", parsedQuery.Id))
	}
	bz, err := json.Marshal(perpetualClobDetails)
	if err != nil {
		return nil, errorsmod.Wrap(err, "Error encoding Clob as JSON")
	}

	return bz, nil
}
