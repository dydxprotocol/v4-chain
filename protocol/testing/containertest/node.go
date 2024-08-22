package containertest

import (
	"context"
	"fmt"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/app"

	comethttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	gogogrpc "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	"github.com/dydxprotocol/v4-chain/protocol/cmd/dydxprotocold/cmd"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/ory/dockertest/v3"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// When polling a node, poll every `pollFrequencyNs` and give up after `pollAttempts` attempts.
	pollFrequencyNs = time.Second
	pollAttempts    = 60
	cometPort       = "26657/tcp"
	grpcPort        = "9090/tcp"
)

type Node struct {
	keyring   *keyring.Keyring
	cometPort string
	grpcPort  string
	resource  *dockertest.Resource
}

func newNode(keyring *keyring.Keyring, resource *dockertest.Resource) (node *Node, err error) {
	node = &Node{
		keyring:  keyring,
		resource: resource,
	}
	node.cometPort = resource.GetHostPort(cometPort)
	node.grpcPort = resource.GetHostPort(grpcPort)
	return node, err
}

func (n *Node) createCometClient() (*comethttp.HTTP, error) {
	return comethttp.New("tcp://"+n.cometPort, "/websocket")
}

func (n *Node) createGrpcConn() (*grpc.ClientConn, error) {
	return grpc.Dial(n.grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials())) //nolint:staticcheck
}

// Wait for current block height has advanced by at least `numBlocks`
func (n *Node) Wait(numBlocks int64) error {
	currBlock, err := n.LatestBlockHeight()
	if err != nil {
		return err
	}
	return n.WaitUntilBlockHeight(currBlock + numBlocks)
}

// Wait for the current block height to reach at least `height`
func (n *Node) WaitUntilBlockHeight(height int64) error {
	for i := 0; i < pollAttempts; i++ {
		latestHeight, err := n.LatestBlockHeight()

		if err == nil && latestHeight >= height {
			return nil
		}
		time.Sleep(pollFrequencyNs)
	}
	return fmt.Errorf("timed out when waiting for block height %d", height)
}

func (n *Node) LatestBlockHeight() (int64, error) {
	cometClient, err := n.createCometClient()
	if err != nil {
		return 0, err
	}
	status, err := cometClient.Status(context.Background())
	if err != nil {
		return 0, err
	}
	return status.SyncInfo.LatestBlockHeight, nil
}

// Get a `Context` to be used for broadcasting tx given a signer address.
// NB: The cosmos client for broadcasting tx seems to be extremely coupled with command, so we have to create a dummy
// command to hijack some of the private functionality to create the context. There may be a better way to get Context,
// such as going a level lower to the tx.Factory level, but this works for now.
func (n *Node) getContextForBroadcastTx(signer string) (*client.Context, *pflag.FlagSet, error) {
	initClientCtx := client.Context{}.
		WithCodec(constants.TestEncodingCfg.Codec).
		WithInterfaceRegistry(constants.TestEncodingCfg.InterfaceRegistry).
		WithTxConfig(constants.TestEncodingCfg.TxConfig).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithViper(cmd.EnvPrefix)

	option := cmd.GetOptionWithCustomStartCmd()
	rootCmd := cmd.NewRootCmd(option, app.DefaultNodeHome)
	flags.AddTxFlagsToCmd(rootCmd)
	flags := rootCmd.Flags()

	hostPort := n.resource.GetHostPort(cometPort)
	if err := flags.Set("node", fmt.Sprintf("tcp://%s", hostPort)); err != nil {
		return nil, nil, err
	}
	if err := flags.Set("from", signer); err != nil {
		return nil, nil, err
	}
	if err := flags.Set("chain-id", "localdydxprotocol"); err != nil {
		return nil, nil, err
	}

	// NB: In `cmd/dydxprotocol/root.go` this step is done before ReadFromClientConfig, but here we choose to
	// do it second because ReadPersistentCommandFlags sets the node address we configured in flags.
	// If we were to do it in reverse, ReadFromClientConfig would overwrite the node address.
	initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, flags)
	if err != nil {
		return nil, nil, err
	}

	initClientCtx = initClientCtx.WithKeyring(*n.keyring).WithSkipConfirmation(true)
	ctx := context.WithValue(context.Background(), client.ClientContextKey, &initClientCtx)
	rootCmd.SetContext(ctx)

	initClientCtx, err = client.GetClientTxContext(rootCmd)
	if err != nil {
		return nil, nil, err
	}
	return &initClientCtx, flags, nil
}

// Broadcast a tx to the node given the message and a signer address.
func BroadcastTx[M cosmos.Msg](n *Node, message M, signer string) (err error) {
	clientContext, flags, err := n.getContextForBroadcastTx(signer)
	if err != nil {
		return err
	}

	txFactory, err := tx.NewFactoryCLI(*clientContext, flags)
	if err != nil {
		return err
	}

	// Use default gas limit and gas fee.
	txFactory = txFactory.WithGas(constants.TestGasLimit).WithFees(constants.TestFee)

	if err = tx.GenerateOrBroadcastTxWithFactory(*clientContext, txFactory, message); err != nil {
		return err
	}
	return nil
}

// Broadcast a tx to the node given the message and a signer address.
func BroadcastTxWithoutValidateBasic[M cosmos.Msg](n *Node, message M, signer string) (err error) {
	clientContext, flags, err := n.getContextForBroadcastTx(signer)
	if err != nil {
		return err
	}

	txFactory, err := tx.NewFactoryCLI(*clientContext, flags)
	if err != nil {
		return err
	}

	// Use default gas limit and gas fee.
	txFactory = txFactory.WithGas(constants.TestGasLimit).WithFees(constants.TestFee)

	if err = tx.BroadcastTx(*clientContext, txFactory, message); err != nil {
		return err
	}
	return nil
}

// Query the node's grpc endpoint given the client constructor, request method, and request
func Query[Request proto.Message, Response proto.Message, Client interface{}](
	n *Node,
	clientConstructor func(gogogrpc.ClientConn) Client,
	requestFn func(Client, context.Context, Request, ...grpc.CallOption) (Response, error),
	request Request) (proto.Message, error) {
	conn, err := n.createGrpcConn()
	if err != nil {
		return nil, err
	}
	client := clientConstructor(conn)
	return requestFn(client, context.Background(), request)
}
