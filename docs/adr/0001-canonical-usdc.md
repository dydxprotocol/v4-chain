# 1. Canonical USDC

* Status: Accepted
* Date: 2026-09-15
* Deciders: protocol
* Tags: ibc, usdc, x/canonicalusdc

## Context

dYdX Chain trading, subaccounts, deposits, and the indexer all use one USDC denom:

`uusdc` = `ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5`

That hash is Noble-origin USDC (`transfer/<noble-channel>/uusdc`). The chain needs the same user-facing asset to be backed by Injective-origin USDC instead, without a denom migration.

ICS20 is asynchronous. Any design that moves backing across chains has to account for in-flight packets (success ack, error ack, timeout) so the ledger cannot double-spend or double-refund.

## Decision

Keep `uusdc` as the only application USDC. Add `x/canonicalusdc` to classify what backs that supply and to rewrite IBC on two bound channels.

The goal of `GRADUAL` is to loop Injective-in / Noble-out until `noble_backing` is zero, then turn Noble withdrawals off so users can only exit via Injective.

An upcoming coordinated upgrade adds the module in `DISABLED` (passthrough). Gov `MsgUpdateControls` later enables `GRADUAL`. On that transition the keeper snapshots:

```
noble_backing = bank.GetSupply(uusdc) - injective_backing - pending_injective
```

Do not snapshot in the upgrade handler. Supply still moves while disabled.

### Modes

| Mode | Behavior |
| --- | --- |
| `DISABLED` | Identity. Existing transfer stack unchanged. |
| `GRADUAL` | Canonical routing. Noble deposits of `uusdc` rejected. Injective deposits mint `uusdc`. |
| `PAUSED` | Canonical `uusdc` flow fails. Pending settlements still complete. Other IBC unchanged. |

`MsgUpdateControls` is the only module message. It binds Noble/Injective channel, client, connection, packet denoms, `max_transfer_amount`, `migration_ceiling`, `max_pending_settlements`, and the Noble withdrawal flag.

`noble_packet_denom` and `injective_packet_denom` are the counterparty's base denom as it arrives in an inbound packet (e.g. `uusdc`), not the full trace. Outbound packets carry the full trace (`transfer/<channel>/<base>`), because ICS20 writes the voucher's full denom path into the packet. Every comparison against an outbound packet uses the trace form.

Snapshot happens on every transition into `GRADUAL` from `DISABLED` or `PAUSED`.

### Ledger

```
supply(uusdc) == noble_backing + injective_backing + pending_injective
```

`pending_noble` is intentionally not in that sum. The two pending fields are different liabilities.

| Field | Meaning |
| --- | --- |
| `noble_backing` | Still Noble-classified. Spent when a Noble `MsgTransfer` is accepted. |
| `injective_backing` | Physical Injective USDC on the module. Spent on Injective `MsgTransfer`. |
| `pending_injective` | In-flight Injective withdraw. Logical `uusdc` is still on the module (locked), so it remains in `supply` and in the equation. |
| `pending_noble` | In-flight Noble withdraw. The user's `uusdc` already left via ICS20, so `supply` dropped with `noble_backing`. The field is only a ticket for ack/timeout, not backing. |

All amounts are `bytes` with `(gogoproto.customtype) = "github.com/dydxprotocol/v4-chain/protocol/dtypes.SerializableInt"`.

IBC stack: `canonicalusdc` → `ratelimit` → `transfer`.

**Inbound (middleware)**

- Injective channel, matching packet denom: physical tokens to the module, mint `uusdc` to the original receiver, `injective_backing += amount`. Cap: `max_transfer_amount`, `migration_ceiling`.
- Noble channel, `uusdc`: reject before the inner app (`ErrNobleDepositsDisabled`).
- Other channel, `uusdc`: passthrough, no ledger change. Before enablement, `uusdc` sent to another chain was escrowed (not burned), so it stayed in `supply` and is already inside the `noble_backing` snapshot. The return unescrows it and `supply` does not move. Escrow bounds what can return, and outbound on other channels is rejected after enablement, so no new escrow is created.
- `DISABLED`: passthrough, including Noble deposits.

**Outbound (`MsgTransfer` decorator)**

Users still send `ibc.applications.transfer.v1.MsgTransfer` of `uusdc`. Both routes reject unless `amount <= max_transfer_amount` and `PendingCount < max_pending_settlements` (`max_pending_settlements` is a count cap, `<= 100_000`).

- Injective channel: require `amount <= injective_backing`. Charge the `uusdc` withdrawal rate limit for `amount`. Lock user `uusdc` on the module (`pending_injective +=`, `injective_backing -=`), send physical Injective USDC. Success ack burns the lock and clears pending. Timeout/error: inner transfer refunds physical to the module, unlock `uusdc` to the user, `injective_backing +=`, undo the `uusdc` rate-limit charge, pending cleared.
- Noble channel: require `NobleWithdrawalsEnabled` and `amount <= noble_backing`. `noble_backing -=`, `pending_noble +=`, forward the original packet (user `uusdc` leaves, `supply` drops). Success ack clears pending only. Timeout/error: inner transfer refunds `uusdc` to the user (`supply` rises), `noble_backing +=`, pending cleared.
- Other channel of `uusdc`: `ErrUnsupportedChannel`.
- Other denoms: passthrough.

`subtract` fails if backing would go negative (`ErrInsufficientBacking`).

**Rate limiting**

