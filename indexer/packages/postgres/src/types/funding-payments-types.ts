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
}

export interface FundingPaymentsCreateObject {
  subaccountId: string;
  createdAt: string;
  createdAtHeight: string;
  perpetualId: string;
  ticker: string;
  oraclePrice: string;
  size: string;
  side: PositionSide;
  rate: string;
  payment: string;
}

export interface FundingPaymentsFromDatabase {
  subaccountId: string;
  createdAt: string;
  createdAtHeight: string;
  perpetualId: string;
  ticker: string;
  oraclePrice: string;
  size: string;
  side: PositionSide;
  rate: string;
  payment: string;
}

export interface FundingPaymentsQueryConfig {
  limit?: number;
  subaccountId?: string[];
  perpetualId?: string[];
  ticker?: string;
  createdAtHeight?: string;
  createdAt?: string;
  createdBeforeOrAtHeight?: string;
  createdBeforeOrAt?: string;
  createdOnOrAfterHeight?: string;
  createdOnOrAfter?: string;
} 