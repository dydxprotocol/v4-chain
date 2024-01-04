package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestPerformStatelessMatchPerpetualDeleveragingValidation(t *testing.T) {
	tests := map[string]struct {
		match types.MatchPerpetualDeleveraging

		expectedError error
	}{
		"Success": {
			match: types.MatchPerpetualDeleveraging{
				Liquidated:  constants.Alice_Num0,
				PerpetualId: constants.ClobPair_Eth.MustGetPerpetualId(),
				Fills: []types.MatchPerpetualDeleveraging_Fill{
					{
						OffsettingSubaccountId: constants.Bob_Num0,
						FillAmount:             10,
					},
				},
			},
			expectedError: nil,
		},
		"Length of fills is zero": {
			match: types.MatchPerpetualDeleveraging{
				Liquidated:  constants.Alice_Num0,
				PerpetualId: constants.ClobPair_Eth.MustGetPerpetualId(),
				Fills:       []types.MatchPerpetualDeleveraging_Fill{},
			},
			expectedError: nil,
		},
		"Deleveraging fill subaccount id the same as liquidation subaccount id": {
			match: types.MatchPerpetualDeleveraging{
				Liquidated:  constants.Alice_Num0,
				PerpetualId: constants.ClobPair_Eth.MustGetPerpetualId(),
				Fills: []types.MatchPerpetualDeleveraging_Fill{
					{
						OffsettingSubaccountId: constants.Alice_Num0,
						FillAmount:             10,
					},
				},
			},
			expectedError: types.ErrDeleveragingAgainstSelf,
		},
		"Duplicate subaccount ids in fills": {
			match: types.MatchPerpetualDeleveraging{
				Liquidated:  constants.Alice_Num0,
				PerpetualId: constants.ClobPair_Eth.MustGetPerpetualId(),
				Fills: []types.MatchPerpetualDeleveraging_Fill{
					{
						OffsettingSubaccountId: constants.Bob_Num0,
						FillAmount:             10,
					},
					{
						OffsettingSubaccountId: constants.Bob_Num0,
						FillAmount:             5,
					},
				},
			},
			expectedError: types.ErrDuplicateDeleveragingFillSubaccounts,
		},
		"Zero fill amount": {
			match: types.MatchPerpetualDeleveraging{
				Liquidated:  constants.Alice_Num0,
				PerpetualId: constants.ClobPair_Eth.MustGetPerpetualId(),
				Fills: []types.MatchPerpetualDeleveraging_Fill{
					{
						OffsettingSubaccountId: constants.Bob_Num0,
						FillAmount:             10,
					},
					{
						OffsettingSubaccountId: constants.Bob_Num1,
						FillAmount:             0,
					},
				},
			},
			expectedError: types.ErrZeroDeleveragingFillAmount,
		},
		"Invalid liquidation subaccount id": {
			match: types.MatchPerpetualDeleveraging{
				Liquidated:  constants.InvalidSubaccountIdNumber,
				PerpetualId: constants.ClobPair_Eth.MustGetPerpetualId(),
				Fills: []types.MatchPerpetualDeleveraging_Fill{
					{
						OffsettingSubaccountId: constants.Bob_Num0,
						FillAmount:             10,
					},
					{
						OffsettingSubaccountId: constants.Bob_Num1,
						FillAmount:             5,
					},
				},
			},
			expectedError: satypes.ErrInvalidSubaccountIdNumber,
		},
		"Invalid deleveraging fill subaccount id": {
			match: types.MatchPerpetualDeleveraging{
				Liquidated:  constants.Alice_Num0,
				PerpetualId: constants.ClobPair_Eth.MustGetPerpetualId(),
				Fills: []types.MatchPerpetualDeleveraging_Fill{
					{
						OffsettingSubaccountId: constants.InvalidSubaccountIdNumber,
						FillAmount:             10,
					},
					{
						OffsettingSubaccountId: constants.Bob_Num1,
						FillAmount:             5,
					},
				},
			},
			expectedError: satypes.ErrInvalidSubaccountIdNumber,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.match.Validate()
			if tc.expectedError != nil {
				require.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
