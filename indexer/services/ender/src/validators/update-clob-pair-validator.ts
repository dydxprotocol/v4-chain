import { perpetualMarketRefresher } from '@klyraprotocol-indexer/postgres';
import { IndexerTendermintEvent, UpdateClobPairEventV1 } from '@klyraprotocol-indexer/v4-protos';

import { Validator } from './validator';
import { Handler } from '../handlers/handler';
import { UpdateClobPairHandler } from '../handlers/update-clob-pair-handler';

export class UpdateClobPairValidator extends Validator<UpdateClobPairEventV1> {
  public validate(): void {
    if (perpetualMarketRefresher.getPerpetualMarketFromClobPairId(
      this.event.clobPairId.toString(),
    ) === undefined) {
      return this.logAndThrowParseMessageError(
        'UpdateClobPairEvent.clobPairId must correspond with an existing perpetual_market.clobPairId',
        { event: this.event },
      );
    }
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
    _: string,
  ): Handler<UpdateClobPairEventV1>[] {
    const handler: Handler<UpdateClobPairEventV1> = new UpdateClobPairHandler(
      this.block,
      this.blockEventIndex,
      indexerTendermintEvent,
      txId,
      this.event,
    );

    return [handler];
  }
}
