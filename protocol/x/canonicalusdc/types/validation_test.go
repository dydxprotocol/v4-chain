package types_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
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

	controls = validControls()
	controls.MaxTransferAmount = dtypes.SerializableInt{}
	require.ErrorIs(t, controls.Validate(), types.ErrInvalidControls)

	controls = validControls()
	controls.MigrationCeiling = dtypes.NewInt(999)
	require.ErrorIs(t, controls.Validate(), types.ErrInvalidControls)

	require.NoError(t, validControls().Validate())
}

func TestAmountOf(t *testing.T) {
	zero, err := types.AmountOf(dtypes.SerializableInt{})
	require.NoError(t, err)
	require.True(t, zero.IsZero())

	amount, err := types.AmountOf(dtypes.NewInt(42))
	require.NoError(t, err)
	require.Equal(t, int64(42), amount.Int64())

	_, err = types.AmountOf(dtypes.NewInt(-1))
	require.ErrorIs(t, err, types.ErrInvalidAmount)

	tooWide := new(big.Int).Lsh(big.NewInt(1), math.MaxBitLen)
	_, err = types.AmountOf(dtypes.NewIntFromBigInt(tooWide))
	require.ErrorIs(t, err, types.ErrInvalidAmount)

	_, err = types.PositiveAmountOf(dtypes.ZeroInt())
	require.ErrorIs(t, err, types.ErrInvalidAmount)

	require.Equal(t, "42", types.NewAmount(math.NewInt(42)).String())
}

func TestControlsTracesMatchOutboundPacketDenoms(t *testing.T) {
	controls := validControls()
	require.Equal(t, "transfer/channel-0/uusdc", controls.NobleTrace())
	require.Equal(t, "transfer/channel-1/uusdc", controls.InjectiveTrace())
	require.Equal(t, controls.LogicalDenom, transfertypes.ParseDenomTrace(controls.NobleTrace()).IBCDenom())
	require.Equal(t, controls.InjectiveDenom, transfertypes.ParseDenomTrace(controls.InjectiveTrace()).IBCDenom())
}

func TestControlsSameRoute(t *testing.T) {
	base := validControls()
	require.True(t, base.SameRoute(base))

	limitsOnly := base
	limitsOnly.MaxTransferAmount = dtypes.NewInt(1)
	limitsOnly.NobleWithdrawalsEnabled = false
	limitsOnly.Mode = types.Mode_MODE_PAUSED
	require.True(t, base.SameRoute(limitsOnly))

	mutations := []func(*types.Controls){
		func(c *types.Controls) { c.NobleChannel = "channel-9" },
		func(c *types.Controls) { c.NobleClient = "07-tendermint-9" },
		func(c *types.Controls) { c.NobleConnection = "connection-9" },
		func(c *types.Controls) { c.NoblePacketDenom = "uother" },
		func(c *types.Controls) { c.InjectiveChannel = "channel-9" },
		func(c *types.Controls) { c.InjectiveClient = "07-tendermint-9" },
		func(c *types.Controls) { c.InjectiveConnection = "connection-9" },
		func(c *types.Controls) { c.InjectivePacketDenom = "uother" },
		func(c *types.Controls) { c.InjectiveDenom = "ibc/OTHER" },
	}
	for _, mutate := range mutations {
		changed := base
		mutate(&changed)
		require.False(t, base.SameRoute(changed))
	}
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
		MaxTransferAmount:       dtypes.NewInt(1000),
		MigrationCeiling:        dtypes.NewInt(10000),
		MaxPendingSettlements:   10,
		NobleClient:             "07-tendermint-0",
		NobleConnection:         "connection-0",
		InjectiveClient:         "07-tendermint-1",
		InjectiveConnection:     "connection-1",
	}
}
