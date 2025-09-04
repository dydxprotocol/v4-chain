#![allow(clippy::doc_overindented_list_items)]
#![allow(clippy::doc_lazy_continuation)]
/// re-export of cosmos-sdk
pub use cosmos_sdk_proto;
use cosmos_sdk_proto::Any;

include!(concat!(env!("CARGO_MANIFEST_DIR"), "/src/_includes.rs"));

use prost::Name;

pub trait ToAny: Name + Sized {
    /// Converts the type to `prost_types::Any`.
    fn to_any(self) -> Any {
        let value = self.encode_to_vec();
        let type_url = Self::type_url();
        Any { type_url, value }
    }
}

impl<M: Name> ToAny for M {}

#[cfg(test)]
mod test {
    use super::ToAny;
    use crate::cosmos_sdk_proto::cosmos::bank::v1beta1::MsgSend;
    use crate::dydxprotocol::clob::MsgCancelOrder;

    #[test]
    /// Tests the conversion of `MsgCancelOrder` to `cosmos_sdk_proto::Any`.
    pub fn test_any_conversion() {
        let msg = MsgCancelOrder {
            order_id: None,
            good_til_oneof: None,
        };
        let any = msg.to_any();
        let url = "/dydxprotocol.clob.MsgCancelOrder";
        assert_eq!(any.type_url, url);
    }

    #[test]
    /// Tests the conversion of `MsgSend` to `cosmos_sdk_proto::Any`.
    pub fn test_any_conversion_wrapped() {
        let msg = MsgSend::default();
        let any = msg.to_any();
        let url = "/cosmos.bank.v1beta1.MsgSend";
        assert_eq!(any.type_url, url);
    }
}
