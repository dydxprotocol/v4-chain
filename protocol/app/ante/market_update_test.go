package ante_test

import (
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	assets "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	perpetualtypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	prices_types "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/skip-mev/slinky/pkg/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/app/ante"
	testante "github.com/dydxprotocol/v4-chain/protocol/testutil/ante"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
)

func TestIsMarketUpdateTx(t *testing.T) {
	tests := []struct {
		name    string
		msgs    []sdk.Msg
		want    bool
		wantErr bool
	}{
		{
			name: "do nothing for a single bank message",
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
			name: "do nothing for multiple bank messages",
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
			name: "return true for single msg update",
			msgs: []sdk.Msg{
				&mmtypes.MsgUpdateMarkets{
					Authority:     constants.BobAccAddress.String(),
					UpdateMarkets: nil,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "return true for single msg upsert",
			msgs: []sdk.Msg{
				&mmtypes.MsgUpsertMarkets{
					Authority: constants.BobAccAddress.String(),
					Markets:   nil,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "return error for multiple msg update",
			msgs: []sdk.Msg{
				&mmtypes.MsgUpdateMarkets{
					Authority:     constants.BobAccAddress.String(),
					UpdateMarkets: nil,
				},
				&mmtypes.MsgUpdateMarkets{
					Authority:     constants.BobAccAddress.String(),
					UpdateMarkets: nil,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "return error for multiple msg upsert",
			msgs: []sdk.Msg{
				&mmtypes.MsgUpsertMarkets{
					Authority: constants.BobAccAddress.String(),
					Markets:   nil,
				},
				&mmtypes.MsgUpsertMarkets{
					Authority: constants.BobAccAddress.String(),
					Markets:   nil,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "return error for multiple mixed market map messages",
			msgs: []sdk.Msg{
				&mmtypes.MsgUpsertMarkets{
					Authority: constants.BobAccAddress.String(),
					Markets:   nil,
				},
				&mmtypes.MsgUpdateMarkets{
					Authority:     constants.BobAccAddress.String(),
					UpdateMarkets: nil,
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

			got, err := ante.IsMarketUpdateTx(tx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if got != tt.want {
				t.Errorf("IsMarketUpdateTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

var (
	testMarket = mmtypes.Market{
		Ticker: mmtypes.Ticker{
			CurrencyPair: types.CurrencyPair{
				Base:  "TEST",
				Quote: "USD",
			},
			Decimals:         1,
			MinProviderCount: 1,
			Enabled:          true,
			Metadata_JSON:    "",
		},
		ProviderConfigs: nil,
	}

	testMarketParams = prices_types.MarketParam{
		Id:                 0,
		Pair:               "TEST-USD",
		Exponent:           -8,
		MinExchanges:       1,
		MinPriceChangePpm:  10,
		ExchangeConfigJson: `{"test_config_placeholder":{}}`,
	}
)

func TestValidateMarketUpdateDecorator_AnteHandle(t *testing.T) {
	type marketPerpPair struct {
		market prices_types.MarketParam
		perp   perpetualtypes.Perpetual
	}

	type args struct {
		msgs        []sdk.Msg
		simulate    bool
		marketPerps []marketPerpPair
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "do nothing for non-market map messages - simulate",
			args: args{
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
				simulate: true,
			},
			wantErr: false,
		},
		{
			name: "do nothing for non-market map messages",
			args: args{
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
				simulate: false,
			},
			wantErr: false,
		},
		{
			name: "reject for multiple market messages - simulate",
			args: args{
				msgs: []sdk.Msg{
					&mmtypes.MsgUpsertMarkets{
						Authority: constants.BobAccAddress.String(),
						Markets:   nil,
					},
					&mmtypes.MsgUpsertMarkets{
						Authority: constants.BobAccAddress.String(),
						Markets:   nil,
					},
				},
				simulate: true,
			},
			wantErr: true,
		},
		{
			name: "reject for multiple market messages",
			args: args{
				msgs: []sdk.Msg{
					&mmtypes.MsgUpsertMarkets{
						Authority: constants.BobAccAddress.String(),
						Markets:   nil,
					},
					&mmtypes.MsgUpsertMarkets{
						Authority: constants.BobAccAddress.String(),
						Markets:   nil,
					},
				},
				simulate: false,
			},
			wantErr: true,
		},
		{
			name: "reject for multiple mixed market messages - simulate",
			args: args{
				msgs: []sdk.Msg{
					&mmtypes.MsgUpsertMarkets{
						Authority: constants.BobAccAddress.String(),
						Markets:   nil,
					},
					&mmtypes.MsgUpdateMarkets{
						Authority:     constants.BobAccAddress.String(),
						UpdateMarkets: nil,
					},
				},
				simulate: true,
			},
			wantErr: true,
		},
		{
			name: "reject for multiple mixed market messages",
			args: args{
				msgs: []sdk.Msg{
					&mmtypes.MsgUpsertMarkets{
						Authority: constants.BobAccAddress.String(),
						Markets:   nil,
					},
					&mmtypes.MsgUpdateMarkets{
						Authority:     constants.BobAccAddress.String(),
						UpdateMarkets: nil,
					},
				},
				simulate: false,
			},
			wantErr: true,
		},
		{
			name: "accept a single message with no markets",
			args: args{
				msgs: []sdk.Msg{
					&mmtypes.MsgUpsertMarkets{
						Authority: constants.BobAccAddress.String(),
						Markets:   nil,
					},
				},
				simulate: false,
			},
			wantErr: false,
		},
		{
			name: "accept a single message with no perps that match the market (disabled market)",
			args: args{
				msgs: []sdk.Msg{
					&mmtypes.MsgUpsertMarkets{
						Authority: constants.BobAccAddress.String(),
						Markets: []mmtypes.Market{
							testMarket,
						},
					},
				},
				simulate:    false,
				marketPerps: []marketPerpPair{},
			},
			wantErr: false,
		},
		{
			name: "accept a single message with no cross markets (only isolated)",
			args: args{
				msgs: []sdk.Msg{
					&mmtypes.MsgUpsertMarkets{
						Authority: constants.BobAccAddress.String(),
						Markets: []mmtypes.Market{
							testMarket,
						},
					},
				},
				simulate: false,
				marketPerps: []marketPerpPair{
					{
						market: testMarketParams,
						perp: perpetualtypes.Perpetual{
							Params: perpetualtypes.PerpetualParams{
								Id:                0,
								Ticker:            "BTC-USD small margin requirement",
								MarketId:          testMarketParams.Id,
								AtomicResolution:  int32(-8),
								DefaultFundingPpm: int32(0),
								LiquidityTier:     uint32(0),
								MarketType:        perpetualtypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED,
							},
							FundingIndex: dtypes.ZeroInt(),
							OpenInterest: dtypes.ZeroInt(),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "reject a single message with cross markets",
			args: args{
				msgs: []sdk.Msg{
					&mmtypes.MsgUpsertMarkets{
						Authority: constants.BobAccAddress.String(),
						Markets: []mmtypes.Market{
							testMarket,
						},
					},
				},
				simulate: false,
				marketPerps: []marketPerpPair{
					{
						market: testMarketParams,
						perp: perpetualtypes.Perpetual{
							Params: perpetualtypes.PerpetualParams{
								Id:                0,
								Ticker:            "TEST-USD small margin requirement",
								MarketId:          testMarketParams.Id,
								AtomicResolution:  int32(-8),
								DefaultFundingPpm: int32(0),
								LiquidityTier:     uint32(0),
								MarketType:        perpetualtypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
							},
							FundingIndex: dtypes.ZeroInt(),
							OpenInterest: dtypes.ZeroInt(),
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()

			// setup initial perps based on test
			for _, pair := range tt.args.marketPerps {
				marketID := rand.Uint32()
				pair.market.Id = marketID

				_, err := tApp.App.PricesKeeper.CreateMarket(
					ctx,
					pair.market,
					prices_types.MarketPrice{
						Id:       marketID,
						Exponent: -8,
						Price:    10,
					},
				)
				require.NoError(t, err)

				_, err = tApp.App.PerpetualsKeeper.CreatePerpetual(
					ctx,
					marketID,
					pair.perp.Params.Ticker,
					marketID,
					pair.perp.Params.AtomicResolution,
					pair.perp.Params.DefaultFundingPpm,
					pair.perp.Params.LiquidityTier,
					pair.perp.Params.MarketType,
				)
				require.NoError(t, err)
			}

			wrappedHandler := ante.NewValidateMarketUpdateDecorator(tApp.App.PerpetualsKeeper, tApp.App.PricesKeeper)
			anteHandler := sdk.ChainAnteDecorators(wrappedHandler)

			// Empty private key, so tx's signature should be empty.
			var (
				privs   []cryptotypes.PrivKey
				accSeqs []uint64
				accNums []uint64
			)

			tx, err := tx.CreateTestTx(
				ctx,
				tt.args.msgs,
				privs,
				accNums,
				accSeqs,
				tApp.App.ChainID(),
				signing.SignMode_SIGN_MODE_DIRECT,
				tApp.App.TxConfig(),
				0, // timeoutHeight
			)
			require.NoError(t, err)

			_, err = anteHandler(ctx, tx, tt.args.simulate)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
