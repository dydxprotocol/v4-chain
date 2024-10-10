/* ------- PERPETUAL MARKET TYPES ------- */

export interface PerpetualMarketCreateObject {
  id: string;
  clobPairId: string;
  ticker: string;
  marketId: number;
  status: PerpetualMarketStatus;
  priceChange24H: string;
  volume24H: string;
  trades24H: number;
  nextFundingRate: string;
  openInterest: string;
  quantumConversionExponent: number;
  atomicResolution: number;
  dangerIndexPpm: number;
  isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock: string;
  subticksPerTick: number;
  stepBaseQuantums: number;
  liquidityTierId: number;
  marketType: PerpetualMarketType;
  baseOpenInterest: string;
  perpYieldIndex: string;
}

export interface PerpetualMarketUpdateObject {
  id?: string;
  clobPairId?: string;
  ticker?: string;
  marketId?: number;
  status?: PerpetualMarketStatus;
  priceChange24H?: string;
  volume24H?: string;
  trades24H?: number;
  nextFundingRate?: string;
  openInterest?: string;
  quantumConversionExponent?: number;
  atomicResolution?: number;
  dangerIndexPpm?: number;
  isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock?: string;
  subticksPerTick?: number;
  stepBaseQuantums?: number;
  liquidityTierId?: number;
  perpYieldIndex?: string;
}

export enum PerpetualMarketColumns {
  id = 'id',
  clobPairId = 'clobPairId',
  ticker = 'ticker',
  marketId = 'marketId',
  status = 'status',
  priceChange24H = 'priceChange24H',
  volume24H = 'volume24H',
  trades24H = 'trades24H',
  nextFundingRate = 'nextFundingRate',
  openInterest = 'openInterest',
  quantumConversionExponent = 'quantumConversionExponent',
  atomicResolution = 'atomicResolution',
  dangerIndexPpm = 'dangerIndexPpm',
  isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock = 'isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock',
  subticksPerTick = 'subticksPerTick',
  stepBaseQuantums = 'stepBaseQuantums',
  liquidityTierId = 'liquidityTierId',
  perpYieldIndex = 'perpYieldIndex',
}

export enum PerpetualMarketStatus {
  ACTIVE = 'ACTIVE',
  PAUSED = 'PAUSED',
  CANCEL_ONLY = 'CANCEL_ONLY',
  POST_ONLY = 'POST_ONLY',
  INITIALIZING = 'INITIALIZING',
  FINAL_SETTLEMENT = 'FINAL_SETTLEMENT',
}

export enum PerpetualMarketType {
  CROSS = 'CROSS',
  ISOLATED = 'ISOLATED',
}
