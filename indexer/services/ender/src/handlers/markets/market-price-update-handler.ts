import { logger } from '@dydxprotocol-indexer/base';
import {
  MarketFromDatabase,
  OraclePriceFromDatabase,
  OraclePriceModel,
  MarketMessageContents, storeHelpers, MarketModel, marketRefresher,
} from '@dydxprotocol-indexer/postgres';
import { MarketEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { generateOraclePriceContents } from '../../helpers/kafka-helper';
import {
  ConsolidatedKafkaEvent,
  MarketPriceUpdateEventMessage,
} from '../../lib/types';
import { Handler } from '../handler';

export class MarketPriceUpdateHandler extends Handler<MarketEventV1> {
  eventType: string = 'MarketEvent';

  public getParallelizationIds(): string[] {
    // MarketEvents with the same market must be handled sequentially
    return [`${this.eventType}_${this.event.marketId}`];
  }

  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    logger.info({
      at: 'MarketPriceUpdateHandler#handle',
      message: 'Received MarketEvent with MarketPriceUpdate.',
      event: this.event,
    });

    const eventDataBinary: Uint8Array = this.indexerTendermintEvent.dataBytes;
    const result: pg.QueryResult = await storeHelpers.rawQuery(
      `SELECT dydx_market_price_update_handler(
        ${this.block.height}, 
        '${this.block.time?.toISOString()}', 
        '${JSON.stringify(MarketEventV1.decode(eventDataBinary))}' 
      ) AS result;`,
      { txId: this.txId },
    ).catch((error: Error) => {
      logger.error({
        at: 'MarketPriceUpdateHandler#internalHandle',
        message: 'Failed to handle MarketEventV1',
        error,
      });

      if (error.message.includes('MarketPriceUpdateEvent contains a non-existent market id')) {
        const castedMarketPriceUpdateMessage:
        MarketPriceUpdateEventMessage = this.event as MarketPriceUpdateEventMessage;
        this.logAndThrowParseMessageError(
          'MarketPriceUpdateEvent contains a non-existent market id',
          { castedMarketPriceUpdateMessage },
        );
      }

      throw error;
    });

    const market: MarketFromDatabase = MarketModel.fromJson(
      result.rows[0].result.market) as MarketFromDatabase;
    const oraclePrice: OraclePriceFromDatabase = OraclePriceModel.fromJson(
      result.rows[0].result.oracle_price) as OraclePriceFromDatabase;

    marketRefresher.updateMarket(market);

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
