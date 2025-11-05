import { cacheControlMiddleware, stats } from '@dydxprotocol-indexer/base';
import {
  SubaccountFromDatabase,
  SubaccountTable,
  SubaccountUsernamesFromDatabase,
  SubaccountUsernamesTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import {
  checkSchema,
  matchedData,
} from 'express-validator';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { defaultRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { NotFoundError } from '../../../lib/errors';
import { checkIfValidDydxAddress, handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { subaccountInfoToTraderSearchResponse } from '../../../request-helpers/request-transformer';
import { TraderSearchRequest, TraderSearchResponse } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'social-trading-controller';
const socialTradingCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_SOCIAL_TRADING,
);

@Route('trader')
class SocialTradingController extends Controller {
  @Get('/search')
  async searchTrader(
    @Query() searchParam: string,
  ): Promise<TraderSearchResponse> {
    if (checkIfValidDydxAddress(searchParam)) {
      const subaccounts: SubaccountFromDatabase[] = await
      SubaccountTable.findAll({
        address: searchParam,
        subaccountNumber: 0,
      }, []);
      const subaccount: SubaccountFromDatabase = subaccounts[0];

      if (!subaccount) {
        throw new NotFoundError(`Subaccount not found:${searchParam}`);
      }

      const subaccountUsernames: SubaccountUsernamesFromDatabase[] = await
      SubaccountUsernamesTable.findAll({
        subaccountId: [subaccount.id],
      }, []);

      return subaccountInfoToTraderSearchResponse(subaccount, subaccountUsernames[0]);
    }

    const subaccountUsername: SubaccountUsernamesFromDatabase | undefined = await
    SubaccountUsernamesTable.findByUsername(searchParam);

    if (!subaccountUsername) {
      throw new NotFoundError(`Subaccount not found:${searchParam}`);
    }
    // subaccount search below cannot be undefined because of foreign key constraint
    const subaccount: SubaccountFromDatabase | undefined = await
    SubaccountTable.findById(subaccountUsername.subaccountId);

    return subaccountInfoToTraderSearchResponse(subaccount!, subaccountUsername);

  }
}

router.get('/search',
  rateLimiterMiddleware(defaultRateLimiter),
  socialTradingCacheControlMiddleware,
  ...checkSchema({
    searchParam: {
      in: 'query',
      isString: true,
      errorMessage: 'searchParam is required',
    },
  }),
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    try {
      const { searchParam } = matchedData(req) as TraderSearchRequest;
      const controller: SocialTradingController = new SocialTradingController();
      const response: TraderSearchResponse = await controller.searchTrader(searchParam);
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'SocialTradingController GET /',
        'User search error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.search_trader.timing`,
        Date.now() - start,
      );
    }
  });

export default router;
