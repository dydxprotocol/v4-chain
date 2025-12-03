# PerpEngine - Developer Guide

A minimal perpetuals trading engine on Base Sepolia. Use the Gateway API to send batch settlements and listen to events.

## Contract Addresses

**Base Sepolia (Chain ID: 84532):**
- PerpEngine: `0x3DBA9Ab9d23664Da7bFEfACc0017d55576c5e4D3`

**Markets:**
- BTC-USD: `0x4254432d55534400000000000000000000000000000000000000000000000000`
- ETH-USD: `0x4554482d55534400000000000000000000000000000000000000000000000000`

**Other Contracts (for reference):**
- CollateralVault: `0xe0f16C0B95098CAD7a8bC7Db966C5464E1a24bEA`
- Oracle: `0xC00209d9Aaf599526e4D0c2edfCbFB385c62Cd5c`
- USDC: `0x036CbD53842c5426634e7929541eC2318f3dCF7e`

**Contract ABI:**
- Verified on [Basescan](https://sepolia.basescan.org/address/0x3DBA9Ab9d23664Da7bFEfACc0017d55576c5e4D3#code)
- Or use: `out/PerpEngine.sol/PerpEngine.json` (after `pnpm run compile`)

## Quick Start

### 1. Send Batch Settlement

**Gateway Endpoint:** `POST http://localhost:8080/settle-batch`

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
  "txHash": "0x...",
  "blockNumber": 12345678
}
```

**Example (JavaScript):**
```javascript
const response = await fetch('http://localhost:8080/settle-batch', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    settlements: [{
      marketId: "BTC-USD",
      userId: "user-1",
      evmAddress: "0xaa9826d668747e60cfe5a57c81dee920bfb95061",
      balanceDelta: "-10000",      // -0.01 USDC (6 decimals)
      sizeDelta: "10000000000000000", // +0.01 BTC (18 decimals)
      reason: "trade"
    }]
  })
});
const result = await response.json();
console.log('Tx Hash:', result.txHash);
```

### 2. Listen to Events

**Event:** `BalanceSettled(bytes32 indexed marketId, address indexed user, int256 balanceDelta, int256 sizeDelta)`

**Example (JavaScript with ethers.js):**
```javascript
const { ethers } = require('ethers');

const provider = new ethers.JsonRpcProvider('https://sepolia.base.org');
const engineAddress = '0x3DBA9Ab9d23664Da7bFEfACc0017d55576c5e4D3';
const engineABI = require('./out/PerpEngine.sol/PerpEngine.json').abi;
const engine = new ethers.Contract(engineAddress, engineABI, provider);

// Listen for events
engine.on('BalanceSettled', (marketId, user, balanceDelta, sizeDelta, event) => {
  console.log('Event:', {
    marketId,
    user,
    balanceDelta: balanceDelta.toString(),
    sizeDelta: sizeDelta.toString(),
    txHash: event.transactionHash,
    blockNumber: event.blockNumber
  });
});
```

**Or parse from transaction receipt:**
```javascript
const receipt = await provider.waitForTransaction(txHash);
const iface = new ethers.Interface(engineABI);
for (const log of receipt.logs) {
  const parsed = iface.parseLog(log);
  if (parsed && parsed.name === 'BalanceSettled') {
    console.log('BalanceSettled:', parsed.args);
  }
}
```

## API Reference

### POST /settle-batch

Settles multiple trades in a single transaction.

**Request Fields:**
- `marketId`: Market identifier (e.g., "BTC-USD")
- `userId`: Off-chain user ID (for logging)
- `evmAddress`: User's Ethereum address
- `balanceDelta`: Balance change in USDC (6 decimals, signed string)
- `sizeDelta`: Position size change (18 decimals, signed string)
- `reason`: Reason for settlement (e.g., "trade", "funding")

**Response:**
- `txHash`: Transaction hash
- `blockNumber`: Block number where transaction was mined
- `error`: Error message (if failed)

## Contract Functions

**Key Function:**
- `settleBatch(Settlement[] calldata settlements)`: Batch settle multiple users

**View Functions:**
- `getPosition(bytes32 marketId, address user)`: Get user position
- `getMargin(address user)`: Get user margin info
- `getUnrealizedPnl(bytes32 marketId, address user)`: Get unrealized PnL

**Contract Mode:**
- `engineEnabled`: `false` (ledger mode - no risk checks)
- Contract trusts operator's `settleBatch()` calls

## Setup

**Prerequisites:**
- Node.js (v18+)
- pnpm
- Gateway service running (see Gateway README)

**Install:**
```bash
pnpm install
pnpm run compile
```

## Examples

See `scripts/demo-batch-order.js` for a complete example.

## Links

- Contract on Basescan: https://sepolia.basescan.org/address/0x3DBA9Ab9d23664Da7bFEfACc0017d55576c5e4D3
- Gateway API: See `../v5-demo-gateway/README.md`
