import {
  CountryHeaders,
  isRestrictedCountryHeaders,
} from '@dydxprotocol-indexer/compliance';

import { IncomingMessage } from '../types';

export class CountryRestrictor {
  public isRestrictedCountry(req: IncomingMessage): boolean {
    if (isRestrictedCountryHeaders(req.headers as CountryHeaders)) {
      return true;
    }

    return false;
  }

  public getCountry(req: IncomingMessage): string | undefined {
    const countryHeaders: CountryHeaders = req.headers as CountryHeaders;
    return countryHeaders['cf-ipcountry'];
  }
}
