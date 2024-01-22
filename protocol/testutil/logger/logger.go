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
func TestLogger() (log.Logger, *bytes.Buffer) {
	var logBuffer bytes.Buffer
	return log.NewLogger(kitlog.NewSyncWriter(&logBuffer)), &logBuffer
}
