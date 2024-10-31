package testutil

import (
	"testing"

	sdaioracletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/client/types"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

// SetupExchangeResponses validates and sets up responses returned by exchange APIs using `gock`.
func SetupSDaiResponse(
	t *testing.T,
) {
	rootUrl := sdaioracletypes.ETHRPC
	response := gock.New(rootUrl).
		Post("/").
		MatchType("json").
		Persist().
		Reply(200).
		JSON(map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"result":  "0x0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000",
		})

	require.NotNil(t, response)
}
