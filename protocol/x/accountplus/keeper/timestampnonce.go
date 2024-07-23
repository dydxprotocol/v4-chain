package keeper

const TimestampNonceSequenceCutoff uint64 = 1 << 40 // 2^40

func IsTimestampNonce(ts uint64) bool {
	return ts >= TimestampNonceSequenceCutoff
}
