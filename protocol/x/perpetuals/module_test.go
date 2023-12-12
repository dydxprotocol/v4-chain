package perpetuals_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	pricetypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil_json "github.com/dydxprotocol/v4-chain/protocol/testutil/json"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	epochs_keeper "github.com/dydxprotocol/v4-chain/protocol/x/epochs/keeper"
	epoch_types "github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals"
	perpetuals_keeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	prices_keeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createAppModule(t *testing.T) perpetuals.AppModule {
	am, _, _, _, _ := createAppModuleWithKeeper(t)
	return am
}

// Returns the keeper and context along with the AppModule.
// This is useful for tests which want to write/read state
// to/from the keeper.
func createAppModuleWithKeeper(t *testing.T) (
	perpetuals.AppModule,
	*perpetuals_keeper.Keeper,
	*prices_keeper.Keeper,
	*epochs_keeper.Keeper,
	sdk.Context,
) {
	interfaceRegistry := types.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(interfaceRegistry)

	pc := keeper.PerpetualsKeepers(t)

	return perpetuals.NewAppModule(
		appCodec,
		pc.PerpetualsKeeper,
		nil,
		nil,
	), pc.PerpetualsKeeper, pc.PricesKeeper, pc.EpochsKeeper, pc.Ctx
}

func createAppModuleBasic(t *testing.T) perpetuals.AppModuleBasic {
	interfaceRegistry := types.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(interfaceRegistry)

	appModule := perpetuals.NewAppModuleBasic(appCodec)
	require.NotNil(t, appModule)

	return appModule
}

func TestAppModuleBasic_Name(t *testing.T) {
	am := createAppModuleBasic(t)

	require.Equal(t, "perpetuals", am.Name())
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
		`{"perpetuals":[],"liquidity_tiers":[],"params":{"funding_rate_clamp_factor_ppm":6000000,`+
			`"premium_vote_clamp_factor_ppm":60000000,"min_num_votes_per_sample":15}}`,
		string(json),
	)
}

func TestAppModuleBasic_ValidateGenesisErrInvalidJSON(t *testing.T) {
	am := createAppModuleBasic(t)

	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	h := json.RawMessage(`{"missingClosingQuote: true}`)

	err := am.ValidateGenesis(cdc, nil, h)
	require.EqualError(t, err, "failed to unmarshal perpetuals genesis state: unexpected EOF")
}

func TestAppModuleBasic_ValidateGenesisErrBadState(t *testing.T) {
	am := createAppModuleBasic(t)

	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	h := json.RawMessage(`{
		"perpetuals":[
		   {
			  "params":{
				"ticker":""
			  }
		   }
		],
		"params":{
		   "funding_rate_clamp_factor_ppm":6000000,
		   "premium_vote_clamp_factor_ppm":60000000,
		   "min_num_votes_per_sample":15
		}
	 }`)

	err := am.ValidateGenesis(cdc, nil, h)
	require.EqualError(t, err, "Ticker must be non-empty string")
}

func TestAppModuleBasic_ValidateGenesis(t *testing.T) {
	am := createAppModuleBasic(t)

	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	h := json.RawMessage(`{
		"perpetuals":[
		   {
			   "params":{
				    "ticker":"EXAM-USD",
				    "market_id":0
			  }
		   }
		],
		"params":{
		   "funding_rate_clamp_factor_ppm":6000000,
		   "premium_vote_clamp_factor_ppm":60000000,
		   "min_num_votes_per_sample":15
		}
	 }`)

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

	// Expect AllPerpetuals route registered
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/dydxprotocol/perpetuals/perpetual", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect Markets route registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/dydxprotocol/perpetuals/perpetual/0", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect unexpected route not registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/dydxprotocol/perpetuals/foo/bar/baz", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Equal(t, 404, recorder.Code)
}

func TestAppModuleBasic_GetTxCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetTxCmd()
	require.Equal(t, "perpetuals", cmd.Use)
	require.Equal(t, 0, len(cmd.Commands()))
}

