package clob_test

import (
	"encoding/json"
	"testing"

	sdkmath "cosmossdk.io/math"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	aptypes "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	"github.com/stretchr/testify/require"
)

type TestBlockWithMsgs struct {
	Block uint32
	Msgs  []TestSdkMsg
}

type TestSdkMsg struct {
	Msg            []sdk.Msg
	Authenticators []uint64
	Fees           sdk.Coins
	Gas            uint64
	AccountNum     []uint64
	SeqNum         []uint64
	Signers        []cryptotypes.PrivKey

	ExpectedRespCode uint32
	ExpectedLog      string
}

func TestPlaceOrder_PermissionedKeys_Failures(t *testing.T) {
	config := []aptypes.SubAuthenticatorInitData{
		{
			Type:   "SignatureVerification",
			Config: constants.AlicePrivateKey.PubKey().Bytes(),
		},
		{
			Type:   "MessageFilter",
			Config: []byte("/cosmos.bank.v1beta1.MsgSend"),
		},
	}
	compositeAuthenticatorConfig, err := json.Marshal(config)
	require.NoError(t, err)

	tests := map[string]struct {
		smartAccountEnabled bool
		blocks              []TestBlockWithMsgs

		expectedOrderIdsInMemclob map[clobtypes.OrderId]bool
	}{
		"Txn has authenticators specified, but smart account is not enabled": {
			smartAccountEnabled: false,
			blocks: []TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								clobtypes.NewMsgPlaceOrder(
									testapp.MustScaleOrder(
										constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
										testapp.DefaultGenesis(),
									),
								),
							},
							Authenticators: []uint64{0},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        0,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{0},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: aptypes.ErrSmartAccountNotActive.ABCICode(),
							ExpectedLog:      aptypes.ErrSmartAccountNotActive.Error(),
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId: false,
			},
		},
		"Txn has authenticators specified, but authenticator is not found": {
			smartAccountEnabled: true,
			blocks: []TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								clobtypes.NewMsgPlaceOrder(
									testapp.MustScaleOrder(
										constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
										testapp.DefaultGenesis(),
									),
								),
							},
							Authenticators: []uint64{0},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        0,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{0},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: aptypes.ErrAuthenticatorNotFound.ABCICode(),
							ExpectedLog:      aptypes.ErrAuthenticatorNotFound.Error(),
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId: false,
			},
		},
		"Txn has authenticators specified, but authenticator was removed": {
			smartAccountEnabled: true,
			blocks: []TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								&aptypes.MsgAddAuthenticator{
									Sender:            constants.BobAccAddress.String(),
									AuthenticatorType: "MessageFilter",
									Data:              []byte("/cosmos.bank.v1beta1.MsgSend"),
								},
							},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        300_000,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{1},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: 0,
						},
					},
				},
				{
					Block: 4,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								&aptypes.MsgRemoveAuthenticator{
									Sender: constants.BobAccAddress.String(),
									Id:     0,
								},
							},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        300_000,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{2},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: 0,
						},
					},
				},
				{
					Block: 6,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								clobtypes.NewMsgPlaceOrder(
									testapp.MustScaleOrder(
										constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
										testapp.DefaultGenesis(),
									),
								),
							},
							Authenticators: []uint64{0},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        0,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{0},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: aptypes.ErrAuthenticatorNotFound.ABCICode(),
							ExpectedLog:      aptypes.ErrAuthenticatorNotFound.Error(),
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId: false,
			},
		},
		"Txn rejected by signature verification authenticator": {
			smartAccountEnabled: true,
			blocks: []TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								&aptypes.MsgAddAuthenticator{
									Sender:            constants.BobAccAddress.String(),
									AuthenticatorType: "SignatureVerification",
									// Allow signature verification using Alice's public key.
									Data: constants.AlicePrivateKey.PubKey().Bytes(),
								},
							},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        300_000,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{1},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: 0,
						},
					},
				},
				{
					Block: 4,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								clobtypes.NewMsgPlaceOrder(
									testapp.MustScaleOrder(
										constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
										testapp.DefaultGenesis(),
									),
								),
							},
							Authenticators: []uint64{0},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        0,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{1},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: aptypes.ErrSignatureVerification.ABCICode(),
							ExpectedLog:      aptypes.ErrSignatureVerification.Error(),
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId: false,
			},
		},
		"Txn rejected by message filter authenticator": {
			smartAccountEnabled: true,
			blocks: []TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								&aptypes.MsgAddAuthenticator{
									Sender:            constants.BobAccAddress.String(),
									AuthenticatorType: "MessageFilter",
									Data:              []byte("/cosmos.bank.v1beta1.MsgSend"),
								},
							},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        300_000,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{1},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: 0,
						},
					},
				},
				{
					Block: 4,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								clobtypes.NewMsgPlaceOrder(
									testapp.MustScaleOrder(
										constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
										testapp.DefaultGenesis(),
									),
								),
							},
							Authenticators: []uint64{0},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        0,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{1},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: aptypes.ErrMessageTypeVerification.ABCICode(),
							ExpectedLog:      aptypes.ErrMessageTypeVerification.Error(),
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId: false,
			},
		},
		"Txn rejected by clob pair id filter authenticator": {
			smartAccountEnabled: true,
			blocks: []TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								&aptypes.MsgAddAuthenticator{
									Sender:            constants.BobAccAddress.String(),
									AuthenticatorType: "ClobPairIdFilter",
									Data:              []byte("0"),
								},
							},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        300_000,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{1},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: 0,
						},
					},
				},
				{
					Block: 4,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								clobtypes.NewMsgPlaceOrder(
									testapp.MustScaleOrder(
										constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
										testapp.DefaultGenesis(),
									),
								),
							},
							Authenticators: []uint64{0},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        0,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{1},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: aptypes.ErrClobPairIdVerification.ABCICode(),
							ExpectedLog:      aptypes.ErrClobPairIdVerification.Error(),
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId: false,
			},
		},
		"Txn rejected by subaccount number filter authenticator": {
			smartAccountEnabled: true,
			blocks: []TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								&aptypes.MsgAddAuthenticator{
									Sender:            constants.BobAccAddress.String(),
									AuthenticatorType: "SubaccountFilter",
									Data:              []byte("1"),
								},
							},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        300_000,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{1},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: 0,
						},
					},
				},
				{
					Block: 4,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								clobtypes.NewMsgPlaceOrder(
									testapp.MustScaleOrder(
										constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
										testapp.DefaultGenesis(),
									),
								),
							},
							Authenticators: []uint64{0},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        0,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{1},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: aptypes.ErrSubaccountVerification.ABCICode(),
							ExpectedLog:      aptypes.ErrSubaccountVerification.Error(),
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId: false,
			},
		},
		"Txn rejected by all of authenticator": {
			smartAccountEnabled: true,
			blocks: []TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								&aptypes.MsgAddAuthenticator{
									Sender:            constants.BobAccAddress.String(),
									AuthenticatorType: "AllOf",
									Data:              compositeAuthenticatorConfig,
								},
							},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        300_000,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{1},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: 0,
						},
					},
				},
				{
					Block: 4,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								clobtypes.NewMsgPlaceOrder(
									testapp.MustScaleOrder(
										constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
										testapp.DefaultGenesis(),
									),
								),
							},
							Authenticators: []uint64{0},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        0,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{1},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: aptypes.ErrAllOfVerification.ABCICode(),
							ExpectedLog:      aptypes.ErrAllOfVerification.Error(),
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId: false,
			},
		},
		"Txn rejected by any of authenticator": {
			smartAccountEnabled: true,
			blocks: []TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								&aptypes.MsgAddAuthenticator{
									Sender:            constants.BobAccAddress.String(),
									AuthenticatorType: "AnyOf",
									Data:              compositeAuthenticatorConfig,
								},
							},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        300_000,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{1},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: 0,
						},
					},
				},
				{
					Block: 4,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								clobtypes.NewMsgPlaceOrder(
									testapp.MustScaleOrder(
										constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
										testapp.DefaultGenesis(),
									),
								),
							},
							Authenticators: []uint64{0},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        0,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{1},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: aptypes.ErrAnyOfVerification.ABCICode(),
							ExpectedLog:      aptypes.ErrAnyOfVerification.Error(),
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId: false,
			},
		},
		"One of the messages in the transaction is rejected": {
			smartAccountEnabled: true,
			blocks: []TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								&aptypes.MsgAddAuthenticator{
									Sender:            constants.BobAccAddress.String(),
									AuthenticatorType: "MessageFilter",
									Data:              []byte("/cosmos.bank.v1beta1.MsgSend"),
								},
							},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        300_000,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{1},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: 0,
						},
					},
				},
				{
					Block: 4,
					Msgs: []TestSdkMsg{
						{
							Msg: []sdk.Msg{
								&banktypes.MsgSend{
									FromAddress: constants.BobAccAddress.String(),
									ToAddress:   constants.AliceAccAddress.String(),
									Amount: sdk.Coins{sdk.Coin{
										Denom:  "foo",
										Amount: sdkmath.OneInt(),
									}},
								},
								&sendingtypes.MsgCreateTransfer{
									Transfer: &sendingtypes.Transfer{
										Sender:    constants.Bob_Num0,
										Recipient: constants.Alice_Num0,
										AssetId:   assettypes.AssetUsdc.Id,
										Amount:    500_000_000, // $500
									},
								},
							},
							Authenticators: []uint64{0, 0},

							Fees:       constants.TestFeeCoins_5Cents,
							Gas:        300_000,
							AccountNum: []uint64{1},
							SeqNum:     []uint64{2},
							Signers:    []cryptotypes.PrivKey{constants.BobPrivateKey},

							ExpectedRespCode: aptypes.ErrMessageTypeVerification.ABCICode(),
							ExpectedLog:      aptypes.ErrMessageTypeVerification.Error(),
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId: false,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *aptypes.GenesisState) {
						genesisState.Params.IsSmartAccountActive = tc.smartAccountEnabled
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			for _, block := range tc.blocks {
				for _, msg := range block.Msgs {
					tx, err := testtx.GenTx(
						ctx,
						tApp.App.TxConfig(),
						msg.Msg,
						msg.Fees,
						msg.Gas,
						tApp.App.ChainID(),
						msg.AccountNum,
						msg.SeqNum,
						msg.Signers,
						msg.Signers,
						msg.Authenticators,
					)
					require.NoError(t, err)

					bytes, err := tApp.App.TxConfig().TxEncoder()(tx)
					if err != nil {
						panic(err)
					}
					checkTxReq := abcitypes.RequestCheckTx{
						Tx:   bytes,
						Type: abcitypes.CheckTxType_New,
					}

					resp := tApp.CheckTx(checkTxReq)
					require.Equal(
						t,
						msg.ExpectedRespCode,
						resp.Code,
						"Response code was not as expected",
					)
					require.Contains(
						t,
						resp.Log,
						msg.ExpectedLog,
						"Response log was not as expected",
					)
				}
				ctx = tApp.AdvanceToBlock(block.Block, testapp.AdvanceToBlockOptions{})
			}

			for orderId, shouldHaveOrder := range tc.expectedOrderIdsInMemclob {
				_, exists := tApp.App.ClobKeeper.MemClob.GetOrder(orderId)
				require.Equal(t, shouldHaveOrder, exists)
			}
		})
	}
}
