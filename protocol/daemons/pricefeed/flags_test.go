package pricefeed_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dydxprotocol/v4/daemons/pricefeed"
	"github.com/dydxprotocol/v4/mocks"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	THIRTY_SECONDS_IN_NANOSEC = 30_000_000_000
)

func TestCheckMaxPriceAge(t *testing.T) {
	require.Equal(t, time.Duration(THIRTY_SECONDS_IN_NANOSEC), pricefeed.MaxPriceAge)
}

func TestAddSharedPriceFeedFlagsToCmd(t *testing.T) {
	cmd := cobra.Command{}

	pricefeed.AddSharedPriceFeedFlagsToCmd(&cmd)
	tests := map[string]struct {
		expectedFlagName string
	}{
		fmt.Sprintf("Has %s flag", pricefeed.FlagPriceFeedUnixSocketAddr): {
			expectedFlagName: "pricefeed-unixsocketaddress",
		},
		fmt.Sprintf("Has %s flag", pricefeed.FlagPriceFeedEnabled): {
			expectedFlagName: "pricefeed-enabled",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Contains(t, cmd.Flags().FlagUsages(), tc.expectedFlagName)
		})
	}
}

func AddClientPriceFeedFlagsToCmd(t *testing.T) {
	cmd := cobra.Command{}

	pricefeed.AddSharedPriceFeedFlagsToCmd(&cmd)
	tests := map[string]struct {
		expectedFlagName string
	}{
		fmt.Sprintf("Has %s flag", pricefeed.FlagPriceFeedPriceUpdaterLoopDelayMs): {
			expectedFlagName: "pricefeed-price-updater-loop-delay-ms",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Contains(t, cmd.Flags().FlagUsages(), tc.expectedFlagName)
		})
	}
}

func TestGetServerPricefeedFlagValuesFromOptions(t *testing.T) {
	tests := map[string]struct {
		// parameters
		pricefeedUnixSocketAddressOpt string
		pricefeedEnabled              bool

		// expectations
		expectedPricefeedUnixSocketAddressOpt string
		expectedPricefeedEnabled              bool
	}{
		"Sets values from options": {
			pricefeedUnixSocketAddressOpt:         "/special-pricefeed.sock",
			pricefeedEnabled:                      false,
			expectedPricefeedUnixSocketAddressOpt: "/special-pricefeed.sock",
			expectedPricefeedEnabled:              false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			optsMap := make(map[string]interface{})
			optsMap[pricefeed.FlagPriceFeedUnixSocketAddr] = tc.pricefeedUnixSocketAddressOpt
			optsMap[pricefeed.FlagPriceFeedEnabled] = tc.pricefeedEnabled
			mockOpts := mocks.AppOptions{}
			mockOpts.On("Get", mock.Anything).
				Return(func(key string) interface{} {
					return optsMap[key]
				})

			pricefeedEnabled,
				pricefeedUnixSocketAddress := pricefeed.GetServerPricefeedFlagValuesFromOptions(
				&mockOpts,
			)

			require.Equal(t, tc.expectedPricefeedUnixSocketAddressOpt, pricefeedUnixSocketAddress)
			require.Equal(t, tc.expectedPricefeedEnabled, pricefeedEnabled)
		})
	}
}

func TestGetClientPricefeedFlagValuesFromOptions(t *testing.T) {
	tests := map[string]struct {
		// parameters
		pricefeedUnixSocketAddressOpt    string
		pricefeedEnabled                 bool
		pricefeedPriceUpdaterLoopDelayMs int

		// expectations
		expectedPricefeedUnixSocketAddressOpt    string
		expectedPricefeedEnabled                 bool
		expectedPricefeedPriceUpdaterLoopDelayMs uint32
	}{
		"Sets values from options": {
			pricefeedUnixSocketAddressOpt:            "/special-pricefeed.sock",
			pricefeedEnabled:                         false,
			pricefeedPriceUpdaterLoopDelayMs:         100,
			expectedPricefeedUnixSocketAddressOpt:    "/special-pricefeed.sock",
			expectedPricefeedEnabled:                 false,
			expectedPricefeedPriceUpdaterLoopDelayMs: uint32(100),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			optsMap := make(map[string]interface{})
			optsMap[pricefeed.FlagPriceFeedUnixSocketAddr] = tc.pricefeedUnixSocketAddressOpt
			optsMap[pricefeed.FlagPriceFeedEnabled] = tc.pricefeedEnabled
			optsMap[pricefeed.FlagPriceFeedPriceUpdaterLoopDelayMs] = tc.pricefeedPriceUpdaterLoopDelayMs
			mockOpts := mocks.AppOptions{}
			mockOpts.On("Get", mock.Anything).
				Return(func(key string) interface{} {
					return optsMap[key]
				})

			pricefeedEnabled,
				pricefeedUnixSocketAddress,
				pricefeedPriceUpdaterLoopDelayMs := pricefeed.GetClientPricefeedFlagValuesFromOptions(
				&mockOpts,
			)

			require.Equal(t, tc.expectedPricefeedUnixSocketAddressOpt, pricefeedUnixSocketAddress)
			require.Equal(t, tc.expectedPricefeedEnabled, pricefeedEnabled)
			require.Equal(t, tc.expectedPricefeedPriceUpdaterLoopDelayMs, pricefeedPriceUpdaterLoopDelayMs)
		})
	}
}
