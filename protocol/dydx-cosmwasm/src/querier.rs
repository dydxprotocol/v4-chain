use cosmwasm_std::{QuerierWrapper, StdResult};
use protobuf::Error;
use crate::proto_structs::{PerpetualClobDetails, LiquidityTier};
use crate::query::{DydxQuery, DydxQueryWrapper};
use crate::route::DydxRoute;
use crate::{MarketPrice, PerpetualClobDetailsResponse, LiquidityTiersResponse, Subaccount, SubaccountResponse, MarketPriceResponse};

/// This is a helper wrapper to easily use our custom queries
pub struct DydxQuerier<'a> {
    querier: &'a QuerierWrapper<'a, DydxQueryWrapper>,
}

impl<'a> DydxQuerier<'a> {
    pub fn new(querier: &'a QuerierWrapper<DydxQueryWrapper>) -> Self {
        DydxQuerier { querier }
    }

    pub fn query_market_price(&self, market_id: u32) -> StdResult<MarketPriceResponse> {
        let request = DydxQueryWrapper {
            route: DydxRoute::MarketPrice,
            query_data: DydxQuery::MarketPrice { id: market_id },
        }
        .into();

        let result: StdResult<MarketPrice> = self.querier.query(&request);
        result.map(|market_price| MarketPriceResponse { market_price })
    }

    pub fn query_subaccount(&self, owner: String, number: u32) -> StdResult<SubaccountResponse> {
        let request = DydxQueryWrapper {
            route: DydxRoute::Subaccount,
            query_data: DydxQuery::Subaccount { 
                owner: owner,
                number: number,
            },
        }
            .into();

        let result: Result<Subaccount, cosmwasm_std::StdError> = self.querier.query::<Subaccount>(&request);
        Ok(SubaccountResponse { subaccount: result? })


    }

    pub fn query_perpetual_clob_details(&self, perpetual_id: u32) -> StdResult<PerpetualClobDetailsResponse> {
        let request = DydxQueryWrapper {
            route: DydxRoute::PerpetualClobDetails,
            query_data: DydxQuery::PerpetualClobDetails { id: perpetual_id },
        }
            .into();
        
        let result: Result<PerpetualClobDetails, cosmwasm_std::StdError> = self.querier.query::<PerpetualClobDetails>(&request);
        Ok(PerpetualClobDetailsResponse { perpetual_clob_details: result? })
    }

    pub fn query_liquidity_tiers(
        &self,
    ) -> StdResult<LiquidityTiersResponse> {
        let request = DydxQueryWrapper {
            route: DydxRoute::LiquidityTiers,
            query_data: DydxQuery::LiquidityTiers,
        }
        .into();

        let result: Result<Vec<LiquidityTier>, cosmwasm_std::StdError> =
            self.querier.query::<Vec<LiquidityTier>>(&request);
        Ok(LiquidityTiersResponse {
            liquidity_tiers: result?,
        })
    }
}