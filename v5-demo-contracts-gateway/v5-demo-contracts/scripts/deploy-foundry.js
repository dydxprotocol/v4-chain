/**
 * Deployment script compatible with Foundry/Anvil
 * Uses ethers.js with Anvil RPC instead of Hardhat Runtime Environment
 */

const { ethers } = require("ethers");
require("dotenv").config();

async function main() {
    // Connect to Anvil (local) or network RPC
    const rpcUrl = process.env.RPC_URL || "http://localhost:8545";
    // Use OPERATOR_PRIVATE_KEY for demo deployments
    const privateKey = process.env.OPERATOR_PRIVATE_KEY || process.env.PRIVATE_KEY || "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"; // Anvil default
    
    const provider = new ethers.JsonRpcProvider(rpcUrl);
    const deployer = new ethers.Wallet(privateKey, provider);
    
    console.log("Deploying contracts with the account:", deployer.address);
    console.log("Balance:", ethers.formatEther(await provider.getBalance(deployer.address)), "ETH");

    // Load contract ABIs (from out/ directory after forge build)
    const MockERC20Artifact = require("../out/MockERC20.sol/MockERC20.json");
    const CollateralVaultArtifact = require("../out/CollateralVault.sol/CollateralVault.json");
    const OracleArtifact = require("../out/Oracle.sol/Oracle.json");
    const PerpEngineArtifact = require("../out/PerpEngine.sol/PerpEngine.json");

    // 1. Use USDC address (real USDC for demo, or MockERC20 for local)
    let collateralTokenAddress = process.env.USDC_ADDRESS || process.env.COLLATERAL_TOKEN_ADDRESS;
    let mockERC20;
    
    if (!collateralTokenAddress) {
        console.log("Deploying MockERC20 (local/testing only)...");
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
        console.log("Using Collateral Token:", collateralTokenAddress);
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

    // Set Engine Mode (default: ledger mode for v3 integration)
    // Set ENGINE_ENABLED=true in .env for standalone mode (UI demo)
    const engineEnabled = process.env.ENGINE_ENABLED === "true";
    await engine.setEngineEnabled(engineEnabled);
    console.log(`Engine Mode: ${engineEnabled ? "Standalone (enabled)" : "Ledger (disabled - for v3 integration)"}`);

    // 7. Export config (if on testnet/mainnet)
    const network = await provider.getNetwork();
    const chainId = network.chainId.toString();
    
    // Save to deployment.json for Gateway/UI
    const fs = require("fs");
    const path = require("path");
    const deploymentPath = path.resolve(__dirname, "../deployment.json");
    let deploymentConfig = {};
    
    if (fs.existsSync(deploymentPath)) {
        try {
            deploymentConfig = JSON.parse(fs.readFileSync(deploymentPath));
        } catch (e) {
            console.log("Could not parse existing deployment.json, starting fresh.");
        }
    }
    
    // Update config for this chain
    deploymentConfig[chainId] = {
        name: chainId === "84532" ? "Base Sepolia" : (chainId === "421614" ? "Arbitrum Sepolia" : `Chain ${chainId}`),
        addresses: {
            mockToken: mockERC20 ? await mockERC20.getAddress() : collateralTokenAddress,
            vault: vaultAddress,
            oracle: oracleAddress,
            engine: engineAddress
        },
        markets: [
            { id: BTC_MARKET, name: "BTC-USD" },
            { id: ETH_MARKET, name: "ETH-USD" }
        ]
    };
    
    fs.writeFileSync(deploymentPath, JSON.stringify(deploymentConfig, null, 2));
    console.log(`Config saved to deployment.json for Chain ID ${chainId}`);

    console.log("\nDeployment Complete!");
    console.log("----------------------------------------------------");
    console.log("Chain ID:        ", chainId);
    console.log("Collateral Token:", collateralTokenAddress);
    console.log("CollateralVault: ", vaultAddress);
    console.log("Oracle:          ", oracleAddress);
    console.log("PerpEngine:      ", engineAddress);
    console.log("Engine Mode:     ", engineEnabled ? "Standalone" : "Ledger");
    console.log("----------------------------------------------------");
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});

