import { logger } from '@dydxprotocol-indexer/base';
import { MarketEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { ConsolidatedKafkaEvent } from '../../lib/types';
import { Handler } from '../handler';

export class MarketModifyHandler extends Handler<MarketEventV1> {
  eventType: string = 'MarketEvent';

  public getParallelizationIds(): string[] {
    // MarketEvents with the same market must be handled sequentially
    return [`${this.eventType}_${this.event.marketId}`];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(_: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    logger.info({
      at: 'MarketModifyHandler#handle',
      message: 'Received MarketEvent with MarketModify.',
      event: this.event,
    });

    return [];
  }
}
