mod querier;
mod query;
mod route;
mod proto_structs;
mod bytes_helper;
mod response_helper;

pub use querier::DydxQuerier;
pub use query::{
    MarketPriceResponse, DydxQuery, DydxQueryWrapper,
};
pub use proto_structs::MarketPrice;
pub use route::DydxRoute;

// This export is added to all contracts that import this package, signifying that they require
// "dydx" support on the chain they run on.
#[no_mangle]
extern "C" fn requires_dydx() {}