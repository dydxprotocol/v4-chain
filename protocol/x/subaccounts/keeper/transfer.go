package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
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
						AssetId:          lib.UsdcAssetId,
						BigQuantumsDelta: bigBalanceDelta,
					},
				},
			},
		}
	} else {
		// TODO(DEC-715): Support non-USDC assets.
		return nil, types.ErrAssetTransferThroughBankNotImplemented
	}

	success, successPerUpdate, err := k.CanUpdateSubaccounts(ctx, updates)
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
) error {
	// Update subaccount to reflect the transfer.
	success, successPerUpdate, err := k.UpdateSubaccounts(ctx, updates)

	// Neither of the two conditions below should be true, since `k.CanUpdateSubaccount()`
	// already succeeded.
	if err != nil {
		return err
	}

	return types.GetErrorFromUpdateResults(success, successPerUpdate, updates)
}

// TransferFundsFromSubaccountToModule returns an error if the call to `k.CanUpdateSubaccounts()`
// fails. Otherwise, deducts the asset quantums from the subaccount, translates the
// `assetId` and `quantums` into a `sdk.Coin`, and calls
// `bankKeeper.SendCoinsFromModuleToModule()`.
func (k Keeper) TransferFundsFromSubaccountToModule(
	ctx sdk.Context,
	fromSubaccountId types.SubaccountId,
	toModule string,
	assetId uint32,
	quantums *big.Int,
) error {
	// TODO(DEC-715): Support non-USDC assets.
	if assetId != lib.UsdcAssetId {
		return types.ErrAssetTransferThroughBankNotImplemented
	}

	if quantums.Sign() <= 0 {
		return errorsmod.Wrap(types.ErrAssetTransferQuantumsNotPositive, lib.Uint32ToString(assetId))
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

	// Send coins from the `subaccounts` module account to `toModule`.
	if err := k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		types.ModuleName, // senderModule
		toModule,         // recipientModule
		[]sdk.Coin{coinToTransfer},
	); err != nil {
		return err
	}

	// Apply subaccount updates.
	if err := k.applyValidSubaccountUpdateForTransfer(
		ctx,
		updates,
	); err != nil {
		panic(fmt.Sprintf("SendCoinsFromModuleToModule() succeeded but UpdateSubaccounts() failed: %v", err))
	}

	return nil
}

// TransferFundsFromModuleToSubaccount returns an error if the call to `k.CanUpdateSubaccounts()`
// fails. Otherwise, increases the asset quantums in the subaccount, translates the
// `assetId` and `quantums` into a `sdk.Coin`, and calls `bankKeeper.SendCoinsFromModuleToModule()`.
func (k Keeper) TransferFundsFromModuleToSubaccount(
	ctx sdk.Context,
	fromModule string,
	toSubaccountId types.SubaccountId,
	assetId uint32,
	quantums *big.Int,
) error {
	// TODO(DEC-715): Support non-USDC assets.
	if assetId != lib.UsdcAssetId {
		return types.ErrAssetTransferThroughBankNotImplemented
	}

	if quantums.Sign() <= 0 {
		return errorsmod.Wrap(types.ErrAssetTransferQuantumsNotPositive, lib.Uint32ToString(assetId))
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

	// Send coins from `fromModule` to the `subaccounts` module account.
	if err := k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		fromModule,       // senderModule
		types.ModuleName, // recipientModule
		[]sdk.Coin{coinToTransfer},
	); err != nil {
		return err
	}

	// Apply subaccount updates.
	return k.applyValidSubaccountUpdateForTransfer(
		ctx,
		updates,
	)
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
	if assetId != lib.UsdcAssetId {
		return types.ErrAssetTransferThroughBankNotImplemented
	}

	if quantums.Sign() <= 0 {
		return errorsmod.Wrap(types.ErrAssetTransferQuantumsNotPositive, lib.Uint32ToString(assetId))
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

	// Send coins from `fromModule` to the `subaccounts` module account.
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		fromAccount,
		types.ModuleName,
		[]sdk.Coin{coinToTransfer},
	); err != nil {
		return err
	}

	// Apply subaccount updates.
	return k.applyValidSubaccountUpdateForTransfer(
		ctx,
		updates,
	)
}

// WithdrawFundsFromSubaccountToAccount returns an error if the call to `k.CanUpdateSubaccounts()`
// fails. Otherwise, deducts the asset quantums from the subaccount, translates the
// `assetId` and `quantums` into a `sdk.Coin`, and calls `bankKeeper.SendCoinsFromModuleToAccount()`.
func (k Keeper) WithdrawFundsFromSubaccountToAccount(
	ctx sdk.Context,
	fromSubaccountId types.SubaccountId,
	toAccount sdk.AccAddress,
	assetId uint32,
	quantums *big.Int,
) error {
	// TODO(DEC-715): Support non-USDC assets.
	if assetId != lib.UsdcAssetId {
		return types.ErrAssetTransferThroughBankNotImplemented
	}

	if quantums.Sign() <= 0 {
		return errorsmod.Wrap(types.ErrAssetTransferQuantumsNotPositive, lib.Uint32ToString(assetId))
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

	// Send coins from `fromModule` to the `subaccounts` module account.
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		toAccount,
		[]sdk.Coin{coinToTransfer},
	); err != nil {
		return err
	}

	// Apply subaccount updates.
	return k.applyValidSubaccountUpdateForTransfer(
		ctx,
		updates,
	)
}

// TransferFeesToFeeCollectorModule translates the assetId and quantums into a sdk.Coin,
// and moves the funds from subaccounts module to the `fee_collector` module account by calling
// bankKeeper.SendCoinsFromModuleToModule(). Does not change any individual subaccount state.
func (k Keeper) TransferFeesToFeeCollectorModule(
	ctx sdk.Context,
	assetId uint32,
	quantums *big.Int,
) error {
	// TODO(DEC-715): Support non-USDC assets.
	if assetId != lib.UsdcAssetId {
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

	// Send coins from `subaccounts` to the `auth` module fee collector account.
	fromModule := types.ModuleName
	toModule := authtypes.FeeCollectorName

	if quantums.Sign() < 0 {
		// In the case of a liquidation, net fees can be negative if the maker gets a rebate.
		fromModule, toModule = toModule, fromModule
	}

	if err := k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		fromModule,
		toModule,
		[]sdk.Coin{coinToTransfer},
	); err != nil {
		return err
	}

	return nil
}

// TransferInsuranceFundPayments transfers funds in and out of the insurance fund to the subaccounts
// module by calling `bankKeeper.SendCoinsFromModuleToModule`.
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
) error {
	if insuranceFundDelta.Sign() == 0 {
		return nil
	}

	_, coinToTransfer, err := k.assetsKeeper.ConvertAssetToCoin(
		ctx,
		lib.UsdcAssetId,
		new(big.Int).Abs(insuranceFundDelta),
	)
	if err != nil {
		// Panic if USDC does not exist.
		panic(err)
	}

	// Determine the sender and receiver.
	// Send coins from `subaccounts` to the `insurance_fund` module account by default.
	fromModule := types.ModuleName
	toModule := clobtypes.InsuranceFundName

	if insuranceFundDelta.Sign() < 0 {
		// Insurance fund needs to cover losses from liquidations.
		// Send coins from the insurance fund to the `subaccounts` module account.
		fromModule, toModule = toModule, fromModule
	}

	return k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		fromModule,
		toModule,
		[]sdk.Coin{coinToTransfer},
	)
}
