import { createKafkaMessage } from '@klyraprotocol-indexer/kafka';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  Timestamp,
  MarketEventV1,
  StatefulOrderEventV1,
} from '@klyraprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';

import { defaultHeight, defaultTime, defaultTxHash } from './constants';
import { createIndexerTendermintBlock, createIndexerTendermintEvent } from './indexer-proto-helpers';
import { KlyraIndexerSubtypes } from '../../src/lib/types';

export function createKafkaMessageFromMarketEvent({
  marketEvents,
  transactionIndex,
  height,
  time,
  txHash,
}: {
  marketEvents: MarketEventV1[],
  transactionIndex: number,
  height: number,
  time: Timestamp,
  txHash: string,
}): KafkaMessage {
  const events: IndexerTendermintEvent[] = [];
  for (let eventIndex: number = 0; eventIndex < marketEvents.length; eventIndex++) {
    events.push(
      createIndexerTendermintEvent(
        KlyraIndexerSubtypes.MARKET,
        MarketEventV1.encode(marketEvents[eventIndex]).finish(),
        transactionIndex,
        eventIndex,
      ),
    );
  }

  const block: IndexerTendermintBlock = createIndexerTendermintBlock(
    height,
    time,
    events,
    [txHash],
  );

  const binaryBlock: Uint8Array = IndexerTendermintBlock.encode(block).finish();
  return createKafkaMessage(Buffer.from(binaryBlock));
}

export function createKafkaMessageFromStatefulOrderEvent(
  event: StatefulOrderEventV1,
  transactionIndex: number = 0,
  height: number = defaultHeight,
  time: Timestamp = defaultTime,
  txHash: string = defaultTxHash,
): KafkaMessage {
  const events: IndexerTendermintEvent[] = [];
  events.push(
    createIndexerTendermintEvent(
      KlyraIndexerSubtypes.STATEFUL_ORDER,
      StatefulOrderEventV1.encode(event).finish(),
      transactionIndex,
      0,
    ),
  );

  const block: IndexerTendermintBlock = createIndexerTendermintBlock(
    height,
    time,
    events,
    [txHash],
  );

  const binaryBlock: Uint8Array = IndexerTendermintBlock.encode(block).finish();
  return createKafkaMessage(Buffer.from(binaryBlock));
}
