//go:build all || integration_test

package cli_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/nullify"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// func networkWithLiquidityTierAndPerpetualObjects(
// 	t *testing.T,
// 	m int,
// 	n int,
// ) (
// 	*network.Network,
// 	[]types.LiquidityTier,
// 	[]types.Perpetual,
// ) {
// 	t.Helper()
// 	cfg := network.DefaultConfig(nil)

// 	// Init Prices state.
// 	pricesState := constants.Prices_DefaultGenesisState
// 	pricesBuf, pricesErr := cfg.Codec.MarshalJSON(&pricesState)
// 	require.NoError(t, pricesErr)
// 	cfg.GenesisState[pricestypes.ModuleName] = pricesBuf

// 	// Init Perpetuals state.
// 	state := types.GenesisState{}
// 	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

// 	// Generate `m` Liquidity Tiers.
// 	for i := 0; i < m; i++ {
// 		liquidityTier := types.LiquidityTier{
// 			Id:                     uint32(i),
// 			Name:                   fmt.Sprintf("test_liquidity_tier_name_%d", i),
// 			InitialMarginPpm:       uint32(1_000_000 / (i + 1)),
// 			MaintenanceFractionPpm: uint32(1_000_000 / (i + 1)),
// 			ImpactNotional:         uint64(500_000_000 * (i + 1)),
// 		}
// 		nullify.Fill(&liquidityTier) //nolint:staticcheck
// 		state.LiquidityTiers = append(state.LiquidityTiers, liquidityTier)
// 	}

// 	// Generate `n` Perpetuals.
// 	for i := 0; i < n; i++ {
// 		perpetual := types.Perpetual{
// 			Params: types.PerpetualParams{
// 				Id:            uint32(i),
// 				Ticker:        fmt.Sprintf("test_query_ticker_%d", i),
// 				LiquidityTier: uint32(i % m),
// 			},
// 			FundingIndex: dtypes.ZeroInt(),
// 		}
// 		nullify.Fill(&perpetual) //nolint:staticcheck
// 		state.Perpetuals = append(state.Perpetuals, perpetual)
// 	}
// 	buf, err := cfg.Codec.MarshalJSON(&state)
// 	require.NoError(t, err)
// 	cfg.GenesisState[types.ModuleName] = buf

// 	return network.New(t, cfg), state.LiquidityTiers, state.Perpetuals
// }

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
		perpetual := types.Perpetual{
			Params: types.PerpetualParams{
				Id:            uint32(i),
				Ticker:        fmt.Sprintf("test_query_ticker_%d", i),
				LiquidityTier: uint32(i % m),
			},
			FundingIndex: dtypes.ZeroInt(),
		}
		nullify.Fill(&perpetual) //nolint:staticcheck
		state.Perpetuals = append(state.Perpetuals, perpetual)
	}

	return state.LiquidityTiers, state.Perpetuals
}

func getPerpetualGenesisShort() string {

	//"{\\\"asset_positions\\\": [{\\\"quantums\\\": \\\"1000\\\"}], \\\"id\\\": {\\\"number\\\": 2, \\\"owner\\\": \\\"0\\\"}, \\\"margin_enabled\\\": true}"
	// {
	// 	"asset_positions": [
	// 	  {
	// 		"asset_id": 0,
	// 		"index": 0,
	// 		"quantums": "100000000000000000"
	// 	  }
	// 	],
	// 	"id": {
	// 	  "number": 0,
	// 	  "owner": "dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4"
	// 	},
	// 	"margin_enabled": true
	//   },
	return "\".app_state.perpetuals.liquidity_tiers = [{\\\"name\\\": \\\"test_liquidity_tier_name_0\\\", \\\"initial_margin_ppm\\\": \\\"1000000\\\", \\\"maintenance_fraction_ppm\\\": \\\"1000000\\\", \\\"impact_notional\\\": \\\"500000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"name\\\": \\\"test_liquidity_tier_name_1\\\", \\\"initial_margin_ppm\\\": \\\"500000\\\", \\\"maintenance_fraction_ppm\\\": \\\"500000\\\", \\\"impact_notional\\\": \\\"1000000000\\\"}] | .app_state.perpetuals.perpetuals = [{\\\"params\\\": {\\\"ticker\\\": \\\"test_query_ticker_0\\\"}, \\\"funding_index\\\": \\\"0\\\"}, {\\\"params\\\": {\\\"id\\\": \\\"1\\\", \\\"ticker\\\": \\\"test_query_ticker_1\\\", \\\"liquidity_tier\\\": \\\"1\\\"}, \\\"funding_index\\\": \\\"0\\\"}] | .app_state.prices.market_prices = [{\\\"exponent\\\": \\\"-5\\\", \\\"price\\\": \\\"5000000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"exponent\\\": \\\"-6\\\", \\\"price\\\": \\\"3000000000\\\"}] | .app_state.prices.market_params = [{\\\"pair\\\": \\\"BTC-USD\\\", \\\"exponent\\\": \\\"-5\\\", \\\"min_exchanges\\\": \\\"2\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tBTCUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC/USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XXBTZUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}, {\\\"id\\\": \\\"1\\\", \\\"pair\\\": \\\"ETH-USD\\\", \\\"exponent\\\": \\\"-6\\\", \\\"min_exchanges\\\": \\\"1\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tETHUSD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH/USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XETHZUSD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}]\" \"\""
}

