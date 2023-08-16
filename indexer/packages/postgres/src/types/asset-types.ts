/* ------- ASSET TYPES ------- */

export interface AssetCreateObject {
  id: string,
  symbol: string,
  atomicResolution: number,
  hasMarket: boolean,
  marketId?: number,
}

export interface AssetUpdateObject {
  id: string,
  symbol?: string,
  atomicResolution?: number,
  hasMarket?: boolean,
  marketId?: number | null,
}

export enum AssetColumns {
  id = 'id',
  symbol = 'symbol',
  atomicResolution = 'atomicResolution',
  hasMarket = 'hasMarket',
  marketId = 'marketId',
}
