package v_0_3_0

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"
	clobmodulekeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	perpetualsmodulekeeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	pricesmodulekeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
)

var (
	Upgrade = upgrades.Upgrade{
		UpgradeName: UpgradeName,
		StoreUpgrades: store.StoreUpgrades{
			Added: []string{
				evidencetypes.StoreKey,
			},
		},
	}
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	clobKeeper *clobmodulekeeper.Keeper,
	perpetualsKeeper *perpetualsmodulekeeper.Keeper,
	pricesKeeper *pricesmodulekeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Running v0.3.0 Upgrade...")

		const PEPE_ID = 28

		// Update subticks_per_tick and quantum_conversion_exponent for PEPE.
		pepeClobPair, found := clobKeeper.GetClobPair(ctx, PEPE_ID)
		if !found {
			panic("PEPE ClobPair not found")
		}

		pepeClobPair.QuantumConversionExponent = -9
		clobKeeper.UnsafeSetClobPair(ctx, pepeClobPair)

		// Update atomic_resolution for PEPE.
		pepePerpetual, err := perpetualsKeeper.GetPerpetual(ctx, PEPE_ID)
		if err != nil {
			panic(err)
		}

		pepePerpetual.Params.AtomicResolution = 1
		perpetualsKeeper.UnsafeSetPerpetual(ctx, pepePerpetual)

		// Update market price exponent for PEPE.
		pepePrice, err := pricesKeeper.GetMarketPrice(ctx, PEPE_ID)
		if err != nil {
			panic(err)
		}
		pepePrice.Exponent = -16
		pricesKeeper.UnsafeSetMarketPrice(ctx, pepePrice)

		// Update market params exponent for PEPE.
		pepeParams, exists := pricesKeeper.GetMarketParam(ctx, PEPE_ID)
		if !exists {
			panic("PEPE MarketParam not found")
		}
		pepeParams.Exponent = -16
		pricesKeeper.UnsafeModifyMarketParam(ctx, pepeParams)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
