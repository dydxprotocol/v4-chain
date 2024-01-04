package main

import (
	"context"
	sdkmath "cosmossdk.io/math"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	libeth "github.com/dydxprotocol/v4-chain/protocol/lib/eth"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"

	eth "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*
main runs an rpc query against an Ethereum node to get any relevant bridge events and then prints out
relevant state that would be modified as a result of these events.

Usage:

	go run scripts/bridge_events/bridge_events.go \
		-denom <token_denom> \
		-totalsupply <total_supply> \
		-rpc <rpc_node_url> \
		-address <bridge_contract_address> \
		-toblock <last_block_inclusive>
*/
func main() {
	ctx := context.Background()

	// ------------ FLAGS ------------

	// Get flags.
	var denom, totalSupply, rpcNode, bridgeAddress string
	var toBlock int64
	var verbose bool
	flag.StringVar(&denom, "denom", "adv4tnt", "token denom")
	flag.StringVar(&totalSupply, "totalsupply", "1000000000000000000000000000", "token's total supply (base 10)")
	flag.StringVar(&rpcNode, "rpc", "https://eth-sepolia.g.alchemy.com/v2/demo", "rpc node url")
	flag.StringVar(&bridgeAddress, "address", "0xcca9D5f0a3c58b6f02BD0985fC7F9420EA24C1f0", "bridge address")
	flag.Int64Var(&toBlock, "toblock", 100_000_000, "last block (inclusive)")
	flag.BoolVar(&verbose, "verbose", false, "print additional JSON")
	flag.Parse()

	// Print the flags used for the user
	fmt.Println("Using the following configuration (modifiable via flags):")
	fmt.Println("denom:", denom)
	fmt.Println("totalsupply:", totalSupply)
	fmt.Println("rpc:", rpcNode)
	fmt.Println("address:", bridgeAddress)
	fmt.Println("toblock:", toBlock)
	fmt.Println()

	// ------------ INPUT VALIDATION ------------
	// Validate `-totalsupply`.
	totalSupplyBigInt, ok := new(big.Int).SetString(totalSupply, 10)
	if !ok {
		log.Fatal("invalid total supply")
	}

	// ------------ LOGIC ------------

	// Create client.
	ethClient, err := ethclient.Dial(rpcNode)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { ethClient.Close() }()

	// Get chain ID from the Ethereum Node.
	chainId, err := ethClient.ChainID(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch logs from Ethereum Node.
	filterQuery := eth.FilterQuery{
		FromBlock: nil,
		ToBlock:   big.NewInt(toBlock),
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(bridgeAddress)},
		Topics: [][]ethcommon.Hash{
			{ethcommon.HexToHash(constants.BridgeEventSignature)},
		},
	}
	logs, err := ethClient.FilterLogs(ctx, filterQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Remember the latest event info and the total balance to credit each account.
	aei := bridgetypes.BridgeEventInfo{
		NextId:         0,
		EthBlockHeight: 0,
	}
	balances := make(map[string]*big.Int)
	totalAmountBridged := big.NewInt(0)

	// Iterate over each event and populate the above fields
	for _, log := range logs {
		event := libeth.BridgeLogToEvent(log, denom)
		aei.NextId = lib.Max(aei.NextId, event.Id+1)
		aei.EthBlockHeight = lib.Max(aei.EthBlockHeight, log.BlockNumber)

		// Add amount to total bridged amount.
		totalAmountBridged.Add(totalAmountBridged, event.Coin.Amount.BigInt())

		// Get the current balance of the account.
		v, exists := balances[event.Address]
		if !exists {
			v = big.NewInt(0)
		}

		// Store the new balance in the balances map.
		balances[event.Address] = v.Add(v, event.Coin.Amount.BigInt())
	}

	// ------------ OUTPUT ------------

	cdc := app.GetEncodingConfig().Codec

	// Print total amount bridged.
	fmt.Printf("Total amount bridged: %s\n", totalAmountBridged.String())

	// Print bridge module account remaining balance.
	bridgeModAccBalance := totalSupplyBigInt.Sub(totalSupplyBigInt, totalAmountBridged)
	fmt.Printf("Remaining bridge module account balance: %s\n", bridgeModAccBalance.String())

	// Print x/bridge event params.
	eventParams := bridgetypes.EventParams{
		Denom:      denom,
		EthChainId: chainId.Uint64(),
		EthAddress: bridgeAddress,
	}
	fmt.Printf("\"bridge.event_params\": %s\n", mustJson(eventParams))

	// Print x/bridge acknowledged event info.
	fmt.Printf("\"bridge.acknowledged_event_info\": %s\n", mustJson(aei))

	// Print x/bank balances sorted by address.
	bankBalances := make([]banktypes.Balance, 0)
	sortedAddresses := lib.GetSortedKeys[sort.StringSlice](balances)
	for _, address := range sortedAddresses {
		bankBalances = append(
			bankBalances,
			banktypes.Balance{
				Address: address,
				Coins:   sdk.NewCoins(sdk.NewCoin(denom, sdkmath.NewIntFromBigInt(balances[address]))),
			},
		)
	}
	bankGenesisJson := cdc.MustMarshalJSON(banktypes.NewGenesisState(
		banktypes.Params{},
		bankBalances,
		sdk.NewCoins(),
		[]banktypes.Metadata{},
		[]banktypes.SendEnabled{},
	))
	fmt.Printf("\"bank.balances\": %s\n", extractFieldFromJson(bankGenesisJson, "balances"))

	// Stop here if not verbose.
	if !verbose {
		return
	}

	// Print x/auth accounts information.
	genesisAccounts := authtypes.GenesisAccounts{}
	for i, address := range sortedAddresses {
		var ba sdk.AccountI = &authtypes.BaseAccount{
			Address:       address,
			PubKey:        nil,
			AccountNumber: uint64(i),
			Sequence:      uint64(0),
		}
		genesisAccounts = append(
			genesisAccounts,
			ba.(authtypes.GenesisAccount),
		)
	}
	authGenesisJson := cdc.MustMarshalJSON(authtypes.NewGenesisState(
		authtypes.Params{},
		genesisAccounts,
	))
	fmt.Printf("\"auth.accounts\": %s\n", extractFieldFromJson(authGenesisJson, "accounts"))
}

// extractFieldFromJson takes a JSON dictionary as an input and returns the json of a single field.
func extractFieldFromJson(input []byte, field string) string {
	var v map[string]interface{}
	err := json.Unmarshal(input, &v)
	if err != nil {
		log.Fatal(err)
	}
	return mustJson(v[field])
}

// mustJson marshals v into indented JSON and exits the program if it fails to do so.
func mustJson(v any) string {
	output, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(output)
}
