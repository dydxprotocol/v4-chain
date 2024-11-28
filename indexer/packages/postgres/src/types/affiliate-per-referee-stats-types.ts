export interface AffiliatePerRefereeStatsCreateObject {
  affiliateAddress: string,
  refereeAddress: string,
  affiliateEarnings: string,
  referredMakerTrades: number,
  referredTakerTrades: number,
  referredTotalVolume: string,
  firstReferralBlockHeight: number,
  totalReferredTakerFees: string,
  totalReferredMakerFees: string,
  totalReferredMakerRebates: string,
  totalReferredLiquidationfees: string,
}

export enum AffiliatePerRefereeStatsColumns {
  affiliateAddress = 'affiliateAddress',
  refereeAddress = 'refereeAddress',
  affiliateEarnings = 'affiliateEarnings',
  referredMakerTrades = 'referredMakerTrades',
  referredTakerTrades = 'referredTakerTrades',
  referredTotalVolume = 'referredTotalVolume',
  firstReferralBlockHeight = 'firstReferralBlockHeight',
  totalReferredTakerFees = 'totalReferredTakerFees',
  totalReferredMakerFees = 'totalReferredMakerFees',
  totalReferredMakerRebates = 'totalReferredMakerRebates',
  totalReferredLiquidationfees = 'totalReferredLiquidationfees',
}
