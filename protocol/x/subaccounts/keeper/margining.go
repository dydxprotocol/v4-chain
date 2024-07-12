package keeper

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	perplib "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/lib"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	salib "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// GetMarginedUpdates calculates the quote balance updates needed
// for the given settled updates.
func GetMarginedUpdates(
	settledUpdates []types.SettledUpdate,
	perpInfos perptypes.PerpInfos,
) (
	marginedUpdates []types.SettledUpdate,
) {
	marginedUpdates = make([]types.SettledUpdate, len(settledUpdates))

	for i, update := range settledUpdates {
		marginedUpdates[i] = getMarginedUpdate(update, perpInfos)
	}

	return marginedUpdates
}

// GetMarginedUpdate calculates the quote balance updates needed
// for the given settled updates.
func getMarginedUpdate(
	update types.SettledUpdate,
	perpInfos perptypes.PerpInfos,
) (
	marginedUpdate types.SettledUpdate,
) {
	marginedAssetUpdates := update.GetAssetUpdates()
	marginedPerpetualUpdates := update.GetPerpetualUpdates()

	// Calculate the updated subaccount.
	updatedSubaccount := salib.CalculateUpdatedSubaccount(update, perpInfos)
	updatedPositionMap := make(map[uint32]*types.PerpetualPosition)
	for _, pos := range updatedSubaccount.PerpetualPositions {
		updatedPositionMap[pos.PerpetualId] = pos
	}

	// For each of the updated positions, check if the position is margined and
	// if we need to move any collateral.
	rebalancingNeeded := false
	for _, u := range update.PerpetualUpdates {
		pos := updatedPositionMap[u.PerpetualId]
		if pos == nil {
			continue
		}

		// case 1: the position is fully closed, but there is still some collateral left.
		// move the remaining collateral to the main quote balance.
		if pos.Quantums.Sign() == 0 {
			moveCollateralToMainQuoteBalance(
				marginedAssetUpdates,
				marginedPerpetualUpdates,
				u.PerpetualId,
				pos.GetQuoteBalance(),
			)
			continue
		}

		perpInfo := perpInfos.MustGet(pos.PerpetualId)
		risk := perplib.GetNetCollateralAndMarginRequirements(
			perpInfo.Perpetual,
			perpInfo.Price,
			perpInfo.LiquidityTier,
			pos.GetBigQuantums(),
			pos.GetQuoteBalance(),
		)

		// case 2: the position is undercollateralized w.r.t. the maintenance margin requirement.
		// In this case, we need to rebalance the collateral across all positions.
		if !risk.IsMaintenanceCollateralized() {
			rebalancingNeeded = true
		}
	}

	// Deal with undercollateralized positions if needed.
	if rebalancingNeeded {
		rebalanceCollateralAcrossPositions(
			update.SettledSubaccount,
			marginedAssetUpdates,
			marginedPerpetualUpdates,
			perpInfos,
		)
	}

	r := types.SettledUpdate{
		SettledSubaccount: update.SettledSubaccount,
	}
	if len(marginedAssetUpdates) > 0 {
		r.AssetUpdates = lib.MapToSortedSlice[lib.Sortable[uint32]](marginedAssetUpdates)
	}
	if len(marginedPerpetualUpdates) > 0 {
		r.PerpetualUpdates = lib.MapToSortedSlice[lib.Sortable[uint32]](marginedPerpetualUpdates)
	}
	return r
}

