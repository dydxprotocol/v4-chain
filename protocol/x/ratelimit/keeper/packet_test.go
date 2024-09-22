package keeper_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdaidaemontypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	delaymsgmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	packetDataUnparsableSender = ibctransfertypes.FungibleTokenPacketData{
		Sender:   "",
		Receiver: testAddress2,
		Denom:    "gsdai",
		Amount:   "100000000",
	}

	packetDataUnparsableReceiver = ibctransfertypes.FungibleTokenPacketData{
		Sender:   testAddress1,
		Receiver: "",
		Denom:    "gsdai",
		Amount:   "100000000",
	}

	packetDataUnparsableAmount = ibctransfertypes.FungibleTokenPacketData{
		Sender:   testAddress1,
		Receiver: testAddress2,
		Denom:    "gsdai",
		Amount:   "1abc2",
	}

	examplePacketDataNoSourcePort = ibctransfertypes.FungibleTokenPacketData{
		Sender:   testAddress1,
		Receiver: testAddress2,
		Denom:    "gsdai",
		Amount:   "100000000",
	}

	examplePacketDataForSourcePortPacket = ibctransfertypes.FungibleTokenPacketData{
		Sender:   testAddress1,
		Receiver: testAddress2,
		Denom:    "gsdai",
		Amount:   "100000000000000",
	}

	examplePacketDataNonSDai = ibctransfertypes.FungibleTokenPacketData{
		Sender:   testAddress1,
		Receiver: testAddress2,
		Denom:    "pdai",
		Amount:   "100000000000000",
	}

	fullSDaiPathPacketData = ibctransfertypes.FungibleTokenPacketData{
		Sender:   testAddress1,
		Receiver: testAddress2,
		Denom:    types.SDaiBaseDenomFullPath,
		Amount:   "100000000000000",
	}

	marshaledExamplePacketDataNoSourcePort = marshalPacketData(examplePacketDataNoSourcePort)

	marshaledexamplePacketDataForSourcePortPacket = marshalPacketData(examplePacketDataForSourcePortPacket)

	marshaledExamplePacketDataNonSDai = marshalPacketData(examplePacketDataNonSDai)

	marshaledExampleFullSDaiPathPacketData = marshalPacketData(fullSDaiPathPacketData)

	examplePacketNoSourcePort = channeltypes.Packet{
		SourceChannel: "channel-0",
		Sequence:      1,
		Data:          marshaledExamplePacketDataNoSourcePort,
	}

	examplePacketWithSourcePort = channeltypes.Packet{
		SourceChannel: "channel-0",
		SourcePort:    "transfer",
		Sequence:      1,
		Data:          marshaledexamplePacketDataForSourcePortPacket,
	}

	examplePacketNonSDai = channeltypes.Packet{
		SourceChannel: "channel-0",
		Sequence:      1,
		Data:          marshaledExamplePacketDataNonSDai,
	}

	packetEmptyData = channeltypes.Packet{
		SourceChannel: "channel-0",
		Sequence:      1,
	}

	packetFullSDaiPathPacketData = channeltypes.Packet{
		SourceChannel: "channel-0",
		SourcePort:    "transfer",
		Sequence:      1,
		Data:          marshaledExampleFullSDaiPathPacketData,
	}

	exampleAckError = errors.New("ABCI code: 1: error handling packet: see events for details")
)

func TestPendingPacket(t *testing.T) {
	testChannelId := "channel-0"
	testSequence := uint64(20)
	testSequence2 := uint64(22)

	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RatelimitKeeper

	// Set pending packet in state
	k.SetPendingSendPacket(ctx, testChannelId, testSequence)
	k.SetPendingSendPacket(ctx, testChannelId, testSequence2)

	// Test HasPendingSendPacket
	require.True(t, k.HasPendingSendPacket(ctx, testChannelId, testSequence))
	require.True(t, k.HasPendingSendPacket(ctx, testChannelId, testSequence2))
	require.False(t, k.HasPendingSendPacket(ctx, "non-existent-channel", testSequence))
	require.False(t,
		k.HasPendingSendPacket(
			ctx, testChannelId,
			42, // non-existent sequence number
		),
	)

	// Remove pending packet from state
	k.RemovePendingSendPacket(
		ctx,
		testChannelId,
		testSequence,
	)

	require.False(t, k.HasPendingSendPacket(ctx, testChannelId, testSequence)) // Removed
	require.True(t, k.HasPendingSendPacket(ctx, testChannelId, testSequence2))
	require.False(t, k.HasPendingSendPacket(ctx, "non-existent-channel", testSequence))
	require.False(t,
		k.HasPendingSendPacket(
			ctx, testChannelId,
			42, // non-existent sequence number
		),
	)
}

