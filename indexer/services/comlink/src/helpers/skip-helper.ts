import { logger } from '@dydxprotocol-indexer/base';
import { PermissionApprovalTable } from '@dydxprotocol-indexer/postgres';
import { route, messages } from '@skip-go/client/cjs';
import type { Adapter } from '@solana/wallet-adapter-base';
import type { Transaction, VersionedTransaction } from '@solana/web3.js';
import { Turnkey } from '@turnkey/sdk-server';
import { TurnkeySigner } from '@turnkey/solana';
import { createAccount } from '@turnkey/viem';
import { deserializePermissionAccount } from '@zerodev/permissions';
import { toECDSASigner } from '@zerodev/permissions/signers';
import { CreateKernelAccountReturnType } from '@zerodev/sdk';
import { KERNEL_V3_1, KERNEL_V3_3 } from '@zerodev/sdk/constants';
import { decode, fromWords } from 'bech32';
import bs58 from 'bs58';
import { min } from 'lodash';
import { encodeFunctionData, type Hex } from 'viem';
import type { EntryPointVersion, SmartAccountImplementation } from 'viem/account-abstraction';
import { avalanche } from 'viem/chains';

import config from '../config';
import {
  dydxChainId,
  entryPoint,
  ETH_USDC_QUANTUM,
  ETH_WEI_QUANTUM,
  ethDenomByChainId, usdcAddressByChainId,
} from '../lib/smart-contract-constants';
import { getAddress, getETHPrice, publicClients } from './alchemy-helpers';

const turnkeyClient = new Turnkey({
  apiBaseUrl: config.TURNKEY_API_BASE_URL as string,
  apiPublicKey: config.TURNKEY_API_PUBLIC_KEY as string,
  apiPrivateKey: config.TURNKEY_API_PRIVATE_KEY as string,
  defaultOrganizationId: config.TURNKEY_ORGANIZATION_ID,
});

