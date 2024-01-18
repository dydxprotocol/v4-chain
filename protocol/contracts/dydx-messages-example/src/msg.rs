use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::{Addr, Coin};
use cw_utils::Expiration;
use dydx_cosmwasm::MarketPrice;

#[cw_serde]
pub struct InstantiateMsg {
    pub arbiter: String,
    pub recipient: String,
    /// When end height set and block height exceeds this value, the escrow is expired.
    /// Once an escrow is expired, it can be returned to the original funder (via "refund").
    ///
    /// When end time (in seconds since epoch 00:00:00 UTC on 1 January 1970) is set and
    /// block time exceeds this value, the escrow is expired.
    /// Once an escrow is expired, it can be returned to the original funder (via "refund").
    pub expiration: Option<Expiration>,
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
    #[returns(MarketPrice)]
    MarketPrice {
        id: u32,
    }
}

#[cw_serde]
pub struct ArbiterResponse {
    pub arbiter: Addr,
}