func TestAcknowledgeIBCTransferPacket(t *testing.T) {
	var (
		amountSent   *big.Int
		denomForTest string
	)
	testCases := map[string]struct {
		packet      channeltypes.Packet
		ack         channeltypes.Acknowledgement
		ackSuccess  bool
		customSetup func(*testapp.TestApp, sdk.Context)
		expectedErr string
	}{
		"Success: Ack Success": {
			packet:     examplePacketNoSourcePort,
			ack:        channeltypes.NewResultAcknowledgement([]byte{1}),
			ackSuccess: true,
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetPendingSendPacket(ctx, examplePacketNoSourcePort.SourceChannel, examplePacketNoSourcePort.Sequence)
				denomCapcity := types.DenomCapacity{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(100000000),
						dtypes.NewInt(200000000),
					},
				}
				app.App.RatelimitKeeper.SetDenomCapacity(ctx, denomCapcity)
				denomForTest = "gsdai"
				amountSent = big.NewInt(100000000)
			},
			expectedErr: "",
		},
		"Success: Ack Error": {
			packet:     examplePacketNoSourcePort,
			ack:        channeltypes.NewErrorAcknowledgement(exampleAckError),
			ackSuccess: false,
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetPendingSendPacket(ctx, examplePacketNoSourcePort.SourceChannel, examplePacketNoSourcePort.Sequence)
				denomCapcity := types.DenomCapacity{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(100000000),
						dtypes.NewInt(200000000),
					},
				}
				app.App.RatelimitKeeper.SetDenomCapacity(ctx, denomCapcity)
				denomForTest = "gsdai"
				amountSent = big.NewInt(100000000)
			},
			expectedErr: "",
		},
		"Failure: Ack Error": {
			packet: channeltypes.Packet{
				SourceChannel: "channel-0",
				Sequence:      2,
				Data:          marshaledExamplePacketDataNoSourcePort,
			},
			ack:        channeltypes.Acknowledgement{},
			ackSuccess: false,
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetPendingSendPacket(ctx, "channel-0", 2)
				denomCapcity := types.DenomCapacity{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(100000000),
						dtypes.NewInt(200000000),
					},
				}
				app.App.RatelimitKeeper.SetDenomCapacity(ctx, denomCapcity)
				denomForTest = "gsdai"
				amountSent = big.NewInt(100000000)
			},
			expectedErr: "unsupported acknowledgement response field",
		},
		"Failure: Packet Error": {
			packet: channeltypes.Packet{
				SourceChannel: "channel-0",
				Sequence:      2,
				Data:          []byte(`invalid`),
			},
			ack:        channeltypes.NewResultAcknowledgement([]byte{1}),
			ackSuccess: true,
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetPendingSendPacket(ctx, "channel-0", 2)
				denomCapcity := types.DenomCapacity{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(100000000),
						dtypes.NewInt(200000000),
					},
				}
				app.App.RatelimitKeeper.SetDenomCapacity(ctx, denomCapcity)
				denomForTest = "gsdai"
				amountSent = big.NewInt(100000000)
			},
			expectedErr: "invalid character 'i' looking for beginning of value",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper

			tc.customSetup(tApp, ctx)

			ackBz, err := ibctransfertypes.ModuleCdc.MarshalJSON(&tc.ack)
			require.NoError(t, err, "no error expected when marshalling ack")

			initialDenomCapacity := k.GetDenomCapacity(ctx, denomForTest)

			err = k.AcknowledgeIBCTransferPacket(ctx, tc.packet, ackBz)

			newDenomCapacity := k.GetDenomCapacity(ctx, denomForTest)

			if tc.expectedErr == "" {
				require.NoError(t, err)
				require.False(t, k.HasPendingSendPacket(ctx, tc.packet.SourceChannel, tc.packet.Sequence))
				if tc.ackSuccess {
					require.Equal(t, initialDenomCapacity, newDenomCapacity)
				} else {
					for i, capacity := range newDenomCapacity.CapacityList {
						expectedCapcity := dtypes.NewIntFromBigInt(new(big.Int).Add(initialDenomCapacity.CapacityList[i].BigInt(), amountSent))
						require.Equal(t, expectedCapcity, capacity)
					}
				}
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
				require.Equal(t, initialDenomCapacity, newDenomCapacity)
			}
		})
	}
}

func TestUndoSendPacket(t *testing.T) {
	var (
		amountSent   *big.Int
		denomForTest string
	)
	denomForTest = "gsdai"
	amountSent = big.NewInt(100000000)

	testCases := map[string]struct {
		channelId    string
		sequence     uint64
		customSetup  func(*testapp.TestApp, sdk.Context)
		expectedUndo bool
	}{
		"Correctly undoes sendPacket": {
			channelId: "channel-0",
			sequence:  1,
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetPendingSendPacket(ctx, "channel-0", 1)
			},
			expectedUndo: true,
		},
		"Channel Id not found": {
			channelId: "channel-1",
			sequence:  1,
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetPendingSendPacket(ctx, "channel-1", 1)
			},
			expectedUndo: false,
		},
		"Sequence not found not found": {
			channelId: "channel-0",
			sequence:  2,
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetPendingSendPacket(ctx, "channel-0", 2)
			},
			expectedUndo: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper
			denomCapcity := types.DenomCapacity{
				Denom: denomForTest,
				CapacityList: []dtypes.SerializableInt{
					dtypes.NewInt(100000000),
					dtypes.NewInt(200000000),
				},
			}
			k.SetDenomCapacity(ctx, denomCapcity)
			tc.customSetup(tApp, ctx)

			initialDenomCapacity := k.GetDenomCapacity(ctx, denomForTest)
			k.UndoSendPacket(ctx, examplePacketNoSourcePort.SourceChannel, examplePacketNoSourcePort.Sequence, denomForTest, amountSent)
			newDenomCapacity := k.GetDenomCapacity(ctx, denomForTest)

			if tc.expectedUndo {
				require.False(t, k.HasPendingSendPacket(ctx, tc.channelId, tc.sequence))
				for i, capacity := range newDenomCapacity.CapacityList {
					expectedCapcity := dtypes.NewIntFromBigInt(new(big.Int).Add(initialDenomCapacity.CapacityList[i].BigInt(), amountSent))
					require.Equal(t, expectedCapcity, capacity)
				}
			} else {
				require.True(t, k.HasPendingSendPacket(ctx, tc.channelId, tc.sequence))
				require.Equal(t, initialDenomCapacity, newDenomCapacity)
			}
		})
	}
}

