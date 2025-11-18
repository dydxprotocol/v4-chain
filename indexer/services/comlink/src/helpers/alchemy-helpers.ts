import { logger } from '@dydxprotocol-indexer/base';
import { findByEvmAddress, findBySmartAccountAddress, findBySvmAddress } from '@dydxprotocol-indexer/postgres/build/src/stores/turnkey-users-table';
import { TurnkeyUserFromDatabase } from '@dydxprotocol-indexer/postgres/build/src/types';
import { getKernelAddressFromECDSA } from '@zerodev/ecdsa-validator';
import { getEntryPoint, KERNEL_V3_1 } from '@zerodev/sdk/constants';
import { Alchemy } from 'alchemy-sdk';
import { decode, encode } from 'bech32';
import express from 'express';
import {
  Address, Chain, createPublicClient, http, checksumAddress,
  PublicClient,
} from 'viem';
import {
  arbitrum, avalanche, base, mainnet, optimism,
} from 'viem/chains';

import config from '../config';
import { create4xxResponse } from '../lib/helpers';
import {
  nobleChainId, osmosisChainId, neutronChainId, dydxChainId,
} from '../lib/smart-contract-constants';

const evmChainIdToAlchemyWebhookId: Record<string, string> = {
  [mainnet.id.toString()]: config.ETHEREUM_WEBHOOK_ID,
  [arbitrum.id.toString()]: config.ARBITRUM_WEBHOOK_ID,
  [avalanche.id.toString()]: config.AVALANCHE_WEBHOOK_ID,
  [base.id.toString()]: config.BASE_WEBHOOK_ID,
  [optimism.id.toString()]: config.OPTIMISM_WEBHOOK_ID,
};

const solanaAlchemyWebhookId = config.SOLANA_WEBHOOK_ID;

export const alchemyNetworkToChainIdMap: Record<string, string> = {
  ARB_MAINNET: arbitrum.id.toString(),
  AVAX_MAINNET: avalanche.id.toString(),
  BASE_MAINNET: base.id.toString(),
  OPT_MAINNET: optimism.id.toString(),
  ETH_MAINNET: mainnet.id.toString(),
  SOLANA_MAINNET: 'solana',
};

export const chains: Record<string, Chain> = {
  [mainnet.id.toString()]: mainnet,
  [arbitrum.id.toString()]: arbitrum,
  [avalanche.id.toString()]: avalanche,
  [base.id.toString()]: base,
  [optimism.id.toString()]: optimism,
};

export const chainInAlchemy: Record<string, string> = {
  [mainnet.id.toString()]: 'eth-mainnet',
  [arbitrum.id.toString()]: 'arb-mainnet',
  [avalanche.id.toString()]: 'avax-mainnet',
  [base.id.toString()]: 'base-mainnet',
  [optimism.id.toString()]: 'opt-mainnet',
};

export const publicClients = Object.keys(chains).reduce((acc, chainId) => {
  acc[chainId] = createPublicClient({
    transport: http(getAlchemyRPCEndpoint(chainId)),
    chain: chains[chainId],
  });
  return acc;
}, {} as Record<string, PublicClient>);

const alchemy = new Alchemy({
  apiKey: config.ALCHEMY_API_KEY,
});

export async function addAddressesToAlchemyWebhook(evm?: string, svm?: string): Promise<void> {
  const errors: string[] = [];

  // Add EVM address to webhook for monitoring
  if (evm) {
    const record: TurnkeyUserFromDatabase | undefined = await findByEvmAddress(evm);
    if (!record) {
      throw new Error(`EVM address does not exist in the database: ${evm}`);
    }
    // Register the address with all EVM webhooks in parallel
    await Promise.allSettled(
      Object.entries(evmChainIdToAlchemyWebhookId).map(async ([chainId, webhookId]) => {
        try {
          await registerAddressWithAlchemyWebhookWithRetry(evm, webhookId);
          logger.info({
            at: 'TurnkeyController#addAddressesToAlchemyWebhook',
            message: `Successfully registered EVM address with webhook for chain ${chainId}`,
            address: evm,
            chainId,
            webhookId,
          });
        } catch (error) {
          logger.error({
            at: 'TurnkeyController#addAddressesToAlchemyWebhook',
            message: `Failed to register EVM address with webhook for chain ${chainId} after retries`,
            error,
            address: evm,
            chainId,
            webhookId,
          });
          errors.push(`Failed to register EVM address with webhook for chain ${chainId} after retries`);
        }
      }),
    );
  }

  // Add SVM address to webhook for monitoring
  if (svm) {
    const record: TurnkeyUserFromDatabase | undefined = await findBySvmAddress(svm);
    if (!record) {
      throw new Error(`SVM address does not exist in the database: ${svm}`);
    }
    try {
      await registerAddressWithAlchemyWebhookWithRetry(svm, solanaAlchemyWebhookId);
    } catch (error) {
      logger.error({
        at: 'TurnkeyController#addAddressesToAlchemyWebhook',
        message: 'Failed to add addresses to Alchemy webhook',
        error,
        evmAddress: evm,
        svmAddress: svm,
      });
      errors.push(`Failed to add addresses to Solana Alchemy webhook: ${error}`);
    }
  }

  // If there were any errors, log them but don't throw - allow partial success
  if (errors.length > 0) {
    logger.warning({
      at: 'TurnkeyController#addAddressesToAlchemyWebhook',
      message: `Some webhook registrations failed: ${errors.join('; ')}`,
      evmAddress: evm,
      svmAddress: svm,
      errors,
    });
  }
}

