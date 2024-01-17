package wasmbinding

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/pkg/errors"

	bindings "github.com/dydxprotocol/v4-chain/protocol/wasmbinding/bindings"

	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
)

type QueryPlugin struct {
	pricesKeeper *priceskeeper.Keeper
}

// NewQueryPlugin returns a reference to a new PriceQueryPlugin.
func NewQueryPlugin(pk *priceskeeper.Keeper) *QueryPlugin {
	return &QueryPlugin{
		pricesKeeper: pk,
	}
}

func (qp QueryPlugin) HandleOracleQuery(ctx sdk.Context, queryData json.RawMessage) ([]byte, error) {
	var parsedQuery bindings.DydxOracleQuery
	if err := json.Unmarshal(queryData, &parsedQuery); err != nil {
		return nil, errorsmod.Wrap(err, "Error parsing DydxOracleQuery")
	}

	marketPrice, err := qp.pricesKeeper.GetMarketPrice(ctx, parsedQuery.MarketId)
	if err != nil {
		return nil, errorsmod.Wrap(err, fmt.Sprintf("Error getting price for market %d", parsedQuery.MarketId))
	}

	res := bindings.WasmOracleQueryResponse{
		Price:    marketPrice.Price,
		Exponent: marketPrice.Exponent,
	}
	bz, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "Error encoding WasmOracleQueryResponse as JSON")
	}

	return bz, nil
}
