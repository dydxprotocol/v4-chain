package lib

import (
	"fmt"
	"github.com/cometbft/cometbft/crypto/tmhash"
)

type TxHash string

func GetTxHash(tx []byte) TxHash {
	return TxHash(fmt.Sprintf("%X", tmhash.Sum(tx)))
}
