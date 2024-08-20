import { logger } from '@dydxprotocol-indexer/base';
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
  TokenTable,
} from '@dydxprotocol-indexer/postgres';

export async function sendOrderFilledNotification(
  order: OrderFromDatabase,
  market: PerpetualMarketFromDatabase,
) {
  try {
    const subaccount = await SubaccountTable.findById(order.subaccountId);
    if (!subaccount) {
      throw new Error(`Subaccount not found for id ${order.subaccountId}`);
    }

    const tokens = (await TokenTable.findAll({ address: subaccount.address }, []))
      .map((token) => token.token);
    if (tokens.length === 0) {
      throw new Error(`No token found for address ${subaccount.address}`);
    }

    const notification = createNotification(
      NotificationType.ORDER_FILLED,
      {
        [NotificationDynamicFieldKey.AMOUNT]: order.size.toString(),
        [NotificationDynamicFieldKey.MARKET]: market.ticker,
        [NotificationDynamicFieldKey.AVERAGE_PRICE]: order.price,
      },
    );

    await sendFirebaseMessage(tokens, notification);
  } catch (error) {
    logger.error({
      at: 'ender#notification-functions',
      message: 'Error sending order filled notification',
      error,
    });
  }
}

export async function sendOrderTriggeredNotification(
  order: OrderFromDatabase,
  market: PerpetualMarketFromDatabase,
  subaccount: SubaccountFromDatabase,
) {
  try {
    const tokens = (await TokenTable.findAll({ address: subaccount.address }, []))
      .map((token) => token.token);
    if (tokens.length === 0) {
      throw new Error(`No tokens found for address ${subaccount.address}`);
    }
    const notification = createNotification(
      NotificationType.ORDER_TRIGGERED,
      {
        [NotificationDynamicFieldKey.MARKET]: market.ticker,
        [NotificationDynamicFieldKey.PRICE]: order.price,
        [NotificationDynamicFieldKey.AMOUNT]: order.size.toString(),
      },
    );
    await sendFirebaseMessage(tokens, notification);
  } catch (error) {
    logger.error({
      at: 'ender#notification-functions',
      message: 'Error sending order triggered notification',
      error,
    });
  }
}
