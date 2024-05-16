//go:build all || integration_test

package cli_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/epochs/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func checkExpectedEpoch(
	t *testing.T,
	networkStartTime time.Time,
	genesis types.EpochInfo,
	actual types.EpochInfo) {
	require.Equal(t,
		genesis.Duration,
		actual.Duration,
	)
	// Check that EpochInfo is initialized.
	require.True(t,
		actual.IsInitialized,
	)
	// Check that NextTick is fast forwarded.
	require.Less(t,
		uint32(networkStartTime.Unix()),
		actual.NextTick,
	)
	// Check that NextTick is fast forwarded by exact multiples of duration.
	require.Zero(t,
		(actual.NextTick-genesis.NextTick)%genesis.Duration,
	)
}

func networkWithEpochInfoObjects(t *testing.T) network.Config {
	t.Helper()
	cfg := network.DefaultConfig(nil)
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	return cfg
}

func getDefaultGenesisEpochById(t *testing.T, id string) types.EpochInfo {
	for _, epochInfo := range types.DefaultGenesis().GetEpochInfoList() {
		if epochInfo.Name == id {
			return epochInfo
		}
	}
	panic(fmt.Errorf("DefaultGenesisEpoch not found (%s)", id))
}

func TestShowEpochInfo(t *testing.T) {
	//cfg := networkWithEpochInfoObjects(t)

	cfg := network.DefaultConfig(nil)

	networkStartTime := time.Now()

	// net := network.New(t, cfg)
	// _, err := net.WaitForHeight(3)
	// require.NoError(t, err)

	// we use the default docker already running

	// ctx := net.Validators[0].ClientCtx
	// common := []string{
	// 	fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	// }
	for _, tc := range []struct {
		desc string
		id   string
		args []string
		err  error
	}{
		{
			desc: "found default funding tick epoch",
			id:   string(types.FundingTickEpochInfoName),
			//args: common,
		},
		{
			desc: "found default funding sample epoch",
			id:   string(types.FundingSampleEpochInfoName),
			//args: common,
		},
		{
			desc: "not found",
			id:   strconv.Itoa(100000),
			//args: common,
			err: status.Error(codes.NotFound, "not found"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			// args := []string{
			// 	tc.id,
			// }
			// args = append(args, tc.args...)
			// out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowEpochInfo(), args)

			cmd := exec.Command("bash", "-c", "docker exec interchain-security-instance-setup interchain-security-cd query epochs show-epoch-info "+tc.id+" --node tcp://7.7.8.4:26658 -o json")
			var queryOut bytes.Buffer
			var stdQueryErr bytes.Buffer
			cmd.Stdout = &queryOut
			cmd.Stderr = &stdQueryErr
			err := cmd.Run()

			if tc.err != nil {
				// stat, ok := status.FromError(tc.err)
				// require.True(t, ok)
				// require.ErrorIs(t, stat.Err(), tc.err)

				// require that stdQueryErr contains the string not found
				require.Contains(t, stdQueryErr.String(), "not found")

			} else {
				require.NoError(t, err)

				genesisEpoch := getDefaultGenesisEpochById(t, tc.id)

				var resp types.QueryEpochInfoResponse
				require.NoError(t, cfg.Codec.UnmarshalJSON(queryOut.Bytes(), &resp))
				require.NotNil(t, resp.EpochInfo)
				checkExpectedEpoch(t, networkStartTime, genesisEpoch, resp.EpochInfo)
			}
		})
	}
}

func TestListEpochInfo(t *testing.T) {
	//cfg := networkWithEpochInfoObjects(t)

	cfg := network.DefaultConfig(nil)
	networkStartTime := time.Now()

	// net := network.New(t, cfg)
	// _, err := net.WaitForHeight(3)
	// require.NoError(t, err)

	objs := types.DefaultGenesis().GetEpochInfoList()

	//ctx := net.Validators[0].ClientCtx
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
			base64Next := base64.StdEncoding.EncodeToString(next)
			args += fmt.Sprintf(" --%s=%s", flags.FlagPageKey, base64Next)
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

			cmd := exec.Command("bash", "-c", "docker exec interchain-security-instance-setup interchain-security-cd query epochs list-epoch-info"+args)
			var queryOut bytes.Buffer
			cmd.Stdout = &queryOut
			err := cmd.Run()

			//out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListEpochInfo(), args)
			require.NoError(t, err)
			var resp types.QueryEpochInfoAllResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(queryOut.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.EpochInfo), step)
			for _, epoch := range resp.EpochInfo {
				genesisEpoch := getDefaultGenesisEpochById(t, epoch.Name)
				checkExpectedEpoch(t, networkStartTime, genesisEpoch, epoch)
			}
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)

			//out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListEpochInfo(), args)

			cmd := exec.Command("bash", "-c", "docker exec interchain-security-instance-setup interchain-security-cd query epochs list-epoch-info"+args)
			var queryOut bytes.Buffer
			cmd.Stdout = &queryOut
			err := cmd.Run()

			require.NoError(t, err)
			var resp types.QueryEpochInfoAllResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(queryOut.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.EpochInfo), step)
			for _, epoch := range resp.EpochInfo {
				genesisEpoch := getDefaultGenesisEpochById(t, epoch.Name)
				checkExpectedEpoch(t, networkStartTime, genesisEpoch, epoch)
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)

		//out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListEpochInfo(), args)

		cmd := exec.Command("bash", "-c", "docker exec interchain-security-instance-setup interchain-security-cd query epochs list-epoch-info"+args)
		var queryOut bytes.Buffer
		cmd.Stdout = &queryOut
		err := cmd.Run()
		require.NoError(t, err)

		var resp types.QueryEpochInfoAllResponse
		require.NoError(t, cfg.Codec.UnmarshalJSON(queryOut.Bytes(), &resp))
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		for _, epoch := range resp.EpochInfo {
			genesisEpoch := getDefaultGenesisEpochById(t, epoch.Name)
			checkExpectedEpoch(t, networkStartTime, genesisEpoch, epoch)
		}
	})
}