func TestAppModuleBasic_GetQueryCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetQueryCmd()
	require.Equal(t, "perpetuals", cmd.Use)
	require.Equal(t, 5, len(cmd.Commands()))
	require.Equal(t, "get-params", cmd.Commands()[0].Name())
	require.Equal(t, "get-premium-samples", cmd.Commands()[1].Name())
	require.Equal(t, "get-premium-votes", cmd.Commands()[2].Name())
	require.Equal(t, "list-perpetual", cmd.Commands()[3].Name())
	require.Equal(t, "show-perpetual", cmd.Commands()[4].Name())
}

func TestAppModule_Name(t *testing.T) {
	am := createAppModule(t)
	require.Equal(t, "perpetuals", am.Name())
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

	// The corresponding `Market` must exist, so create it.
	am, keeper, pricesKeeper, _, ctx := createAppModuleWithKeeper(t)
	if _, err := pricesKeeper.CreateMarket(
		ctx,
		pricetypes.MarketParam{
			Id:                 0,
			Pair:               constants.EthUsdPair,
			Exponent:           -2,
			MinExchanges:       1,
			MinPriceChangePpm:  1_000,
			ExchangeConfigJson: "{}",
		},
		pricetypes.MarketPrice{
			Id:       0,
			Exponent: -2,
			Price:    1_000,
		},
	); err != nil {
		t.Errorf("failed to create a market %s", err)
	}

	msg := `{
		"perpetuals":[
		   {
			  "params": {
				 "ticker":"EXAM-USD",
				 "market_id":0,
				 "liquidity_tier":0
			  }
		   }
		],
		"liquidity_tiers":[
		   {
			  "name":"Large-Cap",
			  "initial_margin_ppm":50000,
			  "maintenance_fraction_ppm":500000,
			  "impact_notional":10000000000
		   }
		],
		"params":{
		   "funding_rate_clamp_factor_ppm":6000000,
		   "premium_vote_clamp_factor_ppm":60000000,
		   "min_num_votes_per_sample":15
		}
	}`
	gs := json.RawMessage(msg)

	result := am.InitGenesis(ctx, cdc, gs)
	require.Equal(t, 0, len(result))

	perpetuals := keeper.GetAllPerpetuals(ctx)
	require.Equal(t, 1, len(perpetuals))

	require.Equal(t, "EXAM-USD", perpetuals[0].Params.Ticker)
	require.Equal(t, uint32(0), perpetuals[0].Params.Id)

	genesisJson := am.ExportGenesis(ctx, cdc)
	expected := `{
		"perpetuals":[
		   {
			  "params":{
				 "id":0,
				 "ticker":"EXAM-USD",
				 "market_id":0,
				 "atomic_resolution":0,
				 "default_funding_ppm":0,
				 "liquidity_tier":0
			  },
			  "funding_index":"0"
		   }
		],
		"liquidity_tiers":[
		   {
			  "id":0,
			  "name":"Large-Cap",
			  "initial_margin_ppm":50000,
			  "maintenance_fraction_ppm":500000,
			  "base_position_notional":"0",
			  "impact_notional":"10000000000"
		   }
		],
		"params":{
		   "funding_rate_clamp_factor_ppm":6000000,
		   "premium_vote_clamp_factor_ppm":60000000,
		   "min_num_votes_per_sample":15
		}
	 }`
	require.Equal(t,
		testutil_json.CompactJsonString(t, expected),
		testutil_json.CompactJsonString(t, string(genesisJson)),
	)
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
	am, perpKeeper, _, epochsKeeper, ctx := createAppModuleWithKeeper(t)

	// Initialize empty samples in storage.
	perpKeeper.SetEmptyPremiumSamples(ctx)

	for _, epochInfo := range epoch_types.DefaultGenesis().EpochInfoList {
		if err := epochsKeeper.CreateEpochInfo(ctx, epochInfo); err != nil {
			t.Errorf("failed to create an epoch %s", err)
		}
	}

	var req abci.RequestEndBlock
	result := am.EndBlock(ctx, req)
	require.Equal(t, 0, len(result))
}
