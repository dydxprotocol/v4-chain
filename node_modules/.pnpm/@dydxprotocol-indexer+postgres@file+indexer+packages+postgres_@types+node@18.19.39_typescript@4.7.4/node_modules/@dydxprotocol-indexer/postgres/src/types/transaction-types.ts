/* ------- TRANSACTION TYPES ------- */

export interface TransactionCreateObject {
  blockHeight: string,
  transactionIndex: number,
  transactionHash: string,
}

export enum TransactionColumns {
  id = 'id',
  blockHeight = 'blockHeight',
  transactionIndex = 'transactionIndex',
  transactionHash = 'transactionHash',
}
