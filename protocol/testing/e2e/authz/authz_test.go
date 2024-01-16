package authz_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/gogoproto/proto"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func newAny(v proto.Message) *codectypes.Any {
	any, err := codectypes.NewAnyWithValue(v)
	if err != nil {
		panic(err)
	}

	return any
}

func TestAuthz(t *testing.T) {
	tests := map[string]struct {
		subaccounts []satypes.Subaccount

		msgGrant *authz.MsgGrant
		msgExec  *authz.MsgExec

		expectedMsgExecCheckTxSuccess   bool
		expectedMsgExecCheckTxCode      uint32
		expectedMsgExecDeliverTxSuccess bool
		expectedMsgExecDeliverTxCode    uint32

		verifyResults func(ctx sdk.Context, tApp *testapp.TestApp)
	}{
		"Success: Alice grants permission to Bob to transfer from her account. Bob transfers from Alice's account.": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Bob_Num0_100_000USD,
			},

			msgGrant: &authz.MsgGrant{
				Granter: constants.AliceAccAddress.String(),
				Grantee: constants.BobAccAddress.String(),
				Grant: authz.Grant{
					Authorization: newAny(
						authz.NewGenericAuthorization(sdk.MsgTypeURL(&sendingtypes.MsgCreateTransfer{})),
					),
				},
			},

			msgExec: &authz.MsgExec{
				Grantee: constants.BobAccAddress.String(),
				Msgs: []*codectypes.Any{
					newAny(
						&sendingtypes.MsgCreateTransfer{
							Transfer: &sendingtypes.Transfer{
								Sender:    constants.Alice_Num0,
								Recipient: constants.Bob_Num0,
								AssetId:   0,
								Amount:    10_000_000_000, // $10,000
							},
						},
					),
				},
			},

			expectedMsgExecCheckTxSuccess:   true,
			expectedMsgExecCheckTxCode:      abcitypes.CodeTypeOK,
			expectedMsgExecDeliverTxSuccess: true,
			expectedMsgExecDeliverTxCode:    abcitypes.CodeTypeOK,

			verifyResults: func(ctx sdk.Context, tApp *testapp.TestApp) {
				expectedSubaccounts := []satypes.Subaccount{
					{
						Id: &constants.Alice_Num0,
						AssetPositions: testutil.CreateUsdcAssetPosition(
							big.NewInt(90_000_000_000),
						),
					},
					{
						Id: &constants.Bob_Num0,
						AssetPositions: testutil.CreateUsdcAssetPosition(
							big.NewInt(110_000_000_000),
						),
					},
				}
				for _, subaccount := range expectedSubaccounts {
					actualSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *subaccount.Id)
					require.Equal(t, subaccount, actualSubaccount)
				}
			},
		},
		`Success: Alice grants permission to Bob to place orders. Bob places some orders (note that
			this does allow Bob to bypass the rate limiter)`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Bob_Num0_100_000USD,
			},

			msgGrant: &authz.MsgGrant{
				Granter: constants.AliceAccAddress.String(),
				Grantee: constants.BobAccAddress.String(),
				Grant: authz.Grant{
					Authorization: newAny(
						authz.NewGenericAuthorization(sdk.MsgTypeURL(&clobtypes.MsgPlaceOrder{})),
					),
				},
			},

			msgExec: &authz.MsgExec{
				Grantee: constants.BobAccAddress.String(),
				Msgs: []*codectypes.Any{
					// The rate limiter is set to 1 order per 10 blocks, but this MsgExec contains
					// two order placements.
					newAny(
						&clobtypes.MsgPlaceOrder{
							Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT15,
						},
					),
					newAny(
						&clobtypes.MsgPlaceOrder{
							Order: constants.LongTermOrder_Alice_Num0_Id1_Clob0_Buy1BTC_Price50000_GTBT15,
						},
					),
				},
			},

			expectedMsgExecCheckTxSuccess:   true,
			expectedMsgExecCheckTxCode:      abcitypes.CodeTypeOK,
			expectedMsgExecDeliverTxSuccess: true,
			expectedMsgExecDeliverTxCode:    abcitypes.CodeTypeOK,

			verifyResults: func(ctx sdk.Context, tApp *testapp.TestApp) {
				require.Equal(t, 2, len(tApp.App.ClobKeeper.GetAllStatefulOrders(ctx)))
				orders, err := tApp.App.ClobKeeper.MemClob.GetSubaccountOrders(
					ctx,
					0,
					constants.Alice_Num0,
					clobtypes.Order_SIDE_BUY,
				)
				require.NoError(t, err)
				require.Equal(t, 2, len(orders))
			},
		},
		"Fail (external): Bob tries to transfer from Alice's account without permission.": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Bob_Num0_100_000USD,
			},

			msgGrant: nil,

			msgExec: &authz.MsgExec{
				Grantee: constants.BobAccAddress.String(),
				Msgs: []*codectypes.Any{
					newAny(
						&sendingtypes.MsgCreateTransfer{
							Transfer: &sendingtypes.Transfer{
								Sender:    constants.Alice_Num0,
								Recipient: constants.Bob_Num0,
								AssetId:   0,
								Amount:    10_000_000_000, // $10,000
							},
						},
					),
				},
			},

			expectedMsgExecCheckTxSuccess:   true,
			expectedMsgExecCheckTxCode:      abcitypes.CodeTypeOK,
			expectedMsgExecDeliverTxSuccess: false,
			expectedMsgExecDeliverTxCode:    authz.ErrNoAuthorizationFound.ABCICode(),

			verifyResults: func(ctx sdk.Context, tApp *testapp.TestApp) {
				expectedSubaccounts := []satypes.Subaccount{
					{
						Id: &constants.Alice_Num0,
						AssetPositions: testutil.CreateUsdcAssetPosition(
							big.NewInt(100_000_000_000),
						),
					},
					{
						Id: &constants.Bob_Num0,
						AssetPositions: testutil.CreateUsdcAssetPosition(
							big.NewInt(100_000_000_000),
						),
					},
				}
				for _, subaccount := range expectedSubaccounts {
					actualSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *subaccount.Id)
					require.Equal(t, subaccount, actualSubaccount)
				}
			},
		},
		"Fail (internal): Granting permissions to execute internal messages doesn't allow execution": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Bob_Num0_100_000USD,
			},

			msgGrant: &authz.MsgGrant{
				Granter: constants.AliceAccAddress.String(),
				Grantee: constants.BobAccAddress.String(),
				Grant: authz.Grant{
					Authorization: newAny(
						authz.NewGenericAuthorization(sdk.MsgTypeURL(&clobtypes.MsgCreateClobPair{})),
					),
				},
			},

			msgExec: &authz.MsgExec{
				Grantee: constants.BobAccAddress.String(),
				Msgs: []*codectypes.Any{
					newAny(
						&clobtypes.MsgCreateClobPair{
							Authority: lib.GovModuleAddress.String(),
							ClobPair:  constants.ClobPair_Btc2,
						},
					),
				},
			},

			expectedMsgExecCheckTxSuccess:   true,
			expectedMsgExecCheckTxCode:      abcitypes.CodeTypeOK,
			expectedMsgExecDeliverTxSuccess: false,
			expectedMsgExecDeliverTxCode:    authz.ErrNoAuthorizationFound.ABCICode(),

			verifyResults: func(ctx sdk.Context, tApp *testapp.TestApp) {
				// Verify no clob pairs were created.
				require.Equal(t, 2, len(tApp.App.ClobKeeper.GetAllClobPairs(ctx)))
			},
		},
		"Fail (internal): Bob tries to create a new clob pair (authority = gov)": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Bob_Num0_100_000USD,
			},

			msgGrant: nil,

			msgExec: &authz.MsgExec{
				Grantee: constants.BobAccAddress.String(),
				Msgs: []*codectypes.Any{
					newAny(
						&clobtypes.MsgCreateClobPair{
							// Authority = gov
							Authority: lib.GovModuleAddress.String(),
							ClobPair:  constants.ClobPair_Btc2,
						},
					),
				},
			},

			expectedMsgExecCheckTxSuccess:   true,
			expectedMsgExecCheckTxCode:      abcitypes.CodeTypeOK,
			expectedMsgExecDeliverTxSuccess: false,
			expectedMsgExecDeliverTxCode:    authz.ErrNoAuthorizationFound.ABCICode(),

			verifyResults: func(ctx sdk.Context, tApp *testapp.TestApp) {
				// Verify no clob pairs were created.
				require.Equal(t, 2, len(tApp.App.ClobKeeper.GetAllClobPairs(ctx)))
			},
		},
		"Fail (internal): Bob tries to create a new clob pair (authority = bob)": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Bob_Num0_100_000USD,
			},

			msgGrant: nil,

			msgExec: &authz.MsgExec{
				Grantee: constants.BobAccAddress.String(),
				Msgs: []*codectypes.Any{
					newAny(
						&clobtypes.MsgCreateClobPair{
							// Authority = bob
							Authority: constants.BobAccAddress.String(),
							ClobPair:  constants.ClobPair_Btc2,
						},
					),
				},
			},

			expectedMsgExecCheckTxSuccess:   true,
			expectedMsgExecCheckTxCode:      abcitypes.CodeTypeOK,
			expectedMsgExecDeliverTxSuccess: false,
			// This fails because Bob is not authorized to create clob pairs (not in the
			// list of authorities defined in clobkeeper.authorities).
			expectedMsgExecDeliverTxCode: govtypes.ErrInvalidSigner.ABCICode(),

			verifyResults: func(ctx sdk.Context, tApp *testapp.TestApp) {
				// Verify no clob pairs were created.
				require.Equal(t, 2, len(tApp.App.ClobKeeper.GetAllClobPairs(ctx)))
			},
		},
		//
		// Below tests fail during CheckTx since the ante handler would reject these transactions.
		//
		"Fail (app injected): Bob tries to propose operations": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Bob_Num0_100_000USD,
			},

			msgGrant: nil,

			msgExec: &authz.MsgExec{
				Grantee: constants.BobAccAddress.String(),
				Msgs: []*codectypes.Any{
					newAny(
						&clobtypes.MsgProposedOperations{},
					),
				},
			},

			expectedMsgExecCheckTxSuccess: false,
			expectedMsgExecCheckTxCode:    sdkerrors.ErrInvalidRequest.ABCICode(),
		},
		"Fail (double nested): Bob wraps another nested message": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Bob_Num0_100_000USD,
			},

			msgGrant: nil,

			msgExec: &authz.MsgExec{
				Grantee: constants.BobAccAddress.String(),
				Msgs: []*codectypes.Any{
					newAny(
						&authz.MsgExec{},
					),
				},
			},

			expectedMsgExecCheckTxSuccess: false,
			expectedMsgExecCheckTxCode:    sdkerrors.ErrInvalidRequest.ABCICode(),
		},
		"Fail (unsupported): Bob wraps unspported transactions": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Bob_Num0_100_000USD,
			},

			msgGrant: nil,

			msgExec: &authz.MsgExec{
				Grantee: constants.BobAccAddress.String(),
				Msgs: []*codectypes.Any{
					newAny(
						&icacontrollertypes.MsgUpdateParams{},
					),
				},
			},

			expectedMsgExecCheckTxSuccess: false,
			expectedMsgExecCheckTxCode:    sdkerrors.ErrInvalidRequest.ABCICode(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
			appOpts := map[string]interface{}{
				indexer.MsgSenderInstanceForTest: msgSender,
			}

			// Initialize test app
			tApp := testapp.NewTestAppBuilder(t).WithAppOptions(appOpts).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = tc.subaccounts
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.BlockRateLimitConfig = clobtypes.BlockRateLimitConfiguration{
							MaxStatefulOrdersPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
								{
									NumBlocks: 10,
									Limit:     1,
								},
							},
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			if tc.msgGrant != nil {
				for _, checkTx := range testapp.MustMakeCheckTxsWithSdkMsg(
					ctx,
					tApp.App,
					testapp.MustMakeCheckTxOptions{
						AccAddressForSigning: tc.msgGrant.Granter, // Granter
						Gas:                  1000000,
						FeeAmt:               constants.TestFeeCoins_5Cents,
					},
					tc.msgGrant,
				) {
					resp := tApp.CheckTx(checkTx)
					require.True(
						t,
						resp.IsOK(),
						"Expected CheckTx to succeed. Response: %+v",
						resp,
					)
				}
			}

			// Give grantee some permissions.
			ctx = tApp.AdvanceToBlock(uint32(ctx.BlockHeight())+1, testapp.AdvanceToBlockOptions{})

			// Grantee executes a msg.
			for _, checkTx := range testapp.MustMakeCheckTxsWithSdkMsg(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: tc.msgExec.Grantee, // Grantee
					Gas:                  1000000,
					FeeAmt:               constants.TestFeeCoins_5Cents,
				},
				tc.msgExec,
			) {
				resp := tApp.CheckTx(checkTx)
				require.Equal(
					t,
					tc.expectedMsgExecCheckTxSuccess,
					resp.IsOK(),
					"Expected CheckTx to succeed. Response: %+v",
					resp,
				)
				require.Equal(t, tc.expectedMsgExecCheckTxCode, resp.Code)
			}

			if tc.expectedMsgExecCheckTxSuccess {
				ctx = tApp.AdvanceToBlock(uint32(ctx.BlockHeight())+1, testapp.AdvanceToBlockOptions{
					ValidateFinalizeBlock: func(
						ctx sdk.Context,
						request abcitypes.RequestFinalizeBlock,
						response abcitypes.ResponseFinalizeBlock,
					) (haltchain bool) {
						// Note the first TX is MsgProposeOperations.
						txResult := response.TxResults[1]
						require.Equal(t, tc.expectedMsgExecDeliverTxSuccess, txResult.IsOK())
						require.Equal(t, tc.expectedMsgExecDeliverTxCode, txResult.Code)
						return false
					},
				})
			}

			// Verify results.
			if tc.verifyResults != nil {
				ctx = tApp.AdvanceToBlock(uint32(ctx.BlockHeight())+1, testapp.AdvanceToBlockOptions{})
				tc.verifyResults(ctx, tApp)
			}
		})
	}
}
