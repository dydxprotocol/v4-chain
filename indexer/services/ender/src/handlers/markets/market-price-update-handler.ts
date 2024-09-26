import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  MarketFromDatabase,
  OraclePriceFromDatabase,
  OraclePriceModel,
  MarketMessageContents, MarketModel,
} from '@dydxprotocol-indexer/postgres';
import { MarketEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import config from '../../config';
import { generateOraclePriceContents } from '../../helpers/kafka-helper';
import {
  ConsolidatedKafkaEvent,
} from '../../lib/types';
import { Handler } from '../handler';

export class MarketPriceUpdateHandler extends Handler<MarketEventV1> {
  eventType: string = 'MarketEvent';

  public getParallelizationIds(): string[] {
    // MarketEvents with the same market must be handled sequentially
    return [`${this.eventType}_${this.event.marketId}`];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    logger.info({
      at: 'MarketPriceUpdateHandler#handle',
      message: 'Received MarketEvent with MarketPriceUpdate.',
      event: this.event,
    });

    const market: MarketFromDatabase = MarketModel.fromJson(
      resultRow.market) as MarketFromDatabase;
    const oraclePrice: OraclePriceFromDatabase = OraclePriceModel.fromJson(
      resultRow.oracle_price) as OraclePriceFromDatabase;

    // Handle latency from resultRow
    stats.timing(
      `${config.SERVICE_NAME}.handle_market_price_update_event.sql_latency`,
      Number(resultRow.latency),
      this.generateTimingStatsOptions(),
    );

    return [
      this.generateKafkaEvent(
        oraclePrice, market.pair,
      ),
    ];
  }

  protected generateKafkaEvent(
    oraclePrice: OraclePriceFromDatabase,
    pair: string,
  ): ConsolidatedKafkaEvent {
    const contents: MarketMessageContents = generateOraclePriceContents(
      oraclePrice, pair,
    );

    return this.generateConsolidatedMarketKafkaEvent(
      JSON.stringify(contents),
    );
  }
}
