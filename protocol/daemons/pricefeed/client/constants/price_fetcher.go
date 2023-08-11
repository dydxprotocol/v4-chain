package constants

const (
	// 5K is chosen to be >> than the number of messages an exchange could send in any period before the
	// price encoder is able to read the messages from the buffer, even if we add O(10-100) markets dynamically,
	// but not large enough to allow more than at most a few minutes of price messages to accumulate.
	FixedBufferSize = 1024 * 5
	// https://stackoverflow.com/questions/37774624/go-http-get-concurrency-and-connection-reset-by-peer.
	// This is a good number to start with based on the above link. Adjustments can/will be made accordingly.
	MaxConnectionsPerExchange = 50
)
