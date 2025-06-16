import { PositionSide } from './position-types';

export enum FundingPaymentsColumns {
  subaccountId = 'subaccountId',
  createdAt = 'createdAt',
  createdAtHeight = 'createdAtHeight',
  perpetualId = 'perpetualId',
  ticker = 'ticker',
  oraclePrice = 'oraclePrice',
  size = 'size',
  side = 'side',
  rate = 'rate',
  payment = 'payment',
  fundingIndex = 'fundingIndex',
}

export interface FundingPaymentsCreateObject {
  subaccountId: string,
  createdAt: string,
  createdAtHeight: string,
  perpetualId: string,
  ticker: string,
  oraclePrice: string,
  size: string,
  side: PositionSide,
  rate: string,
  payment: string,
  fundingIndex: string,
}
