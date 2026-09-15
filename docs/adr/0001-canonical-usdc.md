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

`MsgUpdateControls` is the only module message. It binds Noble/Injective channel, client, connection, packet denoms, `max_transfer_amount`, `migration_ceiling`, `max_pending_settlements`, and the Noble withdrawal flag/cutoff.

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
| `legacy_downstream` | Budget for `uusdc` returning from some other channel. No increment after genesis. |

IBC stack: `canonicalusdc` → `ratelimit` → `transfer`.

**Inbound (middleware)**

- Injective channel, matching packet denom: physical tokens to the module, mint `uusdc` to the original receiver, `injective_backing += amount`. Cap: `max_transfer_amount`, `migration_ceiling`.
- Noble channel, `uusdc`: reject before the inner app (`ErrNobleDepositsDisabled`).
- Other channel, `uusdc`: decrement `legacy_downstream` if it covers the amount; else reject.
- `DISABLED`: passthrough, including Noble deposits.

**Outbound (`MsgTransfer` decorator)**

Users still send `ibc.applications.transfer.v1.MsgTransfer` of `uusdc`.

- Injective channel: lock user `uusdc` on the module (`pending_injective +=`, `injective_backing -=`), send physical Injective USDC. Success ack burns the lock and clears pending. Timeout/error: inner transfer refunds physical to the module, unlock `uusdc` to the user, `injective_backing +=`, pending cleared.
- Noble channel: requires `NobleWithdrawalsEnabled` and cutoff. `noble_backing -=`, `pending_noble +=`, forward the original packet (user `uusdc` leaves, `supply` drops). Success ack clears pending only. Timeout/error: inner transfer refunds `uusdc` to the user (`supply` rises), `noble_backing +=`, pending cleared.
- Other channel of `uusdc`: `ErrUnsupportedChannel`.
- Other denoms: passthrough.

A completed ring (100k slots) makes ack/timeout replay a no-op.

**ICS4**

Packets that are logical Noble USDC (packet denom `uusdc` on the Noble channel, or the Noble denom trace) require the decorator's authorization key.

User-sent physical Injective USDC is allowed so holders can unwind it off-chain and deposit it back through the Injective inbound path (which mints `uusdc`). Module-sent Injective packets (canonical withdraw) still require authorization.

### Replacing Noble backing

While Noble withdrawals are still enabled:

1. Gov enables `GRADUAL` (snapshot).
2. Receive Injective USDC in chunks `<= max_transfer_amount`. Receiver gets `uusdc`; `injective_backing` rises.
3. Send that same amount out Noble. `noble_backing` falls.
4. Repeat until `noble_backing` is ~0.
5. Set `NobleWithdrawalsEnabled = false`. Users then exit only via Injective.

The `uusdc` minted in step 2 must be the `uusdc` withdrawn in step 3.

### Invariants (when not `DISABLED`)

- `supply(uusdc) == noble_backing + injective_backing + pending_injective` (`pending_noble` excluded; those coins already left `supply`)
- Module physical Injective balance `>= injective_backing`
- Module locked `uusdc` `>= pending_injective` (Noble pending does not lock module `uusdc`)
- Pending count and per-route sums match stored settlements (`pending_noble` is still counted here)
- Cannot switch to `DISABLED` or change channel/denom identity while pending exists

## Considered options

**1. Second user-facing denom (Injective IBC hash)**

Users would hold a different `ibc/...`. CLOB, sending, subaccounts, and the indexer would all change. Rejected.

**2. Rename `uusdc` to the Injective hash**

Same blast radius as (1), plus every existing balance and in-flight packet. Rejected.

**3. Memo-funded backing swap**

An allowlisted controller sends Injective USDC to the module with a JSON memo (no mint), then `MsgExecuteBackingSwap` sends Noble USDC to their Noble address. Same ledger substitution as the 1:1 loop if they withdraw immediately, at the cost of extra messages, `restricted_funding`, a participant allowlist, and memo parsing. Rejected. The `MsgTransfer` loop is the mechanism.

**4. Snapshot `noble_backing` in the upgrade that adds the module**

Supply still changes while `DISABLED`. The number would be stale by enablement. Rejected. Snapshot in `MsgUpdateControls` on `DISABLED` → not-disabled.

Chosen: keep `uusdc`. Classify backing. Replace it with Injective-in / Noble-out until `noble_backing` is zero. Snapshot on enable.

## Consequences

Positive:

- Indexer, comlink, and existing `uusdc` balances stay on the same denom.
- The upgrade that adds the module is inert (`DISABLED`). Enablement is a later gov tx.
- Users keep using `MsgTransfer`. No new user message.

Negative:

- After enable, Noble deposits of `uusdc` stop. Physical Injective USDC already on dYdX does not convert in place; it has to leave and come back on the Injective channel.
- `legacy_downstream` is an ops input. At 0, unclassified returns of `uusdc` from other chains fail.
- Enablement should happen with Noble IBC quiet. A timeout that refunds `uusdc` after the snapshot breaks the invariant until the ledger is fixed.
- Once `NobleWithdrawalsEnabled` is false, the 1:1 loop cannot take remaining Noble backing out. A leftover needs a later gov tool or is left as-is.

## Confirmation

`go test ./x/canonicalusdc/...` in `protocol/`. Implementation: `protocol/x/canonicalusdc`.
