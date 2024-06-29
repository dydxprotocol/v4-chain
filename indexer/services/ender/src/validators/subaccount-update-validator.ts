import { IndexerTendermintEvent, SubaccountUpdateEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { Handler } from '../handlers/handler';
import { SubaccountUpdateHandler } from '../handlers/subaccount-update-handler';
import { subaccountUpdateEventV1ToSubaccountUpdate } from '../helpers/translation-helper';
import { SubaccountUpdate } from '../lib/translated-types';
import { Validator } from './validator';

export class SubaccountUpdateValidator extends Validator<SubaccountUpdateEventV1> {
  public validate(): void {
    if (this.event.subaccountId === undefined) {
      this.logAndThrowParseMessageError(
        'SubaccountUpdateEvent must contain a subaccountId',
        { event: this.event },
      );
    }
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
    _: string,
  ): Handler<SubaccountUpdate>[] {
    return [
      new SubaccountUpdateHandler(
        this.block,
        this.blockEventIndex,
        indexerTendermintEvent,
        txId,
        subaccountUpdateEventV1ToSubaccountUpdate(this.event),
      ),
    ];
  }
}
