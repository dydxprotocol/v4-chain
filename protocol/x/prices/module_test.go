package prices_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/app/module"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	marketmapkeeper "github.com/dydxprotocol/slinky/x/marketmap/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices"
	prices_keeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	// Exchange config json is left empty as it is not validated by the server.
	// This genesis state is formatted to export back to itself. It explicitly defines all fields using valid defaults.
	validGenesisState = `{` +
		`"market_params":[{"id":0,"pair":"DENT-USD","exponent":-1,"min_exchanges":1,"min_price_change_ppm":1,` +
		`"exchange_config_json":"{}"}],` +
		`"market_prices":[{"id":0,"exponent":-1,"price":"1"}]` +
		`}`
)

func createAppModule(t *testing.T) prices.AppModule {
	am, _, _ := createAppModuleWithKeeper(t)
	return am
}

// Returns the keeper and context along with the AppModule.
// This is useful for tests which want to write/read state
// to/from the keeper.
func createAppModuleWithKeeper(t *testing.T) (prices.AppModule, *prices_keeper.Keeper, sdk.Context) {
	appCodec := codec.NewProtoCodec(module.InterfaceRegistry)

	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	// Mock the time provider response for market creation.
	mockTimeProvider.On("Now").Return(constants.TimeT)

	return prices.NewAppModule(
		appCodec,
		*keeper,
		nil,
		nil,
		nil,
		nil,
	), keeper, ctx
}

func createAppModuleBasic(t *testing.T) prices.AppModuleBasic {
	appCodec := codec.NewProtoCodec(module.InterfaceRegistry)

	appModule := prices.NewAppModuleBasic(appCodec)
	require.NotNil(t, appModule)

	return appModule
}

func TestAppModuleBasic_Name(t *testing.T) {
	am := createAppModuleBasic(t)

	require.Equal(t, "prices", am.Name())
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
	require.Len(t, fv.MapKeys(), 6)
}

func TestAppModuleBasic_DefaultGenesis(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

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
		"Invalid Json": {
			genesisJson: `{"missingClosingQuote: true}`,
			expectedErr: "failed to unmarshal prices genesis state: unexpected EOF",
		},
		"Bad state: duplicate market param id": {
			genesisJson: `{"market_params": [` +
				`{"id":0,"pair": "DENT-USD","minPriceChangePpm":1},` +
				`{"id":0,"pair": "LINK-USD","minPriceChangePpm":1}` +
				`]}`,
			expectedErr: "duplicated market param id",
		},
		"Bad state: Invalid param": {
			genesisJson: `{"market_params": [{ "pair": "" }]}`,
			expectedErr: errorsmod.Wrap(pricestypes.ErrInvalidInput, "Pair cannot be empty").Error(),
		},
		"Bad state: Mismatch between params and prices": {
			genesisJson: `{"market_params": [{"pair": "DENT-USD","minPriceChangePpm":1}]}`,
			expectedErr: "expected the same number of market prices and market params",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			am := createAppModuleBasic(t)

			cdc := codec.NewProtoCodec(module.InterfaceRegistry)

			err := am.ValidateGenesis(cdc, nil, json.RawMessage(tc.genesisJson))
			require.EqualError(t, err, tc.expectedErr)
		})
	}
}

func TestAppModuleBasic_ValidateGenesis(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	h := json.RawMessage(validGenesisState)

	err := am.ValidateGenesis(cdc, nil, h)
	require.NoError(t, err)
}

func TestAppModuleBasic_RegisterGRPCGatewayRoutes(t *testing.T) {
	am := createAppModuleBasic(t)

	router := runtime.NewServeMux()

	am.RegisterGRPCGatewayRoutes(client.Context{}, router)

	// Expect AllMarkets route registered
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/dydxprotocol/prices/market", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect Markets route registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/dydxprotocol/prices/market/0", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect unexpected route not registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/dydxprotocol/prices/foo/bar/baz", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Equal(t, 404, recorder.Code)
}

func TestAppModuleBasic_GetTxCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetTxCmd()
	require.Equal(t, "prices", cmd.Use)
	require.Equal(t, 0, len(cmd.Commands()))
}

func TestAppModuleBasic_GetQueryCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetQueryCmd()
	require.Equal(t, "prices", cmd.Use)
	require.Equal(t, 4, len(cmd.Commands()))
	require.Equal(t, "list-market-param", cmd.Commands()[0].Name())
	require.Equal(t, "list-market-price", cmd.Commands()[1].Name())
	require.Equal(t, "show-market-param", cmd.Commands()[2].Name())
	require.Equal(t, "show-market-price", cmd.Commands()[3].Name())
}

func TestAppModule_Name(t *testing.T) {
	am := createAppModule(t)
	require.Equal(t, "prices", am.Name())
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
	gs := json.RawMessage(validGenesisState)

	// Create the market in market map
	var genState pricestypes.GenesisState
	cdc.MustUnmarshalJSON(gs, &genState)
	keepertest.CreateMarketsInMarketMapFromParams(
		t,
		ctx,
		keeper.MarketMapKeeper.(*marketmapkeeper.Keeper),
		genState.MarketParams,
	)

	am.InitGenesis(ctx, cdc, gs)

	marketParams := keeper.GetAllMarketParams(ctx)
	require.Equal(t, 1, len(marketParams))

	require.Equal(t, "DENT-USD", marketParams[0].Pair)
	require.Equal(t, uint32(0), marketParams[0].Id)

	genesisJson := am.ExportGenesis(ctx, cdc)
	require.Equal(t, validGenesisState, string(genesisJson))
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
