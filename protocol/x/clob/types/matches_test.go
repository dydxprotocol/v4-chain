package types_test

import (
	"testing"

	types "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestGetMakerSubaccountIds_MatchOrders(t *testing.T) {
	subAccountId1 := satypes.SubaccountId{Owner: "owner1"}
	subAccountId2 := satypes.SubaccountId{Owner: "owner2"}
	m := types.MatchOrders{
		Fills: []types.MakerFill{
			{MakerOrderId: types.OrderId{SubaccountId: subAccountId1}},
			{MakerOrderId: types.OrderId{SubaccountId: subAccountId2}},
		},
	}

	result := m.GetMakerSubaccountIds()
	require.Equal(t, []satypes.SubaccountId{subAccountId1, subAccountId2}, result)
}

func TestGetMakerSubaccountIds_MatchPerpetualLiquidationNew(t *testing.T) {
	subAccountId1 := satypes.SubaccountId{Owner: "owner1"}
	subAccountId2 := satypes.SubaccountId{Owner: "owner2"}
	m := types.MatchPerpetualLiquidation{
		Fills: []types.MakerFill{
			{MakerOrderId: types.OrderId{SubaccountId: subAccountId1}},
			{MakerOrderId: types.OrderId{SubaccountId: subAccountId2}},
		},
	}

	result := m.GetMakerSubaccountIds()
	require.Equal(t, []satypes.SubaccountId{subAccountId1, subAccountId2}, result)
}
