export enum PnlColumns {
  subaccountId = 'subaccountId',
  createdAt = 'createdAt',
  createdAtHeight = 'createdAtHeight',
  equity = 'equity',
  netTransfers = 'netTransfers',
  totalPnl = 'totalPnl',
}

export interface PnlCreateObject {
  subaccountId: string,
  createdAt: string,
  createdAtHeight: string,
  equity: string,
  netTransfers: string,
  totalPnl: string,
}
