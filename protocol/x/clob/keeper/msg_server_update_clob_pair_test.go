package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/app/module"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/codec"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerUpdateClobPair(t *testing.T) {
	tests := map[string]struct {
		msg          *types.MsgUpdateClobPair
		setup        func(ks keepertest.ClobKeepersTestContext, mockIndexerEventManager *mocks.IndexerEventManager)
		expectedResp *types.MsgUpdateClobPairResponse
		expectedErr  error
	}{
		"Success": {
			msg: &types.MsgUpdateClobPair{
				Authority: lib.GovModuleAddress.String(),
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
			setup: func(ks keepertest.ClobKeepersTestContext, mockIndexerEventManager *mocks.IndexerEventManager) {
				cdc := codec.NewProtoCodec(module.InterfaceRegistry)
				store := prefix.NewStore(ks.Ctx.KVStore(ks.StoreKey), []byte(types.ClobPairKeyPrefix))
				// Write clob pair to state with clob pair id 0 and status initializing.
				clobPair := constants.ClobPair_Btc
				clobPair.Status = types.ClobPair_STATUS_INITIALIZING
				b := cdc.MustMarshal(&clobPair)
				store.Set(lib.Uint32ToKey(constants.ClobPair_Btc.Id), b)

				mockIndexerEventManager.On("AddTxnEvent",
					ks.Ctx,
					indexerevents.SubtypeUpdateClobPair,
					indexerevents.UpdateClobPairEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewUpdateClobPairEvent(
							clobPair.GetClobPairId(),
							types.ClobPair_STATUS_ACTIVE,
							clobPair.QuantumConversionExponent,
							types.SubticksPerTick(clobPair.GetSubticksPerTick()),
							satypes.BaseQuantums(clobPair.GetStepBaseQuantums()),
						),
					),
				).Once().Return()
			},
			expectedResp: &types.MsgUpdateClobPairResponse{},
		},
		"Error: unsupported status transition from active to initializing": {
			msg: &types.MsgUpdateClobPair{
				Authority: lib.GovModuleAddress.String(),
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
			setup: func(ks keepertest.ClobKeepersTestContext, mockIndexerEventManager *mocks.IndexerEventManager) {
				cdc := codec.NewProtoCodec(module.InterfaceRegistry)
				store := prefix.NewStore(ks.Ctx.KVStore(ks.StoreKey), []byte(types.ClobPairKeyPrefix))
				// Write clob pair to state with clob pair id 0 and status active.
				clobPair := constants.ClobPair_Btc
				clobPair.Status = types.ClobPair_STATUS_ACTIVE
				b := cdc.MustMarshal(&clobPair)
				store.Set(lib.Uint32ToKey(constants.ClobPair_Btc.Id), b)
			},
			expectedErr: types.ErrInvalidClobPairStatusTransition,
		},
		"Panic: clob pair not found": {
			msg: &types.MsgUpdateClobPair{
				Authority: lib.GovModuleAddress.String(),
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
			expectedErr: types.ErrInvalidClobPairUpdate,
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
			setup: func(ks keepertest.ClobKeepersTestContext, mockIndexerEventManager *mocks.IndexerEventManager) {
				// write default btc clob pair to state
				cdc := codec.NewProtoCodec(module.InterfaceRegistry)
				store := prefix.NewStore(ks.Ctx.KVStore(ks.StoreKey), []byte(types.ClobPairKeyPrefix))
				// Write clob pair to state with clob pair id 0 and status initializing.
				b := cdc.MustMarshal(&constants.ClobPair_Btc)
				store.Set(lib.Uint32ToKey(constants.ClobPair_Btc.Id), b)
			},
			expectedErr: govtypes.ErrInvalidSigner,
		},
		"Error: cannot update metadata with new perpetual id": {
			msg: &types.MsgUpdateClobPair{
				Authority: lib.GovModuleAddress.String(),
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
			setup: func(ks keepertest.ClobKeepersTestContext, mockIndexerEventManager *mocks.IndexerEventManager) {
				// write default btc clob pair to state
				cdc := codec.NewProtoCodec(module.InterfaceRegistry)
				store := prefix.NewStore(ks.Ctx.KVStore(ks.StoreKey), []byte(types.ClobPairKeyPrefix))
				// Write clob pair to state with clob pair id 0 and status initializing.
				b := cdc.MustMarshal(&constants.ClobPair_Btc)
				store.Set(lib.Uint32ToKey(constants.ClobPair_Btc.Id), b)
			},
			expectedErr: types.ErrInvalidClobPairUpdate,
		},
		"Error: cannot update step base quantums": {
			msg: &types.MsgUpdateClobPair{
				Authority: lib.GovModuleAddress.String(),
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
			setup: func(ks keepertest.ClobKeepersTestContext, mockIndexerEventManager *mocks.IndexerEventManager) {
				// write default btc clob pair to state
				cdc := codec.NewProtoCodec(module.InterfaceRegistry)
				store := prefix.NewStore(ks.Ctx.KVStore(ks.StoreKey), []byte(types.ClobPairKeyPrefix))
				// Write clob pair to state with clob pair id 0 and status initializing.
				b := cdc.MustMarshal(&constants.ClobPair_Btc)
				store.Set(lib.Uint32ToKey(constants.ClobPair_Btc.Id), b)
			},
			expectedErr: types.ErrInvalidClobPairUpdate,
		},
		"Error: cannot update subticks per tick": {
			msg: &types.MsgUpdateClobPair{
				Authority: lib.GovModuleAddress.String(),
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
			setup: func(ks keepertest.ClobKeepersTestContext, mockIndexerEventManager *mocks.IndexerEventManager) {
				// write default btc clob pair to state
				cdc := codec.NewProtoCodec(module.InterfaceRegistry)
				store := prefix.NewStore(ks.Ctx.KVStore(ks.StoreKey), []byte(types.ClobPairKeyPrefix))
				// Write clob pair to state with clob pair id 0 and status initializing.
				b := cdc.MustMarshal(&constants.ClobPair_Btc)
				store.Set(lib.Uint32ToKey(constants.ClobPair_Btc.Id), b)
			},
			expectedErr: types.ErrInvalidClobPairUpdate,
		},
		"Error: cannot update quantum conversion exponent": {
			msg: &types.MsgUpdateClobPair{
				Authority: lib.GovModuleAddress.String(),
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
			setup: func(ks keepertest.ClobKeepersTestContext, mockIndexerEventManager *mocks.IndexerEventManager) {
				// write default btc clob pair to state
				cdc := codec.NewProtoCodec(module.InterfaceRegistry)
				store := prefix.NewStore(ks.Ctx.KVStore(ks.StoreKey), []byte(types.ClobPairKeyPrefix))
				// Write clob pair to state with clob pair id 0 and status initializing.
				b := cdc.MustMarshal(&constants.ClobPair_Btc)
				store.Set(lib.Uint32ToKey(constants.ClobPair_Btc.Id), b)
			},
			expectedErr: types.ErrInvalidClobPairUpdate,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)

			ks.MarketMapKeeper.InitGenesis(ks.Ctx, constants.MarketMap_DefaultGenesisState)
			prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			if tc.setup != nil {
				tc.setup(ks, mockIndexerEventManager)
			}

			k := ks.ClobKeeper
			msgServer := keeper.NewMsgServerImpl(k)

			resp, err := msgServer.UpdateClobPair(ks.Ctx, tc.msg)
			require.Equal(t, tc.expectedResp, resp)

			mockIndexerEventManager.AssertExpectations(t)

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				clobPair, found := k.GetClobPair(ks.Ctx, types.ClobPairId(tc.msg.ClobPair.Id))
				require.True(t, found)
				require.Equal(t, clobPair, tc.msg.ClobPair)
			}
		})
	}
}
