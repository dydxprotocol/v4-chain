/* ------- ORDER TYPES ------- */

import { IsoString } from './utility-types';

export enum OrderSide {
  BUY = 'BUY',
  SELL = 'SELL',
}

export enum OrderStatus {
  OPEN = 'OPEN',
  FILLED = 'FILLED',
  CANCELED = 'CANCELED',
  BEST_EFFORT_CANCELED = 'BEST_EFFORT_CANCELED',
  UNTRIGGERED = 'UNTRIGGERED',
  ERROR = 'ERROR',
}

export enum OrderType {
  LIMIT = 'LIMIT',
  MARKET = 'MARKET',
  STOP_LIMIT = 'STOP_LIMIT',
  STOP_MARKET = 'STOP_MARKET',
  TRAILING_STOP = 'TRAILING_STOP',
  TAKE_PROFIT = 'TAKE_PROFIT',
  TAKE_PROFIT_MARKET = 'TAKE_PROFIT_MARKET',
  TWAP = 'TWAP',
  TWAP_SUBORDER = 'TWAP_SUBORDER',
}

export enum TimeInForce {
  // GTT represents Good-Til-Time, where an order will first match with existing orders on the book
  // and any remaining size will be added to the book as a maker order, which will expire at a
  // given expiry time.
  GTT = 'GTT',
  // FOK represents Fill-Or-KILl where it's enforced that an order will either be filled
  // completely and immediately by maker orders on the book or canceled if the entire amount can't
  // be filled.
  FOK = 'FOK',
  // IOC represents Immediate-Or-Cancel, where it's enforced that an order only be matched with
  // maker orders on the book. If the order has remaining size after matching with existing orders
  // on the book, the remaining size is not placed on the book.
  IOC = 'IOC',
  // POST_ONLY is where it's enforced that an order only be placed on the book as a maker order.
  POST_ONLY = 'POST_ONLY',
}

export interface OrderCreateObject {
  subaccountId: string,
  clientId: string,
  clobPairId: string,
  side: OrderSide,
  size: string,
  totalFilled: string,
  price: string,
  type: OrderType,
  status: OrderStatus,
  timeInForce: TimeInForce,
  reduceOnly: boolean,
  orderFlags: string,
  updatedAt: IsoString,
  updatedAtHeight: string,
  goodTilBlock?: string,
  goodTilBlockTime?: string,
  // createdAtHeight is optional because short term orders do not have a createdAtHeight.
  createdAtHeight?: string,
  clientMetadata: string,
  triggerPrice?: string,
  builderAddress?: string,
  feePpm?: string,
  orderRouterAddress?: string,
  duration?: string,
  interval?: string,
  priceTolerance?: string,
}

export interface OrderUpdateObject {
  id: string,
  clobPairId?: string,
  side?: OrderSide,
  size?: string,
  totalFilled?: string,
  price?: string,
  type?: OrderType,
  status?: OrderStatus,
  timeInForce?: TimeInForce,
  reduceOnly?: boolean,
  orderFlags?: string,
  updatedAt?: IsoString,
  updatedAtHeight?: string,
  goodTilBlock?: string | null,
  goodTilBlockTime?: string | null,
  clientMetadata?: string,
  triggerPrice?: string,
  orderRouterAddress?: string,
  duration?: string | null,
  interval?: string | null,
  priceTolerance?: string | null,
}

export enum OrderColumns {
  id = 'id',
  subaccountId = 'subaccountId',
  clientId = 'clientId',
  clobPairId = 'clobPairId',
  side = 'side',
  size = 'size',
  totalFilled = 'totalFilled',
  price = 'price',
  type = 'type',
  status = 'status',
  timeInForce = 'timeInForce',
  reduceOnly = 'reduceOnly',
  goodTilBlock = 'goodTilBlock',
  goodTilBlockTime = 'goodTilBlockTime',
  perpetualId = 'perpetualId',
  openEventId = 'openEventId',
  orderFlags = 'orderFlags',
  updatedAt = 'updatedAt',
  updatedAtHeight = 'updatedAtHeight',
  createdAtHeight = 'createdAtHeight',
  clientMetadata = 'clientMetadata',
  triggerPrice = 'triggerPrice',
  duration = 'duration',
  interval = 'interval',
  priceTolerance = 'priceTolerance',
}
