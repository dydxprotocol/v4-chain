package keeper_test

import (
	"math"
	"testing"

	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

func TestInitializeAccount(t *testing.T) {
	baseTsNonce := uint64(math.Pow(2, 40))
	genesisState := &types.GenesisState{
		Accounts: []types.AccountState{
			{
				Address: constants.AliceAccAddress.String(),
				TimestampNonceDetails: types.TimestampNonceDetails{
					TimestampNonces: []uint64{baseTsNonce + 1, baseTsNonce + 2, baseTsNonce + 3},
					MaxEjectedNonce: baseTsNonce,
				},
			},
			{
				Address: constants.BobAccAddress.String(),
				TimestampNonceDetails: types.TimestampNonceDetails{
					TimestampNonces: []uint64{baseTsNonce + 5, baseTsNonce + 6, baseTsNonce + 7},
					MaxEjectedNonce: baseTsNonce + 1,
				},
			},
		},
	}

	t.Run("Cannot initialize existing account", func(t *testing.T) {
		ctx, k, _, _ := keepertest.TimestampNonceKeepers(t)
		accountplus.InitGenesis(ctx, *k, *genesisState)
		err := k.InitializeAccount(ctx, constants.AliceAccAddress)
		require.NotNil(t, err, "Account should not be able to be initialized if already exists")
	})

	t.Run("Can initialize new account", func(t *testing.T) {
		ctx, k, _, _ := keepertest.TimestampNonceKeepers(t)
		accountplus.InitGenesis(ctx, *k, *genesisState)

		expectedAccount := keeper.DefaultAccountState(constants.CarlAccAddress)

		err := k.InitializeAccount(ctx, constants.CarlAccAddress)
		require.Nil(t, err, "Should be able to initialize account if it did not exist")

		actualAccount, found := k.GetAccountState(ctx, constants.CarlAccAddress)
		require.True(t, found, "Could not find account")
		require.Equal(t, actualAccount, expectedAccount)
	})
}
