export interface AffiliateReferredUsersCreateObject {
  affiliateAddress: string;
  refereeAddress: string;
  referredAtBlock: number;
}

export enum AffiliateReferredUsersColumns {
  affiliateAddress = 'affiliateAddress',
  refereeAddress = 'refereeAddress',
  referredAtBlock = 'referredAtBlock',
}
