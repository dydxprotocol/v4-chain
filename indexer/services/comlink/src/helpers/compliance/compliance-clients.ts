import { ComplianceClient, PlaceHolderProviderClient } from '@dydxprotocol-indexer/compliance';
import { ComplianceProvider } from '@dydxprotocol-indexer/postgres';

export interface ClientAndProvider {
  client: ComplianceClient;
  provider: ComplianceProvider;
}

// TODO(IND-369): Change this to the Elliptic client
export const placeHolderProvider: ClientAndProvider = {
  client: new PlaceHolderProviderClient(),
  provider: ComplianceProvider.ELLIPTIC,
};
