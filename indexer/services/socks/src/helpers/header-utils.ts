import { GeoOriginHeaders } from '@dydxprotocol-indexer/compliance';

import { IncomingMessage } from '../types';

export function getGeoOriginHeaders(req: IncomingMessage): GeoOriginHeaders {
  const geoOriginHeaders = req.headers as GeoOriginHeaders;
  return {
    'geo-origin-country': geoOriginHeaders['geo-origin-country'],
    'geo-origin-region': geoOriginHeaders['geo-origin-region'],
    'geo-origin-status': geoOriginHeaders['geo-origin-status'],
  } as GeoOriginHeaders;
}
