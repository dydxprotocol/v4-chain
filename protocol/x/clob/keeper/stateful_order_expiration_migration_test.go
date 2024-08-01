package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/store/prefix"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestUnsafeMigrateOrderExpirationState(t *testing.T) {
	tests := map[string]struct {
		timeSlicesToOrderIds map[time.Time][]types.OrderId
	}{
		"Multiple time slices": {
			timeSlicesToOrderIds: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15.OrderId,
				},
				constants.Time_21st_Feb_2021.Add(1): {
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTB30.OrderId,
				},
				constants.Time_21st_Feb_2021.Add(77): {
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
				},
			},
		},
		"No time slices": {
			timeSlicesToOrderIds: map[time.Time][]types.OrderId{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			for timestamp, orderIds := range tc.timeSlicesToOrderIds {
				ks.ClobKeeper.LegacySetStatefulOrdersTimeSliceInState(ks.Ctx, timestamp, orderIds)
			}

			ks.ClobKeeper.UnsafeMigrateOrderExpirationState(ks.Ctx)

			oldStore := prefix.NewStore(
				ks.Ctx.KVStore(ks.StoreKey),
				[]byte(types.LegacyStatefulOrdersTimeSlicePrefix), //nolint:staticcheck
			)
			it := oldStore.Iterator(nil, nil)
			defer it.Close()
			require.False(t, it.Valid())

			for goodTilTime, expectedOrderIds := range tc.timeSlicesToOrderIds {
				orderIds := ks.ClobKeeper.GetStatefulOrderIdExpirations(ks.Ctx, goodTilTime)
				require.ElementsMatch(
					t,
					expectedOrderIds,
					orderIds,
					"Mismatch of order IDs for timestamp",
					goodTilTime.String(),
				)
			}
		})
	}
}
