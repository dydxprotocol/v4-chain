//go:build all || integration_test

package cli_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/stretchr/testify/require"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/nullify"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	serializableIntCmpOpts = cmp.Options{
		cmpopts.IgnoreFields(types.Perpetual{}, "FundingIndex"), // existing ignore option
		cmp.FilterValues(
			func(x, y dtypes.SerializableInt) bool {
				// This will apply the custom comparer only to SerializableInt fields
				return true // Apply this filter to all SerializableInt values
			},
			cmp.Comparer(func(x, y dtypes.SerializableInt) bool { return x.Cmp(y) == 0 }),
		),
	}
)

func networkWithLiquidityTierAndPerpetualObjects(
	t *testing.T,
	m int,
	n int,
) (
	[]types.LiquidityTier,
	[]types.Perpetual,
) {
	t.Helper()
	cfg := network.DefaultConfig(nil)

	// Init Perpetuals state.
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	// Generate `m` Liquidity Tiers.
	for i := 0; i < m; i++ {
		liquidityTier := types.LiquidityTier{
			Id:                     uint32(i),
			Name:                   fmt.Sprintf("test_liquidity_tier_name_%d", i),
			InitialMarginPpm:       uint32(1_000_000 / (i + 1)),
			MaintenanceFractionPpm: uint32(1_000_000 / (i + 1)),
			ImpactNotional:         uint64(500_000_000 * (i + 1)),
		}
		nullify.Fill(&liquidityTier) //nolint:staticcheck
		state.LiquidityTiers = append(state.LiquidityTiers, liquidityTier)
	}

	// Generate `n` Perpetuals.
	for i := 0; i < n; i++ {
		marketType := types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS
		IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock := uint64(0)
		if i%2 == 1 {
			marketType = types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED
			IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock = uint64(1_000_000)
		}

		perpetual := types.Perpetual{
			Params: types.PerpetualParams{
				Id:             uint32(i),
				Ticker:         fmt.Sprintf("test_query_ticker_%d", i),
				LiquidityTier:  uint32(i % m),
				MarketType:     marketType,
				DangerIndexPpm: uint32(0),
				IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock: IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
			},
			FundingIndex:    dtypes.ZeroInt(),
			OpenInterest:    dtypes.ZeroInt(),
			YieldIndex:      "0/1",
			LastFundingRate: dtypes.ZeroInt(),
		}
		nullify.Fill(&perpetual) //nolint:staticcheck
		state.Perpetuals = append(state.Perpetuals, perpetual)
	}

	return state.LiquidityTiers, state.Perpetuals
}

func GetPerpetualGenesisShort() string {
	return "\".app_state.perpetuals.liquidity_tiers = [{\\\"name\\\": \\\"test_liquidity_tier_name_0\\\", \\\"initial_margin_ppm\\\": \\\"1000000\\\", \\\"maintenance_fraction_ppm\\\": \\\"1000000\\\", \\\"impact_notional\\\": \\\"500000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"name\\\": \\\"test_liquidity_tier_name_1\\\", \\\"initial_margin_ppm\\\": \\\"500000\\\", \\\"maintenance_fraction_ppm\\\": \\\"500000\\\", \\\"impact_notional\\\": \\\"1000000000\\\"}] | .app_state.perpetuals.perpetuals = [{\\\"params\\\": {\\\"ticker\\\": \\\"test_query_ticker_0\\\", \\\"market_type\\\": \\\"PERPETUAL_MARKET_TYPE_CROSS\\\"}, \\\"funding_index\\\": \\\"0\\\"}, {\\\"params\\\": {\\\"id\\\": \\\"1\\\", \\\"ticker\\\": \\\"test_query_ticker_1\\\", \\\"liquidity_tier\\\": \\\"1\\\", \\\"market_type\\\": \\\"PERPETUAL_MARKET_TYPE_CROSS\\\"}, \\\"funding_index\\\": \\\"0\\\"}] | .app_state.prices.market_prices = [{\\\"exponent\\\": \\\"-5\\\", \\\"spot_price\\\": \\\"5000000000\\\", \\\"pnl_price\\\": \\\"5000000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"exponent\\\": \\\"-6\\\", \\\"spot_price\\\": \\\"3000000000\\\", \\\"pnl_price\\\": \\\"3000000000\\\"}] | .app_state.prices.market_params = [{\\\"pair\\\": \\\"BTC-USD\\\", \\\"exponent\\\": \\\"-5\\\", \\\"min_exchanges\\\": \\\"2\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tBTCUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC/USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XXBTZUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}, {\\\"id\\\": \\\"1\\\", \\\"pair\\\": \\\"ETH-USD\\\", \\\"exponent\\\": \\\"-6\\\", \\\"min_exchanges\\\": \\\"1\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tETHUSD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH/USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XETHZUSD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}]\" \"\""
}

