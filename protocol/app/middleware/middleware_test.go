package middleware_test

import (
	"bytes"
	"fmt"
	"github.com/cometbft/cometbft/libs/log"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/app/middleware"
	"github.com/stretchr/testify/require"
)

func TestRunTxPanicLoggingMiddleware(t *testing.T) {
	tests := map[string]struct {
		function     func()
		expectedLogs []string
	}{
		"no panic": {
			function: func() {
				// Do something that does not panic.
			},
		},
		"panic with string": {
			function: func() {
				panic("test123")
			},
			expectedLogs: []string{
				"E[202",                       // error and date prefix
				"runTx panic'ed with test123", // message
				"middleware_test.go",          // part of stack trace
			},
		},
		"panic with error": {
			function: func() {
				panic(fmt.Errorf("test456"))
			},
			expectedLogs: []string{
				"E[202",                       // error and date prefix
				"runTx panic'ed with test456", // message
				"middleware_test.go",          // part of stack trace
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Restore the old logger after the test runs since middleware.Logger is a global variable.
			oldLogger := middleware.Logger
			defer func() { middleware.Logger = oldLogger }()

			buf := new(bytes.Buffer)
			middleware.Logger = log.NewTMLogger(buf)

			func() {
				defer func() {
					if r := recover(); r != nil {
						handler := middleware.NewRunTxPanicLoggingMiddleware()
						err := handler(r)
						require.Nil(t, err)
					}
				}()
				tc.function()
			}()

			if tc.expectedLogs == nil {
				require.Empty(t, buf.String())
			}
			for _, expectedLog := range tc.expectedLogs {
				require.Contains(t, buf.String(), expectedLog)
			}
		})
	}
}
