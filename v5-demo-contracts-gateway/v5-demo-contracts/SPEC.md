# On-chain Components Spec

## Overview
A minimalist Perpetual Engine supporting **Multi-Market** trading with **Single Collateral** (USDC). It operates in two modes:
1.  **Ledger Mode**: Off-chain system pushes updates via trusted `settle()` calls.
2.  **Standalone Mode**: Users trade directly against the contract with on-chain risk checks.

## Contracts

### `PerpEngine.sol`
The core logic contract.

**Key Functions:**
*   `settle(bytes32 marketId, address user, int256 balanceDelta, int256 sizeDelta)`
    *   **Role**: `onlyOwner` (Gateway)
    *   **Action**: Updates user balance and position size blindly. Emits `BalanceSettled`.
*   `openPosition(bytes32 marketId, int256 sizeDelta, uint256 maxPrice, uint256 deadline)`
    *   **Role**: Public (requires `engineEnabled = true`)
    *   **Action**: Opens/increases position with slippage and margin checks.
*   `closePosition(bytes32 marketId, uint256 maxSlippage, uint256 deadline)`
    *   **Role**: Public (requires `engineEnabled = true`)
    *   **Action**: Closes entire position at Oracle price.
*   `liquidate(bytes32 marketId, address user)`
    *   **Role**: Public
    *   **Action**: Liquidates unsafe positions based on Maintenance Margin.

**Storage:**
*   `markets[marketId]`: Config (Margin Ratios, Fees).
*   `positions[marketId][user]`: Struct `{ size, entryPrice }`.
*   `vault`: Reference to `CollateralVault`.
*   `oracle`: Reference to `Oracle`.

### `Oracle.sol`
Simple centralized oracle for demo purposes.
*   `prices[marketId]`: Current price (8 decimals).
*   `setPrice(bytes32 marketId, int256 price)`: Updates price.

### `CollateralVault.sol`
Holds all system collateral (USDC).
*   `deposit(uint256 amount)`: User deposits USDC.
*   `withdraw(uint256 amount)`: User withdraws USDC.
*   `modifyBalance(address user, int256 amount)`: Called by Engine to credit/debit PnL.

## Events
*   `BalanceSettled(bytes32 marketId, address user, int256 balanceDelta, int256 sizeDelta)`
*   `PositionChanged(bytes32 marketId, address user, int256 newSize, int256 entryPrice, int256 realizedPnL)`
*   `Liquidated(bytes32 marketId, address user, address liquidator, int256 penalty)`
*   `MarketAdded(bytes32 marketId)`

## Deployment
*   **Network**: Base Sepolia
*   **Markets**: BTC-USD, ETH-USD
*   **Collateral**: MockUSDC (18 decimals)