func TestShowPerpetual(t *testing.T) {
	liq, objs := networkWithLiquidityTierAndPerpetualObjects(t, 2, 2)
	genesisChanges := getPerpetualGenesisShort()

	setupCmd := exec.Command("bash", "-c", "cd ../../../../ethos/ethos-chain && ./e2e-setup -setup "+genesisChanges)

	fmt.Println("Running setup command", setupCmd.String())
	var out bytes.Buffer
	var stderr bytes.Buffer
	setupCmd.Stdout = &out
	setupCmd.Stderr = &stderr
	err := setupCmd.Run()
	if err != nil {
		t.Fatalf("Failed to set up environment: %v, stdout: %s, stderr: %s", err, out.String(), stderr.String())
	}
	fmt.Println("Setup output:", out.String())

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
			args += " --node tcp://7.7.8.4:26658 -o json"

			cmd := exec.Command("bash", "-c", "docker exec interchain-security-instance interchain-security-cd query perpetuals show-perpetual "+args)
			var queryOut bytes.Buffer
			cmd.Stdout = &queryOut
			err = cmd.Run()

			require.NoError(t, err)
			var resp types.QueryPerpetualResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(queryOut.Bytes(), &resp))
			require.NotNil(t, resp.Perpetual)
			checkExpectedPerp(t, tc.obj, resp.Perpetual)

		})
	}

	stopCmd := exec.Command("bash", "-c", "docker stop interchain-security-instance")
	if err := stopCmd.Run(); err != nil {
		t.Fatalf("Failed to stop Docker container: %v", err)
	}
	fmt.Println("Stopped Docker container")
	// Remove the Docker container
	removeCmd := exec.Command("bash", "-c", "docker rm interchain-security-instance")
	if err := removeCmd.Run(); err != nil {
		t.Fatalf("Failed to remove Docker container: %v", err)
	}
	fmt.Println("Removed Docker container")
}

// Check the received perpetual matches with expected.
// FundingIndex field is ignored since it can vary depending on funding-tick epoch.
// TODO(DEC-606): Improve end-to-end testing related to ticking epochs.
func checkExpectedPerp(t *testing.T, expected types.Perpetual, received types.Perpetual) {
	if diff := cmp.Diff(expected, received, cmpopts.IgnoreFields(types.Perpetual{}, "FundingIndex")); diff != "" {
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

func TestListPerpetual(t *testing.T) {
	liq, objs := networkWithLiquidityTierAndPerpetualObjects(t, 3, 5)
	jsonObj, _ := json.MarshalIndent(liq, "", "  ")
	fmt.Println("liq: ", string(jsonObj))
	jsonObj, _ = json.MarshalIndent(objs, "", "  ")
	fmt.Println("objs: ", string(jsonObj))

	// request := func(next []byte, offset, limit uint64, total bool) []string {
	// 	args := []string{
	// 		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	// 	}
	// 	if next == nil {
	// 		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
	// 	} else {
	// 		args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
	// 	}
	// 	args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
	// 	if total {
	// 		args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
	// 	}
	// 	return args
	// }
	t.Run("ByOffset", func(t *testing.T) {
		// step := 2
		// for i := 0; i < len(objs); i += step {
		// 	args := request(nil, uint64(i), uint64(step), false)
		// 	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPerpetual(), args)
		// 	require.NoError(t, err)
		// 	var resp types.QueryAllPerpetualsResponse
		// 	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		// 	require.LessOrEqual(t, len(resp.Perpetual), step)
		// 	for _, perp := range resp.Perpetual {
		// 		expectedContainsReceived(t, objs, perp)
		// 	}
		// }

		jsonObj, _ := json.MarshalIndent(liq, "", "  ")
		fmt.Println("liq: ", string(jsonObj))
		jsonObj, _ = json.MarshalIndent(objs, "", "  ")
		fmt.Println("objs: ", string(jsonObj))

		require.Equal(t, 1, 0)

	})
	// t.Run("ByKey", func(t *testing.T) {
	// 	step := 2
	// 	var next []byte
	// 	for i := 0; i < len(objs); i += step {
	// 		args := request(next, 0, uint64(step), false)
	// 		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPerpetual(), args)
	// 		require.NoError(t, err)
	// 		var resp types.QueryAllPerpetualsResponse
	// 		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	// 		require.LessOrEqual(t, len(resp.Perpetual), step)
	// 		for _, perp := range resp.Perpetual {
	// 			expectedContainsReceived(t, objs, perp)
	// 		}
	// 		next = resp.Pagination.NextKey
	// 	}
	// })
	// t.Run("Total", func(t *testing.T) {
	// 	args := request(nil, 0, uint64(len(objs)), true)
	// 	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPerpetual(), args)
	// 	require.NoError(t, err)
	// 	var resp types.QueryAllPerpetualsResponse
	// 	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	// 	require.NoError(t, err)
	// 	require.Equal(t, len(objs), int(resp.Pagination.Total))
	// 	cmpOptions := []cmp.Option{
	// 		cmpopts.IgnoreFields(types.Perpetual{}, "FundingIndex"),
	// 		cmpopts.SortSlices(func(x, y types.Perpetual) bool {
	// 			return x.Params.Id > y.Params.Id
	// 		}),
	// 	}
	// 	if diff := cmp.Diff(objs, resp.Perpetual, cmpOptions...); diff != "" {
	// 		t.Errorf("resp.Perpetual mismatch (-want +received):\n%s", diff)
	// 	}
	// })
}
