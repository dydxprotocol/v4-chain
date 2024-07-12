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

	epochStartBlockNumber, found := k.GetCurrentDAIYieldEpochBlockNumber(ctx, currentEpoch)
	// this case should never be reached but we return true as the epochs are malconfigured
	// perhaps an epoch was missed
	if !found {
		return false, errorsmod.Wrap(
			types.ErrInvalidSDAIConversionRate,
			"DAI yield epoch not found",
		)
	}

	currentBlockNumber := ctx.BlockHeight()

	if k.DAIYieldEpochHasElapsed(uint64(currentBlockNumber), epochStartBlockNumber) {
		return true, nil
	}

	return false, nil
}

func (k Keeper) DAIYieldEpochHasElapsed(currentBlockNumber uint64, epochStartBlockNumber uint64) bool {
	return currentBlockNumber >= epochStartBlockNumber+uint64(types.DAI_YIELD_MIN_EPOCH_BLOCKS)
}

func (k Keeper) CheckFirstDAIYieldEpoch(ctx sdk.Context) (uint64, bool) {

	currentEpoch, found := k.GetCurrentDaiYieldEpochNumber(ctx)
	if !found {
		return 0, true
	}
	return currentEpoch, false
}

func (k Keeper) GetCurrentDAIYieldEpochBlockNumber(ctx sdk.Context, currentEpoch uint64) (uint64, bool) {
	params, err := k.GetDAIYieldEpochParamsForEpoch(ctx, currentEpoch)
	if err != nil {
		return 0, false
	}

	blockNumber := params.BlockNumber
	return blockNumber, true
}

func (k Keeper) PruneOldDAIYieldEpoch(ctx sdk.Context, newEpoch uint64) error {
	params, err := k.GetDAIYieldEpochParamsForEpoch(ctx, newEpoch)
	if err != nil {
		return err
	}

	err = k.TransferRemainingDAIYieldToInsuranceFund(ctx, params.TradingDaiMinted, params.TotalTradingDaiClaimedForEpoch)
	if err != nil {
		return err
	}

	// no need to explicitly delete the epoch params as they will get overwritten
	return nil

}

func (k Keeper) GetDAIYieldEpochParamsForEpoch(
	ctx sdk.Context,
	epoch uint64,
) (
	params types.DaiYieldEpochParams,
	err error,
) {
	isStored, err := k.isEpochStored(ctx, epoch)
	if err != nil {
		return types.DaiYieldEpochParams{}, err
	}
	if !isStored {
		return types.DaiYieldEpochParams{}, errorsmod.Wrap(types.ErrEpochNotStored, "Could not find epoch info when getting yield epoch params.")
	}

	epochIndex := k.getEpochIndexFromEpoch(epoch)

	params, success := k.GetDaiYieldEpochParams(ctx, epochIndex)
	if !success {
		return types.DaiYieldEpochParams{}, errorsmod.Wrap(types.ErrEpochNotRetrieved, "Could not retrieve epoch info when getting yield epoch params.")
	}
	return params, nil
}

func (k Keeper) SetDAIYieldEpochParamsForEpoch(
	ctx sdk.Context,
	epoch uint64,
	params types.DaiYieldEpochParams,
) (
	err error,
) {
	epochIndex := k.getEpochIndexFromEpoch(epoch)
	k.SetDaiYieldEpochParams(ctx, epochIndex, params)
	return nil
}

func (k Keeper) getEpochIndexFromEpoch(
	epoch uint64,
) (
	epochIndex uint64,
) {
	return epoch % types.MAX_NUM_YIELD_EPOCHS_STORED
}

func (k Keeper) isEpochStored(
	ctx sdk.Context,
	epoch uint64,
) (
	isStored bool,
	err error,
) {
	currEpoch, success := k.GetCurrentDaiYieldEpochNumber(ctx)
	if !success {
		return false, errorsmod.Wrap(types.ErrEpochNotRetrieved, "Could not retrieve yield epoch number when checking if epoch stored")
	}

	if epoch > currEpoch {
		return false, nil
	}

	if currEpoch < types.MAX_NUM_YIELD_EPOCHS_STORED {
		return true, nil
	}

	lastInvalidEpoch := currEpoch - types.MAX_NUM_YIELD_EPOCHS_STORED

	if epoch <= lastInvalidEpoch {
		return false, nil
	}

	return true, nil
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

	newEpoch := currentEpoch + 1

	err := k.PruneOldDAIYieldEpoch(ctx, newEpoch)
	if err != nil {
		return err
	}

	tDAISupply, tradingDaiMinted, yieldCollectedByInsuranceFund, newEpoch, err := k.CalculateYieldParamsForNewEpoch(ctx)
	if err != nil {
		return err
	}

	yieldParams := k.CreateNewDaiYieldEpochParams(ctx, tDAISupply, tradingDaiMinted, yieldCollectedByInsuranceFund)

	k.SetDAIYieldEpochParamsForEpoch(ctx, newEpoch, yieldParams)

	k.SetCurrentDaiYieldEpochNumber(ctx, newEpoch)

	return nil
}

func (k Keeper) CalculateYieldParamsForNewEpoch(ctx sdk.Context) (*big.Int, *big.Int, *big.Int, uint64, error) {
	tDAISupply := new(big.Int)
	tradingDaiMinted := new(big.Int)
	yieldCollectedByInsuranceFund := new(big.Int)
	newEpoch := uint64(0)

	currentEpoch, found := k.GetCurrentDaiYieldEpochNumber(ctx)

	if found {
		newEpoch = currentEpoch + 1

		err := k.PruneOldDAIYieldEpoch(ctx, newEpoch)
		if err != nil {
			return nil, nil, nil, 0, err
		}

		tDAISupply, tradingDaiMinted, yieldCollectedByInsuranceFund, err = k.MintYieldGeneratedDuringEpoch(ctx)
		if err != nil {
			return nil, nil, nil, 0, err
		}
	}

	return tDAISupply, tradingDaiMinted, yieldCollectedByInsuranceFund, newEpoch, nil
}

