package ante_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	assets "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/rate_limit"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/stretchr/testify/require"

	customante "github.com/dydxprotocol/v4-chain/protocol/app/ante"
	testante "github.com/dydxprotocol/v4-chain/protocol/testutil/ante"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func TestIsSingleWasmExecTx(t *testing.T) {
	tests := []struct {
		name    string
		msgs    []sdk.Msg
		want    bool
		wantErr bool
	}{
		{
			name: "returns false for a single bank message",
			msgs: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: constants.BobAccAddress.String(),
					ToAddress:   constants.AliceAccAddress.String(),
					Amount: []sdk.Coin{
						sdk.NewCoin(assets.AssetUsdc.Denom, sdkmath.NewInt(1)),
					},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false for multiple bank messages",
			msgs: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: constants.BobAccAddress.String(),
					ToAddress:   constants.AliceAccAddress.String(),
					Amount: []sdk.Coin{
						sdk.NewCoin(assets.AssetUsdc.Denom, sdkmath.NewInt(1)),
					},
				},
				&banktypes.MsgSend{
					FromAddress: constants.BobAccAddress.String(),
					ToAddress:   constants.AliceAccAddress.String(),
					Amount: []sdk.Coin{
						sdk.NewCoin(assets.AssetUsdc.Denom, sdkmath.NewInt(1)),
					},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns true for single wasm exec message",
			msgs: []sdk.Msg{
				&wasmtypes.MsgExecuteContract{
					Sender: constants.BobAccAddress.String(),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "return false for single wasm instantiate",
			msgs: []sdk.Msg{
				&wasmtypes.MsgInstantiateContract{
					Sender: constants.BobAccAddress.String(),
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "return error for multiple wasm execute msg",
			msgs: []sdk.Msg{
				&wasmtypes.MsgExecuteContract{
					Sender: constants.BobAccAddress.String(),
				},
				&wasmtypes.MsgExecuteContract{
					Sender: constants.BobAccAddress.String(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "return error for wasm execute msg + wasm instantiate msg",
			msgs: []sdk.Msg{
				&wasmtypes.MsgExecuteContract{
					Sender: constants.BobAccAddress.String(),
				},
				&wasmtypes.MsgInstantiateContract{
					Sender: constants.BobAccAddress.String(),
				},
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := testante.SetupTestSuite(t, true)
			suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder()

			require.NoError(t, suite.TxBuilder.SetMsgs(tt.msgs...))

			// Empty private key, so tx's signature should be empty.
			var (
				privs   []cryptotypes.PrivKey
				accSeqs []uint64
				accNums []uint64
			)

			tx, err := suite.CreateTestTx(
				suite.Ctx,
				privs,
				accNums,
				accSeqs,
				suite.Ctx.ChainID(),
				signing.SignMode_SIGN_MODE_DIRECT,
			)
			require.NoError(t, err)

			got, err := customante.IsSingleWasmExecTx(tx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if got != tt.want {
				t.Errorf("IsSingleWasmExecTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWasmExecDecorator(t *testing.T) {
	suite := testante.SetupTestSuite(t, true)

	var err error
	suite.ClientCtx.TxConfig, err = authtx.NewTxConfigWithOptions(
		codec.NewProtoCodec(suite.EncCfg.InterfaceRegistry),
		authtx.ConfigOptions{},
	)
	require.NoError(t, err)
	suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder()

	// make block height non-zero to ensure account numbers part of signBytes
	suite.Ctx = suite.Ctx.WithBlockHeight(1)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	feeAmount := testdata.NewTestFeeAmount()

	type testCase struct {
		name           string
		msgs           []sdk.Msg
		gasLimit       uint64
		shouldErr      bool
		recheck        bool
		expectedErrMsg string // supply empty string to ignore this check
		rateLimiter    rate_limit.RateLimiter[string]
	}

	validGasLimit := uint64(1_000_000)
	invalidGasLimit := uint64(customante.WasmExecMaxGasLimit) + 1

	testCases := []testCase{
		{
			name: "Valid tx, no-op rate limiter",
			msgs: []sdk.Msg{
				&wasmtypes.MsgExecuteContract{
					Sender: addr1.String(),
				},
			},
			gasLimit:    validGasLimit,
			rateLimiter: rate_limit.NewNoOpRateLimiter[string](),
		},
		{
			name: "Error when rate limit is exceeded",
			msgs: []sdk.Msg{
				&wasmtypes.MsgExecuteContract{
					Sender: addr1.String(),
				},
			},
			gasLimit: validGasLimit,
			rateLimiter: rate_limit.NewSingleBlockRateLimiter[string](
				"test",
				clobtypes.MaxPerNBlocksRateLimit{
					Limit:     0,
					NumBlocks: 1,
				},
			),
			shouldErr:      true,
			expectedErrMsg: "Rate of 1 exceeds configured block rate limit of {NumBlocks:1 Limit:0}",
		},
		{
			name: "Don't ratelimit in ReCheck",
			msgs: []sdk.Msg{
				&wasmtypes.MsgExecuteContract{
					Sender: addr1.String(),
				},
			},
			gasLimit: validGasLimit,
			rateLimiter: rate_limit.NewSingleBlockRateLimiter[string](
				"test",
				clobtypes.MaxPerNBlocksRateLimit{
					Limit:     0,
					NumBlocks: 1,
				},
			),
			recheck:   true,
			shouldErr: false,
		},
		{
			name: "Valid tx, no-op rate limiter",
			msgs: []sdk.Msg{
				&wasmtypes.MsgExecuteContract{
					Sender: addr1.String(),
				},
			},
			gasLimit:    invalidGasLimit,
			rateLimiter: rate_limit.NewNoOpRateLimiter[string](),
			shouldErr:   true,
			expectedErrMsg: fmt.Sprintf(
				"CosmWasm execution specified gas limit (%v) exceeds `WasmExecMaxGasLimit` "+
					"(%v): invalid gas limit",
				invalidGasLimit,
				customante.WasmExecMaxGasLimit,
			),
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			suite.Ctx = suite.Ctx.WithIsReCheckTx(tc.recheck)
			suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder() // Create new txBuilder for each test

			require.NoError(t, suite.TxBuilder.SetMsgs(tc.msgs...))
			suite.TxBuilder.SetFeeAmount(feeAmount)
			suite.TxBuilder.SetGasLimit(tc.gasLimit)

			suite.RateLimitKeeper.SetRateLimiter(tc.rateLimiter)
			wasmExecDecorator := customante.NewWasmExecDecorator(
				suite.RateLimitKeeper,
			)

			tx, err := suite.CreateTestTx(
				suite.Ctx,
				[]cryptotypes.PrivKey{priv1}, // testPrivKey
				[]uint64{0},                  // testAccountNum
				[]uint64{0},                  // testAccountSeq
				suite.Ctx.ChainID(),
				signing.SignMode_SIGN_MODE_DIRECT,
			)
			require.NoError(t, err)

			antehandler := sdk.ChainAnteDecorators(wasmExecDecorator)

			txBytes, err := suite.ClientCtx.TxConfig.TxEncoder()(tx)
			require.NoError(t, err)
			byteCtx := suite.Ctx.WithTxBytes(txBytes)

			_, err = antehandler(byteCtx, tx, false)
			if tc.shouldErr {
				require.NotNil(t, err, "TestCase %d: %s did not error as expected", i, tc.name)
				if tc.expectedErrMsg != "" {
					require.Contains(t, err.Error(), tc.expectedErrMsg)
				}
			} else {
				require.Nil(t, err, "TestCase %d: %s errored unexpectedly. Err: %v", i, tc.name, err)
			}
		})
	}
}
