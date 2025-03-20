/* ------- SUBACCOUNT TYPES ------- */

import { IsoString } from './utility-types';

export interface SubaccountCreateObject {
  address: string,
  subaccountNumber: number,
  updatedAt: IsoString,
  updatedAtHeight: string,
}

export interface SubaccountUpdateObject {
  id: string,
  updatedAt: IsoString,
  updatedAtHeight: string,
}

export enum SubaccountColumns {
  id = 'id',
  address = 'address',
  subaccountNumber = 'subaccountNumber',
  updatedAt = 'updatedAt',
  updatedAtHeight = 'updatedAtHeight',
}

export interface ParentSubaccount {
  address: string,
  subaccountNumber: number,
}
