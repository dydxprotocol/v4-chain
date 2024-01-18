use cosmwasm_std::{QuerierWrapper, StdResult};

use crate::proto_structs::MarketPrice;
use crate::query::{DydxQuery, DydxQueryWrapper, SubaccountResponse};
use crate::route::DydxRoute;

/// This is a helper wrapper to easily use our custom queries
pub struct DydxQuerier<'a> {
    querier: &'a QuerierWrapper<'a, DydxQueryWrapper>,
}

impl<'a> DydxQuerier<'a> {
    pub fn new(querier: &'a QuerierWrapper<DydxQueryWrapper>) -> Self {
        DydxQuerier { querier }
    }

    pub fn query_market_price(&self, market_id: u32) -> StdResult<MarketPrice> {
        let request = DydxQueryWrapper {
            route: DydxRoute::Prices,
            query_data: DydxQuery::MarketPrice { id: market_id },
        }
        .into();

        self.querier.query(&request)
    }

    pub fn query_subaccount(&self, owner: String, number: u32) -> StdResult<SubaccountResponse> {
        let request = DydxQueryWrapper {
            route: DydxRoute::Subaccount,
            query_data: DydxQuery::Subaccount { owner, number },
        }
            .into();

        self.querier.query(&request)
    }
}