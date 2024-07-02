use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use serde_repr::*;
use cosmwasm_std::{
  to_json_string,
  CosmosMsg,
  CustomMsg,
};

use crate::SubaccountId;


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

#[non_exhaustive]
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum DydxMsg {
  // Maps to https://github.com/dydxprotocol/v4-chain/blob/main/proto/dydxprotocol/sending/transfer.proto#L31 message on protocol
  DepositToSubaccountV1 {
    recipient: SubaccountId,
    asset_id: u32,
    quantums: u64,
  },
  // Maps to https://github.com/dydxprotocol/v4-chain/blob/main/proto/dydxprotocol/sending/transfer.proto#L50 message on protocol
    WithdrawFromSubaccountV1 {
        subaccount_number: u32,
        recipient: String,
        asset_id: u32,
        quantums: u64,
    },
  // Maps to https://github.com/dydxprotocol/v4-chain/blob/main/proto/dydxprotocol/clob/tx.proto#L78 message on protocol
  PlaceOrderV1{
    subaccount_number: u32,
    client_id: u32,
    order_flags: u32,
    clob_pair_id: u32,
    side: OrderSide,
    quantums: u64,
    subticks: u64,
    good_til_block_time: u32,
    time_in_force: OrderTimeInForce,
    reduce_only: bool,
    client_metadata: u32,
    condition_type: OrderConditionType,
    conditional_order_trigger_subticks: u64,
  },
  // Maps to https://github.com/dydxprotocol/v4-chain/blob/main/proto/dydxprotocol/clob/tx.proto#L84 on protocol
  CancelOrderV1 {
    subaccount_number: u32,
    client_id: u32,
    order_flags: u32,
    clob_pair_id: u32,
    good_til_block_time: u32,
  }
}

impl From<DydxMsg> for CosmosMsg<DydxMsg> {
  fn from(original: DydxMsg) -> Self {
    CosmosMsg::Custom(original)
  }
}

impl CustomMsg for DydxMsg {}

#[cfg(test)]
mod tests {
  use super::*;
  
  #[test]
  fn deposit_to_subaccount_msg_json_validation() {
    let msg: DydxMsg = DydxMsg::DepositToSubaccountV1 {
      recipient: SubaccountId {
        owner: "b".to_string(),
        number: 0,
      },
      asset_id: 0,
      quantums: 10000000000,
    };
    let json = to_json_string(&msg).unwrap();
    assert_eq!(
      String::from_utf8_lossy(&json),
      r#"{"deposit_to_subaccount_v1":{"recipient":{"owner":"b","number":0},"asset_id":0,"quantums":10000000000}}"#
    );
  }

  #[test]
  fn withdraw_from_subaccount_msg_json_validation() {
    let msg: DydxMsg = DydxMsg::WithdrawFromSubaccountV1 {
      subaccount_number: 0,
      recipient: "b".to_string(),
      asset_id: 0,
      quantums: 10000000000,
    };
    let json = to_json_string(&msg).unwrap();
    assert_eq!(
      String::from_utf8_lossy(&json),
      r#"{"withdraw_from_subaccount_v1":{"subaccount_number":0,"recipient":"b","asset_id":0,"quantums":10000000000}}"#
    );
  }

  #[test]
  fn place_order_msg_json_validation() {
    let msg: DydxMsg = DydxMsg::PlaceOrderV1 {
      subaccount_number: 0,
      client_id: 0,
      order_flags: 0,
      clob_pair_id: 0,
      side: OrderSide::Buy,
      quantums: 10000000000,
      subticks: 10000000000,
      good_til_block_time: 0,
      time_in_force: OrderTimeInForce::Ioc,
      reduce_only: false,
      client_metadata: 0,
      condition_type: OrderConditionType::StopLoss,
      conditional_order_trigger_subticks: 10000000000,
    };
    let json = to_json_string(&msg).unwrap();
    assert_eq!(
      String::from_utf8_lossy(&json),
      r#"{"place_order_v1":{"subaccount_number":0,"client_id":0,"order_flags":0,"clob_pair_id":0,"side":1,"quantums":10000000000,"subticks":10000000000,"good_til_block_time":0,"time_in_force":1,"reduce_only":false,"client_metadata":0,"condition_type":1,"conditional_order_trigger_subticks":10000000000}}"#
    );
  }
  
  #[test]
  fn cancel_order_msg_json_validation() {
    let msg: DydxMsg = DydxMsg::CancelOrderV1 {
      subaccount_number: 0,
      client_id: 0,
      order_flags: 0,
      clob_pair_id: 0,
      good_til_block_time: 0,
    };
    let json = to_json_string(&msg).unwrap();
    assert_eq!(
      String::from_utf8_lossy(&json),
      r#"{"cancel_order_v1":{"subaccount_number":0,"client_id":0,"order_flags":0,"clob_pair_id":0,"good_til_block_time":0}}"#
    );
  }
}