import { IndexerTendermintEvent, TransferEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { Handler } from '../handlers/handler';
import { TransferHandler } from '../handlers/transfer-handler';
import { Validator } from './validator';

export class TransferValidator extends Validator<TransferEventV1> {
  public validate(): void {
    // There must be exactly 1 sender. Either a subaccount id or a wallet address.
    if (
      (!this.event.sender) ||
      (this.event.sender!.subaccountId === undefined && this.event.sender!.address === undefined) ||
      (this.event.sender!.subaccountId !== undefined && this.event.sender!.address !== undefined)
    ) {
      return this.logAndThrowParseMessageError(
        'TransferEvent must have either a sender subaccount id or sender wallet address',
        { event: this.event },
      );
    }

    // There must be exactly 1 recipient. Either a subaccount id or a wallet address.
    if (
      (!this.event.recipient) ||
      (this.event.recipient!.subaccountId === undefined &&
        this.event.recipient!.address === undefined) ||
      (this.event.recipient!.subaccountId !== undefined &&
        this.event.recipient!.address !== undefined)
    ) {
      return this.logAndThrowParseMessageError(
        'TransferEvent must have either a recipient subaccount id or recipient wallet address',
        { event: this.event },
      );
    }

    if (
      this.event.recipient!.address !== undefined &&
      this.event.sender!.address !== undefined
    ) {
      return this.logAndThrowParseMessageError(
        'TransferEvent cannot have both a sender and recipient wallet address',
        { event: this.event },
      );
    }
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
    _: string,
  ): Handler<TransferEventV1>[] {
    return [
      new TransferHandler(
        this.block,
        this.blockEventIndex,
        indexerTendermintEvent,
        txId,
        this.event,
      ),
    ];
  }
}
