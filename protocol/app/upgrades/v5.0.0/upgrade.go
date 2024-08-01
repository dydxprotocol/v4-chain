package v_5_0_0

import (
	"context"
	"fmt"
	"math/big"

	tenderminttypes "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

// perpetuals upgrade handler for v5.0.0
// - 1. Set all perpetuals to cross market type
// - 2. Initialize perpetual open interest to current OI
func perpetualsUpgrade(
	ctx sdk.Context,
	perpetualsKeeper perptypes.PerpetualsKeeper,
	subaccountsKeeper satypes.SubaccountsKeeper,
) {
	perpOIMap := make(map[uint32]*big.Int)

	// Iterate through all subaccounts and perp positions for each subaccount.
	// Aggregate open interest for each perpetual market.
	subaccounts := subaccountsKeeper.GetAllSubaccount(ctx)
	for _, sa := range subaccounts {
		for _, perpPosition := range sa.PerpetualPositions {
			if perpPosition.Quantums.Sign() <= 0 {
				// Only record positive positions for total open interest.
				// Total negative position size should be equal to total positive position size.
				continue
			}
			if openInterest, exists := perpOIMap[perpPosition.PerpetualId]; exists {
				// Already seen this perpetual. Add to open interest.
				openInterest.Add(
					openInterest,
					perpPosition.Quantums.BigInt(),
				)
			} else {
				// Haven't seen this pereptual. Initialize open interest.
				perpOIMap[perpPosition.PerpetualId] = new(big.Int).Set(
					perpPosition.Quantums.BigInt(),
				)
			}
		}
	}

	perpetuals := perpetualsKeeper.GetAllPerpetuals(ctx)
	for _, perp := range perpetuals {
		// 1. Set all perpetuals to cross market type
		perp.Params.MarketType = perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS

		// 2. Initialize perpetual open interest to current OI
		perpOpenInterst := big.NewInt(0)
		if oi, exists := perpOIMap[perp.GetId()]; exists {
			perpOpenInterst.Set(oi)
		}
		perp.OpenInterest = dtypes.NewIntFromBigInt(perpOpenInterst)
		err := perpetualsKeeper.ValidateAndSetPerpetual(
			ctx,
			perp,
		)
		if err != nil {
			panic(fmt.Sprintf(
				"failed to modify open interest for perpetual, openInterest = %v, perpetual = %+v",
				perpOpenInterst,
				perp,
			))
		}
		ctx.Logger().Info(fmt.Sprintf(
			"Successfully migrated perpetual %d: %+v",
			perp.GetId(),
			perp,
		))
	}
}

// blockRateLimitConfigUpdate upgrades the block rate limit. It searches for the
// 1-block window limit for short term and cancellations, sums them, and creates a new
// combined rate limit.
func blockRateLimitConfigUpdate(
	ctx sdk.Context,
	clobKeeper clobtypes.ClobKeeper,
) {
	oldBlockRateLimitConfig := clobKeeper.GetBlockRateLimitConfiguration(ctx)
	ctx.Logger().Info(
		fmt.Sprintf(
			"Combining the short term order placement and cancellation limits of previous config: %+v\n",
			oldBlockRateLimitConfig,
		),
	)
	numAllowedShortTermOrderPlacementsInOneBlock := 0
	numAllowedShortTermOrderCancellationsInOneBlock := 0
	oldShortTermOrderRateLimits := oldBlockRateLimitConfig.MaxShortTermOrdersPerNBlocks
	for _, limit := range oldShortTermOrderRateLimits {
		if limit.NumBlocks == 1 {
			numAllowedShortTermOrderPlacementsInOneBlock += int(limit.Limit)
			break
		}
	}
	if numAllowedShortTermOrderPlacementsInOneBlock == 0 {
		panic("Failed to find MaxShortTermOrdersPerNBlocks with window 1.")
	}

	oldShortTermOrderCancellationRateLimits := oldBlockRateLimitConfig.MaxShortTermOrderCancellationsPerNBlocks
	for _, limit := range oldShortTermOrderCancellationRateLimits {
		if limit.NumBlocks == 1 {
			numAllowedShortTermOrderCancellationsInOneBlock += int(limit.Limit)
			break
		}
	}
	if numAllowedShortTermOrderCancellationsInOneBlock == 0 {
		panic("Failed to find MaxShortTermOrdersPerNBlocks with window 1.")
	}

	allowedNumShortTermPlaceAndCancelInFiveBlocks :=
		(numAllowedShortTermOrderPlacementsInOneBlock + numAllowedShortTermOrderCancellationsInOneBlock) * 5

	blockRateLimitConfig := clobtypes.BlockRateLimitConfiguration{
		// Kept the same
		MaxStatefulOrdersPerNBlocks: oldBlockRateLimitConfig.MaxStatefulOrdersPerNBlocks,
		// Combine place and cancel, gate over 5 blocks to allow burst
		MaxShortTermOrdersAndCancelsPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
			{
				NumBlocks: 5,
				Limit:     uint32(allowedNumShortTermPlaceAndCancelInFiveBlocks),
			},
		},
	}
	ctx.Logger().Info(
		fmt.Sprintf(
			"Attempting to set rate limiting config to newly combined config: %+v\n",
			blockRateLimitConfig,
		),
	)
	if err := clobKeeper.InitializeBlockRateLimit(ctx, blockRateLimitConfig); err != nil {
		panic(fmt.Sprintf("failed to update the block rate limit configuration: %s", err))
	}
	ctx.Logger().Info(
		fmt.Sprintf(
			"Successfully upgraded block rate limit configuration to: %+v\n",
			clobKeeper.GetBlockRateLimitConfiguration(ctx),
		),
	)
}

