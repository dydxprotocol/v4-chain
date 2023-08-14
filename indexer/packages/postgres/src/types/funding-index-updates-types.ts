/* ------- FUNDING INDEX TYPES ------- */

export interface FundingIndexUpdatesCreateObject {
  perpetualId: string,
  eventId: Buffer,
  rate: string,
  oraclePrice: string,
  fundingIndex: string,
  effectiveAt: string,
  effectiveAtHeight: string,
}

export interface FundingIndexUpdatesUpdateObject {
  id: string,
  perpetualId?: string,
  eventId?: Buffer,
  rate?: string,
  oraclePrice?: string,
  fundingIndex?: string,
  effectiveAt?: string,
  effectiveAtHeight?: string,
}

export enum FundingIndexUpdatesColumns {
  id = 'id',
  perpetualId = 'perpetualId',
  eventId = 'eventId',
  rate = 'rate',
  oraclePrice = 'oraclePrice',
  fundingIndex = 'fundingIndex',
  effectiveAt = 'effectiveAt',
  effectiveAtHeight = 'effectiveAtHeight',
}
