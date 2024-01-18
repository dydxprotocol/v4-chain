use cosmwasm_std::Uint64;
use num_bigint::BigInt;

use crate::bytes_helper::bytes_to_bigint;
use crate::proto_structs::{Subaccount, SubaccountId};

#[derive(Clone, Debug, PartialEq)]
pub struct AssetPositionBigInt {
    pub asset_id: u32,
    pub quantums: BigInt,
    pub index: Uint64,
}

#[derive(Clone, Debug, PartialEq)]
pub struct PerpetualPositionBigInt {
    pub perpetual_id: u32,
    pub quantums: BigInt,
    pub funding_index: BigInt,
}

#[derive(Clone, Debug, PartialEq)]
pub struct SubaccountBigInt {
    pub id: Option<SubaccountId>,
    pub asset_positions: Vec<AssetPositionBigInt>,
    pub perpetual_positions: Vec<PerpetualPositionBigInt>,
    pub margin_enabled: bool,
}

// Function to transform Subaccount to SubaccountBigInt
fn transform_subaccount(subaccount: Subaccount) -> SubaccountBigInt {
    SubaccountBigInt {
        id: subaccount.id,
        asset_positions: subaccount.asset_positions.into_iter().map(|ap| {
            AssetPositionBigInt {
                asset_id: ap.asset_id,
                quantums: bytes_to_bigint(ap.quantums),
                index: ap.index,
            }
        }).collect(),
        perpetual_positions: subaccount.perpetual_positions.into_iter().map(|pp| {
            PerpetualPositionBigInt {
                perpetual_id: pp.perpetual_id,
                quantums: bytes_to_bigint(pp.quantums),
                funding_index: bytes_to_bigint(pp.funding_index),
            }
        }).collect(),
        margin_enabled: subaccount.margin_enabled,
    }
}

#[cfg(test)]
mod tests {
    use crate::proto_structs::{AssetPosition, PerpetualPosition};
    use super::*;

    #[test]
    fn test_transform_subaccount() {
        let mock_subaccount = Subaccount {
            id: Some(SubaccountId {
                owner: "owner_address".to_string(),
                number: 123,
            }),
            asset_positions: vec![
                AssetPosition {
                    asset_id: 1,
                    quantums: vec![0x02, 0x01], // Represents BigInt 1
                    index: Uint64::from(10u64),
                },
                AssetPosition {
                    asset_id: 2,
                    quantums: vec![0x03, 0x01], // Represents BigInt -1
                    index: Uint64::from(20u64),
                },
            ],
            perpetual_positions: vec![
                PerpetualPosition {
                    perpetual_id: 1,
                    quantums: vec![0x02, 0xFF], // Represents BigInt 255
                    funding_index: vec![0x03, 0xFF], // Represents BigInt -255
                },
            ],
            margin_enabled: true,
        };

        // Expected result
        let expected = SubaccountBigInt {
            id: Some(SubaccountId {
                owner: "owner_address".to_string(),
                number: 123,
            }),
            asset_positions: vec![
                AssetPositionBigInt {
                    asset_id: 1,
                    quantums: BigInt::from(1),
                    index: Uint64::from(10u64),
                },
                AssetPositionBigInt {
                    asset_id: 2,
                    quantums: BigInt::from(-1),
                    index: Uint64::from(20u64),
                },
            ],
            perpetual_positions: vec![
                PerpetualPositionBigInt {
                    perpetual_id: 1,
                    quantums: BigInt::from(255),
                    funding_index: BigInt::from(-255),
                },
            ],
            margin_enabled: true,
        };

        // Transform the mock data
        let result = transform_subaccount(mock_subaccount);

        // Assert the transformed data matches the expected result
        assert_eq!(result, expected);
    }
}
