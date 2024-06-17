use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::{Addr, Coin};
use dydx_cosmwasm::{OrderConditionType, OrderSide, OrderTimeInForce, SubaccountId, Order, OrderId, Transfer};
use cw_utils::Expiration;
use dydx_cosmwasm::MarketPriceResponse;
use dydx_cosmwasm::DydxMsg;

#[cw_serde]
pub struct InstantiateMsg {
}

#[cw_serde]
pub enum ExecuteMsg {
    DydxMsg(DydxMsg),
}

#[cw_serde]
pub struct ArbiterResponse {
    pub arbiter: Addr,
}
