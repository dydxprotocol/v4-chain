package testutil

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

// CreateBankGenesisState returns the marshaled genesis state for the bank module.
// It will set the balance of the subaccount module in the genesis.
// If the provided subaccount module balance is negative, this function will panic.
func CreateBankGenesisState(
	t *testing.T,
	cfg network.Config,
	initialSubaccountModuleBalance int64,
) []byte {
	bankGenState := banktypes.GenesisState{
		Balances: []banktypes.Balance{
			{
				Address: satypes.ModuleAddress.String(),
				Coins: []sdk.Coin{
					sdk.NewInt64Coin(
						constants.Usdc.Denom,
						initialSubaccountModuleBalance,
					),
				},
			},
		},
	}
	bankbuf, err := cfg.Codec.MarshalJSON(&bankGenState)
	require.NoError(t, err)

	return bankbuf
}
