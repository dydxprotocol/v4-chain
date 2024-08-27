package simulation_test

import (
	"encoding/json"
	"testing"

	v4module "github.com/StreamFinance-Protocol/stream-chain/protocol/app/module"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
	testutil_rand "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/rand"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/sim_helpers"
	asstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/simulation"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	banksim "github.com/cosmos/cosmos-sdk/x/bank/simulation"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

const (
	numAccounts = 10
)

func TestRandomizedGenState(t *testing.T) {
	cdc := codec.NewProtoCodec(v4module.InterfaceRegistry)

	r := testutil_rand.NewRand()

	simState := module.SimulationState{
		AppParams:    make(simtypes.AppParams),
		Cdc:          cdc,
		Rand:         r,
		NumBonded:    3,
		Accounts:     simtypes.RandomAccounts(r, numAccounts),
		InitialStake: sdkmath.NewInt(1000),
		GenState:     make(map[string]json.RawMessage),
		BondDenom:    sdk.DefaultBondDenom,
	}
	for i := 0; i < 100; i++ {
		banksim.RandomizedGenState(&simState)
		simulation.RandomizedGenState(&simState)

		totalTDaiSupply := sdkmath.NewInt(0)

		var saGenesis types.GenesisState
		simState.Cdc.MustUnmarshalJSON(simState.GenState[types.ModuleName], &saGenesis)

		// at least 1 subaccount per account.
		require.True(t, len(saGenesis.Subaccounts) >= numAccounts)
		// at most 128 subaccounts per account.
		require.True(t, len(saGenesis.Subaccounts) <= numAccounts*sim_helpers.MaxNumSubaccount)

		for _, sa := range saGenesis.Subaccounts {
			require.True(t, sa.GetId().Number < 128)

			// AssetPositions.
			if len(sa.GetAssetPositions()) > 0 {
				require.Len(t, sa.GetAssetPositions(), 1)

				onlyAssetPosition := sa.GetAssetPositions()[0]
				require.True(t, onlyAssetPosition.AssetId == asstypes.AssetTDai.Id)

				bigQuantums := sdkmath.NewIntFromBigInt(onlyAssetPosition.GetBigQuantums())
				totalTDaiSupply = totalTDaiSupply.Add(bigQuantums)
			}

			require.False(t, sa.MarginEnabled)
		}

		var bankGenesis banktypes.GenesisState
		bankGenStateJson := simState.GenState[banktypes.ModuleName]
		simState.Cdc.MustUnmarshalJSON(bankGenStateJson, &bankGenesis)

		foundSubaccountsBalance := false
		subaccountsAddress := types.ModuleAddress.String()

		for _, balance := range bankGenesis.Balances {
			if balance.Address == subaccountsAddress {
				areBalancesEqual := totalTDaiSupply.Equal(balance.Coins[0].Amount)
				require.True(t, areBalancesEqual)
				foundSubaccountsBalance = true
				break
			}
		}

		require.True(t, foundSubaccountsBalance)
		foundTDai := false

		for _, supply := range bankGenesis.Supply {
			if supply.Denom == asstypes.AssetTDai.Denom {
				isSupplyEqual := totalTDaiSupply.Equal(supply.Amount)
				require.True(t, isSupplyEqual)
				foundTDai = true
				break
			}
		}

		require.True(t, foundTDai)
	}
}
