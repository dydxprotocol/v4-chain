import {
  createNotification,
  NotificationDynamicFieldKey,
  NotificationType,
  sendFirebaseMessage,
} from '@dydxprotocol-indexer/notifications';
import { MarketTable, OrderFromDatabase } from '@dydxprotocol-indexer/postgres';

export async function sendOrderFilledNotification(order: OrderFromDatabase) {
  const market = await MarketTable.findById(Number(order.clobPairId));
  if (!market) {
    throw new Error('sendOrderFilledNotification # Market not found');
  }
  const notification = createNotification(
    NotificationType.ORDER_FILLED,
    {
      [NotificationDynamicFieldKey.AMOUNT]: order.size.toString(),
      [NotificationDynamicFieldKey.MARKET]: market.pair,
      [NotificationDynamicFieldKey.AVERAGE_PRICE]: order.price,
    },
  );
  await sendFirebaseMessage(order.subaccountId, notification);
}
