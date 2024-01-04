package prepare_test

import (
	"errors"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/prepare"
	"github.com/stretchr/testify/require"
)

type TestFunction int

const (
	testUpdateMarketPrices TestFunction = iota
	testAcknowledgeBridges
	testAddPremiumVotes
	testProposedOperations
)

func Test_NewPrepareProposalTransactions_Success(t *testing.T) {
	req := &abci.RequestPrepareProposal{
		MaxTxBytes: 123,
	}
	ppt, err := prepare.NewPrepareProposalTxs(req)

	require.NoError(t, err)
	require.Equal(t, uint64(123), ppt.MaxBytes)
	require.Equal(t, uint64(0), ppt.UsedBytes)

	require.Nil(t, ppt.UpdateMarketPricesTx)
	require.Nil(t, ppt.AddPremiumVotesTx)
	require.Nil(t, ppt.AcknowledgeBridgesTx)
	require.Nil(t, ppt.ProposedOperationsTx)
	require.Nil(t, ppt.OtherTxs)
}

func Test_NewPrepareProposalTransactions_Fail(t *testing.T) {
	req := &abci.RequestPrepareProposal{
		MaxTxBytes: 0,
	}
	ppt, err := prepare.NewPrepareProposalTxs(req)

	require.ErrorContains(t, err, "MaxTxBytes must be positive")
	require.Equal(t, prepare.PrepareProposalTxs{}, ppt)
}

func Test_SetUpdateMarketPricesTx(t *testing.T) {
	setterTestCases(t, testUpdateMarketPrices)
}

func Test_SetAcknowledgeBridgesTx(t *testing.T) {
	setterTestCases(t, testAcknowledgeBridges)
}

func Test_SetAddPremiumVotesTx(t *testing.T) {
	setterTestCases(t, testAddPremiumVotes)
}

func Test_SetProposedOperationsTx(t *testing.T) {
	setterTestCases(t, testProposedOperations)
}

func setterTestCases(t *testing.T, tFunc TestFunction) {
	tests := map[string]struct {
		tx []byte

		expectedTx        []byte
		expectedUsedBytes uint64
		expectedErr       error
	}{
		"input is nil": {
			tx:                nil,
			expectedTx:        nil,
			expectedUsedBytes: 0,
		},
		"input is empty": {
			tx:                []byte{},
			expectedTx:        []byte{},
			expectedUsedBytes: 0,
		},
		"input is valid": {
			tx:                []byte{1, 2, 3, 4},
			expectedTx:        []byte{1, 2, 3, 4},
			expectedUsedBytes: 4,
		},
		"input overflows": {
			tx:                []byte{1, 2, 3, 4, 5},
			expectedTx:        nil,
			expectedUsedBytes: 0,
			expectedErr:       errors.New("Exceeds max: max=4, used=0, adding=5"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ppt, err := prepare.NewPrepareProposalTxs(
				&abci.RequestPrepareProposal{
					MaxTxBytes: 4,
				},
			)
			require.NoError(t, err)
			require.Equal(t, uint64(4), ppt.MaxBytes)
			require.Equal(t, uint64(0), ppt.UsedBytes)

			err = setterTestHelper(tFunc, &ppt, tc.tx)

			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expectedTx, getterTestHelper(tFunc, &ppt))
			require.Equal(t, tc.expectedUsedBytes, ppt.UsedBytes)
		})
	}
}

func setterTestHelper(tFunc TestFunction, target *prepare.PrepareProposalTxs, value []byte) error {
	switch tFunc {
	case testUpdateMarketPrices:
		return target.SetUpdateMarketPricesTx(value)
	case testAddPremiumVotes:
		return target.SetAddPremiumVotesTx(value)
	case testAcknowledgeBridges:
		return target.SetAcknowledgeBridgesTx(value)
	case testProposedOperations:
		return target.SetProposedOperationsTx(value)
	default:
		panic("not supported")
	}
}

