/* ------- MARKET TYPES ------- */

export interface MarketCreateObject {
  id: number,
  pair: string,
  exponent: number,
  minPriceChangePpm: number,
  spotPrice?: string,
  pnlPrice?: string,
}

export interface MarketUpdateObject {
  id: number,
  pair?: string,
  minPriceChangePpm?: number,
  spotPrice?: string;
  pnlPrice?: string;
}

export enum MarketColumns {
  id = 'id',
  pair = 'pair',
  exponent = 'exponent',
  minPriceChangePpm = 'minPriceChangePpm',
  spotPrice = 'spotPrice',
  pnlPrice = 'pnlPrice',
}
