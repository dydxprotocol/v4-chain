import { logger, stats } from '@dydxprotocol-indexer/base';
import _ from 'lodash';

import config from '../config';
import { Handler } from '../handlers/handler';
import { KafkaPublisher } from './kafka-publisher';
import { ConsolidatedKafkaEvent, DydxIndexerSubtypes, EventMessage } from './types';

// type alias for an array of handlers.
type HandlerBatch = Handler<EventMessage>[];
type HandlerBatchMap = Partial<{ [indexerSubtype in DydxIndexerSubtypes]: HandlerBatch}>;
const SynchSubtypes = [
  DydxIndexerSubtypes.MARKET,
  DydxIndexerSubtypes.ASSET,
];

/**
 * A class that processes handlers sequentially in the order specified in SynchSubtypes.
 * These events should be handled prior to any events in BatchedHandlers, and are used for
 * processing asset and market events.
 */
export class SyncHandlers {
  syncHandlers: HandlerBatchMap;
  initializationTime: number;

  constructor() {
    this.syncHandlers = {};
    this.initializationTime = Date.now();
  }

  /**
   * Adds a handler which contains an event to be processed. This function should be called in the
   * order in which the events should be processed. The handlers will be processed sequentially
   * by event subtype.
   *
   * @param indexerSubtype The event subtype
   * @param handler The handler to add to the batched handlers
   */
  public addHandler(
    indexerSubtype: DydxIndexerSubtypes,
    handler: Handler<EventMessage>,
  ): void {
    if (!SynchSubtypes.includes(indexerSubtype)) {
      return;
    }
    if (!this.syncHandlers[indexerSubtype]) {
      this.syncHandlers[indexerSubtype] = [];
    }
    // @ts-ignore
    this.syncHandlers[indexerSubtype].push(handler);
  }

  /**
   * Processes all handlers that were passed in through `addHandler` sequentially
   * in the order specified in SynchSubtypes. Adds events to the kafkaPublisher.
   */
  public async process(
    kafkaPublisher: KafkaPublisher,
  ): Promise<void> {
    const start: number = Date.now();
    const handlerCountMapping: { [key: string]: number } = {};
    for (const indexerSubtype of SynchSubtypes) {
      if (this.syncHandlers[indexerSubtype]) {
        const handlerBatch: HandlerBatch = this.syncHandlers[indexerSubtype] as HandlerBatch;
        const consolidatedKafkaEventGroup: ConsolidatedKafkaEvent[][] = await Promise.all(
          _.map(handlerBatch, (handler: Handler<EventMessage>) => {
            const handlerName: string = handler.constructor.name;
            if (!(handlerName in handlerCountMapping)) {
              handlerCountMapping[handlerName] = 0;
            }
            handlerCountMapping[handlerName] += 1;
            return handler.handle();
          }),
        );
        stats.timing(
          `${config.SERVICE_NAME}.synch_handlers.processing_delay.timing`,
          Date.now() - this.initializationTime,
          { eventType: indexerSubtype },
        );

        _.forEach(consolidatedKafkaEventGroup, (events: ConsolidatedKafkaEvent[]) => {
          kafkaPublisher.addEvents(events);
        });
      }
    }
    logger.info({
      at: 'SyncHandlers#process',
      message: 'Finished processing synchronous handlers',
      handlerCountMapping,
      batchProcessTime: Date.now() - start,
    });
    stats.timing(`${config.SERVICE_NAME}.synchronous_events_process_time`, Date.now() - start);
  }
}