func TestShowPerpetual(t *testing.T) {
	_, objs := networkWithLiquidityTierAndPerpetualObjects(t, 2, 2)
	genesisChanges := GetPerpetualGenesisShort()
	network.DeployCustomNetwork(genesisChanges)

	cfg := network.DefaultConfig(nil)

	for _, tc := range []struct {
		desc string
		id   uint32

		args []string
		err  error
		obj  types.Perpetual
	}{
		{
			desc: "found",
			id:   objs[0].Params.Id,

			obj: objs[0],
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := fmt.Sprintf("%v", tc.id)
			perpQuery := "docker exec interchain-security-instance interchain-security-cd" +
				" query perpetuals show-perpetual " + args
			data, _, err := network.QueryCustomNetwork(perpQuery)
			require.NoError(t, err)
			var resp types.QueryPerpetualResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
			require.NotNil(t, resp.Perpetual)
			checkExpectedPerp(t, tc.obj, resp.Perpetual)
		})
	}
	network.CleanupCustomNetwork()
}

// Check the received perpetual matches with expected.
// FundingIndex field is ignored since it can vary depending on funding-tick epoch.
// TODO(DEC-606): Improve end-to-end testing related to ticking epochs.
func checkExpectedPerp(t *testing.T, expected types.Perpetual, received types.Perpetual) {
	if diff := cmp.Diff(expected, received, serializableIntCmpOpts...); diff != "" {
		t.Errorf("resp.Perpetual mismatch (-want +received):\n%s", diff)
	}
}

// Check the received perpetual object matches one of the expected perpetuals.
func expectedContainsReceived(t *testing.T, expectedPerps []types.Perpetual, received types.Perpetual) {
	for _, expected := range expectedPerps {
		if received.Params.Id == expected.Params.Id {
			checkExpectedPerp(t, expected, received)
			return
		}
	}
	t.Errorf("Received perp (%v) not found in expected perps (%v)", received, expectedPerps)
}

func getPerpetualGenesisList() string {
	return "\".app_state.perpetuals.liquidity_tiers = [{\\\"name\\\": \\\"test_liquidity_tier_name_0\\\", \\\"initial_margin_ppm\\\": \\\"1000000\\\", \\\"maintenance_fraction_ppm\\\": \\\"1000000\\\", \\\"impact_notional\\\": \\\"500000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"name\\\": \\\"test_liquidity_tier_name_1\\\", \\\"initial_margin_ppm\\\": \\\"500000\\\", \\\"maintenance_fraction_ppm\\\": \\\"500000\\\", \\\"impact_notional\\\": \\\"1000000000\\\"}, {\\\"id\\\": \\\"2\\\", \\\"name\\\": \\\"test_liquidity_tier_name_2\\\", \\\"initial_margin_ppm\\\": \\\"333333\\\", \\\"maintenance_fraction_ppm\\\": \\\"333333\\\", \\\"impact_notional\\\": \\\"1500000000\\\"}] | .app_state.perpetuals.perpetuals = [{\\\"params\\\": {\\\"ticker\\\": \\\"test_query_ticker_0\\\", \\\"market_type\\\": \\\"PERPETUAL_MARKET_TYPE_CROSS\\\"}, \\\"funding_index\\\": \\\"0\\\", \\\"yield_index\\\": \\\"0/1\\\"}, {\\\"params\\\": {\\\"id\\\": \\\"1\\\", \\\"ticker\\\": \\\"test_query_ticker_1\\\", \\\"liquidity_tier\\\": \\\"1\\\", \\\"market_type\\\": \\\"PERPETUAL_MARKET_TYPE_ISOLATED\\\", \\\"isolated_market_max_cumulative_insurance_fund_delta_per_block\\\": \\\"1000000\\\"}, \\\"funding_index\\\": \\\"0\\\", \\\"yield_index\\\": \\\"0/1\\\"}, {\\\"params\\\": {\\\"id\\\": \\\"2\\\", \\\"ticker\\\": \\\"test_query_ticker_2\\\", \\\"liquidity_tier\\\": \\\"2\\\", \\\"market_type\\\": \\\"PERPETUAL_MARKET_TYPE_CROSS\\\"}, \\\"funding_index\\\": \\\"0\\\", \\\"yield_index\\\": \\\"0/1\\\"}, {\\\"params\\\": {\\\"id\\\": \\\"3\\\", \\\"ticker\\\": \\\"test_query_ticker_3\\\", \\\"market_type\\\": \\\"PERPETUAL_MARKET_TYPE_ISOLATED\\\", \\\"isolated_market_max_cumulative_insurance_fund_delta_per_block\\\": \\\"1000000\\\"}, \\\"funding_index\\\": \\\"0\\\", \\\"yield_index\\\": \\\"0/1\\\"}, {\\\"params\\\": {\\\"id\\\": \\\"4\\\", \\\"ticker\\\": \\\"test_query_ticker_4\\\", \\\"liquidity_tier\\\": \\\"1\\\", \\\"market_type\\\": \\\"PERPETUAL_MARKET_TYPE_CROSS\\\"}, \\\"funding_index\\\": \\\"0\\\", \\\"yield_index\\\": \\\"0/1\\\"}] | .app_state.prices.market_prices = [{\\\"exponent\\\": \\\"-5\\\", \\\"spot_price\\\": \\\"5000000000\\\", \\\"pnl_price\\\": \\\"5000000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"exponent\\\": \\\"-6\\\", \\\"spot_price\\\": \\\"3000000000\\\", \\\"pnl_price\\\": \\\"3000000000\\\"}] | .app_state.prices.market_params = [{\\\"pair\\\": \\\"BTC-USD\\\", \\\"exponent\\\": \\\"-5\\\", \\\"min_exchanges\\\": \\\"2\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tBTCUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC/USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XXBTZUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}, {\\\"id\\\": \\\"1\\\", \\\"pair\\\": \\\"ETH-USD\\\", \\\"exponent\\\": \\\"-6\\\", \\\"min_exchanges\\\": \\\"1\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tETHUSD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH/USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XETHZUSD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}]\" \"\""
}

