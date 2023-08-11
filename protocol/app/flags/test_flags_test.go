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

func TestAddTestFlagsToCommand(t *testing.T) {
	cmd := cobra.Command{}

	flags.AddTestFlagsToCmd(&cmd)
	tests := map[string]struct {
		flagName string
	}{
		fmt.Sprintf("Has %s flag", flags.TestFlagExampleFlag): {
			flagName: flags.TestFlagExampleFlag,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Contains(t, cmd.Flags().FlagUsages(), tc.flagName)
		})
	}
}

func TestGetTestFlagValuesFromOptions(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		isSet       bool
		exampleFlag int

		// Expectations.
		expectedExampleFlagValue int64
	}{
		"Sets defaultQuoteBalance to default if unset": {
			expectedExampleFlagValue: flags.DefaultExampleFlag,
		},
		"Sets values from options": {
			isSet:                    true,
			exampleFlag:              555,
			expectedExampleFlagValue: 555,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			optsMap := make(map[string]interface{})
			if tc.isSet {
				optsMap[flags.TestFlagExampleFlag] = tc.exampleFlag
			}
			mockOpts := mocks.AppOptions{}
			mockOpts.On("Get", mock.AnythingOfType("string")).
				Return(func(key string) interface{} {
					return optsMap[key]
				})

			flags := flags.GetTestFlagValuesFromOptions(&mockOpts)
			require.Equal(
				t,
				tc.expectedExampleFlagValue,
				flags.ExampleFlag,
			)
		})
	}
}