func negativeTncSubaccountSeenAtBlockUpgrade(
	ctx sdk.Context,
	perpetualsKeeper perptypes.PerpetualsKeeper,
	subaccountsKeeper satypes.SubaccountsKeeper,
) {
	// Get block height stored by v4.x.x.
	blockHeight, exists := subaccountsKeeper.LegacyGetNegativeTncSubaccountSeenAtBlock(ctx)
	ctx.Logger().Info(
		fmt.Sprintf(
			"Retrieved block height from store for negative tnc subaccount seen at block: %d, exists: %t\n",
			blockHeight,
			exists,
		),
	)
	// If no block height was stored in the legacy store, no migration needed.
	if !exists {
		return
	}

	// If there are no perpetuals, then no new state needs to be stored, as there can be no
	// negative tnc subaccounts w/o perpetuals.
	perpetuals := perpetualsKeeper.GetAllPerpetuals(ctx)
	ctx.Logger().Info(
		fmt.Sprintf(
			"Retrieved all perpetuals for negative tnc subaccount migration, # of perpetuals is %d\n",
			len(perpetuals),
		),
	)
	if len(perpetuals) == 0 {
		return
	}

	ctx.Logger().Info(
		fmt.Sprintf(
			"Migrating negative tnc subaccount seen store, storing block height %d for perpetual %d\n",
			perpetuals[0].Params.Id,
			blockHeight,
		),
	)
	// Migrate the value from the legacy store to the new store.
	if err := subaccountsKeeper.SetNegativeTncSubaccountSeenAtBlock(
		ctx,
		perpetuals[0].Params.Id, // must be a cross-margined perpetual due to `perpetualsUpgrade`.
		blockHeight,
	); err != nil {
		panic(fmt.Sprintf("failed to set negative tnc subaccount seen at block with value %d: %s", blockHeight, err))
	}
	ctx.Logger().Info(
		fmt.Sprintf(
			"Successfully migrated negative tnc subaccount seen at block with block height %d\n",
			blockHeight,
		),
	)
}

