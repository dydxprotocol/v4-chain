use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::{Addr, Coin};
use dydx_cosmwasm::{OrderConditionType, OrderSide, OrderTimeInForce, SubaccountId, Order, OrderId, Transfer};
use cw_utils::Expiration;
use dydx_cosmwasm::MarketPriceResponse;

#[cw_serde]
pub struct InstantiateMsg {
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
