/* ------- TOKEN TYPES ------- */

type IsoString = string;

export interface TokenCreateObject {
  token: string,
  address: string,
  language: string,
  updatedAt: IsoString,
}

export interface TokenUpdateObject {
  token: string,
  address: string,
  language: string,
  updatedAt: IsoString,
}

export enum TokenColumns {
  token = 'token',
  address = 'address',
  language = 'language',
  updatedAt = 'updatedAt',
}
