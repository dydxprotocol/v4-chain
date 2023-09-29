import { ComplianceClient } from './compliance-client';
import { PlaceHolderProviderClient } from './placeholder-provider';
import { BlocklistProviderClient } from './blocklist-provider';
import { EllipticProviderClient } from './elliptic-provider';
import config from '../config';

enum ClientType {
  PLACEHOLDER = 'PLACEHOLDER',
  BLOCKLIST = 'BLOCKLIST',
  ELLIPTIC = 'ELLIPTIC',
}

// Providers for test-net / development
const placeHolderProvider: ComplianceClient = new PlaceHolderProviderClient();
const blocklistProvider: ComplianceClient = new BlocklistProviderClient();
// Elliptic provider
const ellipticProvider: ComplianceClient  = new EllipticProviderClient();

const COMPLIANCE_CLIENTS: Record<ClientType, ComplianceClient> = {
  [ClientType.PLACEHOLDER]: placeHolderProvider,
  [ClientType.BLOCKLIST]: blocklistProvider,
  [ClientType.ELLIPTIC]: ellipticProvider,
}

export function getComplianceClient(): ComplianceClient {
  let complianceClient: ComplianceClient = ellipticProvider;
  if (COMPLIANCE_CLIENTS[config.COMPLIANCE_DATA_CLIENT as ClientType] !== undefined) {
    complianceClient = COMPLIANCE_CLIENTS[config.COMPLIANCE_DATA_CLIENT as ClientType];
  }
  return complianceClient;
}
