/**
 * Deployment script compatible with Foundry
 * Uses ethers.js with Foundry artifacts instead of Hardhat
 * This is essentially the same as deploy-foundry.js but kept for compatibility
 */

const { ethers } = require("ethers");
require("dotenv").config();

async function main() {
    // Connect to network RPC
    const rpcUrl = process.env.RPC_URL || "https://sepolia.base.org";
    const privateKey = process.env.PRIVATE_KEY;
    
    if (!privateKey) {
        throw new Error("PRIVATE_KEY environment variable is required");
    }
    
    const provider = new ethers.JsonRpcProvider(rpcUrl);
    const deployer = new ethers.Wallet(privateKey, provider);
    
    console.log("Deploying contracts with the account:", deployer.address);
    console.log("Balance:", ethers.formatEther(await provider.getBalance(deployer.address)), "ETH");

    // Load contract ABIs (from out/ directory after forge build)
    const MockERC20Artifact = require("../out/MockERC20.sol/MockERC20.json");
    const CollateralVaultArtifact = require("../out/CollateralVault.sol/CollateralVault.json");
    const OracleArtifact = require("../out/Oracle.sol/Oracle.json");
    const PerpEngineArtifact = require("../out/PerpEngine.sol/PerpEngine.json");

    // 1. Deploy MockERC20 (if needed)
    let collateralTokenAddress = process.env.COLLATERAL_TOKEN_ADDRESS;
    let mockERC20;
    
    if (!collateralTokenAddress) {
        console.log("Deploying MockERC20...");
        const MockERC20Factory = new ethers.ContractFactory(
            MockERC20Artifact.abi,
            MockERC20Artifact.bytecode.object,
            deployer
        );
        mockERC20 = await MockERC20Factory.deploy("Mock USDC", "mUSDC");
        await mockERC20.waitForDeployment();
        collateralTokenAddress = await mockERC20.getAddress();
        console.log("MockERC20 deployed to:", collateralTokenAddress);
    } else {
        console.log("Using existing Collateral Token:", collateralTokenAddress);
    }

    // 2. Deploy CollateralVault
    console.log("Deploying CollateralVault...");
    const CollateralVaultFactory = new ethers.ContractFactory(
        CollateralVaultArtifact.abi,
        CollateralVaultArtifact.bytecode.object,
        deployer
    );
    const vault = await CollateralVaultFactory.deploy(collateralTokenAddress);
    await vault.waitForDeployment();
    const vaultAddress = await vault.getAddress();
    console.log("CollateralVault deployed to:", vaultAddress);

    // 3. Deploy Oracle
    console.log("Deploying Oracle...");
    const OracleFactory = new ethers.ContractFactory(
        OracleArtifact.abi,
        OracleArtifact.bytecode.object,
        deployer
    );
    const oracle = await OracleFactory.deploy();
    await oracle.waitForDeployment();
    const oracleAddress = await oracle.getAddress();
    console.log("Oracle deployed to:", oracleAddress);

    // 4. Deploy PerpEngine
    console.log("Deploying PerpEngine...");
    const PerpEngineFactory = new ethers.ContractFactory(
        PerpEngineArtifact.abi,
        PerpEngineArtifact.bytecode.object,
        deployer
    );
    const engine = await PerpEngineFactory.deploy(vaultAddress, oracleAddress);
    await engine.waitForDeployment();
    const engineAddress = await engine.getAddress();
    console.log("PerpEngine deployed to:", engineAddress);

    // 5. Wire dependencies
    console.log("Setting PerpEngine in Vault...");
    await vault.setPerpEngine(engineAddress);

    // 6. Setup Markets (BTC-USD, ETH-USD)
    const BTC_MARKET = ethers.encodeBytes32String("BTC-USD");
    const ETH_MARKET = ethers.encodeBytes32String("ETH-USD");

    console.log("Setting up markets...");

    // Set Prices
    await oracle.setPrice(BTC_MARKET, ethers.parseUnits("30000", 8));
    await oracle.setPrice(ETH_MARKET, ethers.parseUnits("2000", 8));

    // Add Markets to Engine
    await engine.addMarket(
        BTC_MARKET,
        ethers.parseUnits("0.1", 18),
        ethers.parseUnits("0.05", 18),
        ethers.parseUnits("0.001", 18)
    );
    await engine.addMarket(
        ETH_MARKET,
        ethers.parseUnits("0.1", 18),
        ethers.parseUnits("0.05", 18),
        ethers.parseUnits("0.001", 18)
    );

    // Enable Engine for standalone demo
    await engine.setEngineEnabled(true);

    console.log("Deployment Complete!");
    console.log("----------------------------------------------------");
    console.log("MockToken:       ", mockERC20 ? await mockERC20.getAddress() : collateralTokenAddress);
    console.log("CollateralVault: ", vaultAddress);
    console.log("Oracle:          ", oracleAddress);
    console.log("PerpEngine:      ", engineAddress);
    console.log("----------------------------------------------------");
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
