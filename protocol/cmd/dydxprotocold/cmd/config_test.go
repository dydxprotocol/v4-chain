package cmd_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/cmd/dydxprotocold/cmd"
	"github.com/stretchr/testify/require"
)

func TestMinGasPrice(t *testing.T) {
	require.Equal(t,
		"0.025utdai,25000000000adv4tnt",
		cmd.MinGasPrice,
	)
}
