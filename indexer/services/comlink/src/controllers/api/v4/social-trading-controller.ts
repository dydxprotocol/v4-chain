import {
  SubaccountTable,
  SubaccountUsernamesTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import {
  checkSchema,
  matchedData,
} from 'express-validator';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import { NotFoundError } from '../../../lib/errors';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { SubaccountInfoToTraderSearchResponse } from '../../../request-helpers/request-transformer';

const router: express.Router = express.Router();
const controllerName: string = 'social-trading-controller';

function checkIfValidDydxAddress(address: string): boolean {
  const pattern = /^dydx[0-9a-z]{39}$/;
  return pattern.test(address);
}

async function searchTrader(
  searchParam: string,
) {
  if (checkIfValidDydxAddress(searchParam)) {

    const subaccounts = await SubaccountTable.findAll({
      address: searchParam,
      subaccountNumber: 0,
    }, [], { readReplica: true });
    const subaccount = subaccounts[0];

    if (!subaccount) {
      throw new NotFoundError(`Subaccount not found:${searchParam}`);
    }

    const subaccountUsernames = await SubaccountUsernamesTable.findAll({
      subaccountId: [subaccount.id],
    }, [], { readReplica: true });

    return SubaccountInfoToTraderSearchResponse(subaccount, subaccountUsernames[0]);
  } else {
    const subaccountUsername = await SubaccountUsernamesTable.findByUsername(searchParam);

    if (!subaccountUsername) {
      throw new NotFoundError(`Subaccount not found:${searchParam}`);
    }

    const subaccount = await SubaccountTable.findById(subaccountUsername.subaccountId);
    if (!subaccount) {
      throw new NotFoundError(`Subaccount not found:${subaccountUsername.subaccountId}`);
    }

    return SubaccountInfoToTraderSearchResponse(subaccount, subaccountUsername);
  }
}

router.get('/search', rateLimiterMiddleware(getReqRateLimiter), ...checkSchema({
  searchParam: {
    in: 'query',
    isString: true,
    errorMessage: 'searchParam is required',
  },
}), handleValidationErrors,
ExportResponseCodeStats({ controllerName }),
async (req, res) => {
  try {
    const { searchParam } = matchedData(req);
    const response = await searchTrader(searchParam);
    return res.send(response);
  } catch (error) {
    return handleControllerError(
      'SocialTradingController GET /',
      'User search error',
      error,
      req,
      res,
    );
  }
});

export default router;
