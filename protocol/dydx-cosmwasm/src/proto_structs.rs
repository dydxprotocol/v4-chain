use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cosmwasm_std::Uint64;

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct MarketPrice {
    pub id: u32,
    pub exponent: i32,
    pub price: Uint64,
}
