import { IPlaceOrder, Network, SocketClient } from '@dydxprotocol/v4-client-js';
import Long from 'long';

import { defaultOrder } from './constants';
import { OrderDetails } from './types';

export async function sleep(milliseconds: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, milliseconds));
}

export function createModifiedOrder(
  order: OrderDetails,
): IPlaceOrder {
  const modifiedOrder: IPlaceOrder = defaultOrder;
  modifiedOrder.clientId = Math.floor(Math.random() * 1000000000);
  modifiedOrder.goodTilBlock = 0;
  modifiedOrder.clobPairId = order.clobPairId;
  modifiedOrder.timeInForce = order.timeInForce;
  modifiedOrder.reduceOnly = false;
  modifiedOrder.orderFlags = order.orderFlags;
  modifiedOrder.side = order.side;
  modifiedOrder.quantums = Long.fromNumber(order.quantums);
  modifiedOrder.subticks = Long.fromNumber(order.subticks);
  return modifiedOrder;
}

export function connectAndValidateSocketClient(validateMessage: Function): void {
  const mySocket = new SocketClient(
    Network.local().indexerConfig,
    () => {},
    () => {},
    (message) => {
      if (typeof message.data === 'string') {
        const data = JSON.parse(message.data as string);
        validateMessage(data, mySocket);
      }
    },
  );
  mySocket.connect();
}
