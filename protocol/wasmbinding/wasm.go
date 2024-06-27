package wasmbinding

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	clobKeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	sendingkeeper "github.com/dydxprotocol/v4-chain/protocol/x/sending/keeper"
	subaccountskeeper "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
)

func RegisterCustomPlugins(
	pricesKeeper *priceskeeper.Keeper,
	sendingKeeper *sendingkeeper.Keeper,
	subaccountsKeeper *subaccountskeeper.Keeper,
	clobKeeper *clobKeeper.Keeper,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(pricesKeeper, subaccountsKeeper, clobKeeper)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageEncoders(&wasmkeeper.MessageEncoders{
		Custom: CustomEncoder,
	})

	return []wasmkeeper.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
