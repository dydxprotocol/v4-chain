package memclob

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func TestGetOrderFilledAmount(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	tests := map[string]struct {
		// State.
		orderFilledAmount satypes.BaseQuantums

		// Expectations.
		expectedOrderFilledAmount satypes.BaseQuantums
	}{
		"Returns 0 if the order ID isn't found": {
			orderFilledAmount: 0,

			expectedOrderFilledAmount: 0,
		},
		"Returns the order filled amount if the order ID is marked as filled in state": {
			orderFilledAmount: 10,

			expectedOrderFilledAmount: 10,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClobKeeper := mocks.MemClobKeeper{}
			memclob := NewMemClobPriceTimePriority(false)
			memclob.SetClobKeeper(&memClobKeeper)

			orderId := types.OrderId{
				SubaccountId: constants.Alice_Num0,
				ClientId:     0,
			}

			memClobKeeper.On("GetStatePosition", mock.Anything, mock.Anything, mock.Anything).
				Return(big.NewInt(0))

			memClobKeeper.On("GetOrderFillAmount", mock.Anything, orderId).
				Return(true, tc.orderFilledAmount, uint32(0))

			// Run the test case.
			filledAmount := memclob.GetOrderFilledAmount(ctx, orderId)
			require.Equal(t, tc.orderFilledAmount, filledAmount)
		})
	}
}
