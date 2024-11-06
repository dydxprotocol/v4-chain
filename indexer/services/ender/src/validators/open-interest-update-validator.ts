import { IndexerTendermintEvent, OpenInterestUpdateEventV1 } from '@klyraprotocol-indexer/v4-protos';

import { Validator } from './validator';
import { Handler } from '../handlers/handler';
import { OpenInterestUpdateHandler } from '../handlers/open-interest-update-handler';

export class OpenInterestUpdateValidator extends Validator<OpenInterestUpdateEventV1> {
  public validate(): void {
    if (this.event.openInterestUpdates === null) {
      return this.logAndThrowParseMessageError(
        'OpenInterestUpdateEventV1 openInterestUpdates is not populated',
        { event: this.event },
      );
    }
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
    _: string,
  ): Handler<OpenInterestUpdateEventV1>[] {
    const handler: Handler<OpenInterestUpdateEventV1> = new OpenInterestUpdateHandler(
      this.block,
      this.blockEventIndex,
      indexerTendermintEvent,
      txId,
      this.event,
    );

    return [handler];
  }
}
