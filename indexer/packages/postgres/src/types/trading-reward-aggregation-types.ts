import { IsoString } from './utility-types';

export enum TradingRewardAggregationPeriod {
  DAILY = 'DAILY',
  WEEKLY = 'WEEKLY',
  MONTHLY = 'MONTHLY',
}

export interface TradingRewardAggregationCreateObject {
  address: string,
  startedAt: IsoString,
  startedAtHeight: string,
  endedAt?: IsoString,
  endedAtHeight?: string,
  period: TradingRewardAggregationPeriod,
  amount: string,
}

export interface TradingRewardAggregationUpdateObject {
  id: string,
  endedAt?: IsoString,
  endedAtHeight?: string,
  amount?: string,
}

export enum TradingRewardAggregationColumns {
  id = 'id',
  address = 'address',
  startedAt = 'startedAt',
  startedAtHeight = 'startedAtHeight',
  endedAt = 'endedAt',
  endedAtHeight = 'endedAtHeight',
  period = 'period',
  amount = 'amount',
}
