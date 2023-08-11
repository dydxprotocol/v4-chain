package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAddPremiumVotes{}

func NewFundingPremium(id uint32, premiumPpm int32) *FundingPremium {
	return &FundingPremium{
		PerpetualId: id,
		PremiumPpm:  premiumPpm,
	}
}

func NewMsgAddPremiumVotes(
	proposer string,
	votes []FundingPremium,
) *MsgAddPremiumVotes {
	return &MsgAddPremiumVotes{
		Proposer: proposer,
		Votes:    votes,
	}
}

func (msg *MsgAddPremiumVotes) GetSigners() []sdk.AccAddress {
	proposer, err := sdk.AccAddressFromBech32(msg.Proposer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{proposer}
}

func (msg *MsgAddPremiumVotes) ValidateBasic() error {
	for i, sample := range msg.Votes {
		if i > 0 && msg.Votes[i-1].PerpetualId >= sample.PerpetualId {
			return sdkerrors.Wrap(
				ErrInvalidAddPremiumVotes,
				"premium votes must be sorted by perpetual id in ascending order and cannot contain duplicates",
			)
		}
	}
	return nil
}
