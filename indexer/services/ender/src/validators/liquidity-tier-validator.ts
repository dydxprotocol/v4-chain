import { LiquidityTierUpsertEventV1, IndexerTendermintEvent } from '@dydxprotocol-indexer/v4-protos';

import { Handler } from '../handlers/handler';
import { LiquidityTierHandler } from '../handlers/liquidity-tier-handler';
import { Validator } from './validator';

export class LiquidityTierValidator extends Validator<LiquidityTierUpsertEventV1> {
  public validate(): void {}

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
  ): Handler<LiquidityTierUpsertEventV1>[] {
    const handler: Handler<LiquidityTierUpsertEventV1> = new LiquidityTierHandler(
      this.block,
      indexerTendermintEvent,
      txId,
      this.event,
    );

    return [handler];
  }
}
