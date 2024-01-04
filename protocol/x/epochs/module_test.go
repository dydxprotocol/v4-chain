package epochs_test

import (
	"bytes"
	"encoding/json"
	"github.com/dydxprotocol/v4-chain/protocol/app/module"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/epochs"
	epochs_keeper "github.com/dydxprotocol/v4-chain/protocol/x/epochs/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createAppModule(t *testing.T) epochs.AppModule {
	am, _, _ := createAppModuleWithKeeper(t)
	return am
}

// Returns the keeper and context along with the AppModule.
// This is useful for tests which want to write/read state
// to/from the keeper.
func createAppModuleWithKeeper(t *testing.T) (epochs.AppModule, *epochs_keeper.Keeper, sdk.Context) {
	appCodec := codec.NewProtoCodec(module.InterfaceRegistry)

	ctx, keeper, _ := keeper.EpochsKeeper(t)

	return epochs.NewAppModule(appCodec, *keeper), keeper, ctx
}

func createAppModuleBasic(t *testing.T) epochs.AppModuleBasic {
	appCodec := codec.NewProtoCodec(module.InterfaceRegistry)

	appModule := epochs.NewAppModuleBasic(appCodec)
	require.NotNil(t, appModule)

	return appModule
}

func TestAppModuleBasic_Name(t *testing.T) {
	am := createAppModuleBasic(t)

	require.Equal(t, "epochs", am.Name())
}

func TestAppModuleBasic_RegisterCodecLegacyAmino(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewLegacyAmino()
	am.RegisterLegacyAminoCodec(cdc)

	var buf bytes.Buffer
	err := cdc.Amino.PrintTypes(&buf)
	require.NoError(t, err)
	require.NotContains(t, buf.String(), "Msg") // epochs does not support any messages.
}

func TestAppModuleBasic_RegisterInterfaces(t *testing.T) {
	am := createAppModuleBasic(t)

	registry := codectypes.NewInterfaceRegistry()
	am.RegisterInterfaces(registry)
	// implInterfaces is a map[reflect.Type]reflect.Type that isn't exported and can't be mocked
	// due to it using an unexported method on the interface thus we use reflection to access the field
	// directly that contains the registrations.
	fv := reflect.ValueOf(registry).Elem().FieldByName("implInterfaces")
	require.Len(t, fv.MapKeys(), 0)
}

func TestAppModuleBasic_DefaultGenesis(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	result := am.DefaultGenesis(cdc)
	json, err := result.MarshalJSON()
	require.NoError(t, err)

	expectedJson := `{"epoch_info_list":`
	expectedJson += `[{"name":"funding-sample","next_tick":30,"duration":60,`
	expectedJson += `"current_epoch":0,"current_epoch_start_block":0,"is_initialized":false,`
	expectedJson += `"fast_forward_next_tick":true},{"name":"funding-tick",`
	expectedJson += `"next_tick":0,"duration":3600,"current_epoch":0,"current_epoch_start_block":0,`
	expectedJson += `"is_initialized":false,"fast_forward_next_tick":true},{"name":"stats-epoch",`
	expectedJson += `"next_tick":0,"duration":3600,"current_epoch":0,"current_epoch_start_block":0,`
	expectedJson += `"is_initialized":false,"fast_forward_next_tick":true}]}`
	require.Equal(t, expectedJson, string(json))
}

func TestAppModuleBasic_ValidateGenesisErrInvalidJSON(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	h := json.RawMessage(`{"missingClosingQuote: true}`)

	err := am.ValidateGenesis(cdc, nil, h)
	require.EqualError(t, err, "failed to unmarshal epochs genesis state: unexpected EOF")
}

func TestAppModuleBasic_ValidateGenesisErrBadState_EmptyName(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	h := json.RawMessage(`{"epoch_info_list":[{"name":""}]}`)

	err := am.ValidateGenesis(cdc, nil, h)
	require.ErrorIs(
		t,
		err,
		types.ErrEmptyEpochInfoName,
	)
}

func TestAppModuleBasic_ValidateGenesis(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	validGenesis := `{"epoch_info_list":[`
	validGenesis += `{"name":"funding-sample","next_tick":30,"duration":60,`
	validGenesis += `"current_epoch":0,"current_epoch_start_block":0,"fast_forward_next_tick":true},`
	validGenesis += `{"name":"funding-tick","next_tick":0,"duration":3600,`
	validGenesis += `"current_epoch":0,"current_epoch_start_block":0, "fast_forward_next_tick":true}]}`
	h := json.RawMessage(validGenesis)

	err := am.ValidateGenesis(cdc, nil, h)
	require.NoError(t, err)
}

