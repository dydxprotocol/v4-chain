package authz_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/gogoproto/proto"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	assetstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
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
		"Success: Alice grants permission to Bob to send from her account. Bob sends from Alice's account.": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Bob_Num0_100_000USD,
			},

			msgGrant: &authz.MsgGrant{
				Granter: constants.AliceAccAddress.String(),
				Grantee: constants.BobAccAddress.String(),
				Grant: authz.Grant{
					Authorization: newAny(
						authz.NewGenericAuthorization(sdk.MsgTypeURL(&banktypes.MsgSend{})),
					),
				},
			},

			msgExec: &authz.MsgExec{
				Grantee: constants.BobAccAddress.String(),
				Msgs: []*codectypes.Any{
					newAny(
						&banktypes.MsgSend{
							FromAddress: constants.AliceAccAddress.String(),
							ToAddress:   constants.BobAccAddress.String(),
							Amount: []sdk.Coin{
								sdk.NewCoin(assetstypes.AssetUsdc.Denom, sdkmath.NewInt(1)),
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
				aliceBalance := tApp.App.BankKeeper.GetBalance(
					ctx,
					constants.AliceAccAddress,
					assetstypes.AssetUsdc.Denom,
				)
				require.Equal(
					t,
					sdk.NewCoin(
						assetstypes.AssetUsdc.Denom,
						// Alice paid 5 cents in fees for MsgGrant.
						sdkmath.NewInt(100000000000000000-50000-1),
					),
					aliceBalance,
				)
				bobBalance := tApp.App.BankKeeper.GetBalance(
					ctx,
					constants.BobAccAddress,
					assetstypes.AssetUsdc.Denom,
				)
				require.Equal(
					t,
					sdk.NewCoin(
						assetstypes.AssetUsdc.Denom,
						// Bob paid 5 cents in fees for MsgExec.
						sdkmath.NewInt(100000000000000000-50000+1),
					),
					bobBalance,
				)
			},
		},
		"Fail (external): Bob tries to vote on behalf of Alice without permission.": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Bob_Num0_100_000USD,
			},

			msgGrant: nil,

			msgExec: &authz.MsgExec{
				Grantee: constants.BobAccAddress.String(),
				Msgs: []*codectypes.Any{
					newAny(
						&govtypesv1.MsgVote{
							ProposalId: 0,
							Voter:      constants.AliceAccAddress.String(),
						},
					),
				},
			},

			expectedMsgExecCheckTxSuccess:   true,
			expectedMsgExecCheckTxCode:      abcitypes.CodeTypeOK,
			expectedMsgExecDeliverTxSuccess: false,
			expectedMsgExecDeliverTxCode:    authz.ErrNoAuthorizationFound.ABCICode(),
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
						authz.NewGenericAuthorization(sdk.MsgTypeURL(&banktypes.MsgUpdateParams{})),
					),
				},
			},

			msgExec: &authz.MsgExec{
				Grantee: constants.BobAccAddress.String(),
				Msgs: []*codectypes.Any{
					newAny(
						&banktypes.MsgUpdateParams{
							Authority: lib.GovModuleAddress.String(),
						},
					),
				},
			},

			expectedMsgExecCheckTxSuccess:   true,
			expectedMsgExecCheckTxCode:      abcitypes.CodeTypeOK,
			expectedMsgExecDeliverTxSuccess: false,
			expectedMsgExecDeliverTxCode:    authz.ErrNoAuthorizationFound.ABCICode(),
		},
		"Fail (internal): Bob tries to update gov params (authority = gov)": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Bob_Num0_100_000USD,
			},

			msgGrant: nil,

			msgExec: &authz.MsgExec{
				Grantee: constants.BobAccAddress.String(),
				Msgs: []*codectypes.Any{
					newAny(
						&govtypesv1.MsgUpdateParams{
							// Authority = gov
							Authority: lib.GovModuleAddress.String(),
						},
					),
				},
			},

			expectedMsgExecCheckTxSuccess:   true,
			expectedMsgExecCheckTxCode:      abcitypes.CodeTypeOK,
			expectedMsgExecDeliverTxSuccess: false,
			expectedMsgExecDeliverTxCode:    authz.ErrNoAuthorizationFound.ABCICode(),
		},
		"Fail (internal): Bob tries to update gov params (authority = bob)": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Bob_Num0_100_000USD,
			},

			msgGrant: nil,

			msgExec: &authz.MsgExec{
				Grantee: constants.BobAccAddress.String(),
				Msgs: []*codectypes.Any{
					newAny(
						&govtypesv1.MsgUpdateParams{
							// Authority = bob
							Authority: constants.BobAccAddress.String(),
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
		"Fail (dydx custom): Bob wraps dydx messages": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Bob_Num0_100_000USD,
			},

			msgGrant: nil,

			msgExec: &authz.MsgExec{
				Grantee: constants.BobAccAddress.String(),
				Msgs: []*codectypes.Any{
					newAny(
						&clobtypes.MsgPlaceOrder{},
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
