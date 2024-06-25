package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
)

func TestAssetPosition_DeepCopy(t *testing.T) {
	p := constants.Short_Asset_1ETH
	deepCopy := p.DeepCopy()

	require.Equal(t, p, deepCopy)
	require.NotSame(t, &p, &deepCopy)
}
