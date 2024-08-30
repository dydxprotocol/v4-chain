package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	cometbfttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var (
	accAddrs = []sdk.AccAddress{
		sdk.AccAddress([]byte("cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5")),
		sdk.AccAddress([]byte("invalid_______________")),
	}

	price1, _ = ConvertStringToBigInt("1095368296575849877285046738")

	price2, _ = ConvertStringToBigInt("1095369098523294619828519483")

	price3, _ = ConvertStringToBigInt("1095369387224518455677735570")

	price_two = new(big.Int).Mul(big.NewInt(2), new(big.Int).Exp(big.NewInt(types.BASE_10), big.NewInt(types.SDAI_DECIMALS), nil))
)

type PoolTestTransfer struct {
	// Setup.
	sDAIAmount             *big.Int
	sDAIPrice              *big.Int
	userAddr               sdk.AccAddress
	userInitialSDAIBalance *big.Int
	userInitialTDAIBalance *big.Int
	// Expectations.
	expectedTDAIAmount *big.Int
	expectedErr        error
	expectErr          bool
}

type PoolTestCase struct {
	transfers []PoolTestTransfer
}

func TestGetTradingDAIFromSDAIAmount(t *testing.T) {

	tests := map[string]struct {
		sDAIAmount         *big.Int
		sDAIPrice          *big.Int
		expectedTDAIAmount *big.Int
		expectedErr        error
	}{
		"Zero sDAI amount": {
			sDAIAmount:         big.NewInt(0),
			sDAIPrice:          big.NewInt(1),
			expectedTDAIAmount: big.NewInt(0),
			expectedErr:        nil,
		},
		"Non-zero sDAI amount with valid price": {
			sDAIAmount:         big.NewInt(500),
			sDAIPrice:          price_two,
			expectedTDAIAmount: big.NewInt(1000),
			expectedErr:        nil,
		},
		"sDAI price not found": {
			sDAIAmount:         big.NewInt(1000),
			sDAIPrice:          nil,
			expectedTDAIAmount: nil,
			expectedErr:        errors.New("sDAI price not found"),
		},
		"Division by zero": {
			sDAIAmount:         big.NewInt(1000),
			sDAIPrice:          big.NewInt(0),
			expectedTDAIAmount: nil,
			expectedErr:        errors.New("sDAI price is zero"),
		},
		"Real example": {
			sDAIAmount:         big.NewInt(913),
			sDAIPrice:          price1,
			expectedTDAIAmount: big.NewInt(1000),
			expectedErr:        nil,
		},
		"Real example 2": {
			sDAIAmount:         big.NewInt(913),
			sDAIPrice:          price2,
			expectedTDAIAmount: big.NewInt(1000),
			expectedErr:        nil,
		},
		"Real example 3": {
			sDAIAmount:         big.NewInt(90166324963409613),
			sDAIPrice:          price3,
			expectedTDAIAmount: big.NewInt(98765432123456789),
			expectedErr:        nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis cometbfttypes.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				return genesis
			}).Build()

			ctx := tApp.InitChain()

			k := tApp.App.RatelimitKeeper

			if tc.sDAIPrice != nil {
				k.SetSDAIPrice(ctx, tc.sDAIPrice)
			}

			gotSDAIAmount, err := k.GetTradingDAIFromSDAIAmount(ctx, tc.sDAIAmount)

			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedTDAIAmount, gotSDAIAmount, "SDAI amounts mismatch.")
		})
	}
}

