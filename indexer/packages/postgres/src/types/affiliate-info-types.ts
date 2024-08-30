export interface AffiliateInfoCreateObject {
  address: string,
  affiliateEarnings: number,
  referredMakerTrades: number,
  referredTakerTrades: number,
  totalReferredFees: number,
  totalReferredUsers: number,
  referredNetProtocolEarnings: number,
  firstReferralBlockHeight: number,
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
}
