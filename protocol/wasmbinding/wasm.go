package wasmbinding

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	clobKeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	perpKeeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	sendingkeeper "github.com/dydxprotocol/v4-chain/protocol/x/sending/keeper"
	subaccountskeeper "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
)

func RegisterCustomPlugins(
	pricesKeeper *priceskeeper.Keeper,
	sendingKeeper *sendingkeeper.Keeper,
	subaccountsKeeper *subaccountskeeper.Keeper,
	clobKeeper *clobKeeper.Keeper,
	perpKeeper *perpKeeper.Keeper,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(pricesKeeper, subaccountsKeeper, clobKeeper, perpKeeper)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageEncoders(&wasmkeeper.MessageEncoders{
		Custom: EncodeDydxCustomWasmMessage,
	})

	return []wasmkeeper.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