func TestGetTradingDAIFromSDAIAmountAndRoundUp(t *testing.T) {

	tests := map[string]struct {
		sDAIAmount         *big.Int
		sDAIPrice          *big.Int
		expectedTDAIAmount *big.Int
		expectedErr        error
	}{
		"Zero sDAI amount": {
			sDAIAmount:         big.NewInt(0),
			sDAIPrice:          big.NewInt(1),
			expectedTDAIAmount: big.NewInt(0),
			expectedErr:        nil,
		},
		"Non-zero sDAI amount with valid price": {
			sDAIAmount:         big.NewInt(500),
			sDAIPrice:          price_two,
			expectedTDAIAmount: big.NewInt(1000),
			expectedErr:        nil,
		},
		"sDAI price not found": {
			sDAIAmount:         big.NewInt(1000),
			sDAIPrice:          nil,
			expectedTDAIAmount: nil,
			expectedErr:        errors.New("sDai price not found: Failed to convert sDai amount to corresponding TDai Amount"),
		},
		"Division by zero": {
			sDAIAmount:         big.NewInt(1000),
			sDAIPrice:          big.NewInt(0),
			expectedTDAIAmount: nil,
			expectedErr:        errors.New("sDAI price is zero"),
		},
		"Real example": {
			sDAIAmount:         big.NewInt(913),
			sDAIPrice:          price1,
			expectedTDAIAmount: big.NewInt(1001),
			expectedErr:        nil,
		},
		"Real example 2": {
			sDAIAmount:         big.NewInt(913),
			sDAIPrice:          price2,
			expectedTDAIAmount: big.NewInt(1001),
			expectedErr:        nil,
		},
		"Real example 3": {
			sDAIAmount:         big.NewInt(90166324963409613),
			sDAIPrice:          price3,
			expectedTDAIAmount: big.NewInt(98765432123456790),
			expectedErr:        nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis cometbfttypes.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				return genesis
			}).Build()

			ctx := tApp.InitChain()

			k := tApp.App.RatelimitKeeper

			if tc.sDAIPrice != nil {
				k.SetSDAIPrice(ctx, tc.sDAIPrice)
			}

			gotSDAIAmount, err := k.GetTradingDAIFromSDAIAmountAndRoundUp(ctx, tc.sDAIAmount)

			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedTDAIAmount, gotSDAIAmount, "SDAI amounts mismatch.")
		})
	}
}

