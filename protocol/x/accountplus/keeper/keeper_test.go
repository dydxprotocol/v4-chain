package keeper_test

import (
	"math"
	"testing"

	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
<<<<<<< HEAD
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/testutils"
=======
>>>>>>> main
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

func TestInitializeAccount(t *testing.T) {
	baseTsNonce := uint64(math.Pow(2, 40))
	genesisState := &types.GenesisState{
<<<<<<< HEAD
		Accounts: []*types.AccountState{
			{
				Address: constants.AliceAccAddress.String(),
				TimestampNonceDetails: &types.TimestampNonceDetails{
=======
		Accounts: []types.AccountState{
			{
				Address: constants.AliceAccAddress.String(),
				TimestampNonceDetails: types.TimestampNonceDetails{
>>>>>>> main
					TimestampNonces: []uint64{baseTsNonce + 1, baseTsNonce + 2, baseTsNonce + 3},
					MaxEjectedNonce: baseTsNonce,
				},
			},
			{
				Address: constants.BobAccAddress.String(),
<<<<<<< HEAD
				TimestampNonceDetails: &types.TimestampNonceDetails{
=======
				TimestampNonceDetails: types.TimestampNonceDetails{
>>>>>>> main
					TimestampNonces: []uint64{baseTsNonce + 5, baseTsNonce + 6, baseTsNonce + 7},
					MaxEjectedNonce: baseTsNonce + 1,
				},
			},
		},
	}

	t.Run("Cannot initialize existing account", func(t *testing.T) {
		ctx, k, _, _ := keepertest.TimestampNonceKeepers(t)
		accountplus.InitGenesis(ctx, *k, *genesisState)
<<<<<<< HEAD
		_, err := k.InitializeAccount(ctx, constants.AliceAccAddress)
=======
		err := k.InitializeAccount(ctx, constants.AliceAccAddress)
>>>>>>> main
		require.NotNil(t, err, "Account should not be able to be initialized if already exists")
	})

	t.Run("Can initialize new account", func(t *testing.T) {
		ctx, k, _, _ := keepertest.TimestampNonceKeepers(t)
		accountplus.InitGenesis(ctx, *k, *genesisState)

<<<<<<< HEAD
		expectedAccount := types.AccountState{
			Address:               constants.CarlAccAddress.String(),
			TimestampNonceDetails: keeper.DeepCopyTimestampNonceDetails(keeper.InitialTimestampNonceDetails),
		}

		account, err := k.InitializeAccount(ctx, constants.CarlAccAddress)
		require.Nil(t, err, "Should be able to initialize account if it did not exist")

		isAccountEqual := testutils.CompareAccountStates(&account, &expectedAccount)
		require.True(t, isAccountEqual, "Initialized account does not have correct initial state")
=======
		expectedAccount := keeper.DefaultAccountState(constants.CarlAccAddress)

		err := k.InitializeAccount(ctx, constants.CarlAccAddress)
		require.Nil(t, err, "Should be able to initialize account if it did not exist")

		actualAccount, found := k.GetAccountState(ctx, constants.CarlAccAddress)
		require.True(t, found, "Could not find account")
		require.Equal(t, actualAccount, expectedAccount)
>>>>>>> main
	})
}
