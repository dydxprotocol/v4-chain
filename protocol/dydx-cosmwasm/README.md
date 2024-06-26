# dYdX Bindings for Cosmwasm Contracts

dYdX specific bindings for cosmwasm contract to be able to interact with the dYdX Chain by exposing custom messages, queries, and structs that correspond to custom module functionality.

# Integration
Use this package as a dependency in your cosmwasm contract to interact with the dYdX chain.
This package has bindings to send messages and queries to the dYdX chain.

### Queries
dYdX exposes 3 custom queries
1. Getting market price for a given asset
2. Getting Subaccount details
3. Getting details on the perpetual and clob pair

the json payloads for these queries are defined [here](https://github.com/dydxprotocol/v4-chain/blob/feature/cosmwasm/protocol/dydx-cosmwasm/src/query.rs#L22-L33). The functions which are called are defined [here](https://github.com/dydxprotocol/v4-chain/blob/feature/cosmwasm/protocol/dydx-cosmwasm/src/querier.rs#L18-L54)

### Messages
dYdx exposes 4 custom messages
1. DepositToSubaccount
2. WithdrawFromSubaccount
3. PlaceOrder
4. CancelOrder

The json payloads for custom messages dYdX exposes are defined [here](https://github.com/dydxprotocol/v4-chain/blob/feature/cosmwasm/protocol/dydx-cosmwasm/src/msg.rs#L79-L102)

For an example contract that uses these bindings, see [here](https://github.com/dydxprotocol/v4-chain/tree/feature/cosmwasm/protocol/contracts/dydx-messages-example)

# Compilation
To compile the contract, run the following command, follow the instructions in this [doc](https://docs.cosmwasm.com/docs/getting-started/compile-contract). This will create a .wasm file that can be deployed to the dYdX chain.


# Testing
To test the contract, complete the following steps:
1. Follow Getting Started instructions in this [doc](https://github.com/dydxprotocol/v4-chain/tree/main/protocol#get-started) to set up the environment.
2. Build the dYdX binary using `make build` in the `v4-chain/protocol` directory. This will create a dYdX binary in the `v4-chain/protocol/build/` directory called `dydxprotocold`.
3. Run the chain locally using documentation [here](https://github.com/dydxprotocol/v4-chain/tree/main/protocol#running-the-chain-locally)
4. To deploy the contract, run the following command:
```bash
./build/dydxprotocold tx wasm store /path/to/wasm/binary --from alice --gas-prices 25000000000adv4tnt --gas auto --gas-adjustment 1.5 --chain-id localdydxprotocol
```
5. Check the contract code ID using the following command:
```bash
./build/dydxprotocold query wasm list-code
```
6. After the contract is deployed, instantiate the contract using the following command:
```bash
./build/dydxprotocold tx wasm instantiate 1 '<instantiation parameters>' --from alice --label test --gas-prices 25000000000adv4tnt --gas auto --gas-adjustment 1.5 --chain-id localdydxprotocol
```
7. To query the contract, use the following command:
```bash
./build/dydxprotocold query wasm contract-state smart <contract_address> <query_msg> --chain-id localdydxprotocol
```
8. To execute a contract message, use the following command:
```bash
./build/dydxprotocold tx wasm execute <contract_address> <execute_msg> --from alice --gas-prices 25000000000adv4tnt --gas auto --gas-adjustment 1.5 --chain-id localdydxprotocol
```

## Getting funds for local testing
For local testing, use [faucet](https://github.com/dydxprotocol/faucet) to get funds. To run faucet on a local chain, follow instructions below:
1. Run v4 chain locally by following the [v4 instructions](https://github.com/dydxprotocol/v4/blob/main/README.md#running-the-chain-locally) (i.e. `cd ~/v4 && make localnet-start`).
2. Make sure that `faucet` can be built by running `npm run build` or `npm run compile:watch`.
3. Run `npm run localnet`.
4. You can try hitting the faucet endpoint using `curl`. For example:
```
curl -X POST http://localhost:1234/faucet/tokens -H "Content-Type: application/json" -d '{"address": "dydx1nzuttarf5k2j0nug5yzhr6p74t9avehn9hlh8m", "amount": 200, "subaccountNumber": 0}
```

## Load testing on localnet
To run load testing on localnet, use[load tester](https://github.com/dydxprotocol/load-tester). Follow instructions [here](https://github.com/dydxprotocol/load-tester?tab=readme-ov-file#running-against-your-local-chain-localnet)





