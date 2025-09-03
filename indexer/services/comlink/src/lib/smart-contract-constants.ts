import {
  getEntryPoint,
} from '@zerodev/sdk/constants';
import {
  base, arbitrum, avalanche, mainnet, optimism,
} from 'viem/chains';

export const entryPoint = getEntryPoint('0.7');

export const usdcAddressByChainId: Record<string, string> = {
  [mainnet.id.toString()]: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48', // usdc on ethereum mainnet.
  [arbitrum.id.toString()]: '0xaf88d065e77c8cC2239327C5EDb3A432268e5831', // usdc on arbitrum.
  [avalanche.id.toString()]: '0xb97ef9ef8734c71904d8002f8b6bc66dd9c48a6e', // usdc on avalanche.
  [base.id.toString()]: '0x833589fcd6edb6e08f4c7c32d4f71b54bda02913', // usdc on base.
  [optimism.id.toString()]: '0x0b2c639c533813f4aa9d7837caf62653d097ff85', // usdc on optimism.
  solana: 'EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v', // usdc on solana.
  'dydx-mainnet-1': 'ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5', // usdc on dydx.
};

export const ethDenomByChainId: Record<string, string> = {
  [mainnet.id.toString()]: 'ethereum-native', // eth on ethereum mainnet.
  [arbitrum.id.toString()]: 'arbitrum-native', // eth on arbitrum.
  [base.id.toString()]: 'base-native', // eth on base.
  [optimism.id.toString()]: 'optimism-native', // eth on optimism.
};

export const ARBITRUM_GO_FAST_HANDLER_SMART_CONTRACT = '0x4c58aE019E54D10594F1Aa26ABF385B6fb17A52d';
