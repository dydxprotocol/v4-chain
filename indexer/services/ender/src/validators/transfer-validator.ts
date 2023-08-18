import { IndexerTendermintEvent, TransferEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { Handler } from '../handlers/handler';
import { TransferHandler } from '../handlers/transfer-handler';
import { Validator } from './validator';

export class TransferValidator extends Validator<TransferEventV1> {
  public validate(): void {}

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
  ): Handler<TransferEventV1>[] {
    return [
      new TransferHandler(
        this.block,
        indexerTendermintEvent,
        txId,
        this.event,
      ),
    ];
  }
}
