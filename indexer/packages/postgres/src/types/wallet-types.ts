/* ------- WALLET TYPES ------- */

export interface WalletCreateObject {
  address: string,
  totalTradingRewards: string,
  totalVolume: string,
}

export interface WalletUpdateObject {
  address: string,
  totalTradingRewards: string,
  totalVolume: string,
}

export enum WalletColumns {
  address = 'address',
  totalTradingRewards = 'totalTradingRewards',
  totalVolume = 'totalVolume',
}
