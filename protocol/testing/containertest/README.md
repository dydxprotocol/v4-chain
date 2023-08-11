# Container Tests

The `containertest` package provides a testing framework for testing the network as a whole. It allows writing unit tests that start up a network of nodes. Some interfaces are provided to interact with the chain.

At a high level, each unit test starts docker containers that are each running a node. The test framework interacts with these nodes by implementing functions that wrap CometBFT and gRPC clients for the node. The framework also runs an HTTP server as a goroutine, and the network is configured to use that server as the source for exchange prices.

## Running the tests

The following make targets live in `v4` root:
 - `make test-container-build`: Build the docker image needed to run the tests.
 - `make test-container`: Run the tests.
 - `make test-container-accept`: Run the tests and accept any differing expected output.

## Writing tests

Tests should use this general structure:
```
testnet, err := NewTestnet()
require.NoError(t, err)

// Do things prior to network start if needed, such as setting prices.

// Start the network
err = testnet.Start()
require.NoError(t, err)
defer testnet.MustCleanUp()

// Interact with the chain
node := testnet.Nodes["alice"]
Query(node, ...)
BroadcastTx(node, ...)
```

Use the following methods to interact with the chain. See `testnet.go` and `node.go` for details:
 - `Node.BroadcastTx`: Broadcast a tx to a node.
 - `Node.Query`: Send a request a node's query endpoint.
 - `Testnet.SetPrice`: Set the external exchange price of a market. It's recommended to do this before starting the network because the oracle price uses a moving window.

 See `testnet_test.go` for examples.
