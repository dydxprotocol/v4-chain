package cmd

import "github.com/spf13/cobra"

// RootCmdOption configures root command option.
type RootCmdOption struct {
	startCmdCustomizer func(*cobra.Command)
}

// newRootCmdOption returns an empty RootCmdOption.
func newRootCmdOption() *RootCmdOption {
	return &RootCmdOption{}
}

// setCustomizeStartCmd accepts a handler to customize the start command and set it in the option.
func (o *RootCmdOption) setCustomizeStartCmd(f func(startCmd *cobra.Command)) {
	o.startCmdCustomizer = f
}