func TestRedoMintTradingDAIIfAcknowledgeIBCTransferPacketFails(t *testing.T) {
	testCases := map[string]struct {
		packet             channeltypes.Packet
		packetData         ibctransfertypes.FungibleTokenPacketData
		ack                channeltypes.Acknowledgement
		ackSuccess         bool
		expectedMintedTDai *big.Int
		customSetup        func(*testapp.TestApp, sdk.Context)
		expectedErr        string
	}{
		"Success: mints when ack is not success": {
			packet:             packetFullSDaiPathPacketData,
			packetData:         fullSDaiPathPacketData,
			ack:                channeltypes.NewErrorAcknowledgement(exampleAckError),
			ackSuccess:         true,
			expectedMintedTDai: big.NewInt(100),
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"))
			},
			expectedErr: "",
		},
		"Success: does not mint when ack is success": {
			packet:             examplePacketNoSourcePort,
			packetData:         examplePacketDataNoSourcePort,
			ack:                channeltypes.NewResultAcknowledgement([]byte{1}),
			ackSuccess:         true,
			expectedMintedTDai: big.NewInt(0),
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"))
			},
			expectedErr: "",
		},
		"Success: does not mint when denom is not gsdai": {
			packet:             examplePacketNonSDai,
			packetData:         examplePacketDataNonSDai,
			ack:                channeltypes.NewErrorAcknowledgement(exampleAckError),
			ackSuccess:         true,
			expectedMintedTDai: big.NewInt(0),
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"))
			},
			expectedErr: "",
		},
		"Failure: returns err when ack cannot be unpacked": {
			packet:             examplePacketNoSourcePort,
			packetData:         examplePacketDataNoSourcePort,
			ack:                channeltypes.Acknowledgement{},
			ackSuccess:         true,
			expectedMintedTDai: big.NewInt(0),
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"))
			},
			expectedErr: "unsupported acknowledgement response field",
		},
		"Failure: returns err when packet info cannot be parsed": {
			packet: channeltypes.Packet{
				SourceChannel: "channel-0",
				Sequence:      1,
				Data:          []byte(`invalid`),
			},
			packetData:         examplePacketDataNoSourcePort,
			ack:                channeltypes.NewErrorAcknowledgement(exampleAckError),
			ackSuccess:         true,
			expectedMintedTDai: big.NewInt(0),
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"))
			},
			expectedErr: "invalid character 'i' looking for beginning of value",
		},
		"Failure: invalid packet sender": {
			packet:             packetEmptyData,
			packetData:         packetDataUnparsableSender,
			ack:                channeltypes.NewErrorAcknowledgement(exampleAckError),
			ackSuccess:         true,
			expectedMintedTDai: big.NewInt(0),
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"))
			},
			expectedErr: "Unable to convert sender address",
		},
		"Failure: invalid packet receiver": {
			packet:             packetEmptyData,
			packetData:         packetDataUnparsableReceiver,
			ack:                channeltypes.NewErrorAcknowledgement(exampleAckError),
			ackSuccess:         true,
			expectedMintedTDai: big.NewInt(0),
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"))
			},
			expectedErr: "Unable to convert receiver address",
		},
		"Failure: invalid packet amount": {
			packet:             packetEmptyData,
			packetData:         packetDataUnparsableAmount,
			ack:                channeltypes.NewErrorAcknowledgement(exampleAckError),
			ackSuccess:         true,
			expectedMintedTDai: big.NewInt(0),
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"))
			},
			expectedErr: "Unable to cast packet amount",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper

			// Normal cases
			accountAddr := sdk.AccAddress("")
			if len(tc.packet.Data) > 0 {
				sDAICoins := sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, sdkmath.NewIntFromBigInt(keeper.ConvertStringToBigIntWithPanicOnErr(tc.packetData.Amount))))
				err := tApp.App.BankKeeper.MintCoins(ctx, types.TDaiPoolAccount, sDAICoins)
				require.NoError(t, err)
				accountAddr, err = sdk.AccAddressFromBech32(tc.packetData.Sender)
				require.NoError(t, err)
				err = tApp.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.TDaiPoolAccount, accountAddr, sDAICoins)
				require.NoError(t, err)
			}

			tc.customSetup(tApp, ctx)
			testPacket := tc.packet

			// Malformed packet data cases
			if len(testPacket.Data) == 0 {
				testPacket.Data = marshalPacketData(tc.packetData)
			}

			ackBz, err := ibctransfertypes.ModuleCdc.MarshalJSON(&tc.ack)
			require.NoError(t, err, "no error expected when marshalling ack")
			initialTDaiBalance := tApp.App.BankKeeper.GetBalance(ctx, accountAddr, assettypes.TDaiDenom)
			err = k.RedoMintTradingDAIIfAcknowledgeIBCTransferPacketFails(ctx, testPacket, ackBz)
			postMintTDaiBalance := tApp.App.BankKeeper.GetBalance(ctx, accountAddr, assettypes.TDaiDenom)

			balanceDiff := postMintTDaiBalance.Amount.Sub(initialTDaiBalance.Amount)
			if tc.expectedErr == "" {
				require.NoError(t, err)
				require.Equal(t, 0, balanceDiff.BigInt().Cmp(tc.expectedMintedTDai))
			} else {
				require.Error(t, err)
				require.Equal(t, 0, balanceDiff.BigInt().Cmp(big.NewInt(0)))
			}
		})
	}
}

