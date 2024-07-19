import { ComplianceClientResponse } from '../types';
import { ComplianceClient } from './compliance-client';

export class PlaceHolderProviderClient extends ComplianceClient {
  public getComplianceResponse(address: string): Promise<ComplianceClientResponse> {
    let blocked: boolean = false;
    let riskScore: string | undefined;

    if (
      address.charCodeAt(address.length - 1) > 'a'.charCodeAt(0) &&
      address.charCodeAt(address.length - 1) < 'm'.charCodeAt(0)
    ) {
      blocked = true;
      riskScore = '75';
    }

    return Promise.resolve({
      address,
      blocked,
      riskScore,
    });
  }
}
