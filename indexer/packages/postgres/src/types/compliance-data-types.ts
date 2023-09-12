/* ------- COMPLIANCE DATA TYPES ------- */

type IsoString = string;

export interface ComplianceDataCreateObject {
  address: string;
  chain?: string;
  sanctioned: boolean;
  riskScore?: string;
  updatedAt: IsoString;
}

export interface ComplianceDataUpdateObject {
  address: string;
  chain?: string;
  sanctioned?: boolean;
  riskScore?: string;
  updatedAt?: IsoString;
}

export enum ComplianceDataColumns {
  address = 'address',
  chain = 'chain',
  sanctioned = 'sanctioned',
  riskScore = 'riskScore',
  updatedAt = 'updatedAt',
}
