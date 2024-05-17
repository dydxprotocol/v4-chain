//go:build all || integration_test

package cli_test

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/stretchr/testify/require"

	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/nullify"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithSubaccountObjects(t *testing.T, n int) []types.Subaccount {
	t.Helper()
	cfg := network.DefaultConfig(nil)
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		subaccount := types.Subaccount{
			AssetPositions: keepertest.CreateUsdcAssetPosition(big.NewInt(1_000)),
			Id: &types.SubaccountId{
				Owner:  strconv.Itoa(i),
				Number: uint32(n),
			},
		}
		nullify.Fill(&subaccount) //nolint:staticcheck
		state.Subaccounts = append(state.Subaccounts, subaccount)
	}

	fmt.Println("state.Subaccounts", state.Subaccounts)

	return state.Subaccounts
}

func getSubaccountGenesisShort() string {

	return "\".app_state.subaccounts.subaccounts = [{\\\"asset_positions\\\": [{\\\"quantums\\\": \\\"1000\\\"}], \\\"id\\\": {\\\"number\\\": \\\"2\\\", \\\"owner\\\": \\\"0\\\"}, \\\"margin_enabled\\\": false}, {\\\"asset_positions\\\": [{\\\"quantums\\\": \\\"1000\\\"}], \\\"id\\\": {\\\"number\\\": \\\"2\\\", \\\"owner\\\": \\\"1\\\"}, \\\"margin_enabled\\\": false}]\" \"\""

}

func TestShowSubaccount(t *testing.T) {
	objs := networkWithSubaccountObjects(t, 2)
	cfg := network.DefaultConfig(nil)

	genesisChanges := getSubaccountGenesisShort()
	network.DeployCustomNetwork(genesisChanges)

	// ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc   string
		owner  string
		number uint32

		args []string
		err  error
		obj  types.Subaccount
	}{
		{
			desc:   "found",
			owner:  objs[0].Id.Owner,
			number: objs[0].Id.Number,
			args:   common,
			obj:    objs[0],
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.owner,
				strconv.Itoa(int(tc.number)),
			}

			fmt.Println("Args:", args)
			subQuery := "docker exec interchain-security-instance interchain-security-cd query subaccounts show-subaccount " + args[0] + " " + args[1] + " --node tcp://7.7.8.4:26658 -o json"
			data, _, err := newtork.QueryCustomNetwork(subQuery)
			require.NoError(t, err)
			var resp types.QuerySubaccountResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
			require.NotNil(t, resp.Subaccount)
			require.Equal(t,
				nullify.Fill(&tc.obj),          //nolint:staticcheck
				nullify.Fill(&resp.Subaccount), //nolint:staticcheck
			)

		})
	}

	network.CleanupCustomNetwork()
}

func getSubaccountGenesisList() string {

	return "\".app_state.subaccounts.subaccounts = [{\\\"asset_positions\\\": [{\\\"quantums\\\": \\\"1000\\\"}], \\\"id\\\": {\\\"owner\\\": \\\"0\\\", \\\"number\\\": \\\"5\\\"}, \\\"margin_enabled\\\": false}, {\\\"asset_positions\\\": [{\\\"quantums\\\": \\\"1000\\\"}], \\\"id\\\": {\\\"owner\\\": \\\"1\\\", \\\"number\\\": \\\"5\\\"}, \\\"margin_enabled\\\": false}, {\\\"asset_positions\\\": [{\\\"quantums\\\": \\\"1000\\\"}], \\\"id\\\": {\\\"owner\\\": \\\"2\\\", \\\"number\\\": \\\"5\\\"}, \\\"margin_enabled\\\": false}, {\\\"asset_positions\\\": [{\\\"quantums\\\": \\\"1000\\\"}], \\\"id\\\": {\\\"owner\\\": \\\"3\\\", \\\"number\\\": \\\"5\\\"}, \\\"margin_enabled\\\": false}, {\\\"asset_positions\\\": [{\\\"quantums\\\": \\\"1000\\\"}], \\\"id\\\": {\\\"owner\\\": \\\"4\\\", \\\"number\\\": \\\"5\\\"}, \\\"margin_enabled\\\": false}]\" \"\""

}

func removeNewlines(data []byte) []byte {
	var result []byte
	for _, b := range data {
		if b != '\n' && b != '\r' { // Handle both Unix and Windows line endings
			result = append(result, b)
		}
	}
	return result
}

func TestListSubaccount(t *testing.T) {
	objs := networkWithSubaccountObjects(t, 5)
	cfg := network.DefaultConfig(nil)

	genesisChanges := getSubaccountGenesisList()
	network.DeployCustomNetwork(genesisChanges)

	request := func(next []byte, offset, limit uint64, total bool) string {
		args := ""
		if next == nil {
			args += fmt.Sprintf(" --%s=%d", flags.FlagOffset, offset)
		} else {
			args += fmt.Sprintf(" --%s=%s", flags.FlagPageKey, removeNewlines(next))
		}
		args += fmt.Sprintf(" --%s=%d", flags.FlagLimit, limit)
		if total {
			args += fmt.Sprintf(" --%s", flags.FlagCountTotal)
		}

		args += " --node tcp://7.7.8.4:26658 -o json"
		return args
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			subQuery := "docker exec interchain-security-instance interchain-security-cd query subaccounts list-subaccount" + args
			data, _, err := network.QueryCustomNetwork(subQuery)

			require.NoError(t, err)
			var resp types.QuerySubaccountAllResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
			require.LessOrEqual(t, len(resp.Subaccount), step)
			require.Subset(t,
				nullify.Fill(objs),            //nolint:staticcheck
				nullify.Fill(resp.Subaccount), //nolint:staticcheck
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			subQuery := "docker exec interchain-security-instance interchain-security-cd query subaccounts list-subaccount " + args
			data, _, err := network.QueryCustomNetwork(subQuery)
			require.NoError(t, err)
			var resp types.QuerySubaccountAllResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
			require.LessOrEqual(t, len(resp.Subaccount), step)
			require.Subset(t,
				nullify.Fill(objs),            //nolint:staticcheck
				nullify.Fill(resp.Subaccount), //nolint:staticcheck
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		subQuery := "docker exec interchain-security-instance interchain-security-cd query subaccounts list-subaccount " + args
		data, _, err := network.QueryCustomNetwork(subQuery)
		require.NoError(t, err)
		var resp types.QuerySubaccountAllResponse
		require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),            //nolint:staticcheck
			nullify.Fill(resp.Subaccount), //nolint:staticcheck
		)
	})

	network.ClearupCustomNetwork()
}
