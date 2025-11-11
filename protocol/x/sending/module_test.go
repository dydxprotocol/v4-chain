package sending_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/app/module"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending"
	sending_keeper "github.com/dydxprotocol/v4-chain/protocol/x/sending/keeper"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createAppModule(t *testing.T) sending.AppModule {
	am, _, _ := createAppModuleWithKeeper(t)
	return am
}

// Returns the keeper and context along with the AppModule.
// This is useful for tests which want to write/read state
// to/from the keeper.
func createAppModuleWithKeeper(t *testing.T) (sending.AppModule, *sending_keeper.Keeper, sdk.Context) {
	appCodec := codec.NewProtoCodec(module.InterfaceRegistry)

	ks := keeper.SendingKeepers(t)

	return sending.NewAppModule(
		appCodec,
		*ks.SendingKeeper,
		nil,
		nil,
		ks.SubaccountsKeeper,
	), ks.SendingKeeper, ks.Ctx
}

func createAppModuleBasic(t *testing.T) sending.AppModuleBasic {
	appCodec := codec.NewProtoCodec(module.InterfaceRegistry)

	appModule := sending.NewAppModuleBasic(appCodec)
	require.NotNil(t, appModule)

	return appModule
}

func TestAppModuleBasic_Name(t *testing.T) {
	am := createAppModuleBasic(t)

	require.Equal(t, "sending", am.Name())
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

	registry := types.NewInterfaceRegistry()
	am.RegisterInterfaces(registry)
	// implInterfaces is a map[reflect.Type]reflect.Type that isn't exported and can't be mocked
	// due to it using an unexported method on the interface thus we use reflection to access the field
	// directly that contains the registrations.
	fv := reflect.ValueOf(registry).Elem().FieldByName("implInterfaces")
	require.Len(t, fv.MapKeys(), 10)
}

func TestAppModuleBasic_DefaultGenesis(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	result := am.DefaultGenesis(cdc)
	json, err := result.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `{}`, string(json))
}

func TestAppModuleBasic_ValidateGenesisErrInvalidJSON(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	h := json.RawMessage(`{"missingClosingQuote: true}`)

	err := am.ValidateGenesis(cdc, nil, h)
	require.EqualError(t, err, "failed to unmarshal sending genesis state: unexpected EOF")
}

func TestAppModuleBasic_ValidateGenesis(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	h := json.RawMessage(`{}`)

	err := am.ValidateGenesis(cdc, nil, h)
	require.NoError(t, err)
}

func TestAppModuleBasic_RegisterGRPCGatewayRoutes(t *testing.T) {
	am := createAppModuleBasic(t)

	router := runtime.NewServeMux()

	am.RegisterGRPCGatewayRoutes(client.Context{}, router)

	// Sending currently does not have any routes registered, so expect no routes.
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/dydxprotocol/sending/", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Equal(t, 404, recorder.Code)
}

func TestAppModuleBasic_GetTxCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetTxCmd()
	require.Equal(t, "sending", cmd.Use)
	require.Equal(t, 3, len(cmd.Commands()))
	require.Equal(t, "create-transfer", cmd.Commands()[0].Name())
	require.Equal(t, "deposit-to-subaccount", cmd.Commands()[1].Name())
	require.Equal(t, "withdraw-from-subaccount", cmd.Commands()[2].Name())
}

func TestAppModuleBasic_GetQueryCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetQueryCmd()
	require.Equal(t, "sending", cmd.Use)
	require.Equal(t, 0, len(cmd.Commands()))
}

func TestAppModule_Name(t *testing.T) {
	am := createAppModule(t)
	require.Equal(t, "sending", am.Name())
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

func TestAppModule_InitExportGenesis(t *testing.T) {
	am, _, ctx := createAppModuleWithKeeper(t)
	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	msg := `{}`
	gs := json.RawMessage(msg)

	am.InitGenesis(ctx, cdc, gs)

	genesisJson := am.ExportGenesis(ctx, cdc)
	expected := `{}`
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
