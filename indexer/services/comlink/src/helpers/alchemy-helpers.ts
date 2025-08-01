import { logger } from '@dydxprotocol-indexer/base';
import config from '../config';
import { arbitrum, avalanche, base, mainnet, optimism } from 'viem/chains';

const evmChainIdToAlchemyWebhookId: Record<string, string> = {
  [mainnet.id.toString()]: 'wh_ys5e0lhw2iaq0wge',
  [arbitrum.id.toString()]: 'wh_fvxtvyg2uxh0ylba',
  [avalanche.id.toString()]: 'wh_ycy4khfozgyuir3u',
  [base.id.toString()]: 'wh_8pntnwk3jltyduwe',
  [optimism.id.toString()]: 'wh_99yjvuacl28obf0i',
};

const solanaAlchemyWebhookId = 'wh_vv1go1c7wy53q6zy';

export async function addAddressesToAlchemyWebhook(evmAddress: string, svmAddress?: string): Promise<void> {
  try {
    // Add EVM address to webhook for monitoring
    if (evmAddress) {
      // Iterate over all EVM networks and register the address with each webhook
      for (const [chainId, webhookId] of Object.entries(evmChainIdToAlchemyWebhookId)) {
        try {
          await registerAddressWithAlchemyWebhookWithRetry(evmAddress, webhookId);
          logger.info({
            at: 'TurnkeyController#addAddressesToAlchemyWebhook',
            message: `Successfully registered EVM address with webhook for chain ${chainId}`,
            address: evmAddress,
            chainId,
            webhookId,
          });
        } catch (error) {
          logger.error({
            at: 'TurnkeyController#addAddressesToAlchemyWebhook',
            message: `Failed to register EVM address with webhook for chain ${chainId} after retries`,
            error,
            address: evmAddress,
            chainId,
            webhookId,
          });
        }
      }
    }

    // Add SVM address to webhook for monitoring
    if (svmAddress) {
      await registerAddressWithAlchemyWebhookWithRetry(svmAddress, solanaAlchemyWebhookId);
      logger.info({
        at: 'TurnkeyController#addAddressesToAlchemyWebhook',
        message: 'Successfully added svm address to Alchemy webhook',
        evmAddress,
        svmAddress,
      });
    }

  } catch (error) {
    logger.error({
      at: 'TurnkeyController#addAddressesToAlchemyWebhook',
      message: 'Failed to add addresses to Alchemy webhook',
      error,
      evmAddress,
      svmAddress,
    });
    // Don't throw error to avoid breaking the main flow
  }
}

// Register address with Alchemy webhook using REST API
export async function registerAddressWithAlchemyWebhook(address: string, webhookId: string): Promise<void> {
  const webhookUrl = 'https://dashboard.alchemy.com/api/update-webhook-addresses';

  const response = await fetch(webhookUrl, {
    method: 'PATCH',
    headers: {
      'X-Alchemy-Token': config.ALCHEMY_AUTH_TOKEN,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      webhook_id: webhookId,
      addresses_to_add: [address],
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
async function registerAddressWithAlchemyWebhookWithRetry(address: string, webhookId: string): Promise<void> {
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
      await new Promise(resolve => setTimeout(resolve, delay * (i + 1))); // Exponential backoff
    }
  }
}
