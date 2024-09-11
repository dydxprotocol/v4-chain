import { IsoString } from './utility-types';

export interface TradingRewardCreateObject {
  address: string,
  blockTime: IsoString,
  blockHeight: string,
  amount: string,
}

export enum TradingRewardColumns {
  id = 'id',
  address = 'address',
  blockTime = 'blockTime',
  blockHeight = 'blockHeight',
  amount = 'amount',
}
