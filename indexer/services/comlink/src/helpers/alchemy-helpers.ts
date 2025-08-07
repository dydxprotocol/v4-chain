import { logger } from '@dydxprotocol-indexer/base';
import { findByEvmAddress, findBySvmAddress } from '@dydxprotocol-indexer/postgres/build/src/stores/turnkey-users-table';
import { TurnkeyUserFromDatabase } from '@dydxprotocol-indexer/postgres/build/src/types';
import { getKernelAddressFromECDSA } from '@zerodev/ecdsa-validator';
import { getEntryPoint, KERNEL_V3_1 } from '@zerodev/sdk/constants';
import {
  Address, Chain, createPublicClient, http,
} from 'viem';
import {
  arbitrum, avalanche, base, mainnet, optimism,
} from 'viem/chains';

import config from '../config';

const evmChainIdToAlchemyWebhookId: Record<string, string> = {
  [mainnet.id.toString()]: 'wh_ys5e0lhw2iaq0wge',
  [arbitrum.id.toString()]: 'wh_fvxtvyg2uxh0ylba',
  [avalanche.id.toString()]: 'wh_ycy4khfozgyuir3u',
  [base.id.toString()]: 'wh_8pntnwk3jltyduwe',
  [optimism.id.toString()]: 'wh_99yjvuacl28obf0i',
};

const solanaAlchemyWebhookId = 'wh_vv1go1c7wy53q6zy';

function getRPCEndpoint(chainId: string): string {
  if (!Object.keys(chains).includes(chainId)) {
    throw new Error(`Unsupported chainId: ${chainId}`);
  }
  return `${config.ZERODEV_API_BASE_URL}/${config.ZERODEV_API_KEY}/chain/${chainId}`;
}

const chains: Record<string, Chain> = {
  [mainnet.id.toString()]: mainnet,
  [arbitrum.id.toString()]: arbitrum,
  [avalanche.id.toString()]: avalanche,
  [base.id.toString()]: base,
  [optimism.id.toString()]: optimism,
};

export async function addAddressesToAlchemyWebhook(evm?: string, svm?: string): Promise<void> {
  try {
    // Add EVM address to webhook for monitoring
    if (evm) {
      const record: TurnkeyUserFromDatabase | undefined = await findByEvmAddress(evm);
      if (!record) {
        throw new Error(`EVM address does not exist in the database: ${evm}`);
      }
      // Iterate over all EVM networks and register the address with each webhook
      for (const [chainId, webhookId] of Object.entries(evmChainIdToAlchemyWebhookId)) {
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
        }
      }
    }

    // Add SVM address to webhook for monitoring
    if (svm) {
      const record: TurnkeyUserFromDatabase | undefined = await findBySvmAddress(svm);
      if (!record) {
        throw new Error(`SVM address does not exist in the database: ${svm}`);
      }
      await registerAddressWithAlchemyWebhookWithRetry(svm, solanaAlchemyWebhookId);
    }

  } catch (error) {
    logger.error({
      at: 'TurnkeyController#addAddressesToAlchemyWebhook',
      message: 'Failed to add addresses to Alchemy webhook',
      error,
      evmAddress: evm,
      svmAddress: svm,
    });
  }
}

// Register address with Alchemy webhook using REST API
export async function registerAddressWithAlchemyWebhook(
  address: string,
  webhookId: string,
): Promise<void> {
  const webhookUrl = 'https://dashboard.alchemy.com/api/update-webhook-addresses';
  const addressesToAdd: string[] = [address];
  if (webhookId === evmChainIdToAlchemyWebhookId[avalanche.id.toString()]) {
    // for avalanche, we also should add the smart account address to the webhook.
    const smartAccountAddress = await getSmartAccountAddress(address);
    addressesToAdd.push(smartAccountAddress);
  }
  const response = await fetch(webhookUrl, {
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
      await new Promise((resolve) => setTimeout(resolve, delay * (i + 1))); // Exponential backoff
    }
  }
}

/*
 * Returns the smart account address indexed at 0 with entry point 0.7.
 * Also assumes that the address provided here is a valid address that
 * already exists in our database.
 */
async function getSmartAccountAddress(address: string): Promise<string> {
  const publicAvalancheClient = createPublicClient({
    transport: http(getRPCEndpoint(avalanche.id.toString())),
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
