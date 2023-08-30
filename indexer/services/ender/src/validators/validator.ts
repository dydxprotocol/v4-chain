import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import { IndexerTendermintBlock, IndexerTendermintEvent } from '@dydxprotocol-indexer/v4-protos';

import { Handler } from '../handlers/handler';
import { EventMessage } from '../lib/types';

export type ValidatorInitializer = new (
  event: EventMessage,
  block: IndexerTendermintBlock,
) => Validator<EventMessage>;

export abstract class Validator<T> {
  event: T;
  block: IndexerTendermintBlock;

  constructor(event: T, block: IndexerTendermintBlock) {
    this.event = event;
    this.block = block;
  }

  /**
   * Throws ParseMessageError and logs error if the event is invalid in any way.
   */
  public abstract validate(): void;

  protected logAndThrowParseMessageError(
    message: string,
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    additionalParameters: any = {},
  ): void {
    logger.error({
      at: `${this.constructor.name}#logAndThrowParseMessageError`,
      message,
      blockHeight: this.block.height,
      ...additionalParameters,
    });
    throw new ParseMessageError(message);
  }

  /**
   * Creates all relevant handlers for the Event
   */
  public abstract createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
  ): Handler<EventMessage>[];
}
