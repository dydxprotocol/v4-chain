import { logger } from '@dydxprotocol-indexer/base';
import {
  storeHelpers,
} from '@dydxprotocol-indexer/postgres';
import { MarketEventV1 } from '@dydxprotocol-indexer/v4-protos';

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

    const eventDataBinary: Uint8Array = this.indexerTendermintEvent.dataBytes;
    await storeHelpers.rawQuery(
      `SELECT dydx_market_modify_handler(
        '${JSON.stringify(MarketEventV1.decode(eventDataBinary))}' 
      ) AS result;`,
      { txId: this.txId },
    ).catch((error: Error) => {
      logger.error({
        at: 'MarketModifyHandler#internalHandle',
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

    return [];
  }
}
