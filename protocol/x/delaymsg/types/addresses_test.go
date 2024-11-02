package types_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
	"github.com/stretchr/testify/require"
)

func TestModuleAddress(t *testing.T) {
	require.Equal(t, "klyra1mkkvp26dngu6n8rmalaxyp3gwkjuzztqtnn4rg", types.ModuleAddress.String())
}
