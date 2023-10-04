import {
  CountryHeaders,
  isRestrictedCountryHeaders,
} from '@dydxprotocol-indexer/compliance';
import express from 'express';

import { INDEXER_GEOBLOCKED_PAYLOAD } from '../constants';
import { BlockedCode } from '../types';
import { create4xxResponse } from './helpers';

/**
 * Return an error code for users that access the API from a restricted country
 */
export function rejectRestrictedCountries(
  req: express.Request,
  res: express.Response,
  next: express.NextFunction,
) {
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
