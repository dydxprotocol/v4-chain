import { IsoString } from './utility-types';

export interface VaultCreateObject {
  address: string,
  clobPairId: string,
  status: VaultStatus,
  createdAt: IsoString,
  updatedAt: IsoString,
}

export enum VaultStatus {
  DEACTIVATED = 'DEACTIVATED',
  STAND_BY = 'STAND_BY',
  QUOTING = 'QUOTING',
  CLOSE_ONLY = 'CLOSE_ONLY',
}

export enum VaultColumns {
  address = 'address',
  clobPairId = 'clobPairId',
  status = 'status',
  createdAt = 'createdAt',
  updatedAt = 'updatedAt',
}
