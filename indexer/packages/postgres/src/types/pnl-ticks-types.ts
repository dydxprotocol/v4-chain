/* ------- PNL TICKS TYPES ------- */

type IsoString = string;

export interface PnlTicksCreateObject {
  subaccountId: string,
  equity: string,
  totalPnl: string,
  netTransfers: string,
  createdAt: string,
  blockHeight: string,
  blockTime: IsoString,
}

export enum PnlTicksColumns {
  id = 'id',
  subaccountId = 'subaccountId',
  equity = 'equity',
  totalPnl = 'totalPnl',
  netTransfers = 'netTransfers',
  createdAt = 'createdAt',
  blockHeight = 'blockHeight',
  blockTime = 'blockTime',
}

export enum PnlTickInterval {
  hour = 'hour',
  day = 'day',
}
