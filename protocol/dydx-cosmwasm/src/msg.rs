use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use serde_repr::*;
use cosmwasm_std::{
  CosmosMsg,
  CustomMsg,
};

use crate::SubaccountId;

// TODO(affan): handle issue with `GoodTilOneof` in `PlaceOrder` and `CancelOrder` not serializing correctly

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct Transfer {
  pub sender: SubaccountId,
  pub recipient: SubaccountId,
  pub asset_id: u32,
  pub amount: u64,
}

#[derive(Serialize_repr, Deserialize_repr, Clone, Debug, PartialEq, Eq, JsonSchema)]
#[repr(u32)]
pub enum OrderSide {
  Unspecified = 0,
  Buy = 1,
  Sell = 2,
}

#[derive(Serialize_repr, Deserialize_repr, Clone, Debug, PartialEq, Eq, JsonSchema)]
#[repr(u32)]
pub enum OrderTimeInForce {
  Unspecified = 0,
  Ioc = 1,
  PostOnly = 2,
  FillOrKill = 3,
}

#[derive(Serialize_repr, Deserialize_repr, Clone, Debug, PartialEq, Eq, JsonSchema)]
#[repr(u32)]
pub enum OrderConditionType {
  Unspecified = 0,
  StopLoss = 1,
  TakeProfit = 2,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct OrderId {
  pub subaccount_id: SubaccountId,
  pub client_id: u32,
  pub order_flags: u32,
  pub clob_pair_id: u32,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct Order {
  pub order_id: OrderId,
  pub side: OrderSide,
  pub quantums: u64,
  pub subticks: u64,
  pub good_til_oneof: GoodTilOneof,
  pub time_in_force: OrderTimeInForce,
  pub reduce_only: bool,
  pub client_metadata: u32,
  pub condition_type: OrderConditionType,
  pub conditional_order_trigger_subticks: u64,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum GoodTilOneof {
    GoodTilBlock(u32),
    GoodTilBlockTime(u32),
}


#[non_exhaustive]
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum DydxMsg {
  CreateTransfer {
    transfer: Transfer,
  },
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
    good_til_oneof: GoodTilOneof,
  }
}

impl From<DydxMsg> for CosmosMsg<DydxMsg> {
  fn from(original: DydxMsg) -> Self {
    CosmosMsg::Custom(original)
  }
}

impl CustomMsg for DydxMsg {}
