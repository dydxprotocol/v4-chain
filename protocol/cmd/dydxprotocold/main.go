package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/dydxprotocol/v4/app"
	"github.com/dydxprotocol/v4/app/config"
	"github.com/dydxprotocol/v4/cmd/dydxprotocold/cmd"
)

func main() {
	config.SetupConfig()

	option := cmd.GetOptionWithCustomStartCmd()
	rootCmd := cmd.NewRootCmd(option)

	cmd.AddTendermintSubcommands(rootCmd)
	cmd.AddInitCmdPostRunE(rootCmd)

	if err := svrcmd.Execute(rootCmd, app.AppDaemonName, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
