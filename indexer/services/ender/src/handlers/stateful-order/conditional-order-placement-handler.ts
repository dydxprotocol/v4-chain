import {
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';

import { ConsolidatedKafkaEvent } from '../../lib/types';
import { AbstractStatefulOrderHandler } from '../abstract-stateful-order-handler';

// TODO(IND-334): Implement handler.
export class ConditionalOrderPlacementHandler extends
  AbstractStatefulOrderHandler<StatefulOrderEventV1> {
  eventType: string = 'StatefulOrderEvent';

  public getParallelizationIds(): string[] {
    // Implement parallelization ids
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    // Implement handling conditional order placements
    return [];
  }
}
