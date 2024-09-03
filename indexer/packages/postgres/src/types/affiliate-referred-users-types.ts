export interface AffiliateReferredUsersCreateObject {
  affiliateAddress: string,
  refereeAddress: string,
  referredAtBlock: string,
}

export enum AffiliateReferredUsersColumns {
  affiliateAddress = 'affiliateAddress',
  refereeAddress = 'refereeAddress',
  referredAtBlock = 'referredAtBlock',
}
