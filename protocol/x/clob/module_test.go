package clob_test

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/dydxprotocol/v4-chain/protocol/app/module"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	marketmap_keeper "github.com/dydxprotocol/slinky/x/marketmap/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob"
	clob_keeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	clob_types "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals"
	perp_keeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices"
	prices_keeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func getValidGenesisStr() string {
	gs := `{"clob_pairs":[{"id":0,"perpetual_clob_metadata":{"perpetual_id":0},"subticks_per_tick":100,`
	gs += `"step_base_quantums":5,"status":"STATUS_ACTIVE"}],`
	gs += `"liquidations_config":{`
	gs += `"max_liquidation_fee_ppm":5000,"position_block_limits":{"min_position_notional_liquidated":"1000",`
	gs += `"max_position_portion_liquidated_ppm":1000000},"subaccount_block_limits":`
	gs += `{"max_notional_liquidated":"100000000000000","max_quantums_insurance_lost":"100000000000000"},`
	gs += `"fillable_price_config":{"bankruptcy_adjustment_ppm":1000000,`
	gs += `"spread_to_maintenance_margin_ratio_ppm":100000}},"block_rate_limit_config":`
	gs += `{"max_short_term_orders_and_cancels_per_n_blocks":[{"limit": 400,"num_blocks":1}],`
	gs += `"max_stateful_orders_per_n_blocks":[{"limit": 2,"num_blocks":1},{"limit": 20,"num_blocks":100}]},`
	gs += `"equity_tier_limit_config":{"short_term_order_equity_tiers":[{"limit":0,"usd_tnc_required":"0"},`
	gs += `{"limit":1,"usd_tnc_required":"20"},{"limit":5,"usd_tnc_required":"100"},`
	gs += `{"limit":10,"usd_tnc_required":"1000"},{"limit":100,"usd_tnc_required":"10000"},`
	gs += `{"limit":1000,"usd_tnc_required":"100000"}],"stateful_order_equity_tiers":[`
	gs += `{"limit":0,"usd_tnc_required":"0"},{"limit":1,"usd_tnc_required":"20"},`
	gs += `{"limit":5,"usd_tnc_required":"100"},{"limit":10,"usd_tnc_required":"1000"},`
	gs += `{"limit":100,"usd_tnc_required":"10000"},{"limit":200,"usd_tnc_required":"100000"}]}}`
	return gs
}

func createAppModule(t *testing.T) clob.AppModule {
	am, _, _, _, _, _, _ := createAppModuleWithKeeper(t)
	return am
}

// Returns the keeper and context along with the AppModule.
// This is useful for tests which want to write/read state
// to/from the keeper.
func createAppModuleWithKeeper(t *testing.T) (
	clob.AppModule,
	*clob_keeper.Keeper,
	*prices_keeper.Keeper,
	*perp_keeper.Keeper,
	*marketmap_keeper.Keeper,
	sdk.Context,
	*mocks.IndexerEventManager,
) {
	appCodec := codec.NewProtoCodec(module.InterfaceRegistry)

	memClob := memclob.NewMemClobPriceTimePriority(false)
	mockIndexerEventManager := &mocks.IndexerEventManager{}

	mockBankKeeper := &mocks.BankKeeper{}
	mockBankKeeper.On(
		"GetBalance",
		mock.Anything,
		perptypes.InsuranceFundModuleAddress,
		constants.Usdc.Denom,
	).Return(
		sdk.NewCoin(constants.Usdc.Denom, sdkmath.NewIntFromBigInt(new(big.Int))),
	)
	ks := keeper.NewClobKeepersTestContext(t, memClob, mockBankKeeper, mockIndexerEventManager)

	err := keeper.CreateUsdcAsset(ks.Ctx, ks.AssetsKeeper)
	require.NoError(t, err)

	return clob.NewAppModule(
		appCodec,
		ks.ClobKeeper,
		nil,
		nil,
		nil,
	), ks.ClobKeeper, ks.PricesKeeper, ks.PerpetualsKeeper, ks.MarketMapKeeper, ks.Ctx, mockIndexerEventManager
}

func createAppModuleBasic(t *testing.T) clob.AppModuleBasic {
	appCodec := codec.NewProtoCodec(module.InterfaceRegistry)

	appModule := clob.NewAppModuleBasic(appCodec)
	require.NotNil(t, appModule)

	return appModule
}

