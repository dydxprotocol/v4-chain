package v0_2_2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	clobmodulekeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	clobmodulememclob "github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	clobKeeper *clobmodulekeeper.Keeper,
	indexerEventManager indexer_manager.IndexerEventManager,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Running v0.2.2 hard fork...")

		// Update clob pairs.
		clobPairs := clobKeeper.GetAllClobPairs(ctx)
		for _, clobPair := range clobPairs {
			switch clobPair.Id {
			case 0:
				clobPair.SubticksPerTick = 1e4
				clobPair.QuantumConversionExponent = -8
			case 1:
				clobPair.SubticksPerTick = 1e5
				clobPair.QuantumConversionExponent = -9
			case 2:
				clobPair.SubticksPerTick = 1e6
				clobPair.QuantumConversionExponent = -9
			case 3:
				clobPair.SubticksPerTick = 1e7
				clobPair.QuantumConversionExponent = -10
			case 4:
				clobPair.SubticksPerTick = 1e7
				clobPair.QuantumConversionExponent = -10
			case 5:
				clobPair.SubticksPerTick = 1e5
				clobPair.QuantumConversionExponent = -8
			case 6:
				clobPair.SubticksPerTick = 1e7
				clobPair.QuantumConversionExponent = -10
			case 7:
				clobPair.SubticksPerTick = 1e5
				clobPair.QuantumConversionExponent = -8
			case 8:
				clobPair.SubticksPerTick = 1e6
				clobPair.QuantumConversionExponent = -9
			case 9:
				clobPair.SubticksPerTick = 1e5
				clobPair.QuantumConversionExponent = -8
			case 10:
				clobPair.SubticksPerTick = 1e8
				clobPair.QuantumConversionExponent = -11
			case 11:
				clobPair.SubticksPerTick = 1e6
				clobPair.QuantumConversionExponent = -9
			case 12:
				clobPair.SubticksPerTick = 1e6
				clobPair.QuantumConversionExponent = -9
			case 13:
				clobPair.SubticksPerTick = 1e6
				clobPair.QuantumConversionExponent = -9
			case 14:
				clobPair.SubticksPerTick = 1e4
				clobPair.QuantumConversionExponent = -7
			case 15:
				clobPair.SubticksPerTick = 1e8
				clobPair.QuantumConversionExponent = -11
			case 16:
				clobPair.SubticksPerTick = 1e6
				clobPair.QuantumConversionExponent = -9
			case 17:
				clobPair.SubticksPerTick = 1e3
				clobPair.QuantumConversionExponent = -6
			case 18:
				clobPair.SubticksPerTick = 1e7
				clobPair.QuantumConversionExponent = -10
			case 19:
				clobPair.SubticksPerTick = 1e5
				clobPair.QuantumConversionExponent = -8
			case 20:
				clobPair.SubticksPerTick = 1e5
				clobPair.QuantumConversionExponent = -8
			case 21:
				clobPair.SubticksPerTick = 1e6
				clobPair.QuantumConversionExponent = -9
			case 22:
				clobPair.SubticksPerTick = 1e6
				clobPair.QuantumConversionExponent = -9
			case 23:
				clobPair.SubticksPerTick = 1e6
				clobPair.QuantumConversionExponent = -9
			case 24:
				clobPair.SubticksPerTick = 1e6
				clobPair.QuantumConversionExponent = -9
			case 25:
				clobPair.SubticksPerTick = 1e7
				clobPair.QuantumConversionExponent = -10
			case 26:
				clobPair.SubticksPerTick = 1e6
				clobPair.QuantumConversionExponent = -9
			case 27:
				clobPair.SubticksPerTick = 1e6
				clobPair.QuantumConversionExponent = -9
			case 28:
				clobPair.SubticksPerTick = 1e8
				clobPair.QuantumConversionExponent = -11
			case 29:
				clobPair.SubticksPerTick = 1e7
				clobPair.QuantumConversionExponent = -10
			case 30:
				clobPair.SubticksPerTick = 1e8
				clobPair.QuantumConversionExponent = -11
			case 31:
				clobPair.SubticksPerTick = 1e7
				clobPair.QuantumConversionExponent = -10
			case 32:
				clobPair.SubticksPerTick = 1e7
				clobPair.QuantumConversionExponent = -10
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
		clobKeeper.MemClob = clobmodulememclob.NewMemClobPriceTimePriority(indexerEventManager.Enabled())
		clobKeeper.MemClob.SetClobKeeper(clobKeeper)

		clobKeeper.PerpetualIdToClobPairId = make(map[uint32][]clobtypes.ClobPairId)
		clobKeeper.InitMemClobOrderbooks(ctx)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
