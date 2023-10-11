package flags_test

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/flags"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddFlagsToCommand(t *testing.T) {
	cmd := cobra.Command{}

	flags.AddClobFlagsToCmd(&cmd)
	tests := map[string]struct {
		flagName string
	}{
		fmt.Sprintf("Has %s flag", flags.MaxLiquidationAttemptsPerBlock): {
			flagName: flags.MaxLiquidationAttemptsPerBlock,
		},
		fmt.Sprintf("Has %s flag", flags.MaxDeleveragingAttemptsPerBlock): {
			flagName: flags.MaxDeleveragingAttemptsPerBlock,
		},
		fmt.Sprintf("Has %s flag", flags.MevTelemetryHost): {
			flagName: flags.MevTelemetryHost,
		},
		fmt.Sprintf("Has %s flag", flags.MevTelemetryIdentifier): {
			flagName: flags.MevTelemetryIdentifier,
		}}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Contains(t, cmd.Flags().FlagUsages(), tc.flagName)
		})
	}
}

func TestGetFlagValuesFromOptions(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		optsMap map[string]any

		// Expectations.
		expectedMaxLiquidationAttemptsPerBlock      uint32
		expectedMaxDeleveragingAttemptsPerBlock     uint32
		expectedMaxDeleveragingSubaccountsToIterate uint32
		expectedMevTelemetryHost                    string
		expectedMevTelemetryIdentifier              string
	}{
		"Sets to default if unset": {
			expectedMaxLiquidationAttemptsPerBlock:      flags.DefaultMaxLiquidationAttemptsPerBlock,
			expectedMaxDeleveragingAttemptsPerBlock:     flags.DefaultMaxDeleveragingAttemptsPerBlock,
			expectedMaxDeleveragingSubaccountsToIterate: flags.DefaultMaxDeleveragingSubaccountsToIterate,
			expectedMevTelemetryHost:                    flags.DefaultMevTelemetryHost,
			expectedMevTelemetryIdentifier:              flags.DefaultMevTelemetryIdentifier,
		},
		"Sets values from options": {
			optsMap: map[string]any{
				flags.MaxLiquidationAttemptsPerBlock:      uint32(50),
				flags.MaxDeleveragingAttemptsPerBlock:     uint32(25),
				flags.MaxDeleveragingSubaccountsToIterate: uint32(100),
				flags.MevTelemetryHost:                    "https://localhost:13137",
				flags.MevTelemetryIdentifier:              "node-agent-01",
			},
			expectedMaxLiquidationAttemptsPerBlock:      uint32(50),
			expectedMaxDeleveragingAttemptsPerBlock:     uint32(25),
			expectedMaxDeleveragingSubaccountsToIterate: uint32(100),
			expectedMevTelemetryHost:                    "https://localhost:13137",
			expectedMevTelemetryIdentifier:              "node-agent-01",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockOpts := mocks.AppOptions{}
			mockOpts.On("Get", mock.AnythingOfType("string")).
				Return(func(key string) interface{} {
					return tc.optsMap[key]
				})

			flags := flags.GetClobFlagValuesFromOptions(&mockOpts)
			require.Equal(
				t,
				tc.expectedMevTelemetryHost,
				flags.MevTelemetryHost,
			)
			require.Equal(
				t,
				tc.expectedMevTelemetryIdentifier,
				flags.MevTelemetryIdentifier,
			)
			require.Equal(
				t,
				tc.expectedMaxLiquidationAttemptsPerBlock,
				flags.MaxLiquidationAttemptsPerBlock,
			)
			require.Equal(
				t,
				tc.expectedMaxDeleveragingAttemptsPerBlock,
				flags.MaxDeleveragingAttemptsPerBlock,
			)
			require.Equal(
				t,
				tc.expectedMaxDeleveragingSubaccountsToIterate,
				flags.MaxDeleveragingSubaccountsToIterate,
			)
		})
	}
}
