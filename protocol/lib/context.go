package lib

import (
	"fmt"
	"github.com/cometbft/cometbft/crypto/tmhash"
)

// Custom exec modes
const (
	ExecModeBeginBlock        = 100
	ExecModeEndBlock          = 101
	ExecModePrepareCheckState = 102
)

type TxHash string

func GetTxHash(tx []byte) TxHash {
	return TxHash(fmt.Sprintf("%X", tmhash.Sum(tx)))
}
