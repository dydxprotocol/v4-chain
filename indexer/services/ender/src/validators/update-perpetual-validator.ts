import { perpetualMarketRefresher } from '@klyraprotocol-indexer/postgres';
import { IndexerTendermintEvent, UpdatePerpetualEventV1 } from '@klyraprotocol-indexer/v4-protos';

import { Validator } from './validator';
import { Handler } from '../handlers/handler';
import { UpdatePerpetualHandler } from '../handlers/update-perpetual-handler';

export class UpdatePerpetualValidator extends Validator<UpdatePerpetualEventV1> {
  public validate(): void {
    if (perpetualMarketRefresher.getPerpetualMarketFromId(this.event.id.toString()) === undefined) {
      return this.logAndThrowParseMessageError(
        'UpdatePerpetualEvent.id must correspond with an existing perpetual_market.id',
        { event: this.event },
      );
    }
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
    _: string,
  ): Handler<UpdatePerpetualEventV1>[] {
    const handler: Handler<UpdatePerpetualEventV1> = new UpdatePerpetualHandler(
      this.block,
      this.blockEventIndex,
      indexerTendermintEvent,
      txId,
      this.event,
    );

    return [handler];
  }
}
