package prepare_test

import (
	"errors"
	"fmt"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/prepare"
	"github.com/stretchr/testify/require"
)

// protoTxSize mirrors the unexported `prepare.txProtoSize`: the number of bytes a tx contributes to
// a block's data size, matching how CometBFT accounts for it in `Txs.Validate` (1 tag byte + a
// varint-encoded length prefix + the tx payload). Tests use this (instead of raw `len(tx)`) to derive
// expected budgets so the numbers document the framing rather than being magic constants.
func protoTxSize(tx []byte) uint64 {
	return uint64(cmttypes.ComputeProtoSizeForTxs([]cmttypes.Tx{tx}))
}

// uvarintProtoLen returns the number of bytes a protobuf varint encoding of n occupies, independent
// of cmttypes, so `Test_UsedBytes_ProtoFraming` can cross-check `protoTxSize` against an
// independently-derived formula rather than one that trivially matches by construction.
func uvarintProtoLen(n uint64) uint64 {
	size := uint64(1)
	for n >= 0x80 {
		n >>= 7
		size++
	}
	return size
}

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
	// MaxBytes is set to exactly the proto-framed size of the 4-byte "valid" tx below (1 tag byte +
	// 1 varint-len byte + 4 payload bytes = 6), so that tx exactly fits and a tx one byte longer (5
	// payload bytes, whose varint-len byte pushes it to 7 proto bytes) overflows by exactly 1 byte.
	maxBytes := protoTxSize([]byte{1, 2, 3, 4})
	overflowTxProtoSize := protoTxSize([]byte{1, 2, 3, 4, 5})

	tests := map[string]struct {
		tx []byte

		expectedTx        []byte
		expectedUsedBytes uint64
		expectedErr       error
	}{
		"input is nil": {
			tx:                nil,
			expectedTx:        nil,
			expectedUsedBytes: 0, // an unset slot contributes nothing to the block.
		},
		"input is empty": {
			tx:                []byte{},
			expectedTx:        []byte{},
			expectedUsedBytes: 0,
		},
		"input is valid": {
			tx:                []byte{1, 2, 3, 4},
			expectedTx:        []byte{1, 2, 3, 4},
			expectedUsedBytes: maxBytes, // exactly fits: proto size of a 4-byte tx equals MaxBytes.
		},
		"input overflows": {
			tx:                []byte{1, 2, 3, 4, 5},
			expectedTx:        nil,
			expectedUsedBytes: 0,
			expectedErr: fmt.Errorf(
				"Exceeds max: max=%d, used=0, adding=%d", maxBytes, overflowTxProtoSize,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ppt, err := prepare.NewPrepareProposalTxs(
				&abci.RequestPrepareProposal{
					MaxTxBytes: int64(maxBytes), //nolint:gosec
				},
			)
			require.NoError(t, err)
			require.Equal(t, maxBytes, ppt.MaxBytes)
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
	// Each 1-byte tx costs 1 (tag) + 1 (varint(1)) + 1 (payload) = 3 proto bytes, not 1 raw byte.
	// MaxBytes is set to exactly 4 such txs (12), so 4 one-byte txs exactly fit and a 5th overflows.
	oneByteTxProtoSize := protoTxSize([]byte{1})
	maxBytes := 4 * oneByteTxProtoSize

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
			expectedUsedBytes: 4 * oneByteTxProtoSize, // 4 one-byte txs, exactly fits MaxBytes.
		},
		"input exceeds max on first attempt": {
			txs:               [][]byte{{1}, {2}, {3}, {4}, {5}},
			additionalTxs:     [][]byte{{1}},
			expectedTxs:       [][]byte{{1}},
			expectedUsedBytes: oneByteTxProtoSize, // only the additional tx made it in.
			expectedErr: fmt.Errorf(
				"Exceeds max: max=%d, used=0, adding=%d", maxBytes, 5*oneByteTxProtoSize,
			),
			expectedAdditionalErr: nil,
		},
		"input exceeds max on second attempt": {
			txs:               [][]byte{{1}, {2}, {3}, {4}},
			additionalTxs:     [][]byte{{5}},
			expectedTxs:       [][]byte{{1}, {2}, {3}, {4}},
			expectedUsedBytes: maxBytes, // first 4 txs exactly fill the budget.
			expectedErr:       nil,
			expectedAdditionalErr: fmt.Errorf(
				"Exceeds max: max=%d, used=%d, adding=%d", maxBytes, maxBytes, oneByteTxProtoSize,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ppt, err := prepare.NewPrepareProposalTxs(
				&abci.RequestPrepareProposal{
					MaxTxBytes: int64(maxBytes), //nolint:gosec
				},
			)
			require.NoError(t, err)
			require.Equal(t, maxBytes, ppt.MaxBytes)
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
	// MaxBytes is set to exactly the proto-framed sum of the "all are set" case below (3+4+3+3+3+3
	// = 19), so that case exactly exhausts the budget (0 bytes available) while the smaller cases
	// leave room to spare. An unset (nil/empty) slot contributes nothing to the block and costs 0.
	maxBytes := protoTxSize([]byte{1}) + protoTxSize([]byte{2, 3}) + protoTxSize([]byte{4}) +
		protoTxSize([]byte{5}) + protoTxSize([]byte{6}) + protoTxSize([]byte{7})

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
			expectedUsedBytes:  0,
			expectedAvailBytes: maxBytes,
		},
		"inputs are empty": {
			pricesTx:           []byte{},
			fundingTx:          []byte{},
			operationsTx:       []byte{},
			otherTxs:           [][]byte{},
			otherAdditionalTxs: [][]byte{},

			expectedUsedBytes:  0,
			expectedAvailBytes: maxBytes,
		},
		"some are set": {
			pricesTx:           []byte{1},
			fundingTx:          []byte{},
			operationsTx:       []byte{2, 3},
			otherTxs:           [][]byte{},
			otherAdditionalTxs: [][]byte{{4}, {5, 6}},

			// prices(3) + funding(0, empty) + operations(4) + otherAdditional(3+4=7) = 14.
			expectedUsedBytes: protoTxSize([]byte{1}) +
				protoTxSize([]byte{2, 3}) + protoTxSize([]byte{4}) + protoTxSize([]byte{5, 6}),
			expectedAvailBytes: maxBytes - (protoTxSize([]byte{1}) +
				protoTxSize([]byte{2, 3}) + protoTxSize([]byte{4}) + protoTxSize([]byte{5, 6})),
		},
		"all are set": {
			pricesTx:           []byte{1},
			fundingTx:          []byte{2, 3},
			operationsTx:       []byte{4},
			otherTxs:           [][]byte{{5}},
			otherAdditionalTxs: [][]byte{{6}, {7}},

			expectedUsedBytes:  maxBytes, // exactly exhausts the budget.
			expectedAvailBytes: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ppt, err := prepare.NewPrepareProposalTxs(
				&abci.RequestPrepareProposal{
					MaxTxBytes: int64(maxBytes), //nolint:gosec
				},
			)
			require.NoError(t, err)
			require.Equal(t, maxBytes, ppt.MaxBytes)
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
	// MaxBytes is set to exactly the proto-framed sum of the "all set" case below, so every other
	// (smaller) case fits comfortably and "all set" exactly exhausts the budget:
	//   prices{1}(3) + funding{2,3}(4) + bridge{11}(3) + operations{4,5}(4)
	//   + others{6}(3) + others{7,8}(4) + additional{9}(3) + additional{10}(3) = 27.
	maxBytes := protoTxSize([]byte{1}) + protoTxSize([]byte{2, 3}) + protoTxSize([]byte{11}) +
		protoTxSize([]byte{4, 5}) + protoTxSize([]byte{6}) + protoTxSize([]byte{7, 8}) +
		protoTxSize([]byte{9}) + protoTxSize([]byte{10})

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
					MaxTxBytes: int64(maxBytes), //nolint:gosec
				},
			)
			require.NoError(t, err)
			require.Equal(t, maxBytes, ppt.MaxBytes)
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

// Test_UsedBytes_ProtoFraming asserts that budgeting is done against the CometBFT proto-framed size
// of a tx (1 tag byte + a varint-encoded length prefix + the tx payload), not its raw `len(tx)`, and
// that a tx (or batch of txs) is rejected once that framed total exceeds MaxBytes even when the raw
// byte length(s) alone would have fit.
func Test_UsedBytes_ProtoFraming(t *testing.T) {
	t.Run("proto size follows 1 + varint(len) + len, independent of cmttypes", func(t *testing.T) {
		for _, txLen := range []int{0, 1, 2, 100, 127, 128, 129, 16383, 16384} {
			tx := make([]byte, txLen)
			for i := range tx {
				tx[i] = byte(i)
			}

			want := 1 + uvarintProtoLen(uint64(txLen)) + uint64(txLen)
			require.Equalf(t, want, protoTxSize(tx), "txLen=%d", txLen)
		}
	})

	t.Run("a 1-byte tx consumes 3 bytes of budget", func(t *testing.T) {
		ppt, err := prepare.NewPrepareProposalTxs(&abci.RequestPrepareProposal{MaxTxBytes: 3})
		require.NoError(t, err)

		require.NoError(t, ppt.SetAcknowledgeBridgesTx([]byte{1}))
		require.Equal(t, uint64(3), ppt.UsedBytes) // 1 (tag) + 1 (varint(1)) + 1 (payload).
	})

	t.Run("a 128-byte tx consumes 1+2+128=131 bytes of budget", func(t *testing.T) {
		tx := make([]byte, 128)

		ppt, err := prepare.NewPrepareProposalTxs(&abci.RequestPrepareProposal{MaxTxBytes: 131})
		require.NoError(t, err)

		require.NoError(t, ppt.SetAcknowledgeBridgesTx(tx))
		// 1 (tag) + 2 (varint(128), which needs a 2nd byte since 128 >= 0x80) + 128 (payload).
		require.Equal(t, uint64(131), ppt.UsedBytes)
	})

	t.Run("SetAcknowledgeBridgesTx rejects a tx whose proto size exceeds MaxBytes though its raw "+
		"len fits", func(t *testing.T) {
		tx := make([]byte, 128) // raw len 128, proto size 1+2+128=131 (crosses the 1-byte varint boundary).
		protoSize := protoTxSize(tx)
		maxBytes := uint64(len(tx)) + 1 // strictly between the raw len (128) and the proto size (131).
		require.Greater(t, maxBytes, uint64(len(tx)))
		require.Less(t, maxBytes, protoSize)

		ppt, err := prepare.NewPrepareProposalTxs(
			&abci.RequestPrepareProposal{MaxTxBytes: int64(maxBytes)}, //nolint:gosec
		)
		require.NoError(t, err)

		err = ppt.SetAcknowledgeBridgesTx(tx)
		require.ErrorContains(
			t, err, fmt.Sprintf("Exceeds max: max=%d, used=0, adding=%d", maxBytes, protoSize),
		)
		require.Nil(t, ppt.AcknowledgeBridgesTx)
		require.Equal(t, uint64(0), ppt.UsedBytes)
	})

	t.Run("AddOtherTxs rejects a batch whose proto sum exceeds MaxBytes though the raw len sum fits",
		func(t *testing.T) {
			// Two 128-byte txs: raw sum 256, proto sum 2*(1+2+128)=262.
			txs := [][]byte{make([]byte, 128), make([]byte, 128)}
			rawSum := uint64(len(txs[0]) + len(txs[1]))
			protoSum := protoTxSize(txs[0]) + protoTxSize(txs[1])
			maxBytes := rawSum + 1 // strictly between the raw sum (256) and the proto sum (262).
			require.Greater(t, maxBytes, rawSum)
			require.Less(t, maxBytes, protoSum)

			ppt, err := prepare.NewPrepareProposalTxs(
				&abci.RequestPrepareProposal{MaxTxBytes: int64(maxBytes)}, //nolint:gosec
			)
			require.NoError(t, err)

			err = ppt.AddOtherTxs(txs)
			require.ErrorContains(
				t, err, fmt.Sprintf("Exceeds max: max=%d, used=0, adding=%d", maxBytes, protoSum),
			)
			require.Nil(t, ppt.OtherTxs)
			require.Equal(t, uint64(0), ppt.UsedBytes)
		})
}
