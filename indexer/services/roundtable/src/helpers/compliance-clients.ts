import {
  ComplianceClient,
  getComplianceClient,
} from '@dydxprotocol-indexer/compliance';
import { ComplianceProvider } from '@dydxprotocol-indexer/postgres';

export interface ClientAndProvider {
  client: ComplianceClient,
  provider: ComplianceProvider,
}

export const complianceProvider: ClientAndProvider = {
  client: getComplianceClient(),
  provider: ComplianceProvider.ELLIPTIC,
};
