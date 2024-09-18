package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetDefaultQuotingParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	// Params should have default values at genesis.
	params := k.GetDefaultQuotingParams(ctx)
	require.Equal(t, types.DefaultQuotingParams(), params)

	// Set new params and get.
	newParams := types.QuotingParams{
		Layers:                           3,
		SpreadMinPpm:                     4_000,
		SpreadBufferPpm:                  2_000,
		SkewFactorPpm:                    999_999,
		OrderSizePctPpm:                  200_000,
		OrderExpirationSeconds:           10,
		ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
	}
	err := k.SetDefaultQuotingParams(ctx, newParams)
	require.NoError(t, err)
	require.Equal(t, newParams, k.GetDefaultQuotingParams(ctx))

	// Set invalid params and get.
	invalidParams := types.QuotingParams{
		Layers:                           3,
		SpreadMinPpm:                     4_000,
		SpreadBufferPpm:                  2_000,
		SkewFactorPpm:                    1_000_000,
		OrderSizePctPpm:                  200_000,
		OrderExpirationSeconds:           0, // invalid
		ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
	}
	err = k.SetDefaultQuotingParams(ctx, invalidParams)
	require.Error(t, err)
	require.Equal(t, newParams, k.GetDefaultQuotingParams(ctx))
}

func TestGetSetVaultParams(t *testing.T) {
	tests := map[string]struct {
		// Vault id.
		vaultId types.VaultId
		// Vault params to set.
		vaultParams *types.VaultParams
		// Expected on-chain indexer events
		expectedIndexerEvents []*indexerevents.UpsertVaultEventV1
		// Expected error.
		expectedErr error
	}{
		"Success - Vault Clob 0": {
			vaultId:     constants.Vault_Clob0,
			vaultParams: &constants.VaultParams,
			expectedIndexerEvents: []*indexerevents.UpsertVaultEventV1{
				{
					Address:    constants.Vault_Clob0.ToModuleAccountAddress(),
					ClobPairId: constants.Vault_Clob0.Number,
					Status:     v1.VaultStatusToIndexerVaultStatus(constants.VaultParams.Status),
				},
			},
		},
		"Success - Vault Clob 1": {
			vaultId:     constants.Vault_Clob1,
			vaultParams: &constants.VaultParams,
			expectedIndexerEvents: []*indexerevents.UpsertVaultEventV1{
				{
					Address:    constants.Vault_Clob1.ToModuleAccountAddress(),
					ClobPairId: constants.Vault_Clob1.Number,
					Status:     v1.VaultStatusToIndexerVaultStatus(constants.VaultParams.Status),
				},
			},
		},
		"Success - Non-existent Vault Params": {
			vaultId:     constants.Vault_Clob1,
			vaultParams: nil,
		},
		"Failure - Unspecified Status": {
			vaultId: constants.Vault_Clob0,
			vaultParams: &types.VaultParams{
				QuotingParams: &constants.QuotingParams,
			},
			expectedErr: types.ErrUnspecifiedVaultStatus,
		},
		"Failure - Invalid Quoting Params": {
			vaultId: constants.Vault_Clob0,
			vaultParams: &types.VaultParams{
				Status:        types.VaultStatus_VAULT_STATUS_STAND_BY,
				QuotingParams: &constants.InvalidQuotingParams,
			},
			expectedErr: types.ErrInvalidOrderExpirationSeconds,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
			appOpts := map[string]interface{}{
				indexer.MsgSenderInstanceForTest: msgSender,
			}
			tApp := testapp.NewTestAppBuilder(t).WithAppOptions(appOpts).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			if tc.vaultParams == nil {
				_, exists := k.GetVaultParams(ctx, tc.vaultId)
				require.False(t, exists)
				return
			}

			err := k.SetVaultParams(ctx, tc.vaultId, *tc.vaultParams)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				_, exists := k.GetVaultParams(ctx, tc.vaultId)
				require.False(t, exists)
			} else {
				require.NoError(t, err)
				p, exists := k.GetVaultParams(ctx, tc.vaultId)
				require.True(t, exists)
				require.Equal(t, *tc.vaultParams, p)
			}

			if tc.expectedErr == nil && tc.vaultParams != nil {
				upsertVaultEventsInBlock := getUpsertVaultEventsFromIndexerBlock(ctx, &k)
				require.ElementsMatch(t, tc.expectedIndexerEvents, upsertVaultEventsInBlock)
			}
		})
	}
}

func TestGetVaultQuotingParams(t *testing.T) {
	tests := map[string]struct {
		/* Setup */
		// Vault id.
		vaultId types.VaultId
		// Vault params to set.
		vaultParams *types.VaultParams
		/* Expectations */
		// Whether quoting params should be default.
		shouldBeDefault bool
	}{
		"Default Quoting Params": {
			vaultId: constants.Vault_Clob0,
			vaultParams: &types.VaultParams{
				Status: types.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
			},
			shouldBeDefault: true,
		},
		"Custom Quoting Params": {
			vaultId:         constants.Vault_Clob1,
			vaultParams:     &constants.VaultParams,
			shouldBeDefault: false,
		},
		"Non-existent Vault Params": {
			vaultId:     constants.Vault_Clob1,
			vaultParams: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			if tc.vaultParams != nil {
				err := k.SetVaultParams(ctx, tc.vaultId, *tc.vaultParams)
				require.NoError(t, err)
				p, exists := k.GetVaultQuotingParams(ctx, tc.vaultId)
				require.True(t, exists)
				if tc.shouldBeDefault {
					require.Equal(t, types.DefaultQuotingParams(), p)
				} else {
					require.Equal(t, *tc.vaultParams.QuotingParams, p)
				}
			} else {
				_, exists := k.GetVaultQuotingParams(ctx, tc.vaultId)
				require.False(t, exists)
			}
		})
	}
}

func TestGetSetOperatorParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	// At genesis, operator defaults to gov module account.
	params := k.GetOperatorParams(ctx)
	require.Equal(
		t,
		types.OperatorParams{
			Operator: constants.GovAuthority,
		},
		params,
	)

	// Set operator to Alice.
	newParams := types.OperatorParams{
		Operator: constants.AliceAccAddress.String(),
	}
	err := k.SetOperatorParams(ctx, newParams)
	require.NoError(t, err)
	require.Equal(t, newParams, k.GetOperatorParams(ctx))

	// Set invalid operator and get.
	invalidParams := types.OperatorParams{
		Operator: "",
	}
	err = k.SetOperatorParams(ctx, invalidParams)
	require.Error(t, err)
	require.Equal(t, newParams, k.GetOperatorParams(ctx))
}

func getUpsertVaultEventsFromIndexerBlock(
	ctx sdk.Context,
	keeper *keeper.Keeper,
) []*indexerevents.UpsertVaultEventV1 {
	block := keeper.GetIndexerEventManager().ProduceBlock(ctx)
	var upsertVaultEvents []*indexerevents.UpsertVaultEventV1
	for _, event := range block.Events {
		if event.Subtype != indexerevents.SubtypeUpsertVault {
			continue
		}
		if _, ok := event.OrderingWithinBlock.(*indexer_manager.IndexerTendermintEvent_TransactionIndex); ok {
			var upsertVaultEvent indexerevents.UpsertVaultEventV1
			err := proto.Unmarshal(event.DataBytes, &upsertVaultEvent)
			if err != nil {
				panic(err)
			}
			upsertVaultEvents = append(upsertVaultEvents, &upsertVaultEvent)
		}
	}
	return upsertVaultEvents
}
