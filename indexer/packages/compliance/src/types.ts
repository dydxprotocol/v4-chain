export interface ComplianceClientResponse {
  address: string,
  chain?: string,
  blocked: boolean,
  riskScore?: string,
}

export interface GeoOriginHeaders extends Record<string, string | undefined> {
  'geo-origin-country'?: string,
  'geo-origin-region'?: string,
  'geo-origin-status'?: string,
}

export enum GeoOriginStatus {
  OK = 'ok',
  RESTRICTED = 'restricted',
}
