import { logger } from '@dydxprotocol-indexer/base';
import { route, messages } from '@skip-go/client/cjs';
import type { Adapter } from '@solana/wallet-adapter-base';
import type { Transaction, VersionedTransaction } from '@solana/web3.js';
import { Turnkey } from '@turnkey/sdk-server';
import { TurnkeySigner } from '@turnkey/solana';
import { decode, fromWords } from 'bech32';
import bs58 from 'bs58';
import { encodeFunctionData, type Hex } from 'viem';
import type { SmartAccountImplementation } from 'viem/account-abstraction';

import config from '../config';
import {
  ethDenomByChainId, usdcAddressByChainId,
} from '../lib/smart-contract-constants';
import { getAddress } from './alchemy-helpers';

export const suborgToApproval = new Map<string, string>();

// Grabs the raw skip route data to carry out the bridge on our own.
export async function getSkipCallData(
  sourceAddress: string,
  sourceAssetDenom: string,
  dydxAddress: string,
  amount: string,
  chainId: string,
): Promise<Parameters<SmartAccountImplementation['encodeCalls']>[0]> {
  // support for hex amounts.
  let amountToUse = amount;
  if (amount.startsWith('0x')) {
    amountToUse = parseInt(amount, 16).toString();
  }
  const routeResult = await route({
    amountIn: amountToUse, // Desired amount in smallest denomination (e.g., uatom)
    sourceAssetDenom,
    sourceAssetChainId: chainId,
    destAssetDenom: usdcAddressByChainId['dydx-mainnet-1'],
    destAssetChainId: 'dydx-mainnet-1',
    cumulativeAffiliateFeeBps: '0',
    smartRelay: true, // skip recommended to enable for better routes and less faults.
    smartSwapOptions: {
      splitRoutes: true,
      evmSwaps: true, // needed for native eth bridging.
    },
    allowUnsafe: false,
    goFast: true,
  });
  logger.info({
    at: 'skip-helper#getSkipCallData',
    message: 'Route result obtained',
    routeResult,
  });
  if (!routeResult) {
    throw new Error('Failed to find a route');
  }

  logger.info({
    at: 'skip-helper#getSkipCallData',
    message: 'Route result obtained',
    routeResult,
    dydxAddress,
  });

  const userAddresses = await Promise.all(
    routeResult.requiredChainAddresses.map(async (cid: string) => ({
      chainId: cid,
      address: await getAddress(cid, sourceAddress, dydxAddress),
    })),
  );

  let addressList: string[] = [];
  userAddresses.forEach((userAddress: { chainId: string, address: string }, index: number) => {
    const requiredChainAddress = routeResult.requiredChainAddresses[index];

    if (requiredChainAddress === userAddress?.chainId) {
      addressList.push(userAddress.address);
    }
  });

  if (addressList.length !== routeResult.requiredChainAddresses.length) {
    addressList = userAddresses.map((x: { chainId: string, address: string }) => x.address);
  }

  const validLength = addressList.length === routeResult.requiredChainAddresses.length ||
    addressList.length === routeResult.chainIds?.length;

  if (!validLength) {
    throw new Error('executeRoute error: invalid address list');
  }

  const timeoutSeconds = '60'; // Set a timeout for the messages request
  const response = await messages({
    timeoutSeconds,
    amountIn: routeResult?.amountIn,
    amountOut: routeResult.estimatedAmountOut || '0',
    sourceAssetChainId: routeResult?.sourceAssetChainId,
    sourceAssetDenom: routeResult?.sourceAssetDenom,
    destAssetChainId: routeResult?.destAssetChainId,
    destAssetDenom: routeResult?.destAssetDenom,
    operations: routeResult?.operations,
    addressList,
    slippageTolerancePercent: '1',
  });

  let data = '';
  let toAddress = '';

  response?.msgs?.forEach((msg, index: number) => {
    if ('evmTx' in msg) {
      logger.info({
        at: 'skip-helper#getSkipCallData',
        message: `Message ${index + 1} EVM Transaction`,
        evmTx: msg.evmTx,
      });
      data = msg.evmTx.data || '';
      toAddress = msg.evmTx.to || '';
    }
  });

  // need value to be the amount if native asset.
  let value = BigInt(0);
  if (Object.values(ethDenomByChainId).map(
    (x) => x.toLowerCase(),
  ).includes(sourceAssetDenom.toLowerCase())) {
    value = BigInt(amount);
  }

  // this is the actual call data that is responsible for the bridge.
  const callData = [
    {
      to: (toAddress.startsWith('0x') ? toAddress : (`0x${toAddress}`)) as Hex,
      value,
      data: data.startsWith('0x') ? data as Hex : (`0x${data}`) as Hex, // "0x",
    },
  ];
  // if it's usdc, we need a swap to eth first, this approves a swap.
  if (Object.values(usdcAddressByChainId).map(
    (x) => x.toLowerCase(),
  ).includes(sourceAssetDenom.toLowerCase())) {
    callData.unshift({
      to: usdcAddressByChainId[chainId] as `0x${string}`,
      value: BigInt(0),
      data: encodeFunctionData({
        abi: [
          {
            name: 'approve',
            type: 'function',
            stateMutability: 'nonpayable',
            inputs: [
              { name: 'spender', type: 'address' },
              { name: 'amount', type: 'uint256' },
            ],
            outputs: [{ name: '', type: 'bool' }],
          },
        ],
        functionName: 'approve',
        args: [
          (toAddress.startsWith('0x') ? toAddress : (`0x${toAddress}`)) as Hex,
          BigInt(amount),
        ],
      }), // "0x",
    });
  }

  return callData;
}

