package lib

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgRouter interface {
	Handler(msg sdk.Msg) baseapp.MsgServiceHandler
}
