# Deployment & Setup Guide

This guide explains how to deploy the **PerpEngine Demo** (Contracts + Gateway) from scratch.

## 1. Prerequisites

### Software
*   **Foundry**: Install via `foundryup` (see [Foundry Book](https://book.getfoundry.sh/getting-started/installation))
*   **Node.js** (v18+) & `pnpm`
*   **Go** (v1.21+)
*   **Git**

### Wallets & Funds (Base Sepolia)
You need **1 Private Key** (The "Admin").
*   **Address**: `0x...` (Your Public Address)
*   **Funds**: ~0.1 ETH on **Base Sepolia** (for gas).
    *   *Get funds from [Coinbase Faucet](https://portal.cdp.coinbase.com/products/faucet) or [Alchemy Faucet](https://www.alchemy.com/faucets/base-sepolia).*

## 2. Roles & Architecture

In a production system, we separate concerns. In this demo, we simplify:

| Role | Responsibility | Production Setup | Demo Setup |
| :--- | :--- | :--- | :--- |
| **Deployer** | Deploys contracts. | Cold Wallet (Multisig) | **Admin Wallet** |
| **Operator** | Signs settlement txs. | Hot Wallet (Server) | **Admin Wallet** |
| **User** | Trades/Tests. | User Wallet | **Admin Wallet** (Self-trade) |

> **Why?** The contract restricts `settle()` to `onlyOwner`. Therefore, the Gateway (Operator) must hold the Owner's private key.

## 3. Step-by-Step Setup

### Step 1: Configure On-Chain (`v5-demo-contracts`)

1.  **Navigate to folder**:
    ```bash
    cd v5-demo-contracts
    ```
2.  **Install Dependencies**:
    ```bash
    pnpm install
    ```
3.  **Set Environment Variables**:
    Create a `.env` file:
    ```bash
    # .env
    PRIVATE_KEY="YOUR_PRIVATE_KEY_WITHOUT_0x"
    RPC_URL="https://sepolia.base.org"
    ```
4.  **Deploy Contracts**:
    This script deploys `CollateralVault`, `Oracle`, `PerpEngine`, and sets up markets (BTC/ETH).
    
    **Deploy with Foundry:**
    ```bash
    # Option 1: Use deploy script
    RPC_URL=https://sepolia.base.org node scripts/deploy.js
    
    # Option 2: Use deploy-foundry script (same functionality)
    pnpm run deploy:base
    ```
    
    *   **Output**: Note the `PerpEngine` address (e.g., `0xFeBd...`).
    *   *Note*: The script automatically updates `ui/src/config.json`.
    
    **For local deployment**, see [DEPLOYMENT_STEPS.md](DEPLOYMENT_STEPS.md).

### Step 2: Configure Gateway (`v5-demo-gateway`)

1.  **Navigate to folder**:
    ```bash
    cd ../v5-demo-gateway
    ```
2.  **Set Environment Variables**:
    The Gateway needs the same Private Key to sign settlements.
    ```bash
    export PRIVATE_KEY="YOUR_PRIVATE_KEY_WITHOUT_0x"
    export PERPENGINE_ADDRESS="ADDRESS_FROM_STEP_1" 
    # Or rely on auto-discovery from ../v5-demo-ui/src/config.json
    ```
3.  **Run Gateway**:
    ```bash
    go run main.go
    ```
    *   *Success*: You should see `Gateway starting on port 8080...`

### Step 3: Verify System (Professional Demo)

1.  **Open the Web UI**:
    ```bash
    cd ../v5-demo-ui
    pnpm run dev
    ```
    Open `http://localhost:5173`.

2.  **Full Lifecycle Test**:
    *   **Connect Wallet**: Click "Connect Wallet".
    *   **Mint**: Click "💰 Mint USDC" to get test tokens.
    *   **Deposit**: Click "📥 Deposit 1k".
        *   *Check*: "Vault Balance" should increase.
    *   **Trade**: Click "🔄 Simulate v3 Trade".
        *   *Check*: "Position" should update.
        *   *Check*: "Gateway View" (green box) should appear and match the data.
    *   **Withdraw**: Click "📤 Withdraw 500".
        *   *Check*: "Vault Balance" should decrease.

3.  **CLI Verification (Optional)**:
    ```bash
    cd ../v5-demo-contracts
    pnpm exec hardhat run scripts/verify-batch.js --network baseSepolia
    ```

## 4. Troubleshooting

*   **"tx reverted on chain"**:
    *   Did you redeploy contracts but forget to restart the Gateway? The Gateway caches the old address. **Restart it.**
    *   Does the wallet have ETH?
*   **"State mismatch"**:
    *   RPC nodes can be slow. The scripts have a 5s delay, but if it fails, try running `./scripts/query-state.sh` manually.

## 5. Multi-Chain Support (Base & Arbitrum)

### 1. Deployment
To deploy to a specific chain, use the `--network` flag. The script will automatically update `ui/src/config.json` for that chain without overwriting others.

*   **Base Sepolia**:
    ```bash
    RPC_URL=https://sepolia.base.org node scripts/export-ui-config.js
    ```
*   **Arbitrum Sepolia**:
    ```bash
    RPC_URL=https://sepolia-rollup.arbitrum.io/rpc node scripts/export-ui-config.js
    ```

### 2. Running the Gateway
Run a separate Gateway instance for each chain you want to support.

*   **Base Gateway**:
    ```bash
    export GATEWAY_PORT=8080
    export RPC_URL="https://sepolia.base.org"
    export PERPENGINE_ADDRESS="<BASE_ENGINE_ADDRESS>"
    go run main.go
    ```
*   **Arbitrum Gateway**:
    ```bash
    export GATEWAY_PORT=8081
    export RPC_URL="https://sepolia-rollup.arbitrum.io/rpc"
    export PERPENGINE_ADDRESS="<ARB_ENGINE_ADDRESS>"
    go run main.go
    ```

### 3. UI Switching
The UI automatically detects the connected wallet's network.
*   Switch your wallet (MetaMask) to **Base Sepolia** -> UI loads Base config.
*   Switch to **Arbitrum Sepolia** -> UI loads Arbitrum config.
*   If the network is not supported, the UI will show an error message.
