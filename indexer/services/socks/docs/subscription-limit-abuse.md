# Subscription-limit abuse: dropping connections and the Cloudflare contract

## Background

On 2026-08-26 between roughly 05:15Z and 05:25Z a single client held ~693 `v4_orderbook` and
~706 `v4_trades` subscriptions against one socks instance (10.0.1.170) while continuously
re-sending `subscribe` requests. CPU saturated on socket writes and JSON serialization, Kafka
consumption fell behind, and every other client that happened to land on that instance saw stale
market data. Trading was unaffected; comlink was unaffected.

## Why the existing limits did not stop it

Two separate controls existed, and the incident slipped between them.

1. `CHANNEL_CONNECTION_LIMITS` (`src/lib/subscription.ts`) caps subscriptions per channel per
   connection — 32 for `v4_orderbook` and `v4_trades`. When a connection is over the cap,
   `subscribe` sent an error and returned. **The connection stayed open and nothing was
   charged for the attempt**, so a retry loop was free.
2. `subscribeRateLimiter` did close connections, but it is keyed on `channel + subscriptionId`.
   A client cycling through hundreds of *distinct* markets gets a fresh allowance per market and
   never trips it.

Critically, the limit branch `return`ed *before* the rate limiter was ever consulted. So the
exact pattern in the incident — many distinct markets, all over the cap, retried in a loop — was
the one case that neither control covered.

## What changed

### 1. Repeated limit hits now drop the connection

`Subscriptions` gained `subscriptionLimitReachedRateLimiter`, keyed on a single per-connection
constant (`subscriptionLimitReached`) rather than on `channel + id`. The allowance is therefore
spent across every channel and market a connection touches, which is what makes the incident
pattern trip it.

Being at the limit stays a normal, recoverable condition: the first
`RATE_LIMIT_SUBSCRIPTION_LIMIT_REACHED_POINTS` hits within
`RATE_LIMIT_SUBSCRIPTION_LIMIT_REACHED_DURATION_MS` are answered with the usual error and the
connection is left alone. Past that, the connection is closed with close code
**4008** (`WS_CLOSE_CODE_SUBSCRIPTION_LIMIT_ABUSE`) and a body of:

```json
{ "message": "Subscription limit abuse", "reason": "subscription-limit-abuse", "retryAfterSeconds": 30 }
```

A distinct close code matters: 1008 is already used for the generic rate-limit drops, and we
want dashboards and clients to be able to tell the two apart.

### 2. The client is refused at the handshake for a cooldown

A websocket close frame is invisible to Cloudflare — it proxies the tunnelled connection but does
not inspect its frames, so "this client was dropped" is not something an edge rate-limiting rule
can ever match on. The handshake, however, is ordinary HTTP.

So a dropped client is put in a short in-memory penalty box (`src/lib/reconnect-penalty.ts`)
keyed by client address. While the penalty is live, `verifyClientNotPenalized` refuses the
websocket upgrade with:

```
HTTP/1.1 429 Too Many Requests
Retry-After: 30
x-dydx-ratelimit-reason: subscription-limit-abuse
```

That gives three things at once: a well-behaved client learns exactly how long to wait without
needing any Cloudflare involvement; the offending client cannot immediately re-establish and
resume hammering; and Cloudflare gets a countable, matchable signal.

The penalty is per-instance and in-memory by design. Its job is to stop one instance being
monopolised and to emit the edge signal. Cross-instance enforcement belongs at the edge, below.

### 3. Client address derivation

`getClientIp` (`src/helpers/header-utils.ts`) prefers `cf-connecting-ip`, then the left-most
`x-forwarded-for` entry, then the socket address. This matches comlink's `getIpAddr`.

This ordering is load-bearing rather than cosmetic. `indexer.dydx.trade` is Cloudflare-proxied
(its A records resolve into published Cloudflare ranges), so the ALB — and ALB access logs — only
ever see Cloudflare edge addresses. `cf-connecting-ip` is the only header that carries the real
client. This is also why the incident could not be attributed from S3 ALB logs.

## Configuration