export function getSvmSigner(suborgId: string, signWith: string) {
  const serverClient = new Turnkey({
    apiBaseUrl: config.TURNKEY_API_BASE_URL as string,
    apiPublicKey: config.TURNKEY_API_SENDER_PUBLIC_KEY as string,
    apiPrivateKey: config.TURNKEY_API_SENDER_PRIVATE_KEY as string,
    defaultOrganizationId: suborgId,
  });

  const turnkeySigner = new TurnkeySigner({
    organizationId: suborgId,
    client: serverClient.apiClient(),
  });

  // eslint-disable-next-line @typescript-eslint/require-await
  return async () => ({
    publicKey: {
      toString: () => signWith,
      toBase58: () => signWith,
    },
    signTransaction: async (tx: Transaction) => {
      try {
        return turnkeySigner.signTransaction(tx, signWith);
      } catch (error) {
        throw new Error(`Failed to sign transaction with TurnkeySigner: ${error.message}`);
      }
    },
    signAllTransactions: async (txs: (Transaction | VersionedTransaction)[]) => {
      try {
        return turnkeySigner.signAllTransactions(txs, signWith);
      } catch (error) {
        throw new Error(`Failed to sign transactions with TurnkeySigner: ${error.message}`);
      }
    },
  } as Adapter);
}

/**
 * Convert a Noble bech32 address (20-byte payload) to a Solana base58 pubkey.
 * Rule: prepend 12 zero bytes to the 20-byte payload → 32 bytes → base58.
 */
export function nobleToSolana(nobleAddress: string): string {
  // Decode bech32 (Cosmos-style uses "bech32", not bech32m, for account addrs)
  const dec = decode(nobleAddress.toLowerCase());

  // Optional safety check: ensure HRP is 'noble'
  if (dec.prefix !== 'noble') {
    throw new Error(`Unexpected HRP "${dec.prefix}". Expected "noble".`);
  }

  // Convert 5-bit words back to raw bytes
  const payload = Buffer.from(fromWords(dec.words));

  // Must be exactly 20 bytes (typical Cosmos-style account payload length)
  if (payload.length !== 20) {
    throw new Error(
      `Invalid payload length ${payload.length}, expected 20 bytes.`,
    );
  }

  // Build 32-byte buffer: 12 zero bytes + 20-byte payload
  const solanaBytes = Buffer.concat([Buffer.alloc(12, 0x00), payload]);

  // Base58-encode for Solana pubkey string
  return bs58.encode(solanaBytes);
}
