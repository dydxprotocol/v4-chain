package logger

import (
	"bytes"

	"cosmossdk.io/log"
	kitlog "github.com/go-kit/log"
)

const (
	LoggerInstanceForTest = "logger-instance-for-test"
)

// TestLogger returns a logger instance and a buffer where all logs are written to.
// TODO(CORE-538): See if we can get rid of this method in favor of the Cosmos Logger now
// (which uses zerolog under the hood).
func TestLogger() (log.Logger, *bytes.Buffer) {
	var logBuffer bytes.Buffer
	return log.NewLogger(kitlog.NewSyncWriter(&logBuffer)), &logBuffer
}
