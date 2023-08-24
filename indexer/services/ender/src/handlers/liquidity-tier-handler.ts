import { LiquidityTierUpsertEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class LiquidityTierHandler extends Handler<LiquidityTierUpsertEventV1> {
  eventType: string = 'LiquidityTierUpsertEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    // Implement this
    return [];
  }
}
