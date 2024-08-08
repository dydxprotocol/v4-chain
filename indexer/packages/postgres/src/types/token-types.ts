/* ------- TOKEN TYPES ------- */

type IsoString = string;

export interface TokenCreateObject {
  token: string,
  address: string,
  updatedAt: IsoString,
}

export interface TokenUpdateObject {
  token: string,
  address: string,
  updatedAt: IsoString,
}

export enum TokenColumns {
  token = 'token',
  address = 'address',
  updatedAt = 'updatedAt',
}