| Variable | Default | Meaning |
|---|---|---|
| `RATE_LIMIT_SUBSCRIPTION_LIMIT_REACHED_POINTS` | `5` | Limit hits allowed per connection before it is dropped |
| `RATE_LIMIT_SUBSCRIPTION_LIMIT_REACHED_DURATION_MS` | `10000` | Window over which those hits are counted |
| `RECONNECT_PENALTY_ENABLED` | `true` | Kill switch for handshake refusal only. The connection drop stays active. |
| `RECONNECT_PENALTY_MS` | `30000` | How long a dropped client's upgrades are refused |
| `RECONNECT_PENALTY_MAX_TRACKED_CLIENTS` | `10000` | Bound on the penalty map |

The existing global `RATE_LIMIT_ENABLED` kill switch is honoured inside `RateLimiter` itself, so
setting it to `false` disables the drop as well.

### Note on shared addresses

The penalty is applied per client address. Clients behind a shared NAT or corporate egress can
therefore be caught by another party's behaviour. The 30s default and the requirement that a
single *connection* first exceed its allowance keep this narrow, and
`RECONNECT_PENALTY_ENABLED=false` disables the handshake refusal alone while leaving the
connection drop — the part that actually protects the instance — in place. Start with the
penalty enabled and watch `socks.connection_rejected_penalty`; a rate much higher than
`socks.subscriptions_limit_abuse_disconnect` would suggest shared-address collateral.

## Metrics

| Metric | Tags | Meaning |
|---|---|---|
| `socks.subscriptions_limit_reached` | `channel`, `instance` | Pre-existing. A limit hit. Normal at low rates. |
| `socks.subscriptions_limit_abuse_disconnect` | `channel`, `instance` | A connection was dropped for repeated hits |
| `socks.reconnect_penalty_applied` | `reason`, `instance` | A client entered the penalty box |
| `socks.connection_rejected_penalty` | `instance` | An upgrade was refused with 429 |

Drops also appear on the existing `socks.num_disconnects` counter tagged `code: "4008"`.

Worth alerting on: a sustained non-zero `socks.subscriptions_limit_abuse_disconnect` means a
client is in a retry loop and not backing off, which is the signal that the edge rule below
should be escalating.

## Cloudflare configuration

There is no Cloudflare terraform in `v4-infrastructure`, so these are console changes.

The origin does the per-instance protection; Cloudflare's job is to stop the reconnect storm from
reaching origin at all, and to make the timeout long enough that a client has to notice.

**Rate limiting rule — "socks subscription-limit abuse"**

- Match: `(http.host eq "indexer.dydx.trade" and starts_with(http.request.uri.path, "/v4/ws"))`
- Counting characteristic: client IP
- Count only responses where: `origin response status eq 429` **or**, if response-header matching
  is available on the plan, response header `x-dydx-ratelimit-reason` eq `subscription-limit-abuse`.
  The header form is the more precise of the two — it will not count 429s emitted for any other
  reason.
- Threshold: 5 requests per 60 seconds
- Action: **Block**, duration **300 seconds** (deliberately an order of magnitude longer than the
  origin's 30s penalty — the origin cooldown is a nudge, the edge block is the consequence)
- Custom response: status `429`, content type `application/json`, body:

```json
{
  "error": "subscription_limit_abuse",
  "message": "Your connection was closed because it repeatedly exceeded the per-connection subscription limit and continued retrying. Unsubscribe before subscribing again, or distribute subscriptions across additional connections. Reconnections are blocked for 5 minutes.",
  "docs": "https://docs.dydx.exchange/api_integration-indexer/indexer_websocket"
}
```

Confirm the actual websocket path before applying — the match above assumes `/v4/ws`.

**Verifying it**

1. Confirm `socks.subscriptions_limit_abuse_disconnect` is non-zero for the offender.
2. Confirm the 429s reach the edge — Cloudflare Analytics, filtered to the socks path, grouped by
   origin status.
3. Confirm the rule fires and that the custom body is what the client actually receives.

**Attribution**, for the next incident: because the ALB only sees Cloudflare addresses, identify
the offender from Cloudflare's own analytics/logs, or from the socks-side logs which now carry
the `cf-connecting-ip`-derived address. S3 ALB access logs will not tell you who it was.