// rebalanceCollateralAcrossPositions rebalances the collateral across all positions
// by first withdrawing as much as possible from the other positions without going below
// their maintenance margin requirements, and then moving collateral to the undercollateralized
// positions.
func rebalanceCollateralAcrossPositions(
	subaccount types.Subaccount,
	assetUpdates map[uint32]types.AssetUpdate,
	perpetualUpdates map[uint32]types.PerpetualUpdate,
	perpInfos perptypes.PerpInfos,
) {
	// Withdraw as much as possible from the other positions without going below
	// their maintenance margin requirements.
	withdrawCollateralFromPerpetualPositions(
		subaccount,
		assetUpdates,
		perpetualUpdates,
		perpInfos,
	)

	// Calculate the updated subaccount after withdrawing collateral.
	updatedSubaccount := salib.CalculateUpdatedSubaccount(
		types.SettledUpdate{
			SettledSubaccount: subaccount,
			AssetUpdates:      lib.MapToSortedSlice[lib.Sortable[uint32]](assetUpdates),
			PerpetualUpdates:  lib.MapToSortedSlice[lib.Sortable[uint32]](perpetualUpdates),
		},
		perpInfos,
	)
	mainQuoteBalance := updatedSubaccount.GetUsdcPosition()

	for _, pos := range updatedSubaccount.PerpetualPositions {
		perpInfo := perpInfos.MustGet(pos.PerpetualId)
		risk := perplib.GetNetCollateralAndMarginRequirements(
			perpInfo.Perpetual,
			perpInfo.Price,
			perpInfo.LiquidityTier,
			pos.GetBigQuantums(),
			pos.GetQuoteBalance(),
		)

		// Move collateral to the position if it is undercollateralized.
		if !risk.IsMaintenanceCollateralized() {
			collateralNeeded := new(big.Int).Sub(risk.MMR, risk.NC)
			collateralToTransfer := lib.BigMin(collateralNeeded, mainQuoteBalance)

			moveCollateralToPosition(assetUpdates, perpetualUpdates, pos.PerpetualId, collateralToTransfer)
			mainQuoteBalance.Sub(mainQuoteBalance, collateralToTransfer)
		}
	}
}

// withdrawCollateralFromPerpetualPositions withdraws all extra collateral from all perpetual positions
// associated with the given subaccount.
// Withdraw as much as possible without going below the maintenance margin.
func withdrawCollateralFromPerpetualPositions(
	subaccount types.Subaccount,
	assetUpdates map[uint32]types.AssetUpdate,
	perpetualUpdates map[uint32]types.PerpetualUpdate,
	perpInfos perptypes.PerpInfos,
) {
	for _, pos := range subaccount.PerpetualPositions {
		perpInfo := perpInfos.MustGet(pos.PerpetualId)
		risk := perplib.GetNetCollateralAndMarginRequirements(
			perpInfo.Perpetual,
			perpInfo.Price,
			perpInfo.LiquidityTier,
			pos.GetBigQuantums(),
			pos.GetQuoteBalance(),
		)

		// Calculate the amount of extra collateral that can be withdrawn.
		// Withdraw as much as possible without going below the maintenance margin.
		extraCollateral := new(big.Int).Sub(risk.NC, risk.MMR)

		if extraCollateral.Sign() > 0 {
			moveCollateralToMainQuoteBalance(
				assetUpdates,
				perpetualUpdates,
				pos.PerpetualId,
				extraCollateral,
			)
		}
	}
}

func moveCollateralToMainQuoteBalance(
	assetUpdates map[uint32]types.AssetUpdate,
	perpetualUpdates map[uint32]types.PerpetualUpdate,
	perpetualId uint32,
	collateral *big.Int,
) {
	moveCollateralToPosition(
		assetUpdates,
		perpetualUpdates,
		perpetualId,
		new(big.Int).Neg(collateral),
	)
}

func moveCollateralToPosition(
	assetUpdates map[uint32]types.AssetUpdate,
	perpetualUpdates map[uint32]types.PerpetualUpdate,
	perpetualId uint32,
	collateral *big.Int,
) {
	if collateral.Sign() == 0 {
		return
	}

	usdcAssetUpdate, ok := assetUpdates[assettypes.AssetUsdc.Id]
	if !ok {
		usdcAssetUpdate = types.AssetUpdate{
			AssetId:          assettypes.AssetUsdc.Id,
			BigQuantumsDelta: new(big.Int),
		}
		assetUpdates[assettypes.AssetUsdc.Id] = usdcAssetUpdate
	}
	usdcAssetUpdate.BigQuantumsDelta.Sub(usdcAssetUpdate.BigQuantumsDelta, collateral)

	perpetualUpdate, ok := perpetualUpdates[perpetualId]
	if !ok {
		perpetualUpdate = types.PerpetualUpdate{
			PerpetualId:          perpetualId,
			BigQuantumsDelta:     new(big.Int),
			BigQuoteBalanceDelta: new(big.Int),
		}
		perpetualUpdates[perpetualId] = perpetualUpdate
	}
	perpetualUpdate.BigQuoteBalanceDelta.Add(
		perpetualUpdate.BigQuoteBalanceDelta,
		collateral,
	)
}
