/* ------- ORACLE PRICE TYPES ------- */

export interface OraclePriceCreateObject {
  marketId: number,
  price: string,
  effectiveAt: string,
  effectiveAtHeight: string,
}

export enum OraclePriceColumns {
  id = 'id',
  marketId = 'marketId',
  price = 'price',
  effectiveAt = 'effectiveAt',
  effectiveAtHeight = 'effectiveAtHeight',
}
