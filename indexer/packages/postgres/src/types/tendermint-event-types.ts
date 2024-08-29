/* ------- TENDERMINT EVENT TYPES ------- */

export interface TendermintEventCreateObject {
  blockHeight: string,
  transactionIndex: number,
  eventIndex: number,
}

export enum TendermintEventColumns {
  id = 'id',
  blockHeight = 'blockHeight',
  transactionIndex = 'transactionIndex',
  eventIndex = 'eventIndex',
}