func TestListPerpetual(t *testing.T) {
	_, objs := networkWithLiquidityTierAndPerpetualObjects(t, 3, 5)

	genesisChanges := getPerpetualGenesisList()

	network.DeployCustomNetwork(genesisChanges)

	cfg := network.DefaultConfig(nil)

	request := func(next []byte, offset, limit uint64, total bool) string {
		args := ""
		if next == nil {
			args += fmt.Sprintf(" --%s=%d", flags.FlagOffset, offset)
		} else {
			base64Next := base64.StdEncoding.EncodeToString(next)
			args += fmt.Sprintf(" --%s=%s", flags.FlagPageKey, base64Next)
		}
		args += fmt.Sprintf(" --%s=%d", flags.FlagLimit, limit)
		if total {
			args += fmt.Sprintf(" --%s", flags.FlagCountTotal)
		}
		return args
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			perpQuery := "docker exec interchain-security-instance interchain-security-cd" +
				" query perpetuals list-perpetual " + args
			data, _, err := network.QueryCustomNetwork(perpQuery)

			require.NoError(t, err)
			var resp types.QueryAllPerpetualsResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
			require.LessOrEqual(t, len(resp.Perpetual), step)
			for _, perp := range resp.Perpetual {
				expectedContainsReceived(t, objs, perp)
			}
		}
	})

	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			perpQuery := "docker exec interchain-security-instance interchain-security-cd" +
				" query perpetuals list-perpetual " + args
			data, _, err := network.QueryCustomNetwork(perpQuery)

			require.NoError(t, err)
			var resp types.QueryAllPerpetualsResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
			require.LessOrEqual(t, len(resp.Perpetual), step)
			for _, perp := range resp.Perpetual {
				expectedContainsReceived(t, objs, perp)
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		perpQuery := "docker exec interchain-security-instance interchain-security-cd" +
			" query perpetuals list-perpetual " + args
		data, _, err := network.QueryCustomNetwork(perpQuery)
		require.NoError(t, err)
		var resp types.QueryAllPerpetualsResponse
		require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		cmpOptions := append(
			serializableIntCmpOpts,
			cmpopts.SortSlices(func(x, y types.Perpetual) bool {
				return x.Params.Id > y.Params.Id
			}),
		)
		if diff := cmp.Diff(objs, resp.Perpetual, cmpOptions...); diff != "" {
			t.Errorf("resp.Perpetual mismatch (-want +received):\n%s", diff)
		}
	})

	network.CleanupCustomNetwork()
}
