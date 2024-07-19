import { stats } from '@dydxprotocol-indexer/base';
import {
  IsoString,
  FillTable,
  FillFromDatabase,
  Liquidity,
  QueryableField,
  perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import {
  checkSchema,
  matchedData,
} from 'express-validator';
import {
  Controller, Get, Path, Query, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { NotFoundError } from '../../../lib/errors';
import {
  handleControllerError,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { rejectRestrictedCountries } from '../../../lib/restrict-countries';
import { CheckLimitAndCreatedBeforeOrAtSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { fillToTradeResponseObject } from '../../../request-helpers/request-transformer';
import {
  MarketType,
  TradeRequest,
  TradeResponse,
  TradeResponseObject,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'trades-controller';

@Route('trades')
class TradesController extends Controller {
  @Get('/perpetualMarket/:ticker')
  async getTrades(
    @Path() ticker: string,
      @Query() limit: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
  ): Promise<TradeResponse> {
    const clobPairId: string | undefined = perpetualMarketRefresher
      .getClobPairIdFromTicker(ticker);

    if (clobPairId === undefined) {
      throw new NotFoundError(`${ticker} not found in tickers of type ${MarketType.PERPETUAL}`);
    }

    const fills: FillFromDatabase[] = await FillTable.findAll(
      {
        clobPairId,
        liquidity: Liquidity.TAKER,
        limit,
        createdBeforeOrAtHeight: createdBeforeOrAtHeight
          ? createdBeforeOrAtHeight.toString()
          : undefined,
        createdBeforeOrAt,
      },
      [QueryableField.LIQUIDITY, QueryableField.CLOB_PAIR_ID, QueryableField.LIMIT],
    );

    return {
      trades: fills.map((fill: FillFromDatabase): TradeResponseObject => {
        return fillToTradeResponseObject(fill);
      }),
    };
  }
}

router.get(
  '/perpetualMarket/:ticker',
  rejectRestrictedCountries,
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckLimitAndCreatedBeforeOrAtSchema,
  ...checkSchema({
    ticker: {
      in: ['params'],
      isString: true,
      custom: {
        options: perpetualMarketRefresher.isValidPerpetualMarketTicker,
        errorMessage: 'ticker must be a valid ticker (BTC-USD, etc)',
      },
    },
  }),
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      ticker,
      limit,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
    }: TradeRequest = matchedData(req) as TradeRequest;

    try {
      const controller: TradesController = new TradesController();
      const response: TradeResponse = await controller.getTrades(
        ticker,
        limit,
        createdBeforeOrAtHeight,
        createdBeforeOrAt,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'TradesController GET /perpetualMarket/:ticker',
        'Trades error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_trades.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
