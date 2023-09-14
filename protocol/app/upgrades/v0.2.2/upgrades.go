package v0_2_2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	clobmodulekeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	clobKeeper *clobmodulekeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Running v0.2.2 hard fork...")

		// Update clob pairs.
		clobPairs := clobKeeper.GetAllClobPairs(ctx)
		for _, clobPair := range clobPairs {
			switch clobPair.Id {
			case 0, 1:
				clobPair.SubticksPerTick = 1e5
				clobPair.QuantumConversionExponent = -9
			case 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17,
				18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32:
				clobPair.SubticksPerTick = 1e6
				clobPair.QuantumConversionExponent = -9
			default:
				panic("unknown clob pair id")
			}
			clobKeeper.UnsafeSetClobPair(
				ctx,
				clobPair,
			)
		}

		// Delete all stateful ordes.
		placedStatefulOrders := clobKeeper.GetAllPlacedStatefulOrders(ctx)
		for _, order := range placedStatefulOrders {
			clobKeeper.MustRemoveStatefulOrder(ctx, order.OrderId)
		}

		untriggeredConditionalOrders := clobKeeper.GetAllUntriggeredConditionalOrders(ctx)
		for _, order := range untriggeredConditionalOrders {
			clobKeeper.MustRemoveStatefulOrder(ctx, order.OrderId)
		}
		clobKeeper.UntriggeredConditionalOrders = make(
			map[clobtypes.ClobPairId]*clobmodulekeeper.UntriggeredConditionalOrders,
		)

		// Update memclob.
		clobKeeper.MemClob.UnsafeResetMemclob(ctx)

		clobKeeper.PerpetualIdToClobPairId = make(map[uint32][]clobtypes.ClobPairId)
		clobKeeper.InitMemClobOrderbooks(ctx)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
