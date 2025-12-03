# Mini Go Gateway Service for PerpEngine Ledger Mode (v2)

This service is a **small Go HTTP service** that plays a gateway-like role for the demo: it receives balance and size deltas from v3 (or a wrapper service) over REST and pushes them on-chain via `settle` in PerpEngine ledger mode.

## Goal

- Keep v3 (or a small sidecar) as the source of truth for trades, PnL, funding, and liquidations.
- Use PerpEngine + CollateralVault as an **on-chain balance + position ledger** for a small set of demo users.
- Make the gateway minimal, simple, and performant:
  - Single small Go binary.
  - One main REST endpoint.
  - Uses `go-ethereum` to send EVM transactions.

## Responsibilities

The gateway service must:

1. **Expose a REST endpoint to receive settlements**

   Single endpoint (JSON over HTTP):

   ```http
   POST /settle
   Content-Type: application/json

   {
     "userId": "v3-user-id-or-alias",
     "evmAddress": "0x...",
     "balanceDelta": "123456",   // signed string, USDC/collateral delta
     "sizeDelta": "1000",        // signed string, position size delta (same units as PerpEngine)
     "reason": "trade|funding|liq",
     "reference": "fillId-or-batchId"
   }
   ```

   - `userId` is for logging only (v3 identity).
   - `evmAddress` is the mapped demo wallet for that user.
   - `balanceDelta` is the net collateral delta (PnL + fees + funding):
     - > 0: user gains collateral (profit/deposit)
     - < 0: user loses collateral (loss/withdraw)
   - `sizeDelta` is the change in position size:
     - > 0: user buys / increases long
     - < 0: user sells / closes (for this demo we only support longs, so negative = closing/reducing)
   - `reason` / `reference` are optional metadata fields for debugging.

   This endpoint is called by an existing v3 service or a thin wrapper once balances/positions are finalized.

2. **Send on-chain `settle` calls**

   The PerpEngine contract exposes:

   ```solidity
   function settle(address user, int256 balanceDelta, int256 sizeDelta) external;
   ```

   For each valid `POST /settle` request, the gateway:

   - Builds and signs a transaction calling:

     ```text
     settle(evmAddress, balanceDelta, sizeDelta)
     ```

     on the configured PerpEngine contract.

   - Uses environment variables / config for:
     - `RPC_URL`
     - `PRIVATE_KEY` (hot wallet)
     - `CHAIN_ID`
     - `PERPENGINE_ADDRESS`

   - Waits for the tx to be mined (or a small number of confirmations).

   - Returns a JSON response:

     ```json
     {
       "txHash": "0x...",
       "blockNumber": 1234567
     }
     ```

3. **Logging and minimal safety**

   - Log each settlement with:
     - `userId`, `evmAddress`, `balanceDelta`, `sizeDelta`, `reason`, `reference`, `txHash`, `blockNumber`.
   - Basic validation:
     - `evmAddress` is non-zero, correct length.
     - `balanceDelta` and `sizeDelta` parse as signed integers.
   - Optional: simple authentication (API key header) for demo security.

## Advanced behavior (configurable)

The following behavior SHOULD be implemented but must be controlled by env flags so it can be enabled/disabled as needed:

- `GATEWAY_MANAGE_NONCES` (default: `true`)
  - If `true`: manage nonces internally (e.g., mutex around sign+send, use pending nonce, handle “replacement transaction underpriced” / nonce conflicts).
  - If `false`: rely on the node / external nonce manager; use the node’s suggested nonce with no extra coordination.

- `GATEWAY_STRICT_ONCHAIN_ERRORS` (default: `true`)
  - If `true`: if the tx fails (revert / out-of-gas), return a 4xx/5xx error to the caller with a short error message so the v3 side knows the settlement failed.
  - If `false`: always return HTTP 200 with a payload that includes an `"error"` field when the tx fails; actual retry/alerting is delegated to another system.

- `GATEWAY_GAS_CHECK` (default: `false`)
  - If `true`: before sending, check the hot wallet’s ETH balance on Base Sepolia and reject if insufficient for gas.
  - If `false`: skip the explicit balance check and let the RPC / chain fail naturally.

## Non-goals (for this demo)

- No Kafka / event stream consumption.
- No DB or stateful reconciliation.
- No order matching or risk logic (all done in v3).
- No complex retry/backoff logic beyond basic tx submission (unless enabled via flags above).

**Summary:** this service is a thin, performant Go adapter: `REST -> settle(user, balanceDelta, sizeDelta)` on PerpEngine, so the demo can show v3-driven balances and positions being reflected on-chain with minimal new code.
