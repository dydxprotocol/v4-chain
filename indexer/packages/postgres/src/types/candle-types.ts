import { IsoString } from './utility-types';

export interface CandleCreateObject {
  startedAt: IsoString,
  ticker: string,
  resolution: CandleResolution,
  low: string,
  high: string,
  open: string,
  close: string,
  baseTokenVolume: string,
  usdVolume: string,
  trades: number,
  startingOpenInterest: string,
  orderbookMidPriceOpen: string | undefined,
  orderbookMidPriceClose: string | undefined,
}

export interface CandleUpdateObject {
  id: string,
  low?: string,
  high?: string,
  open?: string,
  close?: string,
  baseTokenVolume?: string,
  usdVolume?: string,
  trades?: number,
  startingOpenInterest?: string,
  orderbookMidPriceOpen?: string,
  orderbookMidPriceClose?: string,
}

export enum CandleResolution {
  ONE_MINUTE = '1MIN',
  FIVE_MINUTES = '5MINS',
  FIFTEEN_MINUTES = '15MINS',
  THIRTY_MINUTES = '30MINS',
  ONE_HOUR = '1HOUR',
  FOUR_HOURS = '4HOURS',
  ONE_DAY = '1DAY',
}

export enum CandleColumns {
  id = 'id',
  startedAt = 'startedAt',
  ticker = 'ticker',
  resolution = 'resolution',
  low = 'low',
  high = 'high',
  open = 'open',
  close = 'close',
  baseTokenVolume = 'baseTokenVolume',
  usdVolume = 'usdVolume',
  trades = 'trades',
  startingOpenInterest = 'startingOpenInterest',
}