func (k Keeper) MintYieldGeneratedDuringEpoch(ctx sdk.Context) (*big.Int, *big.Int, *big.Int, error) {

	sDAISupplyCoins := k.bankKeeper.GetSupply(ctx, types.SDaiDenom)
	sDAISupply := sDAISupplyCoins.Amount.BigInt()

	if sDAISupply.Cmp(big.NewInt(0)) <= 0 {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0), nil
	}

	tDAISupplyCoins := k.bankKeeper.GetSupply(ctx, types.TradingDAIDenom)
	tDAISupply := tDAISupplyCoins.Amount.BigInt()

	tradingDAIAfterYield, err := k.GetTradingDAIFromSDAIAmount(ctx, sDAISupply)
	if err != nil {
		return nil, nil, nil, err
	}

	if tradingDAIAfterYield.Cmp(tDAISupply) <= 0 {
		return nil, nil, nil, errorsmod.Wrap(
			types.ErrInvalidSDAIConversionRate,
			"Trading DAI supply is less than the sDAI supply",
		)
	}

	tradingDaiToMint := tradingDAIAfterYield.Sub(tradingDAIAfterYield, tDAISupply)
	tradingDAICoins := sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewIntFromBigInt(tradingDaiToMint)))

	if err := k.bankKeeper.MintCoins(
		ctx, types.PoolAccount, tradingDAICoins,
	); err != nil {
		return nil, nil, nil, errorsmod.Wrap(err, "failed to mint new trading DAI")
	}

	yieldCollectedByInsuranceFund, err := k.CollectYieldForInsuranceFunds(ctx, tradingDaiToMint, tDAISupply)
	if err != nil {
		return nil, nil, nil, err
	}

	return tDAISupply, tradingDaiToMint, yieldCollectedByInsuranceFund, nil
}

func (k Keeper) CollectYieldForInsuranceFunds(ctx sdk.Context, tradingDaiMinted *big.Int, tradingDaiSupplyBeforeNewEpoch *big.Int) (*big.Int, error) {

	perpetuals := k.perpetualsKeeper.GetAllPerpetuals(ctx)

	totalYieldCollected := big.NewInt(0)

	collectedYieldForCrossMarketInsuranceFund := false

	for _, perpetual := range perpetuals {
		isIsolated, err := k.perpetualsKeeper.IsIsolatedPerpetual(ctx, perpetual.Params.Id)
		if err != nil {
			return nil, err
		}

		if !isIsolated && collectedYieldForCrossMarketInsuranceFund {
			continue
		}

		insuranceFundModuleAddress, err := k.perpetualsKeeper.GetInsuranceFundModuleAddress(ctx, perpetual.Params.Id)
		if err != nil {
			return nil, err
		}

		yieldCollected, err := k.CollectYieldForInsuranceFund(ctx, insuranceFundModuleAddress, tradingDaiMinted, tradingDaiSupplyBeforeNewEpoch)
		if err != nil {
			return nil, err
		}

		totalYieldCollected.Add(totalYieldCollected, yieldCollected)

		if !isIsolated {
			collectedYieldForCrossMarketInsuranceFund = true
		}

	}

	return totalYieldCollected, nil
}

func (k Keeper) CollectYieldForInsuranceFund(ctx sdk.Context, address sdk.AccAddress, tradingDaiMinted *big.Int, tradingDaiSupplyBeforeNewEpoch *big.Int) (*big.Int, error) {

	if tradingDaiSupplyBeforeNewEpoch.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0), nil
	}

	balance := k.bankKeeper.GetBalance(ctx, address, types.TradingDAIDenom)

	bigBalance := balance.Amount.BigInt()

	if bigBalance.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0), nil
	}

	// yield = (balance * tradingDaiMinted) / tradingDaiSupplyBeforeNewEpoch
	yield := bigBalance.Mul(bigBalance, tradingDaiMinted)
	yield = yield.Div(yield, tradingDaiSupplyBeforeNewEpoch)

	yieldCoins := sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewIntFromBigInt(yield)))

	if err := k.bankKeeper.SendCoins(ctx, authtypes.NewModuleAddress(types.PoolAccount), address, yieldCoins); err != nil {
		return nil, errorsmod.Wrap(err, "failed to send yield to the insurance fund")
	}

	return yield, nil
}

func (k Keeper) CreateNewDaiYieldEpochParams(ctx sdk.Context, tradingDaiSupplyBeforeNewEpoch *big.Int, tradingDaiMinted *big.Int, yieldCollectedByInsuranceFund *big.Int) types.DaiYieldEpochParams {

	marketPrices := k.pricesKeeper.GetAllMarketPrices(ctx)

	// Convert []MarketPrice to []*MarketPrice
	marketPricesPtrs := make([]*pricetypes.MarketPrice, len(marketPrices))
	for i := range marketPrices {
		marketPricesPtrs[i] = &marketPrices[i]
	}

	yieldParams := types.DaiYieldEpochParams{
		TradingDaiMinted:               tradingDaiMinted.String(),
		TotalTradingDaiPreMint:         tradingDaiSupplyBeforeNewEpoch.String(),
		TotalTradingDaiClaimedForEpoch: yieldCollectedByInsuranceFund.String(),
		BlockNumber:                    uint64(ctx.BlockHeight()),
		EpochMarketPrices:              marketPricesPtrs,
	}

	return yieldParams
}
