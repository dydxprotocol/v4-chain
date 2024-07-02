package testutil

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/network"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

// CreateBankGenesisState returns the marshaled genesis state for the bank module.
// It will set the balance of the subaccount module in the genesis.
// If the provided subaccount module balance is negative, this function will panic.
func CreateBankGenesisState(
	t testing.TB,
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
