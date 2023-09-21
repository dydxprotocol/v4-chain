package module_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/app/module"
	"github.com/stretchr/testify/require"
)

func TestDefaultGenesis(t *testing.T) {
	encodingConfig := app.GetEncodingConfig()
	defaultGenesis := module.SlashingModuleBasic{}.DefaultGenesis(encodingConfig.Codec)
	humanReadableDefaultGenesisState, unmarshalErr := json.Marshal(&defaultGenesis)

	expectedDefaultGenesisState, fileReadErr := os.ReadFile("testdata/slashing_default_genesis_state.json")

	require.NoError(t, unmarshalErr)
	require.NoError(t, fileReadErr)
	require.JSONEq(t,
		string(expectedDefaultGenesisState), string(humanReadableDefaultGenesisState))
}
