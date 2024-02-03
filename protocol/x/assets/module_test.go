package assets_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/module"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets"
	assets_keeper "github.com/dydxprotocol/v4-chain/protocol/x/assets/keeper"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createAppModule(t *testing.T) assets.AppModule {
	am, _, _ := createAppModuleWithKeeper(t)
	return am
}

// Returns the keeper and context along with the AppModule.
// This is useful for tests which want to write/read state
// to/from the keeper.
func createAppModuleWithKeeper(t *testing.T) (assets.AppModule, *assets_keeper.Keeper, sdk.Context) {
	appCodec := codec.NewProtoCodec(module.InterfaceRegistry)

	ctx, keeper, _, _, _, _ := keeper.AssetsKeepers(t, true)

	return assets.NewAppModule(
		appCodec,
		*keeper,
	), keeper, ctx
}

func createAppModuleBasic(t *testing.T) assets.AppModuleBasic {
	appCodec := codec.NewProtoCodec(module.InterfaceRegistry)

	appModule := assets.NewAppModuleBasic(appCodec)
	require.NotNil(t, appModule)

	return appModule
}

func TestAppModuleBasic_Name(t *testing.T) {
	am := createAppModuleBasic(t)

	require.Equal(t, "assets", am.Name())
}

func TestAppModuleBasic_RegisterCodecLegacyAmino(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewLegacyAmino()
	am.RegisterLegacyAminoCodec(cdc)

	var buf bytes.Buffer
	err := cdc.Amino.PrintTypes(&buf)
	require.NoError(t, err)
	require.NotContains(t, buf.String(), "Msg") // assets does not support any messages.
}

func TestAppModuleBasic_RegisterInterfaces(t *testing.T) {
	am := createAppModuleBasic(t)

	registry := types.NewInterfaceRegistry()
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

	expected := `{"assets":[{"id":0,"symbol":"USDC","denom":`
	expected += `"ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",`
	expected += `"denom_exponent":-6,"has_market":false,`
	expected += `"market_id":0,"atomic_resolution":-6}]}`
	require.Equal(t, expected, string(json))
}

func TestAppModuleBasic_ValidateGenesisErrInvalidJSON(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	h := json.RawMessage(`{"missingClosingQuote: true}`)

	err := am.ValidateGenesis(cdc, nil, h)
	require.EqualError(t, err, "failed to unmarshal assets genesis state: unexpected EOF")
}

func TestAppModuleBasic_ValidateGenesisErrBadState_NoAsset(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	h := json.RawMessage(`{"assets": []}`)

	err := am.ValidateGenesis(cdc, nil, h)
	require.EqualError(t, err, "No asset found in genesis state")
}

func TestAppModuleBasic_ValidateGenesis(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	msg := `{"assets":[{"id":0,"symbol":"USDC","denom":`
	msg += `"ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5"`
	msg += `,"denom_exponent":-6,"has_market":false,"atomic_resolution":-6}]}`
	h := json.RawMessage(msg)

	err := am.ValidateGenesis(cdc, nil, h)
	require.NoError(t, err)
}

func TestAppModuleBasic_RegisterGRPCGatewayRoutes(t *testing.T) {
	am := createAppModuleBasic(t)

	router := runtime.NewServeMux()

	am.RegisterGRPCGatewayRoutes(client.Context{}, router)

	// No query routes defined.
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/dydxprotocol/assets", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Equal(t, 404, recorder.Code)
}

func TestAppModuleBasic_GetTxCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetTxCmd()
	require.Equal(t, "assets", cmd.Use)
	require.Equal(t, 0, len(cmd.Commands()))
}

func TestAppModuleBasic_GetQueryCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetQueryCmd()
	require.Equal(t, "assets", cmd.Use)
	require.Len(t, cmd.Commands(), 0)
}

func TestAppModule_Name(t *testing.T) {
	am := createAppModule(t)
	require.Equal(t, "assets", am.Name())
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
	am, keeper, ctx := createAppModuleWithKeeper(t)
	cdc := codec.NewProtoCodec(module.InterfaceRegistry)
	msg := `{"assets":[{"id":0,"symbol":"USDC","denom":`
	msg += `"ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",`
	msg += `"denom_exponent":-6,"has_market":false,"atomic_resolution":-6}]}`
	gs := json.RawMessage(msg)

	am.InitGenesis(ctx, cdc, gs)

	assets := keeper.GetAllAssets(ctx)
	require.Equal(t, 1, len(assets))

	require.Equal(t, uint32(0), assets[0].Id)
	require.Equal(t,
		"ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",
		assets[0].Denom,
	)
	require.False(t, assets[0].HasMarket)
	require.Equal(t, uint32(0), assets[0].MarketId)
	require.Equal(t, int32(-6), assets[0].AtomicResolution)

	genesisJson := am.ExportGenesis(ctx, cdc)
	expected := `{"assets":[{"id":0,"symbol":"USDC","denom":`
	expected += `"ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",`
	expected += `"denom_exponent":-6,"has_market":false,`
	expected += `"market_id":0,"atomic_resolution":-6}]}`
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
