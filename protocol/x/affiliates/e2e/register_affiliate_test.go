package affiliate_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	constants "github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	"github.com/stretchr/testify/require"
)

func TestRegisterAffiliateInvalidSigner(t *testing.T) {
	testCases := []struct {
		name          string
		referee       string
		affiliate     string
		signer        string
		expectSuccess bool
	}{
		{
			name:          "Valid signer (referee)",
			referee:       constants.BobAccAddress.String(),
			affiliate:     constants.AliceAccAddress.String(),
			signer:        constants.BobAccAddress.String(),
			expectSuccess: true,
		},
		{
			name:          "Invalid signer (affiliate)",
			referee:       constants.BobAccAddress.String(),
			affiliate:     constants.AliceAccAddress.String(),
			signer:        constants.AliceAccAddress.String(),
			expectSuccess: false,
		},
		{
			name:          "Invalid signer (non-related address)",
			referee:       constants.BobAccAddress.String(),
			affiliate:     constants.AliceAccAddress.String(),
			signer:        constants.CarlAccAddress.String(),
			expectSuccess: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()

			msgRegisterAffiliate := types.MsgRegisterAffiliate{
				Referee:   tc.referee,
				Affiliate: tc.affiliate,
			}

			checkTxMsgRegisterAffiliate := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: tc.signer,
					Gas:                  constants.TestGasLimit,
					FeeAmt:               constants.TestFeeCoins_5Cents,
				},
				&msgRegisterAffiliate,
			)
			checkTxResp := tApp.CheckTx(checkTxMsgRegisterAffiliate)

			if tc.expectSuccess {
				require.True(t, checkTxResp.IsOK(), "Expected CheckTx to succeed with valid signer")
			} else {
				require.True(t, checkTxResp.IsErr(), "Expected CheckTx to fail with invalid signer")
				require.Contains(t, checkTxResp.Log, "pubKey does not match signer address")
			}
		})
	}
}
