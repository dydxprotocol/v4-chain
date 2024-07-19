//go:build all || integration_test

package cli_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/encoding"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	GrpcNotFoundError = "NotFound"
)

// Prevent strconv unused error
var _ = strconv.IntSize

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
			cfg := network.DefaultConfig(nil)
			genesisChanges := getDelayedGenesisChanges(name)

			network.DeployCustomNetwork(genesisChanges)
			delaymsgQuery := "docker exec interchain-security-instance interchain-security-cd" +
				" query delaymsg get-next-delayed-message-id"
			data, _, _ := network.QueryCustomNetwork(delaymsgQuery)
			var resp types.QueryNextDelayedMessageIdResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
			require.Equal(t, tc.state.NextDelayedMessageId, resp.NextDelayedMessageId)

			network.CleanupCustomNetwork()
		})
	}
}

func getDelayedGenesisChanges(testCase string) string {
	switch testCase {
	case "Default: 0":
		return "\".app_state.delaymsg.delayed_messages = [] |" +
			" .app_state.delaymsg.next_delayed_message_id = \"0\"\" \"\""
	case "Non-zero":
		return "\".app_state.delaymsg.delayed_messages = [] |" +
			" .app_state.delaymsg.next_delayed_message_id = \"20\"\" \"\""

	default:
		panic("unknown case")
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

			genesisChanges := getGenesisChanges(name)

			network.DeployCustomNetwork(genesisChanges)

			cfg := network.DefaultConfig(nil)
			delaymsgQuery := "docker exec interchain-security-instance interchain-security-cd" +
				" query delaymsg get-message 0"
			data, stdQueryErr, err := network.QueryCustomNetwork(delaymsgQuery)

			if name == "Default: 0" {
				require.True(t, strings.Contains(stdQueryErr, GrpcNotFoundError))
			} else {
				require.NoError(t, err)
				var resp types.QueryMessageResponse

				require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))

				err := resp.Message.UnpackInterfaces(cfg.Codec)
				require.NoError(t, err)
				msg, err := resp.Message.GetMessage()
				require.NoError(t, err)

				require.Equal(t, tc.expectedMsg, msg)
			}

			network.CleanupCustomNetwork()
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
			genesisChanges := getGenesisChanges(name)
			network.DeployCustomNetwork(genesisChanges)

			cfg := network.DefaultConfig(nil)
			delaymsgQuery := "docker exec interchain-security-instance interchain-security-cd" +
				" query delaymsg get-block-message-ids 1000"
			data, stdQueryErr, err := network.QueryCustomNetwork(delaymsgQuery)

			if name == "Default: 0" {
				require.True(t, strings.Contains(stdQueryErr, GrpcNotFoundError))
			} else {
				require.NoError(t, err)
				var resp types.QueryBlockMessageIdsResponse
				require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
				require.Equal(t, tc.expectedBlockMessageIds, resp.MessageIds)
			}
			network.CleanupCustomNetwork()
		})
	}
}

func getGenesisChanges(testCase string) string {
	switch testCase {
	case "Default: 0":
		return "\".app_state.delaymsg.delayed_messages = [] |" +
			" .app_state.delaymsg.next_delayed_message_id = \"0\"\" \"\""
	case "Non-zero":
		return "\".app_state.delaymsg.delayed_messages[0] =" +
			" {\\\"id\\\": \\\"0\\\", \\\"msg\\\": {\\\"@type\\\": \\\"/dydxprotocol.perpetuals.MsgUpdateParams\\\"," +
			" \\\"authority\\\": \\\"dydx1mkkvp26dngu6n8rmalaxyp3gwkjuzztq5zx6tr\\\", \\\"params\\\":" +
			" {\\\"funding_rate_clamp_factor_ppm\\\": \\\"6000000\\\", \\\"premium_vote_clamp_factor_ppm\\\":" +
			" \\\"60000000\\\", \\\"min_num_votes_per_sample\\\": \\\"15\\\"}}, \\\"block_height\\\": \\\"1000\\\"} |" +
			" .app_state.delaymsg.next_delayed_message_id = \\\"20\\\"\" \"\""

	default:
		panic("unknown case")
	}
}
