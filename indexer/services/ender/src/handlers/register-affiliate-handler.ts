import { RegisterAffiliateEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class RegisterAffiliateHandler extends Handler<RegisterAffiliateEventV1> {
  eventType: string = 'RegisterAffiliateEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await, @typescript-eslint/no-unused-vars
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    return [];
  }
}
