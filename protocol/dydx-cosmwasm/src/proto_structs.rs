use cosmwasm_std::Uint64;
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct MarketPrice {
    // TODO: Missing "Id" field compared to pricestypes.MarketPrice
    // Adding `Id` field leads to this error when querying the contract:
    // `Error parsing into type dydx_cosmwasm::proto_structs::MarketPrice: missing field `id`: query wasm contract failed: unknown request`
    // Suspect it's failing due to parse the query since `Id` being tested is `0`.
    // We don't need `Id` in the response, so leaving it out for now.
    pub exponent: i32,
    pub price: i64,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct AssetPosition {
    pub asset_id: u32,
    pub quantums: Vec<u8>,
    pub index: Uint64,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct PerpetualPosition {
    pub perpetual_id: u32,
    pub quantums: Vec<u8>,
    pub funding_index: Vec<u8>,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct SubaccountId {
    pub owner: String,
    pub number: u32,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct Subaccount {
    pub id: Option<SubaccountId>,
    pub asset_positions: Vec<AssetPosition>,
    pub perpetual_positions: Vec<PerpetualPosition>,
    pub margin_enabled: bool,
}
