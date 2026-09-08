package types_test

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"
)

func TestDefaultGenesisIsDisabled(t *testing.T) {
	genesis := types.DefaultGenesis()
	require.Equal(t, types.Mode_MODE_DISABLED, genesis.Controls.Mode)
	require.NoError(t, genesis.Validate())
}

func TestControlsRequireFixedLogicalDenomAndBounds(t *testing.T) {
	controls := types.DefaultControls()
	controls.LogicalDenom = "uusdc"
	require.ErrorIs(t, controls.Validate(), types.ErrInvalidControls)

	controls = validControls()
	controls.MaxPendingSettlements = types.HardMaxPending + 1
	require.ErrorIs(t, controls.Validate(), types.ErrInvalidControls)
}

func TestFundingMemoIsVersionedAndStrict(t *testing.T) {
	memo := `{"canonical_usdc":{"version":1,"action":"fund_backing_swap",` +
		`"controller":"dydx1controller","noble_recipient":"noble1receiver"}}`
	parsed, found, err := types.ParseFundingMemo(memo)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, types.FundBackingSwapAction, parsed.Action)

	_, found, err = types.ParseFundingMemo(`{"forward":{"receiver":"other"}}`)
	require.NoError(t, err)
	require.False(t, found)

	_, found, err = types.ParseFundingMemo(
		`{"canonical_usdc":{"version":1,"action":"fund_backing_swap",` +
			`"controller":"x","noble_recipient":"y","extra":true}}`,
	)
	require.True(t, found)
	require.ErrorIs(t, err, types.ErrInvalidFundingMemo)
}

func TestExistingMsgTransferWireEncodingGolden(t *testing.T) {
	msg := &transfertypes.MsgTransfer{
		SourcePort:       "transfer",
		SourceChannel:    "channel-0",
		Token:            sdk.NewInt64Coin(types.DefaultControls().LogicalDenom, 42),
		Sender:           "dydx1sender",
		Receiver:         "inj1receiver",
		TimeoutTimestamp: 99,
		Memo:             "memo",
	}
	bz := transfertypes.ModuleCdc.MustMarshal(msg)
	const golden = "0a087472616e7366657212096368616e6e656c2d301a4a0a446962632f384532" +
		"374241324435343933414635363336373630453335344534363030343536324334364142374543" +
		"3043433443314341313445394532304532353435423512023432220b647964783173656e646572" +
		"2a0c696e6a3172656365697665723200386342046d656d6f"
	require.Equal(t,
		golden,
		hex.EncodeToString(bz),
	)
	require.Equal(t, "/ibc.applications.transfer.v1.MsgTransfer", sdk.MsgTypeURL(msg))
	var decoded transfertypes.MsgTransfer
	require.NoError(t, transfertypes.ModuleCdc.Unmarshal(bz, &decoded))
	require.Equal(t, *msg, decoded)
}

func TestUpdateParticipantsAcceptsOmittedAccountingFields(t *testing.T) {
	controller := sdk.AccAddress(make([]byte, 20)).String()
	msg := types.MsgUpdateParticipants{
		Authority: controller,
		Participants: []types.Participant{{
			Controller:     controller,
			NobleRecipient: "noble1receiver",
			MaxRelease:     "100",
		}},
	}
	require.NoError(t, msg.ValidateBasic())

	msg.Participants[0].Funded = "1"
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidParticipant)
}

func validControls() types.Controls {
	injectiveDenom := transfertypes.ParseDenomTrace("transfer/channel-1/uusdc").IBCDenom()
	return types.Controls{
		Mode:                    types.Mode_MODE_GRADUAL,
		LogicalDenom:            types.DefaultControls().LogicalDenom,
		NobleChannel:            "channel-0",
		InjectiveChannel:        "channel-1",
		NoblePacketDenom:        "uusdc",
		InjectivePacketDenom:    "uusdc",
		InjectiveDenom:          injectiveDenom,
		NobleWithdrawalsEnabled: true,
		MaxTransferAmount:       "1000",
		MigrationCeiling:        "10000",
		MaxPendingSettlements:   10,
		MemoVersion:             1,
		NobleClient:             "07-tendermint-0",
		NobleConnection:         "connection-0",
		InjectiveClient:         "07-tendermint-1",
		InjectiveConnection:     "connection-1",
	}
}
