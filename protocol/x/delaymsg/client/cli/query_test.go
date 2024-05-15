//go:build all || integration_test

package cli_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/encoding"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	GrpcNotFoundError = "NotFound"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func setupNetwork(
	t *testing.T,
	state *types.GenesisState,
) (
	*network.Network,
	client.Context,
) {
	t.Helper()
	cfg := network.DefaultConfig(nil)

	// Init state.
	// Validate global genesis state contains a delaymsg genesis state.
	configDefaultGenesisState := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &configDefaultGenesisState))

	// Update global genesis state with specified delaymsg genesis state.
	buf, err := cfg.Codec.MarshalJSON(state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	// Create network.
	net := network.New(t, cfg)
	ctx := net.Validators[0].ClientCtx

	return net, ctx
}

func TestQueryNextDelayedMessageId(t *testing.T) {
	tests := map[string]struct {
		state *types.GenesisState
	}{
		"Default: 0": {
			state: types.DefaultGenesis(),
		},
		"Non-zero": {
			state: &types.GenesisState{
				DelayedMessages:      []*types.DelayedMessage{},
				NextDelayedMessageId: 20,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// _, ctx := setupNetwork(t, tc.state)
			// out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryNextDelayedMessageId(), []string{})
			// require.NoError(t, err)
			// var resp types.QueryNextDelayedMessageIdResponse
			// require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
			// require.Equal(t, tc.state.NextDelayedMessageId, resp.NextDelayedMessageId)

			cfg := network.DefaultConfig(nil)

			cmd := exec.Command("docker", "exec", "interchain-security-instance-setup", "interchain-security-cd", "query", "delaymsg", "get-next-delayed-message-id", "--node", "tcp://7.7.8.4:26658", "-o json")
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()

			require.NoError(t, err)
			var resp types.QueryNextDelayedMessageIdResponse
			data := out.Bytes()
			require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
			require.Equal(t, tc.state.NextDelayedMessageId, resp.NextDelayedMessageId)

		})
	}
}

func TestQueryMessage(t *testing.T) {
	tests := map[string]struct {
		state       *types.GenesisState
		expectedMsg sdk.Msg
	}{
		"Default: 0": {
			state: types.DefaultGenesis(),
		},
		"Non-zero": {
			state: &types.GenesisState{
				DelayedMessages: []*types.DelayedMessage{
					{
						Id:  0,
						Msg: encoding.EncodeMessageToAny(t, constants.TestMsg1),
					},
				},
				NextDelayedMessageId: 20,
			},
			expectedMsg: constants.TestMsg1,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// _, ctx := setupNetwork(t, tc.state)
			// out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryMessage(), []string{"0"})
			// if tc.expectedMsg == nil {
			// 	require.ErrorContains(t, err, GrpcNotFoundError)
			// } else {
			// 	require.NoError(t, err)
			// 	var resp types.QueryMessageResponse
			// 	require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))

			// 	err := resp.Message.UnpackInterfaces(ctx.Codec)
			// 	require.NoError(t, err)
			// 	msg, err := resp.Message.GetMessage()
			// 	require.NoError(t, err)

			// 	require.Equal(t, tc.expectedMsg, msg)
			// }

			fmt.Println("PRINTING TEST CASE", name)

			genesisChanges := getGenesisChanges(name)

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

			cmd := exec.Command("bash", "-c", "docker", "exec", "interchain-security-instance-setup", "interchain-security-cd", "query", "delaymsg", "get-message", "0", "--node", "tcp://7.7.8.4:26658", "-o json")
			var queryOut bytes.Buffer
			var stdQueryErr bytes.Buffer
			cmd.Stdout = &queryOut
			cmd.Stderr = &stdQueryErr
			err = cmd.Run()

			if tc.expectedMsg == nil {

				require.True(t, strings.Contains(stdQueryErr.String(), GrpcNotFoundError))
			} else {

				require.NoError(t, err)
				var resp types.QueryMessageResponse
				data := queryOut.Bytes()
				require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))

				err := resp.Message.UnpackInterfaces(cfg.Codec)
				require.NoError(t, err)
				msg, err := resp.Message.GetMessage()
				require.NoError(t, err)

				require.Equal(t, tc.expectedMsg, msg)

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

		})
	}
}

