package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices"
	"github.com/stretchr/testify/require"
)

func TestMsgServerUpdateClobPair(t *testing.T) {
	tests := map[string]struct {
		authority     string
		status        types.ClobPair_Status
		setup         func(ks keepertest.ClobKeepersTestContext)
		expectedResp  *types.MsgUpdateClobPairResponse
		expectedErr   error
		expectedPanic string
	}{
		"Success": {
			authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			status:    types.ClobPair_STATUS_ACTIVE,
			setup: func(ks keepertest.ClobKeepersTestContext) {
				registry := codectypes.NewInterfaceRegistry()
				cdc := codec.NewProtoCodec(registry)
				store := prefix.NewStore(ks.Ctx.KVStore(ks.StoreKey), types.KeyPrefix(types.ClobPairKeyPrefix))
				// Write clob pair to state with clob pair id 0 and status initializing.
				clobPair := constants.ClobPair_Btc
				clobPair.Status = types.ClobPair_STATUS_INITIALIZING
				b := cdc.MustMarshal(&clobPair)
				store.Set(types.ClobPairKey(
					types.ClobPairId(constants.ClobPair_Btc.Id),
				), b)
			},
			expectedResp: &types.MsgUpdateClobPairResponse{},
		},
		"Error: unsupported status transition": {
			authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			status:    types.ClobPair_STATUS_INITIALIZING,
			setup: func(ks keepertest.ClobKeepersTestContext) {
				registry := codectypes.NewInterfaceRegistry()
				cdc := codec.NewProtoCodec(registry)
				store := prefix.NewStore(ks.Ctx.KVStore(ks.StoreKey), types.KeyPrefix(types.ClobPairKeyPrefix))
				// Write clob pair to state with clob pair id 0 and status initializing.
				clobPair := constants.ClobPair_Btc
				clobPair.Status = types.ClobPair_STATUS_ACTIVE
				b := cdc.MustMarshal(&clobPair)
				store.Set(types.ClobPairKey(
					types.ClobPairId(constants.ClobPair_Btc.Id),
				), b)
			},
			expectedErr: types.ErrInvalidClobPairStatusTransition,
		},
		"Panic: clob pair not found": {
			authority:     authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			status:        types.ClobPair_STATUS_ACTIVE,
			expectedPanic: "mustGetClobPair: ClobPair with id 0 not found",
		},
		"Error: invalid authority": {
			authority: "foobar",
			status:    types.ClobPair_STATUS_ACTIVE,
			setup: func(ks keepertest.ClobKeepersTestContext) {
				// write default btc clob pair to state
				registry := codectypes.NewInterfaceRegistry()
				cdc := codec.NewProtoCodec(registry)
				store := prefix.NewStore(ks.Ctx.KVStore(ks.StoreKey), types.KeyPrefix(types.ClobPairKeyPrefix))
				// Write clob pair to state with clob pair id 0 and status initializing.
				b := cdc.MustMarshal(&constants.ClobPair_Btc)
				store.Set(types.ClobPairKey(
					types.ClobPairId(constants.ClobPair_Btc.Id),
				), b)
			},
			expectedErr: govtypes.ErrInvalidSigner,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
			prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			if tc.setup != nil {
				tc.setup(ks)
			}

			k := ks.ClobKeeper
			msgServer := keeper.NewMsgServerImpl(k)
			wrappedCtx := sdk.WrapSDKContext(ks.Ctx)

			clobPair := constants.ClobPair_Btc
			clobPair.Status = tc.status
			msg := types.MsgUpdateClobPair{
				Authority: tc.authority,
				ClobPair:  &clobPair,
			}

			if tc.expectedPanic != "" {
				require.PanicsWithValue(t, tc.expectedPanic, func() {
					_, err := msgServer.UpdateClobPair(wrappedCtx, &msg)
					require.NoError(t, err)
				})
			} else {
				resp, err := msgServer.UpdateClobPair(wrappedCtx, &msg)
				require.Equal(t, tc.expectedResp, resp)

				if tc.expectedErr != nil {
					require.ErrorIs(t, err, tc.expectedErr)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}
}
