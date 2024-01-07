import {
  OrderFromDatabase,
  PerpetualMarketFromDatabase,
  orderTranslations,
  SubaccountFromDatabase, OrderModel, PerpetualMarketModel,
} from '@dydxprotocol-indexer/postgres';
import SubaccountModel from '@dydxprotocol-indexer/postgres/build/src/models/subaccount-model';
import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrder,
  OffChainUpdateV1,
  OrderPlaceV1_OrderPlacementStatus,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { ConsolidatedKafkaEvent } from '../../lib/types';
import { Handler } from '../handler';

export class ConditionalOrderTriggeredHandler extends Handler<StatefulOrderEventV1> {
  eventType: string = 'StatefulOrderEvent';

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const order: OrderFromDatabase = OrderModel.fromJson(resultRow.order) as OrderFromDatabase;
    const perpetualMarket: PerpetualMarketFromDatabase = PerpetualMarketModel.fromJson(
      resultRow.perpetual_market) as PerpetualMarketFromDatabase;
    const subaccount: SubaccountFromDatabase = SubaccountModel.fromJson(
      resultRow.subaccount) as SubaccountFromDatabase;

    const indexerOrder: IndexerOrder = orderTranslations.convertToIndexerOrderWithSubaccount(
      order, perpetualMarket, subaccount);
    return this.createKafkaEvents(indexerOrder);
  }

  private createKafkaEvents(order: IndexerOrder): ConsolidatedKafkaEvent[] {
    const offChainUpdate: OffChainUpdateV1 = OffChainUpdateV1.fromPartial({
      orderPlace: {
        order,
        placementStatus: OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
      },
    });

    return [
      this.generateConsolidatedVulcanKafkaEvent(
        getOrderIdHash(order.orderId!),
        offChainUpdate,
      ),
    ];
  }
}
