package keeper

import (
	"sort"

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

	var entries []*types.PerpetualLeverageEntry
	for perpetualId, CustomImfPpm := range leverageMap {
		entries = append(entries, &types.PerpetualLeverageEntry{
			PerpetualId:  perpetualId,
			CustomImfPpm: CustomImfPpm,
		})
	}

	// Sort the entries to ensure deterministic ordering in state
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].PerpetualId < entries[j].PerpetualId
	})

	leverageData := &types.LeverageData{
		Entries: entries,
	}

	b := k.cdc.MustMarshal(leverageData)
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

	var leverageData types.LeverageData
	k.cdc.MustUnmarshal(b, &leverageData)

	leverageMap := make(map[uint32]uint32)
	for _, entry := range leverageData.Entries {
		leverageMap[entry.PerpetualId] = entry.CustomImfPpm
	}

	return leverageMap, true
}

// UpdateLeverage updates leverage for specific perpetuals for a subaccount.
func (k Keeper) UpdateLeverage(
	ctx sdk.Context,
	subaccountId *types.SubaccountId,
	perpetualLeverage map[uint32]uint32,
) error {
	// Sort the perpetual IDs to ensure deterministic ordering
	sortedPerpIds := lib.GetSortedKeys[lib.Sortable[uint32]](perpetualLeverage)

	// Validate leverage against maximum allowed leverage for each perpetual
	for _, perpetualId := range sortedPerpIds {
		custom_imf_ppm := perpetualLeverage[perpetualId]
		minImfPpm, err := k.GetMinImfForPerpetual(ctx, perpetualId)
		if err != nil {
			return errorsmod.Wrapf(
				types.ErrInvalidLeverage,
				"failed to get max leverage for perpetual %d: %v",
				perpetualId,
				err,
			)
		}

		if custom_imf_ppm < minImfPpm {
			return errorsmod.Wrapf(
				types.ErrLeverageExceedsMaximum,
				"%d is less than minimum allowed imf (%d) for perpetual %d resulting in higher than allowed leverage",
				custom_imf_ppm,
				minImfPpm,
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
	for perpetualId, custom_imf_ppm := range perpetualLeverage {
		existingLeverage[perpetualId] = custom_imf_ppm
	}

	// Store updated leverage
	k.SetLeverage(ctx, subaccountId, existingLeverage)
	return nil
}

// GetMinImfForPerpetual returns the IMF ppm allowed for a perpetual
// based on its liquidity tier's initial margin requirement.
func (k Keeper) GetMinImfForPerpetual(ctx sdk.Context, perpetualId uint32) (uint32, error) {
	// Get the perpetual and its liquidity tier in one call
	_, _, liquidityTier, err := k.perpetualsKeeper.GetPerpetualAndMarketPriceAndLiquidityTier(ctx, perpetualId)
	if err != nil {
		return 0, err
	}

	// Calculate max leverage: 1,000,000 / initial_margin_ppm
	if liquidityTier.InitialMarginPpm == 0 {
		return 0, types.ErrInitialMarginPpmIsZero
	}

	return liquidityTier.InitialMarginPpm, nil
}
