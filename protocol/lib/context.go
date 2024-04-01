package lib

import (
	"fmt"

	"github.com/cometbft/cometbft/crypto/tmhash"
)

// Custom exec modes
const (
	ExecModeBeginBlock        = uint32(100)
	ExecModeEndBlock          = uint32(101)
	ExecModePrepareCheckState = uint32(102)
)

type TxHash string

func GetTxHash(tx []byte) TxHash {
	return TxHash(fmt.Sprintf("%X", tmhash.Sum(tx)))
}
