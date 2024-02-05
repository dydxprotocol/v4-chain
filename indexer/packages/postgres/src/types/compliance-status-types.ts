/* ------- COMPLIANCE STATUS TYPES ------- */

type IsoString = string;

export enum ComplianceReason {
  MANUAL = 'MANUAL',
  US_GEO = 'US_GEO',
  CA_GEO = 'CA_GEO',
  SANCTIONED_GEO = 'SANCTIONED_GEO',
  COMPLIANCE_PROVIDER = 'COMPLIANCE_PROVIDER',
}

export enum ComplianceStatus {
  COMPLIANT = 'COMPLIANT',
  FIRST_STRIKE = 'FIRST_STRIKE',
  CLOSE_ONLY = 'CLOSE_ONLY',
  BLOCKED = 'BLOCKED',
}

export interface ComplianceStatusCreateObject {
  address: string;
  status: ComplianceStatus;
  reason?: ComplianceReason;
  createdAt?: IsoString;
  updatedAt?: IsoString;
}

export interface ComplianceStatusUpsertObject {
  address: string;
  status: ComplianceStatus;
  reason?: ComplianceReason;
}

export interface ComplianceStatusUpdateObject {
  address: string;
  status?: ComplianceStatus;
  reason?: ComplianceReason | null;
  createdAt?: IsoString;
  updatedAt?: IsoString;
}

export enum ComplianceStatusColumns {
  address = 'address',
  status = 'status',
  reason = 'reason',
  createdAt = 'createdAt',
  updatedAt = 'updatedAt',
}
