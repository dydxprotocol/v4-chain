package types_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestModuleAddress(t *testing.T) {
	require.Equal(t, "klyra1v88c3xv9xyv3eetdx0tvcmq7ung3dywptd5ps3", types.ModuleAddress.String())
}
