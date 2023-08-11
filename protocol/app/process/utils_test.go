package process_test

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4/app/process"
	"github.com/stretchr/testify/require"
)

func TestGetAppInjectedMsgIdxMaps(t *testing.T) {
	tests := map[string]struct {
		numTxs int
	}{
		"NumTxs = 3":        {numTxs: 3},
		"NumTxs = 100_0000": {numTxs: 10_000},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			txTypeToIdx, idxToTxType := process.GetAppInjectedMsgIdxMaps(tc.numTxs)

			// Validate txTypeToIdx.
			require.Len(t, txTypeToIdx, 3)

			orderIdx, ok := txTypeToIdx[process.ProposedOperationsTxType]
			require.True(t, ok)
			require.Equal(t, 0, orderIdx)

			addFundingIdx, ok := txTypeToIdx[process.AddPremiumVotesTxType]
			require.True(t, ok)
			require.Equal(t, tc.numTxs-2, addFundingIdx)

			updatePricesIdx, ok := txTypeToIdx[process.UpdateMarketPricesTxType]
			require.True(t, ok)
			require.Equal(t, tc.numTxs-1, updatePricesIdx)

			// Validate idxToTxType.
			require.Len(t, idxToTxType, 3)
			operationsTxType, ok := idxToTxType[0]
			require.True(t, ok)
			require.Equal(t, process.ProposedOperationsTxType, operationsTxType)

			addFundingTxType, ok := idxToTxType[tc.numTxs-2]
			require.True(t, ok)
			require.Equal(t, process.AddPremiumVotesTxType, addFundingTxType)

			updatePricesTxType, ok := idxToTxType[tc.numTxs-1]
			require.True(t, ok)
			require.Equal(t, process.UpdateMarketPricesTxType, updatePricesTxType)
		})
	}
}

func TestGetAppInjectedMsgIdxMaps_Panic(t *testing.T) {
	tests := map[string]struct {
		numTxs int
	}{
		"NumTxs: negative": {numTxs: -10},
		"NumTxs: zero":     {numTxs: 0},
		"NumTxs: two":      {numTxs: 2},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.PanicsWithError(
				t,
				fmt.Errorf("num of txs must be at least 3").Error(),
				func() { _, _ = process.GetAppInjectedMsgIdxMaps(tc.numTxs) },
			)
		})
	}
}
