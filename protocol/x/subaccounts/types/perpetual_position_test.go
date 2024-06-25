package types_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func TestDeepCopy(t *testing.T) {
	p := constants.PerpetualPosition_OneISO2Short
	deepCopy := p.DeepCopy()

	require.Equal(t, p, deepCopy)
	require.NotSame(t, &p, &deepCopy)
}
