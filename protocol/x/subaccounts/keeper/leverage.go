package keeper

import (
	"encoding/json"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// getLeverageStore returns a prefix store where leverage data is stored.
func (k Keeper) getLeverageStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.LeverageKeyPrefix))
}

// leverageKey returns the store key to retrieve leverage data for a subaccount.
func leverageKey(subaccountId *types.SubaccountId) []byte {
	return subaccountId.ToStateKey()
}

// SetLeverage stores leverage data for a subaccount.
func (k Keeper) SetLeverage(ctx sdk.Context, subaccountId *types.SubaccountId, leverageMap map[uint32]uint32) {
	store := k.getLeverageStore(ctx)
	key := leverageKey(subaccountId)

	// Marshal and store the map using JSON
	b, err := json.Marshal(leverageMap)
	if err != nil {
		panic(err)
	}
	store.Set(key, b)
}

// GetLeverage retrieves leverage data for a subaccount.
func (k Keeper) GetLeverage(ctx sdk.Context, subaccountId *types.SubaccountId) (map[uint32]uint32, bool) {
	store := k.getLeverageStore(ctx)
	key := leverageKey(subaccountId)

	b := store.Get(key)
	if b == nil {
		return nil, false
	}

	var leverageMap map[uint32]uint32
	if err := json.Unmarshal(b, &leverageMap); err != nil {
		panic(err)
	}

	return leverageMap, true
}

// UpdateLeverage updates leverage for specific perpetuals for a subaccount.
func (k Keeper) UpdateLeverage(
	ctx sdk.Context,
	subaccountId *types.SubaccountId,
	perpetualLeverage map[uint32]uint32,
) error {
	// Validate leverage against maximum allowed for each perpetual
	for perpetualId, leverage := range perpetualLeverage {
		maxLeverage, err := k.GetMaxLeverageForPerpetual(ctx, perpetualId)
		if err != nil {
			return errorsmod.Wrapf(
				types.ErrInvalidLeverage,
				"failed to get max leverage for perpetual %d: %v",
				perpetualId,
				err,
			)
		}

		if leverage > maxLeverage {
			return errorsmod.Wrapf(
				types.ErrLeverageExceedsMaximum,
				"leverage %d exceeds maximum allowed %d for perpetual %d",
				leverage,
				maxLeverage,
				perpetualId,
			)
		}
	}

	// Get existing leverage data
	existingLeverage, exists := k.GetLeverage(ctx, subaccountId)
	if !exists {
		existingLeverage = make(map[uint32]uint32)
	}

	// Update with new leverage values
	for perpetualId, leverage := range perpetualLeverage {
		existingLeverage[perpetualId] = leverage
	}

	// Check if the new leverage values break margin requirements
	err := k.checkNewLeverageAgainstMarginRequirements(ctx, subaccountId, perpetualLeverage)
	if err != nil {
		return err
	}

	// Store updated leverage
	k.SetLeverage(ctx, subaccountId, existingLeverage)
	return nil
}

// construct empty updates for each perpetual for which leverage is configured
func (k Keeper) checkNewLeverageAgainstMarginRequirements(
	ctx sdk.Context,
	subaccountId *types.SubaccountId, 
	leverageMap map[uint32]uint32,
) (err error) {
	for perpetualId := range leverageMap {
		update := types.Update{
			SubaccountId: *subaccountId,
			PerpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId: perpetualId,
					BigQuantumsDelta: big.NewInt(0),
				},
			},
		}
		
		// check margin requirements with new leverage configuration
		risk, err := k.GetNetCollateralAndMarginRequirementsWithLeverage(ctx, update, leverageMap)
		if err != nil {
			return err
		}
		if !risk.IsInitialCollateralized() {
			return errorsmod.Wrapf(
				types.ErrLeverageViolatesMarginRequirements,
				"subaccount %s violates margin requirements with new leverage",
				subaccountId.String(),
			)
		}
	}
	return nil
}

// GetMaxLeverageForPerpetual calculates the maximum leverage allowed for a perpetual
// based on its liquidity tier's initial margin requirement.
// Max leverage = 1,000,000 / initial_margin_ppm
func (k Keeper) GetMaxLeverageForPerpetual(ctx sdk.Context, perpetualId uint32) (uint32, error) {
	// Get the perpetual and its liquidity tier in one call
	_, _, liquidityTier, err := k.perpetualsKeeper.GetPerpetualAndMarketPriceAndLiquidityTier(ctx, perpetualId)
	if err != nil {
		return 0, err
	}

	// Calculate max leverage: 1,000,000 / initial_margin_ppm
	if liquidityTier.InitialMarginPpm == 0 {
		return 0, types.ErrInitialMarginPpmIsZero
	}

	maxLeverage := lib.OneMillion / liquidityTier.InitialMarginPpm
	return maxLeverage, nil
}
