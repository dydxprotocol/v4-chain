package module_test

import (
	"encoding/json"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"os"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/app/module"
	"github.com/stretchr/testify/require"
)

func TestDefaultGenesis(t *testing.T) {
	dydxApp := testapp.DefaultTestApp(nil)
	defaultGenesis := module.SlashingModuleBasic{}.DefaultGenesis(dydxApp.AppCodec())
	humanReadableDefaultGenesisState, unmarshalErr := json.Marshal(&defaultGenesis)

	expectedDefaultGenesisState, fileReadErr := os.ReadFile("testdata/slashing_default_genesis_state.json")

	require.NoError(t, unmarshalErr)
	require.NoError(t, fileReadErr)
	require.JSONEq(t,
		string(expectedDefaultGenesisState), string(humanReadableDefaultGenesisState))
}
