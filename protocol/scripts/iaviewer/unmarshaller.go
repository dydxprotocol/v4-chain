package main

import (
	"github.com/cosmos/gogoproto/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	app "github.com/dydxprotocol/v4-chain/protocol/app"
	clob "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func protoUnmarshaller[M proto.Message](b []byte) string {
	cdc := app.GetEncodingConfig().Codec
	// var m M
	// // TODO: avoid reflection?
	// msgType := reflect.TypeOf(m).Elem()
	// m = reflect.New(msgType).Interface().(M)
	// cdc.MustUnmarshal(b, m)
	// return m.String()

	var acc sdk.AccountI
	if err := cdc.UnmarshalInterface(b, &acc); err != nil {
		panic(err)
	}
	return string(acc.String())
}

// Maps prefix names for modules to an inner registry map of type map[string]func([]byte) string.
// For iavl key-value pair (K_i, V_i) and registry map key-value pair (K_r, V_r), V_r will be used
// to unmarshal V_i if K_r is a prefix of K_i.
// Thus, keys each inner registry map should not be prefixes of each other.
var unmarshallerRegistry = map[string]map[string]func([]byte) string{
	"s/k:clob/": {
		"Clob:":      protoUnmarshaller[*clob.ClobPair],
		"EqTierCfg":  protoUnmarshaller[*clob.EquityTierLimitConfiguration],
		"LiqCfg":     protoUnmarshaller[*clob.LiquidationsConfig],
		"RateLimCfg": protoUnmarshaller[*clob.BlockRateLimitConfiguration],
	},
	"s/k:bank/": {
		"\x01": protoUnmarshaller[*banktypes.Metadata],
	},
	"s/k:acc/": {
		"\x01": protoUnmarshaller[*authtypes.BaseAccount],
	},
}
