package keeper_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestRefreshAllVaultOrders(t *testing.T) {
	tests := map[string]struct {
		// Vault IDs.
		vaultIds []vaulttypes.VaultId
		// Status of each vault above.
		vaultStatuses []vaulttypes.VaultStatus
		// Asset quantums of each vault ID above.
		assetQuantums []*big.Int
		// Position quantums of each vault ID above.
		positionQuantums []*big.Int
		// Activation threshold (quote quantums) of vaults.
		activationThresholdQuoteQuantums *big.Int
	}{
		"Two Vaults, Both Quoting, Both Above Activation Threshold": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob0,
				constants.Vault_Clob1,
			},
			vaultStatuses: []vaulttypes.VaultStatus{
				vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
				vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			},
			assetQuantums: []*big.Int{
				big.NewInt(1_000_000_000), // 1,000 USDC
				big.NewInt(1_000_000_001),
			},
			positionQuantums: []*big.Int{
				big.NewInt(0),
				big.NewInt(0),
			},
			activationThresholdQuoteQuantums: big.NewInt(1_000_000_000),
		},
		"Two Vaults, One Quoting, One Stand-By, Both Above Activation Threshold": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob0,
				constants.Vault_Clob1,
			},
			vaultStatuses: []vaulttypes.VaultStatus{
				vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
				vaulttypes.VaultStatus_VAULT_STATUS_STAND_BY,
			},
			assetQuantums: []*big.Int{
				big.NewInt(1_000_000_000), // 1,000 USDC
				big.NewInt(1_000_000_001),
			},
			positionQuantums: []*big.Int{
				big.NewInt(0),
				big.NewInt(0),
			},
			activationThresholdQuoteQuantums: big.NewInt(1_000_000_000),
		},
		"Two Vaults, One Stand-By, One Deactivated, One Above Activation Threshold": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob0,
				constants.Vault_Clob1,
			},
			vaultStatuses: []vaulttypes.VaultStatus{
				vaulttypes.VaultStatus_VAULT_STATUS_STAND_BY,
				vaulttypes.VaultStatus_VAULT_STATUS_DEACTIVATED,
			},
			assetQuantums: []*big.Int{
				big.NewInt(1_000_000_000), // 1,000 USDC
				big.NewInt(0),
			},
			positionQuantums: []*big.Int{
				big.NewInt(0),
				big.NewInt(0),
			},
			activationThresholdQuoteQuantums: big.NewInt(1_000_000_000),
		},
		"Two Vaults, Both Quoting, Only One above Activation Threshold": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob0,
				constants.Vault_Clob1,
			},
			vaultStatuses: []vaulttypes.VaultStatus{
				vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
				vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			},
			assetQuantums: []*big.Int{
				big.NewInt(1_000_000_000),
				big.NewInt(999_999_999),
			},
			positionQuantums: []*big.Int{
				big.NewInt(0),
				big.NewInt(0),
			},
			activationThresholdQuoteQuantums: big.NewInt(1_000_000_000),
		},
		"Two Vaults, Both Quoting, Both below Activation Threshold": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob0,
				constants.Vault_Clob1,
			},
			vaultStatuses: []vaulttypes.VaultStatus{
				vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
				vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			},
			assetQuantums: []*big.Int{
				big.NewInt(123_456_788),
				big.NewInt(123_456_787),
			},
			positionQuantums: []*big.Int{
				big.NewInt(0),
				big.NewInt(0),
			},
			activationThresholdQuoteQuantums: big.NewInt(123_456_789),
		},
		"Four Vaults, One stand-by, One close-only and short, One close-only and neutral, " +
			"One close-only and long, All above Activation Threshold": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob0,
				constants.Vault_Clob1,
				{
					Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
					Number: 2,
				},
				{
					Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
					Number: 5,
				},
			},
			vaultStatuses: []vaulttypes.VaultStatus{
				vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
				vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
				vaulttypes.VaultStatus_VAULT_STATUS_STAND_BY,
				vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
			},
			assetQuantums: []*big.Int{
				big.NewInt(123_456_789),
				big.NewInt(123_456_790),
				big.NewInt(123_456_791),
				big.NewInt(123_456_792),
			},
			positionQuantums: []*big.Int{
				big.NewInt(-100_000),
				big.NewInt(0),
				big.NewInt(-100_000),
				big.NewInt(100_000),
			},
			activationThresholdQuoteQuantums: big.NewInt(123_456_789),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Enable testapp's indexer event manager
			msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
			appOpts := map[string]interface{}{
				indexer.MsgSenderInstanceForTest: msgSender,
			}

			// Initialize tApp.
			tApp := testapp.NewTestAppBuilder(t).WithAppOptions(appOpts).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Initialize each vault with enough quote quantums to be actively quoting.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						subaccounts := make([]satypes.Subaccount, len(tc.vaultIds))
						for i, vaultId := range tc.vaultIds {
							perpetualPositions := []*satypes.PerpetualPosition{}
							if tc.positionQuantums[i].Sign() != 0 {
								perpetualPositions = append(
									perpetualPositions,
									&satypes.PerpetualPosition{
										PerpetualId: vaultId.Number,
										Quantums:    dtypes.NewIntFromBigInt(tc.positionQuantums[i]),
									},
								)
							}
							subaccounts[i] = satypes.Subaccount{
								Id: vaultId.ToSubaccountId(),
								AssetPositions: []*satypes.AssetPosition{
									testutil.CreateSingleAssetPosition(
										assettypes.AssetUsdc.Id,
										tc.assetQuantums[i],
									),
								},
								PerpetualPositions: perpetualPositions,
							}
						}
						genesisState.Subaccounts = subaccounts
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vaulttypes.GenesisState) {
						defaultQuotingParams := genesisState.DefaultQuotingParams
						defaultQuotingParams.ActivationThresholdQuoteQuantums = dtypes.NewIntFromBigInt(
							tc.activationThresholdQuoteQuantums,
						)
						genesisState.DefaultQuotingParams = defaultQuotingParams
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain().WithIsCheckTx(false)

			// Set vault params of each vault.
			for i, vaultId := range tc.vaultIds {
				err := tApp.App.VaultKeeper.SetVaultParams(
					ctx,
					vaultId,
					vaulttypes.VaultParams{
						Status: tc.vaultStatuses[i],
					},
				)
				require.NoError(t, err)
			}
			// Clear events from setting vault params
			tApp.App.IndexerEventManager.ClearEvents(ctx)

			// Check that there's no stateful orders yet.
			allStatefulOrders := tApp.App.ClobKeeper.GetAllStatefulOrders(ctx)
			require.Len(t, allStatefulOrders, 0)

			// Refresh all vault orders.
			tApp.App.VaultKeeper.RefreshAllVaultOrders(ctx)

			// Check expected orders and order placement indexer events.
			allExpectedOrders := []clobtypes.Order{}
			expectedIndexerEvents := []*indexer_manager.IndexerTendermintEvent{}
			indexerEventIndex := 0
			for i, vaultId := range tc.vaultIds {
				if tc.vaultStatuses[i] == vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY &&
					tc.positionQuantums[i].Sign() != 0 ||
					tc.vaultStatuses[i] == vaulttypes.VaultStatus_VAULT_STATUS_QUOTING &&
						tc.assetQuantums[i].Cmp(tc.activationThresholdQuoteQuantums) >= 0 {
					expectedOrders, err := tApp.App.VaultKeeper.GetVaultClobOrders(ctx, vaultId)
					require.NoError(t, err)

					for _, order := range expectedOrders {
						allExpectedOrders = append(allExpectedOrders, *order)
						event := indexer_manager.IndexerTendermintEvent{
							Subtype: indexerevents.SubtypeStatefulOrder,
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{
								TransactionIndex: 0,
							},
							EventIndex: uint32(indexerEventIndex),
							Version:    indexerevents.StatefulOrderEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewLongTermOrderPlacementEvent(
									*order,
								),
							),
						}
						indexerEventIndex += 1
						expectedIndexerEvents = append(expectedIndexerEvents, &event)
					}
				}
			}
			allStatefulOrders = tApp.App.ClobKeeper.GetAllStatefulOrders(ctx)
			require.ElementsMatch(t, allExpectedOrders, allStatefulOrders)
			block := tApp.App.IndexerEventManager.ProduceBlock(ctx)
			require.ElementsMatch(t, expectedIndexerEvents, block.Events)
		})
	}
}

