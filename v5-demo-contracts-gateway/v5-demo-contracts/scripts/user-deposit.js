/**
 * User Deposit Script
 * Deposits USDC to the vault
 * Run this ONCE per chain after approval
 */

const { ethers } = require("ethers");
require("dotenv").config();

async function main() {
    const userPrivateKey = process.env.USER_PRIVATE_KEY;
    if (!userPrivateKey) {
        throw new Error("USER_PRIVATE_KEY environment variable is required");
    }

    const rpcUrl = process.env.RPC_URL || "https://sepolia.base.org";
    const usdcAddress = process.env.USDC_ADDRESS || "0x036CbD53842c5426634e7929541eC2318f3dCF7e"; // Base Sepolia USDC
    
    // Load deployment config
    const fs = require("fs");
    const path = require("path");
    const deploymentPath = path.resolve(__dirname, "../deployment.json");
    
    if (!fs.existsSync(deploymentPath)) {
        throw new Error("deployment.json not found. Deploy contracts first.");
    }
    
    const deploymentConfig = JSON.parse(fs.readFileSync(deploymentPath));
    const provider = new ethers.JsonRpcProvider(rpcUrl);
    const network = await provider.getNetwork();
    const chainId = network.chainId.toString();
    
    const chainConfig = deploymentConfig[chainId];
    if (!chainConfig) {
        throw new Error(`No config found for Chain ID ${chainId}. Deploy contracts first.`);
    }
    
    const vaultAddress = chainConfig.addresses.vault;
    
    console.log("User Deposit");
    console.log("=============");
    console.log("Chain ID:", chainId);
    console.log("User Address:", new ethers.Wallet(userPrivateKey, provider).address);
    console.log("Vault Address:", vaultAddress);
    console.log("USDC Address:", usdcAddress);
    console.log("");

    const user = new ethers.Wallet(userPrivateKey, provider);
    const userAddress = user.address;
    
    // Load standard ERC20 ABI
    const erc20Abi = require("../contracts/interfaces/ERC20.json");
    const usdc = new ethers.Contract(usdcAddress, erc20Abi, user);
    
    // Load Vault ABI
    const CollateralVaultArtifact = require("../out/CollateralVault.sol/CollateralVault.json");
    const vault = new ethers.Contract(vaultAddress, CollateralVaultArtifact.abi, user);
    
    // Check USDC balance
    const usdcDecimals = await usdc.decimals();
    const usdcBalance = await usdc.balanceOf(userAddress);
    console.log(`USDC Balance: ${ethers.formatUnits(usdcBalance, usdcDecimals)} USDC`);
    
    if (usdcBalance === 0n) {
        throw new Error("User has no USDC. Please fund the user wallet first.");
    }
    
    // Amount to deposit (1 USDC - allows 100 orders of 0.01 USDC each)
    const depositAmount = ethers.parseUnits("1", usdcDecimals);
    
    if (usdcBalance < depositAmount) {
        throw new Error(`Insufficient USDC. Need ${ethers.formatUnits(depositAmount, usdcDecimals)} USDC, have ${ethers.formatUnits(usdcBalance, usdcDecimals)} USDC`);
    }
    
    // Check allowance
    const allowance = await usdc.allowance(userAddress, vaultAddress);
    if (allowance < depositAmount) {
        throw new Error(`Insufficient allowance. Need ${ethers.formatUnits(depositAmount, usdcDecimals)} USDC, have ${ethers.formatUnits(allowance, usdcDecimals)} USDC. Run user-approve.js first.`);
    }
    
    // Deposit
    console.log(`Depositing ${ethers.formatUnits(depositAmount, usdcDecimals)} USDC...`);
    const depositTx = await vault.deposit(depositAmount);
    console.log("Deposit Tx:", depositTx.hash);
    await depositTx.wait();
    console.log("✅ Deposit confirmed");
    
    // Verify
    const vaultBalance = await vault.balanceOf(userAddress);
    console.log(`Vault Balance: ${ethers.formatUnits(vaultBalance, 18)} USDC`);
    
    console.log("\n✅ Deposit complete!");
    console.log(`You can now place up to 100 orders of 0.01 USDC each.`);
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});

