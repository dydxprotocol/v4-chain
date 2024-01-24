package lib

import (
	"context"
	"fmt"

	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
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
	return ctx
}
