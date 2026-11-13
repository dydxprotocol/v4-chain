package v_9_9_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/app"
	canonicalusdckeeper "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/keeper"
	canonicalusdctypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	ratelimittypes "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
)

const (
	injectiveCounterparty = "channel-258"
	nobleCounterparty     = "channel-33"
	usdc                  = int64(1_000_000)
)

type harness struct {
	t       *testing.T
	app     *app.App
	k       *canonicalusdckeeper.Keeper
	stack   porttypes.IBCModule
	relayer sdk.AccAddress
	logical string
}

func (h *harness) ledger(ctx sdk.Context) string {
	l := h.k.GetLedger(ctx)
	return fmt.Sprintf("mode=%s nb=%s ib=%s pi=%s pn=%s rc=%s pending=%d supply=%s",
		strings.TrimPrefix(h.k.GetControls(ctx).Mode.String(), "MODE_"),
		l.NobleBacking.BigInt(), l.InjectiveBacking.BigInt(), l.PendingInjective.BigInt(), l.PendingNoble.BigInt(),
		l.RefillableCapacity.BigInt(), h.k.PendingCount(ctx), h.supply(ctx))
}

func (h *harness) supply(ctx sdk.Context) *big.Int {
	return h.app.BankKeeper.GetSupply(ctx, h.logical).Amount.BigInt()
}

func (h *harness) bal(ctx sdk.Context, addr sdk.AccAddress, denom string) *big.Int {
	return h.app.BankKeeper.GetBalance(ctx, addr, denom).Amount.BigInt()
}

func (h *harness) capacities(ctx sdk.Context) []*big.Int {
	var out []*big.Int
	for _, c := range h.app.RatelimitKeeper.GetDenomCapacity(ctx, h.logical).CapacityList {
		out = append(out, c.BigInt())
	}
	return out
}

func (h *harness) invariantOK(ctx sdk.Context) {
	msg, broken := canonicalusdckeeper.BackingInvariant(*h.k)(ctx)
	require.False(h.t, broken, msg)
}

// recv delivers an inbound ICS20 packet through the transfer stack (the IBC core
// proof checks sit above the callback and are not exercised here).
func (h *harness) recv(ctx sdk.Context, dstChannel, srcChannel, denom string, amount int64, receiver string) (bool, string) {
	data := transfertypes.NewFungibleTokenPacketData(denom, strconv.FormatInt(amount, 10), "sender-on-counterparty", receiver, "")
	seq, _ := h.app.IBCKeeper.ChannelKeeper.GetNextSequenceRecv(ctx, transfertypes.PortID, dstChannel)
	packet := channeltypes.NewPacket(data.GetBytes(), seq, transfertypes.PortID, srcChannel, transfertypes.PortID, dstChannel,
		clienttypes.ZeroHeight(), uint64(ctx.BlockTime().Add(time.Hour).UnixNano()))
	ack := h.stack.OnRecvPacket(ctx, packet, h.relayer)
	if ack.Success() {
		return true, ""
	}
	var a channeltypes.Acknowledgement
	require.NoError(h.t, transfertypes.ModuleCdc.UnmarshalJSON(ack.Acknowledgement(), &a))
	return false, a.GetError()
}

