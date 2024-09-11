/* ------- COMPLIANCE DATA TYPES ------- */

type IsoString = string;

export enum ComplianceProvider {
  ELLIPTIC = 'ELLIPTIC',
}

export interface ComplianceDataCreateObject {
  address: string,
  provider: string,
  chain?: string,
  blocked: boolean,
  riskScore?: string,
  updatedAt?: IsoString,
}

export interface ComplianceDataUpdateObject {
  address: string,
  provider: string,
  chain?: string,
  blocked?: boolean,
  riskScore?: string,
  updatedAt?: IsoString,
}

export enum ComplianceDataColumns {
  address = 'address',
  provider = 'provider',
  chain = 'chain',
  blocked = 'blocked',
  riskScore = 'riskScore',
  updatedAt = 'updatedAt',
}
