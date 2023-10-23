package types_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestModuleAddress(t *testing.T) {
	require.Equal(t, "dydx1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5upwc6", types.ModuleAddress.String())
}
