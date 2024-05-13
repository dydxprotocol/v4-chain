use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

/// SeiRoute is enum type to represent dydx query route path
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum DydxRoute {
    Oracle,
}