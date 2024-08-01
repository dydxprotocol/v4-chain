package app_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
)

func TestDefaultGenesisState(t *testing.T) {
	app := testapp.DefaultTestApp(nil)
	defaultGenesisState := app.DefaultGenesis()
	humanReadableDefaultGenesisState, jsonUnmarshalErr := json.Marshal(&defaultGenesisState)

	expectedDefaultGenesisState, fileReadErr := os.ReadFile("testdata/default_genesis_state.json")

	require.NoError(t, fileReadErr)
	require.NoError(t, jsonUnmarshalErr)
	require.JSONEq(t, string(expectedDefaultGenesisState), string(humanReadableDefaultGenesisState))
}
