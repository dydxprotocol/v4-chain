package types_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func TestDeepCopy(t *testing.T) {
	p := constants.Short_Asset_1ETH
	deepCopy := p.DeepCopy()

	require.Equal(t, p, deepCopy)
	require.NotSame(t, &p, &deepCopy)
}
