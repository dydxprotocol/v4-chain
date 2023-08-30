package msgsender

import (
	"errors"
)

var ErrKafkaAlreadyClosed = errors.New("IndexerMessageSenderKafka is already closed")