func (h *harness) send(ctx sdk.Context, sender sdk.AccAddress, channel, denom string, amount int64, receiver, memo string) (channeltypes.Packet, error) {
	msg := transfertypes.NewMsgTransfer(transfertypes.PortID, channel, sdk.NewInt64Coin(denom, amount), sender.String(), receiver,
		clienttypes.ZeroHeight(), uint64(ctx.BlockTime().Add(time.Hour).UnixNano()), memo)
	res, err := h.app.MsgServiceRouter().Handler(msg)(ctx, msg)
	if err != nil {
		return channeltypes.Packet{}, err
	}
	for _, ev := range res.Events {
		if ev.Type != channeltypes.EventTypeSendPacket {
			continue
		}
		attrs := map[string]string{}
		for _, a := range ev.Attributes {
			attrs[a.Key] = a.Value
		}
		data, err := hex.DecodeString(attrs[channeltypes.AttributeKeyDataHex])
		require.NoError(h.t, err)
		seq, err := strconv.ParseUint(attrs[channeltypes.AttributeKeySequence], 10, 64)
		require.NoError(h.t, err)
		th, err := clienttypes.ParseHeight(attrs[channeltypes.AttributeKeyTimeoutHeight])
		require.NoError(h.t, err)
		tt, err := strconv.ParseUint(attrs[channeltypes.AttributeKeyTimeoutTimestamp], 10, 64)
		require.NoError(h.t, err)
		return channeltypes.NewPacket(data, seq, attrs[channeltypes.AttributeKeySrcPort], attrs[channeltypes.AttributeKeySrcChannel],
			attrs[channeltypes.AttributeKeyDstPort], attrs[channeltypes.AttributeKeyDstChannel], th, tt), nil
	}
	return channeltypes.Packet{}, fmt.Errorf("no send_packet event")
}

func (h *harness) ack(ctx sdk.Context, packet channeltypes.Packet, success bool) error {
	ack := channeltypes.NewResultAcknowledgement([]byte{1}).Acknowledgement()
	if !success {
		ack = channeltypes.NewErrorAcknowledgement(fmt.Errorf("counterparty rejected")).Acknowledgement()
	}
	return h.stack.OnAcknowledgementPacket(ctx, packet, ack, h.relayer)
}

func (h *harness) timeout(ctx sdk.Context, packet channeltypes.Packet) error {
	return h.stack.OnTimeoutPacket(ctx, packet, h.relayer)
}

func (h *harness) setControls(ctx sdk.Context, mutate func(*canonicalusdctypes.Controls)) error {
	c := h.k.GetControls(ctx)
	mutate(&c)
	return h.k.ApplyControls(ctx, c)
}

// findHolder returns a non-module account holding at least min of denom.
func (h *harness) findHolder(ctx sdk.Context, denom string, min int64) sdk.AccAddress {
	var found sdk.AccAddress
	h.app.BankKeeper.IterateAllBalances(ctx, func(addr sdk.AccAddress, coin sdk.Coin) bool {
		if coin.Denom != denom || coin.Amount.LT(sdkmath.NewInt(min)) {
			return false
		}
		if _, isModule := h.app.AccountKeeper.GetAccount(ctx, addr).(sdk.ModuleAccountI); isModule {
			return false
		}
		found = addr
		return true
	})
	require.NotNil(h.t, found, "no holder of %s with at least %d", denom, min)
	return found
}

func abciCode(err error) string {
	_, code, _ := errorsmod.ABCIInfo(err, false)
	return fmt.Sprintf("ABCI code: %d:", code)
}

