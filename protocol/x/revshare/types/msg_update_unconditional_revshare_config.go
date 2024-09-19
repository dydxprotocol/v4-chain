package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

var _ types.Msg = &MsgUpdateUnconditionalRevShareConfig{}

// ValidateBasic performs validation to check the total percentage
// across all configs is <= 100
func (msg *MsgUpdateUnconditionalRevShareConfig) ValidateBasic() error {
	var totalRevsharePercentagePpm uint32 = 0
	config := msg.Config
	for _, config := range config.Configs {
		totalRevsharePercentagePpm += config.SharePpm
	}

	if totalRevsharePercentagePpm >= lib.OneMillion {
		return errorsmod.Wrapf(
			ErrInvalidRevShareConfig,
			"total revshare percentage ppm %d is greater than 100%%",
			totalRevsharePercentagePpm,
		)
	}
	return nil
}
