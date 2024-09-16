/* ------- TOKEN TYPES ------- */

type IsoString = string;

export interface FirebaseNotificationTokenCreateObject {
  token: string,
  address: string,
  language: string,
  updatedAt: IsoString,
}

export interface FirebaseNotificationTokenUpdateObject {
  token?: string,
  address?: string,
  language?: string,
  updatedAt?: IsoString,
}

export enum FirebaseNotificationTokenColumns {
  token = 'token',
  address = 'address',
  language = 'language',
  updatedAt = 'updatedAt',
}