func TestRefreshVaultClobOrders(t *testing.T) {
	type Instance struct {
		advanceBlock func(ctx sdk.Context, tApp *testapp.TestApp) sdk.Context

		/* --- Expectations --- */
		// Whether each order refreshes.
		ordersRefreshed []bool
		// Sides of orders after `advanceBlock` is called.
		orderSides []clobtypes.Order_Side
		// Whether client IDs are canonical or flipped after `advanceBlock` is called.
		clientIdIsCanonical []bool
	}

	tests := map[string]struct {
		// Vault ID.
		vaultId                   vaulttypes.VaultId
		vaultPositionBaseQuantums *big.Int

		instances []Instance
	}{
		"Success - Orders do not refresh": {
			vaultId: constants.Vault_Clob0,
			instances: []Instance{
				{
					advanceBlock: func(ctx sdk.Context, tApp *testapp.TestApp) sdk.Context {
						return tApp.AdvanceToBlock(
							uint32(tApp.GetBlockHeight())+1,
							testapp.AdvanceToBlockOptions{
								BlockTime: ctx.BlockTime().Add(time.Second),
							},
						)
					},
					ordersRefreshed: []bool{
						false,
						false,
						false,
						false,
					},
					orderSides: []clobtypes.Order_Side{
						clobtypes.Order_SIDE_SELL,
						clobtypes.Order_SIDE_BUY,
						clobtypes.Order_SIDE_SELL,
						clobtypes.Order_SIDE_BUY,
					},
					clientIdIsCanonical: []bool{
						true,
						true,
						true,
						true,
					},
				},
			},
		},
		"Success - Orders refresh due to expiration": {
			vaultId: constants.Vault_Clob0,
			instances: []Instance{
				{
					advanceBlock: func(ctx sdk.Context, tApp *testapp.TestApp) sdk.Context {
						orderExpirationSeconds := vaulttypes.DefaultQuotingParams().OrderExpirationSeconds
						return tApp.AdvanceToBlock(
							uint32(tApp.GetBlockHeight())+5,
							testapp.AdvanceToBlockOptions{
								BlockTime: ctx.BlockTime().Add(
									time.Second * time.Duration(orderExpirationSeconds),
								),
							},
						)
					},
					ordersRefreshed: []bool{
						true,
						true,
						true,
						true,
					},
					orderSides: []clobtypes.Order_Side{
						clobtypes.Order_SIDE_SELL,
						clobtypes.Order_SIDE_BUY,
						clobtypes.Order_SIDE_SELL,
						clobtypes.Order_SIDE_BUY,
					},
					clientIdIsCanonical: []bool{
						false,
						false,
						false,
						false,
					},
				},
			},
		},
		"Success - Orders refresh due to price updates": {
			vaultId: constants.Vault_Clob0,
			instances: []Instance{
				{
					advanceBlock: func(ctx sdk.Context, tApp *testapp.TestApp) sdk.Context {
						marketPrice, err := tApp.App.PricesKeeper.GetMarketPrice(ctx, constants.Vault_Clob0.Number)
						require.NoError(t, err)
						msgUpdateMarketPrices := &pricestypes.MsgUpdateMarketPrices{
							MarketPriceUpdates: []*pricestypes.MsgUpdateMarketPrices_MarketPrice{
								{
									MarketId: constants.Vault_Clob0.Number,
									Price:    marketPrice.Price * 2,
								},
							},
						}
						return tApp.AdvanceToBlock(
							uint32(tApp.GetBlockHeight())+1,
							testapp.AdvanceToBlockOptions{
								BlockTime: ctx.BlockTime().Add(time.Second),
								DeliverTxsOverride: [][]byte{
									testtx.MustGetTxBytes(msgUpdateMarketPrices),
								},
							},
						)
					},
					ordersRefreshed: []bool{
						true,
						true,
						true,
						true,
					},
					orderSides: []clobtypes.Order_Side{
						clobtypes.Order_SIDE_SELL,
						clobtypes.Order_SIDE_BUY,
						clobtypes.Order_SIDE_SELL,
						clobtypes.Order_SIDE_BUY,
					},
					clientIdIsCanonical: []bool{
						false,
						false,
						false,
						false,
					},
				},
				{
					advanceBlock: func(ctx sdk.Context, tApp *testapp.TestApp) sdk.Context {
						marketPrice, err := tApp.App.PricesKeeper.GetMarketPrice(ctx, constants.Vault_Clob0.Number)
						require.NoError(t, err)
						msgUpdateMarketPrices := &pricestypes.MsgUpdateMarketPrices{
							MarketPriceUpdates: []*pricestypes.MsgUpdateMarketPrices_MarketPrice{
								{
									MarketId: constants.Vault_Clob0.Number,
									Price:    marketPrice.Price / 2,
								},
							},
						}
						return tApp.AdvanceToBlock(
							uint32(tApp.GetBlockHeight())+1,
							testapp.AdvanceToBlockOptions{
								BlockTime: ctx.BlockTime().Add(time.Second),
								DeliverTxsOverride: [][]byte{
									testtx.MustGetTxBytes(msgUpdateMarketPrices),
								},
							},
						)
					},
					ordersRefreshed: []bool{
						true,
						true,
						true,
						true,
					},
					orderSides: []clobtypes.Order_Side{
						clobtypes.Order_SIDE_SELL,
						clobtypes.Order_SIDE_BUY,
						clobtypes.Order_SIDE_SELL,
						clobtypes.Order_SIDE_BUY,
					},
					clientIdIsCanonical: []bool{
						true,
						true,
						true,
						true,
					},
				},
			},
		},
		"Success - Orders refresh due to order size increase": {
			vaultId: constants.Vault_Clob0,
			instances: []Instance{
				{
					advanceBlock: func(ctx sdk.Context, tApp *testapp.TestApp) sdk.Context {
						msgAllocateToVault := vaulttypes.MsgAllocateToVault{
							Authority:     constants.AliceAccAddress.String(), // operator
							VaultId:       constants.Vault_Clob0,
							QuoteQuantums: dtypes.NewInt(500_000_000),
						}
						CheckTx_MsgAllocateToVault := testapp.MustMakeCheckTx(
							ctx,
							tApp.App,
							testapp.MustMakeCheckTxOptions{
								AccAddressForSigning: constants.AliceAccAddress.String(),
								Gas:                  constants.TestGasLimit,
								FeeAmt:               constants.TestFeeCoins_5Cents,
							},
							&msgAllocateToVault,
						)
						checkTxResp := tApp.CheckTx(CheckTx_MsgAllocateToVault)
						require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

						return tApp.AdvanceToBlock(
							uint32(tApp.GetBlockHeight())+1,
							testapp.AdvanceToBlockOptions{
								BlockTime: ctx.BlockTime().Add(time.Second * 2),
							},
						)
					},
					ordersRefreshed: []bool{
						true,
						true,
						true,
						true,
					},
					orderSides: []clobtypes.Order_Side{
						clobtypes.Order_SIDE_SELL,
						clobtypes.Order_SIDE_BUY,
						clobtypes.Order_SIDE_SELL,
						clobtypes.Order_SIDE_BUY,
					},
					clientIdIsCanonical: []bool{
						false,
						false,
						false,
						false,
					},
				},
			},
		},
		"Success - Orders refresh due to order size decrease": {
			vaultId: constants.Vault_Clob0,
			instances: []Instance{
				{
					advanceBlock: func(ctx sdk.Context, tApp *testapp.TestApp) sdk.Context {
						msgRetrieveFromVault := vaulttypes.MsgRetrieveFromVault{
							Authority:     constants.AliceAccAddress.String(), // operator
							VaultId:       constants.Vault_Clob0,
							QuoteQuantums: dtypes.NewInt(500_000_000),
						}
						CheckTx_MsgRetrieveFromVault := testapp.MustMakeCheckTx(
							ctx,
							tApp.App,
							testapp.MustMakeCheckTxOptions{
								AccAddressForSigning: constants.AliceAccAddress.String(),
								Gas:                  constants.TestGasLimit,
								FeeAmt:               constants.TestFeeCoins_5Cents,
							},
							&msgRetrieveFromVault,
						)
						checkTxResp := tApp.CheckTx(CheckTx_MsgRetrieveFromVault)
						require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

						return tApp.AdvanceToBlock(
							uint32(tApp.GetBlockHeight())+1,
							testapp.AdvanceToBlockOptions{
								BlockTime: ctx.BlockTime().Add(time.Second * 2),
							},
						)
					},
					ordersRefreshed: []bool{
						true,
						true,
						true,
						true,
					},
					orderSides: []clobtypes.Order_Side{
						clobtypes.Order_SIDE_SELL,
						clobtypes.Order_SIDE_BUY,
						clobtypes.Order_SIDE_SELL,
						clobtypes.Order_SIDE_BUY,
					},
					clientIdIsCanonical: []bool{
						false,
						false,
						false,
						false,
					},
				},
			},
		},
		"Success - Orders refresh due to status changing to close-only and zero leverage. No more orders": {
			vaultId: constants.Vault_Clob0,
			instances: []Instance{
				{
					advanceBlock: func(ctx sdk.Context, tApp *testapp.TestApp) sdk.Context {
						msgSetVaultParams := vaulttypes.MsgSetVaultParams{
							Authority: constants.AliceAccAddress.String(), // operator
							VaultId:   constants.Vault_Clob0,
							VaultParams: vaulttypes.VaultParams{
								Status: vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
							},
						}
						CheckTx_MsgSetVaultParams := testapp.MustMakeCheckTx(
							ctx,
							tApp.App,
							testapp.MustMakeCheckTxOptions{
								AccAddressForSigning: constants.AliceAccAddress.String(),
								Gas:                  constants.TestGasLimit,
								FeeAmt:               constants.TestFeeCoins_5Cents,
							},
							&msgSetVaultParams,
						)
						checkTxResp := tApp.CheckTx(CheckTx_MsgSetVaultParams)
						require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

						return tApp.AdvanceToBlock(
							uint32(tApp.GetBlockHeight())+1,
							testapp.AdvanceToBlockOptions{
								BlockTime: ctx.BlockTime().Add(time.Second * 2),
							},
						)
					},
					ordersRefreshed:     []bool{}, // no orders
					orderSides:          []clobtypes.Order_Side{},
					clientIdIsCanonical: []bool{},
				},
			},
		},
		"Success - Orders refresh due to status changing to close-only. Sell orders only": {
			vaultId:                   constants.Vault_Clob1,
			vaultPositionBaseQuantums: big.NewInt(250_000_000),
			instances: []Instance{
				{
					advanceBlock: func(ctx sdk.Context, tApp *testapp.TestApp) sdk.Context {
						msgSetVaultParams := vaulttypes.MsgSetVaultParams{
							Authority: constants.AliceAccAddress.String(), // operator
							VaultId:   constants.Vault_Clob1,
							VaultParams: vaulttypes.VaultParams{
								Status: vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
							},
						}
						CheckTx_MsgSetVaultParams := testapp.MustMakeCheckTx(
							ctx,
							tApp.App,
							testapp.MustMakeCheckTxOptions{
								AccAddressForSigning: constants.AliceAccAddress.String(),
								Gas:                  constants.TestGasLimit,
								FeeAmt:               constants.TestFeeCoins_5Cents,
							},
							&msgSetVaultParams,
						)
						checkTxResp := tApp.CheckTx(CheckTx_MsgSetVaultParams)
						require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

						return tApp.AdvanceToBlock(
							uint32(tApp.GetBlockHeight())+1,
							testapp.AdvanceToBlockOptions{
								BlockTime: ctx.BlockTime().Add(time.Second * 2),
							},
						)
					},
					ordersRefreshed: []bool{
						false, // this sell order doesn't refresh as it's already desired.
						true,
					},
					orderSides: []clobtypes.Order_Side{
						clobtypes.Order_SIDE_SELL,
						clobtypes.Order_SIDE_SELL,
					}, // Two sell orders only as vault is long.
					clientIdIsCanonical: []bool{
						true,
						false, // flipped as it's refreshed.
					},
				},
			},
		},
		"Success - Orders refresh due to status changing to close-only and back to quoting": {
			vaultId:                   constants.Vault_Clob1,
			vaultPositionBaseQuantums: big.NewInt(-250_000_000),
			instances: []Instance{
				{
					advanceBlock: func(ctx sdk.Context, tApp *testapp.TestApp) sdk.Context {
						msgSetVaultParams := vaulttypes.MsgSetVaultParams{
							Authority: constants.AliceAccAddress.String(), // operator
							VaultId:   constants.Vault_Clob1,
							VaultParams: vaulttypes.VaultParams{
								Status: vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
							},
						}
						CheckTx_MsgSetVaultParams := testapp.MustMakeCheckTx(
							ctx,
							tApp.App,
							testapp.MustMakeCheckTxOptions{
								AccAddressForSigning: constants.AliceAccAddress.String(),
								Gas:                  constants.TestGasLimit,
								FeeAmt:               constants.TestFeeCoins_5Cents,
							},
							&msgSetVaultParams,
						)
						checkTxResp := tApp.CheckTx(CheckTx_MsgSetVaultParams)
						require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

						return tApp.AdvanceToBlock(
							uint32(tApp.GetBlockHeight())+1,
							testapp.AdvanceToBlockOptions{
								BlockTime: ctx.BlockTime().Add(time.Second * 2),
							},
						)
					},
					ordersRefreshed: []bool{
						true,
						true,
					},
					orderSides: []clobtypes.Order_Side{
						clobtypes.Order_SIDE_BUY,
						clobtypes.Order_SIDE_BUY,
					}, // Two buy orders only as vault is short.
					clientIdIsCanonical: []bool{
						false,
						false,
					},
				},
				{
					advanceBlock: func(ctx sdk.Context, tApp *testapp.TestApp) sdk.Context {
						msgSetVaultParams := vaulttypes.MsgSetVaultParams{
							Authority: constants.AliceAccAddress.String(), // operator
							VaultId:   constants.Vault_Clob1,
							VaultParams: vaulttypes.VaultParams{
								Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
							},
						}
						CheckTx_MsgSetVaultParams := testapp.MustMakeCheckTx(
							ctx,
							tApp.App,
							testapp.MustMakeCheckTxOptions{
								AccAddressForSigning: constants.AliceAccAddress.String(),
								Gas:                  constants.TestGasLimit,
								FeeAmt:               constants.TestFeeCoins_5Cents,
							},
							&msgSetVaultParams,
						)
						checkTxResp := tApp.CheckTx(CheckTx_MsgSetVaultParams)
						require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

						return tApp.AdvanceToBlock(
							uint32(tApp.GetBlockHeight())+1,
							testapp.AdvanceToBlockOptions{
								BlockTime: ctx.BlockTime().Add(time.Second * 2),
							},
						)
					},
					ordersRefreshed: []bool{
						true,
						true,
						true,
						true,
					},
					orderSides: []clobtypes.Order_Side{
						clobtypes.Order_SIDE_SELL,
						clobtypes.Order_SIDE_BUY,
						clobtypes.Order_SIDE_SELL,
						clobtypes.Order_SIDE_BUY,
					},
					clientIdIsCanonical: []bool{
						true,
						true,
						true,
						true,
					},
				},
			},
		},
		"Success - Orders refresh due to status changing to stand-by. No more orders": {
			vaultId: constants.Vault_Clob0,
			instances: []Instance{
				{
					advanceBlock: func(ctx sdk.Context, tApp *testapp.TestApp) sdk.Context {
						msgSetVaultParams := vaulttypes.MsgSetVaultParams{
							Authority: constants.AliceAccAddress.String(), // operator
							VaultId:   constants.Vault_Clob0,
							VaultParams: vaulttypes.VaultParams{
								Status: vaulttypes.VaultStatus_VAULT_STATUS_STAND_BY,
							},
						}
						CheckTx_MsgSetVaultParams := testapp.MustMakeCheckTx(
							ctx,
							tApp.App,
							testapp.MustMakeCheckTxOptions{
								AccAddressForSigning: constants.AliceAccAddress.String(),
								Gas:                  constants.TestGasLimit,
								FeeAmt:               constants.TestFeeCoins_5Cents,
							},
							&msgSetVaultParams,
						)
						checkTxResp := tApp.CheckTx(CheckTx_MsgSetVaultParams)
						require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

						return tApp.AdvanceToBlock(
							uint32(tApp.GetBlockHeight())+1,
							testapp.AdvanceToBlockOptions{
								BlockTime: ctx.BlockTime().Add(time.Second * 2),
							},
						)
					},
					ordersRefreshed:     []bool{}, // no orders
					orderSides:          []clobtypes.Order_Side{},
					clientIdIsCanonical: []bool{},
				},
			},
		},
		"Success - Orders refresh due to retrieving to 0 equity and status changing to deactivated. No more orders": {
			vaultId: constants.Vault_Clob0,
			instances: []Instance{
				{
					advanceBlock: func(ctx sdk.Context, tApp *testapp.TestApp) sdk.Context {
						msgRetrieveFromVault := vaulttypes.MsgRetrieveFromVault{
							Authority:     constants.AliceAccAddress.String(), // operator
							VaultId:       constants.Vault_Clob0,
							QuoteQuantums: dtypes.NewInt(2_000_000_000), // retrieve all quote quantums the vault has.
						}
						CheckTx_MsgRetrieveFromVault := testapp.MustMakeCheckTx(
							ctx,
							tApp.App,
							testapp.MustMakeCheckTxOptions{
								AccAddressForSigning: constants.AliceAccAddress.String(),
								Gas:                  constants.TestGasLimit,
								FeeAmt:               constants.TestFeeCoins_5Cents,
							},
							&msgRetrieveFromVault,
						)
						checkTxResp := tApp.CheckTx(CheckTx_MsgRetrieveFromVault)
						require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

						msgSetVaultParams := vaulttypes.MsgSetVaultParams{
							Authority: constants.AliceAccAddress.String(), // operator
							VaultId:   constants.Vault_Clob0,
							VaultParams: vaulttypes.VaultParams{
								Status: vaulttypes.VaultStatus_VAULT_STATUS_DEACTIVATED,
							},
						}
						CheckTx_MsgSetVaultParams := testapp.MustMakeCheckTx(
							ctx,
							tApp.App,
							testapp.MustMakeCheckTxOptions{
								AccAddressForSigning: constants.AliceAccAddress.String(),
								Gas:                  constants.TestGasLimit,
								FeeAmt:               constants.TestFeeCoins_5Cents,
							},
							&msgSetVaultParams,
						)
						checkTxResp = tApp.CheckTx(CheckTx_MsgSetVaultParams)
						require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

						return tApp.AdvanceToBlock(
							uint32(tApp.GetBlockHeight())+1,
							testapp.AdvanceToBlockOptions{
								BlockTime: ctx.BlockTime().Add(time.Second * 2),
							},
						)
					},
					ordersRefreshed:     []bool{}, // no orders
					orderSides:          []clobtypes.Order_Side{},
					clientIdIsCanonical: []bool{},
				},
			},
		},
		"Success - Vault for non-existent Clob Pair 4321": {
			vaultId: vaulttypes.VaultId{
				Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
				Number: 4321,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize tApp.
			defaultQuotingParams := vaulttypes.DefaultQuotingParams()
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vaulttypes.GenesisState) {
						genesisState.Vaults = []vaulttypes.Vault{
							{
								VaultId: tc.vaultId,
								VaultParams: vaulttypes.VaultParams{
									Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
								},
							},
						}
						genesisState.OperatorParams = vaulttypes.OperatorParams{
							Operator: constants.AliceAccAddress.String(),
						}
					},
				)
				// Initialize main vault with some balance and vault with enough quote quantums to be actively quoting.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							{
								Id: &vaulttypes.MegavaultMainSubaccount,
								AssetPositions: []*satypes.AssetPosition{
									testutil.CreateSingleAssetPosition(
										assettypes.AssetUsdc.Id,
										big.NewInt(1_000_000_000),
									),
								},
							},
							{
								Id: tc.vaultId.ToSubaccountId(),
								AssetPositions: []*satypes.AssetPosition{
									testutil.CreateSingleAssetPosition(
										assettypes.AssetUsdc.Id,
										big.NewInt(2_000_000_000),
									),
								},
								PerpetualPositions: []*satypes.PerpetualPosition{
									{
										PerpetualId: tc.vaultId.Number,
										Quantums:    dtypes.NewIntFromBigInt(tc.vaultPositionBaseQuantums),
									},
								},
							},
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			// Get canonical and flipped client IDs of this vault's orders.
			orderIds := tApp.App.VaultKeeper.GetVaultClobOrderIds(ctx, tc.vaultId)
			canonicalClientIds := make([]uint32, len(orderIds))
			flippedClientIds := make([]uint32, len(orderIds))
			for i, orderId := range orderIds {
				canonicalClientIds[i] = orderId.ClientId
				flippedClientIds[i] = orderId.ClientId ^ 1
			}

			// Helper function that verifies that vault orders are as expected.
			verifyVaultOrders := func(
				ordersRefreshed []bool,
				oldGTBT uint32,
				newGTBT uint32,
				clientIdIsCanonical []bool,
				expectedSides []clobtypes.Order_Side,
			) {
				allStatefulOrders := tApp.App.ClobKeeper.GetAllStatefulOrders(ctx)
				// Verify the number of vault orders.
				require.Len(t, allStatefulOrders, len(ordersRefreshed))
				// Verify that GTBT of orders is as expected.
				for i, order := range allStatefulOrders {
					if ordersRefreshed[i] {
						require.Equal(t, newGTBT, order.GetGoodTilBlockTime())
					} else {
						require.Equal(t, oldGTBT, order.GetGoodTilBlockTime())
					}
				}

				// Verify that
				// 1. stateful orders have expected client IDs.
				// 2. most recent client IDs are as expected.
				mostRecentClientIds := tApp.App.VaultKeeper.GetMostRecentClientIds(ctx, tc.vaultId)
				require.Len(t, mostRecentClientIds, len(ordersRefreshed))
				for i, order := range allStatefulOrders {
					if clientIdIsCanonical[i] {
						require.Equal(t, canonicalClientIds[i], order.OrderId.ClientId)
						require.Equal(t, canonicalClientIds[i], mostRecentClientIds[i])
					} else {
						require.Equal(t, flippedClientIds[i], order.OrderId.ClientId)
						require.Equal(t, flippedClientIds[i], mostRecentClientIds[i])
					}
				}

				// Verify order sides.
				for i, order := range allStatefulOrders {
					require.Equal(t, expectedSides[i], order.Side)
				}
			}

			// If corresponding clob pair doesn't exist, the vault should not place any orders.
			_, found := tApp.App.ClobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(tc.vaultId.Number))
			if !found {
				require.Zero(t, len(tApp.App.ClobKeeper.GetAllStatefulOrders(ctx)))
				return
			}
			// Vault should place its initial orders (client IDs should be canonical).
			initialOrderSides := make([]clobtypes.Order_Side, defaultQuotingParams.Layers*2)
			for i := uint32(0); i < defaultQuotingParams.Layers; i++ {
				initialOrderSides[2*i] = clobtypes.Order_SIDE_SELL
				initialOrderSides[2*i+1] = clobtypes.Order_SIDE_BUY
			}
			initialOrdersRefreshed := make([]bool, defaultQuotingParams.Layers*2)
			for i := range initialOrdersRefreshed {
				initialOrdersRefreshed[i] = true
			}
			clientIdIsCanonical := make([]bool, defaultQuotingParams.Layers*2)
			for i := range clientIdIsCanonical {
				clientIdIsCanonical[i] = true
			}
			verifyVaultOrders(
				initialOrdersRefreshed,
				0,
				uint32(ctx.BlockTime().Unix())+defaultQuotingParams.OrderExpirationSeconds,
				clientIdIsCanonical,
				initialOrderSides,
			)

			for _, instance := range tc.instances {
				oldBlockTime := uint32(ctx.BlockTime().Unix())
				ctx = instance.advanceBlock(ctx, tApp)
				verifyVaultOrders(
					instance.ordersRefreshed,
					oldBlockTime+defaultQuotingParams.OrderExpirationSeconds,
					uint32(ctx.BlockTime().Unix())+defaultQuotingParams.OrderExpirationSeconds,
					instance.clientIdIsCanonical,
					instance.orderSides,
				)
			}
		})
	}
}