func TestTimeoutIBCTransferPacket(t *testing.T) {
	var (
		amountSent   *big.Int
		denomForTest string
	)
	denomForTest = "gsdai"
	amountSent = big.NewInt(100000000)

	testCases := map[string]struct {
		packet       channeltypes.Packet
		packetData   ibctransfertypes.FungibleTokenPacketData
		customSetup  func(*testapp.TestApp, sdk.Context)
		expectedUndo bool
		expectedErr  string
	}{
		"Success: Parses and undoes send packet": {
			packet:       examplePacketNoSourcePort,
			packetData:   examplePacketDataNoSourcePort,
			customSetup:  func(app *testapp.TestApp, ctx sdk.Context) {},
			expectedUndo: true,
		},
		"Success: Parses packet but can't find channelId": {
			packet: channeltypes.Packet{
				SourceChannel: "channel-1",
				Sequence:      examplePacketNoSourcePort.Sequence,
				Data:          examplePacketNoSourcePort.Data,
			},
			packetData:   examplePacketDataNoSourcePort,
			customSetup:  func(app *testapp.TestApp, ctx sdk.Context) {},
			expectedUndo: false,
		},
		"Success: Parses packet but can't find sequence id": {
			packet: channeltypes.Packet{
				SourceChannel: examplePacketNoSourcePort.SourceChannel,
				Sequence:      2,
				Data:          examplePacketNoSourcePort.Data,
			},
			packetData:   examplePacketDataNoSourcePort,
			customSetup:  func(app *testapp.TestApp, ctx sdk.Context) {},
			expectedUndo: false,
		},
		"Failure: Cannot parse packet: cannot parse sender address": {
			packet:       packetEmptyData,
			packetData:   packetDataUnparsableSender,
			customSetup:  func(app *testapp.TestApp, ctx sdk.Context) {},
			expectedUndo: false,
			expectedErr:  "Unable to convert sender address",
		},
		"Failure: Cannot parse packet: cannot parse receiver address": {
			packet:       packetEmptyData,
			packetData:   packetDataUnparsableReceiver,
			customSetup:  func(app *testapp.TestApp, ctx sdk.Context) {},
			expectedUndo: false,
			expectedErr:  "Unable to convert receiver address",
		},
		"Failure: Cannot parse packet: cannot parse amount": {
			packet:       packetEmptyData,
			packetData:   packetDataUnparsableAmount,
			customSetup:  func(app *testapp.TestApp, ctx sdk.Context) {},
			expectedUndo: false,
			expectedErr:  "Unable to cast packet amount",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper

			denomCapcity := types.DenomCapacity{
				Denom: denomForTest,
				CapacityList: []dtypes.SerializableInt{
					dtypes.NewInt(100000000),
					dtypes.NewInt(200000000),
				},
			}
			k.SetDenomCapacity(ctx, denomCapcity)
			k.SetPendingSendPacket(ctx, examplePacketNoSourcePort.SourceChannel, examplePacketNoSourcePort.Sequence)

			testPacket := tc.packet
			testPacket.Data = marshalPacketData(tc.packetData)

			tc.customSetup(tApp, ctx)

			initialDenomCapacity := k.GetDenomCapacity(ctx, denomForTest)
			err := k.TimeoutIBCTransferPacket(ctx, testPacket)
			newDenomCapacity := k.GetDenomCapacity(ctx, denomForTest)

			if tc.expectedErr == "" {
				require.NoError(t, err)
				if tc.expectedUndo {
					require.False(t, k.HasPendingSendPacket(ctx, examplePacketNoSourcePort.SourceChannel, examplePacketNoSourcePort.Sequence))
					for i, capacity := range newDenomCapacity.CapacityList {
						expectedCapcity := dtypes.NewIntFromBigInt(new(big.Int).Add(initialDenomCapacity.CapacityList[i].BigInt(), amountSent))
						require.Equal(t, expectedCapcity, capacity)
					}
				} else {
					require.True(t, k.HasPendingSendPacket(ctx, examplePacketNoSourcePort.SourceChannel, examplePacketNoSourcePort.Sequence))
					require.Equal(t, initialDenomCapacity, newDenomCapacity)
				}
			} else {
				require.Error(t, err)
				require.True(t, k.HasPendingSendPacket(ctx, examplePacketNoSourcePort.SourceChannel, examplePacketNoSourcePort.Sequence))
				require.Equal(t, initialDenomCapacity, newDenomCapacity)
			}
		})
	}
}

