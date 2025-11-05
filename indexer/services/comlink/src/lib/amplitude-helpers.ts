import { init, track } from '@amplitude/analytics-node';
import { logger } from '@dydxprotocol-indexer/base';

import config from '../config';

// Initialize Amplitude client
let amplitudeInitialized = false;

export function initializeAmplitude(): void {
  if (amplitudeInitialized) {
    return;
  }

  try {
    // Only initialize if API key is provided
    if (config.AMPLITUDE_API_KEY) {
      init(config.AMPLITUDE_API_KEY, {
        serverUrl: config.AMPLITUDE_SERVER_URL || 'https://api2.amplitude.com/2/httpapi',
        flushQueueSize: 1, // Send events immediately
        flushIntervalMillis: 1000, // Flush every second
        serverZone: 'EU',
      });
      amplitudeInitialized = true;
      logger.info({
        at: 'amplitude-helpers#initializeAmplitude',
        message: 'Amplitude client initialized successfully',
      });
    } else {
      logger.info({
        at: 'amplitude-helpers#initializeAmplitude',
        message: 'Amplitude API key not provided, events will not be tracked',
      });
    }
  } catch (error) {
    logger.error({
      at: 'amplitude-helpers#initializeAmplitude',
      message: 'Failed to initialize Amplitude client',
      error: error as Error,
    });
  }
}

export async function trackAmplitudeEvent(
  eventType: string,
  userId?: string,
  eventProperties?: Record<string, unknown>,
): Promise<void> {
  if (!amplitudeInitialized || !config.AMPLITUDE_API_KEY) {
    logger.debug({
      at: 'amplitude-helpers#trackAmplitudeEvent',
      message: 'Amplitude not initialized or API key not provided, skipping event tracking',
      eventType,
    });
    return;
  }

  try {
    await track(eventType, eventProperties, {
      user_id: userId,
    });

    logger.debug({
      at: 'amplitude-helpers#trackAmplitudeEvent',
      message: 'Amplitude event tracked successfully',
      eventType,
      userId,
      eventProperties,
    });
  } catch (error) {
    logger.error({
      at: 'amplitude-helpers#trackAmplitudeEvent',
      message: 'Failed to track Amplitude event',
      error: error as Error,
      eventType,
      userId,
      eventProperties,
    });
  }
}

// Specific helper for TurnKey deposit events
export async function trackTurnkeyDepositSubmitted(
  userId: string,
  chainId: string,
  amount: string,
  transactionHash: string,
  sourceAssetDenom?: string,
): Promise<void> {
  await trackAmplitudeEvent(
    'TurnkeyDepositSubmitted',
    userId, // Using fromAddress as userId
    {
      chain_id: chainId,
      amount,
      transaction_hash: transactionHash,
      source_asset_denom: sourceAssetDenom,
      timestamp: new Date().toISOString(),
    },
  );
}
