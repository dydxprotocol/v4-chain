package keeper

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/lib/log"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// getValidSubaccountUpdatesForTransfer generates subaccount updates and check
// for validity with `CanUpdateSubaccount()`
// Returns the subaccount updates if check is successful.
func (k Keeper) getValidSubaccountUpdatesForTransfer(
	ctx sdk.Context,
	subaccountId types.SubaccountId,
	assetId uint32,
	quantums *big.Int,
	isToSubaccount bool,
) (
	updates []types.Update,
	err error,
) {
	bigBalanceDelta := new(big.Int).Set(quantums)
	if !isToSubaccount {
		bigBalanceDelta.Neg(bigBalanceDelta)
	}

	if assetId == 0 {
		updates = []types.Update{
			{
				SubaccountId: subaccountId,
				AssetUpdates: []types.AssetUpdate{
					{
						AssetId:          assettypes.AssetUsdc.Id,
						BigQuantumsDelta: bigBalanceDelta,
					},
				},
			},
		}
	} else {
		// TODO(DEC-715): Support non-USDC assets.
		return nil, types.ErrAssetTransferThroughBankNotImplemented
	}

	success, successPerUpdate, err := k.CanUpdateSubaccounts(ctx, updates, types.Transfer)
	if err != nil {
		return nil, err
	}

	// If not successful, return error indicating why.
	if err := types.GetErrorFromUpdateResults(success, successPerUpdate, updates); err != nil {
		return nil, err
	}

	return updates, nil
}

// applyValidSubaccountUpdatesForTransfer updates the subaccount by either
// debiting or crediting the subaccount balance.
// Panics if the update fails, as this function assumes the corresponding
// bankKeeper update was successful.
func (k Keeper) applyValidSubaccountUpdateForTransfer(
	ctx sdk.Context,
	updates []types.Update,
	updateType types.UpdateType,
) error {
	// Update subaccount to reflect the transfer.
	success, successPerUpdate, err := k.UpdateSubaccounts(ctx, updates, updateType)

	// Neither of the two conditions below should be true, since `k.CanUpdateSubaccount()`
	// already succeeded.
	if err != nil {
		return err
	}

	return types.GetErrorFromUpdateResults(success, successPerUpdate, updates)
}

// DepositFundsFromAccountToSubaccount returns an error if the call to `k.CanUpdateSubaccounts()`
// fails. Otherwise, increases the asset quantums in the subaccount, translates the
// `assetId` and `quantums` into a `sdk.Coin`, and calls
// `bankKeeper.SendCoinsFromAccountToModule()`.
// TODO(CORE-168): Change function interface to accept `denom` and `amount` instead of `assetId` and `quantums`.
func (k Keeper) DepositFundsFromAccountToSubaccount(
	ctx sdk.Context,
	fromAccount sdk.AccAddress,
	toSubaccountId types.SubaccountId,
	assetId uint32,
	quantums *big.Int,
) error {
	// TODO(DEC-715): Support non-USDC assets.
	if assetId != assettypes.AssetUsdc.Id {
		return types.ErrAssetTransferThroughBankNotImplemented
	}

	if quantums.Sign() <= 0 {
		return errorsmod.Wrap(types.ErrAssetTransferQuantumsNotPositive, lib.UintToString(assetId))
	}

	convertedQuantums, coinToTransfer, err := k.assetsKeeper.ConvertAssetToCoin(
		ctx,
		assetId,
		quantums,
	)
	if err != nil {
		return err
	}

	// Generate subaccount updates and check whether updates can be applied.
	updates, err := k.getValidSubaccountUpdatesForTransfer(
		ctx,
		toSubaccountId,
		assetId,
		convertedQuantums,
		true, // isToSubaccount
	)
	if err != nil {
		return err
	}

	collateralPoolAddr, err := k.GetCollateralPoolForSubaccount(ctx, toSubaccountId)
	if err != nil {
		return err
	}

	// Send coins from `fromModule` to the `subaccounts` module account.
	if err := k.bankKeeper.SendCoins(
		ctx,
		fromAccount,
		collateralPoolAddr,
		[]sdk.Coin{coinToTransfer},
	); err != nil {
		return err
	}

	// Apply subaccount updates.
	return k.applyValidSubaccountUpdateForTransfer(
		ctx,
		updates,
		types.Deposit,
	)
}

