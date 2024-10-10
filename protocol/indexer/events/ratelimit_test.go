package events

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUpdateYieldParamsEventV1_Success(t *testing.T) {
	updatePerpetualEventV1 := NewUpdateYieldParamsEventV1(
		"100000000",
		"1/1",
	)
	expectedUpdatePerpetualEventV1Proto := &UpdateYieldParamsEventV1{
		SdaiPrice:       "100000000",
		AssetYieldIndex: "1/1",
	}
	require.Equal(t, expectedUpdatePerpetualEventV1Proto, updatePerpetualEventV1)
}