func TestAppModuleBasic_RegisterGRPCGatewayRoutes(t *testing.T) {
	am := createAppModuleBasic(t)

	router := runtime.NewServeMux()

	am.RegisterGRPCGatewayRoutes(client.Context{}, router)

	// Expect list all epoch info route registered
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/dydxprotocol/v4/epochs/epoch_info", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect epoch route registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/dydxprotocol/v4/epochs/epoch_info/deewhydeeex", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect unexpected route not registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/dydxprotocol/v4/epochs/invalid/path", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Equal(t, 404, recorder.Code)
}

func TestAppModuleBasic_GetTxCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetTxCmd()
	require.Equal(t, "epochs", cmd.Use)
	require.Equal(t, 0, len(cmd.Commands()))
}

func TestAppModuleBasic_GetQueryCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetQueryCmd()
	require.Equal(t, "epochs", cmd.Use)
	require.Equal(t, 2, len(cmd.Commands()))
	require.Equal(t, "list-epoch-info", cmd.Commands()[0].Name())
	require.Equal(t, "show-epoch-info", cmd.Commands()[1].Name())
}

func TestAppModule_Name(t *testing.T) {
	am := createAppModule(t)
	require.Equal(t, "epochs", am.Name())
}

func TestAppModule_RegisterServices(t *testing.T) {
	mockConfigurator := new(mocks.Configurator)
	mockQueryServer := new(mocks.Server)
	mockMsgServer := new(mocks.Server)

	mockConfigurator.On("QueryServer").Return(mockQueryServer)
	// Since there's no MsgServer for epochs module, configurator does not call `MsgServer`.
	mockQueryServer.On("RegisterService", mock.Anything, mock.Anything).Return()
	// Since there's no MsgServer for epochs module, MsgServer does not call `RegisterServer`.

	am := createAppModule(t)
	am.RegisterServices(mockConfigurator)

	require.Equal(t, true, mockConfigurator.AssertExpectations(t))
	require.Equal(t, true, mockQueryServer.AssertExpectations(t))
	require.Equal(t, true, mockMsgServer.AssertExpectations(t))
}

func TestAppModule_InitExportGenesis(t *testing.T) {
	am, keeper, ctx := createAppModuleWithKeeper(t)
	fixedTime := time.Unix(1667293200, 0) // 2022-11-01 09:00:00 +0000 UTC
	ctxWithFixedTime := ctx.WithBlockTime(fixedTime)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)
	msg := `{"epoch_info_list":[`
	msg += `{"name":"funding-sample","next_tick":30,"duration":60,`
	msg += `"current_epoch":0,"current_epoch_start_block":0,"fast_forward_next_tick":true},`
	msg += `{"name":"funding-tick","next_tick":0,"duration":3600,`
	msg += `"current_epoch":0,"current_epoch_start_block":0,"fast_forward_next_tick":true}]}`
	gs := json.RawMessage(msg)

	am.InitGenesis(ctxWithFixedTime, cdc, gs)

	epochs := keeper.GetAllEpochInfo(ctxWithFixedTime)
	require.Equal(t, 2, len(epochs))

	require.Equal(t, "funding-sample", epochs[0].Name)
	require.NotEqual(t, 30, epochs[0].NextTick)
	require.Equal(t, uint32(60), epochs[0].Duration)
	require.Equal(t, uint32(0), epochs[0].CurrentEpoch)
	require.Equal(t, uint32(0), epochs[0].CurrentEpochStartBlock)

	require.Equal(t, "funding-tick", epochs[1].Name)
	require.NotEqual(t, 0, epochs[1].NextTick)
	require.Equal(t, uint32(3600), epochs[1].Duration)
	require.Equal(t, uint32(0), epochs[1].CurrentEpoch)
	require.Equal(t, uint32(0), epochs[1].CurrentEpochStartBlock)

	genesisJson := am.ExportGenesis(ctxWithFixedTime, cdc)
	expected := `{"epoch_info_list":[{"name":"funding-sample","next_tick":30,"duration":60,`
	expected += `"current_epoch":0,"current_epoch_start_block":0,"is_initialized":false,`
	expected += `"fast_forward_next_tick":true},{"name":"funding-tick","next_tick":0`
	expected += `,"duration":3600,"current_epoch":0,"current_epoch_start_block":0,`
	expected += `"is_initialized":false,"fast_forward_next_tick":true}]}`
	require.Equal(t, expected, string(genesisJson))
}

func TestAppModule_InitGenesisPanic(t *testing.T) {
	am, _, ctx := createAppModuleWithKeeper(t)
	cdc := codec.NewProtoCodec(module.InterfaceRegistry)
	gs := json.RawMessage(`invalid json`)

	require.Panics(t, func() { am.InitGenesis(ctx, cdc, gs) })
}

func TestAppModule_ConsensusVersion(t *testing.T) {
	am := createAppModule(t)
	require.Equal(t, uint64(1), am.ConsensusVersion())
}

func TestAppModule_BeginBlock(t *testing.T) {
	am, _, ctx := createAppModuleWithKeeper(t)

	require.NoError(t, am.BeginBlock(ctx)) // should not panic
}
