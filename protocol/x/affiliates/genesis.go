package affiliates

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	err := k.UpdateAffiliateTiers(ctx, genState.AffiliateTiers)
	if err != nil {
		panic(err)
	}

	err = k.UpdateAffiliateParameters(ctx, &types.MsgUpdateAffiliateParameters{
		AffiliateParameters: genState.AffiliateParameters,
	})
	if err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	affiliateTiers, err := k.GetAllAffiliateTiers(ctx)
	if err != nil {
		panic(err)
	}

	affiliateParameters, err := k.GetAffiliateParameters(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		AffiliateTiers:      affiliateTiers,
		AffiliateParameters: affiliateParameters,
	}
}
