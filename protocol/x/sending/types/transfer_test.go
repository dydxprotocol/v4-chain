package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib/int256"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"github.com/stretchr/testify/require"
)

func TestGetBigQuoteQuantums(t *testing.T) {
	transfer := constants.Transfer_Carl_Num0_Dave_Num0_Quote_500
	quoteQuantums := transfer.GetQuantums()
	require.Equal(t, new(int256.Int).SetUint64(500_000_000), quoteQuantums)
}

func TestGetSubaccountUpdates(t *testing.T) {
	tests := map[string]struct {
		transfer types.Transfer
		expected []satypes.Update
	}{
		"Test subaccount updates": {
			transfer: constants.Transfer_Carl_Num0_Dave_Num0_Quote_500,
			expected: []satypes.Update{
				{
					SubaccountId: constants.Carl_Num0,
					AssetUpdates: testutil.CreateUsdcAssetUpdate(int256.NewInt(-500_000_000)),
				},
				{
					SubaccountId: constants.Dave_Num0,
					AssetUpdates: testutil.CreateUsdcAssetUpdate(int256.NewInt(500_000_000)),
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := []satypes.Update{
				tc.transfer.GetSenderSubaccountUpdate(),
				tc.transfer.GetRecipientSubaccountUpdate(),
			}
			require.Equal(t, tc.expected, result)
		})
	}
}
