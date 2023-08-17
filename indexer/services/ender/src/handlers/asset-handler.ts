import { AssetCreateEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class AssetCreationHandler extends Handler<AssetCreateEventV1> {
  eventType: string = 'AssetCreateEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    // Implement this
    return [];
  }
}