export async function buildUserAddresses(
  requiredChainAddresses: string[],
  sourceAddress: string,
  dydxAddress: string,
) {
  return Promise.all(
    requiredChainAddresses.map(async (cid: string) => ({
      chainId: cid,
      address: await getAddress(cid, sourceAddress, dydxAddress),
    })),
  );
}
const nobleForwardingModule = 'https://api.noble.xyz/noble/forwarding/v1/address/channel';
const skipMessagesTimeoutSeconds = '60';
const dydxNobleChannel = 33;
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
    destAssetDenom: usdcAddressByChainId[dydxChainId],
    destAssetChainId: dydxChainId,
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

  const userAddresses = await buildUserAddresses(
    routeResult.requiredChainAddresses,
    sourceAddress,
    dydxAddress,
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

  // acceptable slippage is smallest of SKIP_SLIPPAGE_TOLERANCE_USDC (Default $100) divided
  // by the estimatedAmountOut or the SKIP_SLIPPAGE_TOLERANCE_PERCENTAGE.
  const slippageTolerancePercent = getSlippageTolerancePercent(routeResult.estimatedAmountOut);
  const response = await messages({
    timeoutSeconds: skipMessagesTimeoutSeconds,
    amountIn: routeResult?.amountIn,
    amountOut: routeResult.estimatedAmountOut,
    sourceAssetChainId: routeResult?.sourceAssetChainId,
    sourceAssetDenom: routeResult?.sourceAssetDenom,
    destAssetChainId: routeResult?.destAssetChainId,
    destAssetDenom: routeResult?.destAssetDenom,
    operations: routeResult?.operations,
    addressList,
    slippageTolerancePercent,
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
    value = BigInt(amountToUse);
  }

  // this is the actual call data that is responsible for the bridge.
  const callData = [
    {
      to: (toAddress.startsWith('0x') ? toAddress : (`0x${toAddress}`)) as Hex,
      value,
      data: data.startsWith('0x') ? data as Hex : (`0x${data}`) as Hex, // "0x",
    },
  ];

  // swap eth to usdc.
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
          BigInt(amountToUse),
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
 * Convert a Noble bech32 address to a hex address.
 * Rule: prepend with 0s so that the end result is 32 bytes.
 */
export function nobleToHex(addr: string): string {
  // decode bech32
  const decoded = decode(addr);

  if (decoded.prefix !== 'noble' && decoded.prefix !== 'dydx') {
    throw new Error(`Invalid HRP: expected "noble", got "${decoded.prefix}"`);
  }

  // convert back from words to bytes
  const addrBytes = Buffer.from(fromWords(decoded.words));

  if (addrBytes.length !== 20) {
    throw new Error(`Invalid address length: expected 20, got ${addrBytes.length}`);
  }

  // left-pad to 32 bytes (as in your original hex form)
  const padded = Buffer.concat([Buffer.alloc(12, 0), addrBytes]); // 12 zeros + 20 bytes

  return `0x${padded.toString('hex')}`;
}

/**
 * Convert a string to its hex representation then left pads 277 bytes of 0s.
 * @param s - The string to convert.
 * @returns The hex representation of the string.
 */
export function encodeToHexAndPad(s: string): string {
  // 277 is the offset for destination call data in the skip go fast smart contract call.
  const offset = 277;
  const hex = Buffer.from(s).toString('hex');
  const padded = Buffer.concat([Buffer.alloc(offset, 0), Buffer.from(hex, 'hex')]);
  return `0x${padded.toString('hex')}`;
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

export async function getKernelAccount(
  chainId: string,
  fromAddress: string,
  suborgId: string,
): Promise<CreateKernelAccountReturnType<EntryPointVersion>> {
  // Initialize the turnkey delegated signer account.
  const turnkeyAccount = await createAccount({
    // @ts-ignore
    client: turnkeyClient.apiClient(),
    signWith: config.APPROVAL_SIGNER_PUBLIC_ADDRESS,
  });
  // if smart account approval is enabled, use the session key + approval to sign for txs.
  // use the permissioned master key as  a signer.
  const sessionKeySigner = await toECDSASigner({
    signer: turnkeyAccount,
  });
  let kernelVersion = KERNEL_V3_3;
  if (chainId === avalanche.id.toString()) {
    kernelVersion = KERNEL_V3_1;
  }

  const row = await PermissionApprovalTable.findBySuborgIdAndChainId(suborgId, chainId);
  if (!row) {
    throw new Error(`No approval found for suborg ${suborgId} and chain ${chainId}`);
  }
  const sessionKeyAccount = await deserializePermissionAccount(
    publicClients[chainId],
    entryPoint,
    kernelVersion,
    row.approval,
    sessionKeySigner,
  );
  return sessionKeyAccount;
}

// for a dydx address, this returns the noble forwarding address of the dydx address.
export async function getNobleForwardingAddress(dydxAddress: string): Promise<string> {
  const endpoint = `${nobleForwardingModule}-${dydxNobleChannel}/${dydxAddress}/`;

  const ac = new AbortController();
  const timeout = setTimeout(() => ac.abort(), 10_000);
  try {
    const response = await fetch(endpoint, {
      signal: ac.signal,
    });
    if (!response.ok) {
      throw new Error(`failed to get a forwarding address: ${response.statusText}`);
    }
    const data = await response.json();
    if (!data || (data && !data.address)) {
      throw new Error('failed to get a forwarding address');
    }
    return data.address;
  } catch (e) {
    throw new Error(`failed to get a forwarding address: ${e}`);
  } finally {
    clearTimeout(timeout);
  }
}

export async function limitAmount(
  chainId: string,
  amount: string,
  sourceAssetDenom: string,
): Promise<string> {
  let amountToUse = BigInt(amount);
  // calculates the most eth we can bridge in one go and pins it to that.
  if (sourceAssetDenom === ethDenomByChainId[chainId]) {
    try {
      const ethPrice = await getETHPrice();
      const maxDepositInWei = Math.floor(
        (config.MAXIMUM_BRIDGE_AMOUNT_USDC / ethPrice) * ETH_WEI_QUANTUM,
      );
      amountToUse = min([
        amountToUse,
        BigInt(maxDepositInWei),
      ])!;
    } catch (error) {
      logger.error({
        at: 'skip-helper#limitAmount',
        message: 'Failed to get ETH price',
        error,
      });
      throw error;
    }
    return amountToUse.toString();
  }

  // calculates the most usdc we can bridge in one go and pins it to that.
  const maxDepositInUsdc = config.MAXIMUM_BRIDGE_AMOUNT_USDC;
  return min([amountToUse, BigInt(maxDepositInUsdc * ETH_USDC_QUANTUM)])!.toString();
}

// getSlippageTolerancePercent returns the acceptable slippage is smallest of
// SKIP_SLIPPAGE_TOLERANCE_USDC (Default $100) divided by the estimatedAmountOut
// or the SKIP_SLIPPAGE_TOLERANCE_PERCENTAGE.
export function getSlippageTolerancePercent(estAmountOut: string): string {
  return min([
    (100 * (config.SKIP_SLIPPAGE_TOLERANCE_USDC * ETH_USDC_QUANTUM)) / parseInt(estAmountOut, 10),
    parseFloat(config.SKIP_SLIPPAGE_TOLERANCE_PERCENTAGE),
  ])!.toString();
}
