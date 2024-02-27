package app

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/dydxprotocol/v4-chain/protocol/app/constants"
	"github.com/dydxprotocol/v4-chain/protocol/app/flags"
	errorspkg "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func configureDatadogProfilerOptions(
	logger log.Logger,
	ddAgentHost string,
	ddAgentPort uint16,
) (
	ddService string,
	ddVersion string,
	ddAgentHostPort string,
	err error,
) {
	// Use a default application name unless overridden by the DD_SERVICE environment variable.
	ddService = constants.ServiceName
	if found := os.Getenv("DD_SERVICE"); found != "" {
		logger.Info(fmt.Sprintf("DD_SERVICE defined, overriding default of '%s'.", constants.ServiceName))
		ddService = found
	}

	// Use the application version during build time unless overridden by the DD_VERSION environment variable.
	ddVersion = version.Version
	if found := os.Getenv("DD_VERSION"); found != "" {
		if ddVersion != "" {
			logger.Info(fmt.Sprintf("DD_VERSION defined, overriding build time version of '%s'.",
				version.Version))
		}
		ddVersion = found
	}
	if ddVersion == "" {
		logger.Info("Not initializing Datadog profiler. Application version was not defined during " +
			"build time and no override has been provided by environment variable DD_VERSION.")
		return "", "", "", errors.New("application version was not defined")
	}

	// Validates that the environment has been populated to ensure that profiles are grouped appropriately.
	ddEnv := os.Getenv("DD_ENV")
	if ddEnv == "" {
		logger.Info("Not initializing Datadog profiler. Application environment has not been " +
			"provided by the DD_ENV environment variable.")
		return "", "", "", errors.New("environment is not defined")
	}

	// Use the command line flag passed in during runtime unless overridden the by the DD_AGENT_HOST
	// environment variable.
	if found := os.Getenv("DD_AGENT_HOST"); found != "" {
		if ddAgentHost != flags.DefaultDdAgentHost {
			logger.Info(fmt.Sprintf("DD_AGENT_HOST defined, overriding --%s flag value of '%s'.",
				flags.DdAgentHost, ddAgentHost))
		}
		ddAgentHost = found
	}
	if ddAgentHost == "" {
		logger.Info(fmt.Sprintf("Not initializing Datadog profiler. Datadog agent host was not specified "+
			"either via flag --%s or environment variable DD_AGENT_HOST.", flags.DdAgentHost))
		return "", "", "", errors.New("datadog agent host was not specified")
	}

	// Override the port specified on the command line with any provided by the DD_TRACE_AGENT_PORT environment
	// variable.
	ddAgentHostPort = net.JoinHostPort(ddAgentHost, strconv.Itoa(int(ddAgentPort)))
	if found := os.Getenv("DD_TRACE_AGENT_PORT"); found != "" {
		if ddAgentPort != flags.DefaultDdTraceAgentPort {
			logger.Info(fmt.Sprintf("DD_TRACE_AGENT_PORT defined, overriding --%s flag value of '%d'.",
				flags.DdTraceAgentPort, ddAgentPort))
		}
		ddAgentHostPort = net.JoinHostPort(ddAgentHost, found)
	}
	return
}

// initDatadogProfiler initializes datadog continuous profiling.
//
// The profiles will be configured as follows:
//   - service: "validator" unless overridden by the DD_SERVICE environment variable.
//   - version: version specified during go build unless overridden by the DD_VERSION environment variable.
//   - env: must be specified by the DD_ENV environment variable.
//
// The profiles will be sent to the Datadog agent based upon the --dd-agent-host and --dd-trace-agent-port command
// line flags unless overridden by the DD_AGENT_HOST and DD_TRACE_AGENT_PORT environment variables respectively.
//
// See https://docs.datadoghq.com/profiler/enabling/go/ for more details.
func initDatadogProfiler(logger log.Logger, ddAgentHost string, ddAgentPort uint16) {
	ddService, ddVersion, ddAgentHostPort, err := configureDatadogProfilerOptions(logger, ddAgentHost, ddAgentPort)
	if err != nil {
		return
	}
	err = profiler.Start(
		profiler.WithService(ddService),
		profiler.WithVersion(ddVersion),
		profiler.WithAgentAddr(ddAgentHostPort),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
			profiler.MutexProfile,
		),
	)
	if err != nil {
		panic(err)
	}
}

type DatadogErrorTrackingObject struct {
	Stack   []map[string]string
	Message string
	Kind    string
}

func (obj DatadogErrorTrackingObject) MarshalZerologObject(e *zerolog.Event) {
	e.Interface("stack", obj.Stack).
		Str("message", obj.Message).
		Str("kind", obj.Kind)
}

var (
	zerologFormatterOnce sync.Once
)

// SetZerologDatadogErrorTrackingFormat sets custom error formatting for log tag
// values that are errors for the zerolog library. Converts them to a format that
// is compatible with datadog error tracking.
func SetZerologDatadogErrorTrackingFormat() {
	zerologFormatterOnce.Do(func() {
		// Error fields are default set under `error`
		// Extract + add the kind and message field
		zerolog.ErrorMarshalFunc = func(err error) interface{} {
			stackArr, ok := pkgerrors.MarshalStack(errorspkg.WithStack(err)).([]map[string]string)
			if !ok {
				return struct{}{}
			}
			objectToReturn := DatadogErrorTrackingObject{
				// Discard common stack prefixes
				// TODO(CLOB-1049) Write test for common stack prefix truncation
				Stack:   stackArr[5:],
				Kind:    "Exception",
				Message: err.Error(),
			}
			return objectToReturn
		}
	})
}
