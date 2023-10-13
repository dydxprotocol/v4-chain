/** Response from a compliance provider */
export interface ComplianceClientResponse {
  address: string,
  chain?: string,
  blocked: boolean,
  riskScore?: string,
}

export enum BlockedCode {
  GEOBLOCKED = 'GEOBLOCKED',
  COMPLIANCE_BLOCKED = 'COMPLIANCE_BLOCKED',
}
