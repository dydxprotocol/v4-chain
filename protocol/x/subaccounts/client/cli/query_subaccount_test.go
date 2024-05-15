//go:build all || integration_test

package cli_test

import (
	"bytes"
	"fmt"
	"math/big"
	"os/exec"
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

// func networkWithSubaccountObjects(t *testing.T, n int) (*network.Network, []types.Subaccount) {
// 	t.Helper()
// 	cfg := network.DefaultConfig(nil)
// 	state := types.GenesisState{}
// 	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

// 	for i := 0; i < n; i++ {
// 		subaccount := types.Subaccount{
// 			Id: &types.SubaccountId{
// 				Owner:  strconv.Itoa(i),
// 				Number: uint32(n),
// 			},
// 			AssetPositions: keepertest.CreateUsdcAssetPosition(big.NewInt(1_000)),
// 		}
// 		nullify.Fill(&subaccount) //nolint:staticcheck
// 		state.Subaccounts = append(state.Subaccounts, subaccount)
// 	}
// 	buf, err := cfg.Codec.MarshalJSON(&state)
// 	require.NoError(t, err)
// 	cfg.GenesisState[types.ModuleName] = buf

// 	return network.New(t, cfg), state.Subaccounts
// }

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
	return "\".app_state.subaccounts.subaccounts = [{\\\"asset_positions\\\": [{\\\"quantums\\\": \\\"1000\\\"}], \\\"id\\\": {\\\"number\\\": \\\"2\\\", \\\"owner\\\": \\\"0\\\"}, \\\"margin_enabled\\\": false}, {\\\"asset_positions\\\": [{\\\"quantums\\\": \\\"1000\\\"}], \\\"id\\\": {\\\"number\\\": \\\"2\\\", \\\"owner\\\": \\\"1\\\"}, \\\"margin_enabled\\\": false}]\" \"\""

}

func TestShowSubaccount(t *testing.T) {
	objs := networkWithSubaccountObjects(t, 2)
	cfg := network.DefaultConfig(nil)

	genesisChanges := getSubaccountGenesisShort()

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

			cmd := exec.Command("bash", "-c", "docker exec interchain-security-instance interchain-security-cd query subaccounts show-subaccount "+args[0]+" "+args[1]+" --node tcp://7.7.8.4:26658 -o json")
			var queryOut bytes.Buffer
			var stdQueryErr bytes.Buffer
			cmd.Stdout = &queryOut
			cmd.Stderr = &stdQueryErr
			err = cmd.Run()

			fmt.Println("Query output:", queryOut.String())
			fmt.Println("Query error:", stdQueryErr.String())

			require.NoError(t, err)
			var resp types.QuerySubaccountResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(queryOut.Bytes(), &resp))
			require.NotNil(t, resp.Subaccount)
			require.Equal(t,
				nullify.Fill(&tc.obj),          //nolint:staticcheck
				nullify.Fill(&resp.Subaccount), //nolint:staticcheck
			)

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

func getSubaccountGenesisList() string {

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

			cmd := exec.Command("bash", "-c", "docker exec interchain-security-instance interchain-security-cd query subaccounts list-subaccount"+args)
			var out bytes.Buffer
			cmd.Stdout = &out
			err = cmd.Run()

			require.NoError(t, err)
			var resp types.QuerySubaccountAllResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(out.Bytes(), &resp))
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

			cmd := exec.Command("bash", "-c", "docker exec interchain-security-instance interchain-security-cd query subaccounts list-subaccount "+args)
			var out bytes.Buffer
			cmd.Stdout = &out

			fmt.Println("Running command", cmd.String())
			err = cmd.Run()

			require.NoError(t, err)
			var resp types.QuerySubaccountAllResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(out.Bytes(), &resp))
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

		cmd := exec.Command("bash", "-c", "docker exec interchain-security-instance interchain-security-cd query subaccounts list-subaccount "+args)
		var out bytes.Buffer
		cmd.Stdout = &out

		fmt.Println("Running command", cmd.String())

		err = cmd.Run()

		require.NoError(t, err)
		var resp types.QuerySubaccountAllResponse
		require.NoError(t, cfg.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),            //nolint:staticcheck
			nullify.Fill(resp.Subaccount), //nolint:staticcheck
		)
	})

	// stopCmd := exec.Command("bash", "-c", "docker stop interchain-security-instance")
	// if err := stopCmd.Run(); err != nil {
	// 	t.Fatalf("Failed to stop Docker container: %v", err)
	// }
	// fmt.Println("Stopped Docker container")
	// // Remove the Docker container
	// removeCmd := exec.Command("bash", "-c", "docker rm interchain-security-instance")
	// if err := removeCmd.Run(); err != nil {
	// 	t.Fatalf("Failed to remove Docker container: %v", err)
	// }
	// fmt.Println("Removed Docker container")
}
