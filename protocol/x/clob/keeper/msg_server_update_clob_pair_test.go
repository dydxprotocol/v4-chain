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
		msg           *types.MsgUpdateClobPair
		setup         func(ks keepertest.ClobKeepersTestContext)
		expectedResp  *types.MsgUpdateClobPairResponse
		expectedErr   error
		expectedPanic string
	}{
		"Success": {
			msg: &types.MsgUpdateClobPair{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				ClobPair: types.ClobPair{
					Id: 0,
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 0,
						},
					},
					StepBaseQuantums:          5,
					SubticksPerTick:           5,
					QuantumConversionExponent: -8,
					Status:                    types.ClobPair_STATUS_ACTIVE,
				},
			},
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
		"Error: unsupported status transition from active to initializing": {
			msg: &types.MsgUpdateClobPair{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				ClobPair: types.ClobPair{
					Id: 0,
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 0,
						},
					},
					StepBaseQuantums:          5,
					SubticksPerTick:           5,
					QuantumConversionExponent: -8,
					Status:                    types.ClobPair_STATUS_INITIALIZING,
				},
			},
			setup: func(ks keepertest.ClobKeepersTestContext) {
				registry := codectypes.NewInterfaceRegistry()
				cdc := codec.NewProtoCodec(registry)
				store := prefix.NewStore(ks.Ctx.KVStore(ks.StoreKey), types.KeyPrefix(types.ClobPairKeyPrefix))
				// Write clob pair to state with clob pair id 0 and status active.
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
			msg: &types.MsgUpdateClobPair{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				ClobPair: types.ClobPair{
					Id: 0,
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 0,
						},
					},
					StepBaseQuantums:          5,
					SubticksPerTick:           5,
					QuantumConversionExponent: -8,
					Status:                    types.ClobPair_STATUS_ACTIVE,
				},
			},
			expectedPanic: "mustGetClobPair: ClobPair with id 0 not found",
		},
		"Error: invalid authority": {
			msg: &types.MsgUpdateClobPair{
				Authority: "foobar",
				ClobPair: types.ClobPair{
					Id: 0,
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 0,
						},
					},
					StepBaseQuantums:          5,
					SubticksPerTick:           5,
					QuantumConversionExponent: -8,
					Status:                    types.ClobPair_STATUS_ACTIVE,
				},
			},
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
		"Error: cannot update metadata with new perpetual id": {
			msg: &types.MsgUpdateClobPair{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				ClobPair: types.ClobPair{
					Id: 0,
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 1,
						},
					},
					StepBaseQuantums:          5,
					SubticksPerTick:           5,
					QuantumConversionExponent: -8,
					Status:                    types.ClobPair_STATUS_ACTIVE,
				},
			},
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
			expectedErr: types.ErrInvalidClobPairUpdate,
		},
		"Error: cannot update step base quantums": {
			msg: &types.MsgUpdateClobPair{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				ClobPair: types.ClobPair{
					Id: 0,
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 1,
						},
					},
					StepBaseQuantums:          10,
					SubticksPerTick:           5,
					QuantumConversionExponent: -8,
					Status:                    types.ClobPair_STATUS_ACTIVE,
				},
			},
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
			expectedErr: types.ErrInvalidClobPairUpdate,
		},
		"Error: cannot update subticks per tick": {
			msg: &types.MsgUpdateClobPair{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				ClobPair: types.ClobPair{
					Id: 0,
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 1,
						},
					},
					StepBaseQuantums:          5,
					SubticksPerTick:           10,
					QuantumConversionExponent: -8,
					Status:                    types.ClobPair_STATUS_ACTIVE,
				},
			},
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
			expectedErr: types.ErrInvalidClobPairUpdate,
		},
		"Error: cannot update quantum converstion exponent": {
			msg: &types.MsgUpdateClobPair{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				ClobPair: types.ClobPair{
					Id: 0,
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 1,
						},
					},
					StepBaseQuantums:          5,
					SubticksPerTick:           5,
					QuantumConversionExponent: -4,
					Status:                    types.ClobPair_STATUS_ACTIVE,
				},
			},
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
			expectedErr: types.ErrInvalidClobPairUpdate,
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

			if tc.expectedPanic != "" {
				require.PanicsWithValue(t, tc.expectedPanic, func() {
					_, err := msgServer.UpdateClobPair(wrappedCtx, tc.msg)
					require.NoError(t, err)
				})
			} else {
				resp, err := msgServer.UpdateClobPair(wrappedCtx, tc.msg)
				require.Equal(t, tc.expectedResp, resp)

				if tc.expectedErr != nil {
					require.ErrorIs(t, err, tc.expectedErr)
				} else {
					require.NoError(t, err)
					clobPair, found := k.GetClobPair(ks.Ctx, types.ClobPairId(tc.msg.ClobPair.Id))
					require.True(t, found)
					require.Equal(t, clobPair, tc.msg.ClobPair)
				}
			}
		})
	}
}
