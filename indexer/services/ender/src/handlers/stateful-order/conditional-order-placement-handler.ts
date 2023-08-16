import { logger } from '@dydxprotocol-indexer/base';
import {
  OrderTable,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  protocolTranslations,
} from '@dydxprotocol-indexer/postgres';
import {
  IndexerOrder,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';

import { getTriggerPrice } from '../../lib/helper';
import { ConsolidatedKafkaEvent } from '../../lib/types';
import { AbstractStatefulOrderHandler } from '../abstract-stateful-order-handler';

// TODO(IND-334): Implement handler.
export class ConditionalOrderPlacementHandler extends
  AbstractStatefulOrderHandler<StatefulOrderEventV1> {
  eventType: string = 'StatefulOrderEvent';

  public getParallelizationIds(): string[] {
    const orderId: string = OrderTable.orderIdToUuid(
      this.event.conditionalOrderPlacement!.order!.orderId!,
    );
    return this.getParallelizationIdsFromOrderId(orderId);
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    const order: IndexerOrder = this.event.conditionalOrderPlacement!.order!;
    const clobPairId: string = order.orderId!.clobPairId.toString();
    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(clobPairId);
    if (perpetualMarket === undefined) {
      logger.error({
        at: 'conditionalOrderPlacementHandler#internalHandle',
        message: 'Unable to find perpetual market',
        clobPairId,
        order,
      });
      throw new Error(`Unable to find perpetual market with clobPairId: ${clobPairId}`);
    }

    await this.runFuncWithTimingStatAndErrorLogging(
      this.upsertOrder(
        perpetualMarket!,
        order,
        protocolTranslations.protocolConditionTypeToOrderType(order.conditionType),
        getTriggerPrice(order, perpetualMarket),
      ),
      this.generateTimingStatsOptions('upsert_order'),
    );
    return [];
  }
}
