import { CountryHeaders } from '@dydxprotocol-indexer/compliance';

import { IncomingMessage } from '../types';

export function getCountry(req: IncomingMessage): string | undefined {
  const countryHeaders: CountryHeaders = req.headers as CountryHeaders;
  return countryHeaders['cf-ipcountry'];
}
