/// re-export of cosmos-sdk
pub use cosmos_sdk_proto;

include!(concat!(env!("CARGO_MANIFEST_DIR"), "/src/_includes.rs"));

use prost::Name;

pub trait ToAny: Name + Sized {
    /// Converts the type to `prost_types::Any`.
    fn to_any(self) -> prost_types::Any {
        let value = self.encode_to_vec();
        let type_url = Self::type_url();
        prost_types::Any { type_url, value }
    }
}

impl<M: Name> ToAny for M {}

#[cfg(test)]
mod test {
    use super::ToAny;
    use crate::cosmos_sdk_proto::cosmos::bank::v1beta1::MsgSend;
    use crate::dydxprotocol::clob::MsgCancelOrder;

    #[test]
    pub fn test_any_conversion() {
        /// Tests the conversion of `MsgCancelOrder` to `prost_types::Any`.
        let msg = MsgCancelOrder {
            order_id: None,
            good_til_oneof: None,
        };
        let any = msg.to_any();
        let url = "/dydxprotocol.clob.MsgCancelOrder";
        assert_eq!(any.type_url, url);
    }

    #[test]
    pub fn test_any_conversion_wrapped() {
        /// Tests the conversion of `MsgSend` to `prost_types::Any`.
        let msg = MsgSend::default();
        let any = msg.to_any();
        let url = "/cosmos.bank.v1beta1.MsgSend";
        assert_eq!(any.type_url, url);
    }
}