func TestAppModuleBasic_Name(t *testing.T) {
	am := createAppModuleBasic(t)

	require.Equal(t, "clob", am.Name())
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

	registry := types.NewInterfaceRegistry()
	am.RegisterInterfaces(registry)
	// implInterfaces is a map[reflect.Type]reflect.Type that isn't exported and can't be mocked
	// due to it using an unexported method on the interface thus we use reflection to access the field
	// directly that contains the registrations.
	fv := reflect.ValueOf(registry).Elem().FieldByName("implInterfaces")
	require.Len(t, fv.MapKeys(), 20)
}

func TestAppModuleBasic_DefaultGenesis(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	result := am.DefaultGenesis(cdc)
	json, err := result.MarshalJSON()
	require.NoError(t, err)

	expected := `{"clob_pairs":[],"liquidations_config":{`
	expected += `"max_liquidation_fee_ppm":5000,"position_block_limits":{"min_position_notional_liquidated":"1000",`
	expected += `"max_position_portion_liquidated_ppm":1000000},"subaccount_block_limits":`
	expected += `{"max_notional_liquidated":"100000000000000","max_quantums_insurance_lost":"100000000000000"},`
	expected += `"fillable_price_config":{"bankruptcy_adjustment_ppm":1000000,`
	expected += `"spread_to_maintenance_margin_ratio_ppm":100000}},"block_rate_limit_config":`
	expected += `{"max_leverage_updates_per_n_blocks":[],`
	expected += `"max_short_term_orders_and_cancels_per_n_blocks":[],"max_stateful_orders_per_n_blocks":[],`
	expected += `"max_short_term_order_cancellations_per_n_blocks":[],"max_short_term_orders_per_n_blocks":[]},`
	expected += `"equity_tier_limit_config":{"short_term_order_equity_tiers":[], "stateful_order_equity_tiers":[]}}`

	require.JSONEq(t, expected, string(json))
}

func TestAppModuleBasic_ValidateGenesisErrInvalidJSON(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	h := json.RawMessage(`{"missingClosingQuote: true}`)

	err := am.ValidateGenesis(cdc, nil, h)
	require.EqualError(t, err, "failed to unmarshal clob genesis state: unexpected EOF")
}

func TestAppModuleBasic_ValidateGenesisErrBadState(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	h := json.RawMessage(`{"clob_pairs":[{"id":0},{"id":0}]}`)

	err := am.ValidateGenesis(cdc, nil, h)
	require.EqualError(t, err, "duplicated id for clobPair")
}

func TestAppModuleBasic_ValidateGenesis(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	h := json.RawMessage(getValidGenesisStr())

	err := am.ValidateGenesis(cdc, nil, h)
	require.NoError(t, err)
}

func TestAppModuleBasic_RegisterGRPCGatewayRoutes(t *testing.T) {
	am := createAppModuleBasic(t)

	router := runtime.NewServeMux()

	am.RegisterGRPCGatewayRoutes(client.Context{}, router)

	// Expect AllClobPairs route registered
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/dydxprotocol/clob/clob_pair", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect ClobPair route registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/dydxprotocol/clob/clob_pair/0", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Contains(t, recorder.Body.String(), "no RPC client is defined in offline mode")

	// Expect unexpected route not registered
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/clob/foo/bar/baz", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Equal(t, 404, recorder.Code)
}

func TestAppModuleBasic_GetTxCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetTxCmd()
	require.Equal(t, "clob", cmd.Use)
	require.Equal(t, 4, len(cmd.Commands()))
	require.Equal(t, "batch-cancel", cmd.Commands()[0].Name())
	require.Equal(t, "cancel-order", cmd.Commands()[1].Name())
	require.Equal(t, "place-order", cmd.Commands()[2].Name())
	require.Equal(t, "update-leverage", cmd.Commands()[3].Name())
}

func TestAppModuleBasic_GetQueryCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetQueryCmd()
	require.Equal(t, "clob", cmd.Use)
	require.Equal(t, 7, len(cmd.Commands()))
	require.Equal(t, "get-block-rate-limit-config", cmd.Commands()[0].Name())
	require.Equal(t, "get-equity-tier-limit-config", cmd.Commands()[1].Name())
	require.Equal(t, "get-liquidations-config", cmd.Commands()[2].Name())
	require.Equal(t, "leverage", cmd.Commands()[3].Name())
	require.Equal(t, "list-clob-pair", cmd.Commands()[4].Name())
	require.Equal(t, "show-clob-pair", cmd.Commands()[5].Name())
	require.Equal(t, "stateful-order", cmd.Commands()[6].Name())
}

