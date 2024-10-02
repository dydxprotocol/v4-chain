export interface AffiliateInfoCreateObject {
  address: string,
  affiliateEarnings: string,
  referredMakerTrades: number,
  referredTakerTrades: number,
  totalReferredMakerFees: string,
  totalReferredTakerFees: string,
  totalReferredUsers: number,
  firstReferralBlockHeight: string,
  referredTotalVolume: string,
}

export enum AffiliateInfoColumns {
  address = 'address',
  affiliateEarnings = 'affiliateEarnings',
  referredMakerTrades = 'referredMakerTrades',
  referredTakerTrades = 'referredTakerTrades',
  totalReferredMakerFees = 'totalReferredMakerFees',
  totalReferredTakerFees = 'totalReferredTakerFees',
  totalReferredUsers = 'totalReferredUsers',
  firstReferralBlockHeight = 'firstReferralBlockHeight',
  referredTotalVolume = 'referredTotalVolume',
}
