package subaccounts_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sample"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts"
	sa_keeper "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createAppModule(t *testing.T) subaccounts.AppModule {
	am, _, _ := createAppModuleWithKeeper(t)
	return am
}

// Returns the keeper and context along with the AppModule.
// This is useful for tests which want to write/read state
// to/from the keeper.
func createAppModuleWithKeeper(t *testing.T) (subaccounts.AppModule, *sa_keeper.Keeper, sdk.Context) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(interfaceRegistry)

	ctx, keeper, _, _, _, _, _, _ := keeper.SubaccountsKeepers(t, true)

	return subaccounts.NewAppModule(
		appCodec,
		*keeper,
	), keeper, ctx
}

func createAppModuleBasic(t *testing.T) subaccounts.AppModuleBasic {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(interfaceRegistry)

	appModule := subaccounts.NewAppModuleBasic(appCodec)
	require.NotNil(t, appModule)

	return appModule
}

func TestAppModuleBasic_Name(t *testing.T) {
	am := createAppModuleBasic(t)

	require.Equal(t, "subaccounts", am.Name())
}

func TestAppModuleBasic_RegisterCodec(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewLegacyAmino()
	am.RegisterCodec(cdc)

	var buf bytes.Buffer
	err := cdc.Amino.PrintTypes(&buf)
	require.NoError(t, err)
	require.NotContains(t, buf.String(), "Msg") // subaccounts does not support any messages.
}

func TestAppModuleBasic_RegisterCodecLegacyAmino(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewLegacyAmino()
	am.RegisterLegacyAminoCodec(cdc)

	var buf bytes.Buffer
	err := cdc.Amino.PrintTypes(&buf)
	require.NoError(t, err)
	require.NotContains(t, buf.String(), "Msg") // subaccounts does not support any messages.
}

func TestAppModuleBasic_RegisterInterfaces(t *testing.T) {
	am := createAppModuleBasic(t)

	mockRegistry := new(mocks.InterfaceRegistry)
	am.RegisterInterfaces(mockRegistry)
	mockRegistry.AssertNumberOfCalls(t, "RegisterImplementations", 0)
	mockRegistry.AssertExpectations(t)
}

func TestAppModuleBasic_DefaultGenesis(t *testing.T) {
	am := createAppModuleBasic(t)

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	result := am.DefaultGenesis(cdc)
	json, err := result.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `{"subaccounts":[]}`, string(json))
}

func TestAppModuleBasic_ValidateGenesisErrInvalidJSON(t *testing.T) {
	am := createAppModuleBasic(t)

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	h := json.RawMessage(`{"missingClosingQuote: true}`)

	err := am.ValidateGenesis(cdc, nil, h)
	require.EqualError(t, err, "failed to unmarshal subaccounts genesis state: unexpected EOF")
}

func TestAppModuleBasic_ValidateGenesisErrBadState_OwnerEmpty(t *testing.T) {
	am := createAppModuleBasic(t)

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	h := json.RawMessage(`{"subaccounts": [{ "id": {"owner": "" } }]}`)

	err := am.ValidateGenesis(cdc, nil, h)
	require.ErrorContains(t, err, "invalid SubaccountId Owner")
}

func TestAppModuleBasic_ValidateGenesisErrBadState_OwnerInvalid(t *testing.T) {
	am := createAppModuleBasic(t)

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	h := json.RawMessage(`{"subaccounts": [{ "id": {"owner": "invalid" } }]}`)

	err := am.ValidateGenesis(cdc, nil, h)
	require.ErrorContains(t, err, "invalid SubaccountId Owner")
}

func TestAppModuleBasic_ValidateGenesisErrBadState_Number(t *testing.T) {
	am := createAppModuleBasic(t)

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	msg := fmt.Sprintf(`{"subaccounts": [{ "id": {"owner": "%s", "number": 128 } }]}`, sample.AccAddress())
	h := json.RawMessage(msg)

	err := am.ValidateGenesis(cdc, nil, h)
	require.EqualError(t, err, "subaccount id number cannot exceed 127")
}

func TestAppModuleBasic_ValidateGenesis(t *testing.T) {
	am := createAppModuleBasic(t)

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	msg := fmt.Sprintf(`{"subaccounts": [{ "id": {"owner": "%s", "number": 127 } }]}`, sample.AccAddress())
	h := json.RawMessage(msg)

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

	// Expect SubaccountAll route registered
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/dydxprotocol/subaccounts/subaccount", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect Subaccount route registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/dydxprotocol/subaccounts/subaccount/foo/127", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect unexpected route not registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/dydxprotocol/subaccounts/foo/bar/baz", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Equal(t, 404, recorder.Code)
}

func TestAppModuleBasic_GetTxCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetTxCmd()
	require.Equal(t, "subaccounts", cmd.Use)
	require.Equal(t, 0, len(cmd.Commands()))
}

func TestAppModuleBasic_GetQueryCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetQueryCmd()
	require.Equal(t, "subaccounts", cmd.Use)
	require.Equal(t, 2, len(cmd.Commands()))
	require.Equal(t, "list-subaccount", cmd.Commands()[0].Name())
	require.Equal(t, "show-subaccount", cmd.Commands()[1].Name())
}

func TestAppModule_Name(t *testing.T) {
	am := createAppModule(t)
	require.Equal(t, "subaccounts", am.Name())
}

func TestAppModule_RegisterServices(t *testing.T) {
	mockConfigurator := new(mocks.Configurator)
	mockQueryServer := new(mocks.Server)
	mockMsgServer := new(mocks.Server)

	mockConfigurator.On("QueryServer").Return(mockQueryServer)
	// Since there's no MsgServer for Subaccounts module, configurator does not call `MsgServer`.
	mockQueryServer.On("RegisterService", mock.Anything, mock.Anything).Return()
	// Since there's no MsgServer for Subaccounts module, MsgServer does not call `RegisterServer`.

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
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	msg := `{"subaccounts": [{ "id": {"owner": "foo", "number": 127 },`
	msg += `"asset_positions":[{"asset_id": 0, "index": 0, "quantums": "1000" }] }]}`
	gs := json.RawMessage(msg)

	result := am.InitGenesis(ctx, cdc, gs)
	require.Equal(t, 0, len(result))

	subaccounts := keeper.GetAllSubaccount(ctx)
	require.Equal(t, 1, len(subaccounts))

	require.Equal(t, "foo", subaccounts[0].Id.Owner)
	require.Equal(t, uint32(127), subaccounts[0].Id.Number)

	genesisJson := am.ExportGenesis(ctx, cdc)
	expected := `{"subaccounts":[{"id":{"owner":"foo","number":127},`
	expected += `"asset_positions":[{"asset_id":0,"quantums":"1000","index":"0"}],`
	expected += `"perpetual_positions":[],"margin_enabled":false}]}`
	require.Equal(t, expected, string(genesisJson))
}

func TestAppModule_InitGenesisPanic(t *testing.T) {
	am, _, ctx := createAppModuleWithKeeper(t)
	interfaceRegistry := codectypes.NewInterfaceRegistry()
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
