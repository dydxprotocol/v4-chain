package app

import (
	"github.com/dydxprotocol/v4-chain/protocol/testutil/logger"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestConfigureDatadogProfilerOptions(t *testing.T) {
	type Result struct {
		ddService       string
		ddVersion       string
		ddAgentHostPort string
	}

	tests := map[string]struct {
		agentHost   string
		agentPort   uint16
		envVars     map[string]string
		expectError bool
		result      Result
	}{
		"Success: default options": {
			agentHost: "test_host",
			agentPort: uint16(9999),
			envVars: map[string]string{
				"DD_VERSION": "v1",
				"DD_ENV":     "staging",
			},
			result: Result{
				ddService:       "validator",
				ddAgentHostPort: "test_host:9999",
				ddVersion:       "v1",
			},
		},
		"Success: environment agent host overrides parameter": {
			agentHost: "test_host",
			agentPort: uint16(9999),
			envVars: map[string]string{
				"DD_VERSION":    "v1",
				"DD_ENV":        "staging",
				"DD_AGENT_HOST": "alternative_host",
			},
			result: Result{
				ddService:       "validator",
				ddAgentHostPort: "alternative_host:9999",
				ddVersion:       "v1",
			},
		},
		"Success: environment trace agent port overrides parameter": {
			agentHost: "test_host",
			agentPort: uint16(9999),
			envVars: map[string]string{
				"DD_VERSION":          "v1",
				"DD_ENV":              "staging",
				"DD_TRACE_AGENT_PORT": "8888",
			},
			result: Result{
				ddService:       "validator",
				ddAgentHostPort: "test_host:8888",
				ddVersion:       "v1",
			},
		},
		"Success: service environment variable overrides default": {
			agentHost: "test_host",
			agentPort: uint16(9999),
			envVars: map[string]string{
				"DD_SERVICE": "test_service",
				"DD_VERSION": "v1",
				"DD_ENV":     "staging",
			},
			result: Result{
				ddService:       "test_service",
				ddAgentHostPort: "test_host:9999",
				ddVersion:       "v1",
			},
		},
		"Failure: missing version": {
			agentHost: "test_host",
			agentPort: uint16(9999),
			envVars: map[string]string{
				"DD_ENV": "staging",
			},
			expectError: true,
		},
		"Failure: missing datadog environment": {
			agentHost: "test_host",
			agentPort: uint16(9999),
			envVars: map[string]string{
				"DD_VERSION": "v1",
			},
			expectError: true,
		},
		"Failure: agent host unspecified": {
			agentHost: "",
			agentPort: uint16(9999),
			envVars: map[string]string{
				"DD_VERSION": "v1",
				"DD_ENV":     "staging",
			},
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			logger, _ := logger.TestLogger()

			// Optional configure environment variables
			if len(tc.envVars) > 0 {
				for key, value := range tc.envVars {
					os.Setenv(key, value)
				}
				defer func() {
					for key := range tc.envVars {
						os.Unsetenv(key)
					}
				}()
			}

			ddService, ddVersion, ddAgentHostPort, err := configureDatadogProfilerOptions(logger, tc.agentHost, tc.agentPort)
			if tc.expectError {
				require.Zero(t, ddService)
				require.Zero(t, ddVersion)
				require.Zero(t, ddAgentHostPort)
				require.NotNil(t, err)
			} else {
				require.Equal(t, tc.result.ddService, ddService)
				require.Equal(t, tc.result.ddVersion, ddVersion)
				require.Equal(t, tc.result.ddAgentHostPort, ddAgentHostPort)
				require.Nil(t, err)
			}
			require.NotPanics(t, func() { initDatadogProfiler(logger, tc.agentHost, tc.agentPort) })
		})
	}
}
