package streaming

// Constants for FullNodeStreamingManager.
const (
	// Transient store key for storing staged events.
	StreamingManagerTransientStoreKey = "tmp_streaming"

	// Key for storing the count of staged events.
	StagedEventsCountKey = "EvtCnt"

	// Key prefix for staged events.
	StagedEventsKeyPrefix = "Evt:"
)
