import { stats } from '@dydxprotocol-indexer/base';
import {
  KafkaTopics,
  MARKETS_WEBSOCKET_MESSAGE_VERSION,
  SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
  TRADES_WEBSOCKET_MESSAGE_VERSION,
} from '@dydxprotocol-indexer/kafka';
import {
  IndexerTendermintEvent,
  MarketMessage,
  OffChainUpdateV1,
  SubaccountId,
  SubaccountMessage,
} from '@dydxprotocol-indexer/v4-protos';

import config from '../config';
import { indexerTendermintEventToTransactionIndex } from '../lib/helper';
import { ConsolidatedKafkaEvent, SingleTradeMessage } from '../lib/types';

export function generateConsolidatedSubaccountKafkaEvent(
  contents: string,
  subaccountId: SubaccountId,
  blockHeight: string,
  indexerTendermintEvent: IndexerTendermintEvent,
): ConsolidatedKafkaEvent {
  stats.increment(`${config.SERVICE_NAME}.create_subaccount_kafka_event`, 1);
  const subaccountMessage: SubaccountMessage = {
    blockHeight,
    transactionIndex: indexerTendermintEventToTransactionIndex(indexerTendermintEvent),
    eventIndex: indexerTendermintEvent.eventIndex,
    contents,
    subaccountId,
    version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
  };

  return {
    topic: KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS,
    message: subaccountMessage,
  };
}

export function generateConsolidatedMarketKafkaEvent(
  contents: string,
): ConsolidatedKafkaEvent {
  stats.increment(`${config.SERVICE_NAME}.create_market_kafka_event`, 1);
  const marketMessage: MarketMessage = {
    contents,
    version: MARKETS_WEBSOCKET_MESSAGE_VERSION,
  };

  return {
    topic: KafkaTopics.TO_WEBSOCKETS_MARKETS,
    message: marketMessage,
  };
}

export function generateConsolidatedTradeKafkaEvent(
  contents: string,
  clobPairId: string,
  blockHeight: string,
  indexerTendermintEvent: IndexerTendermintEvent,
): ConsolidatedKafkaEvent {
  stats.increment(`${config.SERVICE_NAME}.create_trade_kafka_event`, 1);
  const tradeMessage: SingleTradeMessage = {
    blockHeight,
    transactionIndex: indexerTendermintEventToTransactionIndex(indexerTendermintEvent),
    eventIndex: indexerTendermintEvent.eventIndex,
    contents,
    clobPairId,
    version: TRADES_WEBSOCKET_MESSAGE_VERSION,
  };

  return {
    topic: KafkaTopics.TO_WEBSOCKETS_TRADES,
    message: tradeMessage,
  };
}

export function generateConsolidatedVulcanKafkaEvent(
  key: Buffer,
  offChainUpdate: OffChainUpdateV1,
): ConsolidatedKafkaEvent {
  stats.increment(`${config.SERVICE_NAME}.create_vulcan_kafka_event`, 1);

  return {
    topic: KafkaTopics.TO_VULCAN,
    message: {
      key,
      value: offChainUpdate,
    },
  };
}