func TestUndoMintTradingDAIIfAfterTimeoutIBCTransferPacket(t *testing.T) {
	testCases := map[string]struct {
		packet             channeltypes.Packet
		packetData         ibctransfertypes.FungibleTokenPacketData
		expectedMintedTDai *big.Int
		customSetup        func(*testapp.TestApp, sdk.Context)
		expectedErr        string
	}{
		"Success: mints succesfully": {
			packet:             packetFullSDaiPathPacketData,
			packetData:         fullSDaiPathPacketData,
			expectedMintedTDai: big.NewInt(100),
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"))
			},
			expectedErr: "",
		},
		"Success: nonSDai input does not mint": {
			packet:             examplePacketNonSDai,
			packetData:         examplePacketDataNonSDai,
			expectedMintedTDai: big.NewInt(0),
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"))
			},
			expectedErr: "",
		},
		"Failure: returns err when packet info cannot be parsed": {
			packet: channeltypes.Packet{
				SourceChannel: "channel-0",
				Sequence:      1,
				Data:          []byte(`invalid`),
			},
			packetData:         examplePacketDataNoSourcePort,
			expectedMintedTDai: big.NewInt(0),
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"))
			},
			expectedErr: "invalid character 'i' looking for beginning of value",
		},
		"Failure: invalid packet sender": {
			packet:             packetEmptyData,
			packetData:         packetDataUnparsableSender,
			expectedMintedTDai: big.NewInt(0),
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"))
			},
			expectedErr: "Unable to convert sender address",
		},
		"Failure: invalid packet receiver": {
			packet:             packetEmptyData,
			packetData:         packetDataUnparsableReceiver,
			expectedMintedTDai: big.NewInt(0),
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"))
			},
			expectedErr: "Unable to convert receiver address",
		},
		"Failure: invalid packet amount": {
			packet:             packetEmptyData,
			packetData:         packetDataUnparsableAmount,
			expectedMintedTDai: big.NewInt(0),
			customSetup: func(app *testapp.TestApp, ctx sdk.Context) {
				app.App.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"))
			},
			expectedErr: "Unable to cast packet amount",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper

			// Normal cases
			accountAddr := sdk.AccAddress("")
			if len(tc.packet.Data) > 0 {
				sDAICoins := sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, sdkmath.NewIntFromBigInt(keeper.ConvertStringToBigIntWithPanicOnErr(tc.packetData.Amount))))
				err := tApp.App.BankKeeper.MintCoins(ctx, types.TDaiPoolAccount, sDAICoins)
				require.NoError(t, err)
				accountAddr, err = sdk.AccAddressFromBech32(tc.packetData.Sender)
				require.NoError(t, err)
				err = tApp.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.TDaiPoolAccount, accountAddr, sDAICoins)
				require.NoError(t, err)
			}

			tc.customSetup(tApp, ctx)
			testPacket := tc.packet

			// Malformed packet data cases
			if len(testPacket.Data) == 0 {
				testPacket.Data = marshalPacketData(tc.packetData)
			}

			initialTDaiBalance := tApp.App.BankKeeper.GetBalance(ctx, accountAddr, assettypes.TDaiDenom)
			err := k.UndoMintTradingDAIIfAfterTimeoutIBCTransferPacket(ctx, testPacket)
			postMintTDaiBalance := tApp.App.BankKeeper.GetBalance(ctx, accountAddr, assettypes.TDaiDenom)

			balanceDiff := postMintTDaiBalance.Amount.Sub(initialTDaiBalance.Amount)
			if tc.expectedErr == "" {
				require.NoError(t, err)
				require.Equal(t, 0, balanceDiff.BigInt().Cmp(tc.expectedMintedTDai))
			} else {
				require.Error(t, err)
				require.Equal(t, 0, balanceDiff.BigInt().Cmp(big.NewInt(0)))
			}
		})
	}
}
func TestSendPacket(t *testing.T) {
	var (
		amountSent   = keeper.ConvertStringToBigIntWithPanicOnErr("100000000000000")
		denomForTest = "gsdai"
	)
	testCases := map[string]struct {
		channelCap       *capabilitytypes.Capability
		sourcePort       string
		sourceChannel    string
		timeoutHeight    clienttypes.Height
		timeoutTimestamp uint64
		data             []byte
		denomCapacities  []types.DenomCapacity
		mockICS4Wrapper  func(*mocks.ICS4Wrapper)
		expectedSequence uint64
		expectedErr      string
	}{
		"Success: Basic": {
			channelCap:       &capabilitytypes.Capability{},
			sourcePort:       "transfer",
			sourceChannel:    "channel-0",
			timeoutHeight:    clienttypes.NewHeight(1, 1000),
			timeoutTimestamp: 1000000,
			data:             marshaledexamplePacketDataForSourcePortPacket,
			denomCapacities: []types.DenomCapacity{
				{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(200000000000000),
						dtypes.NewInt(300000000000000),
					},
				},
			},
			mockICS4Wrapper: func(m *mocks.ICS4Wrapper) {
				m.On("SendPacket", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(uint64(1), nil)
			},
			expectedSequence: 1,
			expectedErr:      "",
		},
		"Success: Capacity goes to 0": {
			channelCap:       &capabilitytypes.Capability{},
			sourcePort:       "transfer",
			sourceChannel:    "channel-0",
			timeoutHeight:    clienttypes.NewHeight(1, 1000),
			timeoutTimestamp: 1000000,
			data:             marshaledexamplePacketDataForSourcePortPacket,
			denomCapacities: []types.DenomCapacity{
				{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(100000000000000),
						dtypes.NewInt(200000000000000),
					},
				},
				{
					Denom: "utdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(200000000000000),
						dtypes.NewInt(300000000000000),
					},
				},
			},
			mockICS4Wrapper: func(m *mocks.ICS4Wrapper) {
				m.On("SendPacket", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(uint64(1), nil)
			},
			expectedSequence: 1,
			expectedErr:      "",
		},
		"Success: Multiple capacities": {
			channelCap:       &capabilitytypes.Capability{},
			sourcePort:       "transfer",
			sourceChannel:    "channel-0",
			timeoutHeight:    clienttypes.NewHeight(1, 1000),
			timeoutTimestamp: 1000000,
			data:             marshaledexamplePacketDataForSourcePortPacket,
			denomCapacities: []types.DenomCapacity{
				{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(100000000000000),
					},
				},
			},
			mockICS4Wrapper: func(m *mocks.ICS4Wrapper) {
				m.On("SendPacket", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(uint64(1), nil)
			},
			expectedSequence: 1,
			expectedErr:      "",
		},
		"Failure: ICS4Wrapper.SendPacket fails": {
			channelCap:       &capabilitytypes.Capability{},
			sourcePort:       "transfer",
			sourceChannel:    "channel-0",
			timeoutHeight:    clienttypes.NewHeight(1, 1000),
			timeoutTimestamp: 1000000,
			data:             marshaledexamplePacketDataForSourcePortPacket,
			denomCapacities: []types.DenomCapacity{
				{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(200000000000000),
						dtypes.NewInt(300000000000000),
					},
				},
			},
			mockICS4Wrapper: func(m *mocks.ICS4Wrapper) {
				m.On("SendPacket", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(uint64(0), fmt.Errorf("ICS4Wrapper error"))
			},
			expectedSequence: 0,
			expectedErr:      "ICS4Wrapper error",
		},
		"Failure: Capacity too low": {
			channelCap:       &capabilitytypes.Capability{},
			sourcePort:       "transfer",
			sourceChannel:    "channel-0",
			timeoutHeight:    clienttypes.NewHeight(1, 1000),
			timeoutTimestamp: 1000000,
			data:             marshaledexamplePacketDataForSourcePortPacket,
			denomCapacities: []types.DenomCapacity{
				{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(500000),
					},
				},
			},
			mockICS4Wrapper: func(m *mocks.ICS4Wrapper) {
				m.On("SendPacket", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(uint64(1), nil)
			},
			expectedSequence: 0,
			expectedErr:      "withdrawal amount would exceed rate-limit capacity",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()

			// Create a new mock ICS4Wrapper
			mockICS4Wrapper := mocks.NewICS4Wrapper(t)
			tc.mockICS4Wrapper(mockICS4Wrapper)

			// Set the mock ICS4Wrapper in the keeper.
			// This is akin to the keeper initialization in app.go
			sDAIEventManager := sdaidaemontypes.NewsDAIEventManager()
			tApp.App.RatelimitKeeper = *keeper.NewKeeper(
				tApp.App.AppCodec(),
				tApp.App.GetKey(types.StoreKey),
				sDAIEventManager,
				tApp.App.IndexerEventManager,
				tApp.App.BankKeeper,
				tApp.App.BlockTimeKeeper,
				&tApp.App.PerpetualsKeeper,
				&tApp.App.AssetsKeeper,
				mockICS4Wrapper,
				[]string{
					lib.GovModuleAddress.String(),
					delaymsgmoduletypes.ModuleAddress.String(),
				},
			)

			k := tApp.App.RatelimitKeeper

			for _, denomCapacity := range tc.denomCapacities {
				k.SetDenomCapacity(ctx, denomCapacity)
			}

			sequence, err := k.SendPacket(
				ctx,
				tc.channelCap,
				tc.sourcePort,
				tc.sourceChannel,
				tc.timeoutHeight,
				tc.timeoutTimestamp,
				tc.data,
			)

			if tc.expectedErr == "" {
				require.NoError(t, err)
				require.Equal(t, tc.expectedSequence, sequence)
				require.True(t, k.HasPendingSendPacket(ctx, tc.sourceChannel, sequence))

				for _, initialDenomCapacity := range tc.denomCapacities {
					resultDenomCapacity := k.GetDenomCapacity(ctx, initialDenomCapacity.Denom)
					if initialDenomCapacity.Denom == denomForTest {
						for i, resultCapacity := range resultDenomCapacity.CapacityList {
							expectedCapacity := new(big.Int).Sub(initialDenomCapacity.CapacityList[i].BigInt(), amountSent)
							require.Equal(t, 0, expectedCapacity.Cmp(resultCapacity.BigInt()),
								"Expected capacity %v. Got %v.", expectedCapacity, resultCapacity,
							)
						}
					} else {
						require.Equal(t, initialDenomCapacity, resultDenomCapacity)
					}
				}
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
				require.Equal(t, tc.expectedSequence, sequence)

				if tc.expectedSequence > 0 {
					require.False(t, k.HasPendingSendPacket(ctx, tc.sourceChannel, tc.expectedSequence))
				}

				for _, initialDenomCapacity := range tc.denomCapacities {
					resultDenomCapacity := k.GetDenomCapacity(ctx, initialDenomCapacity.Denom)
					require.Equal(t, initialDenomCapacity, resultDenomCapacity)
				}
			}

			mockICS4Wrapper.AssertExpectations(t)
		})
	}
}

