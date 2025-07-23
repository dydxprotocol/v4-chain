package keeper

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
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

	if k.bankKeeper.BlockedAddr(toAccount) {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", toAccount)
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
	revSharesForFill revsharetypes.RevSharesForFill,
	fill clobtypes.FillForProcess,
) error {
	// get perpetual
	totalFeeQuoteQuantums := new(big.Int).Add(fill.TakerFeeQuoteQuantums, fill.MakerFeeQuoteQuantums)
	perpetual, err := k.perpetualsKeeper.GetPerpetual(ctx, fill.ProductId)
	if err != nil {
		return err
	}

	collateralPoolAddr, err := k.GetCollateralPoolFromPerpetualId(ctx, fill.ProductId)
	if err != nil {
		return err
	}

	// Transfer fees to rev share recipients
	for _, revShare := range revSharesForFill.AllRevShares {
		// transfer fees to the recipient
		recipientAddress, err := sdk.AccAddressFromBech32(revShare.Recipient)
		if err != nil {
			return err
		}

		if err := k.TransferFees(
			ctx,
			assetId,
			collateralPoolAddr,
			recipientAddress,
			revShare.QuoteQuantums,
		); err != nil {
			return err
		}

		// Emit revenue share
		metrics.AddSampleWithLabels(
			metrics.RevenueShareDistribution,
			metrics.GetMetricValueFromBigInt(revShare.QuoteQuantums),
			metrics.GetLabelForStringValue(metrics.RevShareType, revShare.RevShareType.String()),
			metrics.GetLabelForStringValue(metrics.RecipientAddress, revShare.Recipient),
		)

		// Old metric which is being kept for now to ensure data continuity
		if revShare.RevShareType == revsharetypes.REV_SHARE_TYPE_MARKET_MAPPER {
			labels := []metrics.Label{
				metrics.GetLabelForIntValue(metrics.MarketId, int(perpetual.Params.MarketId)),
			}
			metrics.AddSampleWithLabels(
				metrics.MarketMapperRevenueDistribution,
				metrics.GetMetricValueFromBigInt(revShare.QuoteQuantums),
				labels...,
			)
		}
	}

	totalTakerFeeRevShareQuantums := big.NewInt(0)
	if value, ok := revSharesForFill.FeeSourceToQuoteQuantums[revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE]; ok {
		totalTakerFeeRevShareQuantums = value
	}
	totalMakerFeeRevShareQuantums := big.NewInt(0)
	if value, ok := revSharesForFill.FeeSourceToQuoteQuantums[revsharetypes.REV_SHARE_FEE_SOURCE_MAKER_FEE]; ok {
		totalMakerFeeRevShareQuantums = value
	}
	totalNetFeeRevShareQuantums := big.NewInt(0)
	if value, ok :=
		revSharesForFill.FeeSourceToQuoteQuantums[revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE]; ok {
		totalNetFeeRevShareQuantums = value
	}

	totalRevShareQuoteQuantums := big.NewInt(0).Add(
		totalTakerFeeRevShareQuantums,
		totalMakerFeeRevShareQuantums,
	)
	totalRevShareQuoteQuantums = big.NewInt(0).Add(
		totalRevShareQuoteQuantums,
		totalNetFeeRevShareQuantums,
	)

	// Remaining amount goes to the fee collector
	feeCollectorShare := new(big.Int).Sub(
		totalFeeQuoteQuantums,
		totalRevShareQuoteQuantums,
	)

	// If Collector fee share is < 0, panic
	if feeCollectorShare.Sign() < 0 {
		panic("fee collector share is < 0")
	}

	// Emit fee colletor metric
	metrics.AddSample(
		metrics.NetFeesPostRevenueShareDistribution,
		metrics.GetMetricValueFromBigInt(feeCollectorShare),
	)

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

// TransferBuilderFees transfers builder code fees from the collateral pool to the builder address.
// Prior to the transfer, the builder fees are subtracted from the trader's subaccount quote balance.
// This function will panic if the builder fee quantums is negative - which is never expected.
func (k Keeper) TransferBuilderFees(
	ctx sdk.Context,
	productId uint32,
	builderFeeQuantums *big.Int,
	builderAddress string,
) error {
	collateralPoolAddr, err := k.GetCollateralPoolFromPerpetualId(ctx, productId)
	if err != nil {
		return err
	}

	if builderFeeQuantums.Sign() < 0 {
		panic(fmt.Sprintf("builder fee quantums is negative: address: %s, quantums: %s",
			builderAddress,
			builderFeeQuantums.String(),
		))
	}

	_, coinToTransfer, err := k.assetsKeeper.ConvertAssetToCoin(
		ctx,
		assettypes.AssetUsdc.Id,
		builderFeeQuantums,
	)
	if err != nil {
		// Panic if USDC does not exist.
		panic(err)
	}
	recipient, err := sdk.AccAddressFromBech32(builderAddress)
	if err != nil {
		return err
	}

	return k.bankKeeper.SendCoins(
		ctx,
		collateralPoolAddr,
		recipient,
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

// TransferIsolatedInsuranceFundToCross transfers funds from an isolated perpetual's
// insurance fund to the cross-perpetual insurance fund.
// Note: This uses the `x/bank` keeper and modifies `x/bank` state.
func (k Keeper) TransferIsolatedInsuranceFundToCross(ctx sdk.Context, perpetualId uint32) error {
	// Validate perpetual exists
	if _, err := k.perpetualsKeeper.GetPerpetual(ctx, perpetualId); err != nil {
		return err
	}

	isolatedInsuranceFundBalance := k.GetInsuranceFundBalance(ctx, perpetualId)

	// Skip if balance is zero
	if isolatedInsuranceFundBalance.Sign() == 0 {
		return nil
	}

	_, exists := k.assetsKeeper.GetAsset(ctx, assettypes.AssetUsdc.Id)
	if !exists {
		return fmt.Errorf("USDC asset not found in state")
	}

	_, coinToTransfer, err := k.assetsKeeper.ConvertAssetToCoin(
		ctx,
		assettypes.AssetUsdc.Id,
		isolatedInsuranceFundBalance,
	)
	if err != nil {
		return err
	}

	isolatedInsuranceFundAddr, err := k.perpetualsKeeper.GetInsuranceFundModuleAddress(ctx, perpetualId)
	if err != nil {
		return err
	}

	crossInsuranceFundAddr := perptypes.InsuranceFundModuleAddress

	return k.bankKeeper.SendCoins(
		ctx,
		isolatedInsuranceFundAddr,
		crossInsuranceFundAddr,
		[]sdk.Coin{coinToTransfer},
	)
}

// TransferIsolatedCollateralToCross transfers the collateral balance from an isolated perpetual's
// collateral pool to the cross-margin collateral pool. This is used during the upgrade process
// from isolated perpetuals to cross-margin.
// Note: This uses the `x/bank` keeper and modifies `x/bank` state.
func (k Keeper) TransferIsolatedCollateralToCross(ctx sdk.Context, perpetualId uint32) error {
	// Validate perpetual exists
	if _, err := k.perpetualsKeeper.GetPerpetual(ctx, perpetualId); err != nil {
		return err
	}

	isolatedCollateralPoolAddr, err := k.GetCollateralPoolFromPerpetualId(ctx, perpetualId)
	if err != nil {
		return err
	}

	crossCollateralPoolAddr := types.ModuleAddress

	usdcAsset, exists := k.assetsKeeper.GetAsset(ctx, assettypes.AssetUsdc.Id)
	if !exists {
		panic("TransferIsolatedCollateralToCross: Usdc asset not found in state")
	}

	isolatedCollateralPoolBalance := k.bankKeeper.GetBalance(
		ctx,
		isolatedCollateralPoolAddr,
		usdcAsset.Denom,
	)

	// Skip if balance is zero
	if isolatedCollateralPoolBalance.IsZero() {
		return nil
	}

	return k.bankKeeper.SendCoins(
		ctx,
		isolatedCollateralPoolAddr,
		crossCollateralPoolAddr,
		[]sdk.Coin{isolatedCollateralPoolBalance},
	)
}
