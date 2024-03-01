import { IndexerTendermintEvent, DeleveragingEventV2 } from '@dydxprotocol-indexer/v4-protos';

import { Handler } from '../handlers/handler';
import { DeleveragingHandlerV2 } from '../handlers/order-fills/deleveraging-handler-v2';
import { Validator } from './validator';

export class DeleveragingValidatorV2 extends Validator<DeleveragingEventV2> {
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

    if (this.event.price.eq(0)) {
      return this.logAndThrowParseMessageError(
        'DeleveragingEvent price cannot equal 0',
        { event: this.event },
      );
    }
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
  ): Handler<DeleveragingEventV2>[] {
    return [
      new DeleveragingHandlerV2(
        this.block,
        this.blockEventIndex,
        indexerTendermintEvent,
        txId,
        this.event,
      ),
    ];
  }
}
