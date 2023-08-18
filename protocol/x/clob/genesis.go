package clob

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k *keeper.Keeper, genState types.GenesisState) {
	k.InitializeForGenesis(ctx)

	// Create all `ClobPair` structs.
	for _, elem := range genState.ClobPairs {
		perpetualId, err := elem.GetPerpetualId()
		if err != nil {
			panic(sdkerrors.Wrap(types.ErrInvalidClobPairParameter, err.Error()))
		}
		_, err = k.CreatePerpetualClobPair(
			ctx,
			perpetualId,
			satypes.BaseQuantums(elem.StepBaseQuantums),
			elem.QuantumConversionExponent,
			elem.SubticksPerTick,
			elem.Status,
			elem.MakerFeePpm,
			elem.TakerFeePpm,
		)
		if err != nil {
			panic(err)
		}
	}

	// Create the `LiquidationsConfig` in state, and panic if the genesis state is invalid.
	if err := k.InitializeLiquidationsConfig(ctx, genState.LiquidationsConfig); err != nil {
		panic(err)
	}

	if err := k.InitializeBlockRateLimit(ctx, genState.BlockRateLimitConfig); err != nil {
		panic(err)
	}

	if err := k.InitializeEquityTierLimit(ctx, genState.EquityTierLimitConfig); err != nil {
		panic(err)
	}

	k.InitializeProcessProposerMatchesEvents(ctx)

	// Set the last committed block-time to the genesis time.
	k.SetBlockTimeForLastCommittedBlock(ctx)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Read the CLOB pairs from state.
	genesis.ClobPairs = k.GetAllClobPair(ctx)

	// Read the liquidations config from state.
	genesis.LiquidationsConfig = k.GetLiquidationsConfig(ctx)

	// Read the block rate limit configuration from state.
	genesis.BlockRateLimitConfig = k.GetBlockRateLimitConfiguration(ctx)

	// Read the equity tier limit configuration from state.
	genesis.EquityTierLimitConfig = k.GetEquityTierLimitConfiguration(ctx)

	return genesis
}
