import { logger } from '@dydxprotocol-indexer/base';

import { Handler } from '../handlers/handler';
import {DydxIndexerSubtypes, EventMessage, EventProtoWithTypeAndVersion} from './types';

export const BULK_UPDATE_SUBTYPES: DydxIndexerSubtypes[] = [
  DydxIndexerSubtypes.TRANSFER,
];

export const BULK_UPDATE_SUBTYPE_TO_HANDLER: Record<DydxIndexerSubtypes, Handler> = {

}

export class BulkUpdateHandlers {

  handlerBatch: <EventMessage>[] = [];

  public addEvent(
    event: EventProtoWithTypeAndVersion,
  ): void {
    const indexerSubtype: DydxIndexerSubtypes = event.type;
    if (!BULK_UPDATE_SUBTYPES.includes(indexerSubtype)) {
      logger.error({
        at: 'BulkUpdateHandlers#addHandler',
        message: `Invalid indexerSubtype: ${indexerSubtype}`,
      });
      return;
    }
    // @ts-ignore
    this.handlerBatch.push(handler);
  }

  public process(): void {

  }

}
