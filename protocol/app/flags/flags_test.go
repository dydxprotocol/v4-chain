package flags_test

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4/app/flags"
	"github.com/dydxprotocol/v4/mocks"
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
		expectedNonValidatingFullNodeFlag bool
		expectedDdAgentHost               string
		expectedDdTraceAgentPort          uint16
	}{
		"Sets to default if unset": {
			expectedNonValidatingFullNodeFlag: false,
			expectedDdAgentHost:               "",
			expectedDdTraceAgentPort:          8126,
		},
		"Sets values from options": {
			optsMap: map[string]any{
				flags.NonValidatingFullNodeFlag: true,
				flags.DdAgentHost:               "agentHostTest",
				flags.DdTraceAgentPort:          uint16(777),
			},
			expectedNonValidatingFullNodeFlag: true,
			expectedDdAgentHost:               "agentHostTest",
			expectedDdTraceAgentPort:          777,
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
		})
	}
}
