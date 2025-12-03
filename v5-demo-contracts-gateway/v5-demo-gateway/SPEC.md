# Gateway Service Spec

## Overview
A minimalist Go service that bridges an off-chain matching engine (v3-style) to the on-chain Perpetual Engine (v5). It exposes a REST API to settle trades, funding, and liquidations on-chain.

## API

### `POST /settle`
Settle a balance and position update for a user.

**Request:**
```json
{
  "marketId": "BTC-USD",       // Market Identifier (e.g. "BTC-USD", "ETH-USD")
  "userId": "uuid-123",        // Off-chain User ID (for logging)
  "evmAddress": "0x123...",    // User's On-chain Address
  "balanceDelta": "500000",    // Amount to credit/debit (18 decimals, signed)
  "sizeDelta": "100000000",    // Position size change (18 decimals, signed)
  "reason": "trade",           // "trade" | "funding" | "liquidation"
  "reference": "fill-123"      // Off-chain reference ID
}
```

**Response:**
```json
{
  "txHash": "0xabc...",        // Ethereum Transaction Hash
  "blockNumber": 1234567,      // Block number where tx was mined
  "error": ""                  // Error message if failed
}
```

## Configuration
Configured via Environment Variables. Supports auto-discovery from sibling directories for zero-config dev.

| Variable | Default | Description |
| :--- | :--- | :--- |
| `PORT` | `8080` | Service port |
| `RPC_URL` | `https://sepolia.base.org` | Ethereum RPC Endpoint |
| `PRIVATE_KEY` | *(Auto-discovered)* | Operator Private Key (Hex) |
| `PERPENGINE_ADDRESS` | *(Auto-discovered)* | Deployed Contract Address |
| `GATEWAY_MANAGE_NONCES` | `true` | Handle nonce concurrency locally |
| `GATEWAY_STRICT_ONCHAIN_ERRORS` | `true` | Return 500 on tx revert |

## Behavior
1.  **Validation**: Checks EVM address format and parses BigInts.
2.  **Transaction**: Signs and sends a transaction to `PerpEngine.settle()`.
3.  **Confirmation**: Waits for transaction mining (simple polling).
4.  **Concurrency**: Uses a mutex for local nonce management to support concurrent requests.
