export enum PnlColumns {
  subaccountId = 'subaccountId',
  createdAt = 'createdAt',
  createdAtHeight = 'createdAtHeight',
  deltaFundingPayments = 'deltaFundingPayments',
  deltaPositionEffects = 'deltaPositionEffects',
  totalPnl = 'totalPnl',
}

export interface PnlCreateObject {
  subaccountId: string,
  createdAt: string,
  createdAtHeight: string,
  deltaFundingPayments: string,
  deltaPositionEffects: string,
  totalPnl: string,
}
