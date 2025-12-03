# Gateway Service

REST API service that bridges off-chain systems to the on-chain PerpEngine contract. Receives batch settlements and submits them to the blockchain.

## Quick Start

**Build:**
```bash
go build -o gateway
```

**Run:**
```bash
./gateway
```

Gateway starts on `http://localhost:8080` (default).

## Configuration

**Auto-Discovery (Recommended):**
Gateway automatically discovers config from `../v5-demo-contracts/.env` and `../v5-demo-contracts/deployment.json`.

**Manual Config:**
Set environment variables:
```bash
export OPERATOR_PRIVATE_KEY=0x...
export PERPENGINE_ADDRESS=0x3DBA9Ab9d23664Da7bFEfACc0017d55576c5e4D3
export COLLATERAL_VAULT_ADDRESS=0xe0f16C0B95098CAD7a8bC7Db966C5464E1a24bEA
export RPC_URL=https://sepolia.base.org
export GATEWAY_PORT=8080
```

## API Endpoints

### POST /settle-batch

Settles multiple trades in a single transaction.

**Request:**
```json
{
  "settlements": [
    {
      "marketId": "BTC-USD",
      "userId": "user-1",
      "evmAddress": "0xaa9826d668747e60cfe5a57c81dee920bfb95061",
      "balanceDelta": "-10000",
      "sizeDelta": "10000000000000000",
      "reason": "trade"
    }
  ]
}
```

**Response:**
```json
{
  "txHash": "0x49f787543b6fd9b42e50db6f3f531459d58ca7f8fb93a3e54fe212b82c26c15b",
  "blockNumber": 34491971
}
```

**Error Response:**
```json
{
  "error": "Insufficient balance"
}
```

### POST /settle

Settles a single trade (legacy endpoint).

**Request:**
```json
{
  "marketId": "BTC-USD",
  "userId": "user-1",
  "evmAddress": "0xaa9826d668747e60cfe5a57c81dee920bfb95061",
  "balanceDelta": "-10000",
  "sizeDelta": "10000000000000000",
  "reason": "trade"
}
```

### GET /user-state

Get user's on-chain state.

**Request:**
```
GET /user-state?evmAddress=0xaa9826d668747e60cfe5a57c81dee920bfb95061&marketId=BTC-USD
```

**Response:**
```json
{
  "balance": "980000",
  "position": {
    "size": "30000000000000000",
    "entryPrice": "0"
  }
}
```

## ChainID Detection

Gateway automatically detects ChainID from RPC URL:
- `https://sepolia.base.org` → 84532 (Base Sepolia)
- `https://sepolia-rollup.arbitrum.io/rpc` → 421614 (Arbitrum Sepolia)
- `http://localhost:8545` → 31337 (Anvil)

Override via `CHAIN_ID` environment variable if needed.

## Operator

Gateway uses the operator's private key to sign transactions. The operator must be the contract owner (deployed contracts).

**Current Operator (Base Sepolia):**
- Address: `0xc1b62e7cb5ddbc88b651d8b920329dbac22485b2`

## Events

After settlement, the contract emits `BalanceSettled` event:
```solidity
event BalanceSettled(
  bytes32 indexed marketId,
  address indexed user,
  int256 balanceDelta,
  int256 sizeDelta
);
```

Listen to events using your Ethereum client (ethers.js, web3.js, etc.).

## Examples

See `../v5-demo-contracts/scripts/demo-batch-order.js` for a complete example.
