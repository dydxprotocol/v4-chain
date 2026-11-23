import {
  getEntryPoint,
} from '@zerodev/sdk/constants';
import {
  base, arbitrum, avalanche, mainnet, optimism,
} from 'viem/chains';

export const entryPoint = getEntryPoint('0.7');

export const dydxChainId = 'dydx-mainnet-1';

export const usdcAddressByChainId: Record<string, string> = {
  [mainnet.id.toString()]: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48', // usdc on ethereum mainnet.
  [arbitrum.id.toString()]: '0xaf88d065e77c8cC2239327C5EDb3A432268e5831', // usdc on arbitrum.
  [avalanche.id.toString()]: '0xB97EF9Ef8734C71904D8002F8b6Bc66Dd9c48a6E', // usdc on avalanche.
  [base.id.toString()]: '0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913', // usdc on base.
  [optimism.id.toString()]: '0x0b2C639c533813f4Aa9D7837CAf62653d097Ff85', // usdc on optimism.
  solana: 'EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v', // usdc on solana.
  [dydxChainId]: 'ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5', // usdc on dydx.
};

export const ethDenomByChainId: Record<string, string> = {
  [mainnet.id.toString()]: 'ethereum-native', // eth on ethereum mainnet.
  [arbitrum.id.toString()]: 'arbitrum-native', // eth on arbitrum.
  [base.id.toString()]: 'base-native', // eth on base.
  [optimism.id.toString()]: 'optimism-native', // eth on optimism.
};

// used to convert the amount to the actual usdc amount.
export const SOLANA_USDC_QUANTUM = 1_000_000;
export const ETH_USDC_QUANTUM = 1_000_000;
export const ETH_WEI_QUANTUM = 1_000_000_000_000_000_000;

export const ARBITRUM_GO_FAST_HANDLER_SMART_CONTRACT = '0x4c58aE019E54D10594F1Aa26ABF385B6fb17A52d';

export const nobleChainId = 'noble-1';
export const osmosisChainId = 'osmosis-1';
export const neutronChainId = 'neutron-1';

// Deposits are bridged over CCTP to USDC on Injective, and the user sends them on to dYdX from
// there. Skip addresses Injective by its EVM chain id when it is the destination of a CCTP
// transfer; `injective-1` and `1776` are the same chain.
export const injectiveChainId = 'injective-1';
export const injectiveEvmChainId = '1776';
export const INJECTIVE_EVM_USDC_ADDRESS = '0xa00C59fF5a080D2b954d0c75e46E22a0c371235a';
// The CCTP domain of Injective.
export const INJECTIVE_CCTP_DOMAIN = 29;
// Skip's CCTP v2 relayer contract. It has the same address on every supported EVM chain.
export const SKIP_CCTP_V2_RELAYER_ADDRESS = '0x400BB58033a7763A834199190B68F66A2661aE73';
// The most USDC, in its smallest unit, that an EVM bridge call may leave with Skip's relayer as
// its fee: 5 USDC, which has 6 decimals on every supported chain. The fee is a few hundred
// units; this leaves room for it to rise without letting a call name the whole amount as fee.
export const EVM_BRIDGE_MAX_RELAY_FEE = 5_000_000;

// The CCTP v2 token messenger program on Solana, and its `deposit_for_burn` instruction.
export const SOLANA_CCTP_V2_PROGRAM = 'CCTPV2vPZJS2u2BBsUoscuikbYjnpFmbFsvVuJdgUMQe';
export const SOLANA_CCTP_V2_DEPOSIT_FOR_BURN_DISCRIMINATOR = 'd73c3d2e723780b0';
export const SOLANA_TOKEN_PROGRAM = 'TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA';
// The most USDC, in its smallest unit, that a Solana bridge transaction may pay Skip's relayer.
// The fee is a few hundred units; this leaves room for it to rise without allowing a transfer
// of any size.
export const SOLANA_BRIDGE_MAX_RELAY_FEE = 1_000_000;
