import { logger } from '@dydxprotocol-indexer/base';
import {
  MarketFromDatabase,
  MarketModel,
  marketRefresher,
  storeHelpers,
} from '@dydxprotocol-indexer/postgres';
import { MarketEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { ConsolidatedKafkaEvent, MarketCreateEventMessage } from '../../lib/types';
import { Handler } from '../handler';

export class MarketCreateHandler extends Handler<MarketEventV1> {
  eventType: string = 'MarketEvent';

  public getParallelizationIds(): string[] {
    // MarketEvents with the same market must be handled sequentially
    return [`${this.eventType}_${this.event.marketId}`];
  }

  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    logger.info({
      at: 'MarketCreateHandler#handle',
      message: 'Received MarketEvent with MarketCreate.',
      event: this.event,
    });

    const eventDataBinary: Uint8Array = this.indexerTendermintEvent.dataBytes;
    const result: pg.QueryResult = await storeHelpers.rawQuery(
      `SELECT dydx_market_create_handler(
        '${JSON.stringify(MarketEventV1.decode(eventDataBinary))}'
      ) AS result;`,
      { txId: this.txId },
    ).catch((error: Error) => {
      logger.error({
        at: 'MarketCreateHandler#handleViaSqlFunction',
        message: 'Failed to handle MarketEventV1',
        error,
      });

      if (error.message.includes('Market in MarketCreate already exists')) {
        const marketCreate: MarketCreateEventMessage = this.event as MarketCreateEventMessage;
        this.logAndThrowParseMessageError(
          'Market in MarketCreate already exists',
          { marketCreate },
        );
      }

      throw error;
    });

    const market: MarketFromDatabase = MarketModel.fromJson(
      result.rows[0].result.market) as MarketFromDatabase;
    marketRefresher.updateMarket(market);
    return [];
  }
}
