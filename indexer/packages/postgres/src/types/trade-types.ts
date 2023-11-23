import { FillType } from './fill-types';

export enum TradeType {
  // LIMIT is the trade type for a fill with a limit taker order.
  LIMIT = 'LIMIT',
  // LIQUIDATED is the trade type for a fill with a liquidated taker order.
  LIQUIDATED = 'LIQUIDATED',
  // DELEVERAGED is the trade type for a fill with a deleveraged taker order.
  DELEVERAGED = 'DELEVERAGED',
}

export function fillTypeToTradeType(fillType: FillType): TradeType {
  switch (fillType) {
    case FillType.LIMIT:
      return TradeType.LIMIT;
    case FillType.LIQUIDATED:
      return TradeType.LIQUIDATED;
    case FillType.DELEVERAGED:
      return TradeType.DELEVERAGED;
    default:
      throw new Error(`Unknown fill type: ${fillType}`);
  }
}
