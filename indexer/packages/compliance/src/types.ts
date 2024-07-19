/** Response from a compliance provider */
export interface ComplianceClientResponse {
  address: string,
  chain?: string,
  blocked: boolean,
  riskScore?: string,
}
