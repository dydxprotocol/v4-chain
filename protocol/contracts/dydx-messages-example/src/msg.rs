use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::{Addr, Coin};
use cw_utils::Expiration;
use dydx_cosmwasm::{OrderConditionType, OrderSide, OrderTimeInForce, SubaccountId, MarketPrice, Subaccount};

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
    PlaceOrder {
        subaccount_id: SubaccountId,
        client_id: u32,
        order_flags: u32,
        clob_pair_id: u32,
        side: OrderSide,
        quantums: u64,
        subticks: u64,
        good_til_block: Option<u32>,
        good_til_block_time: Option<u32>,
        time_in_force: OrderTimeInForce,
        reduce_only: bool,
        client_metadata: u32,
        condition_type: OrderConditionType,
        conditional_order_trigger_subticks: u64,
    }
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
    },
    #[returns(SubaccountResponse)]
    Subaccount {
        address: String,
        subaccountNumber: u32,
    }
}

#[cw_serde]
pub struct ArbiterResponse {
    pub arbiter: Addr,
}
