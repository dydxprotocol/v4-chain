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

A subscribe request that is already awaiting its initial response from comlink when the drop
lands is invalidated rather than allowed to complete. Every connection carries an epoch that is
captured before the await and re-checked after it; teardown drops the epoch. Without that check
a late-resolving request would recreate the exact structures teardown had just deleted, and
because teardown has already run, nothing would ever clean them up — leaving a dead connection
being serialized for indefinitely. The same race exists on the ordinary disconnect path and is
covered by the same guard.

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

The fixed cooldown is deliberate. A self-clearing, rate-proportional penalty is generally
preferable — AWS WAF's rate-based rules and Cloudflare's `mitigation_timeout: 0` both work that
way — but it is not available at this layer: a rate cannot be measured for a connection that is
being refused before it is accepted. A hard stop is required, and Kraken makes the same call at
the same layer, banning a client IP for 10 minutes after too many connection attempts. 30s is
deliberately at the conservative end of the range implied by comparable systems.

The penalty is per-instance and in-memory by design. Its job is to stop one instance being
monopolised and to emit the edge signal. Cross-instance enforcement belongs at the edge, below.

### 3. Client address derivation

`getClientIp` (`src/helpers/header-utils.ts`) prefers `cf-connecting-ip`, then the left-most
`x-forwarded-for` entry, then the socket address. This matches comlink's `getIpAddr`.

This ordering is load-bearing rather than cosmetic. `indexer.dydx.trade` is Cloudflare-proxied
(its A records resolve into published Cloudflare ranges), so the ALB — and ALB access logs — only
ever see Cloudflare edge addresses. `cf-connecting-ip` is the only header that carries the real
client. This is also why the incident could not be attributed from S3 ALB logs.

## Threshold sizing: measured, and why the drop is still disabled

Measured on mainnet (dYdX Indexer org, ap1) for 2026-08-26 04:30-06:30Z. Note that the socks
metrics there have tag grouping disabled for cardinality control, so these are fleet-wide totals;
no per-channel or per-instance breakdown is available.

| Signal | Baseline | Peak in window |
|---|---|---|
| `socks.message_received_subscribe` | ~300k/min | 1.31M/min (05:22) |
| `socks.subscriptions_limit_reached` | ~200k/min | 635k/min (05:11) |
| Rejected share of all subscribes | ~45-65% | ~57% |
| `socks.num_concurrent_connections` | ~3,700 | ~4,270 |
| Implied per-connection rejections | ~54/min (~9 per 10s) | ~171/min (~29 per 10s) |

The per-connection row is a fleet average, and how it should be read depends entirely on how the
rejections are distributed -- which these metrics cannot show, because tag grouping is disabled
on them.

**The working assumption is that rejections are heavily concentrated.** A well-behaved client
subscribes once per market when it connects and then stays subscribed; in steady state it
produces no rejections at all. Producing them continuously requires continuously re-subscribing,
which is the pathological pattern by definition. The baseline is also flat rather than tracking
connection churn, which is what a bulk-subscribe-at-connect explanation would predict. On that
reading the ~200k/min comes from a small number of offenders, the typical connection sits at or
near zero, and an allowance of 5 per 10s does not endanger healthy clients.

The alternative reading -- that the load is spread evenly -- would mean the average connection
sits at roughly twice the allowance while nothing is wrong, and enabling the drop would
disconnect a large share of the customer base on a rolling basis. The two readings are the same
numbers and lead to opposite decisions, so the distribution has to be measured rather than
assumed.

**How to measure it before enabling anything.** The abuse allowance is evaluated even when
`SUBSCRIPTION_LIMIT_ABUSE_DROP_ENABLED` is off, and every connection that exceeds it is counted.
`socks.subscriptions_limit_abuse_connections` reports the number of *distinct* connections that
would have been dropped in each metrics interval, tagged `enforcing:false` while the drop is
disabled. Deploy with the drop off and compare that gauge against
`socks.num_concurrent_connections`:

- A handful of connections against several thousand confirms concentration. Enable the drop.
- A large fraction means the behaviour is normal for this system, and dropping connections is the
  wrong control -- raising the per-channel cap, making the refusal cacheable so clients stop
  retrying, or limiting subscribe *attempts* regardless of outcome would each address more of the
  load without disconnecting anyone.

This costs nothing to run and turns the decision into a reading rather than an argument.

## Configuration

| Variable | Default | Meaning |
|---|---|---|
| `SUBSCRIPTION_LIMIT_ABUSE_DROP_ENABLED` | `false` | Kill switch for the connection drop, independent of `RATE_LIMIT_ENABLED`. Ships off. Do not enable without first reading the threshold-sizing section above. |
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

Scoping abuse enforcement to a client address is a known-imperfect fallback, and both Cloudflare
and Google Cloud publish that caveat: a shared NAT or VPN egress pools unrelated users into one
bucket, and an attacker with an address range can rotate out of it. The accepted mitigation is to
key on an authenticated principal where one exists and fall back to address only where none does.
socks has no credential at the handshake, so address is the only key available at that point —
which is why the cooldown is deliberately short and confined to connection establishment rather
than applied to ongoing usage.

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
| `socks.subscription_abandoned_stale_connection` | `channel`, `instance` | An initial response resolved after its connection was torn down |
| `socks.subscriptions_limit_abuse_connections` | `enforcing`, `instance` | Distinct connections over the abuse allowance per interval, counted whether or not the drop is enabled |

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

- Match: `(http.host eq "indexer.dydx.trade" and http.request.uri.path eq "/v4/ws")`
- Counting characteristic: client IP
- **Counting expression — must repeat host and path:**

  ```
  http.host eq "indexer.dydx.trade"
  and http.request.uri.path eq "/v4/ws"
  and http.response.code eq 429
  and any(http.response.headers["x-dydx-ratelimit-reason"][*] eq "subscription-limit-abuse")
  ```

  A custom counting expression does **not** inherit the rule's matching expression. Cloudflare
  evaluates it independently, so host and path have to be restated or the counter will increment
  on unrelated traffic.

  Exact path equality, not a prefix. The ALB routes exactly `/v4/ws` to socks and everything else
  under `/v4/*` to comlink, which runs its own rate limiter and returns its own 429s. A prefix
  match would still cover `/v4/ws...` paths the ALB hands to comlink; equality mirrors the ALB
  rule exactly and keeps the blast radius to the one endpoint at issue.

  The reason header is the second, independent guard: socks is the only origin that emits it, so
  even a misconfigured path condition cannot make the counter fire on comlink's rate-limit
  responses. Counting on status alone would block those clients out of the websocket entirely.
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

The path is `/v4/ws`, matching `aws_lb_listener_rule.public_https_socks` in `v4-infrastructure`
(`indexer/load_balancer.tf`), which forwards exactly that path to the socks target group.

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
