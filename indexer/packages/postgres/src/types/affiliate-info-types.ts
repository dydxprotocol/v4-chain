export interface AffiliateInfoCreateObject {
  address: string,
  affiliateEarnings: string,
  referredMakerTrades: number,
  referredTakerTrades: number,
  totalReferredFees: string,
  totalReferredUsers: number,
  referredNetProtocolEarnings: string,
  firstReferralBlockHeight: string,
  referredTotalVolume: string,
}

export enum AffiliateInfoColumns {
  address = 'address',
  affiliateEarnings = 'affiliateEarnings',
  referredMakerTrades = 'referredMakerTrades',
  referredTakerTrades = 'referredTakerTrades',
  totalReferredFees = 'totalReferredFees',
  totalReferredUsers = 'totalReferredUsers',
  referredNetProtocolEarnings = 'referredNetProtocolEarnings',
  firstReferralBlockHeight = 'firstReferralBlockHeight',
  referredTotalVolume = 'referredTotalVolume',
}
