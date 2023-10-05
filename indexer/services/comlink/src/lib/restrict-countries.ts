import {
  CountryHeaders,
  isRestrictedCountryHeaders,
  INDEXER_GEOBLOCKED_PAYLOAD,
} from '@dydxprotocol-indexer/compliance';
import express from 'express';

import { BlockedCode } from '../types';
import { create4xxResponse } from './helpers';
import { getIpAddr, isIndexerIp } from './utils';

/**
 * Return an error code for users that access the API from a restricted country
 */
export function rejectRestrictedCountries(
  req: express.Request,
  res: express.Response,
  next: express.NextFunction,
) {
  const ipAddr: string | undefined = getIpAddr(req);

  // Don't enforce geo-blocking for internal IPs as they don't go through a proxy
  if (ipAddr !== undefined && isIndexerIp(ipAddr)) {
    return next();
  }

  if (isRestrictedCountryHeaders(req.headers as CountryHeaders)) {
    return create4xxResponse(
      res,
      INDEXER_GEOBLOCKED_PAYLOAD,
      403,
      { code: BlockedCode.GEOBLOCKED },
    );
  }

  return next();
}
