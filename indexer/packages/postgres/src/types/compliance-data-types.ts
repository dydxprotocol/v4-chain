/* ------- COMPLIANCE DATA TYPES ------- */

type IsoString = string;

export enum ComplianceProvider {
  ELLIPTIC = 'ELLIPTIC',
}

export interface ComplianceDataCreateObject {
  address: string;
  provider: string;
  chain?: string;
  sanctioned: boolean;
  riskScore?: string;
  updatedAt: IsoString;
}

export interface ComplianceDataUpdateObject {
  address: string;
  provider: string;
  chain?: string;
  sanctioned?: boolean;
  riskScore?: string;
  updatedAt?: IsoString;
}

export enum ComplianceDataColumns {
  address = 'address',
  provider = 'provider',
  chain = 'chain',
  sanctioned = 'sanctioned',
  riskScore = 'riskScore',
  updatedAt = 'updatedAt',
}
