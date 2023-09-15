package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

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

	go run scripts/bridge_events.go \
		-denom <token_denom> \
		-rpc <rpc_node_url> \
		-address <bridge_contract_address> \
		-toblock <last_block_inclusive>
*/
func main() {
	ctx := context.Background()

	// ------------ FLAGS ------------

	// Get flags.
	var denom, rpcNode, bridgeAddress string
	var toBlock int64
	flag.StringVar(&denom, "denom", "dv4tnt", "token denom")
	flag.StringVar(&rpcNode, "rpc", "https://eth-sepolia.g.alchemy.com/v2/demo", "rpc node url")
	flag.StringVar(&bridgeAddress, "address", "0xcca9D5f0a3c58b6f02BD0985fC7F9420EA24C1f0", "bridge address")
	flag.Int64Var(&toBlock, "toblock", 100_000_000, "last block (inclusive)")
	flag.Parse()

	// Print the flags used for the user
	fmt.Println("Using the following configuration (modifiable via flags):")
	fmt.Println("denom:", denom)
	fmt.Println("rpc:", rpcNode)
	fmt.Println("address:", bridgeAddress)
	fmt.Println("toblock:", toBlock)
	fmt.Println()

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

	// Iterate over each event and populate the above fields
	for _, log := range logs {
		event := libeth.BridgeLogToEvent(log, denom)
		aei.NextId = lib.Max(aei.NextId, event.Id)
		aei.EthBlockHeight = lib.Max(aei.EthBlockHeight, log.BlockNumber)

		// Get the current balance of the account.
		v, exists := balances[event.Address]
		if !exists {
			v = big.NewInt(0)
		}

		// Store the new balance in the balances map.
		balances[event.Address] = v.Add(v, event.Coin.Amount.BigInt())
	}

	// ------------ OUTPUT ------------

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
		newBalance := balances[address]
		bankBalances = append(bankBalances, banktypes.Balance{
			Address: address,
			Coins:   sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromBigInt(newBalance))),
		})
	}
	fmt.Printf("\"bank.balances\": %s\n", mustJson(bankBalances))
}

// mustJson marshals v into indented JSON and exits the program if it fails to do so.
func mustJson(v any) string {
	output, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(output)
}
