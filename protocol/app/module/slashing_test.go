package module_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/dydxprotocol/v4/app/basic_manager"
	"github.com/dydxprotocol/v4/app/module"
	"github.com/dydxprotocol/v4/lib/encoding"
	"github.com/dydxprotocol/v4/testutil/stringutils"
	"github.com/stretchr/testify/require"
)

func TestDefaultGenesis(t *testing.T) {
	encodingConfig := encoding.MakeEncodingConfig(basic_manager.ModuleBasics)
	defaultGenesis := module.SlashingModuleBasic{}.DefaultGenesis(encodingConfig.Codec)
	humanReadableDefaultGenesisState, unmarshalErr := json.Marshal(&defaultGenesis)

	expectedDefaultGenesisState, fileReadErr := os.ReadFile("testdata/slashing_default_genesis_state.json")

	require.NoError(t, unmarshalErr)
	require.NoError(t, fileReadErr)
	require.Equal(t,
		stringutils.StripSpaces(expectedDefaultGenesisState), string(humanReadableDefaultGenesisState))
}
