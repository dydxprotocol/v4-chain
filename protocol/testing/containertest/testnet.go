package containertest

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/testexchange"
	pricefeed "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	pricefeed_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/pricefeed"
	"github.com/ory/dockertest/v3"
)

// For now all this config data like peers and monikers are hard coded to match the local net.
// In the future we'll pull stuff from a config.
const persistentPeers = "17e5e45691f0d01449c84fd4ae87279578cdd7ec@testnet-local-alice:26656," +
	"b69182310be02559483e42c77b7b104352713166@testnet-local-bob:26656," +
	"47539956aaa8e624e0f1d926040e54908ad0eb44@testnet-local-carl:26656," +
	"5882428984d83b03d0c907c1f0af343534987052@testnet-local-dave:26656"

// Resources will expire in 10 minutes
const resourceLifetimeSecs = 600

// The version of that we're upgrading to (aka the current commit)
const UpgradeToVersion = "v3.0.0"

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
}

// NewTestnet returns a new Testnet. If creation fails, an error is returned.
// In some cases, resources could be initialized but not properly cleaned up. The error will reflect this.
func NewTestnet() (testnet *Testnet, err error) {
	testnet = &Testnet{
		Nodes:   make(map[string]*Node),
		keyring: keyring.NewInMemory(constants.TestEncodingCfg.Codec),
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
	// NB: Docker lets you create multiple networks with the same name. ID, however, is unique.
	// Consider not using the same name in the future if it proves to be a problem.
	t.network, err = t.pool.CreateNetwork("test-network")
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
			return fmt.Errorf("could not connect to node: %s", moniker)
		}
	}
	return nil
}

func (t *Testnet) initializeNode(moniker string) (*Node, error) {
	var entrypointCommand string
	if t.isPreupgradeGenesis {
		entrypointCommand = "/dydxprotocol/preupgrade_entrypoint.sh"
	} else {
		entrypointCommand = "dydxprotocold"
	}

	resource, err := t.pool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       fmt.Sprintf("testnet-local-%s", moniker),
			Repository: "dydxprotocol-container-test",
			Tag:        "",
			NetworkID:  t.network.Network.ID,
			ExposedPorts: []string{
				"26657",
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
				fmt.Sprintf("UPGRADE_TO_VERSION=%s", UpgradeToVersion),
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
