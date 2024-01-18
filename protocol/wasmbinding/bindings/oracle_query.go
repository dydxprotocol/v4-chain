package bindings

// DydxOracleQuery contains oracle price from x/prices module.
// See https://github.com/osmosis-labs/osmosis-bindings/blob/main/packages/bindings/src/query.rs
type DydxOracleQuery struct {
	MarketId uint32 `json:"market_id,omitempty"`
}

type WasmOracleQueryResponse struct {
	Price    uint64 `json:"price,omitempty"`
	Exponent int32  `json:"exponent,omitempty"`
}
