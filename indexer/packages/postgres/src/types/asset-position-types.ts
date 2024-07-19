/* ------- ASSET POSITION TYPES ------- */

export interface AssetPositionCreateObject {
  subaccountId: string,
  assetId: string,
  size: string,
  isLong: boolean,
}

export enum AssetPositionColumns {
  id = 'id',
  subaccountId = 'subaccountId',
  assetId = 'assetId',
  size = 'size',
  isLong = 'isLong',
}
