package v_8_0_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	v_8_0 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v8.0"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	accountplustypes "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	"github.com/stretchr/testify/suite"
)

type UpgradeTestSuite struct {
	suite.Suite

	tApp *testapp.TestApp
	Ctx  sdk.Context
}

func TestMigrateAccountplusAccountState(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

func (s *UpgradeTestSuite) SetupTest() {
	s.tApp = testapp.NewTestAppBuilder(s.T()).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *accountplustypes.GenesisState) {
				genesisState.Params.IsSmartAccountActive = false
			},
		)
		return genesis
	}).Build()
	s.Ctx = s.tApp.InitChain()
}

func (s *UpgradeTestSuite) TestUpgrade_MigrateAccountplusAccountState() {
	ctx := s.Ctx
	store := ctx.KVStore(s.tApp.App.AccountPlusKeeper.GetStoreKey())
	prefixStore := prefix.NewStore(store, []byte(accountplustypes.AccountStateKeyPrefix))

	// Create some AccountState with no prefixes
	addresses := []string{"address1", "address2", "address3"}
	for _, addr := range addresses {
		accAddress := sdk.AccAddress([]byte(addr))
		accountState := accountplustypes.AccountState{
			Address: addr,
			TimestampNonceDetails: accountplustypes.TimestampNonceDetails{
				TimestampNonces: []uint64{1, 2, 3},
				MaxEjectedNonce: 0,
			},
		}
		bz := s.tApp.App.AccountPlusKeeper.GetCdc().MustMarshal(&accountState)
		store.Set(accAddress, bz)
	}

	// Verify unprefixed keys were successfully created
	for _, addr := range addresses {
		accAddress := sdk.AccAddress([]byte(addr))
		bz := store.Get(accAddress)
		s.Require().NotNil(bz, "Unprefixed key not created for %s", addr)
	}

	// Migrate
	v_8_0.MigrateAccountplusAccountState(ctx, s.tApp.App.AccountPlusKeeper)

	// Verify that unprefixed keys are deleted and prefixed keys exist
	for _, addr := range addresses {
		accAddress := sdk.AccAddress([]byte(addr))

		// Check that the old key no longer exists
		bzOld := store.Get(accAddress)
		s.Require().Nil(bzOld, "Unprefixed AccountState should be deleted for %s", addr)

		// Check that the new prefixed key exists
		bzNew := prefixStore.Get(accAddress)
		var actualAccountState accountplustypes.AccountState
		s.tApp.App.AccountPlusKeeper.GetCdc().MustUnmarshal(bzNew, &actualAccountState)
		expectedAccountState := accountplustypes.AccountState{
			Address: addr,
			TimestampNonceDetails: accountplustypes.TimestampNonceDetails{
				TimestampNonces: []uint64{1, 2, 3},
				MaxEjectedNonce: 0,
			},
		}
		s.Require().NotNil(bzNew, "Prefixed AccountState should exist for %s", addr)
		s.Require().Equal(expectedAccountState, actualAccountState, "Incorrect AccountState after migration for %s", addr)
	}
}
