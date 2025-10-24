import { stats } from '@dydxprotocol-indexer/base';

import config from '../config';
import { GeoOriginHeaders, GeoOriginStatus } from '../types';

export function isRestrictedCountryHeaders(headers: GeoOriginHeaders): boolean {
  if (config.INDEXER_LEVEL_GEOBLOCKING_ENABLED === false) {
    return false;
  }

  const geoStatus: string | undefined = headers['geo-origin-status'];

  if (
    geoStatus === undefined ||
    geoStatus !== GeoOriginStatus.OK
  ) {
    stats.increment(
      `${config.SERVICE_NAME}.rejected_restricted_country`,
      1,
      undefined,
      {
        country: headers['geo-origin-country'] || '',
        region: headers['geo-origin-region'] || '',
        status: headers['geo-origin-status'] || '',
      },
    );
    return true;
  }

  return false;
}

export function isWhitelistedAddress(address: string): boolean {
  return config.WHITELISTED_ADDRESSES.split(',').includes(address);
}