func TestMintTradingDAIToUserAccount(t *testing.T) {
	// Test Case Definition
	tests := map[string]PoolTestCase{
		"User has more sDAI than transfer amount": {
			transfers: []PoolTestTransfer{
				{
					sDAIAmount:             big.NewInt(250),
					sDAIPrice:              price_two,
					userAddr:               accAddrs[0],
					userInitialSDAIBalance: big.NewInt(1000),
					expectedTDAIAmount:     big.NewInt(500),
					expectedErr:            nil,
					expectErr:              false,
				},
			},
		},
		"User has exactly the sDAI transfer amount": {
			transfers: []PoolTestTransfer{
				{
					sDAIAmount:             big.NewInt(500),
					sDAIPrice:              price_two,
					userAddr:               accAddrs[0],
					userInitialSDAIBalance: big.NewInt(1000),
					expectedTDAIAmount:     big.NewInt(1000),
					expectedErr:            nil,
					expectErr:              false,
				},
			},
		},
		"User has less sDAI than transfer amount": {
			transfers: []PoolTestTransfer{
				{
					sDAIAmount:             big.NewInt(1000),
					sDAIPrice:              price_two,
					userAddr:               accAddrs[0],
					userInitialSDAIBalance: big.NewInt(500),
					expectedTDAIAmount:     nil,
					expectedErr:            errors.New("failed to send sDAI to ratelimit module"),
					expectErr:              true,
				},
			},
		},
		"User has zero sDAI balance": {
			transfers: []PoolTestTransfer{
				{
					sDAIAmount:             big.NewInt(1000),
					sDAIPrice:              price_two,
					userAddr:               accAddrs[0],
					userInitialSDAIBalance: big.NewInt(0),
					expectedTDAIAmount:     nil,
					expectedErr:            errors.New("failed to send sDAI to ratelimit module"),
					expectErr:              true,
				},
			},
		},
		"User has large sDAI balance and small transfer amount": {
			transfers: []PoolTestTransfer{
				{
					sDAIAmount:             big.NewInt(1),
					sDAIPrice:              price_two,
					userAddr:               accAddrs[0],
					userInitialSDAIBalance: big.NewInt(1000000),
					expectedTDAIAmount:     big.NewInt(2),
					expectedErr:            nil,
					expectErr:              false,
				},
			},
		},
		"User has small sDAI balance and large transfer amount": {
			transfers: []PoolTestTransfer{
				{
					sDAIAmount:             big.NewInt(1000000),
					sDAIPrice:              price_two,
					userAddr:               accAddrs[0],
					userInitialSDAIBalance: big.NewInt(1),
					expectedTDAIAmount:     nil,
					expectedErr:            errors.New("failed to send sDAI to ratelimit module"),
					expectErr:              true,
				},
			},
		},
		"Real price will round down": {
			transfers: []PoolTestTransfer{
				{
					sDAIAmount:             big.NewInt(913),
					sDAIPrice:              price2,
					userAddr:               accAddrs[0],
					userInitialSDAIBalance: big.NewInt(2000),
					expectedTDAIAmount:     big.NewInt(1000),
					expectErr:              false,
				},
			},
		},
		"User has an invalid address": {
			transfers: []PoolTestTransfer{
				{
					sDAIAmount:             big.NewInt(1000000),
					sDAIPrice:              price_two,
					userAddr:               accAddrs[1],
					userInitialSDAIBalance: big.NewInt(1),
					expectedTDAIAmount:     nil,
					expectErr:              true,
				},
			},
		},
	}

	// Test Case Execution
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis cometbfttypes.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				return genesis
			}).Build()

			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper
			bankKeeper := tApp.App.BankKeeper
			accountKeeper := tApp.App.AccountKeeper

			// Process each minting transfer
			for _, transfer := range tc.transfers {

				/* Setup state */
				k.SetSDAIPrice(ctx, transfer.sDAIPrice)
				sDAICoins := sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, sdkmath.NewIntFromBigInt(transfer.userInitialSDAIBalance)))

				mintingErr := bankKeeper.MintCoins(ctx, types.TDaiPoolAccount, sDAICoins)
				require.NoError(t, mintingErr)
				sendingErr := bankKeeper.SendCoinsFromModuleToAccount(ctx, types.TDaiPoolAccount, transfer.userAddr, sDAICoins)
				require.NoError(t, sendingErr)

				// Check initial balance
				initialPoolSDAIBalance := bankKeeper.GetBalance(
					ctx,
					accountKeeper.GetModuleAddress(types.SDaiPoolAccount),
					types.SDaiDenom,
				).Amount.BigInt()
				initialPoolTDAIBalance := bankKeeper.GetBalance(
					ctx,
					accountKeeper.GetModuleAddress(types.TDaiPoolAccount),
					types.TDaiDenom,
				).Amount.BigInt()
				initialUserSDAIBalance := bankKeeper.GetBalance(
					ctx,
					transfer.userAddr,
					types.SDaiDenom,
				).Amount.BigInt()
				initialUserTDAIBalance := bankKeeper.GetBalance(
					ctx,
					transfer.userAddr,
					types.TDaiDenom,
				).Amount.BigInt()

				// Execute Minting
				err := k.MintTradingDAIToUserAccount(ctx, transfer.userAddr, transfer.sDAIAmount)

				// Verify success
				if transfer.expectErr {
					require.Error(t, err)
					if transfer.expectedErr != nil {
						require.ErrorContains(t, err, transfer.expectedErr.Error())
					}
				} else {
					require.NoError(t, err)

					// Verify state change
					endingPoolSDAIBalance := bankKeeper.GetBalance(
						ctx,
						accountKeeper.GetModuleAddress(types.SDaiPoolAccount),
						types.SDaiDenom,
					).Amount.BigInt()
					endingPoolTDAIBalance := bankKeeper.GetBalance(
						ctx,
						accountKeeper.GetModuleAddress(types.TDaiPoolAccount),
						types.TDaiDenom,
					).Amount.BigInt()
					endingUserSDAIBalance := bankKeeper.GetBalance(
						ctx,
						transfer.userAddr,
						types.SDaiDenom,
					).Amount.BigInt()
					endingUserTDAIBalance := bankKeeper.GetBalance(
						ctx,
						transfer.userAddr,
						types.TDaiDenom,
					).Amount.BigInt()

					deltaPoolSDAI := new(big.Int).Sub(endingPoolSDAIBalance, initialPoolSDAIBalance)
					require.Equal(t, transfer.sDAIAmount, deltaPoolSDAI, "Change in pool SDAI balance incorrect.")

					deltaPoolTDAI := new(big.Int).Sub(endingPoolTDAIBalance, initialPoolTDAIBalance)
					require.Equal(t, big.NewInt(0), deltaPoolTDAI, "Change in pool TDAI balance incorrect. Should always be 0 when minting.")

					deltaUserSDAI := new(big.Int).Sub(initialUserSDAIBalance, endingUserSDAIBalance)
					require.Equal(t, transfer.sDAIAmount, deltaUserSDAI, "Change in user SDAI balance incorrect.")

					deltaUserTDAI := new(big.Int).Sub(endingUserTDAIBalance, initialUserTDAIBalance)
					require.Equal(t, transfer.expectedTDAIAmount, deltaUserTDAI, "Change in user TDAI balance incorrect.")
				}
			}
		})
	}
}

