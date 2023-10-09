package cmd_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/cmd/dydxprotocold/cmd"
	"github.com/stretchr/testify/require"
)

func TestMinGasPrice(t *testing.T) {
	require.Equal(t,
		"0.025ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5,25000000000adv4tnt",
		cmd.MinGasPrice,
	)
}
