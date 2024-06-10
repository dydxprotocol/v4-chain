use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cosmwasm_std::Uint64;

// TODO(affan): standardize proto compilation

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct MarketPrice {
    pub id: u32,
    pub exponent: i32,
    pub price: Uint64,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct AssetPosition {
    pub asset_id: u32,
    pub quantums: Vec<u8>,
    pub index: Uint64,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct PerpetualPosition {
    pub perpetual_id: u32,
    pub quantums: Vec<u8>,
    pub funding_index: Vec<u8>,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct SubaccountId {
    pub owner: String,
    pub number: u32,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct Subaccount {
    pub id: Option<SubaccountId>,
    pub asset_positions: Vec<AssetPosition>,
    pub perpetual_positions: Vec<PerpetualPosition>,
    pub margin_enabled: bool,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct ClobPair {
    pub id: u32,
    pub metadata: Metadata,
    pub step_base_quantums: u64,
    pub subticks_per_tick: u32,
    pub quantum_conversion_exponent: i32,
    pub status: Status,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
#[serde(tag = "type", content = "value")]
pub enum Metadata {
    PerpetualClobMetadata(PerpetualClobMetadata),
    SpotClobMetadata(SpotClobMetadata),
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct PerpetualClobMetadata {
    pub perpetual_id: u32,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct SpotClobMetadata {
    pub base_asset_id: u32,
    pub quote_asset_id: u32,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub enum Status {
    #[serde(rename = "STATUS_UNSPECIFIED")]
    Unspecified = 0,
    #[serde(rename = "STATUS_ACTIVE")]
    Active = 1,
    #[serde(rename = "STATUS_PAUSED")]
    Paused = 2,
    #[serde(rename = "STATUS_CANCEL_ONLY")]
    CancelOnly = 3,
    #[serde(rename = "STATUS_POST_ONLY")]
    PostOnly = 4,
    #[serde(rename = "STATUS_INITIALIZING")]
    Initializing = 5,
    #[serde(rename = "STATUS_FINAL_SETTLEMENT")]
    FinalSettlement = 6,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct Perpetual {
    pub params: PerpetualParams,
    pub funding_index: Vec<u8>,
    pub open_interest: Vec<u8>,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct PerpetualParams {
    pub id: u32,
    pub ticker: String,
    pub market_id: u32,
    pub atomic_resolution: i32,
    pub default_funding_ppm: i32,
    pub liquidity_tier: u32,
    pub market_type: PerpetualMarketType,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub enum PerpetualMarketType {
    #[serde(rename = "PERPETUAL_MARKET_TYPE_UNSPECIFIED")]
    Unspecified = 0,
    #[serde(rename = "PERPETUAL_MARKET_TYPE_CROSS")]
    Cross = 1,
    #[serde(rename = "PERPETUAL_MARKET_TYPE_ISOLATED")]
    Isolated = 2,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct MarketPremiums {
    pub perpetual_id: u32,
    pub premiums: Vec<i32>,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct PremiumStore {
    pub all_market_premiums: Vec<MarketPremiums>,
    pub num_premiums: u32,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct LiquidityTier {
    pub id: u32,
    pub name: String,
    pub initial_margin_ppm: u32,
    pub maintenance_fraction_ppm: u32,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub base_position_notional: Option<u64>,
    pub impact_notional: u64,
    pub open_interest_lower_cap: u64,
    pub open_interest_upper_cap: u64,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct PerpetualClobDetails {
    pub clob_pair: ClobPair,
    pub perpetual: Perpetual,
}