func TestWithdrawSDaiFromTDai(t *testing.T) {
	// Test Case Definition
	tests := map[string]PoolTestCase{
		"User has more tDAI than transfer amount": {
			transfers: []PoolTestTransfer{
				{
					sDAIAmount:             big.NewInt(250),
					sDAIPrice:              price_two,
					userAddr:               accAddrs[0],
					userInitialTDAIBalance: big.NewInt(1000),
					expectedTDAIAmount:     big.NewInt(500),
					expectedErr:            nil,
					expectErr:              false,
				},
			},
		},
		"User has exactly the tDAI transfer amount": {
			transfers: []PoolTestTransfer{
				{
					sDAIAmount:             big.NewInt(500),
					sDAIPrice:              price_two,
					userAddr:               accAddrs[0],
					userInitialTDAIBalance: big.NewInt(1000),
					expectedTDAIAmount:     big.NewInt(1000),
					expectedErr:            nil,
					expectErr:              false,
				},
			},
		},
		"User has less tDAI than transfer amount": {
			transfers: []PoolTestTransfer{
				{
					sDAIAmount:             big.NewInt(1000),
					sDAIPrice:              price_two,
					userAddr:               accAddrs[0],
					userInitialTDAIBalance: big.NewInt(250),
					expectedTDAIAmount:     nil,
					expectedErr:            errors.New("failed to send tDAI from user account to tDai pool account: spendable balance 250utdai is smaller than 2000utdai: insufficient funds"),
					expectErr:              true,
				},
			},
		},
		"User has zero tDAI balance": {
			transfers: []PoolTestTransfer{
				{
					sDAIAmount:             big.NewInt(1000),
					sDAIPrice:              price_two,
					userAddr:               accAddrs[0],
					userInitialTDAIBalance: big.NewInt(0),
					expectedTDAIAmount:     nil,
					expectedErr:            errors.New("failed to send tDAI from user account to tDai pool account: spendable balance 0utdai is smaller than 2000utdai: insufficient funds"),
					expectErr:              true,
				},
			},
		},
		"User has large tDAI balance and small transfer amount": {
			transfers: []PoolTestTransfer{
				{
					sDAIAmount:             big.NewInt(1),
					sDAIPrice:              price_two,
					userAddr:               accAddrs[0],
					userInitialTDAIBalance: big.NewInt(1000000),
					expectedTDAIAmount:     big.NewInt(2),
					expectedErr:            nil,
					expectErr:              false,
				},
			},
		},
		"User has small tDAI balance and large transfer amount": {
			transfers: []PoolTestTransfer{
				{
					sDAIAmount:             big.NewInt(1000000),
					sDAIPrice:              price_two,
					userAddr:               accAddrs[0],
					userInitialTDAIBalance: big.NewInt(1),
					expectedTDAIAmount:     nil,
					expectedErr:            errors.New("failed to send tDAI from user account to tDai pool account: spendable balance 1utdai is smaller than 2000000utdai: insufficient funds"),
					expectErr:              true,
				},
			},
		},
		"Real price will round up": {
			transfers: []PoolTestTransfer{
				{
					sDAIAmount:             big.NewInt(913),
					sDAIPrice:              price2,
					userAddr:               accAddrs[0],
					userInitialTDAIBalance: big.NewInt(2000),
					expectedTDAIAmount:     big.NewInt(1001),
					expectErr:              false,
				},
			},
		},
		"User has an invalid address": {
			transfers: []PoolTestTransfer{
				{
					sDAIAmount:             big.NewInt(1000000),
					sDAIPrice:              price_two,
					userAddr:               accAddrs[1],
					userInitialTDAIBalance: big.NewInt(1),
					expectedTDAIAmount:     nil,
					expectErr:              true,
				},
			},
		},
	}

	// Test Case Execution
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis cometbfttypes.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				return genesis
			}).Build()

			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper
			bankKeeper := tApp.App.BankKeeper
			accountKeeper := tApp.App.AccountKeeper

			// Process each minting transfer
			for _, transfer := range tc.transfers {

				/* Setup state */
				k.SetSDAIPrice(ctx, transfer.sDAIPrice)

				sDAICoins := sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, sdkmath.NewIntFromBigInt(transfer.sDAIAmount)))
				mintingErr := bankKeeper.MintCoins(ctx, types.TDaiPoolAccount, sDAICoins)
				require.NoError(t, mintingErr)
				sendingErr := bankKeeper.SendCoinsFromModuleToModule(ctx, types.TDaiPoolAccount, types.SDaiPoolAccount, sDAICoins)
				require.NoError(t, sendingErr)

				tDAICoins := sdk.NewCoins(sdk.NewCoin(types.TDaiDenom, sdkmath.NewIntFromBigInt(transfer.userInitialTDAIBalance)))

				// Simulate user having appropriate amount of tDAI in their account
				// TODO: Make sure that we also test cases, where the user does not have enought tDAI to mint the given amount of sDAI
				mintingErr = bankKeeper.MintCoins(ctx, types.TDaiPoolAccount, tDAICoins)
				require.NoError(t, mintingErr)
				sendingErr = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.TDaiPoolAccount, transfer.userAddr, tDAICoins)
				require.NoError(t, sendingErr)

				// Check initial balance
				initialPoolSDAIBalance := bankKeeper.GetBalance(
					ctx,
					accountKeeper.GetModuleAddress(types.SDaiPoolAccount),
					types.SDaiDenom,
				).Amount.BigInt()
				initialPoolTDAIBalance := bankKeeper.GetBalance(
					ctx,
					accountKeeper.GetModuleAddress(types.TDaiPoolAccount),
					types.TDaiDenom,
				).Amount.BigInt()
				initialUserSDAIBalance := bankKeeper.GetBalance(
					ctx,
					transfer.userAddr,
					types.SDaiDenom,
				).Amount.BigInt()
				initialUserTDAIBalance := bankKeeper.GetBalance(
					ctx,
					transfer.userAddr,
					types.TDaiDenom,
				).Amount.BigInt()

				// Execute Minting
				err := k.WithdrawSDaiFromTDai(ctx, transfer.userAddr, transfer.sDAIAmount)

				// Verify success
				if transfer.expectErr {
					require.Error(t, err)
					if transfer.expectedErr != nil {
						require.ErrorContains(t, err, transfer.expectedErr.Error())
					}
				} else {
					require.NoError(t, err)

					// Verify state change
					endingPoolSDAIBalance := bankKeeper.GetBalance(
						ctx,
						accountKeeper.GetModuleAddress(types.SDaiPoolAccount),
						types.SDaiDenom,
					).Amount.BigInt()
					endingPoolTDAIBalance := bankKeeper.GetBalance(
						ctx,
						accountKeeper.GetModuleAddress(types.TDaiPoolAccount),
						types.TDaiDenom,
					).Amount.BigInt()
					endingUserSDAIBalance := bankKeeper.GetBalance(
						ctx,
						transfer.userAddr,
						types.SDaiDenom,
					).Amount.BigInt()
					endingUserTDAIBalance := bankKeeper.GetBalance(
						ctx,
						transfer.userAddr,
						types.TDaiDenom,
					).Amount.BigInt()

					deltaPoolSDAI := new(big.Int).Sub(initialPoolSDAIBalance, endingPoolSDAIBalance)
					require.Equal(t, transfer.sDAIAmount, deltaPoolSDAI, "Change in pool SDAI balance incorrect.")

					deltaPoolTDAI := new(big.Int).Sub(endingPoolTDAIBalance, initialPoolTDAIBalance)
					require.Equal(t, big.NewInt(0), deltaPoolTDAI, "Change in pool TDAI balance incorrect. Should always be 0 when minting.")

					deltaUserSDAI := new(big.Int).Sub(endingUserSDAIBalance, initialUserSDAIBalance)
					require.Equal(t, transfer.sDAIAmount, deltaUserSDAI, "Change in user SDAI balance incorrect.")

					deltaUserTDAI := new(big.Int).Sub(initialUserTDAIBalance, endingUserTDAIBalance)
					require.Equal(t, transfer.expectedTDAIAmount, deltaUserTDAI, "Change in user TDAI balance incorrect.")
				}
			}
		})
	}
}

func ConvertStringToBigInt(str string) (*big.Int, error) {

	bigint, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return nil, errorsmod.Wrap(
			types.ErrUnableToDecodeBigInt,
			"Unable to convert the sDAI conversion rate to a big int",
		)
	}

	return bigint, nil
}
