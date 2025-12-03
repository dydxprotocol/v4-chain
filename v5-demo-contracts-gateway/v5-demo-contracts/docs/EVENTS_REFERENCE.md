# v5 Demo Events Reference

This document lists all events emitted by the smart contracts. Indexers and frontends can subscribe to these events to track state changes.

## 1. PerpEngine Events
**Address**: `0xFeBd8490C545466DC4b0d566c33ff02beDF6f303`

### `BalanceSettled`
Emitted when a user's balance or position is updated via `settle` or `settleBatch`. This is the primary event for tracking ledger updates from the Gateway.

**Signature**
```solidity
event BalanceSettled(
    bytes32 indexed marketId, 
    address indexed user, 
    int256 balanceDelta, 
    int256 sizeDelta
);
```

**Parameters**
- `marketId`: The market identifier.
- `user`: The user's address.
- `balanceDelta`: Change in USDC balance (6 decimals).
- `sizeDelta`: Change in position size (18 decimals).

---

### `PositionChanged`
Emitted when a position is modified via direct engine actions (e.g., `openPosition`, `closePosition`, `liquidate`).
*Note: In Ledger Mode (Gateway), `BalanceSettled` is more common, but `PositionChanged` is used for internal logic updates.*

**Signature**
```solidity
event PositionChanged(
    bytes32 indexed marketId, 
    address indexed user, 
    int256 newSize, 
    int256 entryPrice, 
    int256 realizedPnL
);
```

**Parameters**
- `marketId`: The market identifier.
- `user`: The user's address.
- `newSize`: The updated total position size.
- `entryPrice`: The updated average entry price.
- `realizedPnL`: PnL realized during this update.

---

### `Liquidated`
Emitted when a user is liquidated.

**Signature**
```solidity
event Liquidated(
    bytes32 indexed marketId, 
    address indexed user, 
    address indexed liquidator, 
    int256 penaltyPaid
);
```

**Parameters**
- `marketId`: The market identifier.
- `user`: The liquidated user.
- `liquidator`: The address that performed the liquidation.
- `penaltyPaid`: The amount of collateral taken as penalty.

---

### `MarketAdded`
Emitted when a new market is initialized.

**Signature**
```solidity
event MarketAdded(bytes32 indexed marketId);
```

---

## 2. CollateralVault Events
**Address**: `0xE33B69eEFaaDC717242791a57f548358b64D01F9`

### `Deposit`
Emitted when a user deposits collateral.

**Signature**
```solidity
event Deposit(address indexed user, uint256 amount);
```

### `Withdraw`
Emitted when a user withdraws collateral.

**Signature**
```solidity
event Withdraw(address indexed user, uint256 amount);
```
