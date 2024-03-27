package lib

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
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

// UnwrapSDKContext is a thin wrapper around cosmos sdk's unwrap function
// that extracts the cosmos context from the standard golang context.
// If moduleName is provided, it appends the persistent log tag with
// the module name to the logger in the context.
func UnwrapSDKContext(
	goCtx context.Context,
	moduleName string,
) sdk.Context {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if moduleName != "" {
		ctx = log.AddPersistentTagsToLogger(
			ctx,
			log.Module,
			fmt.Sprintf("x/%s", moduleName),
		)
	}
	// Generate a 20-length random hex string for request id.
	bytes := make([]byte, 10)
	_, err := rand.Read(bytes)
	if err != nil {
		return ctx
	}
	requestId := hex.EncodeToString(bytes)
	ctx = log.AddPersistentTagsToLogger(
		ctx,
		log.RequestId, requestId,
	)
	return ctx
}
