package flags_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/dydxprotocol/v4-chain/protocol/app/flags"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddFlagsToCommand(t *testing.T) {
	cmd := cobra.Command{}

	flags.AddFlagsToCmd(&cmd)
	tests := map[string]struct {
		flagName string
	}{
		fmt.Sprintf("Has %s flag", flags.NonValidatingFullNodeFlag): {
			flagName: flags.NonValidatingFullNodeFlag,
		},
		fmt.Sprintf("Has %s flag", flags.DdAgentHost): {
			flagName: flags.DdAgentHost,
		},
		fmt.Sprintf("Has %s flag", flags.DdTraceAgentPort): {
			flagName: flags.DdTraceAgentPort,
		},
		fmt.Sprintf("Has %s flag", flags.GrpcStreamingEnabled): {
			flagName: flags.GrpcStreamingEnabled,
		},
		fmt.Sprintf("Has %s flag", flags.GrpcStreamingFlushIntervalMs): {
			flagName: flags.GrpcStreamingFlushIntervalMs,
		},
		fmt.Sprintf("Has %s flag", flags.GrpcStreamingMaxBatchSize): {
			flagName: flags.GrpcStreamingMaxBatchSize,
		},
		fmt.Sprintf("Has %s flag", flags.FullNodeStreamingSnapshotInterval): {
			flagName: flags.FullNodeStreamingSnapshotInterval,
		},
		fmt.Sprintf("Has %s flag", flags.GrpcStreamingMaxChannelBufferSize): {
			flagName: flags.GrpcStreamingMaxChannelBufferSize,
		},
		fmt.Sprintf("Has %s flag", flags.WebsocketStreamingEnabled): {
			flagName: flags.WebsocketStreamingEnabled,
		},
		fmt.Sprintf("Has %s flag", flags.WebsocketStreamingPort): {
			flagName: flags.WebsocketStreamingPort,
		},
		fmt.Sprintf("Has %s flag", flags.OptimisticExecutionEnabled): {
			flagName: flags.OptimisticExecutionEnabled,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Contains(t, cmd.Flags().FlagUsages(), tc.flagName)
		})
	}
}

