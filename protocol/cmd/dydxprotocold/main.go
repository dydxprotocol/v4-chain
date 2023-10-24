package main

import (
	"github.com/cosmos/cosmos-sdk/server"
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/app/config"
	"github.com/dydxprotocol/v4-chain/protocol/cmd/dydxprotocold/cmd"
)

func main() {
	config.SetupConfig()

	option := cmd.GetOptionWithCustomStartCmd()
	rootCmd := cmd.NewRootCmd(
		option,
		func(serverCtxPtr *server.Context) {},
		func(s string, appConfig *cmd.DydxAppConfig) (string, *cmd.DydxAppConfig) {
			return s, appConfig
		},
		func(app *app.App) {},
	)

	cmd.AddTendermintSubcommands(rootCmd)

	if err := svrcmd.Execute(rootCmd, app.AppDaemonName, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