// Register address with Alchemy webhook using REST API
export async function registerAddressWithAlchemyWebhook(
  address: string,
  webhookId: string,
): Promise<void> {
  if (!config.ALCHEMY_AUTH_TOKEN) {
    throw new Error('ALCHEMY_AUTH_TOKEN is not set: cannot register address with Alchemy webhook');
  }
  const addressesToAdd: string[] = [];
  if (webhookId === evmChainIdToAlchemyWebhookId[avalanche.id.toString()]) {
    // for avalanche, we also should add the smart account address to the webhook.
    const smartAccountAddress = await getSmartAccountAddress(address);
    addressesToAdd.push(smartAccountAddress);
  } else {
    addressesToAdd.push(address);
  }
  const response = await fetch(config.ALCHEMY_WEBHOOK_UPDATE_URL, {
    method: 'PATCH',
    headers: {
      'X-Alchemy-Token': config.ALCHEMY_AUTH_TOKEN,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      webhook_id: webhookId,
      addresses_to_add: addressesToAdd,
      addresses_to_remove: [],
    }),
  });

  if (!response.ok) {
    const errorText = await response.text();
    logger.error({
      at: 'TurnkeyController#registerAddressWithAlchemyWebhook',
      message: `Failed to register address with Alchemy webhook: ${response.statusText} - ${errorText}`,
      address,
      webhookId,
    });
    throw new Error(`Failed to register address with Alchemy webhook: ${response.statusText} - ${errorText}`);
  }

  logger.info({
    at: 'TurnkeyController#registerAddressWithAlchemyWebhook',
    message: `Address ${address} successfully added to Alchemy webhook`,
    address,
  });
}

// Register address with Alchemy webhook using REST API with retry logic
async function registerAddressWithAlchemyWebhookWithRetry(
  address: string,
  webhookId: string,
): Promise<void> {
  const maxRetries = 3;
  const delay = 1000; // 1 second

  for (let i = 0; i < maxRetries; i++) {
    try {
      await registerAddressWithAlchemyWebhook(address, webhookId);
      return; // Success, exit retry loop
    } catch (error) {
      if (i === maxRetries - 1) {
        logger.error({
          at: 'TurnkeyController#registerAddressWithAlchemyWebhookWithRetry',
          message: `Failed to register address with Alchemy webhook after ${maxRetries} retries`,
          error,
          address,
          webhookId,
        });
        throw error; // Re-throw the error after all retries fail
      }
      logger.warning({
        at: 'TurnkeyController#registerAddressWithAlchemyWebhookWithRetry',
        message: `Retrying Alchemy webhook registration for address ${address} (attempt ${i + 1}/${maxRetries})`,
        error,
        address,
        webhookId,
      });
      await new Promise((resolve) => setTimeout(resolve, delay * (i + 1))); // linear backoff
    }
  }
}

/*
 * Returns the smart account address indexed at 0 with entry point 0.7.
 * Also assumes that the address provided here is a valid address that
 * already exists in our database.
 */
