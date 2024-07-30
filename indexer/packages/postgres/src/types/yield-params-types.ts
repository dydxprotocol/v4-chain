/* ------- YIELD PARAMS TYPES ------- */

import { IsoString } from './utility-types';

export interface YieldParamsCreateObject {
  sDAIPrice: string,
  assetYieldIndex: string,
  createdAt: IsoString,
  createdAtHeight: string,
}

export enum YieldParamsColumns {
  id = 'id',
  sDAIPrice = 'sDAIPrice',
  assetYieldIndex = 'assetYieldIndex',
  createdAt = 'createdAt',
  createdAtHeight = 'createdAtHeight',
}
