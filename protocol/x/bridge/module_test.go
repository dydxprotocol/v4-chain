package bridge_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	bank_keeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bridge_servertypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/bridge"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge"
	bridge_keeper "github.com/dydxprotocol/v4-chain/protocol/x/bridge/keeper"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createAppModule(t *testing.T) bridge.AppModule {
	am, _, _, _, _ := createAppModuleWithKeeper(t)
	return am
}

// Returns the keeper and context along with the AppModule.
// This is useful for tests which want to write/read state
// to/from the keeper.
func createAppModuleWithKeeper(t *testing.T) (
	bridge.AppModule,
	*bridge_keeper.Keeper,
	*bridge_servertypes.BridgeEventManager,
	*bank_keeper.BaseKeeper,
	sdk.Context,
) {
	interfaceRegistry := types.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(interfaceRegistry)

	ctx, keeper, _, _, bridgeEventManager, bankKeeper, _ := keeper.BridgeKeepers(t)

	return bridge.NewAppModule(
		appCodec,
		*keeper,
	), keeper, bridgeEventManager, bankKeeper, ctx
}

func createAppModuleBasic(t *testing.T) bridge.AppModuleBasic {
	interfaceRegistry := types.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(interfaceRegistry)

	appModule := bridge.NewAppModuleBasic(appCodec)
	require.NotNil(t, appModule)

	return appModule
}

func TestAppModuleBasic_Name(t *testing.T) {
	am := createAppModuleBasic(t)

	require.Equal(t, "bridge", am.Name())
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
	mockRegistry.AssertNumberOfCalls(t, "RegisterImplementations", 10)
	mockRegistry.AssertExpectations(t)
}

func TestAppModuleBasic_DefaultGenesis(t *testing.T) {
	am := createAppModuleBasic(t)

	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	result := am.DefaultGenesis(cdc)
	json, err := result.MarshalJSON()
	require.NoError(t, err)
	require.Equal(
		t,
		`{"event_params":{"denom":"bridge-token","eth_chain_id":"11155111",`+
			`"eth_address":"0xEf01c3A30eB57c91c40C52E996d29c202ae72193"},"propose_params":`+
			`{"max_bridges_per_block":10,"propose_delay_duration":"60s","skip_rate_ppm":800000,`+
			`"skip_if_block_delayed_by_duration":"5s"},"safety_params":{"is_disabled":false,`+
			`"delay_blocks":86400},"acknowledged_event_info":{"next_id":0,"eth_block_height":"0"}}`,
		string(json),
	)
}

func TestAppModuleBasic_ValidateGenesisErrInvalidJSON(t *testing.T) {
	am := createAppModuleBasic(t)

	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	h := json.RawMessage(`{"missingClosingQuote: true}`)

	err := am.ValidateGenesis(cdc, nil, h)
	require.EqualError(t, err, "failed to unmarshal bridge genesis state: unexpected EOF")
}

func TestAppModuleBasic_ValidateGenesisErrBadState(t *testing.T) {
	am := createAppModuleBasic(t)

	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	// bad JSON - extra { at the beginning.
	h := json.RawMessage(`{{"event_params":{"denom":"bridge-token","eth_chain_id":"11155111",
		"eth_address":"0xEf01c3A30eB57c91c40C52E996d29c202ae72193"},"propose_params":{"max_bridges_per_block":10,
		"propose_delay_duration":"60s","skip_rate_ppm":800000,"skip_if_block_delayed_by_duration":"5s"},
		"safety_params":{"is_disabled":false,"delay_blocks":86400},"acknowledged_event_info":{"next_id":0,
		"eth_block_height":"0"}}`)

	err := am.ValidateGenesis(cdc, nil, h)
	require.EqualError(t, err, `failed to unmarshal bridge genesis state: invalid character '{' `+
		`looking for beginning of object key string`)
}

