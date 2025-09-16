export interface AffiliateRefereeStatsCreateObject {
  affiliateAddress: string,
  refereeAddress: string,
  affiliateEarnings: string,
  referredMakerTrades: number,
  referredTakerTrades: number,
  referredTotalVolume: string,
  referralBlockHeight: string,
  referredTakerFees: string,
  referredMakerFees: string,
  referredMakerRebates: string,
  referredLiquidationFees: string,
}

export enum AffiliateRefereeStatsColumns {
  affiliateAddress = 'affiliateAddress',
  refereeAddress = 'refereeAddress',
  affiliateEarnings = 'affiliateEarnings',
  referredMakerTrades = 'referredMakerTrades',
  referredTakerTrades = 'referredTakerTrades',
  referredTotalVolume = 'referredTotalVolume',
  referralBlockHeight = 'referralBlockHeight',
  referredTakerFees = 'referredTakerFees',
  referredMakerFees = 'referredMakerFees',
  referredMakerRebates = 'referredMakerRebates',
  referredLiquidationFees = 'referredLiquidationFees',
}
