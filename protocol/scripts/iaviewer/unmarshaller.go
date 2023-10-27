package main

import (
	"reflect"

	"github.com/cosmos/cosmos-sdk/codec"
	app "github.com/dydxprotocol/v4-chain/protocol/app"
	clob "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func protoUnmarshaller[M codec.ProtoMarshaler](b []byte) string {
	cdc := app.GetEncodingConfig().Codec
	var m M
	// TODO: avoid reflection?
	msgType := reflect.TypeOf(m).Elem()
	m = reflect.New(msgType).Interface().(M)
	cdc.MustUnmarshal(b, m)
	return m.String()
}

var unmarshallerRegistry = map[string]map[string]func([]byte) string{
	"s/k:clob/": {
		"Clob:":      protoUnmarshaller[*clob.ClobPair],
		"EqTierCfg":  protoUnmarshaller[*clob.EquityTierLimitConfiguration],
		"LiqCfg":     protoUnmarshaller[*clob.LiquidationsConfig],
		"RateLimCfg": protoUnmarshaller[*clob.BlockRateLimitConfiguration],
	},
}