func TestTrySendRateLimitedPacket(t *testing.T) {
	var (
		amountSent   = keeper.ConvertStringToBigIntWithPanicOnErr("100000000000000")
		denomForTest = "gsdai"
	)
	testCases := map[string]struct {
		packet          channeltypes.Packet
		packetData      ibctransfertypes.FungibleTokenPacketData
		denomCapacities []types.DenomCapacity
		expectedErr     string
	}{
		"Success: Basic": {
			packet: examplePacketWithSourcePort,
			denomCapacities: []types.DenomCapacity{
				{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(200000000000000),
						dtypes.NewInt(300000000000000),
					},
				},
			},
			expectedErr: "",
		},
		"Success: Multiple denom capacities": {
			packet: examplePacketWithSourcePort,
			denomCapacities: []types.DenomCapacity{
				{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(200000000000000),
						dtypes.NewInt(300000000000000),
					},
				},
				{
					Denom: "uusdc",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(150000000),
						dtypes.NewInt(450000000),
					},
				},
			},
			expectedErr: "",
		},
		"Success: Capacity goes to 0": {
			packet: examplePacketWithSourcePort,
			denomCapacities: []types.DenomCapacity{
				{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(100000000000000),
					},
				},
			},
			expectedErr: "",
		},
		"Failure: Capacity too low": {
			packet: examplePacketNoSourcePort,
			denomCapacities: []types.DenomCapacity{
				{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(500000),
					},
				},
			},
			expectedErr: "capacity too low",
		},
		"Failure: Invalid packet": {
			packet: channeltypes.Packet{
				SourceChannel: "channel-0",
				Sequence:      1,
				Data:          []byte(`invalid`),
			},
			denomCapacities: []types.DenomCapacity{
				{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(500000),
					},
				},
			},
			expectedErr: "invalid character 'i' looking for beginning of value",
		},
		"Failure: invalid packet sender": {
			packet:     packetEmptyData,
			packetData: packetDataUnparsableSender,
			denomCapacities: []types.DenomCapacity{
				{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(200000000),
						dtypes.NewInt(300000000),
					},
				},
				{
					Denom: "uusdc",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(150000000),
						dtypes.NewInt(450000000),
					},
				},
			},
			expectedErr: "Unable to convert sender address",
		},
		"Failure: invalid packet receiver": {
			packet:     packetEmptyData,
			packetData: packetDataUnparsableReceiver,
			denomCapacities: []types.DenomCapacity{
				{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(200000000),
						dtypes.NewInt(300000000),
					},
				},
				{
					Denom: "uusdc",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(150000000),
						dtypes.NewInt(450000000),
					},
				},
			},
			expectedErr: "Unable to convert receiver address",
		},
		"Failure: invalid packet amount": {
			packet:     packetEmptyData,
			packetData: packetDataUnparsableAmount,
			denomCapacities: []types.DenomCapacity{
				{
					Denom: "gsdai",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(200000000),
						dtypes.NewInt(300000000),
					},
				},
				{
					Denom: "uusdc",
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(150000000),
						dtypes.NewInt(450000000),
					},
				},
			},
			expectedErr: "Unable to cast packet amount",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper
			for _, denomCapacity := range tc.denomCapacities {
				k.SetDenomCapacity(ctx, denomCapacity)
			}

			// Malformed packet data cases
			testPacket := tc.packet
			if len(testPacket.Data) == 0 {
				testPacket.Data = marshalPacketData(tc.packetData)
			}

			err := k.TrySendRateLimitedPacket(ctx, tc.packet)

			hasPendingPacket := k.HasPendingSendPacket(ctx, testPacket.SourceChannel, testPacket.Sequence)

			if tc.expectedErr == "" {
				require.NoError(t, err)
				require.True(t, hasPendingPacket)
				for _, initialDenomCapacity := range tc.denomCapacities {
					resultDenomCapacity := k.GetDenomCapacity(ctx, initialDenomCapacity.Denom)
					if initialDenomCapacity.Denom == denomForTest {
						for i, resultCapacity := range resultDenomCapacity.CapacityList {
							expectedCapacity := new(big.Int).Sub(initialDenomCapacity.CapacityList[i].BigInt(), amountSent)
							require.Equal(t, 0, expectedCapacity.Cmp(resultCapacity.BigInt()),
								"Expected capacity %v. Got %v.", expectedCapacity, resultCapacity,
							)
						}
					} else {
						require.Equal(t, initialDenomCapacity, resultDenomCapacity)
					}
				}
			} else {
				require.Error(t, err)
				require.False(t, hasPendingPacket)
				for _, initialDenomCapacity := range tc.denomCapacities {
					resultDenomCapacity := k.GetDenomCapacity(ctx, initialDenomCapacity.Denom)
					require.Equal(t, initialDenomCapacity, resultDenomCapacity)
				}
			}
		})
	}
}

