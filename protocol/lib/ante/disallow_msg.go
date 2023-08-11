package ante

import sdk "github.com/cosmos/cosmos-sdk/types"

// IsDisallowExternalSubmitMsg returns true if the msg is not allowed to be submitted externally.
func IsDisallowExternalSubmitMsg(msg sdk.Msg) bool {
	if IsAppInjectedMsg(msg) || IsInternalMsg(msg) || IsUnsupportedMsg(msg) {
		return true
	}
	if IsNestedMsg(msg) {
		if err := ValidateNestedMsg(msg); err != nil {
			return true
		}
	}
	return false
}
