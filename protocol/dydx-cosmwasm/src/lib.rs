mod msg;
mod querier;
mod query;
mod route;
mod proto_structs;
mod bytes_helper;
mod response_helper;

pub use msg::{DydxMsg, SubaccountId, Transfer, Order, OrderSide, OrderTimeInForce, OrderConditionType, OrderId};
pub use querier::DydxQuerier;
pub use query::{
    DydxQuery, DydxQueryWrapper, SubaccountResponse,
};
// TODO: Export MarketPriceResponse instead for style consistency.
pub use proto_structs::MarketPrice;
pub use proto_structs::Subaccount;
pub use route::DydxRoute;

// This export is added to all contracts that import this package, signifying that they require
// "dydx" support on the chain they run on.
#[no_mangle]
extern "C" fn requires_dydxprotocol() {}