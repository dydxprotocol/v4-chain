import { log } from 'console';

import {
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';
import {
  Controller, Get, Route, Path,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import { NotFoundError } from '../../../lib/errors';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { CheckTickerParamSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import { HistoricalFundingRequest, HistoricalFundingPaymentResponse, MarketType } from '../../../types';

const router: express.Router = express.Router();

@Route('historicalFundingPayment')
class HistoricalFundingPaymentController extends Controller {
  @Get('/:ticker')
  async getHistoricalFundingPayment(
    @Path() ticker: string,
  //   @Query() limit?: number,
  //   @Query() effectiveBeforeOrAtHeight?: number,
  //   @Query() effectiveBeforeOrAt?: IsoString,
  ): Promise<HistoricalFundingPaymentResponse> {
    const perpetualMarket: (
      PerpetualMarketFromDatabase | undefined
    ) = await PerpetualMarketTable.findByTicker(ticker);

    if (perpetualMarket === undefined) {
      throw new NotFoundError(`${ticker} not found in markets of type ${MarketType.PERPETUAL}`);
    }

    return { historicalFundingPayments: [{ ticker }] };
  }
}

router.get(
  '/:ticker',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckTickerParamSchema,
  handleValidationErrors,
  async (req: express.Request, res: express.Response) => {
    const {
      ticker,
    }: {
      ticker: string,
    } = matchedData(req) as HistoricalFundingRequest;

    try {
      const controller:
      HistoricalFundingPaymentController = new HistoricalFundingPaymentController();
      const response:
      HistoricalFundingPaymentResponse = await controller.getHistoricalFundingPayment(ticker);

      log(`Response is ${response}`);
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'HistoricalFundingPaymentController GET /',
        'HistoricalFundingPayment error',
        error,
        req,
        res,
      );
    }
  },
);

export default router;
