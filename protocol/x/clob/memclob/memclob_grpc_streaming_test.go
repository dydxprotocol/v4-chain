package memclob

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetOffchainUpdatesForOrderbookSnapshot_Buy(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()

	clobKeeper := &mocks.MemClobKeeper{}
	clobKeeper.On(
		"GetOrderFillAmount",
		mock.Anything,
		mock.Anything,
	).Return(false, satypes.BaseQuantums(0), uint32(0))
	clobKeeper.On("SendOrderbookUpdates", mock.Anything, mock.Anything).Return()

	memclob := NewMemClobPriceTimePriority(false)
	memclob.SetClobKeeper(clobKeeper)

	memclob.CreateOrderbook(constants.ClobPair_Btc)

	orders := []types.Order{
		constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO,
		constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16,
		constants.Order_Bob_Num0_Id12_Clob0_Buy5_Price40_GTB20,
	}

	for _, order := range orders {
		memclob.mustAddOrderToOrderbook(ctx, order, false)
	}

	offchainUpdates := memclob.GetOffchainUpdatesForOrderbookSnapshot(
		ctx,
		constants.ClobPair_Btc.GetClobPairId(),
	)

	expected := types.NewOffchainUpdates()
	// Buy orders are in descending order.
	expected.Append(memclob.GetOrderbookUpdatesForOrderPlacement(ctx, orders[2]))
	expected.Append(memclob.GetOrderbookUpdatesForOrderPlacement(ctx, orders[0]))
	expected.Append(memclob.GetOrderbookUpdatesForOrderPlacement(ctx, orders[1]))

	require.Equal(t, expected, offchainUpdates)
}

func TestGetOffchainUpdatesForOrderbookSnapshot_Sell(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()

	clobKeeper := &mocks.MemClobKeeper{}
	clobKeeper.On(
		"GetOrderFillAmount",
		mock.Anything,
		mock.Anything,
	).Return(false, satypes.BaseQuantums(0), uint32(0))
	clobKeeper.On("SendOrderbookUpdates", mock.Anything, mock.Anything).Return()

	memclob := NewMemClobPriceTimePriority(false)
	memclob.SetClobKeeper(clobKeeper)

	memclob.CreateOrderbook(constants.ClobPair_Btc)

	orders := []types.Order{
		constants.Order_Bob_Num0_Id12_Clob0_Sell20_Price35_GTB32,
		constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
		constants.Order_Alice_Num0_Id1_Clob0_Sell15_Price10_GTB18_PO,
	}

	for _, order := range orders {
		memclob.mustAddOrderToOrderbook(ctx, order, false)
	}

	offchainUpdates := memclob.GetOffchainUpdatesForOrderbookSnapshot(
		ctx,
		constants.ClobPair_Btc.GetClobPairId(),
	)

	expected := types.NewOffchainUpdates()
	// Sell orders are in ascending order.
	expected.Append(memclob.GetOrderbookUpdatesForOrderPlacement(ctx, orders[1]))
	expected.Append(memclob.GetOrderbookUpdatesForOrderPlacement(ctx, orders[2]))
	expected.Append(memclob.GetOrderbookUpdatesForOrderPlacement(ctx, orders[0]))

	require.Equal(t, expected, offchainUpdates)
}
