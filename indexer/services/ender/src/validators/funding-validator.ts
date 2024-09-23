import { logger } from '@dydxprotocol-indexer/base';
import { perpetualMarketRefresher } from '@dydxprotocol-indexer/postgres';
import {
  FundingEventV1,
  FundingEventV1_Type,
  IndexerTendermintEvent,
  FundingUpdateV1,
} from '@dydxprotocol-indexer/v4-protos';

import { FundingHandler } from '../handlers/funding-handler';
import { Handler } from '../handlers/handler';
import { FundingEventMessage } from '../lib/types';
import { Validator } from './validator';

export class FundingValidator extends Validator<FundingEventV1> {
  public validate(): void {
    if (
      this.event.type !== FundingEventV1_Type.TYPE_PREMIUM_SAMPLE &&
      this.event.type !== FundingEventV1_Type.TYPE_FUNDING_RATE_AND_INDEX
    ) {
      return this.logAndThrowParseMessageError(
        'Invalid FundingEvent, type must be TYPE_PREMIUM_SAMPLE or TYPE_FUNDING_RATE_AND_INDEX',
        { event: this.event },
      );
    }
    this.event.updates.forEach((fundingUpdate: FundingUpdateV1) => {
      const perpetualId = fundingUpdate.perpetualId;
      if (
        perpetualMarketRefresher.getPerpetualMarketFromId(
          perpetualId.toString(),
        ) === undefined) {
        logger.error({
          at: `${this.constructor.name}#validate`,
          message: 'Invalid FundingEvent, perpetualId does not exist',
          blockHeight: this.block.height,
          event: this.event,
        });
      }
    });
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
    _: string,
  ): Handler<FundingEventMessage>[] {
    const handler: Handler<FundingEventMessage> = new FundingHandler(
      this.block,
      this.blockEventIndex,
      indexerTendermintEvent,
      txId,
      this.event as FundingEventMessage,
    );

    return [handler];
  }
}
