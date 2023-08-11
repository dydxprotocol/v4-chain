package genprivkey

import (
	"errors"
	"fmt"
	"path/filepath"

	tmed25519 "github.com/cometbft/cometbft/crypto/ed25519"
	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/cometbft/cometbft/privval"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"
)

var (
	FlagHome     = "home"
	FlagMnemonic = "mnemonic"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen-priv-key",
		Short: "Generate a Tendermint private key file from a given mnemonic.",
		RunE: func(cmd *cobra.Command, args []string) error {
			homeDir, _ := cmd.Flags().GetString(FlagHome)
			if homeDir == "" {
				return errors.New("--home flag is required")
			}

			mnemonic, _ := cmd.Flags().GetString(FlagMnemonic)
			if mnemonic == "" {
				return errors.New("--mnemonic flag is required")
			}

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(homeDir)

			if len(mnemonic) > 0 && !bip39.IsMnemonicValid(mnemonic) {
				return fmt.Errorf("invalid mnemonic")
			}

			pvKeyFile := config.PrivValidatorKeyFile()
			if err := tmos.EnsureDir(filepath.Dir(pvKeyFile), 0777); err != nil {
				return err
			}

			pvStateFile := config.PrivValidatorStateFile()
			if err := tmos.EnsureDir(filepath.Dir(pvStateFile), 0777); err != nil {
				return err
			}

			var filePV *privval.FilePV
			privKey := tmed25519.GenPrivKeyFromSecret([]byte(mnemonic))
			filePV = privval.NewFilePV(privKey, pvKeyFile, pvStateFile)

			filePV.Save()

			return nil
		},
	}

	cmd.Flags().String(FlagHome, "", "The application home directory")
	cmd.Flags().String(FlagMnemonic, "", "A bip39 mnemonic from which to generate the private key file")

	return cmd
}
