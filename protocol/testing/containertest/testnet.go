package containertest

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/testexchange"
	pricefeed "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/testing/version"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	pricefeed_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/pricefeed"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// Resources will expire in 10 minutes
const resourceLifetimeSecs = 600

func monikers() map[string]string {
	return map[string]string{
		"alice": constants.AliceMnenomic,
		"bob":   constants.BobMnenomic,
		"carl":  constants.CarlMnenomic,
		"dave":  constants.DaveMnenomic,
	}
}

type Testnet struct {
	Nodes map[string]*Node

	isPreupgradeGenesis bool
	keyring             keyring.Keyring
	pool                *dockertest.Pool
	network             *dockertest.Network
	exchangeServer      *pricefeed_testutil.ExchangeServer
	uniqueId            string
}

// NewTestnet returns a new Testnet. If creation fails, an error is returned.
// In some cases, resources could be initialized but not properly cleaned up. The error will reflect this.
func NewTestnet() (testnet *Testnet, err error) {
	// Generate unique ID for this testnet instance to avoid conflicts
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return nil, err
	}
	uniqueId := fmt.Sprintf("%d", n.Int64())

	testnet = &Testnet{
		Nodes:    make(map[string]*Node),
		keyring:  keyring.NewInMemory(constants.TestEncodingCfg.Codec),
		uniqueId: uniqueId,
	}
	testnet.pool, err = dockertest.NewPool("")
	if err != nil {
		return nil, err
	}
	err = testnet.pool.Client.Ping()
	if err != nil {
		return nil, err
	}

	for moniker, mnemonic := range monikers() {
		_, err = testnet.keyring.NewAccount(moniker, mnemonic, "", sdk.GetConfig().GetFullBIP44Path(), hd.Secp256k1)
		if err != nil {
			return nil, err
		}
	}

	testnet.exchangeServer = pricefeed_testutil.NewExchangeServer()
	return testnet, nil
}

func NewTestnetWithPreupgradeGenesis() (testnet *Testnet, err error) {
	testnet, err = NewTestnet()
	if err != nil {
		return nil, err
	}
	testnet.isPreupgradeGenesis = true
	return testnet, err
}

func (t *Testnet) Start() (err error) {
	err = t.initialize()
	if err != nil {
		cleanUpErr := t.CleanUp()
		if cleanUpErr == nil {
			return fmt.Errorf("testnet initialization failed with error: %w; resources successfully cleaned up", err)
		}
		return fmt.Errorf(
			"testnet initialization failed with error: %s; failed to clean-up resources with error %s",
			err,        //nolint:errorLint
			cleanUpErr) //nolint:errorLint
	}
	return nil
}

// initialize sets up all state that needs to be cleaned up. Returns error immediately upon a failure.
func (t *Testnet) initialize() (err error) {
	// Clean up any existing containers/networks with the same names first
	t.cleanupExistingResources()

	// NB: Docker lets you create multiple networks with the same name. ID, however, is unique.
	// Consider not using the same name in the future if it proves to be a problem.
	networkName := fmt.Sprintf("test-network-%s", t.uniqueId)
	t.network, err = t.pool.CreateNetwork(networkName)
	if err != nil {
		return err
	}

	for moniker := range monikers() {
		node, err := t.initializeNode(moniker)
		if err != nil {
			return err
		}
		t.Nodes[moniker] = node
	}

	for moniker, node := range t.Nodes {
		if err := t.pool.Retry(func() error {
			return node.WaitUntilBlockHeight(2)
		}); err != nil {
			return fmt.Errorf(
				"could not connect to node: %s, %w",
				moniker,
				err,
			)
		}
	}
	return nil
}

