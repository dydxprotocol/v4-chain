/**
 * Export UI config script compatible with Foundry
 * Uses ethers.js with Foundry artifacts (out/) instead of Hardhat
 */

const { ethers } = require("ethers");
const fs = require("fs");
const path = require("path");
require("dotenv").config();

async function main() {
    console.log("Deploying contracts for UI...");

    // Connect to network RPC
    const rpcUrl = process.env.RPC_URL || "https://sepolia.base.org";
    const privateKey = process.env.PRIVATE_KEY;
    
    if (!privateKey) {
        throw new Error("PRIVATE_KEY environment variable is required");
    }

    const provider = new ethers.JsonRpcProvider(rpcUrl);
    const deployer = new ethers.Wallet(privateKey, provider);

    console.log("Deploying with account:", deployer.address);
    console.log("Balance:", ethers.formatEther(await provider.getBalance(deployer.address)), "ETH");

    // Load contract artifacts (from out/ directory after forge build)
    const MockERC20Artifact = require("../out/MockERC20.sol/MockERC20.json");
    const CollateralVaultArtifact = require("../out/CollateralVault.sol/CollateralVault.json");
    const OracleArtifact = require("../out/Oracle.sol/Oracle.json");
    const PerpEngineArtifact = require("../out/PerpEngine.sol/PerpEngine.json");

    // 1. Deploy MockERC20
    console.log("Deploying MockERC20...");
    const MockERC20Factory = new ethers.ContractFactory(
        MockERC20Artifact.abi,
        MockERC20Artifact.bytecode.object,
        deployer
    );
    const usdc = await MockERC20Factory.deploy("Mock USDC", "mUSDC");
    await usdc.waitForDeployment();
    const usdcAddr = await usdc.getAddress();
    console.log("MockERC20:", usdcAddr);

    // 2. Deploy CollateralVault
    console.log("Deploying CollateralVault...");
    const CollateralVaultFactory = new ethers.ContractFactory(
        CollateralVaultArtifact.abi,
        CollateralVaultArtifact.bytecode.object,
        deployer
    );
    const vault = await CollateralVaultFactory.deploy(usdcAddr);
    await vault.waitForDeployment();
    const vaultAddr = await vault.getAddress();
    console.log("CollateralVault:", vaultAddr);

    // 3. Deploy Oracle
    console.log("Deploying Oracle...");
    const OracleFactory = new ethers.ContractFactory(
        OracleArtifact.abi,
        OracleArtifact.bytecode.object,
        deployer
    );
    const oracle = await OracleFactory.deploy();
    await oracle.waitForDeployment();
    const oracleAddr = await oracle.getAddress();
    console.log("Oracle:", oracleAddr);

    // 4. Deploy PerpEngine
    console.log("Deploying PerpEngine...");
    const PerpEngineFactory = new ethers.ContractFactory(
        PerpEngineArtifact.abi,
        PerpEngineArtifact.bytecode.object,
        deployer
    );
    const engine = await PerpEngineFactory.deploy(vaultAddr, oracleAddr);
    await engine.waitForDeployment();
    const engineAddr = await engine.getAddress();
    console.log("PerpEngine deployed to:", engineAddr);

    // 5. Wire dependencies
    await vault.setPerpEngine(engineAddr);

    // Enable Engine Mode for UI Demo
    await engine.setEngineEnabled(true);
    console.log("Engine Mode Enabled");

    // 6. Setup Markets
    const BTC_MARKET = ethers.encodeBytes32String("BTC-USD");
    const ETH_MARKET = ethers.encodeBytes32String("ETH-USD");

    console.log("Setting up markets...");
    await oracle.setPrice(BTC_MARKET, ethers.parseUnits("100000", 8)); // Start high for demo
    await oracle.setPrice(ETH_MARKET, ethers.parseUnits("4000", 8));

    await engine.addMarket(BTC_MARKET, ethers.parseUnits("0.1", 18), ethers.parseUnits("0.05", 18), ethers.parseUnits("0.001", 18));
    await engine.addMarket(ETH_MARKET, ethers.parseUnits("0.1", 18), ethers.parseUnits("0.05", 18), ethers.parseUnits("0.001", 18));

    await engine.setEngineEnabled(true);

    // 7. Export Config
    const network = await provider.getNetwork();
    const chainId = network.chainId.toString();
    console.log(`Detected Chain ID: ${chainId}`);

    const uiConfigPath = path.resolve(__dirname, "../deployment.json");
    let existingConfig = {};
    if (fs.existsSync(uiConfigPath)) {
        try {
            existingConfig = JSON.parse(fs.readFileSync(uiConfigPath));
        } catch (e) {
            console.log("Could not parse existing config, starting fresh.");
        }
    }

    // Shared ABIs (from Foundry artifacts)
    const abis = {
        mockToken: MockERC20Artifact.abi,
        vault: CollateralVaultArtifact.abi,
        oracle: OracleArtifact.abi,
        engine: PerpEngineArtifact.abi
    };

    // Update specific chain config
    existingConfig[chainId] = {
        name: chainId === "84532" ? "Base Sepolia" : (chainId === "421614" ? "Arbitrum Sepolia" : "Localhost"),
        addresses: {
            mockToken: usdcAddr,
            vault: vaultAddr,
            oracle: oracleAddr,
            engine: engineAddr
        },
        markets: [
            { id: BTC_MARKET, name: "BTC-USD" },
            { id: ETH_MARKET, name: "ETH-USD" }
        ]
    };

    // Ensure ABIs are at top level
    existingConfig.abis = abis;

    fs.writeFileSync(uiConfigPath, JSON.stringify(existingConfig, null, 2));
    console.log(`Config updated for Chain ID ${chainId} at ${uiConfigPath}`);
    
    // Also copy to UI if path exists
    const uiConfigDest = path.resolve(__dirname, "../../v5-demo-ui/src/config.json");
    if (fs.existsSync(path.dirname(uiConfigDest))) {
        fs.writeFileSync(uiConfigDest, JSON.stringify(existingConfig, null, 2));
        console.log(`Config also copied to ${uiConfigDest}`);
    }
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
