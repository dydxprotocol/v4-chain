import {
    YieldParamsFromDatabase,
    YieldParamsModel,
  } from '@dydxprotocol-indexer/postgres';

  import { UpdateYieldParamsEventV1 } from '@dydxprotocol-indexer/v4-protos';

  import _ from 'lodash';
  import * as pg from 'pg';
  
  import { ConsolidatedKafkaEvent } from '../lib/types';
  import { Handler } from './handler';
  
  export class YieldParamsHandler extends Handler<UpdateYieldParamsEventV1> {
    eventType: string = 'UpdateYieldParamsEvent';
  
    public getParallelizationIds(): string[] {
      return [];
    }

    public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
      const yieldParams: YieldParamsFromDatabase = YieldParamsModel.fromJson(
        resultRow.yield_params) as YieldParamsFromDatabase;
      return this.generateKafkaEvents(yieldParams);
    }

    /** Generates a kafka websocket event for yieldParams.
     *
     * @param yieldParams
     * @protected
     */
    protected generateKafkaEvents(
      yieldParams: YieldParamsFromDatabase,
    ): ConsolidatedKafkaEvent[] {
      // TODO: [YBCP-28] Consider adding a websocket message for updated yield params
      return [];
    }
}
  