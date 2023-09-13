package v0_2_2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	clobmodulekeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	clobKeeper clobmodulekeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Running v0.2.2 hard fork...")

		// Update clob pairs.
		clobPairs := clobKeeper.GetAllClobPairs(ctx)
		for _, clobPair := range clobPairs {
			switch clobPair.Id {
			case 0:
				clobPair.SubticksPerTick = 1
			case 1:
				clobPair.SubticksPerTick = 1
			case 2:
				clobPair.SubticksPerTick = 1
			case 3:
				clobPair.SubticksPerTick = 1
			case 4:
				clobPair.SubticksPerTick = 1
			case 5:
				clobPair.SubticksPerTick = 1
			case 6:
				clobPair.SubticksPerTick = 1
			case 7:
				clobPair.SubticksPerTick = 1
			case 8:
				clobPair.SubticksPerTick = 1
			case 9:
				clobPair.SubticksPerTick = 1
			case 10:
				clobPair.SubticksPerTick = 1
			case 11:
				clobPair.SubticksPerTick = 1
			case 12:
				clobPair.SubticksPerTick = 1
			case 13:
				clobPair.SubticksPerTick = 1
			case 14:
				clobPair.SubticksPerTick = 1
			case 15:
				clobPair.SubticksPerTick = 1
			case 16:
				clobPair.SubticksPerTick = 1
			case 17:
				clobPair.SubticksPerTick = 1
			case 18:
				clobPair.SubticksPerTick = 1
			case 19:
				clobPair.SubticksPerTick = 1
			case 20:
				clobPair.SubticksPerTick = 1
			case 21:
				clobPair.SubticksPerTick = 1
			case 22:
				clobPair.SubticksPerTick = 1
			case 23:
				clobPair.SubticksPerTick = 1
			case 24:
				clobPair.SubticksPerTick = 1
			case 25:
				clobPair.SubticksPerTick = 1
			case 26:
				clobPair.SubticksPerTick = 1
			case 27:
				clobPair.SubticksPerTick = 1
			case 28:
				clobPair.SubticksPerTick = 1
			case 29:
				clobPair.SubticksPerTick = 1
			case 30:
				clobPair.SubticksPerTick = 1
			case 31:
				clobPair.SubticksPerTick = 1
			case 32:
				clobPair.SubticksPerTick = 1
			}
			clobKeeper.UnsafeSetClobPair(
				ctx,
				clobPair,
			)
		}

		// Delete all stateful ordes.
		removedOrderIds := make([]clobtypes.OrderId, 0)
		placedStatefulOrders := clobKeeper.GetAllPlacedStatefulOrders(ctx)
		for _, order := range placedStatefulOrders {
			clobKeeper.MustRemoveStatefulOrder(ctx, order.OrderId)
			removedOrderIds = append(removedOrderIds, order.OrderId)
		}

		untriggeredConditionalOrders := clobKeeper.GetAllUntriggeredConditionalOrders(ctx)
		for _, order := range untriggeredConditionalOrders {
			clobKeeper.MustRemoveStatefulOrder(ctx, order.OrderId)
			removedOrderIds = append(removedOrderIds, order.OrderId)
		}

		// Purge invalid orders from memclob.
		offchainUpdates := clobKeeper.MemClob.PurgeInvalidMemclobState(
			ctx,
			[]clobtypes.OrderId{},
			[]clobtypes.OrderId{},
			[]clobtypes.OrderId{},
			removedOrderIds,
			clobtypes.NewOffchainUpdates(),
		)
		clobKeeper.SendOffchainMessages(offchainUpdates, nil, metrics.SendPrepareCheckStateOffchainUpdates)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
