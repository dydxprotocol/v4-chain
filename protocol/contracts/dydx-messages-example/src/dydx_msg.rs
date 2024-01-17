use core::fmt;
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use cosmwasm_std::{
  to_json_binary
};

pub trait CustomMsg: Serialize + Clone + fmt::Debug + PartialEq + JsonSchema {}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct SubaccountId {
  pub owner: String,
  pub number: u32,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct Transfer {
  pub sender: SubaccountId,
  pub recipient: SubaccountId,
  pub asset_id: u32,
  pub amount: u64,
}

#[non_exhaustive]
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum SendingMsg {
  CreateTransfer {
    transfer: Transfer,
  }
}

#[cfg(test)]
mod tests {
  use super::*;
  
  #[test]
  fn sending_msg_serializes_to_correct_json() {
    let msg: SendingMsg = SendingMsg::CreateTransfer {
      transfer: Transfer {
        sender: SubaccountId {
          owner: "a".to_string(),
          number: 0,
        },
        recipient: SubaccountId {
          owner: "b".to_string(),
          number: 0,
        },
        asset_id: 0,
        amount: 10000000000,
      },
    };
    let json = to_json_binary(&msg).unwrap();
    assert_eq!(
      String::from_utf8_lossy(&json),
      r#"{"create_transfer":{"transfer":{"sender":{"owner":"a","number":0},"recipient":{"owner":"b","number":0},"asset_id":0,"amount":10000000000}}}"#
    );
  }
}
