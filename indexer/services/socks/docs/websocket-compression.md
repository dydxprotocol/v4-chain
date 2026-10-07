# Websocket compression (permessage-deflate)

## Background

Outbound websocket data is most of the indexer's internet egress. Messages are repetitive JSON,
which compresses well, but socks has always run with the `ws` permessage-deflate extension
disabled. Enabling it trades egress for CPU and memory on socks, so it is staged behind a
percentage rollout and is off by default.

Compression only applies to clients that offer the extension in their handshake
(`Sec-WebSocket-Extensions: permessage-deflate`). Browsers and most Python and Node clients do;
many Go and Rust clients do not. How much egress it can save therefore depends on the client
mix, which the metrics below measure before anything is enabled.

## Metrics

Both are emitted at every setting, including with compression disabled, and are tagged with
`state`:

- `not_offered` — the client did not offer the extension and can never be compressed.
- `offered` — the client offered it, but the connection is outside the rollout.
- `compressed` — the extension was negotiated for the connection.

| Metric | Meaning |
| --- | --- |
| `socks.ws_compression.connections` | New connections, by state. |
| `socks.ws_compression.payload_bytes` | Outbound message bytes before compression, by state. |

The share of `payload_bytes` in `offered` plus `compressed` is the share of websocket egress that
compression can reduce. `payload_bytes` is the uncompressed size in every state, so it does not
show the compression ratio; read that from the service's network bytes sent as the rollout
ramps.

## Configuration

| Config | Default | Meaning |
| --- | --- | --- |
| `WS_COMPRESSION_ROLLOUT_PERCENT` | `0` | Percentage of offering connections that are compressed. `0` disables the extension. |
| `WS_COMPRESSION_THRESHOLD_BYTES` | `1024` | Messages smaller than this are sent uncompressed. |
| `WS_COMPRESSION_METRIC_INTERVAL_MS` | `10000` | How often `payload_bytes` is emitted. |

A connection outside the rollout has the client's offer removed before the handshake is
answered, so it gets exactly the handshake it would get with compression disabled. The decision
is made once per connection and is random, so a reconnecting client may land on either side.

## Rolling out

1. Deploy at `0` and read the `offered` share of `payload_bytes`. If it is small, stop here.
2. Raise `WS_COMPRESSION_ROLLOUT_PERCENT` in steps, watching on each socks task:
   - CPU. Compression runs off the event loop, but the incident described in
     `subscription-limit-abuse.md` saturated a task on socket writes alone.
   - Memory. Each compressed connection holds a compression context of a few hundred KB, and `ws`
     documents memory fragmentation under high concurrency.
   - `socks.message_time_since_received`, for delivery delay.
3. To roll back, set the percentage to `0` and redeploy.

Most messages may be under the default threshold. If the `compressed` share of `payload_bytes`
is high but network bytes sent barely move, lower `WS_COMPRESSION_THRESHOLD_BYTES`.
