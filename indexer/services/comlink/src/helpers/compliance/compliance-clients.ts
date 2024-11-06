import {
  ComplianceClient,
  getComplianceClient,
} from '@klyraprotocol-indexer/compliance';
import { ComplianceProvider } from '@klyraprotocol-indexer/postgres';

export interface ClientAndProvider {
  client: ComplianceClient;
  provider: ComplianceProvider;
}

export const complianceProvider: ClientAndProvider = {
  client: getComplianceClient(),
  provider: ComplianceProvider.ELLIPTIC,
};
