import { IndexerTendermintEvent, RegisterAffiliateEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { Handler } from '../handlers/handler';
import { RegisterAffiliateHandler } from '../handlers/register-affiliate-handler';
import { Validator } from './validator';

export class RegisterAffiliateValidator extends Validator<RegisterAffiliateEventV1> {
  public validate(): void {
    if (this.event.affiliate === null || this.event.referee === null) {
      return this.logAndThrowParseMessageError(
        'RegisterAffiliateEventV1 contains null values for fields',
        { event: this.event },
      );
    }
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
    _: string,
  ): Handler<RegisterAffiliateEventV1>[] {
    const handler: Handler<RegisterAffiliateEventV1> = new RegisterAffiliateHandler(
      this.block,
      this.blockEventIndex,
      indexerTendermintEvent,
      txId,
      this.event,
    );

    return [handler];
  }
}
