export interface PersistentCacheCreateObject {
  key: string,
  value: string,
}

export enum PersistentCacheColumns {
  key = 'key',
  value = 'value',
}

export enum PersistentCacheKeys {
  TOTAL_VOLUME_UPDATE_TIME = 'totalVolumeUpdateTime',
  AFFILIATE_INFO_UPDATE_TIME = 'affiliateInfoUpdateTime',
  FUNDING_PAYMENTS_LAST_PROCESSED_HEIGHT = 'fundingPaymentsLastProcessedHeight',
  PNL_LAST_PROCESSED_HEIGHT = 'pnlLastProcessedHeight',
}
