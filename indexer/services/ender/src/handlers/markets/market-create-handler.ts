import { logger } from '@dydxprotocol-indexer/base';
import { MarketFromDatabase, MarketTable, marketRefresher } from '@dydxprotocol-indexer/postgres';
import { MarketEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { ConsolidatedKafkaEvent, MarketCreateEventMessage } from '../../lib/types';
import { Handler } from '../handler';

export class MarketCreateHandler extends Handler<MarketEventV1> {
  eventType: string = 'MarketEvent';

  public getParallelizationIds(): string[] {
    // MarketEvents with the same market must be handled sequentially
    return [`${this.eventType}_${this.event.marketId}`];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    logger.info({
      at: 'MarketCreateHandler#handle',
      message: 'Received MarketEvent with MarketCreate.',
      event: this.event,
    });
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
