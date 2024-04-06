import {
  stats,
} from '@dydxprotocol-indexer/base';
import {
  BlockFromDatabase,
  BlockTable,
  OrderFromDatabase,
  OrderStatus,
  OrderTable,
  PaginationFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import { ORDER_FLAG_SHORT_TERM } from '@dydxprotocol-indexer/v4-proto-parser';

import config from '../config';

// TODO(IND-227): Change this task to a sanity check instead once canceled order ids are cached
// and the orders have their status correctly set
export default async function runTask(): Promise<void> {
  const queryStart: number = Date.now();
  const latestBlock: BlockFromDatabase = await BlockTable.getLatest({
    readReplica: true,
  });

  const latestBlockHeight: number = parseInt(latestBlock.blockHeight, 10);

  const { results: staleOpenOrders }: PaginationFromDatabase<OrderFromDatabase> = await OrderTable
    .findAll(
      {
        statuses: [OrderStatus.OPEN],
        orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
        // goodTilBlock needs to be < latest block height to be guaranteed to be CANCELED
        goodTilBlockBeforeOrAt: (latestBlockHeight - 1).toString(),
        limit: config.CANCEL_STALE_ORDERS_QUERY_BATCH_SIZE,
      },
      [],
      {
        readReplica: true,
      },
    );
  stats.timing(`${config.SERVICE_NAME}.cancel_stale_orders.query.timing`, Date.now() - queryStart);

  const updateStart: number = Date.now();
  const orderIds: string[] = staleOpenOrders.map((order: OrderFromDatabase) => order.id);
  const canceledOrders: OrderFromDatabase[] = await OrderTable.updateStaleOrderStatusByIds(
    OrderStatus.OPEN,
    OrderStatus.CANCELED,
    latestBlock.blockHeight,
    orderIds,
  );

  stats.timing(
    `${config.SERVICE_NAME}.cancel_stale_orders.update.timing`,
    Date.now() - updateStart,
  );
  stats.gauge(`${config.SERVICE_NAME}.num_stale_orders.count`, orderIds.length);
  stats.gauge(`${config.SERVICE_NAME}.num_stale_orders_canceled.count`, canceledOrders.length);
}