func TestGetVaultClobOrders(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Vault status.
		vaultStatus vaulttypes.VaultStatus
		// Vault quoting params.
		vaultQuotingParams vaulttypes.QuotingParams
		// Vault ID.
		vaultId vaulttypes.VaultId
		// Vault asset.
		vaultAssetQuoteQuantums *big.Int
		// Vault inventory.
		vaultInventoryBaseQuantums *big.Int
		// Clob pair.
		clobPair clobtypes.ClobPair
		// Market param.
		marketParam pricestypes.MarketParam
		// Market price.
		marketPrice pricestypes.MarketPrice
		// Perpetual.
		perpetual perptypes.Perpetual

		/* --- Expectations --- */
		expectedOrderLayers   []uint8
		expectedOrderSides    []clobtypes.Order_Side
		expectedOrderSubticks []uint64
		expectedOrderQuantums []uint64
		expectedTimeInForce   []clobtypes.Order_TimeInForce
		expectedErr           error
	}{
		"Success - Vault Clob 0, quoting status, 2 layers, leverage 0, doesn't cross oracle price": {
			vaultStatus: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			vaultQuotingParams: vaulttypes.QuotingParams{
				Layers:                           2,       // 2 layers
				SpreadMinPpm:                     3_123,   // 31.23 bps
				SpreadBufferPpm:                  1_500,   // 15 bps
				SkewFactorPpm:                    554_321, // 0.554321
				OrderSizePctPpm:                  100_000, // 10%
				OrderExpirationSeconds:           2,       // 2 seconds
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
			},
			vaultId:                    constants.Vault_Clob0,
			vaultAssetQuoteQuantums:    big.NewInt(1_000_000_000), // 1,000 USDC
			vaultInventoryBaseQuantums: big.NewInt(0),
			clobPair:                   constants.ClobPair_Btc,
			marketParam:                constants.TestMarketParams[0],
			marketPrice: pricestypes.MarketPrice{
				Id:       0,
				Exponent: -5,
				Price:    5_000_000, // $50
			},
			perpetual:           constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			expectedOrderLayers: []uint8{0, 0, 1, 1},
			expectedOrderSides: []clobtypes.Order_Side{
				clobtypes.Order_SIDE_SELL,
				clobtypes.Order_SIDE_BUY,
				clobtypes.Order_SIDE_SELL,
				clobtypes.Order_SIDE_BUY,
			},
			// To calculate order subticks:
			// 1. spread = max(spread_min, spread_buffer + min_price_change)
			// 2. leverage = open_notional / equity
			// 3. leverage_i = leverage +/- i * order_size_pct (- for ask and + for bid)
			// 4. skew_i
			//    * for ask when long / bid when short: -skew_factor * leverage_i
			//    * for ask when short: (skew_factor * leverage_i - 1)^2 - 1
			//    * for bid when long: -((skew_factor * leverage_i + 1)^2 - 1)
			// 5. ask_spread_i = (1 + skew_i) * spread
			//    bid_spread_i = (1 - skew_i) * spread
			// 6. a_i = oraclePrice * (1 + ask_spread_i)
			//    b_i = oraclePrice * (1 - bid_spread_i)
			// 7. subticks needs to be a multiple of subticks_per_tick (round up for asks, round down for bids)
			// To calculate size of each order
			// 1. `order_size_pct_ppm * equity / oracle_price`.
			expectedOrderSubticks: []uint64{
				// spreadPpm = max(3_123, 1_500 + 50) = 3_123
				// spread = 0.003123
				// leverage = 0 / 1_000 = 0
				// oracleSubticks = 5_000_000_000 * 10^(-5 - (-8) + (-10) - (-6)) = 5e8
				// leverage_0 = leverage = 0
				// skew_0_ask = -0.554321 * 0 = 0
				// ask_spread_0 = (1 + 0) * 0.003123 = 0.003123
				// a_0 = 5e5 * (1 + 0.003123) = 501_561.5 = 501_565 (rounded up to 5)
				501_565,
				// skew_0_bid = -((0.554321 * 0 + 1)^2 - 1) = 0
				// bid_spread_0 = (1 - 0) * 0.003123 = 0.003123
				// b_0 = 5e5 * (1 - 0.003123) = 498_438.5 = 498435 (rounded down to 5)
				498_435,
				// leverage_1 = leverage - 0.1 = -0.1
				// skew_1 = 0.1 * 0.003123 * 0.554321 ~= 0.000173
				// a_1 = 5e5 * (1 + 0.0554321 + 0.003123*2) = 503209.5 ~= 503_210 (rounded up to 5)
				// skew_1_ask = -0.554321 * -0.1 = 0.0554321
				// ask_spread_1 = (1 + 0.0554321) * 0.003123 = 0.003296114448 ~= 0.003297 (rounded up to 6 decimcal places)
				// a_1 = 5e5 * (1 + 0.003296114448) = 501_648.057224 ~= 501_650 (rounded up to 5)
				501_650,
				// leverage_1 = leverage + 0.1 = 0.1
				// skew_1 = -0.1 * 0.003123 * 0.554321 = -0.000173
				// b_2 = 5e5 * (1 - 0.000173 - 0.003123*2) = 496790.5 ~= 496_790 (rounded down to 5)
				// skew_1_bid = -((0.554321 * 0.1 + 1)^2 - 1) = -0.1139369177
				// bid_spread_1 = (1 - -0.1139369177) * 0.003123 = 0.003478824994
				// b_1 = 5e5 * (1 - 0.003478824994) = 498_260.587503 ~= 498_260 (rounded down to 5)
				498_260,
			},
			// order_size = 10% * $1_000 / $50 = 2
			// order_size_base_quantums = 2 * 10^10 = 20_000_000_000
			expectedOrderQuantums: []uint64{
				20_000_000_000,
				20_000_000_000,
				20_000_000_000,
				20_000_000_000,
			},
			// post-only if increases inventory
			// vault is flat, all orders should be post-only
			expectedTimeInForce: []clobtypes.Order_TimeInForce{
				clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
				clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
				clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
				clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
			},
		},
		"Success - Vault Clob 1, quoting status, 3 layers, leverage -0.6, doesn't cross oracle price": {
			vaultStatus: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			vaultQuotingParams: vaulttypes.QuotingParams{
				Layers:                           3,         // 3 layers
				SpreadMinPpm:                     7_654,     // 76.54 bps
				SpreadBufferPpm:                  2_900,     // 29 bps
				SkewFactorPpm:                    1_234_000, // 1.234
				OrderSizePctPpm:                  100_000,   // 10%
				OrderExpirationSeconds:           4,         // 4 seconds
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
			},
			vaultId:                    constants.Vault_Clob1,
			vaultAssetQuoteQuantums:    big.NewInt(2_000_000_000), // 2,000 USDC
			vaultInventoryBaseQuantums: big.NewInt(-250_000_000),  // -0.25 ETH
			clobPair:                   constants.ClobPair_Eth,
			marketParam:                constants.TestMarketParams[1],
			marketPrice:                constants.TestMarketPrices[1],
			perpetual:                  constants.EthUsd_0DefaultFunding_9AtomicResolution,
			expectedOrderLayers:        []uint8{0, 0, 1, 1, 2, 2},
			expectedOrderSides: []clobtypes.Order_Side{
				clobtypes.Order_SIDE_SELL,
				clobtypes.Order_SIDE_BUY,
				clobtypes.Order_SIDE_SELL,
				clobtypes.Order_SIDE_BUY,
				clobtypes.Order_SIDE_SELL,
				clobtypes.Order_SIDE_BUY,
			},
			// To calculate order subticks:
			// 1. spread = max(spread_min, spread_buffer + min_price_change)
			// 2. leverage = open_notional / equity
			// 3. leverage_i = leverage +/- i * order_size_pct (- for ask and + for bid)
			// 4. skew_i
			//    * for ask when long / bid when short: -skew_factor * leverage_i
			//    * for ask when short: (skew_factor * leverage_i - 1)^2 - 1
			//    * for bid when long: -((skew_factor * leverage_i + 1)^2 - 1)
			// 5. ask_spread_i = (1 + skew_i) * spread
			//    bid_spread_i = (1 - skew_i) * spread
			// 6. a_i = oraclePrice * (1 + ask_spread_i)
			//    b_i = oraclePrice * (1 - bid_spread_i)
			// 7. subticks needs to be a multiple of subticks_per_tick (round up for asks, round down for bids)
			// To calculate size of each order
			// 1. `order_size_pct_ppm * equity / oracle_price`.
			expectedOrderSubticks: []uint64{
				// spreadPpm = max(7_654, 2_900 + 50) = 7_654
				// spread = 0.007654
				// open_notional = -250_000_000 * 10^-9 * 3_000 * 10^6 = -750_000_000
				// leverage = -750_000_000 / (2_000_000_000 - 750_000_000) = -0.6
				// oracleSubticks = 3_000_000_000 * 10^(-6 - (-9) + (-9) - (-6)) = 3e9
				// leverage_0 = leverage - 0 * 0.1 = -0.6
				// skew_ask_0 = (1.234 * -0.6 - 1)^2 - 1 = 2.02899216
				// ask_spread_0 = (1 + 2.02899216) * 0.007654 = 0.02318390599
				// a_0 = 3e9 * (1 + 0.02318390599) = 3_069_551_717.97 ~= 3_069_552_000 (round up to 1000)
				3_069_552_000,
				// skew_bid_0 = -1.234 * -0.6 = 0.7404
				// bid_spread_0 = (1 - 0.7404) * 0.007654 = 0.0019869784
				// b_0 = 3e9 * (1 - 0.0019869784) = 2_994_039_064.8 ~= 2_994_039_000 (round down to 1000)
				2_994_039_000,
				// leverage_1 = leverage - 1 * 0.1 = -0.7
				// skew_ask_1 = (1.234 * -0.7 - 1)^2 - 1 = 2.47375044
				// ask_spread_1 = (1 + 2.47375044) * 0.007654 = 0.02658808587
				// a_1 = 3e9 * (1 + 0.02658808587) = 3_079_764_257.61 ~= 3_079_765_000 (round up to 1000)
				3_079_765_000,
				// leverage_1 = leverage + 1 * 0.1 = -0.5
				// skew_bid_1 = -1.234 * -0.5 = 0.617
				// bid_spread_1 = (1 - 0.617) * 0.007654 = 0.002931482
				// b_1 = 3e9 * (1 - 0.002931482) = 2_991_205_554 ~= 2_991_205_000 (round down to 1000)
				2_991_205_000,
				// leverage_2 = leverage - 2 * 0.1 = -0.8
				// skew_ask_2 = (1.234 * -0.8 - 1)^2 - 1 = 2.94896384
				// ask_spread_2 = (1 + 2.94896384) * 0.007654 = 0.03022536923
				// a_2 = 3e9 * (1 + 0.03022536923) = 3_090_676_107.69 ~= 3_090_677_000 (round up to 1000)
				3_090_677_000,
				// leverage_2 = leverage + 2 * 0.1 = -0.4
				// skew_bid_2 = -1.234 * -0.4 = 0.4936
				// bid_spread_2 = (1 - 0.4936) * 0.007654 = 0.0038759856
				// b_2 = 3e9 * (1 - 0.0038759856) = 2_988_372_043.2 ~= 2_988_372_000 (round down to 1000)
				2_988_372_000,
			},
			// order_size = 10% * 1250 / 3000 ~= 0.04166666667
			// order_size_base_quantums = 0.04166666667e9 ~= 41_666_667
			// round down to nearest multiple of step_base_quantums=1_000.
			expectedOrderQuantums: []uint64{
				41_666_000,
				41_666_000,
				41_666_000,
				41_666_000,
				41_666_000,
				41_666_000,
			},
			// post-only if increases inventory
			// vault is short, sell orders should be post-only
			expectedTimeInForce: []clobtypes.Order_TimeInForce{
				clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
				clobtypes.Order_TIME_IN_FORCE_UNSPECIFIED,
				clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
				clobtypes.Order_TIME_IN_FORCE_UNSPECIFIED,
				clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
				clobtypes.Order_TIME_IN_FORCE_UNSPECIFIED,
			},
		},
		"Success - Vault Clob 1, quoting status, 3 layers, leverage -3, crosses oracle price": {
			vaultStatus: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			vaultQuotingParams: vaulttypes.QuotingParams{
				Layers:                           3,       // 3 layers
				SpreadMinPpm:                     3_000,   // 30 bps
				SpreadBufferPpm:                  8_500,   // 85 bps
				SkewFactorPpm:                    900_000, // 0.9
				OrderSizePctPpm:                  200_000, // 20%
				OrderExpirationSeconds:           4,       // 4 seconds
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
			},
			vaultId:                    constants.Vault_Clob1,
			vaultAssetQuoteQuantums:    big.NewInt(2_000_000_000), // 2,000 USDC
			vaultInventoryBaseQuantums: big.NewInt(-500_000_000),  // -0.5 ETH
			clobPair:                   constants.ClobPair_Eth,
			marketParam:                constants.TestMarketParams[1],
			marketPrice:                constants.TestMarketPrices[1],
			perpetual:                  constants.EthUsd_0DefaultFunding_9AtomicResolution,
			expectedOrderLayers:        []uint8{0, 0, 1, 1, 2, 2},
			expectedOrderSides: []clobtypes.Order_Side{
				clobtypes.Order_SIDE_SELL,
				clobtypes.Order_SIDE_BUY,
				clobtypes.Order_SIDE_SELL,
				clobtypes.Order_SIDE_BUY,
				clobtypes.Order_SIDE_SELL,
				clobtypes.Order_SIDE_BUY,
			},
			// To calculate order subticks:
			// 1. spread = max(spread_min, spread_buffer + min_price_change)
			// 2. leverage = open_notional / equity
			// 3. leverage_i = leverage +/- i * order_size_pct (- for ask and + for bid)
			// 4. skew_i
			//    * for ask when long / bid when short: -skew_factor * leverage_i
			//    * for ask when short: (skew_factor * leverage_i - 1)^2 - 1
			//    * for bid when long: -((skew_factor * leverage_i + 1)^2 - 1)
			// 5. ask_spread_i = (1 + skew_i) * spread
			//    bid_spread_i = (1 - skew_i) * spread
			// 6. a_i = oraclePrice * (1 + ask_spread_i)
			//    b_i = oraclePrice * (1 - bid_spread_i)
			// 7. subticks needs to be a multiple of subticks_per_tick (round up for asks, round down for bids)
			// To calculate size of each order
			// 1. `order_size_pct_ppm * equity / oracle_price`.
			expectedOrderSubticks: []uint64{
				// spreadPpm = max(3_000, 8_500 + 50) = 8_550
				// spread = 0.00855
				// open_notional = -500_000_000 * 10^-9 * 3_000 * 10^6 = -1_500_000_000
				// leverage = -1_500_000_000 / (2_000_000_000 - 1_500_000_000) = -3
				// oracleSubticks = 3_000_000_000 * 10^(-6 - (-9) + (-9) - (-6)) = 3e9
				// leverage_0 = leverage - 0 * 0.2 = -3
				// skew_ask_0 = (0.9 * -3 - 1)^2 - 1 = 12.69
				// ask_spread_0 = (1 + 12.69) * 0.00855 = 0.1170495
				// a_0 = 3e9 * (1 + 0.1170495) = 3_351_148_500 ~= 3_351_149_000 (round up to 1000)
				3_351_149_000,
				// skew_bid_0 = -0.9 * -3 = 2.7
				// bid_spread_0 = (1 - 2.7) * 0.00855 = -0.014535
				// b_0 = 3e9 * (1 - -0.014535) = 3_043_605_000
				3_043_605_000,
				// leverage_1 = leverage - 1 * 0.2
				// skew_ask_1 = (0.9 * -3.2 - 1)^2 - 1 = 14.0544
				// ask_spread_1 = (1 + 14.0544) * 0.00855 = 0.12871512
				// a_1 = 3e9 * (1 + 0.12871512) = 3_386_145_360 ~= 3_386_146_000 (round up to 1000)
				3_386_146_000,
				// leverage_1 = leverage + 1 * 0.2
				// skew_bid_1 = -0.9 * -2.8 = 2.52
				// bid_spread_1 = (1 - 2.52) * 0.00855 = -0.012996
				// b_1 = 3e9 * (1 - -0.012996) = 3_038_988_000
				3_038_988_000,
				// leverage_2 = leverage - 2 * 0.2
				// skew_ask_2 = (0.9 * -3.4 - 1)^2 - 1 = 15.4836
				// ask_spread_2 = (1 + 15.4836) * 0.00855 = 0.14093478
				// a_2 = 3e9 * (1 + 0.14093478) = 3_422_804_340 ~= 3_422_805_000 (round up to 1000)
				3_422_805_000,
				// leverage_2 = leverage + 2 * 0.2
				// skew_bid_2 = -0.9 * -2.6 = 2.34
				// bid_spread_2 = (1 - 2.34) * 0.00855 = -0.011457
				// b_2 = 3e9 * (1 - -0.011457) = 3_034_371_000
				3_034_371_000,
			},
			// order_size = 20% * 500 / 3000 ~= 0.0333333333
			// order_size_base_quantums = 0.0333333333e9 ~= 33_333_333.33
			// round down to nearest multiple of step_base_quantums=1_000.
			expectedOrderQuantums: []uint64{
				33_333_000,
				33_333_000,
				33_333_000,
				33_333_000,
				33_333_000,
				33_333_000,
			},
			// post-only if increases inventory
			// vault is short, sell orders should be post-only
			expectedTimeInForce: []clobtypes.Order_TimeInForce{
				clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
				clobtypes.Order_TIME_IN_FORCE_UNSPECIFIED,
				clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
				clobtypes.Order_TIME_IN_FORCE_UNSPECIFIED,
				clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
				clobtypes.Order_TIME_IN_FORCE_UNSPECIFIED,
			},
		},
		"Success - Vault Clob 1, quoting status, 2 layers, leverage 3, crosses oracle price": {
			vaultStatus: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			vaultQuotingParams: vaulttypes.QuotingParams{
				Layers:                           2,         // 2 layers
				SpreadMinPpm:                     3_000,     // 30 bps
				SpreadBufferPpm:                  1_500,     // 15 bps
				SkewFactorPpm:                    500_000,   // 0.5
				OrderSizePctPpm:                  1_000_000, // 100%
				OrderExpirationSeconds:           4,         // 4 seconds
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
			},
			vaultId:                    constants.Vault_Clob1,
			vaultAssetQuoteQuantums:    big.NewInt(-2_000_000_000), // -2,000 USDC
			vaultInventoryBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			clobPair:                   constants.ClobPair_Eth,
			marketParam:                constants.TestMarketParams[1],
			marketPrice:                constants.TestMarketPrices[1],
			perpetual:                  constants.EthUsd_0DefaultFunding_9AtomicResolution,
			expectedOrderLayers:        []uint8{0, 0, 1, 1},
			expectedOrderSides: []clobtypes.Order_Side{
				clobtypes.Order_SIDE_SELL,
				clobtypes.Order_SIDE_BUY,
				clobtypes.Order_SIDE_SELL,
				clobtypes.Order_SIDE_BUY,
			},
			// To calculate order subticks:
			// 1. spread = max(spread_min, spread_buffer + min_price_change)
			// 2. leverage = open_notional / equity
			// 3. leverage_i = leverage +/- i * order_size_pct (- for ask and + for bid)
			// 4. skew_i
			//    * for ask when long / bid when short: -skew_factor * leverage_i
			//    * for ask when short: (skew_factor * leverage_i - 1)^2 - 1
			//    * for bid when long: -((skew_factor * leverage_i + 1)^2 - 1)
			// 5. ask_spread_i = (1 + skew_i) * spread
			//    bid_spread_i = (1 - skew_i) * spread
			// 6. a_i = oraclePrice * (1 + ask_spread_i)
			//    b_i = oraclePrice * (1 - bid_spread_i)
			// 7. subticks needs to be a multiple of subticks_per_tick (round up for asks, round down for bids)
			// To calculate size of each order
			// 1. `order_size_pct_ppm * equity / oracle_price`.
			expectedOrderSubticks: []uint64{
				// spreadPpm = max(3_000, 1_500 + 50) = 3_000
				// spread = 0.003
				// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
				// leverage = 3_000_000_000 / (-2_000_000_000 + 3_000_000_000) = 3
				// oracleSubticks = 3_000_000_000 * 10^(-6 - (-9) + (-9) - (-6)) = 3e9
				// leverage_0 = leverage - 0 * 1 = 3
				// skew_ask_0 = -0.5 * 3 = -1.5
				// ask_spread_0 = (1 + -1.5) * 0.003 = -0.0015
				// a_0 = 3e9 * (1 + -0.0015) = 2_995_500_000
				2_995_500_000,
				// skew_bid_0 = -((0.5 * 3 + 1)^2 - 1) = -5.25
				// bid_spread_0 = (1 - -5.25) * 0.003 = 0.01875
				// b_0 = 3e9 * (1 - 0.01875) = 2_943_750_000
				2_943_750_000,
				// leverage_1 = leverage - 1 * 1 = 2
				// skew_ask_1 = -0.5 * 2 = -1
				// ask_spread_1 = (1 + -1) * 0.003 = 0
				// a_1 = 3e9 * (1 + 0) = 3_000_000_000
				3_000_000_000,
				// leverage_1 = leverage + 1 * 1 = 4
				// skew_bid_1 = -((0.5 * 4 + 1)^2 - 1) = -8
				// bid_spread_1 = (1 - -8) * 0.003 = 0.027
				// b_1 = 3e9 * (1 - 0.027) = 2_919_000_000
				2_919_000_000,
			},
			// order_size = 100% * 1000 / 3000 ~= 0.333333333
			// order_size_base_quantums = 0.333333333e9 ~= 333_333_333.33
			// round down to nearest multiple of step_base_quantums=1_000.
			expectedOrderQuantums: []uint64{
				333_333_000,
				333_333_000,
				333_333_000,
				333_333_000,
			},
			// post-only if increases inventory
			// vault is long, buy orders should be post-only
			expectedTimeInForce: []clobtypes.Order_TimeInForce{
				clobtypes.Order_TIME_IN_FORCE_UNSPECIFIED,
				clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
				clobtypes.Order_TIME_IN_FORCE_UNSPECIFIED,
				clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
			},
		},
		"Success - Vault Clob 1, close-only status, 3 layers, leverage -0.6, buy orders only": {
			vaultStatus: vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
			vaultQuotingParams: vaulttypes.QuotingParams{
				Layers:                           3,         // 3 layers
				SpreadMinPpm:                     7_654,     // 76.54 bps
				SpreadBufferPpm:                  2_900,     // 29 bps
				SkewFactorPpm:                    1_234_000, // 1.234
				OrderSizePctPpm:                  100_000,   // 10%
				OrderExpirationSeconds:           4,         // 4 seconds
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
			},
			vaultId:                    constants.Vault_Clob1,
			vaultAssetQuoteQuantums:    big.NewInt(2_000_000_000), // 2,000 USDC
			vaultInventoryBaseQuantums: big.NewInt(-250_000_000),  // -0.25 ETH
			clobPair:                   constants.ClobPair_Eth,
			marketParam:                constants.TestMarketParams[1],
			marketPrice:                constants.TestMarketPrices[1],
			perpetual:                  constants.EthUsd_0DefaultFunding_9AtomicResolution,
			expectedOrderLayers:        []uint8{0, 1, 2},
			expectedOrderSides: []clobtypes.Order_Side{
				clobtypes.Order_SIDE_BUY,
				clobtypes.Order_SIDE_BUY,
				clobtypes.Order_SIDE_BUY,
			},
			// To calculate order subticks:
			// 1. spread = max(spread_min, spread_buffer + min_price_change)
			// 2. leverage = open_notional / equity
			// 3. leverage_i = leverage +/- i * order_size_pct (- for ask and + for bid)
			// 4. skew_i
			//    * for ask when long / bid when short: -skew_factor * leverage_i
			//    * for ask when short: (skew_factor * leverage_i - 1)^2 - 1
			//    * for bid when long: -((skew_factor * leverage_i + 1)^2 - 1)
			// 5. ask_spread_i = (1 + skew_i) * spread
			//    bid_spread_i = (1 - skew_i) * spread
			// 6. a_i = oraclePrice * (1 + ask_spread_i)
			//    b_i = oraclePrice * (1 - bid_spread_i)
			// 7. subticks needs to be a multiple of subticks_per_tick (round up for asks, round down for bids)
			// To calculate size of each order
			// 1. `order_size_pct_ppm * equity / oracle_price`.
			expectedOrderSubticks: []uint64{
				// spreadPpm = max(7_654, 2_900 + 50) = 7_654
				// spread = 0.007654
				// open_notional = -250_000_000 * 10^-9 * 3_000 * 10^6 = -750_000_000
				// leverage = -750_000_000 / (2_000_000_000 - 750_000_000) = -0.6
				// oracleSubticks = 3_000_000_000 * 10^(-6 - (-9) + (-9) - (-6)) = 3e9
				// leverage_0 = leverage - 0 * 0.1 = -0.6
				// skew_bid_0 = -1.234 * -0.6 = 0.7404
				// bid_spread_0 = (1 - 0.7404) * 0.007654 = 0.0019869784
				// b_0 = 3e9 * (1 - 0.0019869784) = 2_994_039_064.8 ~= 2_994_039_000 (round down to 1000)
				2_994_039_000,
				// leverage_1 = leverage + 1 * 0.1 = -0.5
				// skew_bid_1 = -1.234 * -0.5 = 0.617
				// bid_spread_1 = (1 - 0.617) * 0.007654 = 0.002931482
				// b_1 = 3e9 * (1 - 0.002931482) = 2_991_205_554 ~= 2_991_205_000 (round down to 1000)
				2_991_205_000,
				// leverage_2 = leverage + 2 * 0.1 = -0.4
				// skew_bid_2 = -1.234 * -0.4 = 0.4936
				// bid_spread_2 = (1 - 0.4936) * 0.007654 = 0.0038759856
				// b_2 = 3e9 * (1 - 0.0038759856) = 2_988_372_043.2 ~= 2_988_372_000 (round down to 1000)
				2_988_372_000,
			},
			// order_size = 10% * 1250 / 3000 ~= 0.04166666667
			// order_size_base_quantums = 0.04166666667e9 ~= 41_666_667
			// round down to nearest multiple of step_base_quantums=1_000.
			expectedOrderQuantums: []uint64{
				41_666_000,
				41_666_000,
				41_666_000,
			},
			expectedTimeInForce: []clobtypes.Order_TimeInForce{
				clobtypes.Order_TIME_IN_FORCE_UNSPECIFIED,
				clobtypes.Order_TIME_IN_FORCE_UNSPECIFIED,
				clobtypes.Order_TIME_IN_FORCE_UNSPECIFIED,
			},
		},
		"Success - Vault Clob 1, close-only status, 3 layers, leverage -0.6, no orders as size is rounded to 0": {
			vaultStatus: vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
			vaultQuotingParams: vaulttypes.QuotingParams{
				Layers:                           3,         // 3 layers
				SpreadMinPpm:                     7_654,     // 76.54 bps
				SpreadBufferPpm:                  2_900,     // 29 bps
				SkewFactorPpm:                    1_234_000, // 1.234
				OrderSizePctPpm:                  100,       // 0.01%
				OrderExpirationSeconds:           4,         // 4 seconds
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
			},
			vaultId:                    constants.Vault_Clob1,
			vaultAssetQuoteQuantums:    big.NewInt(760_000_000),  // 1,000 USDC
			vaultInventoryBaseQuantums: big.NewInt(-250_000_000), // -0.25 ETH
			clobPair:                   constants.ClobPair_Eth,
			marketParam:                constants.TestMarketParams[1],
			marketPrice:                constants.TestMarketPrices[1],
			perpetual:                  constants.EthUsd_0DefaultFunding_9AtomicResolution,
			expectedOrderLayers:        []uint8{},
			expectedOrderSides:         []clobtypes.Order_Side{},
			expectedOrderSubticks:      []uint64{},
			// open_notional = -250_000_000 * 10^-9 * 3_000 * 10^6 = -750_000_000
			// equity = 760_000_000 - 750_000_000 = 10_000_000
			// order_size = `order_size_pct_ppm * equity / oracle_price`
			// = 0.01% * 10 / 3000 ~= 3.3333e-7
			// order_size_base_quantums = 3.3333e-7 * 10^-9 ~= 333 ~= 0 after rounding
			// down to nearest multiple of step_base_quantums=1_000.
			// Thus no orders.
			expectedOrderQuantums: []uint64{},
		},
		"Success - Vault Clob 1, close-only status, 2 layers, leverage 0.6, sell orders only, " +
			"total size of orders capped to position size": {
			vaultStatus: vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
			vaultQuotingParams: vaulttypes.QuotingParams{
				Layers:                           2,         // 2 layers
				SpreadMinPpm:                     7_654,     // 76.54 bps
				SpreadBufferPpm:                  2_900,     // 29 bps
				SkewFactorPpm:                    1_234_000, // 1.234
				OrderSizePctPpm:                  500_000,   // 50%
				OrderExpirationSeconds:           4,         // 4 seconds
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
			},
			vaultId:                    constants.Vault_Clob1,
			vaultAssetQuoteQuantums:    big.NewInt(1_000_000_000), // 1000 USDC
			vaultInventoryBaseQuantums: big.NewInt(500_000_000),   // 0.5 ETH
			clobPair:                   constants.ClobPair_Eth,
			marketParam:                constants.TestMarketParams[1],
			marketPrice:                constants.TestMarketPrices[1],
			perpetual:                  constants.EthUsd_0DefaultFunding_9AtomicResolution,
			expectedOrderLayers:        []uint8{0, 1},
			expectedOrderSides: []clobtypes.Order_Side{
				clobtypes.Order_SIDE_SELL,
				clobtypes.Order_SIDE_SELL,
			},
			// To calculate order subticks:
			// 1. spread = max(spread_min, spread_buffer + min_price_change)
			// 2. leverage = open_notional / equity
			// 3. leverage_i = leverage +/- i * order_size_pct (- for ask and + for bid)
			// 4. skew_i
			//    * for ask when long / bid when short: -skew_factor * leverage_i
			//    * for ask when short: (skew_factor * leverage_i - 1)^2 - 1
			//    * for bid when long: -((skew_factor * leverage_i + 1)^2 - 1)
			// 5. ask_spread_i = (1 + skew_i) * spread
			//    bid_spread_i = (1 - skew_i) * spread
			// 6. a_i = oraclePrice * (1 + ask_spread_i)
			//    b_i = oraclePrice * (1 - bid_spread_i)
			// 7. subticks needs to be a multiple of subticks_per_tick (round up for asks, round down for bids)
			// To calculate size of each order
			// 1. `order_size_pct_ppm * equity / oracle_price`.
			expectedOrderSubticks: []uint64{
				// spreadPpm = max(7_654, 2_900 + 50) = 7_654
				// spread = 0.007654
				// open_notional = 500_000_000 * 10^-9 * 3_000 * 10^6 = 1_500_000_000
				// leverage = 1_500_000_000 / (1_000_000_000 + 1_500_000_000) = 0.6
				// oracleSubticks = 3_000_000_000 * 10^(-6 - (-9) + (-9) - (-6)) = 3e9
				// leverage_0 = leverage - 0 * 0.5 = 0.6
				// skew_ask_0 = -1.234 * 0.6 = -0.7404
				// ask_spread_0 = (1 - 0.7404) * 0.007654 = 0.0019869784
				// a_0 = 3e9 * (1 + 0.0019869784) = 3_005_960_935.2 ~= 3_005_961_000 (round up to 1000)
				3_005_961_000,
				// leverage_1 = leverage - 1 * 0.5 = 0.1
				// skew_ask_0 = -1.234 * 0.1 = -0.1234
				// ask_spread_0 = (1 - 0.1234) * 0.007654 = 0.0067094964
				// a_0 = 3e9 * (1 + 0.0067094964) = 3020128489.2 ~= 3_020_129_000 (round up to 1000)
				3_020_129_000,
			},
			// order_size = 50% * 2500 / 3000 ~= 0.4166666667
			// order_size_base_quantums = 0.4166666667e9 ~= 416_666_667
			// round down to nearest multiple of step_base_quantums=1_000.
			expectedOrderQuantums: []uint64{
				416_666_000,
				// this order size is capped at `position size - sum of previous order sizes`
				// = 500_000_000 - 416_666_000 = 83_334_000
				83_334_000,
			},
			expectedTimeInForce: []clobtypes.Order_TimeInForce{
				clobtypes.Order_TIME_IN_FORCE_UNSPECIFIED,
				clobtypes.Order_TIME_IN_FORCE_UNSPECIFIED,
			},
		},
		"Success - Get orders from Vault for Clob Pair 1, No Orders due to Zero Order Size": {
			vaultStatus: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			vaultQuotingParams: vaulttypes.QuotingParams{
				Layers:                           2,       // 2 layers
				SpreadMinPpm:                     3_000,   // 30 bps
				SpreadBufferPpm:                  1_500,   // 15 bps
				SkewFactorPpm:                    500_000, // 0.5
				OrderSizePctPpm:                  1_000,   // 0.1%
				OrderExpirationSeconds:           2,       // 2 seconds
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
			},
			vaultId:                    constants.Vault_Clob1,
			vaultAssetQuoteQuantums:    big.NewInt(1_000_000), // 1 USDC
			vaultInventoryBaseQuantums: big.NewInt(0),
			clobPair:                   constants.ClobPair_Eth,
			marketParam:                constants.TestMarketParams[1],
			marketPrice:                constants.TestMarketPrices[1],
			perpetual:                  constants.EthUsd_0DefaultFunding_9AtomicResolution,
			expectedOrderLayers:        []uint8{},
			expectedOrderSides:         []clobtypes.Order_Side{},
			expectedOrderSubticks:      []uint64{},
			// order_size = 0.1% * 1 / 3_000 ~= 0.00000033333
			// order_size_base_quantums = 0.000033333e9 = 333
			// round down to nearest multiple of step_base_quantums=1_000.
			// order size is 0.
			expectedOrderQuantums: []uint64{},
			expectedTimeInForce:   []clobtypes.Order_TimeInForce{},
		},
		"Success - Clob Pair doesn't exist, Empty orders": {
			vaultStatus:           vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			vaultQuotingParams:    vaulttypes.DefaultQuotingParams(),
			vaultId:               constants.Vault_Clob0,
			clobPair:              constants.ClobPair_Eth,
			marketParam:           constants.TestMarketParams[1],
			marketPrice:           constants.TestMarketPrices[1],
			perpetual:             constants.EthUsd_NoMarginRequirement,
			expectedOrderLayers:   []uint8{},
			expectedOrderSides:    []clobtypes.Order_Side{},
			expectedOrderSubticks: []uint64{},
			expectedOrderQuantums: []uint64{},
			expectedTimeInForce:   []clobtypes.Order_TimeInForce{},
		},
		"Success - Clob Pair in status final settlement, Empty orders": {
			vaultStatus:        vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			vaultQuotingParams: vaulttypes.DefaultQuotingParams(),
			vaultId:            constants.Vault_Clob1,
			clobPair: clobtypes.ClobPair{
				Id: 1,
				Metadata: &clobtypes.ClobPair_PerpetualClobMetadata{
					PerpetualClobMetadata: &clobtypes.PerpetualClobMetadata{
						PerpetualId: 1,
					},
				},
				StepBaseQuantums:          1000,
				SubticksPerTick:           1000,
				QuantumConversionExponent: -9,
				Status:                    clobtypes.ClobPair_STATUS_FINAL_SETTLEMENT,
			},
			marketParam:           constants.TestMarketParams[1],
			marketPrice:           constants.TestMarketPrices[1],
			perpetual:             constants.EthUsd_NoMarginRequirement,
			expectedOrderLayers:   []uint8{},
			expectedOrderSides:    []clobtypes.Order_Side{},
			expectedOrderSubticks: []uint64{},
			expectedOrderQuantums: []uint64{},
		},
		"Success - Vault Clob 1, stand-by status, 3 layers, leverage -3, Empty orders": {
			vaultStatus:                vaulttypes.VaultStatus_VAULT_STATUS_STAND_BY,
			vaultQuotingParams:         vaulttypes.DefaultQuotingParams(),
			vaultId:                    constants.Vault_Clob1,
			vaultAssetQuoteQuantums:    big.NewInt(2_000_000_000), // 2,000 USDC
			vaultInventoryBaseQuantums: big.NewInt(-500_000_000),  // -0.5 ETH
			clobPair:                   constants.ClobPair_Eth,
			marketParam:                constants.TestMarketParams[1],
			marketPrice:                constants.TestMarketPrices[1],
			perpetual:                  constants.EthUsd_0DefaultFunding_9AtomicResolution,
			expectedOrderLayers:        []uint8{},
			expectedOrderSides:         []clobtypes.Order_Side{},
			expectedOrderSubticks:      []uint64{},
			expectedOrderQuantums:      []uint64{},
		},
		"Success - Vault Clob 0, close-only status, 2 layers, leverage 0, Empty orders": {
			vaultStatus: vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
			vaultQuotingParams: vaulttypes.QuotingParams{
				Layers:                           2,       // 2 layers
				SpreadMinPpm:                     3_123,   // 31.23 bps
				SpreadBufferPpm:                  1_500,   // 15 bps
				SkewFactorPpm:                    554_321, // 0.554321
				OrderSizePctPpm:                  100_000, // 10%
				OrderExpirationSeconds:           2,       // 2 seconds
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
			},
			vaultId:                    constants.Vault_Clob0,
			vaultAssetQuoteQuantums:    big.NewInt(1_000_000_000), // 1,000 USDC
			vaultInventoryBaseQuantums: big.NewInt(0),
			clobPair:                   constants.ClobPair_Btc,
			marketParam:                constants.TestMarketParams[0],
			marketPrice: pricestypes.MarketPrice{
				Id:       0,
				Exponent: -5,
				Price:    5_000_000, // $50
			},
			perpetual:             constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			expectedOrderLayers:   []uint8{},
			expectedOrderSides:    []clobtypes.Order_Side{},
			expectedOrderSubticks: []uint64{},
			expectedOrderQuantums: []uint64{},
			expectedTimeInForce:   []clobtypes.Order_TimeInForce{},
		},
		"Error - Vault equity is zero": {
			vaultStatus:                vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			vaultQuotingParams:         vaulttypes.DefaultQuotingParams(),
			vaultId:                    constants.Vault_Clob0,
			vaultAssetQuoteQuantums:    big.NewInt(0),
			vaultInventoryBaseQuantums: big.NewInt(0),
			clobPair:                   constants.ClobPair_Btc,
			marketParam:                constants.TestMarketParams[0],
			marketPrice:                constants.TestMarketPrices[0],
			perpetual:                  constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			expectedErr:                vaulttypes.ErrNonPositiveEquity,
		},
		"Error - Vault equity is negative": {
			vaultStatus:                vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			vaultQuotingParams:         vaulttypes.DefaultQuotingParams(),
			vaultId:                    constants.Vault_Clob0,
			vaultAssetQuoteQuantums:    big.NewInt(5_000_000), // 5 USDC
			vaultInventoryBaseQuantums: big.NewInt(-10_000_000),
			clobPair:                   constants.ClobPair_Btc,
			marketParam:                constants.TestMarketParams[0],
			marketPrice:                constants.TestMarketPrices[0],
			perpetual:                  constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			expectedErr:                vaulttypes.ErrNonPositiveEquity,
		},
		"Error - Market price is zero": {
			vaultStatus:                vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			vaultQuotingParams:         vaulttypes.DefaultQuotingParams(),
			vaultId:                    constants.Vault_Clob0,
			vaultAssetQuoteQuantums:    big.NewInt(1_000_000_000), // 1,000 USDC
			vaultInventoryBaseQuantums: big.NewInt(0),
			clobPair:                   constants.ClobPair_Btc,
			marketParam:                constants.TestMarketParams[0],
			marketPrice: pricestypes.MarketPrice{
				Id:       0,
				Exponent: -5,
				Price:    0,
			},
			perpetual:   constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			expectedErr: vaulttypes.ErrZeroMarketPrice,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize tApp and ctx.
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Initialize prices module with test market param and market price.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *pricestypes.GenesisState) {
						genesisState.MarketParams = []pricestypes.MarketParam{tc.marketParam}
						genesisState.MarketPrices = []pricestypes.MarketPrice{tc.marketPrice}
					},
				)
				// Initialize perpetuals module with test perpetual.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.LiquidityTiers = constants.LiquidityTiers
						genesisState.Perpetuals = []perptypes.Perpetual{tc.perpetual}
					},
				)
				// Initialize clob module with test clob pair.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.ClobPairs = []clobtypes.ClobPair{tc.clobPair}
					},
				)
				// Initialize subaccounts module with vault's equity and inventory.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						assetPositions := []*satypes.AssetPosition{}
						if tc.vaultAssetQuoteQuantums != nil && tc.vaultAssetQuoteQuantums.Sign() != 0 {
							assetPositions = append(
								assetPositions,
								&satypes.AssetPosition{
									AssetId:  assettypes.AssetUsdc.Id,
									Quantums: dtypes.NewIntFromBigInt(tc.vaultAssetQuoteQuantums),
								},
							)
						}
						perpPositions := []*satypes.PerpetualPosition{}
						if tc.vaultInventoryBaseQuantums != nil && tc.vaultInventoryBaseQuantums.Sign() != 0 {
							perpPositions = append(
								perpPositions,
								testutil.CreateSinglePerpetualPosition(
									tc.perpetual.Params.Id,
									tc.vaultInventoryBaseQuantums,
									big.NewInt(0),
									big.NewInt(0),
								),
							)
						}
						genesisState.Subaccounts = []satypes.Subaccount{
							{
								Id:                 tc.vaultId.ToSubaccountId(),
								AssetPositions:     assetPositions,
								PerpetualPositions: perpPositions,
							},
						}
					},
				)
				// Initialize vault module with vault's parameters.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vaulttypes.GenesisState) {
						genesisState.Vaults = []vaulttypes.Vault{
							{
								VaultId: tc.vaultId,
								VaultParams: vaulttypes.VaultParams{
									Status:        tc.vaultStatus,
									QuotingParams: &tc.vaultQuotingParams,
								},
							},
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			// Get vault orders.
			orders, err := tApp.App.VaultKeeper.GetVaultClobOrders(ctx, tc.vaultId)
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
				return
			}
			require.NoError(t, err)

			// Get expected orders.
			buildVaultClobOrder := func(
				layer uint8,
				side clobtypes.Order_Side,
				quantums uint64,
				subticks uint64,
				timeInForce clobtypes.Order_TimeInForce,
			) *clobtypes.Order {
				return &clobtypes.Order{
					OrderId: clobtypes.OrderId{
						SubaccountId: *tc.vaultId.ToSubaccountId(),
						ClientId:     vaulttypes.GetVaultClobOrderClientId(side, layer),
						OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
						ClobPairId:   tc.vaultId.Number,
					},
					Side:     side,
					Quantums: quantums,
					Subticks: subticks,
					GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(ctx.BlockTime().Unix()) + tc.vaultQuotingParams.OrderExpirationSeconds,
					},
					TimeInForce: timeInForce,
				}
			}
			expectedOrders := make([]*clobtypes.Order, 0)
			for i := 0; i < len(tc.expectedOrderQuantums); i++ {
				expectedOrders = append(
					expectedOrders,
					buildVaultClobOrder(
						tc.expectedOrderLayers[i],
						tc.expectedOrderSides[i],
						tc.expectedOrderQuantums[i],
						tc.expectedOrderSubticks[i],
						tc.expectedTimeInForce[i],
					),
				)
			}

			// Compare expected orders with actual orders.
			require.Equal(
				t,
				expectedOrders,
				orders,
			)
		})
	}
}

