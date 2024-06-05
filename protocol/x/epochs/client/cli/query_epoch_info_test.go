//go:build all || integration_test

package cli_test

import (
	"encoding/base64"
	"fmt"
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
	genesis := "\".app_state.subaccounts.subaccounts = [{\\\"asset_positions\\\": [{\\\"quantums\\\": \\\"1000\\\"}], \\\"id\\\": {\\\"number\\\": \\\"2\\\", \\\"owner\\\": \\\"0\\\"}, \\\"margin_enabled\\\": false}, {\\\"asset_positions\\\": [{\\\"quantums\\\": \\\"1000\\\"}], \\\"id\\\": {\\\"number\\\": \\\"2\\\", \\\"owner\\\": \\\"1\\\"}, \\\"margin_enabled\\\": false}]\" \"\""
	network.DeployCustomNetwork(genesis)
	cfg := network.DefaultConfig(nil)

	networkStartTime := time.Now()

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
			epochQuery := "docker exec interchain-security-instance interchain-security-cd query epochs show-epoch-info " + tc.id
			data, stdQueryErr, err := network.QueryCustomNetwork(epochQuery)

			if tc.err != nil {
				require.Contains(t, stdQueryErr, "not found")

			} else {
				require.NoError(t, err)

				genesisEpoch := getDefaultGenesisEpochById(t, tc.id)

				var resp types.QueryEpochInfoResponse
				require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
				require.NotNil(t, resp.EpochInfo)
				checkExpectedEpoch(t, networkStartTime, genesisEpoch, resp.EpochInfo)
			}
		})
	}
	network.CleanupCustomNetwork()
}

func TestListEpochInfo(t *testing.T) {
	//cfg := networkWithEpochInfoObjects(t)
	genesis := "\".app_state.subaccounts.subaccounts = [{\\\"asset_positions\\\": [{\\\"quantums\\\": \\\"1000\\\"}], \\\"id\\\": {\\\"number\\\": \\\"2\\\", \\\"owner\\\": \\\"0\\\"}, \\\"margin_enabled\\\": false}, {\\\"asset_positions\\\": [{\\\"quantums\\\": \\\"1000\\\"}], \\\"id\\\": {\\\"number\\\": \\\"2\\\", \\\"owner\\\": \\\"1\\\"}, \\\"margin_enabled\\\": false}]\" \"\""
	network.DeployCustomNetwork(genesis)
	cfg := network.DefaultConfig(nil)
	networkStartTime := time.Now()

	objs := types.DefaultGenesis().GetEpochInfoList()

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
			epochQuery := "docker exec interchain-security-instance interchain-security-cd query epochs list-epoch-info " + args
			data, _, err := network.QueryCustomNetwork(epochQuery)

			require.NoError(t, err)
			var resp types.QueryEpochInfoAllResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
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

			epochQuery := "docker exec interchain-security-instance interchain-security-cd query epochs list-epoch-info " + args
			data, _, err := network.QueryCustomNetwork(epochQuery)
			require.NoError(t, err)
			var resp types.QueryEpochInfoAllResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
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

		epochQuery := "docker exec interchain-security-instance interchain-security-cd query epochs list-epoch-info " + args
		data, _, err := network.QueryCustomNetwork(epochQuery)
		require.NoError(t, err)

		var resp types.QueryEpochInfoAllResponse
		require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		for _, epoch := range resp.EpochInfo {
			genesisEpoch := getDefaultGenesisEpochById(t, epoch.Name)
			checkExpectedEpoch(t, networkStartTime, genesisEpoch, epoch)
		}
	})
	network.CleanupCustomNetwork()
}
