package flags

import (
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
)

// A struct containing the values of all test flags.
type TestFlags struct {
	ExampleFlag int64
}

// List of test CLI flags.
const (
	TestFlagExampleFlag = "testflags-example-flag"
)

// Default values.
const (
	DefaultExampleFlag = int64(0)
)

// AddTestFlagsToCmd adds the flags to app initialization that are used for testing.
// These flags should be applied to the `start` command of the V4 Cosmos application.
// E.g. `dydxprotocold start --testflags-example-flag <value>`.
func AddTestFlagsToCmd(cmd *cobra.Command) {
	cmd.
		Flags().
		Int64(
			TestFlagExampleFlag,
			DefaultExampleFlag,
			"Test flag used as an example of how to add additional test-flags in the future.",
		)
}

// GetTestFlagValuesFromOptions gets values used for testing from the `AppOptions` struct which contains values
// from the command-line flags.
func GetTestFlagValuesFromOptions(
	appOpts servertypes.AppOptions,
) TestFlags {
	// Create default result.
	result := TestFlags{
		ExampleFlag: DefaultExampleFlag,
	}

	// Populate the testflags if they exist.

	// Note the cast to `int` instead of `int64` (`int64` conversion fails for some reason).
	exampleFlagValue, ok := appOpts.Get(TestFlagExampleFlag).(int)
	if ok {
		result.ExampleFlag = int64(exampleFlagValue)
	}

	return result
}
