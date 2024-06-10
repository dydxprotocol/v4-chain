use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::{Addr, Coin};
use dydx_cosmwasm::{OrderConditionType, OrderSide, OrderTimeInForce, SubaccountId, Order, OrderId, Transfer};
use cw_utils::Expiration;

#[cw_serde]
pub struct InstantiateMsg {
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
    DepositToSubaccount {
        sender: String,
        recipient: SubaccountId,
        asset_id: u32,
        quantums: u64,
      },
        WithdrawFromSubaccount {
            sender: SubaccountId,
            recipient: String,
            asset_id: u32,
            quantums: u64,
        },
      PlaceOrder {
        order: Order,
      },
      CancelOrder {
        order_id: OrderId,
        good_til_block: Option<u32>,
        good_til_block_time: Option<u32>,
      },
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    /// Returns a human-readable representation of the arbiter.
    #[returns(ArbiterResponse)]
    Arbiter {},
}

#[cw_serde]
pub struct ArbiterResponse {
    pub arbiter: Addr,
}
