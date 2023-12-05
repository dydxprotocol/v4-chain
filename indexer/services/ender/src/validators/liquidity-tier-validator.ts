import { logger } from '@dydxprotocol-indexer/base';
import { IndexerTendermintEvent, LiquidityTierUpsertEventV1 } from '@dydxprotocol-indexer/v4-protos';
import Long from 'long';

import { Handler } from '../handlers/handler';
import { LiquidityTierHandler } from '../handlers/liquidity-tier-handler';
import { Validator } from './validator';

export class LiquidityTierValidator extends Validator<LiquidityTierUpsertEventV1> {
  public validate(): void {
    if (this.event.name === '') {
      logger.error({
        at: `${this.constructor.name}#validate`,
        message: 'LiquidityTierUpsertEventV1 name is not populated',
        blockHeight: this.block.height,
        event: this.event,
      });
    }

    if (this.event.basePositionNotional.eq(Long.fromValue(0))) {
      logger.error({
        at: `${this.constructor.name}#validate`,
        message: 'LiquidityTierUpsertEventV1 basePositionNotional is not populated',
        blockHeight: this.block.height,
        event: this.event,
      });
    }

    if (this.event.initialMarginPpm === 0) {
      return this.logAndThrowParseMessageError(
        'LiquidityTierUpsertEventV1 initialMarginPpm is not populated',
        { event: this.event },
      );
    }

    if (this.event.maintenanceFractionPpm === 0) {
      return this.logAndThrowParseMessageError(
        'LiquidityTierUpsertEventV1 maintenanceFractionPpm is not populated',
        { event: this.event },
      );
    }
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
  ): Handler<LiquidityTierUpsertEventV1>[] {
    const handler: Handler<LiquidityTierUpsertEventV1> = new LiquidityTierHandler(
      this.block,
      this.blockEventIndex,
      indexerTendermintEvent,
      txId,
      this.event,
    );

    return [handler];
  }
}
