use cosmwasm_std::CustomQuery;
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use crate::proto_structs::MarketPrice;
use crate::proto_structs::Subaccount;
use crate::route::DydxRoute;

/// SeiQueryWrapper is an override of QueryRequest::Custom to access Sei-specific modules
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub struct DydxQueryWrapper {
    pub route: DydxRoute,
    pub query_data: DydxQuery,
}

// implement custom query
impl CustomQuery for DydxQueryWrapper {}

/// SeiQuery is defines available query datas
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum DydxQuery {
    MarketPrice {
        id: u32,
    },
    Subaccount {
        owner: String,
        number: u32,
    }
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct SubaccountResponse {
    pub subaccount: Subaccount,
}
