import { logger } from '@dydxprotocol-indexer/base';
import { MarketFromDatabase, MarketTable } from '@dydxprotocol-indexer/postgres';
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
    // TODO(DEC-1752): Removed the height check once database seeding comes from V4 events.
    if (market !== undefined && this.block.height !== 0 && this.block.height !== 1) {
      this.logAndThrowParseMessageError(
        'Market in MarketCreate already exists',
        { marketCreate },
      );
    }
    if (market === undefined) {
      await this.runFuncWithTimingStatAndErrorLogging(
        MarketTable.create({
          id: marketCreate.marketId,
          pair: marketCreate.marketCreate.base!.pair,
          exponent: marketCreate.marketCreate.exponent,
          minPriceChangePpm: marketCreate.marketCreate.base!.minPriceChangePpm,
        }),
        this.generateTimingStatsOptions('create_market'),
      );
    }
    return [];
  }
}
