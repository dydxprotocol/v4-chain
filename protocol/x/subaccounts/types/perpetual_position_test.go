package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
)

func TestPerpetualPosition_DeepCopy(t *testing.T) {
	p := constants.PerpetualPosition_OneISO2Short
	deepCopy := p.DeepCopy()

	require.Equal(t, p, deepCopy)
	require.NotSame(t, &p, &deepCopy)
}
