import { cacheControlMiddleware, stats } from '@dydxprotocol-indexer/base';
import {
  FillColumns,
  FillFromDatabase,
  FillTable,
  IsoString,
  OrderTable,
  OrderType,
  Ordering,
  SubaccountTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import {
  checkSchema,
  matchedData,
  query,
} from 'express-validator';
import _ from 'lodash';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { fillsRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceAndGeoCheck } from '../../../lib/compliance-and-geo-check';
import { NotFoundError } from '../../../lib/errors';
import {
  buildClobPairIdToMarket, getClobPairId, handleControllerError, isDefined,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { computeTradeHistory, paginateTradeHistory } from '../../../lib/trade-history';
import {
  CheckLimitAndCreatedBeforeOrAtSchema,
  CheckPaginationSchema,
  CheckParentSubaccountSchema,
  CheckSubaccountSchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import {
  MarketType,
  ParentSubaccountTradeHistoryRequest,
  TradeHistoryRequest,
  TradeHistoryResponse,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'trade-history-controller';
const tradeHistoryCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_FILLS,
);

// Shared helper: batch-fetch order types for a set of fills
async function buildOrderTypeMap(
  fills: FillFromDatabase[],
): Promise<Record<string, OrderType>> {
  const orderIds: string[] = _.uniq(
    fills
      .map((f) => f.orderId)
      .filter((id): id is string => id !== undefined && id !== null),
  );

  if (orderIds.length === 0) {
    return {};
  }

  const { results: orders } = await OrderTable.findAll(
    { id: orderIds },
    [],
  );

  const orderTypeMap: Record<string, OrderType> = {};
  for (const order of orders) {
    orderTypeMap[order.id] = order.type;
  }
  return orderTypeMap;
}

// Shared: ordering for fill queries (chronological ASC)
const fillOrderBy: [string, Ordering][] = [
  [FillColumns.createdAt, Ordering.ASC],
  [FillColumns.eventId, Ordering.ASC],
];

// Shared: market/marketType validation schema
const marketTypeCheckSchema = checkSchema({
  market: {
    in: ['query'],
    isString: true,
    optional: true,
  },
  marketType: {
    in: ['query'],
    isIn: {
      options: [Object.values(MarketType)],
    },
    optional: true,
    errorMessage: 'marketType must be a valid market type (PERPETUAL/SPOT)',
  },
});

@Route('tradeHistory')
class TradeHistoryController extends Controller {
  @Get('/')
  async getTradeHistory(
    @Query() address: string,
      @Query() subaccountNumber: number,
      @Query() market?: string,
      @Query() marketType?: MarketType,
      @Query() limit?: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
      @Query() page?: number,
  ): Promise<TradeHistoryResponse> {
    let clobPairId: string | undefined;
    if (isDefined(market) && isDefined(marketType)) {
      clobPairId = await getClobPairId(market!, marketType!);
      if (clobPairId === undefined) {
        throw new NotFoundError(`${market} not found in markets of type ${marketType}`);
      }
    }

    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);

    const { results: fills } = await FillTable.findAll(
      {
        subaccountId: [subaccountId],
        clobPairId,
        createdBeforeOrAtHeight: createdBeforeOrAtHeight?.toString(),
        createdBeforeOrAt,
      },
      [],
      { orderBy: fillOrderBy },
    );

    const orderTypeMap = await buildOrderTypeMap(fills);
    const clobPairIdToMarket = buildClobPairIdToMarket();
    const allRows = computeTradeHistory(fills, orderTypeMap, clobPairIdToMarket);

    const effectiveLimit = limit ?? config.API_LIMIT_V4;
    const paginated = paginateTradeHistory(allRows, effectiveLimit, page);

    return {
      tradeHistory: paginated.tradeHistory,
      pageSize: paginated.pageSize,
      totalResults: paginated.totalResults,
      offset: paginated.offset,
    };
  }

  @Get('/parentSubaccountNumber')
  async getTradeHistoryForParentSubaccount(
    @Query() address: string,
      @Query() parentSubaccountNumber: number,
      @Query() market?: string,
      @Query() marketType?: MarketType,
      @Query() limit?: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
      @Query() page?: number,
  ): Promise<TradeHistoryResponse> {
    let clobPairId: string | undefined;
    if (isDefined(market) && isDefined(marketType)) {
      clobPairId = await getClobPairId(market!, marketType!);
      if (clobPairId === undefined) {
        throw new NotFoundError(`${market} not found in markets of type ${marketType}`);
      }
    }

    const { results: fills } = await FillTable.findAll(
      {
        parentSubaccount: {
          address,
          subaccountNumber: parentSubaccountNumber,
        },
        clobPairId,
        createdBeforeOrAtHeight: createdBeforeOrAtHeight?.toString(),
        createdBeforeOrAt,
      },
      [],
      { orderBy: fillOrderBy },
    );

    const orderTypeMap = await buildOrderTypeMap(fills);
    const clobPairIdToMarket = buildClobPairIdToMarket();
    const allRows = computeTradeHistory(fills, orderTypeMap, clobPairIdToMarket);

    const effectiveLimit = limit ?? config.API_LIMIT_V4;
    const paginated = paginateTradeHistory(allRows, effectiveLimit, page);

    return {
      tradeHistory: paginated.tradeHistory,
      pageSize: paginated.pageSize,
      totalResults: paginated.totalResults,
      offset: paginated.offset,
    };
  }
}

// ---------- Route 1: GET / (single subaccount) ----------

router.get(
  '/',
  rateLimiterMiddleware(fillsRateLimiter),
  tradeHistoryCacheControlMiddleware,
  ...CheckSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtSchema,
  ...CheckPaginationSchema,
  query('market').if(query('marketType').exists()).notEmpty()
    .withMessage('market must be provided if marketType is provided'),
  query('marketType').if(query('market').exists()).notEmpty()
    .withMessage('marketType must be provided if market is provided'),
  ...marketTypeCheckSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
      market,
      marketType,
      limit,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
      page,
    }: TradeHistoryRequest = matchedData(req) as TradeHistoryRequest;

    const subaccountNum: number = +subaccountNumber;

    try {
      const controller: TradeHistoryController = new TradeHistoryController();
      const response: TradeHistoryResponse = await controller.getTradeHistory(
        address,
        subaccountNum,
        market,
        marketType,
        limit,
        createdBeforeOrAtHeight,
        createdBeforeOrAt,
        page,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'TradeHistoryController GET /',
        'Trade history error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_trade_history.timing`,
        Date.now() - start,
      );
    }
  },
);

