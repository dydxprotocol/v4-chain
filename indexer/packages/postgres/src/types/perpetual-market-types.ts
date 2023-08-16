/* ------- PERPETUAL MARKET TYPES ------- */

export interface PerpetualMarketCreateObject {
  id: string;
  clobPairId: string;
  ticker: string;
  marketId: number;
  status: PerpetualMarketStatus;
  baseAsset: string;
  quoteAsset: string;
  lastPrice: string;
  priceChange24H: string;
  volume24H: string;
  trades24H: number;
  nextFundingRate: string;
  basePositionSize: string;
  incrementalPositionSize: string;
  maxPositionSize: string;
  openInterest: string;
  quantumConversionExponent: number;
  atomicResolution: number;
  subticksPerTick: number;
  minOrderBaseQuantums: number;
  stepBaseQuantums: number;
  liquidityTierId: number;
}

export interface PerpetualMarketUpdateObject {
  id?: string;
  clobPairId?: string;
  ticker?: string;
  marketId?: number;
  status?: PerpetualMarketStatus;
  baseAsset?: string;
  quoteAsset?: string;
  lastPrice?: string;
  priceChange24H?: string;
  volume24H?: string;
  trades24H?: number;
  nextFundingRate?: string;
  basePositionSize?: string;
  incrementalPositionSize?: string;
  maxPositionSize?: string;
  openInterest?: string;
  quantumConversionExponent?: number;
  atomicResolution?: number;
  subticksPerTick?: number;
  minOrderBaseQuantums?: number;
  stepBaseQuantums?: number;
  liquidityTierId?: number;
}

export enum PerpetualMarketColumns {
  id = 'id',
  clobPairId = 'clobPairId',
  ticker = 'ticker',
  marketId = 'marketId',
  status = 'status',
  baseAsset = 'baseAsset',
  quoteAsset = 'quoteAsset',
  lastPrice = 'lastPrice',
  priceChange24H = 'priceChange24H',
  volume24H = 'volume24H',
  trades24H = 'trades24H',
  nextFundingRate = 'nextFundingRate',
  basePositionSize = 'basePositionSize',
  incrementalPositionSize = 'incrementalPositionSize',
  maxPositionSize = 'maxPositionSize',
  openInterest = 'openInterest',
  quantumConversionExponent = 'quantumConversionExponent',
  atomicResolution = 'atomicResolution',
  subticksPerTick = 'subticksPerTick',
  minOrderBaseQuantums = 'minOrderBaseQuantums',
  stepBaseQuantums = 'stepBaseQuantums',
  liquidityTierId = 'liquidityTierId',
}

export enum PerpetualMarketStatus {
  ACTIVE = 'ACTIVE',
  PAUSED = 'PAUSED',
  CANCEL_ONLY = 'CANCEL_ONLY',
  POST_ONLY = 'POST_ONLY',
}
