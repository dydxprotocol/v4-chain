/* ------- MARKET TYPES ------- */

export interface MarketCreateObject {
  id: number,
  pair: string,
  exponent: number,
  minPriceChangePpm: number,
  oraclePrice?: string,
}

export interface MarketUpdateObject {
  id: number,
  pair?: string,
  // Unsafe to update as it can lead to inaccurate price calculations for a market
  exponent?: number,
  minPriceChangePpm?: number,
  oraclePrice?: string;
}

export enum MarketColumns {
  id = 'id',
  pair = 'pair',
  exponent = 'exponent',
  minPriceChangePpm = 'minPriceChangePpm',
  oraclePrice = 'oraclePrice',
}
