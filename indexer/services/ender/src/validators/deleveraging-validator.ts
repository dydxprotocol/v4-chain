import { IndexerTendermintEvent, DeleveragingEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { DeleveragingHandler } from '../handlers/deleveraging-handler';
import { Handler } from '../handlers/handler';
import { Validator } from './validator';

export class DeleveragingValidator extends Validator<DeleveragingEventV1> {
  public validate(): void {
    if (this.event.liquidated === undefined) {
      return this.logAndThrowParseMessageError(
        'DeleveragingEvent must have a liquidated subaccount id',
        { event: this.event },
      );
    }

    if (this.event.offsetting === undefined) {
      return this.logAndThrowParseMessageError(
        'DeleveragingEvent must have an offsetting subaccount id',
        { event: this.event },
      );
    }

  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
  ): Handler<DeleveragingEventV1>[] {
    return [
      new DeleveragingHandler(
        this.block,
        indexerTendermintEvent,
        txId,
        this.event,
      ),
    ];
  }
}
