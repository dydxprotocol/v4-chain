/* ------- TRANSFER TYPES ------- */

export interface TransferCreateObject {
  senderSubaccountId: string,
  recipientSubaccountId: string,
  assetId: string,
  size: string,
  eventId: Buffer,
  transactionHash: string,
  createdAt: string,
  createdAtHeight: string,
}

export enum TransferColumns {
  id = 'id',
  senderSubaccountId = 'senderSubaccountId',
  recipientSubaccountId = 'recipientSubaccountId',
  assetId = 'assetId',
  size = 'size',
  eventId = 'eventId',
  transactionHash = 'transactionHash',
  createdAt = 'createdAt',
  createdAtHeight = 'createdAtHeight',
}
