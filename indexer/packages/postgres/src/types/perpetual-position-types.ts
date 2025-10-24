/* ------- PERPETUAL POSITION TYPES ------- */

import { PositionSide } from './position-types';

type IsoString = string;

export type MarketOpenInterest = {
  perpetualMarketId: string,
  openInterest: string,
};

export enum PerpetualPositionStatus {
  OPEN = 'OPEN',
  CLOSED = 'CLOSED',
  LIQUIDATED = 'LIQUIDATED',
}

export interface PerpetualPositionCreateObject {
  subaccountId: string,
  perpetualId: string,
  side: PositionSide,
  status: PerpetualPositionStatus,
  size: string,
  maxSize: string,
  sumOpen?: string,
  sumClose?: string,
  entryPrice?: string,
  createdAt: IsoString,
  createdAtHeight: string,
  openEventId: Buffer,
  lastEventId: Buffer,
  settledFunding: string,
  closedAt?: IsoString,
  closedAtHeight?: string,
  closeEventId?: Buffer,
  exitPrice?: string,
  totalRealizedPnl?: string,
}

export interface PerpetualPositionUpdateObject {
  id: string,
  side?: PositionSide,
  status?: PerpetualPositionStatus,
  size?: string,
  maxSize?: string,
  entryPrice?: string,
  exitPrice?: string | null,
  sumOpen?: string,
  sumClose?: string,
  createdAt?: IsoString,
  closedAt?: IsoString | null,
  createdAtHeight?: string,
  closedAtHeight?: string | null,
  closeEventId?: Buffer | null,
  lastEventId?: Buffer,
  settledFunding?: string,
  totalRealizedPnl?: string,
}

// Object used to update a subaccount's perpetual position in the SubaccountUpdateHandler
export interface PerpetualPositionSubaccountUpdateObject {
  id: string,
  closedAt?: IsoString | null,
  closedAtHeight?: string | null,
  closeEventId?: Buffer | null,
  lastEventId: Buffer,
  settledFunding: string,
  status: PerpetualPositionStatus,
  size: string,
}

/*
This is all of the fields in PerpetualPositionFromDatabase with the exception of:
- subaccountId
- createdAt
- createdAtHeight
- openEventId
closedAt, closedAtHeight, and closeEventId are nullable.
*/
export interface UpdatedPerpetualPositionSubaccountKafkaObject {
  id: string,
  perpetualId: string,
  side: PositionSide,
  status: PerpetualPositionStatus,
  size: string,
  maxSize: string,
  entryPrice: string,
  exitPrice?: string,
  sumOpen: string,
  sumClose: string,
  closedAt?: IsoString | null,
  closedAtHeight?: string | null,
  lastEventId: Buffer,
  closeEventId?: Buffer | null,
  settledFunding: string,
  realizedPnl?: string,
  unrealizedPnl?: string,
}

export interface PerpetualPositionCloseObject {
  id: string,
  closedAt: IsoString,
  closedAtHeight: string,
  closeEventId: Buffer,
  settledFunding: string,
}

export enum PerpetualPositionColumns {
  id = 'id',
  subaccountId = 'subaccountId',
  perpetualId = 'perpetualId',
  side = 'side',
  status = 'status',
  size = 'size',
  maxSize = 'maxSize',
  entryPrice = 'entryPrice',
  exitPrice = 'exitPrice',
  sumOpen = 'sumOpen',
  sumClose = 'sumClose',
  createdAt = 'createdAt',
  closedAt = 'closedAt',
  createdAtHeight = 'createdAtHeight',
  closedAtHeight = 'closedAtHeight',
  openEventId = 'openEventId',
  closeEventId = 'closeEventId',
  lastEventId = 'lastEventId',
  settledFunding = 'settledFunding',
}
