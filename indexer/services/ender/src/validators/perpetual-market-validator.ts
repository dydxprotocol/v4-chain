import { perpetualMarketRefresher } from '@dydxprotocol-indexer/postgres';
import { PerpetualMarketCreateEventV1, IndexerTendermintEvent, PerpetualMarketCreateEventV2 } from '@dydxprotocol-indexer/v4-protos';
import Long from 'long';

import { Handler } from '../handlers/handler';
import { PerpetualMarketCreationHandler } from '../handlers/perpetual-market-handler';
import { Validator } from './validator';

export class PerpetualMarketValidator extends Validator<
  PerpetualMarketCreateEventV1 | PerpetualMarketCreateEventV2
> {
  public validate(): void {
    if (perpetualMarketRefresher.getPerpetualMarketFromId(this.event.id.toString()) !== undefined) {
      return this.logAndThrowParseMessageError(
        'PerpetualMarketCreateEvent id already exists',
        { event: this.event },
      );
    }
    if (this.event.ticker === '') {
      return this.logAndThrowParseMessageError(
        'PerpetualMarketCreateEvent ticker is not populated',
        { event: this.event },
      );
    }
    if (this.event.subticksPerTick === 0) {
      return this.logAndThrowParseMessageError(
        'PerpetualMarketCreateEvent subticksPerTick is not populated',
        { event: this.event },
      );
    }

    if (this.event.stepBaseQuantums.eq(Long.fromValue(0))) {
      return this.logAndThrowParseMessageError(
        'PerpetualMarketCreateEvent stepBaseQuantums is not populated',
        { event: this.event },
      );
    }
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
    _: string,
  ): Handler<PerpetualMarketCreateEventV1>[] {
    const handler: Handler<PerpetualMarketCreateEventV1> = new PerpetualMarketCreationHandler(
      this.block,
      this.blockEventIndex,
      indexerTendermintEvent,
      txId,
      this.event,
    );

    return [handler];
  }
}
