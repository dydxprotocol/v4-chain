import { logger } from '@dydxprotocol-indexer/base';
import {
  MarketFromDatabase,
  MarketUpdateObject,
  MarketTable,
  OraclePriceCreateObject,
  OraclePriceFromDatabase,
  OraclePriceModel,
  OraclePriceTable,
  protocolTranslations,
  MarketMessageContents, storeHelpers, MarketModel, marketRefresher,
} from '@dydxprotocol-indexer/postgres';
import { MarketEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { updatePriceCacheWithPrice } from '../../caches/price-cache';
import config from '../../config';
import { generateOraclePriceContents } from '../../helpers/kafka-helper';
import {
  ConsolidatedKafkaEvent,
  MarketPriceUpdateEventMessage,
} from '../../lib/types';
import { Handler } from '../handler';

type OraclePriceWithTicker = {
  oraclePrice: OraclePriceFromDatabase,
  pair: string,
};

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
    if (config.USE_MARKET_MODIFY_HANDLER_SQL_FUNCTION) {
      return this.handleViaSqlFunction();
    }
    return this.handleViaKnexQueries();
  }

  private async handleViaKnexQueries(): Promise<ConsolidatedKafkaEvent[]> {
    // MarketHandler already makes sure the event has 'priceUpdate' as the oneofKind.
    const castedMarketPriceUpdateMessage:
    MarketPriceUpdateEventMessage = this.event as MarketPriceUpdateEventMessage;

    const { oraclePrice, pair }:
    OraclePriceWithTicker = await this.runFuncWithTimingStatAndErrorLogging(
      this.createOraclePriceAndUpdateMarketFromEvent(castedMarketPriceUpdateMessage),
      this.generateTimingStatsOptions('create_and_update_oracle_prices'),
    );
    return [
      this.generateKafkaEvent(
        oraclePrice, pair,
      ),
    ];
  }

  private async handleViaSqlFunction(): Promise<ConsolidatedKafkaEvent[]> {
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
        at: 'MarketPriceUpdateHandler#handleViaSqlFunction',
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
    updatePriceCacheWithPrice(oraclePrice);

    return [
      this.generateKafkaEvent(
        oraclePrice, market.pair,
      ),
    ];
  }

  protected async updateMarketFromEvent(
    castedMarketPriceUpdateMessage: MarketPriceUpdateEventMessage,
    humanPrice: string,
  ): Promise<MarketFromDatabase> {

    const market: MarketFromDatabase | undefined = await MarketTable.findById(
      castedMarketPriceUpdateMessage.marketId,
      { txId: this.txId },
    );

    if (market === undefined) {
      this.logAndThrowParseMessageError(
        'Market in MarketPriceUpdateEventMessage doesn\'t exist',
        { castedMarketModifyMessage: castedMarketPriceUpdateMessage },
      );
    }

    const updateObject: MarketUpdateObject = {
      id: castedMarketPriceUpdateMessage.marketId,
      oraclePrice: humanPrice,
    };

    const updatedMarket:
    MarketFromDatabase | undefined = await MarketTable
      .update(updateObject, { txId: this.txId });
    if (updatedMarket === undefined) {
      this.logAndThrowParseMessageError(
        'Failed to update market in markets table',
        { castedMarketModifyMessage: castedMarketPriceUpdateMessage },
      );
    }
    return updatedMarket as MarketFromDatabase;
  }

  protected async createOraclePriceAndUpdateMarketFromEvent(
    castedMarketPriceUpdateMessage: MarketPriceUpdateEventMessage,
  ): Promise<{oraclePrice: OraclePriceFromDatabase, pair: string}> {
    const market: MarketFromDatabase | undefined = await MarketTable
      .findById(castedMarketPriceUpdateMessage.marketId, { txId: this.txId });
    if (market === undefined) {
      this.logAndThrowParseMessageError(
        'MarketPriceUpdateEvent contains a non-existent market id',
        { castedMarketPriceUpdateMessage },
      );
    }
    const humanPrice: string = protocolTranslations.protocolPriceToHuman(
      castedMarketPriceUpdateMessage.priceUpdate.priceWithExponent.toString(),
      market!.exponent,
    );
    await this.updateMarketFromEvent(castedMarketPriceUpdateMessage, humanPrice);
    const oraclePriceToCreate: OraclePriceCreateObject = {
      marketId: castedMarketPriceUpdateMessage.marketId,
      price: humanPrice,
      effectiveAt: this.timestamp.toISO(),
      effectiveAtHeight: this.block.height.toString(),
    };
    const oraclePriceFromDatabase: OraclePriceFromDatabase = await OraclePriceTable.create(
      oraclePriceToCreate,
      { txId: this.txId },
    );
    updatePriceCacheWithPrice(oraclePriceFromDatabase);
    return { oraclePrice: oraclePriceFromDatabase, pair: market!.pair };
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
