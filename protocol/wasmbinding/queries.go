package wasmbinding

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/pkg/errors"

	bindings "github.com/dydxprotocol/v4-chain/protocol/wasmbinding/bindings"
	clobKeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
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
	var parsedQuery bindings.MarketPriceRequestWrapper

	if err := json.Unmarshal(queryData, &parsedQuery); err != nil {
		return nil, errorsmod.Wrap(err, "Error parsing DydxMarketPriceQuery")
	}

	marketPrice, err := qp.pricesKeeper.GetMarketPrice(ctx, parsedQuery.MarketPrice.Id)
	if err != nil {
		return nil, errorsmod.Wrap(err, fmt.Sprintf("Error getting price for market %d", parsedQuery.MarketPrice.Id))
	}

	bz, err := json.Marshal(marketPrice)
	if err != nil {
		return nil, errorsmod.Wrap(err, "Error encoding MarketPrice as JSON")
	}

	return bz, nil
}

func (qp QueryPlugin) HandleSubaccountsQuery(ctx sdk.Context, queryData json.RawMessage) ([]byte, error) {
	var parsedQuery bindings.SubaccountRequestWrapper
	if err := json.Unmarshal(queryData, &parsedQuery); err != nil {
		return nil, errorsmod.Wrap(err, "Error parsing DydxGetClobQuery")
	}

	parsedSubaccountId := subaccountstypes.SubaccountId{
		Owner:  parsedQuery.Subaccount.Owner,
		Number: parsedQuery.Subaccount.Number,
	}
	subaccount := qp.subaccountsKeeper.GetSubaccount(ctx, parsedSubaccountId)

	bz, err := json.Marshal(subaccount)
	if err != nil {
		return nil, errorsmod.Wrap(err, "Error encoding Subaccount as JSON")
	}

	fmt.Println("bz values", string(bz))
	if err != nil {
		return nil, errorsmod.Wrap(err, "Error encoding Subaccount as JSON")
	}

	return bz, nil
}

func (qp QueryPlugin) HandlePerpetualClobDetailsQuery(ctx sdk.Context, queryData json.RawMessage) ([]byte, error) {
	var parsedQuery bindings.PerpeutalClobDetailsRequestWrapper
	if err := json.Unmarshal(queryData, &parsedQuery); err != nil {
		return nil, errorsmod.Wrap(err, "Error parsing DydxGetClobQuery")
	}

	perpetualClobDetails, err := qp.clobKeeper.GetPerpetualClobDetails(ctx, clobtypes.ClobPairId(parsedQuery.PerpetualClobDetails.Id))

	if err != nil {
		return nil, errorsmod.Wrap(err, fmt.Sprintf("Error getting clob details for pair %d", parsedQuery.PerpetualClobDetails.Id))
	}
	bz, err := json.Marshal(perpetualClobDetails)
	if err != nil {
		return nil, errorsmod.Wrap(err, "Error encoding Clob as JSON")
	}

	return bz, nil
}
