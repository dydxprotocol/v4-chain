import {
  OrderFromDatabase,
  OrderTable,
  PerpetualMarketFromDatabase,
  SubaccountFromDatabase,
  SubaccountMessageContents,
} from '@dydxprotocol-indexer/postgres';
import {
  IndexerSubaccountId,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';

import { generateOrderSubaccountMessage } from '../../helpers/kafka-helper';
import { ConsolidatedKafkaEvent } from '../../lib/types';
import { AbstractStatefulOrderHandler } from '../abstract-stateful-order-handler';

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
    const result:
    [OrderFromDatabase,
      PerpetualMarketFromDatabase,
      SubaccountFromDatabase | undefined] = await this.handleEventViaSqlFunction();

    const subaccountId:
    IndexerSubaccountId = this.event.conditionalOrderPlacement!.order!.orderId!.subaccountId!;
    return this.createKafkaEvents(subaccountId, result[0], result[1]);
  }

  private createKafkaEvents(
    subaccountId: IndexerSubaccountId,
    conditionalOrder: OrderFromDatabase,
    perpetualMarket: PerpetualMarketFromDatabase): ConsolidatedKafkaEvent[] {

    // Since the order isn't placed on the book, no message is sent to vulcan
    // ender needs to send the websocket message indicating the conditional order was placed
    const message: SubaccountMessageContents = {
      orders: [
        generateOrderSubaccountMessage(conditionalOrder, perpetualMarket.ticker),
      ],
    };

    return [
      this.generateConsolidatedSubaccountKafkaEvent(
        JSON.stringify(message),
        subaccountId,
      ),
    ];
  }
}