`x/ratelimit` keys capacity on the packet denom. The Injective route's packet carries the physical Injective denom, so the decorator charges `uusdc` capacity on that route and undoes it on failed ack or timeout. The Noble route packet is `uusdc` and is charged by `x/ratelimit` directly. Both routes draw from one `uusdc` bucket.

**Settlement**

IBC core delivers exactly one ack or timeout per packet. Ack and timeout match the pending entry by `(source_channel, sequence)`, then check sender, receiver, amount, and the full-trace denom against the packet.

Settlement checks only what it touches: the pending entry exists, route fields stay nonnegative, and on an Injective refund the module holds the physical coins. The global supply equation is not evaluated during settlement, so drift never blocks in-flight withdrawals. It is reported through `x/crisis` and telemetry.

**ICS4**

Packets that are logical Noble USDC (the Noble denom trace, on any channel) require the decorator's authorization key. This is also what rejects `uusdc` outbound on other channels after enablement.

User-sent physical Injective USDC is allowed so holders can unwind it off-chain and deposit it back through the Injective inbound path (which mints `uusdc`). Module-sent Injective packets (canonical withdraw) still require authorization.

### Replacing Noble backing

While Noble withdrawals are still enabled:

1. Gov enables `GRADUAL` (snapshot).
2. Receive Injective USDC in a batch `<= max_transfer_amount`. Receiver gets `uusdc`; `injective_backing` rises.
3. Send that same amount out Noble. `noble_backing` falls.
4. Repeat in batches until `noble_backing == 0`. The full amount does not need to be funded up front; the migrator only needs one batch of Injective USDC at a time.
5. Set `NobleWithdrawalsEnabled = false`. `MsgUpdateControls` rejects this unless `noble_backing == 0` and no Noble settlement is pending. Users then exit only via Injective.

The `uusdc` minted in step 2 must be the `uusdc` withdrawn in step 3. This is an operational requirement of the migrator, not enforced by the module.

### Invariants (when not `DISABLED`)

- `supply(uusdc) == noble_backing + injective_backing + pending_injective` (`pending_noble` excluded; those coins already left `supply`)
- Module physical Injective balance `>= injective_backing`
- Module locked `uusdc` `>= pending_injective` (Noble pending does not lock module `uusdc`)
- Backing fields stay nonnegative; outbound amount cannot exceed the selected route's backing
- Pending count and per-route sums match stored settlements (`pending_noble` is still counted here); count `<= max_pending_settlements`
- Cannot switch to `DISABLED` or change channel/client/connection/denom identity while pending exists
- `NobleWithdrawalsEnabled` cannot go `true` → `false` unless `noble_backing == 0` and no Noble settlement is pending

`MsgUpdateControls` enforces the transition rules (last two bullets). The supply equation is enforced only on transitions into `GRADUAL`, right after the snapshot, where it holds by construction. A transition into `PAUSED` is never blocked by the supply equation, so governance can always stop canonical flow.

### Queries

- `State`: controls and ledger.
- `PendingSettlements`: paginated (`cosmos.base.query.v1beta1.PageRequest`).
- `PendingSettlement`: one entry by `(source_channel, sequence)`.

## Considered options

**1. Second user-facing denom (Injective IBC hash)**

Users would hold a different `ibc/...`. CLOB, sending, subaccounts, and the indexer would all change. Rejected.

**2. Rename `uusdc` to the Injective hash**

Same blast radius as (1), plus every existing balance and in-flight packet. Rejected.

**3. Memo-funded backing swap**

An allowlisted controller sends Injective USDC to the module with a JSON memo (no mint), then `MsgExecuteBackingSwap` sends Noble USDC to their Noble address. Same ledger substitution as the 1:1 loop if they withdraw immediately, at the cost of extra messages, `restricted_funding`, a participant allowlist, and memo parsing. Rejected. The `MsgTransfer` loop is the mechanism.

**4. Snapshot `noble_backing` in the upgrade that adds the module**

Supply still changes while `DISABLED`. The number would be stale by enablement. Rejected. Snapshot in `MsgUpdateControls` on transition into `GRADUAL`.

**5. Separate rate-limit bucket for the Injective denom**

Governance would set `LimitParams` for the physical Injective denom. Two independent buckets roughly double the effective `uusdc` outflow cap, and the Injective bucket's baseline would track the module's physical balance rather than user supply. Rejected in favor of charging `uusdc` capacity on both routes.

Chosen: keep `uusdc`. Classify backing. Replace it with Injective-in / Noble-out until `noble_backing` is zero. Snapshot on enable.

## Consequences

Positive:

- Indexer, comlink, and existing `uusdc` balances stay on the same denom.
- The upgrade that adds the module is inert (`DISABLED`). Enablement is a later gov tx.
- Users keep using `MsgTransfer`. No new user message.

Negative:

- After enable, Noble deposits of `uusdc` stop. Physical Injective USDC already on dYdX does not convert in place; it has to leave and come back on the Injective channel.
- Enablement should happen with Noble IBC quiet. A Noble packet sent before the snapshot that times out afterward refunds `uusdc` without a matching `noble_backing` credit, so `supply` exceeds the ledger. Governance recovers with `GRADUAL` → `PAUSED` → `GRADUAL`, which re-snapshots `noble_backing`.
- Noble withdrawals cannot be turned off while `noble_backing > 0`, so the migration cannot be finalized until the loop drains it exactly. If the last residue cannot be drained (for example a Noble timeout refund re-adds backing), Noble withdrawals stay enabled until it can.

## Confirmation

`go test ./x/canonicalusdc/...` in `protocol/`. Implementation: `protocol/x/canonicalusdc`.
