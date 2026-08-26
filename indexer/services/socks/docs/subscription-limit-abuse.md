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

The close frame is only the start of the drop. `ws` keeps an uncooperative peer's socket alive
for up to 30 seconds waiting for a close frame that may never come, and goes on delivering that
peer's messages for the whole window. So the drop also, in order: unsubscribes the connection
immediately rather than waiting for the close event, so the forwarder stops serializing market
data for it; destroys the socket outright if `close()` itself throws; and arms a short
forced-termination fallback. Inbound messages are discarded without being parsed once the socket
is no longer `OPEN`, so a client that ignores the close frame gets no further work done on its
behalf.

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
| `SUBSCRIPTION_LIMIT_ABUSE_DROP_ENABLED` | `false` | Kill switch for the connection drop, independent of `RATE_LIMIT_ENABLED`. Ships off; enable once thresholds are sized. |
| `RATE_LIMIT_SUBSCRIPTION_LIMIT_REACHED_POINTS` | `5` | Limit hits allowed per connection before it is dropped |
| `RATE_LIMIT_SUBSCRIPTION_LIMIT_REACHED_DURATION_MS` | `10000` | Window over which those hits are counted |
| `SUBSCRIPTION_LIMIT_ABUSE_FORCE_TERMINATE_MS` | `1000` | Grace period after the close frame before the socket is destroyed |
| `RECONNECT_PENALTY_ENABLED` | `true` | Kill switch for handshake refusal only. The connection drop stays active. |
| `RECONNECT_PENALTY_MS` | `30000` | How long a dropped client's upgrades are refused |
| `RECONNECT_PENALTY_MAX_TRACKED_CLIENTS` | `100000` | Bound on the penalty map |

`SUBSCRIPTION_LIMIT_ABUSE_DROP_ENABLED` gates the drop on its own, so it can be rolled back
without touching the subscribe, ping and invalid-message limiters. The global `RATE_LIMIT_ENABLED`
switch still disables everything, including the drop. Disabling the drop also disables the
cooldown, since a client only ever enters it by being dropped.

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

| `socks.subscriptions_limit_abuse_force_terminate` | `reason`, `instance` | A dropped socket was destroyed without a completed close handshake |
| `socks.message_dropped_not_open` | `readyState`, `instance` | An inbound message was discarded because the socket was no longer open |

Drops also appear on the existing `socks.num_disconnects` counter tagged `code: "4008"`.

A steady `socks.subscriptions_limit_abuse_force_terminate` tagged `close_timeout` means clients
are ignoring close frames rather than disconnecting, which is itself a sign of a badly behaved
integration.

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
- **Counting expression — must repeat host and path:**

  ```
  http.host eq "indexer.dydx.trade"
  and starts_with(http.request.uri.path, "/v4/ws")
  and http.response.code eq 429
  and any(http.response.headers["x-dydx-ratelimit-reason"][*] eq "subscription-limit-abuse")
  ```

  A custom counting expression does **not** inherit the rule's matching expression. Cloudflare
  evaluates it independently, so host and path have to be restated or the counter will increment
  on unrelated traffic.

  The reason header is not optional hardening here, it is what makes the rule correct. The ALB
  routes exactly `/v4/ws` to socks and everything else under `/v4/*` to comlink, which runs its
  own rate limiter and returns its own 429s. A counting expression of "status 429" alone — or one
  scoped only by a `/v4/ws` prefix, which still covers `/v4/ws...` paths the ALB hands to comlink
  — would count comlink's rate-limit responses and block those clients out of the websocket
  entirely. socks is the only origin that emits this header.
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

**Prerequisite.** Counting on response fields requires a Business plan or above. Confirm the
zone's tier before building the rule. If response-based counting is unavailable, the fallback is
a request-rate rule on the `/v4/ws` path per client IP with no response condition — cruder, since
it cannot distinguish an abusive reconnect storm from a legitimate one, so it needs a threshold
set well above normal reconnect behaviour.

**Verifying it**

1. Confirm `socks.subscriptions_limit_abuse_disconnect` is non-zero for the offender.
2. Confirm the 429s reach the edge — Cloudflare Analytics, filtered to the socks path, grouped by
   origin status.
3. Confirm the rule fires and that the custom body is what the client actually receives.

**Attribution**, for the next incident: because the ALB only sees Cloudflare addresses, identify
the offender from Cloudflare's own analytics/logs, or from the socks-side logs which now carry
the `cf-connecting-ip`-derived address. S3 ALB access logs will not tell you who it was.