func TestMainnetStateScenarios(t *testing.T) {
	if os.Getenv("MAINNET_HOME") == "" {
		t.Skip("MAINNET_HOME not set")
	}
	dydxApp, base, _, _ := loadMainnet(t)
	dydxApp.CapabilityKeeper.InitMemStore(base)
	stack, ok := dydxApp.IBCKeeper.Router.GetRoute(transfertypes.ModuleName)
	require.True(t, ok)
	controls := dydxApp.CanonicalUsdcKeeper.GetControls(base)
	h := &harness{t: t, app: dydxApp, k: dydxApp.CanonicalUsdcKeeper, stack: stack, relayer: sdk.AccAddress(make([]byte, 20)), logical: controls.LogicalDenom}
	user := sdk.AccAddress([]byte("mainnet-dryrun-user!"))
	module := dydxApp.AccountKeeper.GetModuleAddress(canonicalusdctypes.ModuleName)
	voucher := controls.InjectiveDenom
	injReceiver := "inj1cml96vmptgw99syqrrz8az79xer2pcgp0a885r"
	nobleReceiver := "noble1cml96vmptgw99syqrrz8az79xer2pcgpsp2jnk"

	limits := dydxApp.RatelimitKeeper.GetLimitParams(base, h.logical)
	fmt.Printf("SCENARIO limiters on %s: %v\n", h.logical, limits.Limiters)
	fmt.Printf("SCENARIO capacities now: %v\n", h.capacities(base))
	fmt.Printf("SCENARIO base: %s\n", h.ledger(base))
	h.invariantOK(base)
	supply0 := h.supply(base)

	t.Run("noble_deposit_bounces", func(t *testing.T) {
		ctx, _ := base.CacheContext()
		ok, errMsg := h.recv(ctx, controls.NobleChannel, nobleCounterparty, controls.NoblePacketDenom, 1000*usdc, user.String())
		fmt.Printf("SCENARIO noble deposit: ok=%v ack=%q | %s\n", ok, errMsg, h.ledger(ctx))
		require.False(t, ok)
		require.Contains(t, errMsg, abciCode(canonicalusdctypes.ErrNobleDepositsDisabled))
		require.Equal(t, 0, h.supply(ctx).Cmp(supply0))
		require.Equal(t, int64(0), h.bal(ctx, user, h.logical).Int64())
	})

	deposited, _ := base.CacheContext()
	t.Run("injective_deposit", func(t *testing.T) {
		caps := h.capacities(deposited)
		ok, errMsg := h.recv(deposited, controls.InjectiveChannel, injectiveCounterparty, controls.InjectivePacketDenom, 1000*usdc, user.String())
		fmt.Printf("SCENARIO injective deposit 1000: ok=%v ack=%q | %s\n", ok, errMsg, h.ledger(deposited))
		require.True(t, ok, errMsg)
		require.Equal(t, 1000*usdc, h.bal(deposited, user, h.logical).Int64())
		require.Equal(t, 1000*usdc, h.bal(deposited, module, voucher).Int64())
		require.Equal(t, 1000*usdc, h.k.GetLedger(deposited).InjectiveBacking.BigInt().Int64())
		require.Equal(t, 0, h.supply(deposited).Cmp(new(big.Int).Add(supply0, big.NewInt(1000*usdc))))
		require.Equal(t, caps, h.capacities(deposited), "a deposit with no refillable capacity leaves the limiter alone")
		h.invariantOK(deposited)
	})

	t.Run("injective_withdraw_ack", func(t *testing.T) {
		ctx, _ := deposited.CacheContext()
		caps := h.capacities(ctx)
		packet, err := h.send(ctx, user, controls.InjectiveChannel, h.logical, 400*usdc, injReceiver, "")
		require.NoError(t, err)
		fmt.Printf("SCENARIO withdraw 400 sent seq=%d | %s | capacities %v -> %v\n", packet.Sequence, h.ledger(ctx), caps, h.capacities(ctx))
		l := h.k.GetLedger(ctx)
		require.Equal(t, 600*usdc, l.InjectiveBacking.BigInt().Int64())
		require.Equal(t, 400*usdc, l.PendingInjective.BigInt().Int64())
		require.Equal(t, 400*usdc, l.RefillableCapacity.BigInt().Int64())
		require.Equal(t, 600*usdc, h.bal(ctx, user, h.logical).Int64())
		require.Equal(t, 600*usdc, h.bal(ctx, module, voucher).Int64())
		p, found := h.k.GetPendingForPacket(ctx, packet.SourceChannel, packet.Sequence)
		require.True(t, found)
		require.Equal(t, canonicalusdctypes.Route_ROUTE_INJECTIVE, p.Route)
		for i, c := range h.capacities(ctx) {
			require.Equal(t, 0, c.Cmp(new(big.Int).Sub(caps[i], big.NewInt(400*usdc))), "limiter %d charged", i)
		}
		h.invariantOK(ctx)
		require.NoError(t, h.ack(ctx, packet, true))
		fmt.Printf("SCENARIO withdraw 400 acked | %s\n", h.ledger(ctx))
		l = h.k.GetLedger(ctx)
		require.Equal(t, int64(0), l.PendingInjective.BigInt().Int64())
		require.Equal(t, 600*usdc, l.InjectiveBacking.BigInt().Int64())
		require.Equal(t, uint32(0), h.k.PendingCount(ctx))
		h.invariantOK(ctx)
	})

	for _, tc := range []struct {
		name string
		undo func(sdk.Context, channeltypes.Packet) error
	}{
		{"injective_withdraw_timeout", h.timeout},
		{"injective_withdraw_error_ack", func(ctx sdk.Context, p channeltypes.Packet) error { return h.ack(ctx, p, false) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _ := deposited.CacheContext()
			caps := h.capacities(ctx)
			packet, err := h.send(ctx, user, controls.InjectiveChannel, h.logical, 400*usdc, injReceiver, "")
			require.NoError(t, err)
			require.NoError(t, tc.undo(ctx, packet))
			fmt.Printf("SCENARIO %s | %s | capacities %v\n", tc.name, h.ledger(ctx), h.capacities(ctx))
			l := h.k.GetLedger(ctx)
			require.Equal(t, 1000*usdc, l.InjectiveBacking.BigInt().Int64())
			require.Equal(t, int64(0), l.PendingInjective.BigInt().Int64())
			require.Equal(t, int64(0), l.RefillableCapacity.BigInt().Int64())
			require.Equal(t, 1000*usdc, h.bal(ctx, user, h.logical).Int64())
			require.Equal(t, 1000*usdc, h.bal(ctx, module, voucher).Int64())
			require.Equal(t, caps, h.capacities(ctx), "limiter charge undone")
			h.invariantOK(ctx)
		})
	}

	t.Run("injective_over_withdraw_and_memo", func(t *testing.T) {
		ctx, _ := deposited.CacheContext()
		_, err := h.send(ctx, user, controls.InjectiveChannel, h.logical, 2000*usdc, injReceiver, "")
		require.ErrorIs(t, err, canonicalusdctypes.ErrInsufficientBacking)
		_, err = h.send(ctx, user, controls.InjectiveChannel, h.logical, 100*usdc, injReceiver, "hello")
		require.ErrorIs(t, err, canonicalusdctypes.ErrMemoNotAllowed)
		_, err = h.send(ctx, user, controls.InjectiveChannel, h.logical, 100*usdc, injReceiver, `{"wasm":{"contract":"inj1x","msg":{}}}`)
		require.ErrorIs(t, err, canonicalusdctypes.ErrMemoNotAllowed)
		fmt.Printf("SCENARIO over-withdraw and memos rejected | %s\n", h.ledger(ctx))
		require.Equal(t, 1000*usdc, h.bal(ctx, user, h.logical).Int64())
	})

	t.Run("rate_limit_on_injective_withdrawal", func(t *testing.T) {
		ctx, _ := base.CacheContext()
		caps := h.capacities(ctx)
		if len(caps) == 0 {
			t.Skip("no limiters")
		}
		smallest := caps[0]
		for _, c := range caps[1:] {
			if c.Cmp(smallest) < 0 {
				smallest = c
			}
		}
		over := new(big.Int).Add(smallest, big.NewInt(usdc))
		deposit := new(big.Int).Add(over, big.NewInt(usdc))
		ok, errMsg := h.recv(ctx, controls.InjectiveChannel, injectiveCounterparty, controls.InjectivePacketDenom, deposit.Int64(), user.String())
		require.True(t, ok, errMsg)
		_, err := h.send(ctx, user, controls.InjectiveChannel, h.logical, over.Int64(), injReceiver, "")
		fmt.Printf("SCENARIO withdraw %s over the smallest capacity %s: %v\n", over, smallest, err)
		require.ErrorIs(t, err, ratelimittypes.ErrWithdrawalExceedsCapacity)
		within := new(big.Int).Sub(smallest, big.NewInt(usdc))
		_, err = h.send(ctx, user, controls.InjectiveChannel, h.logical, within.Int64(), injReceiver, "")
		require.NoError(t, err)
		fmt.Printf("SCENARIO withdraw %s within capacity ok | %s | capacities %v\n", within, h.ledger(ctx), h.capacities(ctx))
		h.invariantOK(ctx)
	})

	whale := h.findHolder(base, h.logical, 50_000*usdc)
	fmt.Printf("SCENARIO noble withdrawer %s holds %s\n", whale, h.bal(base, whale, h.logical))
	t.Run("noble_withdraw_ack_and_timeout", func(t *testing.T) {
		ctx, _ := base.CacheContext()
		caps := h.capacities(ctx)
		before := h.bal(ctx, whale, h.logical)
		_, err := h.send(ctx, whale, controls.NobleChannel, h.logical, 10_000*usdc, nobleReceiver, `{"forward":{"receiver":"osmo1x","port":"transfer","channel":"channel-1"}}`)
		require.ErrorIs(t, err, canonicalusdctypes.ErrMemoNotAllowed)
		packet, err := h.send(ctx, whale, controls.NobleChannel, h.logical, 10_000*usdc, nobleReceiver, "plain memo is fine")
		require.NoError(t, err)
		fmt.Printf("SCENARIO noble withdraw 10000 sent seq=%d | %s | capacities %v -> %v\n", packet.Sequence, h.ledger(ctx), caps, h.capacities(ctx))
		l := h.k.GetLedger(ctx)
		require.Equal(t, 10_000*usdc, l.PendingNoble.BigInt().Int64())
		require.Equal(t, 0, l.NobleBacking.BigInt().Cmp(new(big.Int).Sub(supply0, big.NewInt(10_000*usdc))))
		require.Equal(t, 0, h.supply(ctx).Cmp(new(big.Int).Sub(supply0, big.NewInt(10_000*usdc))), "Noble vouchers burn at send")
		require.Equal(t, 0, h.bal(ctx, whale, h.logical).Cmp(new(big.Int).Sub(before, big.NewInt(10_000*usdc))))
		require.Equal(t, int64(0), l.RefillableCapacity.BigInt().Int64(), "Noble withdrawals never accrue refill")
		for i, c := range h.capacities(ctx) {
			require.Equal(t, 0, c.Cmp(new(big.Int).Sub(caps[i], big.NewInt(10_000*usdc))), "limiter %d charged", i)
		}
		h.invariantOK(ctx)
		acked, _ := ctx.CacheContext()
		require.NoError(t, h.ack(acked, packet, true))
		fmt.Printf("SCENARIO noble withdraw acked | %s\n", h.ledger(acked))
		require.Equal(t, int64(0), h.k.GetLedger(acked).PendingNoble.BigInt().Int64())
		h.invariantOK(acked)
		require.NoError(t, h.timeout(ctx, packet))
		fmt.Printf("SCENARIO noble withdraw timed out | %s | capacities %v\n", h.ledger(ctx), h.capacities(ctx))
		require.Equal(t, 0, h.supply(ctx).Cmp(supply0))
		require.Equal(t, 0, h.k.GetLedger(ctx).NobleBacking.BigInt().Cmp(supply0))
		require.Equal(t, 0, h.bal(ctx, whale, h.logical).Cmp(before))
		require.Equal(t, caps, h.capacities(ctx))
		h.invariantOK(ctx)
	})

	t.Run("governance_pause_and_deposit_toggle", func(t *testing.T) {
		ctx, _ := deposited.CacheContext()
		require.NoError(t, h.setControls(ctx, func(c *canonicalusdctypes.Controls) { c.Mode = canonicalusdctypes.Mode_MODE_PAUSED }))
		ok, errMsg := h.recv(ctx, controls.InjectiveChannel, injectiveCounterparty, controls.InjectivePacketDenom, 10*usdc, user.String())
		require.False(t, ok)
		require.Contains(t, errMsg, abciCode(canonicalusdctypes.ErrPaused))
		_, err := h.send(ctx, user, controls.InjectiveChannel, h.logical, 10*usdc, injReceiver, "")
		require.ErrorIs(t, err, canonicalusdctypes.ErrPaused)
		_, err = h.send(ctx, whale, controls.NobleChannel, h.logical, 10*usdc, nobleReceiver, "")
		require.ErrorIs(t, err, canonicalusdctypes.ErrPaused)
		fmt.Printf("SCENARIO paused: deposit ack=%q, withdrawals rejected | %s\n", errMsg, h.ledger(ctx))
		require.NoError(t, h.setControls(ctx, func(c *canonicalusdctypes.Controls) { c.Mode = canonicalusdctypes.Mode_MODE_GRADUAL }))
		require.NoError(t, h.setControls(ctx, func(c *canonicalusdctypes.Controls) { c.InjectiveDepositsDisabled = true }))
		ok, errMsg = h.recv(ctx, controls.InjectiveChannel, injectiveCounterparty, controls.InjectivePacketDenom, 10*usdc, user.String())
		require.False(t, ok)
		require.Contains(t, errMsg, abciCode(canonicalusdctypes.ErrInjectiveDepositsDisabled))
		_, err = h.send(ctx, user, controls.InjectiveChannel, h.logical, 10*usdc, injReceiver, "")
		require.NoError(t, err, "withdrawals stay open with deposits disabled")
		require.NoError(t, h.setControls(ctx, func(c *canonicalusdctypes.Controls) { c.InjectiveDepositsDisabled = false }))
		ok, errMsg = h.recv(ctx, controls.InjectiveChannel, injectiveCounterparty, controls.InjectivePacketDenom, 10*usdc, user.String())
		require.True(t, ok, errMsg)
		fmt.Printf("SCENARIO restored: deposit ok | %s\n", h.ledger(ctx))
		h.invariantOK(ctx)
	})

	t.Run("pending_cap", func(t *testing.T) {
		ctx, _ := deposited.CacheContext()
		require.NoError(t, h.setControls(ctx, func(c *canonicalusdctypes.Controls) { c.MaxPendingSettlements = 2 }))
		p1, err := h.send(ctx, user, controls.InjectiveChannel, h.logical, 10*usdc, injReceiver, "")
		require.NoError(t, err)
		_, err = h.send(ctx, user, controls.InjectiveChannel, h.logical, 10*usdc, injReceiver, "")
		require.NoError(t, err)
		_, err = h.send(ctx, user, controls.InjectiveChannel, h.logical, 10*usdc, injReceiver, "")
		require.ErrorIs(t, err, canonicalusdctypes.ErrPendingLimit)
		err = h.setControls(ctx, func(c *canonicalusdctypes.Controls) { c.MaxPendingSettlements = 1 })
		require.Error(t, err, "cap below the live count")
		require.NoError(t, h.ack(ctx, p1, true))
		_, err = h.send(ctx, user, controls.InjectiveChannel, h.logical, 10*usdc, injReceiver, "")
		require.NoError(t, err, "settlement frees a slot")
		fmt.Printf("SCENARIO pending cap: third refused, cap 1 rejected (%v), slot freed by ack | %s\n", err, h.ledger(ctx))
		h.invariantOK(ctx)
	})

	t.Run("endblock_auto_pause", func(t *testing.T) {
		ctx, _ := deposited.CacheContext()
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		require.NoError(t, h.app.BankKeeper.MintCoins(ctx, canonicalusdctypes.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(h.logical, usdc))))
		_, broken := canonicalusdckeeper.BackingInvariant(*h.k)(ctx)
		require.True(t, broken)
		h.k.EnforceBacking(ctx)
		breach := false
		for _, ev := range ctx.EventManager().Events() {
			if ev.Type == canonicalusdctypes.EventTypeBackingBreach {
				breach = true
			}
		}
		fmt.Printf("SCENARIO 1 uusdc minted outside the ledger -> %s breach_event=%v\n", h.ledger(ctx), breach)
		require.Equal(t, canonicalusdctypes.Mode_MODE_PAUSED, h.k.GetControls(ctx).Mode)
		require.True(t, breach)
		nbBefore := h.k.GetLedger(ctx).NobleBacking.BigInt()
		// Resuming re-snapshots: the stray supply is reclassified as Noble-backed.
		require.NoError(t, h.setControls(ctx, func(c *canonicalusdctypes.Controls) { c.Mode = canonicalusdctypes.Mode_MODE_GRADUAL }))
		fmt.Printf("SCENARIO resumed after the breach -> %s\n", h.ledger(ctx))
		require.Equal(t, 0, h.k.GetLedger(ctx).NobleBacking.BigInt().Cmp(new(big.Int).Add(nbBefore, big.NewInt(usdc))))
		h.invariantOK(ctx)
	})

	t.Run("legacy_injective_voucher_returns_as_ics20", func(t *testing.T) {
		ctx, _ := base.CacheContext()
		holder := h.findHolder(ctx, voucher, usdc)
		before := h.bal(ctx, holder, voucher)
		packet, err := h.send(ctx, holder, controls.InjectiveChannel, voucher, usdc, injReceiver, "")
		require.NoError(t, err)
		fmt.Printf("SCENARIO legacy voucher holder %s (%s) sent 1 USDC seq=%d | %s\n", holder, before, packet.Sequence, h.ledger(ctx))
		require.Equal(t, 0, h.bal(ctx, holder, voucher).Cmp(new(big.Int).Sub(before, big.NewInt(usdc))))
		require.Equal(t, uint32(0), h.k.PendingCount(ctx), "plain ICS20, no settlement")
		require.Equal(t, int64(0), h.k.GetLedger(ctx).InjectiveBacking.BigInt().Int64())
		require.Equal(t, 0, h.supply(ctx).Cmp(supply0))
		h.invariantOK(ctx)
	})

	t.Run("other_channels_cannot_bypass_or_counterfeit", func(t *testing.T) {
		ctx, _ := base.CacheContext()
		var other channeltypes.IdentifiedChannel
		for _, ch := range h.app.IBCKeeper.ChannelKeeper.GetAllChannels(ctx) {
			if ch.PortId == transfertypes.PortID && ch.State == channeltypes.OPEN && ch.ChannelId != controls.NobleChannel && ch.ChannelId != controls.InjectiveChannel {
				other = ch
				break
			}
		}
		require.NotEmpty(t, other.ChannelId)
		_, err := h.send(ctx, whale, other.ChannelId, h.logical, 10*usdc, "cosmos1cml96vmptgw99syqrrz8az79xer2pcgp0wls6c", "")
		fmt.Printf("SCENARIO uusdc over %s (counterparty %s): %v\n", other.ChannelId, other.Counterparty.ChannelId, err)
		require.Error(t, err)
		ok, errMsg := h.recv(ctx, other.ChannelId, other.Counterparty.ChannelId, controls.InjectivePacketDenom, 1000*usdc, user.String())
		fmt.Printf("SCENARIO erc20 USDC arriving over %s: ok=%v ack=%q | %s\n", other.ChannelId, ok, errMsg, h.ledger(ctx))
		require.True(t, ok, errMsg)
		require.Equal(t, int64(0), h.bal(ctx, user, h.logical).Int64(), "no canonical USDC minted")
		require.Equal(t, int64(0), h.k.GetLedger(ctx).InjectiveBacking.BigInt().Int64())
		require.Equal(t, 0, h.supply(ctx).Cmp(supply0))
		h.invariantOK(ctx)
	})
}
