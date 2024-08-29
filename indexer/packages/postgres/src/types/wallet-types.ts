/* ------- WALLET TYPES ------- */

export interface WalletCreateObject {
  address: string,
  totalTradingRewards: string,
  totalVolume: string,
  isWhitelistAffiliate: boolean,
}

export interface WalletUpdateObject {
  address: string,
  totalTradingRewards: string,
  totalVolume: string,
  isWhitelistAffiliate: boolean,
}

export enum WalletColumns {
  address = 'address',
  totalTradingRewards = 'totalTradingRewards',
  totalVolume = 'totalVolume',
  isWhitelistAffiliate = 'isWhitelistAffiliate',
}
