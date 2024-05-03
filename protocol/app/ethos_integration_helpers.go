package app

import (
	storetypes "cosmossdk.io/store/types"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	clobkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/keeper"
	perpetualskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/keeper"
	"github.com/cosmos/cosmos-sdk/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	ibctestingtypes "github.com/cosmos/ibc-go/v8/testing/types"
	testutil "github.com/ethos-works/ethos/ethos-chain/testutil/integration"
	ibcconsumerkeeper "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/keeper"
)

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// ConsumerApp interface implementations for integration tests

// GetConsumerKeeper implements the ConsumerApp interface.
func (app *App) GetConsumerKeeper() ibcconsumerkeeper.Keeper {
	return app.ConsumerKeeper
}

// GetTestBankKeeper implements the ConsumerApp interface.
func (app *App) GetTestBankKeeper() testutil.TestBankKeeper {
	return app.BankKeeper
}

// GetTestAccountKeeper implements the ConsumerApp interface.
func (app *App) GetTestAccountKeeper() testutil.TestAccountKeeper {
	return app.AccountKeeper
}

// GetTestSlashingKeeper implements the ConsumerApp interface.
func (app *App) GetTestSlashingKeeper() testutil.TestSlashingKeeper {
	return app.SlashingKeeper
}

// GetTestEvidenceKeeper implements the ConsumerApp interface.
func (app *App) GetTestEvidenceKeeper() evidencekeeper.Keeper {
	return app.EvidenceKeeper
}

func (app *App) GetTestClobKeeper() clobkeeper.Keeper {
	return *app.ClobKeeper
}

// GetIBCKeeper implements the TestingApp interface.
func (app *App) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

// GetScopedIBCKeeper implements the TestingApp interface.
func (app *App) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

// GetTxConfig implements the TestingApp interface.
func (app *App) GetTxConfig() client.TxConfig {
	return app.txConfig
}

// GetStakingKeeper implements the TestingApp interface.
func (app *App) GetStakingKeeper() ibctestingtypes.StakingKeeper {
	return app.ConsumerKeeper
}

func (app *App) GetTestPerpetualsKeeper() perpetualskeeper.Keeper {
	return *app.PerpetualsKeeper
}
