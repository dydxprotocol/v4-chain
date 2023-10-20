package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

// populateStatefulOrders adds stateful orders to the keeper: a long-term order, a triggered conditional order, and an
// untriggered conditional order.
func populateStatefulOrders(ks keepertest.ClobKeepersTestContext) {
	// Long term order.
	ks.ClobKeeper.SetLongTermOrderPlacement(
		ks.Ctx,
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
		1,
	)

	// Triggered conditional order.
	ks.ClobKeeper.SetLongTermOrderPlacement(
		ks.Ctx,
		constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
		1,
	)
	ks.ClobKeeper.MustTriggerConditionalOrder(
		ks.Ctx,
		constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId,
	)

	// Untriggered conditional order.
	ks.ClobKeeper.SetLongTermOrderPlacement(
		ks.Ctx,
		constants.ConditionalOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT15_StopLoss20,
		1,
	)
}

func TestAllStatefulOrders(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
	wctx := sdk.WrapSDKContext(ks.Ctx)

	populateStatefulOrders(ks)

	for name, tc := range map[string]struct {
		request  *types.QueryAllStatefulOrdersRequest
		response *types.QueryAllStatefulOrdersResponse
		err      error
	}{
		"Nil request returns an error": {
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
		"Success": {
			request: &types.QueryAllStatefulOrdersRequest{},
			response: &types.QueryAllStatefulOrdersResponse{
				StatefulOrders: []types.Order{
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
					constants.ConditionalOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT15_StopLoss20,
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			response, err := ks.ClobKeeper.AllStatefulOrders(wctx, tc.request)
			if tc.err != nil {
				require.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestStatefulOrderCount(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
	wctx := sdk.WrapSDKContext(ks.Ctx)

	populateStatefulOrders(ks)

	for name, tc := range map[string]struct {
		request  *types.QueryStatefulOrderCountRequest
		response *types.QueryStatefulOrderCountResponse
		err      error
	}{
		"Nil request returns an error": {
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
		"Success": {
			request: &types.QueryStatefulOrderCountRequest{
				SubaccountId: &constants.Alice_Num0,
			},
			response: &types.QueryStatefulOrderCountResponse{
				Count: 3,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			response, err := ks.ClobKeeper.StatefulOrderCount(wctx, tc.request)
			if tc.err != nil {
				require.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}
