package delaymsg_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed"
	testutildelaymsg "github.com/dydxprotocol/v4-chain/protocol/testutil/delaymsg"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg"
	delaymsg_keeper "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/keeper"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createAppModule(t *testing.T) delaymsg.AppModule {
	am, _, _ := createAppModuleWithKeeper(t)
	return am
}

// Returns the keeper and context along with the AppModule.
// This is useful for tests which want to write/read state
// to/from the keeper.
func createAppModuleWithKeeper(t *testing.T) (delaymsg.AppModule, *delaymsg_keeper.Keeper, sdk.Context) {
	interfaceRegistry := types.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(interfaceRegistry)

	ctx, keeper, _, _, _, _ := keeper.DelayMsgKeepers(t)

	return delaymsg.NewAppModule(
		appCodec,
		*keeper,
	), keeper, ctx
}

func createAppModuleBasic(t *testing.T) delaymsg.AppModuleBasic {
	interfaceRegistry := types.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(interfaceRegistry)

	appModule := delaymsg.NewAppModuleBasic(appCodec)
	require.NotNil(t, appModule)

	return appModule
}

func TestAppModuleBasic_Name(t *testing.T) {
	am := createAppModuleBasic(t)

	require.Equal(t, "delaymsg", am.Name())
}

func TestAppModuleBasic_RegisterCodec(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewLegacyAmino()
	am.RegisterCodec(cdc)

	var buf bytes.Buffer
	err := cdc.Amino.PrintTypes(&buf)
	require.NoError(t, err)
}

func TestAppModuleBasic_RegisterCodecLegacyAmino(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewLegacyAmino()
	am.RegisterLegacyAminoCodec(cdc)

	var buf bytes.Buffer
	err := cdc.Amino.PrintTypes(&buf)
	require.NoError(t, err)
}

func TestAppModuleBasic_RegisterInterfaces(t *testing.T) {
	am := createAppModuleBasic(t)

	mockRegistry := new(mocks.InterfaceRegistry)
	mockRegistry.On("RegisterImplementations", (*sdk.Msg)(nil), mock.Anything).Return()
	mockRegistry.On("RegisterImplementations", (*tx.MsgResponse)(nil), mock.Anything).Return()
	am.RegisterInterfaces(mockRegistry)
	mockRegistry.AssertNumberOfCalls(t, "RegisterImplementations", 2)
	mockRegistry.AssertExpectations(t)
}

func TestAppModuleBasic_DefaultGenesis(t *testing.T) {
	am := createAppModuleBasic(t)

	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	expectedGenesisJsonString := pricefeed.ReadJsonTestFile(t, "expected_default_genesis.json")

	result := am.DefaultGenesis(cdc)
	json, err := result.MarshalJSON()
	require.NoError(t, err)

	require.Equal(t, expectedGenesisJsonString, string(json))
}

func TestAppModuleBasic_ValidateGenesisErr(t *testing.T) {
	tests := map[string]struct {
		genesisJson string
		expectedErr string
	}{
		"Invalid json": {
			genesisJson: `{"missingClosingQuote: true}`,
			expectedErr: "failed to unmarshal delaymsg genesis state: unexpected EOF",
		},
		"Invalid state": {
			genesisJson: `{"next_delayed_message_id":1,` +
				`"delayed_messages":[{"id": 1,"block_height":1}]}`,
			expectedErr: "invalid delayed message at index 0 with id 1: Delayed msg is nil: Invalid genesis state",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			am := createAppModuleBasic(t)

			interfaceRegistry := types.NewInterfaceRegistry()
			cdc := codec.NewProtoCodec(interfaceRegistry)

			err := am.ValidateGenesis(cdc, nil, json.RawMessage(tc.genesisJson))
			require.EqualError(t, err, tc.expectedErr)
		})
	}
}

func TestAppModuleBasic_ValidateGenesis(t *testing.T) {
	am := createAppModuleBasic(t)

	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	bridgetypes.RegisterInterfaces(interfaceRegistry)

	validGenesisState := pricefeed.ReadJsonTestFile(t, "valid_genesis_state.json")

	h := json.RawMessage(validGenesisState)

	err := am.ValidateGenesis(cdc, nil, h)
	require.NoError(t, err)
}

func TestAppModuleBasic_RegisterRESTRoutes(t *testing.T) {
	am := createAppModuleBasic(t)

	router := mux.NewRouter()

	am.RegisterRESTRoutes(client.Context{}, router)

	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		return errors.New("No Routes Expected")
	})

	require.NoError(t, err)
}

