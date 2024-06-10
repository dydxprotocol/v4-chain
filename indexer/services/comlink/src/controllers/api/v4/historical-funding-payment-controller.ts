import {
  BlockFromDatabase,
  BlockTable,
  DEFAULT_POSTGRES_OPTIONS,
  FundingIndexMap,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  PerpetualPositionFromDatabase,
  PerpetualPositionStatus,
  PerpetualPositionTable,
  SubaccountFromDatabase,
  SubaccountTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';
import {
  Controller, Get, Route, Path, Query,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import { NotFoundError } from '../../../lib/errors';
import {
  getFundingIndexMaps,
  getPerpetualPositionsWithUpdatedFunding,
  handleControllerError,
  initializePerpetualPositionsWithFunding,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import {
  CheckSubaccountSchema,
  CheckTickerParamSchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import {
  HistoricalFundingRequest,
  HistoricalFundingPaymentResponse,
  MarketType,
  PerpetualPositionWithFunding,
} from '../../../types';

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
    const subaccountId: string = SubaccountTable.uuid(
      address,
      subaccountNumber,
    );
    const [
      subaccount,
      perpetualMarket,
      latestBlock,
    ] : [
      SubaccountFromDatabase | undefined,
      PerpetualMarketFromDatabase | undefined,
      BlockFromDatabase,
    ] = await Promise.all([
      SubaccountTable.findById(subaccountId),
      PerpetualMarketTable.findByTicker(ticker),
      BlockTable.getLatest(),
    ]);

    if (perpetualMarket === undefined) {
      throw new NotFoundError(
        `${ticker} not found in markets of type ${MarketType.PERPETUAL}`,
      );
    }

    if (subaccount === undefined) {
      throw new NotFoundError(
        `No subaccount found with address ${address} and subaccountNumber ${subaccountNumber}`,
      );
    }

    // Add tests for pagination, limits & date ranges

    const settledPositions = await PerpetualPositionTable.findAll(
      {
        subaccountId: [subaccountId],
        perpetualId: [perpetualMarket?.id],
      },
      [],
      {
        ...DEFAULT_POSTGRES_OPTIONS,
      },
    );

    const closedPositions = settledPositions.filter((position) => {
      return position.status === PerpetualPositionStatus.CLOSED ||
      position.status === PerpetualPositionStatus.LIQUIDATED;
    });
    const settledFundingPayments = mapSettledFundingPayments(closedPositions);

    const {
      lastUpdatedFundingIndexMap,
      latestFundingIndexMap,
    }: {
      lastUpdatedFundingIndexMap: FundingIndexMap;
      latestFundingIndexMap: FundingIndexMap;
    } = await getFundingIndexMaps(subaccount, latestBlock);

    const openPositions = settledPositions.filter((position) => {
      return position.status === PerpetualPositionStatus.OPEN;
    });
    const positionsWithUnsettledFunding = getPerpetualPositionsWithUpdatedFunding(
      initializePerpetualPositionsWithFunding(openPositions),
      latestFundingIndexMap,
      lastUpdatedFundingIndexMap,
    );

    const unsettledFundingPayments = mapUnsettledFundingPayments(
      positionsWithUnsettledFunding,
      subaccount.updatedAt,
    );

    const combined = unsettledFundingPayments.concat(settledFundingPayments);

    return {
      ticker,
      fundingPayments: combined,
    };
  }
}

function mapUnsettledFundingPayments(
  positions: PerpetualPositionWithFunding[],
  subaccountUpdatedAt: string) {
  return positions.map((position) => ({
    payment: position.unsettledFunding,
    effectiveAt: subaccountUpdatedAt,
  }));
}

function mapSettledFundingPayments(positions: PerpetualPositionFromDatabase[]) {
  return positions.map((position) => ({
    payment: position.settledFunding,
    effectiveAt: position.closedAt!,
  }));
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
      ticker: string;
      address: string;
      subaccountNumber: number;
    } = matchedData(req) as HistoricalFundingRequest;

    try {
      const controller:
      HistoricalFundingPaymentController = new HistoricalFundingPaymentController();

      const response:
      HistoricalFundingPaymentResponse = await controller.getHistoricalFundingPayment(
        ticker,
        address,
        subaccountNumber,
      );

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
