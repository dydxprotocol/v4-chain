package app_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/dydxprotocol/v4/app"
	"github.com/dydxprotocol/v4/app/basic_manager"
	"github.com/dydxprotocol/v4/lib/encoding"
	"github.com/stretchr/testify/require"
)

func TestDefaultGenesisState(t *testing.T) {
	encodingConfig := encoding.MakeEncodingConfig(basic_manager.ModuleBasics)
	defaultGenesisState := app.NewDefaultGenesisState(encodingConfig.Codec)
	humanReadableDefaultGenesisState, jsonUnmarshalErr := json.Marshal(&defaultGenesisState)

	expectedDefaultGenesisState, fileReadErr := os.ReadFile("testdata/default_genesis_state.json")

	require.NoError(t, fileReadErr)
	require.NoError(t, jsonUnmarshalErr)
	require.JSONEq(t, string(expectedDefaultGenesisState), string(humanReadableDefaultGenesisState))
}