func getterTestHelper(tFunc TestFunction, target *prepare.PrepareProposalTxs) []byte {
	switch tFunc {
	case testUpdateMarketPrices:
		return target.UpdateMarketPricesTx
	case testAddPremiumVotes:
		return target.AddPremiumVotesTx
	case testAcknowledgeBridges:
		return target.AcknowledgeBridgesTx
	case testProposedOperations:
		return target.ProposedOperationsTx
	default:
		panic("not supported")
	}
}

func Test_AddOtherTxs(t *testing.T) {
	tests := map[string]struct {
		txs           [][]byte
		additionalTxs [][]byte

		expectedTxs           [][]byte
		expectedUsedBytes     uint64
		expectedErr           error
		expectedAdditionalErr error
	}{
		"input is nil": {
			txs:           nil,
			additionalTxs: nil,

			expectedTxs:           nil,
			expectedUsedBytes:     0,
			expectedErr:           errors.New("No txs to add."),
			expectedAdditionalErr: errors.New("No txs to add."),
		},
		"input is empty": {
			txs:           [][]byte{{}},
			additionalTxs: [][]byte{{1}, {}},

			expectedTxs:           nil,
			expectedUsedBytes:     0,
			expectedErr:           errors.New("Cannot add zero length tx: []"),
			expectedAdditionalErr: errors.New("Cannot add zero length tx: []"),
		},
		"input is valid": {
			txs:               [][]byte{{1}, {2}, {3}},
			additionalTxs:     [][]byte{{4}},
			expectedTxs:       [][]byte{{1}, {2}, {3}, {4}},
			expectedUsedBytes: 4,
		},
		"input exceeds max on first attempt": {
			txs:                   [][]byte{{1}, {2}, {3}, {4}, {5}},
			additionalTxs:         [][]byte{{1}},
			expectedTxs:           [][]byte{{1}},
			expectedUsedBytes:     1,
			expectedErr:           errors.New("Exceeds max: max=4, used=0, adding=5"),
			expectedAdditionalErr: nil,
		},
		"input exceeds max on second attempt": {
			txs:                   [][]byte{{1}, {2}, {3}, {4}},
			additionalTxs:         [][]byte{{5}},
			expectedTxs:           [][]byte{{1}, {2}, {3}, {4}},
			expectedUsedBytes:     4,
			expectedErr:           nil,
			expectedAdditionalErr: errors.New("Exceeds max: max=4, used=4, adding=1"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ppt, err := prepare.NewPrepareProposalTxs(
				&abci.RequestPrepareProposal{
					MaxTxBytes: 4,
				},
			)
			require.NoError(t, err)
			require.Equal(t, uint64(4), ppt.MaxBytes)
			require.Equal(t, uint64(0), ppt.UsedBytes)

			// initial txs.
			err = ppt.AddOtherTxs(tc.txs)
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			// additional txs.
			err = ppt.AddOtherTxs(tc.additionalTxs)
			if tc.expectedAdditionalErr != nil {
				require.ErrorContains(t, err, tc.expectedAdditionalErr.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expectedTxs, ppt.OtherTxs)
			require.Equal(t, tc.expectedUsedBytes, ppt.UsedBytes)
		})
	}
}

func Test_UpdateUsedBytes(t *testing.T) {
	tests := map[string]struct {
		usedBytes     uint64
		bytesToRemove uint64
		bytesToAdd    uint64

		expectedErr error
	}{
		"Valid: replaced > add": {
			usedBytes:     4,
			bytesToRemove: 4,
			bytesToAdd:    2,
		},
		"Valid: replaced = add": {
			usedBytes:     5,
			bytesToRemove: 3,
			bytesToAdd:    3,
		},
		"Valid: replaced < add": {
			usedBytes:     5,
			bytesToRemove: 3,
			bytesToAdd:    5,
		},
		"Cannot be Negative": {
			usedBytes:     0,
			bytesToRemove: 3,
			bytesToAdd:    2,
			expectedErr:   errors.New("Result cannot be negative"),
		},
		"Exceeds max": {
			usedBytes:     0,
			bytesToRemove: 0,
			bytesToAdd:    11,
			expectedErr:   errors.New("Exceeds max: max=10, used=0, adding=11"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ppt, err := prepare.NewPrepareProposalTxs(
				&abci.RequestPrepareProposal{
					MaxTxBytes: 10,
				},
			)
			require.NoError(t, err)
			require.Equal(t, uint64(10), ppt.MaxBytes)
			require.Equal(t, uint64(0), ppt.UsedBytes)

			ppt.UsedBytes = tc.usedBytes

			err = ppt.UpdateUsedBytes(tc.bytesToRemove, tc.bytesToAdd)
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_GetAvailableBytes(t *testing.T) {
	tests := map[string]struct {
		pricesTx           []byte
		fundingTx          []byte
		operationsTx       []byte
		otherTxs           [][]byte
		otherAdditionalTxs [][]byte

		expectedUsedBytes  uint64
		expectedAvailBytes uint64
	}{
		"inputs are nil": {
			expectedAvailBytes: 10,
		},
		"inputs are empty": {
			pricesTx:           []byte{},
			fundingTx:          []byte{},
			operationsTx:       []byte{},
			otherTxs:           [][]byte{},
			otherAdditionalTxs: [][]byte{},

			expectedUsedBytes:  0,
			expectedAvailBytes: 10,
		},
		"some are set": {
			pricesTx:           []byte{1},
			fundingTx:          []byte{},
			operationsTx:       []byte{2, 3},
			otherTxs:           [][]byte{},
			otherAdditionalTxs: [][]byte{{4}, {5, 6}},

			expectedUsedBytes:  6,
			expectedAvailBytes: 4,
		},
		"all are set": {
			pricesTx:           []byte{1},
			fundingTx:          []byte{2, 3},
			operationsTx:       []byte{4},
			otherTxs:           [][]byte{{5}},
			otherAdditionalTxs: [][]byte{{6}, {7}},

			expectedUsedBytes:  7,
			expectedAvailBytes: 3,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ppt, err := prepare.NewPrepareProposalTxs(
				&abci.RequestPrepareProposal{
					MaxTxBytes: 10,
				},
			)
			require.NoError(t, err)
			require.Equal(t, uint64(10), ppt.MaxBytes)
			require.Equal(t, uint64(0), ppt.UsedBytes)

			err = ppt.SetUpdateMarketPricesTx(tc.pricesTx)
			require.NoError(t, err)

			err = ppt.SetAddPremiumVotesTx(tc.fundingTx)
			require.NoError(t, err)

			err = ppt.SetProposedOperationsTx(tc.operationsTx)
			require.NoError(t, err)

			// initial txs.
			err = ppt.AddOtherTxs(tc.otherTxs)
			if len(tc.otherTxs) == 0 {
				require.ErrorContains(t, err, "No txs to add.")
			} else {
				require.NoError(t, err)
			}

			// additional txs.
			err = ppt.AddOtherTxs(tc.otherAdditionalTxs)
			if len(tc.otherAdditionalTxs) == 0 {
				require.ErrorContains(t, err, "No txs to add.")
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expectedUsedBytes, ppt.UsedBytes)
			require.Equal(t, tc.expectedAvailBytes, ppt.GetAvailableBytes())
		})
	}
}

func Test_GetTxsInOrder(t *testing.T) {
	tests := map[string]struct {
		operationsTx       []byte
		otherTxs           [][]byte
		otherAdditionalTxs [][]byte
		bridgeTx           []byte
		fundingTx          []byte
		pricesTx           []byte

		expectedTxs [][]byte
		expectedErr error
	}{
		"update market prices is not set": {
			operationsTx:       []byte{},
			otherTxs:           [][]byte{},
			otherAdditionalTxs: [][]byte{},
			bridgeTx:           []byte{},
			fundingTx:          []byte{},
			pricesTx:           []byte{},

			expectedTxs: nil,
			expectedErr: errors.New("UpdateMarketPricesTx must be set"),
		},
		"add funding samples is not set": {
			operationsTx:       []byte{},
			otherTxs:           [][]byte{},
			otherAdditionalTxs: [][]byte{},
			bridgeTx:           []byte{},
			fundingTx:          []byte{},
			pricesTx:           []byte{1},

			expectedTxs: nil,
			expectedErr: errors.New("AddPremiumVotesTx must be set"),
		},
		"acknowledge bridges is not set": {
			operationsTx:       []byte{},
			otherTxs:           [][]byte{},
			otherAdditionalTxs: [][]byte{},
			bridgeTx:           []byte{},
			fundingTx:          []byte{2},
			pricesTx:           []byte{1},

			expectedTxs: nil,
			expectedErr: errors.New("AcknowledgeBridgesTx must be set"),
		},
		"prices, funding, and bridge only": {
			operationsTx:       []byte{},
			otherTxs:           [][]byte{},
			otherAdditionalTxs: [][]byte{},
			bridgeTx:           []byte{4},
			fundingTx:          []byte{2, 3},
			pricesTx:           []byte{1},

			expectedTxs: [][]byte{{4}, {2, 3}, {1}},
			expectedErr: nil,
		},
		"prices, funding, bridge + matched orders": {
			operationsTx:       []byte{4, 5, 6},
			otherTxs:           [][]byte{},
			otherAdditionalTxs: [][]byte{},
			bridgeTx:           []byte{3},
			fundingTx:          []byte{2},
			pricesTx:           []byte{1},

			expectedTxs: [][]byte{{4, 5, 6}, {3}, {2}, {1}},
			expectedErr: nil,
		},
		"prices, funding, bridge + others": {
			operationsTx:       []byte{},
			otherTxs:           [][]byte{{4}, {5, 6}},
			otherAdditionalTxs: [][]byte{},
			bridgeTx:           []byte{3},
			fundingTx:          []byte{2},
			pricesTx:           []byte{1},

			expectedTxs: [][]byte{{4}, {5, 6}, {3}, {2}, {1}},
			expectedErr: nil,
		},
		"partially set": {
			operationsTx:       []byte{4, 5, 6},
			otherTxs:           [][]byte{{7, 8}, {9, 10}},
			otherAdditionalTxs: [][]byte{},
			bridgeTx:           []byte{11},
			fundingTx:          []byte{2, 3},
			pricesTx:           []byte{1},

			expectedTxs: [][]byte{{4, 5, 6}, {7, 8}, {9, 10}, {11}, {2, 3}, {1}},
			expectedErr: nil,
		},
		"all set": {
			operationsTx:       []byte{4, 5},
			otherTxs:           [][]byte{{6}, {7, 8}},
			otherAdditionalTxs: [][]byte{{9}, {10}},
			bridgeTx:           []byte{11},
			fundingTx:          []byte{2, 3},
			pricesTx:           []byte{1},

			expectedTxs: [][]byte{{4, 5}, {6}, {7, 8}, {9}, {10}, {11}, {2, 3}, {1}},
			expectedErr: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ppt, err := prepare.NewPrepareProposalTxs(
				&abci.RequestPrepareProposal{
					MaxTxBytes: 11,
				},
			)
			require.NoError(t, err)
			require.Equal(t, uint64(11), ppt.MaxBytes)
			require.Equal(t, uint64(0), ppt.UsedBytes)

			err = ppt.SetUpdateMarketPricesTx(tc.pricesTx)
			require.NoError(t, err)

			err = ppt.SetAddPremiumVotesTx(tc.fundingTx)
			require.NoError(t, err)

			err = ppt.SetAcknowledgeBridgesTx(tc.bridgeTx)
			require.NoError(t, err)

			err = ppt.SetProposedOperationsTx(tc.operationsTx)
			require.NoError(t, err)

			// initial txs.
			err = ppt.AddOtherTxs(tc.otherTxs)
			if len(tc.otherTxs) == 0 {
				require.ErrorContains(t, err, "No txs to add.")
			} else {
				require.NoError(t, err)
			}

			// additional txs.
			err = ppt.AddOtherTxs(tc.otherAdditionalTxs)
			if len(tc.otherAdditionalTxs) == 0 {
				require.ErrorContains(t, err, "No txs to add.")
			} else {
				require.NoError(t, err)
			}

			txs, err := ppt.GetTxsInOrder()
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedTxs, txs)
		})
	}
}
