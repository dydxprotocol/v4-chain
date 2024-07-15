package util

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func CreateUsdcAssetPositions(
	quoteBalance *big.Int,
) []*satypes.AssetPosition {
	return []*satypes.AssetPosition{
		CreateSingleAssetPosition(
			assettypes.AssetUsdc.Id,
			quoteBalance,
		),
	}
}

func CreateSingleAssetPosition(
	assetId uint32,
	quoteBalance *big.Int,
) *satypes.AssetPosition {
	return &satypes.AssetPosition{
		AssetId:  assetId,
		Quantums: dtypes.NewIntFromBigInt(quoteBalance),
	}
}

func CreateUsdcAssetUpdates(
	deltaQuoteBalance *big.Int,
) []satypes.AssetUpdate {
	return []satypes.AssetUpdate{
		{
			AssetId:          assettypes.AssetUsdc.Id,
			BigQuantumsDelta: deltaQuoteBalance,
		},
	}
}

func CreateSinglePerpetualPosition(
	perpetualId uint32,
	quantums *big.Int,
	fundingIndex *big.Int,
	quoteBalance *big.Int,
) *satypes.PerpetualPosition {
	return &satypes.PerpetualPosition{
		PerpetualId:  perpetualId,
		Quantums:     dtypes.NewIntFromBigInt(quantums),
		FundingIndex: dtypes.NewIntFromBigInt(fundingIndex),
		QuoteBalance: dtypes.NewIntFromBigInt(quoteBalance),
	}
}

// Creates a copy of a subaccount and changes the USDC asset position size of the subaccount by the passed in delta.
func ChangeUsdcBalance(subaccount satypes.Subaccount, deltaQuantums int64) satypes.Subaccount {
	subaccountId := satypes.SubaccountId{
		Owner:  subaccount.Id.Owner,
		Number: subaccount.Id.Number,
	}
	assetPositions := make([]*satypes.AssetPosition, 0)
	for _, ap := range subaccount.AssetPositions {
		if ap.AssetId != assettypes.AssetUsdc.Id {
			assetPositions = append(
				assetPositions,
				CreateSingleAssetPosition(ap.AssetId, ap.Quantums.BigInt()),
			)
		} else {
			assetPositions = append(
				assetPositions,
				CreateSingleAssetPosition(
					ap.AssetId,
					new(big.Int).Add(
						ap.Quantums.BigInt(),
						new(big.Int).SetInt64(deltaQuantums),
					),
				),
			)
		}
	}
	if len(assetPositions) == 0 {
		assetPositions = nil
	}
	perpetualPositions := make([]*satypes.PerpetualPosition, 0)
	for _, pp := range subaccount.PerpetualPositions {
		perpetualPositions = append(
			perpetualPositions,
			CreateSinglePerpetualPosition(
				pp.PerpetualId,
				pp.Quantums.BigInt(),
				pp.FundingIndex.BigInt(),
				pp.GetQuoteBalance(),
			),
		)
	}
	if len(perpetualPositions) == 0 {
		perpetualPositions = nil
	}
	return satypes.Subaccount{
		Id:                 &subaccountId,
		AssetPositions:     assetPositions,
		PerpetualPositions: perpetualPositions,
		MarginEnabled:      subaccount.MarginEnabled,
	}
}
