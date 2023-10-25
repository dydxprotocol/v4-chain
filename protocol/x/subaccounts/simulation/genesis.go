package simulation

import (
	"math/rand"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sim_helpers"
	asstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// genSubaccountIdNumbers returns randomized slice of numbers to use for `Subaccount.SubaccountId.Number`.
func genSubaccountIdNumbers(r *rand.Rand) []uint32 {
	allSubaccountNums := sim_helpers.MakeRange(uint32(sim_helpers.MaxNumSubaccount))
	randomizedSubaccountNums := sim_helpers.RandSliceShuffle(r, allSubaccountNums)
	numSubaccounts := simtypes.RandIntBetween(r, sim_helpers.MinNumSubaccount, sim_helpers.MaxNumSubaccount+1)
	return randomizedSubaccountNums[:numSubaccounts]
}

// RandomizedGenState generates a random GenesisState for `Subaccounts`.
func RandomizedGenState(simState *module.SimulationState) {
	// TODO(DEC-1049): update genesis state for other modules (i.e. bank and auth)
	// so that invariant (i.e. subaccounts module account balance of USDC should
	// always be > than total net collateral of all Subaccounts) is respected.

	r := simState.Rand

	// For each simulator account, create an associated subaccount.
	allSubaccounts := make([]types.Subaccount, 0)

	// Define the total USDC supply as the sum of all USDC quantums in all subaccounts.
	totalUsdcSupply := sdkmath.NewInt(0)

	for _, acc := range simState.Accounts {
		saIdNumbers := genSubaccountIdNumbers(r)

		// For each subaccount id, associate random assets and perpetuals.
		for _, idNum := range saIdNumbers {
			subacct := types.Subaccount{
				Id: &types.SubaccountId{
					Owner:  acc.Address.String(),
					Number: idNum,
				},
			}

			// Toggle adding USDC asset position.
			if sim_helpers.RandBool(r) {
				quantums := r.Uint64()
				subacct.AssetPositions = []*types.AssetPosition{
					{
						AssetId:  asstypes.AssetUsdc.Id,
						Quantums: dtypes.NewIntFromUint64(quantums),
					},
				}

				bigQuantums := sdkmath.NewIntFromUint64(quantums)
				totalUsdcSupply = totalUsdcSupply.Add(bigQuantums)
			}

			// Purposely do NOT add PerpetualPositions. These positions should be created naturally
			// as orders are placed and matched via weighted operations. In order to add these
			// PerpetualPositions as part of genesis, we need to ensure that the short/long
			// PerpetualPositions are perfectly balanced.

			// TODO(DEC-582): randomly toggle `MarginEnabled` once we support margin trading.

			allSubaccounts = append(allSubaccounts, subacct)
		}
	}

	subaccountsGenesis := types.GenesisState{
		Subaccounts: allSubaccounts,
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&subaccountsGenesis)

	updateBankModuleGenesisState(simState, totalUsdcSupply)
}

// updateBankModuleGenesisState updates the bank module's genesis state by
// assigning the total supply of USDC to the balance of the `subaccounts` module.
// This is necessary as the protocol assumes that that the sum of quantums in all USDC
// AssetPositions is <= the total USDC balance of the subaccounts module, and `panic`s
// will occur when transferring fees to the `fee-collector` module during order processing
// if this is not true.
// This method assumes that USDC as a `Coin` in the bank module does not yet exist.
func updateBankModuleGenesisState(
	simState *module.SimulationState,
	totalUsdcSupply sdkmath.Int,
) {
	var bankGenesis banktypes.GenesisState
	bankGenStateJson := simState.GenState[banktypes.ModuleName]
	simState.Cdc.MustUnmarshalJSON(bankGenStateJson, &bankGenesis)

	// Define the balance of the `subaccounts` module.
	subaccountsUsdcBalance := banktypes.Balance{
		Address: types.ModuleAddress.String(),
		Coins: []sdk.Coin{{
			Denom:  asstypes.AssetUsdc.Denom,
			Amount: totalUsdcSupply,
		}},
	}

	// Set the balance of the `subaccounts` module on the bank genesis.
	bankGenesis.Balances = append(bankGenesis.Balances, subaccountsUsdcBalance)

	// Set the total supply of USDC on the bank genesis.
	bankGenesis.Supply = append(bankGenesis.Supply,
		sdk.NewCoin(asstypes.AssetUsdc.Denom, totalUsdcSupply),
	)

	// Update the bank module's genesis state.
	simState.GenState[banktypes.ModuleName] = simState.Cdc.MustMarshalJSON(&bankGenesis)
}
