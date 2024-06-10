package wasmbinding

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	clobkeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	sendingkeeper "github.com/dydxprotocol/v4-chain/protocol/x/sending/keeper"
)

func RegisterCustomPlugins(
	pricesKeeper *priceskeeper.Keeper,
	sendingKeeper *sendingkeeper.Keeper,
	clobKeeper *clobkeeper.Keeper,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(pricesKeeper)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(sendingKeeper, clobKeeper),
	)

	return []wasmkeeper.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