// WithdrawFundsFromSubaccountToAccount returns an error if the call to `k.CanUpdateSubaccounts()`
// fails. Otherwise, deducts the asset quantums from the subaccount, translates the
// `assetId` and `quantums` into a `sdk.Coin`, and calls `bankKeeper.SendCoins()`.
func (k Keeper) WithdrawFundsFromSubaccountToAccount(
	ctx sdk.Context,
	fromSubaccountId types.SubaccountId,
	toAccount sdk.AccAddress,
	assetId uint32,
	quantums *big.Int,
) error {
	// TODO(DEC-715): Support non-USDC assets.
	if assetId != assettypes.AssetUsdc.Id {
		return types.ErrAssetTransferThroughBankNotImplemented
	}

	if quantums.Sign() <= 0 {
		return errorsmod.Wrap(types.ErrAssetTransferQuantumsNotPositive, lib.UintToString(assetId))
	}

	convertedQuantums, coinToTransfer, err := k.assetsKeeper.ConvertAssetToCoin(
		ctx,
		assetId,
		quantums,
	)
	if err != nil {
		return err
	}

	// Generate subaccount updates and check whether updates can be applied.
	updates, err := k.getValidSubaccountUpdatesForTransfer(
		ctx,
		fromSubaccountId,
		assetId,
		convertedQuantums,
		false, // isToSubaccount
	)
	if err != nil {
		return err
	}

	collateralPoolAddr, err := k.GetCollateralPoolForSubaccount(ctx, fromSubaccountId)
	if err != nil {
		return err
	}

	// Send coins from `fromModule` to the `subaccounts` module account.
	if err := k.bankKeeper.SendCoins(
		ctx,
		collateralPoolAddr,
		toAccount,
		[]sdk.Coin{coinToTransfer},
	); err != nil {
		return err
	}

	// Apply subaccount updates.
	return k.applyValidSubaccountUpdateForTransfer(
		ctx,
		updates,
		types.Withdrawal,
	)
}

// DistributeFees calculates the market mapper revenue share and fee collector share
// based on the quantums and perpetual parameters, and transfers the fees to the
// market mapper and fee collector.
func (k Keeper) DistributeFees(
	ctx sdk.Context,
	assetId uint32,
	quantums *big.Int,
	perpetualId uint32,
) error {
	// get perpetual
	perpetual, err := k.perpetualsKeeper.GetPerpetual(ctx, perpetualId)
	if err != nil {
		return err
	}

	collateralPoolAddr, err := k.GetCollateralPoolFromPerpetualId(ctx, perpetualId)
	if err != nil {
		return err
	}

	// calculate market mapper rev share
	marketMapperShare := big.NewInt(0)
	revShareAddr, revSharePpm, err := k.revShareKeeper.GetMarketMapperRevenueShareForMarket(
		ctx,
		perpetual.Params.MarketId,
	)
	if err == nil && revShareAddr != nil {
		if revSharePpm >= 1e6 {
			log.ErrorLog(
				ctx,
				"DistributeFees: revSharePpm is greater than or equal to 100%",
				"revSharePpm",
				revSharePpm,
			)
		} else {
			// marketMapperShare = quantums * revSharePpm / 1e6
			marketMapperShare.Div(
				new(big.Int).Mul(quantums, big.NewInt(int64(revSharePpm))),
				big.NewInt(1e6),
			)
		}
	}

	// Remaining amount goes to the fee collector
	feeCollectorShare := new(big.Int).Sub(quantums, marketMapperShare)

	// Transfer fees to the market mapper
	// TODO: add monitoring to record the amount of fees transferred to the market mapper
	if err := k.TransferFees(
		ctx,
		assetId,
		collateralPoolAddr,
		revShareAddr,
		marketMapperShare,
	); err != nil {
		return err
	}

	// Transfer fees to the fee collector
	if err := k.TransferFees(
		ctx,
		assetId,
		collateralPoolAddr,
		authtypes.NewModuleAddress(authtypes.FeeCollectorName),
		feeCollectorShare,
	); err != nil {
		return err
	}

	return nil
}

// TransferFees translates the assetId and quantums into a sdk.Coin, and moves the funds from
// `fromAddr` to `toAddr` by calling `bankKeeper.SendCoins()`
func (k Keeper) TransferFees(
	ctx sdk.Context,
	assetId uint32,
	fromAddr sdk.AccAddress,
	toAddr sdk.AccAddress,
	quantums *big.Int,
) error {
	// TODO(DEC-715): Support non-USDC assets.
	if assetId != assettypes.AssetUsdc.Id {
		return types.ErrAssetTransferThroughBankNotImplemented
	}

	if quantums.Sign() == 0 {
		return nil
	}

	_, coinToTransfer, err := k.assetsKeeper.ConvertAssetToCoin(
		ctx,
		assetId,
		new(big.Int).Abs(quantums),
	)
	if err != nil {
		return err
	}

	if quantums.Sign() < 0 {
		// In the case of a liquidation, net fees can be negative if the maker gets a rebate.
		fromAddr, toAddr = toAddr, fromAddr
	}

	// Send coins from `fromAddr` to `toAddr`
	if err := k.bankKeeper.SendCoins(
		ctx,
		fromAddr,
		toAddr,
		[]sdk.Coin{coinToTransfer},
	); err != nil {
		return err
	}

	return nil
}

