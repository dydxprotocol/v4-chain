use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::{Addr, Coin};
use cw_utils::Expiration;
use dydx_cosmwasm::MarketPriceResponse;

#[cw_serde]
pub struct InstantiateMsg {
}

#[cw_serde]
pub enum ExecuteMsg {
    Approve {
        // release some coins - if quantity is None, release all coins in balance
        quantity: Option<u64>,
        //quantity: Option<Vec<Coin>>,
    },
    Refund {},
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    /// Returns a human-readable representation of the arbiter.
    #[returns(ArbiterResponse)]
    Arbiter {},
    /// Returns the current market price for the given market id.
    /// This is a custom query that is not part of the cosmwasm standard queries.
    /// It is used to demonstrate how to query the dydx module.
    #[returns(MarketPriceResponse)]
    MarketPrice { id: u32 },
}

#[cw_serde]
pub struct ArbiterResponse {
    pub arbiter: Addr,
}
