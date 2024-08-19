import {
  createNotification,
  NotificationDynamicFieldKey,
  NotificationType,
  sendFirebaseMessage,
} from '@dydxprotocol-indexer/notifications';
import {
  OrderFromDatabase,
  PerpetualMarketFromDatabase,
  SubaccountFromDatabase,
  SubaccountTable,
} from '@dydxprotocol-indexer/postgres';

export async function sendOrderFilledNotification(
  order: OrderFromDatabase,
  market: PerpetualMarketFromDatabase,
) {
  const subaccount = await SubaccountTable.findById(order.subaccountId);
  if (!subaccount) {
    throw new Error(`Subaccount not found for id ${order.subaccountId}`);
  }
  const notification = createNotification(
    NotificationType.ORDER_FILLED,
    {
      [NotificationDynamicFieldKey.AMOUNT]: order.size.toString(),
      [NotificationDynamicFieldKey.MARKET]: market.ticker,
      [NotificationDynamicFieldKey.AVERAGE_PRICE]: order.price,
    },
  );
  await sendFirebaseMessage(subaccount.address, notification);
}

export async function sendOrderTriggeredNotification(
  order: OrderFromDatabase,
  market: PerpetualMarketFromDatabase,
  subaccount: SubaccountFromDatabase,
) {
  const notification = createNotification(
    NotificationType.ORDER_TRIGGERED,
    {
      [NotificationDynamicFieldKey.MARKET]: market.ticker,
      [NotificationDynamicFieldKey.PRICE]: order.price,
      [NotificationDynamicFieldKey.AMOUNT]: order.size.toString(),
    },
  );
  await sendFirebaseMessage(subaccount.address, notification);
}