// ---------- Route 2: GET /parentSubaccountNumber (parent subaccount) ----------

router.get(
  '/parentSubaccountNumber',
  rateLimiterMiddleware(fillsRateLimiter),
  tradeHistoryCacheControlMiddleware,
  ...CheckParentSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtSchema,
  ...CheckPaginationSchema,
  query('market').if(query('marketType').exists()).notEmpty()
    .withMessage('market must be provided if marketType is provided'),
  query('marketType').if(query('market').exists()).notEmpty()
    .withMessage('marketType must be provided if market is provided'),
  ...marketTypeCheckSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      parentSubaccountNumber,
      market,
      marketType,
      limit,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
      page,
    }: ParentSubaccountTradeHistoryRequest = matchedData(
      req,
    ) as ParentSubaccountTradeHistoryRequest;

    const parentSubaccountNum: number = +parentSubaccountNumber;

    try {
      const controller: TradeHistoryController = new TradeHistoryController();
      const response: TradeHistoryResponse = await controller
        .getTradeHistoryForParentSubaccount(
          address,
          parentSubaccountNum,
          market,
          marketType,
          limit,
          createdBeforeOrAtHeight,
          createdBeforeOrAt,
          page,
        );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'TradeHistoryController GET /parentSubaccountNumber',
        'Trade history error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_trade_history_parent.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
