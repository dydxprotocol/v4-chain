/**
 * Demo script compatible with Foundry/Anvil
 * Uses ethers.js with Foundry artifacts instead of Hardhat
 * Designed for local Anvil network
 */

const { ethers } = require("ethers");
require("dotenv").config();

async function main() {
    // Connect to Anvil (local) or network RPC
    const rpcUrl = process.env.RPC_URL || "http://localhost:8545";
    
    // For Anvil, use default accounts. For other networks, use PRIVATE_KEY
    const privateKey = process.env.PRIVATE_KEY || "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"; // Anvil default account 0
    
    const provider = new ethers.JsonRpcProvider(rpcUrl);
    
    // Create wallets from Anvil default accounts (for local testing)
    // Anvil provides 10 accounts, we'll use first 3
    const deployer = new ethers.Wallet(privateKey, provider);
    
    // Derive other accounts from Anvil defaults (for demo)
    // Account 1: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
    const userKey = process.env.USER_PRIVATE_KEY || "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d";
    const user = new ethers.Wallet(userKey, provider);
    
    // Account 2: 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC
    const liquidatorKey = process.env.LIQUIDATOR_PRIVATE_KEY || "0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a";
    const liquidator = new ethers.Wallet(liquidatorKey, provider);

    console.log("Running demo with user:", user.address);
    console.log("Liquidator:", liquidator.address);

    // Load contract artifacts (from out/ directory after forge build)
    const MockERC20Artifact = require("../out/MockERC20.sol/MockERC20.json");
    const CollateralVaultArtifact = require("../out/CollateralVault.sol/CollateralVault.json");
    const OracleArtifact = require("../out/Oracle.sol/Oracle.json");
    const PerpEngineArtifact = require("../out/PerpEngine.sol/PerpEngine.json");

    // --- Setup ---
    console.log("\n1. Setup & Deployment");
    
    // Deploy MockERC20
    const MockERC20Factory = new ethers.ContractFactory(
        MockERC20Artifact.abi,
        MockERC20Artifact.bytecode.object,
        deployer
    );
    const usdc = await MockERC20Factory.deploy("Mock USDC", "mUSDC");
    await usdc.waitForDeployment();
    const usdcAddr = await usdc.getAddress();

    // Deploy CollateralVault
    const CollateralVaultFactory = new ethers.ContractFactory(
        CollateralVaultArtifact.abi,
        CollateralVaultArtifact.bytecode.object,
        deployer
    );
    const vault = await CollateralVaultFactory.deploy(usdcAddr);
    await vault.waitForDeployment();
    const vaultAddr = await vault.getAddress();

    // Deploy Oracle
    const OracleFactory = new ethers.ContractFactory(
        OracleArtifact.abi,
        OracleArtifact.bytecode.object,
        deployer
    );
    const oracle = await OracleFactory.deploy();
    await oracle.waitForDeployment();
    const oracleAddr = await oracle.getAddress();

    // Deploy PerpEngine
    const PerpEngineFactory = new ethers.ContractFactory(
        PerpEngineArtifact.abi,
        PerpEngineArtifact.bytecode.object,
        deployer
    );
    const engine = await PerpEngineFactory.deploy(vaultAddr, oracleAddr);
    await engine.waitForDeployment();
    const engineAddr = await engine.getAddress();

    await vault.setPerpEngine(engineAddr);
    console.log("System deployed.");

    // --- Mint & Approve ---
    console.log("\n2. Mint & Approve Collateral");
    let tx = await usdc.mint(user.address, ethers.parseEther("10000"));
    console.log("Mint Tx:", tx.hash);
    await tx.wait();

    tx = await usdc.mint(liquidator.address, ethers.parseEther("10000"));
    await tx.wait();

    // Connect contracts to user/liquidator wallets
    const usdcUser = usdc.connect(user);
    const usdcLiquidator = usdc.connect(liquidator);
    const vaultUser = vault.connect(user);
    const engineUser = engine.connect(user);
    const engineLiquidator = engine.connect(liquidator);

    tx = await usdcUser.approve(vaultAddr, ethers.MaxUint256);
    console.log("Approve Tx:", tx.hash);
    await tx.wait();

    tx = await usdcLiquidator.approve(vaultAddr, ethers.MaxUint256);
    await tx.wait();
    console.log("User minted 10,000 USDC and approved Vault.");

    // --- Deposit ---
    console.log("\n3. Deposit Collateral");
    tx = await vaultUser.deposit(ethers.parseEther("1000"));
    console.log("Deposit Tx:", tx.hash);
    await tx.wait();
    console.log("User deposited 1,000 USDC.");
    console.log("User Vault Balance:", ethers.formatEther(await vault.balanceOf(user.address)));

    // --- Set Oracle ---
    console.log("\n4. Set Initial Price");
    const BTC_MARKET = ethers.encodeBytes32String("BTC-USD");
    const initialPrice = 30000n * 100000000n; // $30k
    tx = await oracle.setPrice(BTC_MARKET, initialPrice);
    console.log("Set Price Tx:", tx.hash);
    await tx.wait();
    console.log("Price set to $30,000");

    // --- Open Position ---
    console.log("\n5. Open Long Position (5x Leverage)");
    // User has 1000. 5x = 5000 notional.
    // Price = 30,000.
    // Size = 5000 / 30000 = 0.1666 BTC.
    // Let's do 0.1 BTC = 3000 notional (3x).
    const size = ethers.parseEther("0.1");
    const block = await provider.getBlock('latest');
    const deadline = block.timestamp + 60;
    const maxPrice = initialPrice + 100000000n; // Allow some slippage

    tx = await engineUser.openPosition(BTC_MARKET, size, maxPrice, deadline);
    console.log("Open Position Tx:", tx.hash);
    await tx.wait();
    console.log("Opened Long 0.1 BTC.");

    let pos = await engine.getPosition(BTC_MARKET, user.address);
    console.log("Position Size:", ethers.formatEther(pos.size));
    console.log("Entry Price:", Number(pos.entryPrice) / 1e8);

    // --- Price Up ---
    console.log("\n6. Price Moves Up (Profit)");
    const highPrice = 35000n * 100000000n; // $35k
    tx = await oracle.setPrice(BTC_MARKET, highPrice);
    await tx.wait();
    console.log("Price set to $35,000");

    let pnl = await engine.getUnrealizedPnl(BTC_MARKET, user.address);
    console.log("Unrealized PnL:", ethers.formatEther(pnl), "USDC");
    // Expected: 0.1 * (35000 - 30000) = 500 USDC.

    // --- Price Down (Liquidation) ---
    console.log("\n7. Price Moves Down (Liquidation)");
    // Entry 30k. Size 0.1. Notional 3000.
    // Equity ~1000.
    // Maintenance 5% of 3000 = 150.
    // Need Loss > 850.
    // 0.1 * (30000 - P) > 850
    // 3000 - 0.1P > 850
    // 2150 > 0.1P
    // P < 21,500.

    const crashPrice = 21000n * 100000000n; // $21k
    tx = await oracle.setPrice(BTC_MARKET, crashPrice);
    await tx.wait();
    console.log("Price set to $21,000");

    const margin = await engine.getMargin(BTC_MARKET, user.address);
    console.log("Margin Ratio:", ethers.formatEther(margin.marginRatio));
    
    // Get market config for maintenance ratio
    const marketConfig = await engine.markets(BTC_MARKET);
    console.log("Maintenance Ratio:", ethers.formatEther(marketConfig.maintenanceMarginRatio));

    // --- Liquidate ---
    console.log("\n8. Liquidate");
    tx = await engineLiquidator.liquidate(BTC_MARKET, user.address);
    console.log("Liquidate Tx:", tx.hash);
    await tx.wait();
    console.log("Liquidated!");

    pos = await engine.getPosition(BTC_MARKET, user.address);
    console.log("Position Size:", ethers.formatEther(pos.size)); // Should be 0

    const liqBalance = await vault.balanceOf(liquidator.address);
    console.log("Liquidator Balance (Reward):", ethers.formatEther(liqBalance));
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
