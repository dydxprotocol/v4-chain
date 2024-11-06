import { logger, stats } from '@klyraprotocol-indexer/base';
import _ from 'lodash';
import * as pg from 'pg';

import config from '../config';
import { KafkaPublisher } from './kafka-publisher';
import { ConsolidatedKafkaEvent, KlyraIndexerSubtypes, EventMessage } from './types';
import { Handler } from '../handlers/handler';

// type alias for an array of handlers.
type HandlerBatch = Handler<EventMessage>[];
export const SYNCHRONOUS_SUBTYPES: KlyraIndexerSubtypes[] = [
  KlyraIndexerSubtypes.MARKET,
  KlyraIndexerSubtypes.ASSET,
  KlyraIndexerSubtypes.LIQUIDITY_TIER,
  KlyraIndexerSubtypes.PERPETUAL_MARKET,
  KlyraIndexerSubtypes.UPDATE_PERPETUAL,
  KlyraIndexerSubtypes.UPDATE_CLOB_PAIR,
  KlyraIndexerSubtypes.FUNDING,
];

/**
 * A class that processes handlers sequentially.
 *
 * During genesis, these events should be handled prior to any events in BatchedHandlers.
 * After genesis block, these events should be handled after events in BatchedHandlers.
 * It is used for processing asset and market events.
 */
// TODO(IND-514): Remove the batch and sync handlers completely by moving all redis updates into
// a pipeline similar to how we return kafka events and then batch and emit them.
export class SyncHandlers {
  handlerBatch: HandlerBatch;
  initializationTime: number;

  constructor() {
    this.handlerBatch = [];
    this.initializationTime = Date.now();
  }

  /**
   * Adds a handler which contains an event to be processed. This function should be called in the
   * order in which the events should be processed. The handlers will be processed sequentially.
   *
   * @param indexerSubtype The event subtype
   * @param handler The handler to add to the batched handlers
   */
  public addHandler(
    indexerSubtype: KlyraIndexerSubtypes,
    handler: Handler<EventMessage>,
  ): void {
    if (!SYNCHRONOUS_SUBTYPES.includes(indexerSubtype)) {
      logger.error({
        at: 'SyncHandlers#addHandler',
        message: `Invalid indexerSubtype: ${indexerSubtype}`,
      });
      return;
    }
    // @ts-ignore
    this.handlerBatch.push(handler);
  }

  /**
   * Processes all handlers that were passed in through `addHandler` sequentially.
   * Adds events to the kafkaPublisher.
   */
  public async process(
    kafkaPublisher: KafkaPublisher, resultRow: pg.QueryResultRow,
  ): Promise<void> {
    const start: number = Date.now();
    const handlerCountMapping: { [key: string]: number } = {};
    const consolidatedKafkaEventGroup: ConsolidatedKafkaEvent[][] = [];
    for (const handler of this.handlerBatch) {
      const handlerName: string = handler.constructor.name;
      if (!(handlerName in handlerCountMapping)) {
        handlerCountMapping[handlerName] = 0;
      }
      handlerCountMapping[handlerName] += 1;
      const events: ConsolidatedKafkaEvent[] = await handler.handle(
        resultRow[handler.blockEventIndex]);
      consolidatedKafkaEventGroup.push(events);
    }

    _.forEach(consolidatedKafkaEventGroup, (events: ConsolidatedKafkaEvent[]) => {
      kafkaPublisher.addEvents(events);
    });
    logger.info({
      at: 'SyncHandlers#process',
      message: 'Finished processing synchronous handlers',
      handlerCountMapping,
      batchProcessTime: Date.now() - start,
    });
    stats.timing(`${config.SERVICE_NAME}.synchronous_events_process_time`, Date.now() - start);
  }
}
