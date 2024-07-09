package keeper

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	pricetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (k Keeper) CheckCurrentDAIYieldEpochElapsed(ctx sdk.Context) (bool, error) {

	currentEpoch, first := k.CheckFirstDAIYieldEpoch(ctx)
	if first {
		return true, nil
	}

	blockNumber, found := k.GetCurrentDAIYieldEpochBlockNumber(ctx, currentEpoch)
	// this case should never be reached but we return true as the epochs are malconfigured
	// perhaps an epoch was missed
	if !found {
		return false, errorsmod.Wrap(
			types.ErrInvalidSDAIConversionRate,
			"DAI yield epoch not found",
		)
	}

	// get the current block number
	currentBlockNumber := ctx.BlockHeight()

	// check if the current block number is greater than the epoch block number
	if uint64(currentBlockNumber) < blockNumber+uint64(types.DAI_YIELD_MIN_EPOCH_BLOCKS) {
		return true, nil
	}

	return false, nil
}

func (k Keeper) CheckFirstDAIYieldEpoch(ctx sdk.Context) (*big.Int, bool) {

	currentEpoch, found := k.GetCurrentDaiYieldEpochNumber(ctx)
	if !found {
		return nil, true
	}
	return currentEpoch, false
}

func (k Keeper) GetCurrentDAIYieldEpochBlockNumber(ctx sdk.Context, currentEpoch *big.Int) (uint64, bool) {

	params, found := k.GetDaiYieldEpochParams(ctx, currentEpoch.Uint64()%types.DAI_YIELD_ARRAY_SIZE)
	if !found {
		return 0, false
	}

	blockNumber := params.BlockNumber

	return blockNumber, true

}

func (k Keeper) PruneOldDAIYieldEpoch(ctx sdk.Context, newEpoch uint64) error {

	params, found := k.GetDaiYieldEpochParams(ctx, newEpoch%types.DAI_YIELD_ARRAY_SIZE)
	if !found {
		return nil
	}

	err := k.TransferRemainingDAIYieldToInsuranceFund(ctx, params.TradingDaiMinted, params.TotalTradingDaiClaimedForEpoch)
	if err != nil {
		return err
	}

	// no need to explicitly delete the epoch params as they will get overwritten
	return nil

}

func (k Keeper) TransferRemainingDAIYieldToInsuranceFund(ctx sdk.Context, TradingDaiMinted string, TotalTradingDaiClaimedForEpoch string) error {

	tradingDaiMintedAtEpoch, err := ConvertStringToBigInt(TradingDaiMinted)
	if err != nil {
		return err
	}

	tradingDaiClaimedAtEpoch, err := ConvertStringToBigInt(TotalTradingDaiClaimedForEpoch)
	if err != nil {
		return err
	}

	if tradingDaiMintedAtEpoch.Cmp(tradingDaiClaimedAtEpoch) <= 0 {
		return nil
	}

	remainingDai := tradingDaiMintedAtEpoch.Sub(tradingDaiMintedAtEpoch, tradingDaiClaimedAtEpoch)

	tradingDAICoins := sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewIntFromBigInt(remainingDai)))

	if err := k.bankKeeper.SendCoins(ctx, authtypes.NewModuleAddress(types.PoolAccount), perptypes.InsuranceFundModuleAddress, tradingDAICoins); err != nil {
		return errorsmod.Wrap(err, "failed to send trading DAI to the insurance fund")
	}

	return nil
}

func (k Keeper) CreateAndStoreNewDaiYieldEpochParams(ctx sdk.Context) error {

	currentEpoch, found := k.GetCurrentDaiYieldEpochNumber(ctx)
	if !found {
		return errorsmod.Wrap(
			types.ErrInvalidSDAIConversionRate,
			"DAI yield epoch not found",
		)
	}

	newEpoch := currentEpoch.Uint64() + 1

	err := k.PruneOldDAIYieldEpoch(ctx, newEpoch)
	if err != nil {
		return err
	}

	tDAISupply, tradingDaiMinted, err := k.MintYieldGeneratedDuringEpoch(ctx)
	if err != nil {
		return err
	}

	yieldParams := k.CreateNewDaiYieldEpochParams(ctx, tDAISupply, tradingDaiMinted)

	k.SetDaiYieldEpochParams(ctx, newEpoch%uint64(types.DAI_YIELD_ARRAY_SIZE), yieldParams)

	k.SetCurrentDaiYieldEpochNumber(ctx, big.NewInt(int64(newEpoch)))

	return nil
}

func (k Keeper) MintYieldGeneratedDuringEpoch(ctx sdk.Context) (*big.Int, *big.Int, error) {

	// get sDAI supply
	sDAISupplyCoins := k.bankKeeper.GetSupply(ctx, types.SDaiDenom)
	sDAISupply := sDAISupplyCoins.Amount.BigInt()

	tDAISupplyCoins := k.bankKeeper.GetSupply(ctx, types.TradingDAIDenom)
	tDAISupply := tDAISupplyCoins.Amount.BigInt()

	tradingDAIAfterYield, err := k.GetTradingDAIFromSDAIAmount(ctx, sDAISupply)
	if err != nil {
		return nil, nil, err
	}

	if tradingDAIAfterYield.Cmp(tDAISupply) <= 0 {
		return nil, nil, errorsmod.Wrap(
			types.ErrInvalidSDAIConversionRate,
			"Trading DAI supply is less than the sDAI supply",
		)
	}

	tradingDaiToMint := tradingDAIAfterYield.Sub(tradingDAIAfterYield, tDAISupply)
	tradingDAICoins := sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewIntFromBigInt(tradingDaiToMint)))

	if err := k.bankKeeper.MintCoins(
		ctx, types.PoolAccount, tradingDAICoins,
	); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to mint new trading DAI")
	}

	return tDAISupply, tradingDaiToMint, nil
}

func (k Keeper) CreateNewDaiYieldEpochParams(ctx sdk.Context, tradingDaiSupplyBeforeNewEpoch *big.Int, tradingDaiMinted *big.Int) types.DaiYieldEpochParams {

	marketPrices := k.pricesKeeper.GetAllMarketPrices(ctx)

	// Convert []MarketPrice to []*MarketPrice
	marketPricesPtrs := make([]*pricetypes.MarketPrice, len(marketPrices))
	for i := range marketPrices {
		marketPricesPtrs[i] = &marketPrices[i]
	}

	yieldParams := types.DaiYieldEpochParams{
		TradingDaiMinted:               tradingDaiMinted.String(),
		TotalTradingDaiPreMint:         tradingDaiSupplyBeforeNewEpoch.String(),
		TotalTradingDaiClaimedForEpoch: "0",
		BlockNumber:                    uint64(ctx.BlockHeight()),
		EpochMarketPrices:              marketPricesPtrs,
	}

	return yieldParams
}
