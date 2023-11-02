import { logger } from '@dydxprotocol-indexer/base';
import {
  MarketFromDatabase,
  MarketModel,
  MarketTable,
  marketRefresher,
  storeHelpers,
} from '@dydxprotocol-indexer/postgres';
import { MarketEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import config from '../../config';
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
    if (config.USE_MARKET_CREATE_HANDLER_SQL_FUNCTION) {
      return this.handleViaSqlFunction();
    }
    return this.handleViaKnexQueries();
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async handleViaKnexQueries(): Promise<ConsolidatedKafkaEvent[]> {
    // MarketHandler already makes sure the event has 'marketCreate' as the oneofKind.
    const marketCreate: MarketCreateEventMessage = this.event as MarketCreateEventMessage;

    const market: MarketFromDatabase | undefined = await MarketTable.findById(
      marketCreate.marketId,
    );
    if (market !== undefined) {
      this.logAndThrowParseMessageError(
        'Market in MarketCreate already exists',
        { marketCreate },
      );
    }
    await this.runFuncWithTimingStatAndErrorLogging(
      this.createMarket(marketCreate),
      this.generateTimingStatsOptions('create_market'),
    );
    return [];
  }

  private async handleViaSqlFunction(): Promise<ConsolidatedKafkaEvent[]> {
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

  private async createMarket(marketCreate: MarketCreateEventMessage): Promise<void> {
    await MarketTable.create({
      id: marketCreate.marketId,
      pair: marketCreate.marketCreate.base!.pair,
      exponent: marketCreate.marketCreate.exponent,
      minPriceChangePpm: marketCreate.marketCreate.base!.minPriceChangePpm,
    }, { txId: this.txId });
    await marketRefresher.updateMarkets({ txId: this.txId });
  }
}
