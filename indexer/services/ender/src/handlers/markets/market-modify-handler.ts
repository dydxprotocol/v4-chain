import { logger } from '@dydxprotocol-indexer/base';
import {
  MarketFromDatabase, MarketUpdateObject, MarketTable, marketRefresher, storeHelpers, MarketModel,
} from '@dydxprotocol-indexer/postgres';
import { MarketEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import config from '../../config';
import { ConsolidatedKafkaEvent, MarketModifyEventMessage } from '../../lib/types';
import { Handler } from '../handler';

export class MarketModifyHandler extends Handler<MarketEventV1> {
  eventType: string = 'MarketEvent';

  public getParallelizationIds(): string[] {
    // MarketEvents with the same market must be handled sequentially
    return [`${this.eventType}_${this.event.marketId}`];
  }

  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    logger.info({
      at: 'MarketModifyHandler#handle',
      message: 'Received MarketEvent with MarketCreate.',
      event: this.event,
    });
    if (config.USE_MARKET_MODIFY_HANDLER_SQL_FUNCTION) {
      return this.handleViaSqlFunction();
    }
    return this.handleViaKnexQueries();
  }

  private async handleViaKnexQueries(): Promise<ConsolidatedKafkaEvent[]> {
    logger.info({
      at: 'MarketModifyHandler#handle',
      message: 'Received MarketEvent with MarketModify.',
      event: this.event,
    });
    // MarketHandler already makes sure the event has 'marketModify' as the oneofKind.
    const castedMarketModifyMessage:
    MarketModifyEventMessage = this.event as MarketModifyEventMessage;

    await this.runFuncWithTimingStatAndErrorLogging(
      this.updateMarketFromEvent(castedMarketModifyMessage),
      this.generateTimingStatsOptions('update_market'),
    );
    return [];
  }

  private async handleViaSqlFunction(): Promise<ConsolidatedKafkaEvent[]> {
    const eventDataBinary: Uint8Array = this.indexerTendermintEvent.dataBytes;
    const result: pg.QueryResult = await storeHelpers.rawQuery(
      `SELECT dydx_market_modify_handler(
        '${JSON.stringify(MarketEventV1.decode(eventDataBinary))}' 
      ) AS result;`,
      { txId: this.txId },
    ).catch((error: Error) => {
      logger.error({
        at: 'MarketModifyHandler#handleViaSqlFunction',
        message: 'Failed to handle MarketEventV1',
        error,
      });

      const castedMarketModifyMessage:
      MarketModifyEventMessage = this.event as MarketModifyEventMessage;

      if (error.message.includes('Market in MarketModify doesn\'t exist')) {
        this.logAndThrowParseMessageError(
          'Market in MarketModify doesn\'t exist',
          { castedMarketModifyMessage },
        );
      }

      this.logAndThrowParseMessageError(
        'Failed to update market in markets table',
        { castedMarketModifyMessage },
      );
    });

    const market: MarketFromDatabase = MarketModel.fromJson(
      result.rows[0].result.market) as MarketFromDatabase;
    marketRefresher.updateMarket(market);
    return [];
  }

  protected async updateMarketFromEvent(
    castedMarketModifyMessage: MarketModifyEventMessage,
  ): Promise<MarketFromDatabase> {

    const market: MarketFromDatabase | undefined = await MarketTable.findById(
      castedMarketModifyMessage.marketId,
    );
    if (market === undefined) {
      this.logAndThrowParseMessageError(
        'Market in MarketModify doesn\'t exist',
        { castedMarketModifyMessage },
      );
    }

    const updateObject: MarketUpdateObject = {
      id: castedMarketModifyMessage.marketId,
      pair: castedMarketModifyMessage.marketModify.base!.pair!,
      minPriceChangePpm: castedMarketModifyMessage.marketModify.base!.minPriceChangePpm!,
    };

    const updatedMarket:
    MarketFromDatabase | undefined = await MarketTable
      .update(updateObject, { txId: this.txId });
    if (updatedMarket === undefined) {
      this.logAndThrowParseMessageError(
        'Failed to update market in markets table',
        { castedMarketModifyMessage },
      );
    }
    await marketRefresher.updateMarkets({ txId: this.txId });
    return updatedMarket as MarketFromDatabase;
  }
}
