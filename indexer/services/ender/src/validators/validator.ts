import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import { IndexerTendermintBlock, IndexerTendermintEvent } from '@dydxprotocol-indexer/v4-protos';

import { Handler } from '../handlers/handler';
import { EventMessage } from '../lib/types';

export type ValidatorInitializer = new (
  event: EventMessage,
  block: IndexerTendermintBlock,
  eventBlockIndex: number,
) => Validator<EventMessage>;

export abstract class Validator<T extends object> {
  event: T;
  block: IndexerTendermintBlock;
  blockEventIndex: number;

  constructor(event: T, block: IndexerTendermintBlock, blockEventIndex: number) {
    this.event = event;
    this.blockEventIndex = blockEventIndex;
    this.block = block;
  }

  /**
   * Throws ParseMessageError and logs error if the event is invalid in any way.
   */
  public abstract validate(): void;

  /**
   * Returns the decoded event and any additional details the SQL function needs to process.
   */
  // TODO(IND-513): Convert handlers to have a 1-1 relationship with each validator by merging
  // the order fill handlers together into a single handler (and the respective SQL function)
  // and push this method down into the handler itself.
  public getEventForBlockProcessor(): Promise<T> {
    return Promise.resolve(this.event);
  }

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
    messageReceivedTimestamp: string,
  ): Handler<EventMessage>[];

  /**
   * Allows aribtrary logic to skip SQL processing for an event.
   * Defaults to no.
   * @returns
   */
  public shouldSkipSql(): boolean {
    return false;
  }

  /**
   * Allows arbitrary logic to skip handlers for an event.
   * Defaults to no.
   * @returns
   */
  public shouldSkipHandlers(): boolean {
    return false;
  }
}
