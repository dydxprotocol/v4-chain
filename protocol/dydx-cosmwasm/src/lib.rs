mod msg;
mod querier;
mod query;
mod route;
mod dydx_types;
mod proto_structs;
mod serializable_int;

pub use msg::{DydxMsg, Transfer, Order, OrderSide, OrderTimeInForce, OrderConditionType, OrderId};
pub use querier::DydxQuerier;
pub use query::{
    MarketPriceResponse, PerpetualClobDetailsResponse, SubaccountResponse, DydxQuery, DydxQueryWrapper};
pub use proto_structs::{
    AssetPosition, ClobPair, MarketPrice, Metadata, PerpetualClobMetadata, PerpetualPosition, SpotClobMetadata, Status, Subaccount, SubaccountId
};
pub use serializable_int::SerializableInt;
pub use route::DydxRoute;

// This export is added to all contracts that import this package, signifying that they require
// "dydx" support on the chain they run on.
#[no_mangle]
extern "C" fn requires_dydx() {}