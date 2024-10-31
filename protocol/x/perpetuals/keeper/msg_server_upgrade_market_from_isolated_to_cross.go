package keeper

import (
	"context"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

func (k msgServer) UpgradeMarketFromIsolatedToCross(
	goCtx context.Context,
	msg *types.MsgUpgradeMarketFromIsolatedToCross,
) (*types.MsgUpgradeMarketFromIsolatedToCrossResponse, error) {
	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	isolatedInsuranceFundAddress, err := k.Keeper.GetInsuranceFundModuleAddress(ctx, msg.PerpetualId)
	if err != nil {
		return nil, err
	}

	_, err = k.Keeper.SetPerpetualMarketType(
		ctx,
		msg.PerpetualId,
		types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
	)
	if err != nil {
		return nil, err
	}

	crossInsuranceFundAddress, err := k.Keeper.GetInsuranceFundModuleAddress(ctx, msg.PerpetualId)
	if err != nil {
		return nil, err
	}

	_, coinToTransfer, err := k.assetsKeeper.ConvertAssetToCoin(
		ctx,
		assettypes.AssetUsdc.Id,
		new(big.Int).Abs(insuranceFundDelta),
	)

	// TODO Move insurance fund for perpetual to primary insurance fund
	return k.bankKeeper.SendCoins(
		ctx,
		isolatedInsuranceFundAddress,
		crossInsuranceFundAddress,
		[]sdk.Coin{coinToTransfer},
	)

	// clob/keeper/deleveraging.go func (k Keeper) GetInsuranceFundBalance(ctx sdk.Context, perpetualId uint32) (balance *big.Int) {

	// TODO Move collateral pool for perpetual to subaccounts module

	// TODO Propagate changes to indexer

	return &types.MsgUpgradeMarketFromIsolatedToCrossResponse{}, nil
}
