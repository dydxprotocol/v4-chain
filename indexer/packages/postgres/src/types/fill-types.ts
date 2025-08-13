/* ------- FILL TYPES ------- */

import { OrderSide } from './order-types';

export type Market24HourTradeVolumes = {
  clobPairId: string,
  trades24H: string,
  volume24H: string,
};

export enum Liquidity {
  TAKER = 'TAKER',
  MAKER = 'MAKER',
}

export enum FillType {
  // LIMIT is the fill type for a fill with a limit taker order.
  LIMIT = 'LIMIT',
  // LIQUIDATED is for the taker side of the fill where the subaccount was liquidated.
  // The subaccountId associated with this fill is the liquidated subaccount.
  LIQUIDATED = 'LIQUIDATED',
  // LIQUIDATION is for the maker side of the fill, never used for orders
  LIQUIDATION = 'LIQUIDATION',
  // DELEVERAGED is for the subaccount that was deleveraged in a deleveraging event.
  // The fill type will be set to taker.
  DELEVERAGED = 'DELEVERAGED',
  // OFFSETTING is for the offsetting subaccount in a deleveraging event.
  // The fill type will be set to maker.
  OFFSETTING = 'OFFSETTING',
  // TWAP is for a TWAP fill.
  TWAP_SUBORDER = 'TWAP_SUBORDER',
}

export interface FillCreateObject {
  subaccountId: string,
  side: OrderSide,
  liquidity: Liquidity,
  type: FillType,
  clobPairId: string,
  orderId?: string,
  size: string,
  price: string,
  quoteAmount: string,
  eventId: Buffer,
  transactionHash: string,
  createdAt: string,
  createdAtHeight: string,
  clientMetadata?: string,
  fee: string,
  affiliateRevShare: string,
  builderAddress?: string,
  builderFee?: string,
}

export interface FillUpdateObject {
  id: string,
  side?: OrderSide,
  type?: FillType,
  clobPairId?: string,
  orderId?: string | null,
  size?: string,
  price?: string,
  quoteAmount?: string,
}

export enum FillColumns {
  id = 'id',
  subaccountId = 'subaccountId',
  side = 'side',
  liquidity = 'liquidity',
  type = 'type',
  clobPairId = 'clobPairId',
  orderId = 'orderId',
  size = 'size',
  price = 'price',
  quoteAmount = 'quoteAmount',
  eventId = 'eventId',
  transactionHash = 'transactionHash',
  createdAt = 'createdAt',
  createdAtHeight = 'createdAtHeight',
  clientMetadata = 'clientMetadata',
  fee = 'fee',
  affiliateRevShare = 'affiliateRevShare',
  builderAddress = 'builderAddress',
  builderFee = 'builderFee',
}

export type CostOfFills = {
  cost: number,
};

export interface OrderedFillsWithFundingIndices {
  id: string,
  subaccountId: string,
  side: OrderSide,
  size: string,
  createdAtHeight: string,
  fundingIndex: string,
  lastFillId: string,
  lastFillSide: OrderSide,
  lastFillSize: string,
  lastFillCreatedAtHeight: string,
  lastFillFundingIndex: string,
}

export interface OpenSizeWithFundingIndex {
  clobPairId: string,
  openSize: string,
  lastFillHeight: string,
  fundingIndex: string,
  fundingIndexHeight: string,
}
