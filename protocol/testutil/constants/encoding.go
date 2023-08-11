package constants

import "github.com/dydxprotocol/v4/testutil/encoding"

var (
	TestEncodingCfg = encoding.GetTestEncodingCfg()
	TestTxBuilder   = TestEncodingCfg.TxConfig.NewTxBuilder()
)
