import { IndexerTendermintEvent, DeleveragingEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { Handler } from '../handlers/handler';
import { DeleveragingHandler } from '../handlers/order-fills/deleveraging-handler';
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

    if (this.event.fillAmount.eq(0)) {
      return this.logAndThrowParseMessageError(
        'DeleveragingEvent fillAmount cannot equal 0',
        { event: this.event },
      );
    }

    if (this.event.totalQuoteQuantums.eq(0)) {
      return this.logAndThrowParseMessageError(
        'DeleveragingEvent totalQuoteQuantums cannot equal 0',
        { event: this.event },
      );
    }
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
    _: string,
  ): Handler<DeleveragingEventV1>[] {
    return [
      new DeleveragingHandler(
        this.block,
        this.blockEventIndex,
        indexerTendermintEvent,
        txId,
        this.event,
      ),
    ];
  }
}
