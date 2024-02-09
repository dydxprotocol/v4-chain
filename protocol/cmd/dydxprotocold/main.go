package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/app/config"
	"github.com/dydxprotocol/v4-chain/protocol/app/constants"
	"github.com/dydxprotocol/v4-chain/protocol/cmd/dydxprotocold/cmd"
)

func main() {
	config.SetupConfig()

	option := cmd.GetOptionWithCustomStartCmd()
	// [ Home Dir Temp Fix ] (also see protocol/cmd/dydxprotocold/cmd/root.go)
	// We pass in a tempdir as a temporary hack until fixed in cosmos.
	// If not, and we pass in a custom home dir, it will try read or create the default home dir
	// before using the custom home dir.
	// See https://github.com/cosmos/cosmos-sdk/issues/18868.
	rootCmd := cmd.NewRootCmd(option, tempDir())

	cmd.AddTendermintSubcommands(rootCmd)
	cmd.AddInitCmdPostRunE(rootCmd)

	if err := svrcmd.Execute(rootCmd, constants.AppDaemonName, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}

var tempDir = func() string {
	dir, err := os.MkdirTemp("", "dydxprotocol")
	if err != nil {
		dir = app.DefaultNodeHome
	}
	defer os.RemoveAll(dir)

	return dir
}
