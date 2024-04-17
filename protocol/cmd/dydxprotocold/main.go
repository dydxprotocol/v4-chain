package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/config"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/cmd/dydxprotocold/cmd"
)

func main() {
	config.SetupConfig()

	option := cmd.GetOptionWithCustomStartCmd()
	rootCmd := cmd.NewRootCmd(option, app.DefaultNodeHome)

	cmd.AddTendermintSubcommands(rootCmd)
	cmd.AddInitCmdPostRunE(rootCmd)

	if err := svrcmd.Execute(rootCmd, constants.AppDaemonName, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
