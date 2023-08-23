import { perpetualMarketRefresher } from '@dydxprotocol-indexer/postgres';
import { PerpetualMarketCreateEventV1, IndexerTendermintEvent } from '@dydxprotocol-indexer/v4-protos';

import { Handler } from '../handlers/handler';
import { PerpetualMarketCreationHandler } from '../handlers/perpetual-market-handler';
import { Validator } from './validator';

export class PerpetualMarketValidator extends Validator<PerpetualMarketCreateEventV1> {
  public validate(): void {
    if (perpetualMarketRefresher.getPerpetualMarketFromId(this.event.id.toString()) !== undefined) {
      return this.logAndThrowParseMessageError(
        'PerpetualMarketCreateEvent id already exists',
        { event: this.event },
      );
    }
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
  ): Handler<PerpetualMarketCreateEventV1>[] {
    const handler: Handler<PerpetualMarketCreateEventV1> = new PerpetualMarketCreationHandler(
      this.block,
      indexerTendermintEvent,
      txId,
      this.event,
    );

    return [handler];
  }
}
