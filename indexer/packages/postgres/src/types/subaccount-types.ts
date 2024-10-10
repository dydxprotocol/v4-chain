/* ------- SUBACCOUNT TYPES ------- */

import { IsoString } from './utility-types';

export interface SubaccountCreateObject {
  address: string,
  assetYieldIndex: string,
  subaccountNumber: number,
  updatedAt: IsoString,
  updatedAtHeight: string,
}

export interface SubaccountUpdateObject {
  assetYieldIndex: string,
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
  assetYieldIndex = 'assetYieldIndex',
}
