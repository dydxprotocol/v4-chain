package lib_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/stretchr/testify/require"
)

func TestGovModuleAddress(t *testing.T) {
	require.Equal(t, "klyra10d07y265gmmuvt4z0w9aw880jnsr700jv2gw70", lib.GovModuleAddress.String())
}
