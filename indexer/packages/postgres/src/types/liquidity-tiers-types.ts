/* ------- LIQUIDITY TIERS TYPES ------- */

export interface LiquidityTiersCreateObject {
  id: number,
  name: string,
  initialMarginPpm: string,
  maintenanceFractionPpm: string,
  basePositionNotional: string,
}

export interface LiquidityTiersUpdateObject {
  id: number,
  name?: string,
  initialMarginPpm?: string,
  maintenanceFractionPpm?: string,
  basePositionNotional?: string,
}

export enum LiquidityTiersColumns {
  id = 'id',
  name = 'name',
  initialMarginPpm = 'initialMarginPpm',
  maintenanceFractionPpm = 'maintenanceFractionPpm',
  basePositionNotional = 'basePositionNotional',
}