func TestPreprocessSendPacket(t *testing.T) {
	testCases := map[string]struct {
		packetData                 ibctransfertypes.FungibleTokenPacketData
		initialTDaiBalance         sdkmath.Int
		initialSDAIBalanceInModule sdkmath.Int
		expectedTDaiBurned         sdkmath.Int
		expectedSDAISent           sdkmath.Int
		sDAIPrice                  *big.Int
		customSetup                func(*testapp.TestApp, sdk.Context)
		expectedErr                string
	}{
		"Success: burns tDAI and mints sDAI": {
			packetData: ibctransfertypes.FungibleTokenPacketData{
				Denom:    types.SDaiBaseDenomFullPath,
				Amount:   "1000000000000000000",
				Sender:   testAddress1,
				Receiver: testAddress2,
			},
			initialTDaiBalance:         sdkmath.NewInt(2000000),
			initialSDAIBalanceInModule: sdkmath.NewInt(2000000000000000000),
			expectedTDaiBurned:         sdkmath.NewInt(1000000),
			expectedSDAISent:           sdkmath.NewInt(1000000000000000000),
			sDAIPrice:                  keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedErr:                "",
		},
		"Failure: insufficient tDAI balance": {
			packetData: ibctransfertypes.FungibleTokenPacketData{
				Denom:    types.SDaiBaseDenomFullPath,
				Amount:   "2000000000000000000",
				Sender:   testAddress1,
				Receiver: testAddress2,
			},
			initialTDaiBalance:         sdkmath.NewInt(1000000),
			initialSDAIBalanceInModule: sdkmath.NewInt(2000000000000000000),
			expectedTDaiBurned:         sdkmath.NewInt(0),
			expectedSDAISent:           sdkmath.NewInt(0),
			sDAIPrice:                  keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedErr:                "insufficient funds",
		},
		"Success: burns fractional tDAI and mints sDAI": {
			packetData: ibctransfertypes.FungibleTokenPacketData{
				Denom:    types.SDaiBaseDenomFullPath,
				Amount:   "500000000000000000",
				Sender:   testAddress1,
				Receiver: testAddress2,
			},
			initialTDaiBalance:         sdkmath.NewInt(1000000),
			initialSDAIBalanceInModule: sdkmath.NewInt(5000000000000000000),
			expectedTDaiBurned:         sdkmath.NewInt(500000),
			expectedSDAISent:           sdkmath.NewInt(500000000000000000),
			sDAIPrice:                  keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedErr:                "",
		},
		"Success: no change if denom is not the full sDAI base denom path": {
			packetData: ibctransfertypes.FungibleTokenPacketData{
				Denom:    types.TDaiDenom,
				Amount:   "500000000000000000",
				Sender:   testAddress1,
				Receiver: testAddress2,
			},
			initialTDaiBalance:         sdkmath.NewInt(2000000),
			initialSDAIBalanceInModule: sdkmath.NewInt(1000000000000000000),
			expectedTDaiBurned:         sdkmath.NewInt(0),
			expectedSDAISent:           sdkmath.NewInt(0),
			sDAIPrice:                  keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedErr:                "",
		},
		"Failure: invalid amount in packet": {
			packetData: ibctransfertypes.FungibleTokenPacketData{
				Denom:    types.SDaiBaseDenomFullPath,
				Amount:   "invalid_amount",
				Sender:   testAddress1,
				Receiver: testAddress2,
			},
			initialTDaiBalance:         sdkmath.NewInt(2000000),
			initialSDAIBalanceInModule: sdkmath.NewInt(1000000000000000000),
			expectedTDaiBurned:         sdkmath.NewInt(0),
			expectedSDAISent:           sdkmath.NewInt(0),
			sDAIPrice:                  keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedErr:                "Unable to cast packet amount ",
		},
		"Failure: invalid sender address in packet": {
			packetData: ibctransfertypes.FungibleTokenPacketData{
				Denom:    types.SDaiBaseDenomFullPath,
				Amount:   "1000000000000000000",
				Sender:   "invalid_address",
				Receiver: testAddress2,
			},
			initialTDaiBalance:         sdkmath.NewInt(0),
			initialSDAIBalanceInModule: sdkmath.NewInt(1000000000000000000),
			expectedTDaiBurned:         sdkmath.NewInt(1000000),
			expectedSDAISent:           sdkmath.NewInt(1000000000000000000),
			sDAIPrice:                  keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedErr:                "Unable to convert sender address",
		},
		"Failure: invalid receiver address in packet": {
			packetData: ibctransfertypes.FungibleTokenPacketData{
				Denom:    types.SDaiBaseDenomFullPath,
				Amount:   "1000000000000000000",
				Sender:   testAddress1,
				Receiver: "invalid_sender",
			},
			initialTDaiBalance:         sdkmath.NewInt(2000000),
			initialSDAIBalanceInModule: sdkmath.NewInt(1000000000000000000),
			expectedTDaiBurned:         sdkmath.NewInt(0),
			expectedSDAISent:           sdkmath.NewInt(0),
			sDAIPrice:                  keeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedErr:                "Unable to convert receiver address",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper

			burnAllCoinsOfDenom(t, ctx, tApp.App.BankKeeper, types.TDaiDenom)
			burnAllCoinsOfDenom(t, ctx, tApp.App.BankKeeper, types.SDaiDenom)

			// Setup
			var accountAddr sdk.AccAddress
			var err error
			if tc.packetData.Sender != "invalid_address" {
				accountAddr, err = sdk.AccAddressFromBech32(tc.packetData.Sender)
				require.NoError(t, err)

				// Mint initial tDAI balance
				tDAICoins := sdk.NewCoins(sdk.NewCoin(assettypes.TDaiDenom, tc.initialTDaiBalance))
				err = tApp.App.BankKeeper.MintCoins(ctx, types.TDaiPoolAccount, tDAICoins)
				require.NoError(t, err)
				err = tApp.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.TDaiPoolAccount, accountAddr, tDAICoins)
				require.NoError(t, err)
			}

			// Mint initial sDAI balance to module
			sDAICoins := sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, tc.initialSDAIBalanceInModule))
			err = tApp.App.BankKeeper.MintCoins(ctx, types.TDaiPoolAccount, sDAICoins)
			require.NoError(t, err)
			err = tApp.App.BankKeeper.SendCoinsFromModuleToModule(ctx, types.TDaiPoolAccount, types.SDaiPoolAccount, sDAICoins)
			require.NoError(t, err)

			// Set sDAI price
			k.SetSDAIPrice(ctx, tc.sDAIPrice)

			// Get initial supplies
			initialTDaiSupply := tApp.App.BankKeeper.GetSupply(ctx, assettypes.TDaiDenom)
			initialSDAISupply := tApp.App.BankKeeper.GetSupply(ctx, types.SDaiDenom)

			// Execute
			err = k.PreprocessSendPacket(ctx, marshalPacketData(tc.packetData))

			// Get final supplies
			finalTDaiSupply := tApp.App.BankKeeper.GetSupply(ctx, assettypes.TDaiDenom)
			finalSDAISupply := tApp.App.BankKeeper.GetSupply(ctx, types.SDaiDenom)

			if tc.expectedErr == "" {
				require.NoError(t, err)

				// Check tDAI balance
				tDAIBalance := tApp.App.BankKeeper.GetBalance(ctx, accountAddr, assettypes.TDaiDenom)
				expectedTDAIBalance := tc.initialTDaiBalance.Sub(tc.expectedTDaiBurned)
				require.Equal(t, expectedTDAIBalance, tDAIBalance.Amount)

				// Check sDAI balance
				sDAIBalance := tApp.App.BankKeeper.GetBalance(ctx, accountAddr, types.SDaiDenom)
				require.Equal(t, tc.expectedSDAISent, sDAIBalance.Amount)

				// Check supplies
				require.Equal(t, initialTDaiSupply.Amount.Sub(tc.expectedTDaiBurned), finalTDaiSupply.Amount)
				require.Equal(t, initialSDAISupply.Amount, finalSDAISupply.Amount)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)

				// Check that balances haven't changed
				tDAIBalance := tApp.App.BankKeeper.GetBalance(ctx, accountAddr, assettypes.TDaiDenom)
				require.Equal(t, tc.initialTDaiBalance, tDAIBalance.Amount)

				sDAIBalance := tApp.App.BankKeeper.GetBalance(ctx, accountAddr, types.SDaiDenom)
				require.Equal(t, 0, big.NewInt(0).Cmp(sDAIBalance.Amount.BigInt()))

				// Check that supplies haven't changed
				require.Equal(t, initialTDaiSupply.Amount, finalTDaiSupply.Amount)
				require.Equal(t, initialSDAISupply.Amount, finalSDAISupply.Amount)
			}
		})
	}
}

func marshalPacketData(packetData ibctransfertypes.FungibleTokenPacketData) []byte {
	marshaledPacketData, err := json.Marshal(packetData)
	if err != nil {
		panic("Could not set up test")
	}
	return marshaledPacketData
}