func TestAppModule_Name(t *testing.T) {
	am := createAppModule(t)
	require.Equal(t, "clob", am.Name())
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
	am, keeper, pricesKeeper, perpetualsKeeper, marketMapKeeper, ctx, mockIndexerEventManager := createAppModuleWithKeeper(t) //nolint:lll
	ctx = ctx.WithBlockTime(constants.TimeT)
	cdc := codec.NewProtoCodec(module.InterfaceRegistry)
	gs := json.RawMessage(getValidGenesisStr())

	// PerpetualMarketCreateEvents are emitted when initializing the genesis state, so we need to mock
	// the indexer event manager to expect these events.
	mockIndexerEventManager.On("AddTxnEvent",
		ctx,
		indexerevents.SubtypePerpetualMarket,
		indexerevents.PerpetualMarketEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewPerpetualMarketCreateEvent(
				uint32(0),
				uint32(0),
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.Ticker,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.MarketId,
				clob_types.ClobPair_STATUS_ACTIVE,
				0,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.AtomicResolution,
				uint32(100),
				uint64(5),
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.LiquidityTier,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.MarketType,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.DefaultFundingPpm,
			),
		),
	).Once().Return()

	marketMapKeeper.InitGenesis(ctx, constants.MarketMap_DefaultGenesisState)
	prices.InitGenesis(ctx, *pricesKeeper, constants.Prices_DefaultGenesisState)
	perpetuals.InitGenesis(ctx, *perpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

	result := am.InitGenesis(ctx, cdc, gs)
	require.Equal(t, 0, len(result))

	clobPairs := keeper.GetAllClobPairs(ctx)
	require.Equal(t, 1, len(clobPairs))
	require.Equal(t, uint32(0), clobPairs[0].Id)
	require.Equal(t, uint32(0), clobPairs[0].GetPerpetualClobMetadata().PerpetualId)
	require.Equal(t, uint32(100), clobPairs[0].SubticksPerTick)
	require.Equal(t, uint64(5), clobPairs[0].StepBaseQuantums)
	require.Equal(t, clob_types.ClobPair_STATUS_ACTIVE, clobPairs[0].Status)

	liquidationsConfig := keeper.GetLiquidationsConfig(ctx)
	require.Equal(t, uint32(5_000), liquidationsConfig.MaxLiquidationFeePpm)
	require.Equal(t, uint32(1_000_000), liquidationsConfig.FillablePriceConfig.BankruptcyAdjustmentPpm)
	require.Equal(t, uint32(100_000), liquidationsConfig.FillablePriceConfig.SpreadToMaintenanceMarginRatioPpm)
	require.Equal(t, uint64(1_000), liquidationsConfig.PositionBlockLimits.MinPositionNotionalLiquidated)
	require.Equal(t, uint32(1_000_000), liquidationsConfig.PositionBlockLimits.MaxPositionPortionLiquidatedPpm)
	require.Equal(t, uint64(100_000_000_000_000), liquidationsConfig.SubaccountBlockLimits.MaxNotionalLiquidated)
	require.Equal(t, uint64(100_000_000_000_000), liquidationsConfig.SubaccountBlockLimits.MaxQuantumsInsuranceLost)

	blockRateLimitConfig := keeper.GetBlockRateLimitConfiguration(ctx)
	require.Equal(
		t,
		clob_types.BlockRateLimitConfiguration{
			MaxShortTermOrdersAndCancelsPerNBlocks: []clob_types.MaxPerNBlocksRateLimit{
				{
					Limit:     400,
					NumBlocks: 1,
				},
			},
			MaxStatefulOrdersPerNBlocks: []clob_types.MaxPerNBlocksRateLimit{
				{
					Limit:     2,
					NumBlocks: 1,
				},
				{
					Limit:     20,
					NumBlocks: 100,
				},
			},
		},
		blockRateLimitConfig,
	)

	equityTierLimitConfig := keeper.GetEquityTierLimitConfiguration(ctx)
	require.Equal(
		t,
		clob_types.EquityTierLimitConfiguration{
			ShortTermOrderEquityTiers: []clob_types.EquityTierLimit{
				{
					UsdTncRequired: dtypes.NewInt(0),
					Limit:          0,
				},
				{
					UsdTncRequired: dtypes.NewInt(20),
					Limit:          1,
				},
				{
					UsdTncRequired: dtypes.NewInt(100),
					Limit:          5,
				},
				{
					UsdTncRequired: dtypes.NewInt(1000),
					Limit:          10,
				},
				{
					UsdTncRequired: dtypes.NewInt(10000),
					Limit:          100,
				},
				{
					UsdTncRequired: dtypes.NewInt(100000),
					Limit:          1000,
				},
			},
			StatefulOrderEquityTiers: []clob_types.EquityTierLimit{
				{
					UsdTncRequired: dtypes.NewInt(0),
					Limit:          0,
				},
				{
					UsdTncRequired: dtypes.NewInt(20),
					Limit:          1,
				},
				{
					UsdTncRequired: dtypes.NewInt(100),
					Limit:          5,
				},
				{
					UsdTncRequired: dtypes.NewInt(1000),
					Limit:          10,
				},
				{
					UsdTncRequired: dtypes.NewInt(10000),
					Limit:          100,
				},
				{
					UsdTncRequired: dtypes.NewInt(100000),
					Limit:          200,
				},
			},
		},
		equityTierLimitConfig,
	)

	genesisJson := am.ExportGenesis(ctx, cdc)
	expected := `{"clob_pairs":[{"id":0,"perpetual_clob_metadata":{"perpetual_id":0},`
	expected += `"step_base_quantums":"5","subticks_per_tick":100,`
	expected += `"quantum_conversion_exponent":0,"status":"STATUS_ACTIVE"}],`
	expected += `"liquidations_config":{`
	expected += `"max_liquidation_fee_ppm":5000,"position_block_limits":{"min_position_notional_liquidated":"1000",`
	expected += `"max_position_portion_liquidated_ppm":1000000},"subaccount_block_limits":`
	expected += `{"max_notional_liquidated":"100000000000000","max_quantums_insurance_lost":"100000000000000"},`
	expected += `"fillable_price_config":{"bankruptcy_adjustment_ppm":1000000,`
	expected += `"spread_to_maintenance_margin_ratio_ppm":100000}},"block_rate_limit_config":`
	expected += `{"max_leverage_updates_per_n_blocks":[],`
	expected += `"max_short_term_orders_and_cancels_per_n_blocks":[{"limit": 400,"num_blocks":1}],`
	expected += `"max_stateful_orders_per_n_blocks":[{"limit": 2,"num_blocks":1},`
	expected += `{"limit": 20,"num_blocks":100}],"max_short_term_order_cancellations_per_n_blocks":[],`
	expected += `"max_short_term_orders_per_n_blocks":[]},`
	expected += `"equity_tier_limit_config":{"short_term_order_equity_tiers":[{"limit":0,"usd_tnc_required":"0"},`
	expected += `{"limit":1,"usd_tnc_required":"20"},{"limit":5,"usd_tnc_required":"100"},`
	expected += `{"limit":10,"usd_tnc_required":"1000"},{"limit":100,"usd_tnc_required":"10000"},`
	expected += `{"limit":1000,"usd_tnc_required":"100000"}],"stateful_order_equity_tiers":[`
	expected += `{"limit":0,"usd_tnc_required":"0"},{"limit":1,"usd_tnc_required":"20"},`
	expected += `{"limit":5,"usd_tnc_required":"100"},{"limit":10,"usd_tnc_required":"1000"},`
	expected += `{"limit":100,"usd_tnc_required":"10000"},{"limit":200,"usd_tnc_required":"100000"}]}}`
	require.JSONEq(t, expected, string(genesisJson))
}

func TestAppModule_InitGenesisPanic(t *testing.T) {
	am, _, _, _, _, ctx, _ := createAppModuleWithKeeper(t)
	cdc := codec.NewProtoCodec(module.InterfaceRegistry)
	gs := json.RawMessage(`invalid json`)

	require.Panics(t, func() { am.InitGenesis(ctx, cdc, gs) })
}

func TestAppModule_ConsensusVersion(t *testing.T) {
	am := createAppModule(t)
	require.Equal(t, uint64(1), am.ConsensusVersion())
}

func TestAppModule_BeginBlock(t *testing.T) {
	am, _, _, _, _, ctx, _ := createAppModuleWithKeeper(t)

	require.NoError(t, am.BeginBlock(ctx)) // should not panic
}

func TestAppModule_EndBlock(t *testing.T) {
	am, _, _, _, _, ctx, _ := createAppModuleWithKeeper(t)
	ctx = ctx.WithBlockTime(constants.TimeT)

	require.NoError(t, am.EndBlock(ctx))
}