// cleanupExistingResources removes any existing containers/networks that might conflict
func (t *Testnet) cleanupExistingResources() {
	// Try to remove containers with our unique names if they exist
	for moniker := range monikers() {
		containerName := fmt.Sprintf("testnet-local-%s-%s", moniker, t.uniqueId)
		if container, err := t.pool.Client.InspectContainer(containerName); err == nil {
			_ = t.pool.Client.RemoveContainer(docker.RemoveContainerOptions{
				ID:    container.ID,
				Force: true,
			})
		}
	}

	// Try to remove network if it exists
	networkName := fmt.Sprintf("test-network-%s", t.uniqueId)
	if networks, err := t.pool.Client.ListNetworks(); err == nil {
		for _, network := range networks {
			if network.Name == networkName {
				_ = t.pool.Client.RemoveNetwork(network.ID)
				break
			}
		}
	}
}

func (t *Testnet) initializeNode(moniker string) (*Node, error) {
	var entrypointCommand string
	if t.isPreupgradeGenesis {
		entrypointCommand = "/dydxprotocol/preupgrade_entrypoint.sh"
	} else {
		entrypointCommand = "dydxprotocold"
	}

	// Generate dynamic persistent peers using the unique ID
	persistentPeers := fmt.Sprintf(
		"17e5e45691f0d01449c84fd4ae87279578cdd7ec@testnet-local-alice-%s:26656,"+
			"b69182310be02559483e42c77b7b104352713166@testnet-local-bob-%s:26656,"+
			"47539956aaa8e624e0f1d926040e54908ad0eb44@testnet-local-carl-%s:26656,"+
			"5882428984d83b03d0c907c1f0af343534987052@testnet-local-dave-%s:26656",
		t.uniqueId, t.uniqueId, t.uniqueId, t.uniqueId)

	resource, err := t.pool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       fmt.Sprintf("testnet-local-%s-%s", moniker, t.uniqueId),
			Repository: "dydxprotocol-container-test",
			Tag:        "",
			NetworkID:  t.network.Network.ID,
			ExposedPorts: []string{
				"26657/tcp",
				"9090/tcp",
			},
			Entrypoint: []string{
				entrypointCommand,
				"start",
				"--home",
				fmt.Sprintf("/dydxprotocol/chain/.%s", moniker),
				"--p2p.persistent_peers",
				persistentPeers,
				"--bridge-daemon-eth-rpc-endpoint",
				"https://eth-sepolia.g.alchemy.com/v2/demo",
			},
			Env: []string{
				"DAEMON_NAME=dydxprotocold",
				fmt.Sprintf("DAEMON_HOME=/dydxprotocol/chain/.%s", moniker),
				fmt.Sprintf("UPGRADE_TO_VERSION=%s", version.CurrentVersion),
			},
			ExtraHosts: []string{
				fmt.Sprintf("%s:host-gateway", testexchange.TestExchangeHost),
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if err = resource.Expire(resourceLifetimeSecs); err != nil {
		return nil, err
	}

	return newNode(&t.keyring, resource)
}

func (t *Testnet) removeNode(moniker string) error {
	if err := t.pool.Purge(t.Nodes[moniker].resource); err != nil {
		return err
	}
	delete(t.Nodes, moniker)
	return nil
}

// CleanUp cleans up any resources used by a Testnet. This should always be called to avoid leaking docker resources.
func (t *Testnet) CleanUp() error {
	var err error
	errs := []string{}
	for moniker := range t.Nodes {
		if err := t.removeNode(moniker); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if t.network != nil {
		if err = t.pool.RemoveNetwork(t.network); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if t.exchangeServer != nil {
		if err = t.exchangeServer.CleanUp(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("cleanup failed with error(s): %s", strings.Join(errs, ";"))
	}
	return nil
}

func (t *Testnet) MustCleanUp() {
	err := t.CleanUp()
	if err != nil {
		panic(fmt.Errorf("testnet failed to clean up: %w", err))
	}
}

func (t *Testnet) setPrice(marketId pricefeed.MarketId, price float64) {
	t.exchangeServer.SetPrice(marketId, price)
}
