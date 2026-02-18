import { cacheControlMiddleware, stats } from '@dydxprotocol-indexer/base';
import {
  FillColumns,
  FillFromDatabase,
  FillTable,
  OrderTable,
  OrderType,
  Ordering,
  SubaccountTable,
  SubaccountFromDatabase,
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
  CheckLimitSchema,
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
const TRADE_HISTORY_MAX_FILLS: number = 100_000;
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

// Shared: resolve clobPairId from market/marketType params
async function resolveClobPairId(
  market?: string,
  marketType?: MarketType,
): Promise<string | undefined> {
  if (isDefined(market) && isDefined(marketType)) {
    const clobPairId = await getClobPairId(market!, marketType!);
    if (clobPairId === undefined) {
      throw new NotFoundError(`${market} not found in markets of type ${marketType}`);
    }
    return clobPairId;
  }
  return undefined;
}

// Shared: fast fill count check using pagination (COUNT(*) only, loads 1 row)
async function checkFillCount(
  queryParams: object,
): Promise<void> {
  const { total } = await FillTable.findAll(
    { ...queryParams, limit: 1, page: 1 },
    [],
    { orderBy: fillOrderBy },
  );
  if (total !== undefined && total > TRADE_HISTORY_MAX_FILLS) {
    throw new Error(
      `Too many fills (${total}) to compute trade history.`,
    );
  }
}

// Shared: compute trade history from fills and return paginated response
async function buildTradeHistoryResponse(
  fills: FillFromDatabase[],
  subaccountIdToNumber: Record<string, number>,
  limit?: number,
  page?: number,
): Promise<TradeHistoryResponse> {
  const orderTypeMap = await buildOrderTypeMap(fills);
  const clobPairIdToMarket = buildClobPairIdToMarket();
  const allRows = computeTradeHistory(fills, orderTypeMap, clobPairIdToMarket,
    subaccountIdToNumber);

  const effectiveLimit = limit ?? config.API_LIMIT_V4;
  const paginated = paginateTradeHistory(allRows, effectiveLimit, page);

  return {
    tradeHistory: paginated.tradeHistory,
    pageSize: paginated.pageSize,
    totalResults: paginated.totalResults,
    offset: paginated.offset,
  };
}

@Route('tradeHistory')
class TradeHistoryController extends Controller {
  @Get('/')
  async getTradeHistory(
    @Query() address: string,
      @Query() subaccountNumber: number,
      @Query() market?: string,
      @Query() marketType?: MarketType,
      @Query() limit?: number,
      @Query() page?: number,
  ): Promise<TradeHistoryResponse> {
    const clobPairId = await resolveClobPairId(market, marketType);
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);

    await checkFillCount({ subaccountId: [subaccountId], clobPairId });

    const { results: fills } = await FillTable.findAll(
      { subaccountId: [subaccountId], clobPairId },
      [],
      { orderBy: fillOrderBy },
    );

    return buildTradeHistoryResponse(
      fills, { [subaccountId]: subaccountNumber }, limit, page,
    );
  }

  @Get('/parentSubaccountNumber')
  async getTradeHistoryForParentSubaccount(
    @Query() address: string,
      @Query() parentSubaccountNumber: number,
      @Query() market?: string,
      @Query() marketType?: MarketType,
      @Query() limit?: number,
      @Query() page?: number,
  ): Promise<TradeHistoryResponse> {
    const clobPairId = await resolveClobPairId(market, marketType);

    await checkFillCount({
      parentSubaccount: { address, subaccountNumber: parentSubaccountNumber },
      clobPairId,
    });

    const { results: fills } = await FillTable.findAll(
      {
        parentSubaccount: { address, subaccountNumber: parentSubaccountNumber },
        clobPairId,
      },
      [],
      { orderBy: fillOrderBy },
    );

    // Build subaccountId -> subaccountNumber map from the fills
    const uniqueSubaccountIds = _.uniq(fills.map((f) => f.subaccountId));
    const subaccounts: SubaccountFromDatabase[] = uniqueSubaccountIds.length > 0
      ? await SubaccountTable.findAll({ id: uniqueSubaccountIds }, [])
      : [];
    const subaccountIdToNumber: Record<string, number> = {};
    for (const sa of subaccounts) {
      subaccountIdToNumber[sa.id] = sa.subaccountNumber;
    }

    return buildTradeHistoryResponse(fills, subaccountIdToNumber, limit, page);
  }
}

// ---------- Route 1: GET / (single subaccount) ----------

router.get(
  '/',
  rateLimiterMiddleware(fillsRateLimiter),
  tradeHistoryCacheControlMiddleware,
  ...CheckSubaccountSchema,
  ...CheckLimitSchema,
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
  ...CheckLimitSchema,
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