// TransferInsuranceFundPayments transfers funds in and out of the insurance fund to the subaccounts
// module by calling `bankKeeper.SendCoins`.
// This function transfers funds
//   - from the insurance fund to the subaccounts module when `insuranceFundDelta` is negative.
//   - from the subaccounts module to the insurance fund when `insuranceFundDelta` is positive.
//   - does nothing if `insuranceFundDelta` is zero.
//
// If the sender account does not have enough balance for the transfer, an error is returned.
// Note this function does not change any individual subaccount state.
func (k Keeper) TransferInsuranceFundPayments(
	ctx sdk.Context,
	insuranceFundDelta *big.Int,
	perpetualId uint32,
) error {
	if insuranceFundDelta.Sign() == 0 {
		return nil
	}

	_, coinToTransfer, err := k.assetsKeeper.ConvertAssetToCoin(
		ctx,
		assettypes.AssetUsdc.Id,
		new(big.Int).Abs(insuranceFundDelta),
	)
	if err != nil {
		// Panic if USDC does not exist.
		panic(err)
	}

	// Determine the sender and receiver.
	// Send coins from `subaccounts` to the `insurance_fund` module account by default.
	fromModule, err := k.GetCollateralPoolFromPerpetualId(ctx, perpetualId)
	if err != nil {
		panic(err)
	}
	toModule, err := k.perpetualsKeeper.GetInsuranceFundModuleAddress(ctx, perpetualId)
	if err != nil {
		panic(err)
	}

	if insuranceFundDelta.Sign() < 0 {
		// Insurance fund needs to cover losses from liquidations.
		// Send coins from the insurance fund to the `subaccounts` module account.
		fromModule, toModule = toModule, fromModule
	}

	// Use SendCoins API instead of SendCoinsFromModuleToModule since we don't need the
	// module account features
	return k.bankKeeper.SendCoins(
		ctx,
		fromModule,
		toModule,
		[]sdk.Coin{coinToTransfer},
	)
}

// TransferFundsFromSubaccountToSubaccount returns an error if the call to `k.CanUpdateSubaccounts()`
// fails. Otherwise, updates the asset quantums in the subaccounts, translates the
// `assetId` and `quantums` into a `sdk.Coin`, and call `bankKeeper.SendCoins()` if the collateral
// pools for the two subaccounts are different.
// TODO(CORE-168): Change function interface to accept `denom` and `amount` instead of `assetId` and
// `quantums`.
func (k Keeper) TransferFundsFromSubaccountToSubaccount(
	ctx sdk.Context,
	senderSubaccountId types.SubaccountId,
	recipientSubaccountId types.SubaccountId,
	assetId uint32,
	quantums *big.Int,
) error {
	// TODO(DEC-715): Support non-USDC assets.
	if assetId != assettypes.AssetUsdc.Id {
		return types.ErrAssetTransferThroughBankNotImplemented
	}

	updates := []types.Update{
		{
			SubaccountId: senderSubaccountId,
			AssetUpdates: []types.AssetUpdate{
				{
					AssetId:          assettypes.AssetUsdc.Id,
					BigQuantumsDelta: new(big.Int).Neg(quantums),
				},
			},
		},
		{
			SubaccountId: recipientSubaccountId,
			AssetUpdates: []types.AssetUpdate{
				{
					AssetId:          assettypes.AssetUsdc.Id,
					BigQuantumsDelta: new(big.Int).Set(quantums),
				},
			},
		},
	}
	success, successPerUpdate, err := k.CanUpdateSubaccounts(ctx, updates, types.Transfer)
	if err != nil {
		return err
	}
	if err := types.GetErrorFromUpdateResults(success, successPerUpdate, updates); err != nil {
		return err
	}

	_, coinToTransfer, err := k.assetsKeeper.ConvertAssetToCoin(
		ctx,
		assetId,
		quantums,
	)
	if err != nil {
		return err
	}

	senderCollateralPoolAddr, err := k.GetCollateralPoolForSubaccount(ctx, senderSubaccountId)
	if err != nil {
		return err
	}

	recipientCollateralPoolAddr, err := k.GetCollateralPoolForSubaccount(ctx, recipientSubaccountId)
	if err != nil {
		return err
	}

	// Different collateral pool address, need to do a bank send.
	if !senderCollateralPoolAddr.Equals(recipientCollateralPoolAddr) {
		// Use SendCoins API instead of SendCoinsFromModuleToModule since we don't need the
		// module account feature
		if err := k.bankKeeper.SendCoins(
			ctx,
			senderCollateralPoolAddr,
			recipientCollateralPoolAddr,
			[]sdk.Coin{coinToTransfer},
		); err != nil {
			return err
		}
	}

	// Apply subaccount updates.
	return k.applyValidSubaccountUpdateForTransfer(
		ctx,
		updates,
		types.Transfer,
	)
}
