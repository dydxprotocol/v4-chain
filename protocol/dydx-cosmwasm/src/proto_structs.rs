use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use crate::serializable_int::SerializableInt;
use serde_repr::{Deserialize_repr, Serialize_repr};

// TODO(OTE-408): standardize proto compilation

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct MarketPrice {
    #[serde(default)]
    pub id: u32,
    pub exponent: i32,
    pub price: u64,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct AssetPosition {
    #[serde(default)]
    pub asset_id: u32,
    pub quantums: SerializableInt,
    #[serde(default)]
    pub index: u64,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct PerpetualPosition {
    #[serde(default)]
    pub perpetual_id: u32,
    pub quantums: SerializableInt,
    pub funding_index: SerializableInt,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct SubaccountId {
    pub owner: String,
    // go uses omit empty, so we need to provide a default value if not set(which is 0 for u32)
    #[serde(default)]
    pub number: u32,
}
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct Subaccount {
    pub id: Option<SubaccountId>,
    #[serde(default)]
    pub asset_positions: Vec<AssetPosition>,
    #[serde(default)]
    pub perpetual_positions: Vec<PerpetualPosition>,
    #[serde(default)]
    pub margin_enabled: bool,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct ClobPair {
    #[serde(default)]
    pub id: u32,
    // metadata first letter is capitalized to match JSON
    #[serde(rename = "Metadata")]
    pub metadata: Metadata,
    pub step_base_quantums: u64,
    pub subticks_per_tick: u32,
    pub quantum_conversion_exponent: i32,
    pub status: Status,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
#[serde(rename_all = "snake_case")] // Ensure field names match JSON case
pub enum Metadata {
    PerpetualClobMetadata(PerpetualClobMetadata),
    SpotClobMetadata(SpotClobMetadata),
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct PerpetualClobMetadata {
    #[serde(default)]
    pub perpetual_id: u32,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct SpotClobMetadata {
    pub base_asset_id: u32,
    pub quote_asset_id: u32,
}

#[derive(Serialize_repr, Deserialize_repr, Clone, Debug, PartialEq, Eq, JsonSchema)]
#[repr(u8)]
#[serde(rename_all = "lowercase")]
pub enum Status {
    Unspecified = 0,
    Active = 1,
    Paused = 2,
    CancelOnly = 3,
    PostOnly = 4,
    Initializing = 5,
    FinalSettlement = 6,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct Perpetual {
    pub params: PerpetualParams,
    pub funding_index: SerializableInt,
    pub open_interest: SerializableInt,
}
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct PerpetualParams {
    #[serde(default)]
    pub id: u32,
    pub ticker: String,
    #[serde(default)]
    pub market_id: u32,
    pub atomic_resolution: i32,
    #[serde(default)]
    pub default_funding_ppm: i32,
    #[serde(default)]
    pub liquidity_tier: u32,
    pub market_type: PerpetualMarketType,
}

#[derive(Serialize_repr, Deserialize_repr, Clone, Debug, PartialEq, Eq, JsonSchema)]
#[repr(u8)]
#[serde(rename_all = "lowercase")]
pub enum PerpetualMarketType {
    Unspecified = 0,
    Cross = 1,
    Isolated = 2,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct PerpetualClobDetails {
    pub perpetual: Perpetual,
    pub clob_pair: ClobPair,
}

/// LiquidityTier stores margin information.
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct LiquidityTier {
    /// Unique id.
    #[serde(default)]
    pub id: u32,
    /// The name of the tier purely for mnemonic purposes, e.g. "Gold".
    #[serde(default)]
    pub name: String,
    /// The margin fraction needed to open a position.
    /// In parts-per-million.
    #[serde(default)]
    pub initial_margin_ppm: u32,
    /// The fraction of the initial-margin that the maintenance-margin is,
    /// e.g. 50%. In parts-per-million.
    #[serde(default)]
    pub maintenance_fraction_ppm: u32,
    /// The impact notional amount (in quote quantums) is used to determine impact
    /// bid/ask prices and its recommended value is 500 USDC / initial margin
    /// fraction.
    /// - Impact bid price = average execution price for a market sell of the
    /// impact notional value.
    /// - Impact ask price = average execution price for a market buy of the
    /// impact notional value.
    #[serde(default)]
    pub impact_notional: u64,
    /// Lower cap for Open Interest Margin Fraction (OIMF), in quote quantums.
    /// IMF is not affected when OI <= open_interest_lower_cap.
    #[serde(default)]
    pub open_interest_lower_cap: u64,
    /// Upper cap for Open Interest Margin Fraction (OIMF), in quote quantums.
    /// IMF scales linearly to 100% as OI approaches open_interest_upper_cap.
    /// If zero, then the IMF does not scale with OI.
    #[serde(default)]
    pub open_interest_upper_cap: u64,
}
