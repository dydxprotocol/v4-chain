package types_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestInsuranceFundModuleAddress(t *testing.T) {
	require.Equal(t, "klyra1c7ptc87hkd54e3r7zjy92q29xkq7t79w9y9stt", types.InsuranceFundModuleAddress.String())
}