func TestValidate(t *testing.T) {
	tests := map[string]struct {
		flags       flags.Flags
		expectedErr error
	}{
		"success (default values)": {
			flags: flags.Flags{
				NonValidatingFullNode:             flags.DefaultNonValidatingFullNode,
				DdAgentHost:                       flags.DefaultDdAgentHost,
				DdTraceAgentPort:                  flags.DefaultDdTraceAgentPort,
				GrpcAddress:                       config.DefaultGRPCAddress,
				GrpcEnable:                        true,
				FullNodeStreamingSnapshotInterval: flags.DefaultFullNodeStreamingSnapshotInterval,
				OptimisticExecutionEnabled:        false,
			},
		},
		"success - full node & gRPC disabled": {
			flags: flags.Flags{
				GrpcEnable:            false,
				NonValidatingFullNode: true,
			},
		},
		"success - gRPC streaming enabled for validating nodes": {
			flags: flags.Flags{
				NonValidatingFullNode:             false,
				GrpcEnable:                        true,
				GrpcStreamingEnabled:              true,
				GrpcStreamingFlushIntervalMs:      100,
				GrpcStreamingMaxBatchSize:         10000,
				GrpcStreamingMaxChannelBufferSize: 10000,
				WebsocketStreamingEnabled:         false,
			},
		},
		"success - both grpc and websocket streaming enabled for validating nodes": {
			flags: flags.Flags{
				NonValidatingFullNode:             false,
				GrpcEnable:                        true,
				GrpcStreamingEnabled:              true,
				GrpcStreamingFlushIntervalMs:      100,
				GrpcStreamingMaxBatchSize:         10000,
				GrpcStreamingMaxChannelBufferSize: 10000,
				WebsocketStreamingEnabled:         true,
				WebsocketStreamingPort:            8989,
			},
		},
		"success - optimistic execution": {
			flags: flags.Flags{
				NonValidatingFullNode:      false,
				GrpcEnable:                 true,
				OptimisticExecutionEnabled: true,
			},
		},
		"success - optimistic execution canbe  enabled with gRPC streaming": {
			flags: flags.Flags{
				NonValidatingFullNode:             false,
				GrpcEnable:                        true,
				GrpcStreamingEnabled:              true,
				OptimisticExecutionEnabled:        true,
				GrpcStreamingMaxBatchSize:         10000,
				GrpcStreamingFlushIntervalMs:      100,
				GrpcStreamingMaxChannelBufferSize: 10000,
				WebsocketStreamingPort:            8989,
			},
		},
		"failure - gRPC disabled": {
			flags: flags.Flags{
				GrpcEnable: false,
			},
			expectedErr: fmt.Errorf("grpc.enable must be set to true - validating requires gRPC server"),
		},
		"failure - gRPC streaming enabled with gRPC disabled": {
			flags: flags.Flags{
				NonValidatingFullNode: true,
				GrpcEnable:            false,
				GrpcStreamingEnabled:  true,
			},
			expectedErr: fmt.Errorf("grpc.enable must be set to true - grpc streaming requires gRPC server"),
		},
		"failure - websocket streaming enabled with gRPC streaming disabled": {
			flags: flags.Flags{
				NonValidatingFullNode:             true,
				GrpcEnable:                        true,
				GrpcStreamingEnabled:              false,
				WebsocketStreamingEnabled:         true,
				GrpcStreamingFlushIntervalMs:      100,
				GrpcStreamingMaxBatchSize:         10000,
				GrpcStreamingMaxChannelBufferSize: 10000,
			},
			expectedErr: fmt.Errorf("websocket full node streaming requires grpc streaming to be enabled"),
		},
		"success - websocket streaming enabled with gRPC enabled for validating node": {
			flags: flags.Flags{
				NonValidatingFullNode:             true,
				GrpcEnable:                        true,
				WebsocketStreamingEnabled:         true,
				GrpcStreamingEnabled:              true,
				GrpcStreamingFlushIntervalMs:      100,
				GrpcStreamingMaxBatchSize:         10000,
				GrpcStreamingMaxChannelBufferSize: 10000,
				WebsocketStreamingPort:            8989,
			},
		},
		"failure - gRPC streaming enabled with zero batch size": {
			flags: flags.Flags{
				NonValidatingFullNode:        true,
				GrpcEnable:                   true,
				GrpcStreamingEnabled:         true,
				GrpcStreamingFlushIntervalMs: 100,
				GrpcStreamingMaxBatchSize:    0,
			},
			expectedErr: fmt.Errorf("full node streaming batch size must be positive number"),
		},
		"failure - gRPC streaming enabled with zero flush interval ms": {
			flags: flags.Flags{
				NonValidatingFullNode:        true,
				GrpcEnable:                   true,
				GrpcStreamingEnabled:         true,
				GrpcStreamingFlushIntervalMs: 0,
				GrpcStreamingMaxBatchSize:    2000,
			},
			expectedErr: fmt.Errorf("full node streaming flush interval must be positive number"),
		},
		"failure - gRPC streaming enabled with zero channel size ms": {
			flags: flags.Flags{
				NonValidatingFullNode:             true,
				GrpcEnable:                        true,
				GrpcStreamingEnabled:              true,
				GrpcStreamingFlushIntervalMs:      100,
				GrpcStreamingMaxBatchSize:         2000,
				GrpcStreamingMaxChannelBufferSize: 0,
			},
			expectedErr: fmt.Errorf("full node streaming channel size must be positive number"),
		},
		"failure - websocket streaming enabled with zero batch size": {
			flags: flags.Flags{
				NonValidatingFullNode:        true,
				GrpcEnable:                   true,
				GrpcStreamingEnabled:         true,
				GrpcStreamingFlushIntervalMs: 100,
				GrpcStreamingMaxBatchSize:    0,
				WebsocketStreamingEnabled:    true,
			},
			expectedErr: fmt.Errorf("full node streaming batch size must be positive number"),
		},
		"success - full node streaming enabled with 20 snapshot interval": {
			flags: flags.Flags{
				NonValidatingFullNode:             true,
				GrpcEnable:                        true,
				GrpcStreamingEnabled:              true,
				GrpcStreamingFlushIntervalMs:      100,
				GrpcStreamingMaxBatchSize:         2000,
				GrpcStreamingMaxChannelBufferSize: 2000,
				FullNodeStreamingSnapshotInterval: 20,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.flags.Validate()
			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestGetFlagValuesFromOptions(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		optsMap map[string]any

		// Expectations.
		expectedNonValidatingFullNodeFlag         bool
		expectedDdAgentHost                       string
		expectedDdTraceAgentPort                  uint16
		expectedGrpcAddress                       string
		expectedGrpcEnable                        bool
		expectedGrpcStreamingEnable               bool
		expectedGrpcStreamingFlushMs              uint32
		expectedGrpcStreamingBatchSize            uint32
		expectedGrpcStreamingMaxChannelBufferSize uint32
		expectedWebsocketEnabled                  bool
		expectedWebsocketPort                     uint16
		expectedFullNodeStreamingSnapshotInterval uint32
		expectedOptimisticExecutionEnabled        bool
	}{
		"Sets to default if unset": {
			expectedNonValidatingFullNodeFlag:         false,
			expectedDdAgentHost:                       "",
			expectedDdTraceAgentPort:                  8126,
			expectedGrpcAddress:                       "localhost:9090",
			expectedGrpcEnable:                        true,
			expectedGrpcStreamingEnable:               false,
			expectedGrpcStreamingFlushMs:              50,
			expectedGrpcStreamingBatchSize:            100_000,
			expectedGrpcStreamingMaxChannelBufferSize: 100_000,
			expectedWebsocketEnabled:                  false,
			expectedWebsocketPort:                     9092,
			expectedFullNodeStreamingSnapshotInterval: 0,
			expectedOptimisticExecutionEnabled:        false,
		},
		"Sets values from options": {
			optsMap: map[string]any{
				flags.NonValidatingFullNodeFlag:         true,
				flags.DdAgentHost:                       "agentHostTest",
				flags.DdTraceAgentPort:                  uint16(777),
				flags.GrpcEnable:                        false,
				flags.GrpcAddress:                       "localhost:1234",
				flags.GrpcStreamingEnabled:              "true",
				flags.GrpcStreamingFlushIntervalMs:      uint32(408),
				flags.GrpcStreamingMaxBatchSize:         uint32(650),
				flags.GrpcStreamingMaxChannelBufferSize: uint32(972),
				flags.WebsocketStreamingEnabled:         "true",
				flags.WebsocketStreamingPort:            8989,
				flags.FullNodeStreamingSnapshotInterval: uint32(123),
				flags.OptimisticExecutionEnabled:        "true",
			},
			expectedNonValidatingFullNodeFlag:         true,
			expectedDdAgentHost:                       "agentHostTest",
			expectedDdTraceAgentPort:                  777,
			expectedGrpcEnable:                        false,
			expectedGrpcAddress:                       "localhost:1234",
			expectedGrpcStreamingEnable:               true,
			expectedGrpcStreamingFlushMs:              408,
			expectedGrpcStreamingBatchSize:            650,
			expectedGrpcStreamingMaxChannelBufferSize: 972,
			expectedWebsocketEnabled:                  true,
			expectedWebsocketPort:                     8989,
			expectedFullNodeStreamingSnapshotInterval: 123,
			expectedOptimisticExecutionEnabled:        true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockOpts := mocks.AppOptions{}
			mockOpts.On("Get", mock.AnythingOfType("string")).
				Return(func(key string) interface{} {
					return tc.optsMap[key]
				})

			flags := flags.GetFlagValuesFromOptions(&mockOpts)
			require.Equal(
				t,
				tc.expectedNonValidatingFullNodeFlag,
				flags.NonValidatingFullNode,
			)
			require.Equal(
				t,
				tc.expectedDdAgentHost,
				flags.DdAgentHost,
			)
			require.Equal(
				t,
				tc.expectedDdTraceAgentPort,
				flags.DdTraceAgentPort,
			)
			require.Equal(
				t,
				tc.expectedGrpcEnable,
				flags.GrpcEnable,
			)
			require.Equal(
				t,
				tc.expectedGrpcAddress,
				flags.GrpcAddress,
			)
			require.Equal(
				t,
				tc.expectedGrpcStreamingEnable,
				flags.GrpcStreamingEnabled,
			)
			require.Equal(
				t,
				tc.expectedGrpcStreamingFlushMs,
				flags.GrpcStreamingFlushIntervalMs,
			)
			require.Equal(
				t,
				tc.expectedGrpcStreamingBatchSize,
				flags.GrpcStreamingMaxBatchSize,
			)
			require.Equal(
				t,
				tc.expectedFullNodeStreamingSnapshotInterval,
				flags.FullNodeStreamingSnapshotInterval,
			)
			require.Equal(
				t,
				tc.expectedGrpcStreamingMaxChannelBufferSize,
				flags.GrpcStreamingMaxChannelBufferSize,
			)
			require.Equal(
				t,
				tc.expectedWebsocketEnabled,
				flags.WebsocketStreamingEnabled,
			)
			require.Equal(
				t,
				tc.expectedWebsocketPort,
				flags.WebsocketStreamingPort,
			)
		})
	}
}