export async function getSmartAccountAddress(address: string): Promise<string> {
  const publicAvalancheClient = createPublicClient({
    transport: http(getAlchemyRPCEndpoint(avalanche.id.toString())),
    chain: avalanche,
  });

  const kernelAddress = await getKernelAddressFromECDSA({
    publicClient: publicAvalancheClient,
    entryPoint: getEntryPoint('0.7'),
    kernelVersion: KERNEL_V3_1,
    eoaAddress: address as Address,
    index: BigInt(0),
  });
  return kernelAddress;
}

export async function getEOAAddressFromSmartAccountAddress(
  smartAccountAddress: Address,
): Promise<Address> {
  const smartAccountAddressToUse = checksumAddress(smartAccountAddress);
  const record = await findBySmartAccountAddress(smartAccountAddressToUse);
  if (!record || !record.evm_address) {
    throw new Error('Failed to find a turnkey user for address');
  }
  return record.evm_address as Address;
}

export enum CosmosPrefix {
  OSMO = 'osmo',
  NEUTRON = 'neutron',
  NOBLE = 'noble',
}

// Prefix is one of osmosis, neutron, noble. This is how we convert dydx addresses
// to other chain addresses on cosmos. Address here is dydx address.
export function toClientAddressWithPrefix(prefix: CosmosPrefix, address: string): string {
  try {
    const decoded = decode(address);
    if (decoded.prefix !== 'dydx') {
      throw new Error('Incoming address is not a dydx address');
    }
    return encode(prefix, decoded.words);
  } catch (e) {
    throw new Error('Failed to convert dydx address to client address');
  }
}

export function isSupportedEVMChainId(chainId: string): boolean {
  return Object.keys(chains).includes(chainId);
}

export function getZeroDevRPCEndpoint(chainId: string): string {
  if (!isSupportedEVMChainId(chainId)) {
    throw new Error(`Unsupported chainId: ${chainId}`);
  }
  return `${config.ZERODEV_API_BASE_URL}/${config.ZERODEV_API_KEY}/chain/${chainId}`;
}

export function getAlchemyRPCEndpoint(chainId: string): string {
  if (!isSupportedEVMChainId(chainId)) {
    throw new Error(`Unsupported chainId: ${chainId}`);
  }
  return `https://${chainInAlchemy[chainId]}.g.alchemy.com/v2/${config.ALCHEMY_API_KEY}`;
}

// TODO: Verify that this function is 1000% correct. @RUI and @TYLER and @JARED
export function getAddress(
  chainId: string,
  sourceAddress: string,
  dydxAddress: string,
): string {
  if (isSupportedEVMChainId(chainId) || chainId === 'solana') {
    return sourceAddress;
  }
  switch (chainId) {
    case nobleChainId:
      return toClientAddressWithPrefix(CosmosPrefix.NOBLE, dydxAddress);
    case osmosisChainId:
      return toClientAddressWithPrefix(CosmosPrefix.OSMO, dydxAddress);
    case neutronChainId:
      return toClientAddressWithPrefix(CosmosPrefix.NEUTRON, dydxAddress);
    case dydxChainId:
      return dydxAddress;
    default:
      throw new Error(`Unsupported chain ID: ${chainId}`);
  }
}

// middleware to verify alchemy webhook signature
export function verifyAlchemyWebhook(
  req: express.Request,
  res: express.Response,
  next: express.NextFunction,
) {
  const token = req.header('x-alchemy-signature') || '';
  if (!token || token !== config.ALCHEMY_AUTH_TOKEN) {
    return create4xxResponse(res, 'unauthorized webhook', 401);
  }
  return next();
}

export async function getETHPrice(): Promise<number> {
  try {
    const price = await alchemy.prices.getTokenPriceBySymbol(['ETH']);

    // Check if we have valid price data
    if (!price.data || price.data.length === 0) {
      throw new Error('No price data returned from Alchemy API');
    }

    const priceData = price.data[0];
    if (!priceData.prices || priceData.prices.length === 0) {
      if (priceData.error) {
        throw new Error(`Alchemy API error: ${priceData.error.message}`);
      }
      throw new Error('No price data available for ETH');
    }

    return parseFloat(priceData.prices[0].value);
  } catch (error) {
    // Properly serialize error for logging
    const errorDetails = {
      message: error instanceof Error ? error.message : String(error),
      name: error instanceof Error ? error.name : 'Unknown',
      stack: error instanceof Error ? error.stack : undefined,
    };

    logger.error({
      at: 'alchemy-helpers#getETHPrice',
      message: 'Failed to get ETH price',
      error: errorDetails,
    });
    throw error;
  }
}
