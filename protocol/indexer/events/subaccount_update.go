package events

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	salib "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/lib"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// NewSubaccountUpdateEvent creates a SubaccountUpdateEvent representing a subaccount update
// containing its updated perpetual/asset positions.
func NewSubaccountUpdateEvent(
	subaccountId *satypes.SubaccountId,
	updatedPerpetualPositions []*satypes.PerpetualPosition,
	updatedAssetPositions []*satypes.AssetPosition,
	fundingPayments map[uint32]dtypes.SerializableInt,
) *SubaccountUpdateEventV1 {
	indexerSubaccountId := v1.SubaccountIdToIndexerSubaccountId(*subaccountId)
	return &SubaccountUpdateEventV1{
		SubaccountId: &indexerSubaccountId,
		UpdatedPerpetualPositions: v1.PerpetualPositionsToIndexerPerpetualPositions(
			updatedPerpetualPositions,
			fundingPayments,
		),
		UpdatedAssetPositions: v1.AssetPositionsToIndexerAssetPositions(
			AddQuoteBalanceFromPerpetualPositions(
				updatedPerpetualPositions,
				updatedAssetPositions,
			),
		),
	}
}

func AddQuoteBalanceFromPerpetualPositions(
	perpetualPositions []*satypes.PerpetualPosition,
	assetPositions []*satypes.AssetPosition,
) []*satypes.AssetPosition {
	quoteBalance := new(big.Int)
	for _, position := range perpetualPositions {
		quoteBalance.Add(quoteBalance, position.GetQuoteBalance())
	}

	if quoteBalance.Sign() == 0 {
		return assetPositions
	}

	// Add the quote balance to asset positions.
	return salib.CalculateUpdatedAssetPositions(
		assetPositions,
		[]satypes.AssetUpdate{
			{
				AssetId:          assettypes.AssetUsdc.Id,
				BigQuantumsDelta: quoteBalance,
			},
		},
	)
}
