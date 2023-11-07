import { DeleveragingEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class DeleveragingHandler extends Handler<DeleveragingEventV1> {
  eventType: string = 'DeleveragingEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    // Implement this
    return [];
  }
}
