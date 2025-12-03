/**
 * User Approve Script
 * Approves vault to spend USDC
 * Run this ONCE per chain after deployment
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
    
    console.log("User Approve");
    console.log("============");
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
    
    // Check USDC balance
    const usdcDecimals = await usdc.decimals();
    const usdcBalance = await usdc.balanceOf(userAddress);
    console.log(`USDC Balance: ${ethers.formatUnits(usdcBalance, usdcDecimals)} USDC`);
    
    if (usdcBalance === 0n) {
        throw new Error("User has no USDC. Please fund the user wallet first.");
    }
    
    // Amount to approve (1 USDC - allows 100 orders of 0.01 USDC each)
    const approveAmount = ethers.parseUnits("1", usdcDecimals);
    
    // Check current allowance
    const currentAllowance = await usdc.allowance(userAddress, vaultAddress);
    console.log(`Current Allowance: ${ethers.formatUnits(currentAllowance, usdcDecimals)} USDC`);
    
    // Approve if needed
    if (currentAllowance < approveAmount) {
        console.log(`Approving vault to spend ${ethers.formatUnits(approveAmount, usdcDecimals)} USDC...`);
        const approveTx = await usdc.approve(vaultAddress, approveAmount);
        console.log("Approve Tx:", approveTx.hash);
        await approveTx.wait();
        console.log("✅ Approval confirmed");
        
        // Verify
        const newAllowance = await usdc.allowance(userAddress, vaultAddress);
        console.log(`New Allowance: ${ethers.formatUnits(newAllowance, usdcDecimals)} USDC`);
    } else {
        console.log("✅ Already approved");
    }
    
    console.log("\n✅ Approval complete!");
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});