func TestAppModuleBasic_RegisterGRPCGatewayRoutes(t *testing.T) {
	am := createAppModuleBasic(t)

	router := runtime.NewServeMux()

	am.RegisterGRPCGatewayRoutes(client.Context{}, router)

	// Expect NextDelayedMessageId route registered
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/dydxprotocol/v4/delaymsg/next_id", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect Messages route registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/dydxprotocol/v4/delaymsg/message/0", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect BlockMessageIds route registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/dydxprotocol/v4/delaymsg/block/message_ids/100", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect unexpected route not registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/dydxprotocol/v4/delaymsg/foo/bar/baz", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Equal(t, 404, recorder.Code)
}

func TestAppModuleBasic_GetTxCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetTxCmd()
	require.Equal(t, "delaymsg", cmd.Use)
	require.Equal(t, 0, len(cmd.Commands()))
}

func TestAppModuleBasic_GetQueryCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetQueryCmd()
	require.Equal(t, "delaymsg", cmd.Use)
	require.Equal(t, 3, len(cmd.Commands()))
	require.Equal(t, "get-block-message-ids", cmd.Commands()[0].Name())
	require.Equal(t, "get-message", cmd.Commands()[1].Name())
	require.Equal(t, "get-next-delayed-message-id", cmd.Commands()[2].Name())
}

func TestAppModule_Name(t *testing.T) {
	am := createAppModule(t)
	require.Equal(t, "delaymsg", am.Name())
}

func TestAppModule_RegisterServices(t *testing.T) {
	mockConfigurator := new(mocks.Configurator)
	mockQueryServer := new(mocks.Server)
	mockMsgServer := new(mocks.Server)

	mockConfigurator.On("QueryServer").Return(mockQueryServer)
	mockConfigurator.On("MsgServer").Return(mockMsgServer)
	mockQueryServer.On("RegisterService", mock.Anything, mock.Anything).Return()
	mockMsgServer.On("RegisterService", mock.Anything, mock.Anything).Return()

	am := createAppModule(t)
	am.RegisterServices(mockConfigurator)

	require.Equal(t, true, mockConfigurator.AssertExpectations(t))
	require.Equal(t, true, mockQueryServer.AssertExpectations(t))
	require.Equal(t, true, mockMsgServer.AssertExpectations(t))
}

func TestAppModule_RegisterInvariants(t *testing.T) {
	am := createAppModule(t)
	am.RegisterInvariants(nil)
}

func TestAppModule_InitExportGenesis(t *testing.T) {
	am, keeper, ctx := createAppModuleWithKeeper(t)
	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	bridgetypes.RegisterInterfaces(interfaceRegistry)

	validGenesisState := pricefeed.ReadJsonTestFile(t, "valid_genesis_state.json")

	gs := json.RawMessage(validGenesisState)

	result := am.InitGenesis(ctx, cdc, gs)
	require.Equal(t, 0, len(result))

	nextDelayedMessageId := keeper.GetNextDelayedMessageId(ctx)
	require.Equal(t, uint32(2), nextDelayedMessageId)

	delayedMessage, found := keeper.GetMessage(ctx, 1)
	require.True(t, found)
	require.Equal(t, uint32(1), delayedMessage.Id)
	require.Equal(t, uint32(100), delayedMessage.BlockHeight)
	require.Equal(t, testutildelaymsg.CreateTestAnyMsg(t), delayedMessage.Msg)

	blockIds, found := keeper.GetBlockMessageIds(ctx, 100)
	require.True(t, found)
	require.Equal(t, []uint32{1}, blockIds.Ids)

	genesisJson := am.ExportGenesis(ctx, cdc)
	require.Equal(t, validGenesisState, string(genesisJson))
}

func TestAppModule_InitGenesisPanic(t *testing.T) {
	am, _, ctx := createAppModuleWithKeeper(t)
	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	gs := json.RawMessage(`invalid json`)

	require.Panics(t, func() { am.InitGenesis(ctx, cdc, gs) })
}

func TestAppModule_ConsensusVersion(t *testing.T) {
	am := createAppModule(t)
	require.Equal(t, uint64(1), am.ConsensusVersion())
}

func TestAppModule_BeginBlock(t *testing.T) {
	am := createAppModule(t)

	var ctx sdk.Context
	var req abci.RequestBeginBlock
	am.BeginBlock(ctx, req) // should not panic
}

func TestAppModule_EndBlock(t *testing.T) {
	am, _, ctx := createAppModuleWithKeeper(t)

	var req abci.RequestEndBlock
	result := am.EndBlock(ctx, req)
	require.Equal(t, 0, len(result))
}
