package lib_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGovModuleAddress(t *testing.T) {
	require.Equal(t, "dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky", lib.GovModuleAddress.String())
}