func TestQueryBlockMessageIds(t *testing.T) {
	tests := map[string]struct {
		state                   *types.GenesisState
		expectedBlockMessageIds []uint32
	}{
		"Default: 0": {
			state: types.DefaultGenesis(),
		},
		"Non-zero": {
			state: &types.GenesisState{
				DelayedMessages: []*types.DelayedMessage{
					{
						Id:          0,
						Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg1),
						BlockHeight: 10,
					},
				},
				NextDelayedMessageId: 20,
			},
			expectedBlockMessageIds: []uint32{0},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// _, ctx := setupNetwork(t, tc.state)
			// out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryBlockMessageIds(), []string{"10"})
			// if tc.expectedBlockMessageIds == nil {
			// 	require.ErrorContains(t, err, GrpcNotFoundError)
			// } else {
			// 	require.NoError(t, err)
			// 	var resp types.QueryBlockMessageIdsResponse
			// 	require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
			// 	require.Equal(t, tc.expectedBlockMessageIds, resp.MessageIds)
			// }

			// fmt.Println("Starting test", types.ModuleAddress.String())

			// cmd := exec.Command("bash", "-c", "cd ../../../../ && ls") // You can add arguments inside the quotes, e.g., "ls", "-la", "/some/path"
			// var out bytes.Buffer
			// var stderr bytes.Buffer
			// cmd.Stdout = &out
			// cmd.Stderr = &stderr
			// err := cmd.Run()
			// if err != nil {
			// 	fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			// 	return
			// }
			// fmt.Println("Command Successfully Executed")
			// output := out.String()
			// fmt.Println(output)

			fmt.Println(encoding.EncodeMessageToAny(t, constants.TestMsg1))

			genesisChanges := getGenesisChanges(name)
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

			testCmd := exec.Command("bash", "-c", "docker ps")
			var testOut bytes.Buffer
			var testErr bytes.Buffer
			testCmd.Stdout = &testOut
			testCmd.Stderr = &testErr
			err = testCmd.Run()
			if err != nil {
				t.Fatalf("Failed to list Docker containers: %v, stdout: %s, stderr: %s", err, testOut.String(), testErr.String())
			}
			fmt.Println("Docker containers:", testOut.String())
			cfg := network.DefaultConfig(nil)

			cmd := exec.Command("bash", "-c", "docker exec interchain-security-instance interchain-security-cd query delaymsg get-block-message-ids 1000 --node tcp://7.7.8.4:26658 -o json")

			var queryOut bytes.Buffer
			var stdQueryErr bytes.Buffer
			cmd.Stdout = &queryOut
			cmd.Stderr = &stdQueryErr
			err = cmd.Run()

			// if err != nil {
			// 	fmt.Printf("Error running query command, stdout: %s, stderr: %s", queryOut.String(), stdQueryErr.String())
			// }

			if name == "Default: 0" {
				require.True(t, strings.Contains(stdQueryErr.String(), GrpcNotFoundError))

			} else {

				require.NoError(t, err)
				var resp types.QueryBlockMessageIdsResponse
				data := queryOut.Bytes()
				require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
				require.Equal(t, tc.expectedBlockMessageIds, resp.MessageIds)
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

		})
	}
}

func getGenesisChanges(testCase string) string {
	switch testCase {
	case "Default: 0":
		return "\".app_state.delaymsg.delayed_messages = [] | .app_state.delaymsg.next_delayed_message_id = \"0\"\" \"\""
	case "Non-zero":
		// setup(".app_state.delaymsg.delayed_messages[0] = {\"id\": \"0\", \"msg\": {\"@type\": \"/dydxprotocol.perpetuals.MsgUpdateParams\", \"authority\": \"dydx1mkkvp26dngu6n8rmalaxyp3gwkjuzztq5zx6tr\", \"params\": {\"funding_rate_clamp_factor_ppm\": \"6000000\", \"premium_vote_clamp_factor_ppm\": \"60000000\", \"min_num_votes_per_sample\": \"15\"}}, \"block_height\": \"10\"} | .app_state.delaymsg.next_delayed_message_id = \"20\"", "")
		return "\".app_state.delaymsg.delayed_messages[0] = {\\\"id\\\": \\\"0\\\", \\\"msg\\\": {\\\"@type\\\": \\\"/dydxprotocol.perpetuals.MsgUpdateParams\\\", \\\"authority\\\": \\\"dydx1mkkvp26dngu6n8rmalaxyp3gwkjuzztq5zx6tr\\\", \\\"params\\\": {\\\"funding_rate_clamp_factor_ppm\\\": \\\"6000000\\\", \\\"premium_vote_clamp_factor_ppm\\\": \\\"60000000\\\", \\\"min_num_votes_per_sample\\\": \\\"15\\\"}}, \\\"block_height\\\": \\\"1000\\\"} | .app_state.delaymsg.next_delayed_message_id = \\\"20\\\"\" \"\""

	default:
		panic("unknown case")
	}
}
