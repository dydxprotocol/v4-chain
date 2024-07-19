import { ComplianceClientResponse } from '../types';

export abstract class ComplianceClient {
  public abstract getComplianceResponse(address: string): Promise<ComplianceClientResponse>;
}
