import { stats } from '@dydxprotocol-indexer/base';
import _ from 'lodash';
import * as pg from 'pg';

import config from '../config';
import { Handler } from '../handlers/handler';
import { KafkaPublisher } from './kafka-publisher';
import { ConsolidatedKafkaEvent, EventMessage } from './types';

// type alias for an array of handlers, handlers in a `HandlerBatch` can be processed in parallel.
type HandlerBatch = Handler<EventMessage>[];

// TODO(IND-514): Remove the batch and sync handlers completely by moving all redis updates into
// a pipeline similar to how we return kafka events and then batch and emit them.
export class BatchedHandlers {
  // An array of `HandlerBatch`s. Handlers in a `HandlerBatch` can be processed in parallel, and
  // `HandlerBatch`s are processed in a sequential order following the order in `batchedHandlers`.
  batchedHandlers: HandlerBatch[];
  // An array of sets of parallization ids. Each array of ids is the parallelization ids for the
  // corresponding array of handlers in this.batchedHandlers as well as all the parallelization ids
  // of the previous batches. Handlers with overlapping parallelization ids cannot be processed in
  // parallel.
  pIdBatches: Set<string>[];
  initializationTime: number;

  constructor() {
    this.batchedHandlers = [];
    this.pIdBatches = [];
    this.initializationTime = Date.now();
  }

  /**
   * Adds a handler which contains an event to be processed. This function should be called in the
   * order in which the events should be processed. The handlers will be added to batches in which
   * all handlers in the batch can be processed in parallel by using the parallelization ids.
   *
   * The handler will be added to the first batch that does not contain any of the parallelization
   * ids. The parallelization ids of the handler will be added to the batch the handler is added to
   * and all batches before it, because all handlers with overlapping parallelization ids cannot be
   * added to any previous batches.
   * @param handler The handler to add to the batched handlers
   */
  public addHandler(handler: Handler<EventMessage>): void {
    const pIds: string[] = handler.getParallelizationIds();
    let createNewBatch: boolean = true;
    for (let batchIndex: number = 0; batchIndex < this.batchedHandlers.length; batchIndex++) {
      const pIdBatch: Set<string> = this.pIdBatches[batchIndex];
      const arePIdsInBatch: boolean = this.idsInSet(pIds, pIdBatch);
      this.addPIdsToSet(pIds, pIdBatch);

      if (!arePIdsInBatch) {
        this.batchedHandlers[batchIndex].push(handler);
        createNewBatch = false;
        break;
      }
    }

    if (createNewBatch) {
      this.batchedHandlers.push([handler]);
      this.pIdBatches.push(new Set(pIds));
    }
  }

  private idsInSet(pIds: string[], set: Set<string>): boolean {
    return pIds.some((pId: string) => set.has(pId));
  }

  private addPIdsToSet(pIds: string[], set: Set<string>): void {
    pIds.map((pId: string) => set.add(pId));
  }

  /**
   * Processes all handlers that were passed in through `addHandler` parallelizing handlers.handle
   * and ensuring that handlers with overlapping parallelization ids are not processed in parallel.
   * Adds events to the kafkaPublisher.
   */
  public async process(
    kafkaPublisher: KafkaPublisher,
    resultRow: pg.QueryResultRow,
  ): Promise<void> {
    for (let batchIndex = 0; batchIndex < this.batchedHandlers.length; batchIndex++) {
      const start: number = Date.now();
      const handlerCountMapping: { [key: string]: number } = {};
      const consolidatedKafkaEventGroup: ConsolidatedKafkaEvent[][] = await Promise.all(
        _.map(this.batchedHandlers[batchIndex], (handler: Handler<EventMessage>) => {
          const handlerName: string = handler.constructor.name;
          if (!(handlerName in handlerCountMapping)) {
            handlerCountMapping[handlerName] = 0;
          }
          handlerCountMapping[handlerName] += 1;
          stats.timing(
            `${config.SERVICE_NAME}.batched_handlers.processing_delay.timing`,
            Date.now() - this.initializationTime,
            { eventType: handler.eventType },
          );
          return handler.handle(resultRow[handler.blockEventIndex]);
        }),
      );

      _.forEach(consolidatedKafkaEventGroup, (events: ConsolidatedKafkaEvent[]) => {
        kafkaPublisher.addEvents(events);
      });
      stats.timing(`${config.SERVICE_NAME}.batch_process_time`, Date.now() - start);
      stats.histogram(`${config.SERVICE_NAME}.batch_size`, this.batchedHandlers[batchIndex].length);
    }
    stats.histogram(`${config.SERVICE_NAME}.num_batches_in_block`, this.batchedHandlers.length);
  }
}
