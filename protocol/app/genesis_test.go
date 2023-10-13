package app_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/stretchr/testify/require"
)

func TestDefaultGenesisState(t *testing.T) {
	encodingConfig := app.GetEncodingConfig()
	defaultGenesisState := app.NewDefaultGenesisState(encodingConfig.Codec)
	humanReadableDefaultGenesisState, jsonUnmarshalErr := json.Marshal(&defaultGenesisState)

	expectedDefaultGenesisState, fileReadErr := os.ReadFile("testdata/default_genesis_state.json")

	require.NoError(t, fileReadErr)
	require.NoError(t, jsonUnmarshalErr)
	require.JSONEq(t, string(expectedDefaultGenesisState), string(humanReadableDefaultGenesisState))
}
