/* ------- BLOCK TYPES ------- */

type IsoString = string;

export interface BlockCreateObject {
  blockHeight: string,
  time: IsoString,
}

export enum BlockColumns {
  blockHeight = 'blockHeight',
  time = 'time',
}
