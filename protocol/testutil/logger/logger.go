package logger

import (
	"bytes"
	"github.com/cometbft/cometbft/libs/log"
)

const (
	LoggerInstanceForTest = "logger-instance-for-test"
)

// TestLogger returns a logger instance and a buffer where all logs are written to.
func TestLogger() (log.Logger, *bytes.Buffer) {
	var logBuffer bytes.Buffer
	return log.NewTMLogger(log.NewSyncWriter(&logBuffer)), &logBuffer
}
