# Local Deployment Steps (Foundry/Anvil)

Complete step-by-step guide for deploying contracts locally.

## Prerequisites

1. **Foundry installed** (see main README)
2. **Node.js & pnpm** installed
3. **Dependencies installed**: `pnpm install`

## Step-by-Step Deployment

### Step 1: Start Anvil (Local Blockchain)

**First, check if Anvil is already running:**
```bash
cd v5-demo-contracts
pnpm run anvil:check
```

**If Anvil is already running**, you have two options:

**Option A: Use the existing Anvil instance** (Recommended)
- Skip starting Anvil and proceed to Step 2
- The existing instance will work fine

**Option B: Restart Anvil** (Fresh state)
```bash
pnpm run anvil:restart
```

**If Anvil is NOT running**, start it:
```bash
pnpm run anvil
```

**What this does:**
- Starts a local Ethereum node on `http://localhost:8545`
- Creates 10 test accounts with 10,000 ETH each
- Default account: `0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266`

**Keep this terminal open** - Anvil must be running for deployment.

**Expected output:**
```
Listening on 127.0.0.1:8545
Available Accounts
==================
(0) 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 (10000 ETH)
...
```

### Step 2: Compile Contracts

In a **different terminal window**, compile the contracts:

```bash
cd v5-demo-contracts
pnpm run compile
```

**What this does:**
- Compiles all Solidity contracts using `forge build`
- Generates artifacts in `out/` directory
- These artifacts are needed by the deployment script

**Expected output:**
```
Compiling 5 files with 0.8.20
Compiler run successful
```

### Step 3: Deploy Contracts

In the **same terminal** (where you compiled), run:

```bash
pnpm run deploy:local
```

**What this does:**
- Connects to Anvil at `http://localhost:8545`
- Uses the default Anvil account (no `.env` needed for local)
- Deploys contracts in order:
  1. MockERC20 (USDC token)
  2. CollateralVault
  3. Oracle
  4. PerpEngine
- Sets up BTC-USD and ETH-USD markets
- Enables the engine for standalone mode

**Expected output:**
```
Deploying contracts with the account: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Balance: 10000.0 ETH
Deploying MockERC20...
MockERC20 deployed to: 0x5FbDB2315678afecb367f032d93F642f64180aa3
Deploying CollateralVault...
CollateralVault deployed to: 0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
Deploying Oracle...
Oracle deployed to: 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0
Deploying PerpEngine...
PerpEngine deployed to: 0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9
Setting PerpEngine in Vault...
Setting up markets...
Deployment Complete!
----------------------------------------------------
MockToken:       0x5FbDB2315678afecb367f032d93F642f64180aa3
CollateralVault: 0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
Oracle:          0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0
PerpEngine:      0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9
----------------------------------------------------
```

## Troubleshooting

### Error: "Address already in use" or "port 8545 already in use"
**Problem**: Anvil is already running on port 8545.

**Solution**: 
1. **Check if Anvil is running**: `pnpm run anvil:check`
2. **Option A**: Use the existing instance (recommended) - just proceed with deployment
3. **Option B**: Stop and restart: `pnpm run anvil:restart`
4. **Option C**: Stop only: `pnpm run anvil:stop`

### Error: "Cannot connect to Anvil"
**Problem**: Anvil is not running or wrong port.

**Solution**: 
1. Check Anvil is running: `pnpm run anvil:check` or `lsof -i :8545`
2. Start Anvil: `pnpm run anvil`
3. Verify it's listening on `http://localhost:8545`

### Error: "Cannot find module '../out/...'"
**Problem**: Contracts not compiled.

**Solution**: Run `pnpm run compile` first.

### Error: "nonce too low"
**Problem**: Anvil state is stale from previous deployment.

**Solution**: 
1. Stop Anvil (Ctrl+C)
2. Restart Anvil: `pnpm run anvil`
3. Deploy again: `pnpm run deploy:local`

### Error: "insufficient funds"
**Problem**: Account doesn't have ETH (shouldn't happen with Anvil defaults).

**Solution**: Restart Anvil to reset accounts.

## What Happens After Deployment?

The contracts are now deployed and ready to use:

- **MockERC20**: Test USDC token (mintable)
- **CollateralVault**: Manages user collateral
- **Oracle**: Price feed (settable by owner)
- **PerpEngine**: Main trading engine (enabled for standalone mode)

You can now:
- Interact with contracts via scripts
- Connect the Gateway service
- Use the React UI (point it to local Anvil)

## Next Steps

- **Run tests**: `pnpm test`
- **Deploy to Base Sepolia**: `pnpm run deploy:base` (requires `.env` with RPC_URL and PRIVATE_KEY)
- **Start Gateway**: See `v5-demo-gateway` README

