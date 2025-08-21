import { logger } from '@dydxprotocol-indexer/base';
import { route, messages } from '@skip-go/client/cjs';
import type { Adapter } from '@solana/wallet-adapter-base';
import type { Transaction, VersionedTransaction } from '@solana/web3.js';
import { Turnkey } from '@turnkey/sdk-server';
import { TurnkeySigner } from '@turnkey/solana';
import { encodeFunctionData, type Hex } from 'viem';
import type { SmartAccountImplementation } from 'viem/account-abstraction';
import {
  mainnet, arbitrum, avalanche, base, optimism,
} from 'viem/chains';

import config from '../config';
import { getAddress } from './alchemy-helpers';

const controllerName: string = 'skip-helper';

export const usdcAddressByChainId: Record<string, string> = {
  [mainnet.id.toString()]: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48', // usdc on ethereum mainnet.
  [arbitrum.id.toString()]: '0xaf88d065e77c8cC2239327C5EDb3A432268e5831', // usdc on arbitrum.
  [avalanche.id.toString()]: '0xb97ef9ef8734c71904d8002f8b6bc66dd9c48a6e', // usdc on avalanche.
  [base.id.toString()]: '0x833589fcd6edb6e08f4c7c32d4f71b54bda02913', // usdc on base.
  [optimism.id.toString()]: '0x0b2c639c533813f4aa9d7837caf62653d097ff85', // usdc on optimism.
  solana: 'EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v', // usdc on solana.
  'dydx-mainnet-1': 'ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5', // usdc on dydx.
};

const ethDenomByChainId: Record<string, string> = {
  [mainnet.id.toString()]: 'ethereum-native', // eth on ethereum mainnet.
  [arbitrum.id.toString()]: 'arbitrum-native', // eth on arbitrum.
  [base.id.toString()]: 'base-native', // eth on base.
  [optimism.id.toString()]: 'optimism-native', // eth on optimism.
};

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
    goFast: true,
  });
  logger.info({
    at: `${controllerName}#getSkipCallData`,
    message: 'Route result obtained',
    routeResult,
  });
  if (!routeResult) {
    throw new Error('Failed to find a route');
  }

  logger.info({
    at: `${controllerName}#getSkipCallData`,
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
        at: `${controllerName}#getSkipCallData`,
        message: `Message ${index + 1} EVM Transaction`,
        evmTx: msg.evmTx,
      });
      data = msg.evmTx.data || '';
      toAddress = msg.evmTx.to || '';
    } else if ('svmTx' in msg) {
      logger.info({
        at: `${controllerName}#getSkipCallData`,
        message: `Message ${index + 1} SVM Transaction`,
        svmTx: msg.svmTx,
      });
    } else if ('multiChainMsg' in msg) {
      logger.info({
        at: `${controllerName}#getSkipCallData`,
        message: `Message ${index + 1} Multi-Chain Message`,
        multiChainMsg: msg.multiChainMsg,
      });
    }
  });

  // need value to be the amount if native asset.
  let value = BigInt(0);
  if (Object.values(ethDenomByChainId).map(
    (x) => x.toLowerCase(),
  ).includes(sourceAssetDenom.toLowerCase())) {
    value = BigInt(amount);
  }

  const callData = [
    {
      to: (toAddress.startsWith('0x') ? toAddress : (`0x${toAddress}`)) as Hex,
      value,
      data: data.startsWith('0x') ? data as Hex : (`0x${data}`) as Hex, // "0x",
    },
  ];
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
        // @ts-ignore
        return await turnkeySigner.signTransaction(tx, signWith);
      } catch (error) {
        throw new Error(`Failed to sign transaction with TurnkeySigner: ${error.message}`);
      }
    },
    signAllTransactions: async (txs: (Transaction | VersionedTransaction)[]) => {
      try {
        // @ts-ignore
        return await turnkeySigner.signAllTransactions(txs, signWith);
      } catch (error) {
        throw new Error(`Failed to sign transactions with TurnkeySigner: ${error.message}`);
      }
    },
  } as Adapter);
}
