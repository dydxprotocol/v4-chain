package flags

import (
	"github.com/spf13/cobra"
)

const (
	FlagAuthenticators = "authenticators"
)

// AddTxPermissionedKeyFlagsToCmd adds common flags to a module tx command.
func AddTxPermissionedKeyFlagsToCmd(cmd *cobra.Command) {
	f := cmd.Flags()
	f.UintSlice(FlagAuthenticators, nil, "Authenticators to use for authenticating this transaction.")
}

// GetPermisionedKeyAuthenticatorsForExtOptions returns the authenticators from the provided command flags.
func GetPermisionedKeyAuthenticatorsForExtOptions(cmd *cobra.Command) ([]uint64, error) {
	flags := cmd.Flags()
	values, err := flags.GetUintSlice(FlagAuthenticators)
	if err == nil {
		authenticators := make([]uint64, len(values))
		for i, v := range values {
			authenticators[i] = uint64(v)
		}
		return authenticators, nil
	}
	return nil, err
}
