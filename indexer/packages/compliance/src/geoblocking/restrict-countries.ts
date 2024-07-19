import { stats } from '@dydxprotocol-indexer/base';

import config from '../config';
import { isRestrictedCountry } from './util';

export interface CountryHeaders {
  'cf-ipcountry'?: string,
}

export function isRestrictedCountryHeaders(headers: CountryHeaders): boolean {
  if (config.INDEXER_LEVEL_GEOBLOCKING_ENABLED === false) {
    return false;
  }

  const ipCountry: string | undefined = headers['cf-ipcountry'];

  if (
    ipCountry === undefined ||
    isRestrictedCountry(ipCountry)
  ) {
    stats.increment(
      `${config.SERVICE_NAME}.rejected_restricted_country`,
      1,
      undefined,
      {
        country: String(ipCountry),
      },
    );
    return true;
  }

  return false;
}
