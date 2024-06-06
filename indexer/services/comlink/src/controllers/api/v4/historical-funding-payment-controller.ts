import { log } from 'console';

import {
  DEFAULT_POSTGRES_OPTIONS,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  PerpetualPositionStatus,
  PerpetualPositionTable,
  SubaccountTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';
import {
  Controller, Get, Route, Path,
  Query,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import { NotFoundError } from '../../../lib/errors';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { CheckSubaccountSchema, CheckTickerParamSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import { HistoricalFundingRequest, HistoricalFundingPaymentResponse, MarketType } from '../../../types';

const router: express.Router = express.Router();

@Route('historicalFundingPayment')
class HistoricalFundingPaymentController extends Controller {
  @Get('/:ticker')
  async getHistoricalFundingPayment(
    @Path() ticker: string,
      @Query() address: string,
      @Query() subaccountNumber: number,
  //   @Query() limit?: number,
  //   @Query() effectiveBeforeOrAtHeight?: number,
  //   @Query() effectiveBeforeOrAt?: IsoString,
  ): Promise<HistoricalFundingPaymentResponse> {
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);

    const perpetualMarket: (
      PerpetualMarketFromDatabase | undefined
    ) = await PerpetualMarketTable.findByTicker(ticker);

    if (perpetualMarket === undefined) {
      throw new NotFoundError(`${ticker} not found in markets of type ${MarketType.PERPETUAL}`);
    }

    const settledPositions = await PerpetualPositionTable.findAll({
      subaccountId: [subaccountId],
      perpetualId: [perpetualMarket?.id],
      status: [PerpetualPositionStatus.CLOSED],
    },
    [],
    {
      ...DEFAULT_POSTGRES_OPTIONS,
    },
    );

    return {
      ticker,
      fundingPayments: settledPositions.map((position) => {
        return {
          payment: position.settledFunding,
          effectiveAt: position.closedAt!,
        };
      }),
    };
  }
}

router.get(
  '/:ticker',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckTickerParamSchema,
  ...CheckSubaccountSchema,
  handleValidationErrors,
  async (req: express.Request, res: express.Response) => {
    const {
      ticker,
      address,
      subaccountNumber,
    }: {
      ticker: string,
      address: string,
      subaccountNumber: number,
    } = matchedData(req) as HistoricalFundingRequest;

    try {
      const controller:
      HistoricalFundingPaymentController = new HistoricalFundingPaymentController();

      const response: HistoricalFundingPaymentResponse = await
      controller.getHistoricalFundingPayment(ticker, address, subaccountNumber);

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
