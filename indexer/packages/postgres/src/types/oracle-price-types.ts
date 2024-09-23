/* ------- ORACLE PRICE TYPES ------- */

export interface OraclePriceCreateObject {
  marketId: number,
  spotPrice: string,
  pnlPrice: string,
  effectiveAt: string,
  effectiveAtHeight: string,
}

export enum OraclePriceColumns {
  id = 'id',
  marketId = 'marketId',
  spotPrice = 'spotPrice',
  pnlPrice = 'pnlPrice',
  effectiveAt = 'effectiveAt',
  effectiveAtHeight = 'effectiveAtHeight',
}
