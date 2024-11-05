import { UpdateYieldParamsEventV1 } from '@klyraprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { Handler } from './handler';
import { ConsolidatedKafkaEvent } from '../lib/types';

export class YieldParamsHandler extends Handler<UpdateYieldParamsEventV1> {
  eventType: string = 'UpdateYieldParamsEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  public async internalHandle(_resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    return Promise.resolve(this.generateKafkaEvents());
  }

  /** Generates a kafka websocket event for yieldParams.
   *
   * @param yieldParams
   * @protected
   */
  protected generateKafkaEvents(): ConsolidatedKafkaEvent[] {
    // TODO: [YBCP-28] Consider adding a websocket message for updated yield params
    return [];
  }
}
