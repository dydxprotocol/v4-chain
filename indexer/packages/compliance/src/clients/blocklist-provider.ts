import config from '../config';
import { ComplianceClientResponse } from '../types';
import { ComplianceClient } from './compliance-client';

const BLOCKED_ADDRESSES: Set<string> = new Set(config.BLOCKED_ADDRESSES.split(','));

export class BlocklistProviderClient extends ComplianceClient {
  public getComplianceResponse(address: string): Promise<ComplianceClientResponse> {
    let blocked: boolean = false;
    if (BLOCKED_ADDRESSES.has(address)) {
      blocked = true;
    }

    return Promise.resolve({
      address,
      blocked,
    });
  }
}