func TestAppModuleBasic_ValidateGenesis(t *testing.T) {
	am := createAppModuleBasic(t)

	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	h := json.RawMessage(`{"event_params":{"denom":"bridge-token","eth_chain_id":"11155111",
		"eth_address":"0xEf01c3A30eB57c91c40C52E996d29c202ae72193"},"propose_params":{"max_bridges_per_block":10,
		"propose_delay_duration":"60s","skip_rate_ppm":800000,"skip_if_block_delayed_by_duration":"5s"},
		"safety_params":{"is_disabled":false,"delay_blocks":86400},"acknowledged_event_info":{"next_id":0,
		"eth_block_height":"0"}}`)

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

	// Expect EventParams route registered
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/dydxprotocol/v4/bridge/event_params", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect ProposeParams route registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/dydxprotocol/v4/bridge/propose_params", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect SafetyParams route registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/dydxprotocol/v4/bridge/safety_params", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect AcknowledgedEventInfo route registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/dydxprotocol/v4/bridge/acknowledged_event_info", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect RecognizedEventInfo route registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/dydxprotocol/v4/bridge/recognized_event_info", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")
}

func TestAppModuleBasic_GetTxCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetTxCmd()
	require.Equal(t, "bridge", cmd.Use)
	require.Equal(t, 0, len(cmd.Commands()))
}

func TestAppModuleBasic_GetQueryCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetQueryCmd()
	require.Equal(t, "bridge", cmd.Use)
	require.Equal(t, 6, len(cmd.Commands()))
	require.Equal(t, "get-acknowledged-event-info", cmd.Commands()[0].Name())
	require.Equal(t, "get-delayed-complete-bridge-messages", cmd.Commands()[1].Name())
	require.Equal(t, "get-event-params", cmd.Commands()[2].Name())
	require.Equal(t, "get-propose-params", cmd.Commands()[3].Name())
	require.Equal(t, "get-recognized-event-info", cmd.Commands()[4].Name())
	require.Equal(t, "get-safety-params", cmd.Commands()[5].Name())
}

func TestAppModule_Name(t *testing.T) {
	am := createAppModule(t)
	require.Equal(t, "bridge", am.Name())
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
	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	am, keeper, _, _, ctx := createAppModuleWithKeeper(t)
	msg := `{"event_params": {"denom": "bridge-token", "eth_chain_id": "77",
	"eth_address": "0xEf01c3A30eB57c91c40C52E996d29c202ae72193"}, "propose_params": {"max_bridges_per_block": 10,
	"propose_delay_duration": "60s","skip_rate_ppm": 800000, "skip_if_block_delayed_by_duration": "5s"},
	"safety_params": {"is_disabled": false,"delay_blocks": 86400}, "acknowledged_event_info": {"next_id": 0,
	"eth_block_height": "0"}}`
	gs := json.RawMessage(msg)

	result := am.InitGenesis(ctx, cdc, gs)
	require.Equal(t, 0, len(result))

	require.Equal(t, uint64(77), keeper.GetEventParams(ctx).EthChainId)
	require.Equal(t, time.Second*60, keeper.GetProposeParams(ctx).ProposeDelayDuration)
	require.Equal(t, uint32(86400), keeper.GetSafetyParams(ctx).DelayBlocks)
	require.Equal(t, uint32(0), keeper.GetAcknowledgedEventInfo(ctx).NextId)

	genesisJson := am.ExportGenesis(ctx, cdc)
	expected := `{"event_params":{"denom":"bridge-token","eth_chain_id":"77",`
	expected += `"eth_address":"0xEf01c3A30eB57c91c40C52E996d29c202ae72193"},"propose_params":{`
	expected += `"max_bridges_per_block":10,"propose_delay_duration":"60s","skip_rate_ppm":800000,`
	expected += `"skip_if_block_delayed_by_duration":"5s"},"safety_params":{"is_disabled":false,"delay_blocks":86400},`
	expected += `"acknowledged_event_info":{"next_id":0,"eth_block_height":"0"}}`
	require.Equal(t, expected, string(genesisJson))
}

func TestAppModule_InitGenesisPanic(t *testing.T) {
	am, _, _, _, ctx := createAppModuleWithKeeper(t)
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
	am := createAppModule(t)

	var ctx sdk.Context
	var req abci.RequestEndBlock
	result := am.EndBlock(ctx, req)
	require.Equal(t, 0, len(result))
}
