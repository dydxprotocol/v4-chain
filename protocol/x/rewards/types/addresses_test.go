package types_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestTreasuryModuleAddress(t *testing.T) {
	require.Equal(t, "klyra16wrau2x4tsg033xfrrdpae6kxfn9kyueujpa62", types.TreasuryModuleAddress.String())
}
