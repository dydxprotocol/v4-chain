import { stats } from '@dydxprotocol-indexer/base';
import express from 'express';

import config from '../config';
import { create4xxResponse } from './helpers';
import { isRestrictedCountry } from './utils';

/**
 * Return an error code for users that access the API from a restricted country
 */
export function rejectRestrictedCountries(
  req: express.Request,
  res: express.Response,
  next: express.NextFunction,
) {
  const {
    'cf-ipcountry': ipCountry,
  } = req.headers as {
    'cf-ipcountry'?: string,
  };

  if (
    ipCountry !== undefined &&
    isRestrictedCountry(ipCountry)
  ) {
    stats.increment(
      `${config.SERVICE_NAME}.rejected_restricted_country_request`,
      1,
      undefined,
      {
        country: ipCountry,
      },
    );
    return create4xxResponse(
      res,
      'Forbidden',
      403,
    );
  }

  return next();
}
