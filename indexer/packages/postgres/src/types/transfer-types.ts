/* ------- TRANSFER TYPES ------- */

export interface TransferCreateObject {
  senderSubaccountId?: string,
  recipientSubaccountId?: string,
  senderWalletAddress?: string,
  recipientWalletAddress?: string,
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
  senderWalletAddress = 'senderWalletAddress',
  recipientWalletAddress = 'recipientWalletAddress',
  assetId = 'assetId',
  size = 'size',
  eventId = 'eventId',
  transactionHash = 'transactionHash',
  createdAt = 'createdAt',
  createdAtHeight = 'createdAtHeight',
}

export enum TransferType {
  TRANSFER_IN = 'TRANSFER_IN',
  TRANSFER_OUT = 'TRANSFER_OUT',
  DEPOSIT = 'DEPOSIT',
  WITHDRAWAL = 'WITHDRAWAL',
}