// Initialize soft and upper caps for OIMF
func initializeOIMFCaps(
	ctx sdk.Context,
	perpetualsKeeper perptypes.PerpetualsKeeper,
) {
	allLiquidityTiers := perpetualsKeeper.GetAllLiquidityTiers(ctx)
	for _, tier := range allLiquidityTiers {
		if tier.Id == 0 {
			// For large cap, no OIMF caps
			tier.OpenInterestLowerCap = 0
			tier.OpenInterestUpperCap = 0
		} else if tier.Id == 1 {
			// Mid cap
			tier.OpenInterestLowerCap = 20_000_000_000_000 // 20 million USDC
			tier.OpenInterestUpperCap = 50_000_000_000_000 // 50 million USDC
		} else if tier.Id == 2 {
			// Long tail
			tier.OpenInterestLowerCap = 5_000_000_000_000  // 5 million USDC
			tier.OpenInterestUpperCap = 10_000_000_000_000 // 10 million USDC
		} else {
			// Safety
			tier.OpenInterestLowerCap = 2_000_000_000_000 // 2 million USDC
			tier.OpenInterestUpperCap = 5_000_000_000_000 // 5 million USDC
		}

		lt, err := perpetualsKeeper.SetLiquidityTier(
			ctx,
			tier.Id,
			tier.Name,
			tier.InitialMarginPpm,
			tier.MaintenanceFractionPpm,
			tier.ImpactNotional,
			tier.OpenInterestLowerCap,
			tier.OpenInterestUpperCap,
		)
		if err != nil {
			panic(fmt.Sprintf("failed to set liquidity tier: %+v,\n err: %s", tier.Id, err))
		}
		// TODO(OTE-248): Optional - emit indexer events that for updated liquidity tier
		ctx.Logger().Info(
			fmt.Sprintf(
				"Successfully set liqiquidity tier with `OpenInterestLower/UpperCap`: %+v\n",
				lt,
			),
		)
		// TODO(OTE-249): Add upgrade test that checks if the OIMF caps are set correctly
	}
}

func voteExtensionsUpgrade(
	ctx sdk.Context,
	keeper consensusparamkeeper.Keeper,
) {
	currentParams, err := keeper.Params(ctx, &consensustypes.QueryParamsRequest{})
	if err != nil || currentParams == nil || currentParams.Params == nil {
		panic(fmt.Sprintf("failed to retrieve existing consensus params in VE upgrade handler: %s", err))
	}
	currentParams.Params.Abci = &tenderminttypes.ABCIParams{
		VoteExtensionsEnableHeight: ctx.BlockHeight() + VEEnableHeightDelta,
	}
	_, err = keeper.UpdateParams(ctx, &consensustypes.MsgUpdateParams{
		Authority: keeper.GetAuthority(),
		Block:     currentParams.Params.Block,
		Evidence:  currentParams.Params.Evidence,
		Validator: currentParams.Params.Validator,
		Abci:      currentParams.Params.Abci,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to update consensus params : %s", err))
	}
	ctx.Logger().Info(
		"Successfully set VoteExtensionsEnableHeight",
		"consensus_params",
		currentParams.Params.String(),
	)
}

// Deprecated: Initialize vault module parameters.
// This function is deprecated because `Params` in `x/vault` is replaced with `DefaultQuotingParams` in v6.x.
func initializeVaultModuleParams(
	ctx sdk.Context,
	vaultKeeper vaulttypes.VaultKeeper,
) {
	// // Set vault module parameters to default values.
	// vaultParams := vaulttypes.DefaultParams()
	// if err := vaultKeeper.SetParams(ctx, vaultParams); err != nil {
	// 	panic(fmt.Sprintf("failed to set vault module parameters: %s", err))
	// }
	// ctx.Logger().Info("Successfully set vault module parameters")
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	perpetualsKeeper perptypes.PerpetualsKeeper,
	clobKeeper clobtypes.ClobKeeper,
	subaccountsKeeper satypes.SubaccountsKeeper,
	consensusParamKeeper consensusparamkeeper.Keeper,
	vaultKeeper vaulttypes.VaultKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		// Migrate pruneable orders to new format
		clobKeeper.MigratePruneableOrders(sdkCtx)
		sdkCtx.Logger().Info("Successfully migrated pruneable orders")

		// Set all perpetuals to cross market type and initialize open interest
		perpetualsUpgrade(sdkCtx, perpetualsKeeper, subaccountsKeeper)

		// Set block rate limit configuration
		blockRateLimitConfigUpdate(sdkCtx, clobKeeper)

		// Migrate state from legacy store for negative tnc subaccount seen to new store for
		// negative tnc subaccount seen.
		// Note, must be done after the upgrade to perpetuals to cross market type.
		negativeTncSubaccountSeenAtBlockUpgrade(sdkCtx, perpetualsKeeper, subaccountsKeeper)
		// Initialize liquidity tier with lower and upper OI caps.
		initializeOIMFCaps(sdkCtx, perpetualsKeeper)

		// Set vote extension enable height
		voteExtensionsUpgrade(sdkCtx, consensusParamKeeper)

		// Initialize `x/vault` module params.
		initializeVaultModuleParams(sdkCtx, vaultKeeper)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