func TestGetVaultClobOrderIds(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Vault ID.
		vaultId vaulttypes.VaultId
		// Layers.
		layers uint32

		/* --- Expectations --- */
		expectedNumOrders uint32
	}{
		"Vault Clob 0, 2 layers": {
			vaultId:           constants.Vault_Clob0,
			layers:            2,
			expectedNumOrders: 4,
		},
		"Vault Clob 1, 7 layers": {
			vaultId:           constants.Vault_Clob1,
			layers:            7,
			expectedNumOrders: 14,
		},
		"Vault Clob 0, 0 layers": {
			vaultId:           constants.Vault_Clob0,
			layers:            0,
			expectedNumOrders: 0,
		},
		"Vault Clob 797, 2 layers": {
			vaultId: vaulttypes.VaultId{
				Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
				Number: 797,
			},
			layers:            2,
			expectedNumOrders: 4,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			k := tApp.App.VaultKeeper
			ctx := tApp.InitChain()

			// Set number of layers.
			quotingParams := constants.QuotingParams
			quotingParams.Layers = tc.layers
			err := k.SetVaultParams(ctx, tc.vaultId, vaulttypes.VaultParams{
				Status:        vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
				QuotingParams: &quotingParams,
			})
			require.NoError(t, err)

			// Construct expected order IDs.
			expectedOrderIds := make([]*clobtypes.OrderId, tc.layers*2)
			for i := uint32(0); i < tc.layers; i++ {
				expectedOrderIds[2*i] = &clobtypes.OrderId{
					SubaccountId: *tc.vaultId.ToSubaccountId(),
					ClientId:     vaulttypes.GetVaultClobOrderClientId(clobtypes.Order_SIDE_SELL, uint8(i)),
					OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
					ClobPairId:   tc.vaultId.Number,
				}
				expectedOrderIds[2*i+1] = &clobtypes.OrderId{
					SubaccountId: *tc.vaultId.ToSubaccountId(),
					ClientId:     vaulttypes.GetVaultClobOrderClientId(clobtypes.Order_SIDE_BUY, uint8(i)),
					OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
					ClobPairId:   tc.vaultId.Number,
				}
			}

			// Verify order IDs.
			require.Equal(t, tc.expectedNumOrders, uint32(len(expectedOrderIds)))
			require.Equal(t, expectedOrderIds, k.GetVaultClobOrderIds(ctx, tc.vaultId))
		})
	}
}

func TestGetSetMostRecentClientIds(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Vault ID.
		vaultId vaulttypes.VaultId
		// Client IDs.
		clientIds []uint32
	}{
		"Vault Clob 0, non-existent client IDs": {
			vaultId: constants.Vault_Clob0,
		},
		"Vault Clob 0, empty client IDs": {
			vaultId:   constants.Vault_Clob0,
			clientIds: []uint32{},
		},
		"Vault Clob 0, 4 client IDs": {
			vaultId:   constants.Vault_Clob0,
			clientIds: []uint32{111, 222, 333, 444},
		},
		"Vault Clob 1, 6 client IDs": {
			vaultId:   constants.Vault_Clob0,
			clientIds: []uint32{0, 1, 987654321, 555666, 3453, 1010101010},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			// Set most recent client IDs if provided.
			if tc.clientIds != nil {
				k.SetMostRecentClientIds(ctx, tc.vaultId, tc.clientIds)
			}

			// Verify most recent client IDs.
			if tc.clientIds == nil {
				require.Empty(t, k.GetMostRecentClientIds(ctx, tc.vaultId))
			} else {
				require.Equal(t, tc.clientIds, k.GetMostRecentClientIds(ctx, tc.vaultId))
			}
		})
	}
}
