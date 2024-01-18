mod msg;
mod querier;
mod query;
mod route;
mod proto_structs;
mod bytes_helper;
mod response_helper;

pub use msg::{SendingMsg, SubaccountId, Transfer};
pub use querier::DydxQuerier;
pub use query::{
    DydxQuery, DydxQueryWrapper,
};
// TODO: Export MarketPriceResponse instead for style consistency.
pub use proto_structs::MarketPrice;
pub use route::DydxRoute;

// This export is added to all contracts that import this package, signifying that they require
// "dydx" support on the chain they run on.
#[no_mangle]
extern "C" fn requires_dydxprotocol() {}