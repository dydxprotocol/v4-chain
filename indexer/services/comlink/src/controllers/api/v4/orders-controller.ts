import { cacheControlMiddleware, logger, stats } from '@dydxprotocol-indexer/base';
import {
  APIOrderStatus,
  APIOrderStatusEnum,
  DEFAULT_POSTGRES_OPTIONS,
  IsoString,
  OrderColumns,
  OrderFromDatabase,
  Ordering,
  OrderQueryConfig,
  OrderSide,
  OrderStatus,
  OrderTable,
  OrderType,
  PaginationFromDatabase,
  ParentSubaccount,
  perpetualMarketRefresher,
  protocolTranslations,
  SubaccountTable,
} from '@dydxprotocol-indexer/postgres';
import { OrdersCache, SubaccountOrderIdsCache } from '@dydxprotocol-indexer/redis';
import { RedisOrder } from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import express from 'express';
import { checkSchema, matchedData, query } from 'express-validator';
import _ from 'lodash';
import { DateTime } from 'luxon';
import {
  Controller, Get, Path, Query, Route,
} from 'tsoa';

import { ordersRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { redisReadOnlyClient } from '../../../helpers/redis/redis-controller';
import { complianceAndGeoCheck } from '../../../lib/compliance-and-geo-check';
import { NotFoundError } from '../../../lib/errors';
import {
  getChildSubaccountNums,
  handleControllerError,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import {
  CheckLimitSchema, CheckParentSubaccountSchema,
  CheckSubaccountSchema,
  CheckTickerOptionalQuerySchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import {
  mergePostgresAndRedisOrdersToResponseObjects,
  postgresAndRedisOrderToResponseObject,
} from '../../../request-helpers/request-transformer';
import { sanitizeArray } from '../../../request-helpers/sanitizers';
import { validateArray } from '../../../request-helpers/validators';
import {
  GetOrderRequest,
  ListOrderRequest,
  OrderResponseObject,
  ParentSubaccountListOrderRequest,
  PostgresOrderMap,
  RedisOrderMap,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'orders-controller';
const ordersCacheControlMiddleware = cacheControlMiddleware(config.CACHE_CONTROL_DIRECTIVE_ORDERS);

/**
 * Lists orders for a set of subaccounts based on various filters.
 * @param subaccountIdToNumber A mapping of subaccount IDs to their corresponding numbers.
 * @param limit The maximum number of orders to return.
 * @param ticker Optional ticker to filter orders by.
 * @param side Optional order side to filter orders by.
 * @param type Optional order type to filter orders by.
 * @param status Optional array of order statuses to filter orders by.
 * @param goodTilBlockBeforeOrAt Optional filter for orders good until a specific block.
 * @param goodTilBlockTimeBeforeOrAt Optional filter for orders good until a specific time.
 * @param returnLatestOrders Flag indicating whether to return the latest orders.
 * @returns An array of order response objects.
 */
async function listOrdersCommon(
  subaccountIdToNumber: Record<string, number>,
  parentSubaccount?: ParentSubaccount,
  limit?: number,
  ticker?: string,
  side?: OrderSide,
  type?: OrderType,
  includeTypes?: OrderType[],
  excludeTypes?: OrderType[],
  status?: APIOrderStatus[],
  goodTilBlockBeforeOrAt?: number,
  goodTilBlockAfter?: number,
  goodTilBlockTimeBeforeOrAt?: IsoString,
  goodTilBlockTimeAfter?: IsoString,
  returnLatestOrders?: boolean,
): Promise<OrderResponseObject[]> {
  let clobPairId: string | undefined;
  if (ticker !== undefined) {
    clobPairId = perpetualMarketRefresher.getClobPairIdFromTicker(ticker);
  }

  const subaccountIds: string[] = Object.keys(subaccountIdToNumber);

  const orderQueryConfig: OrderQueryConfig = {
    limit,
    clobPairId,
    side,
    type,
    includeTypes,
    excludeTypes,
    goodTilBlockBeforeOrAt: goodTilBlockBeforeOrAt?.toString(),
    goodTilBlockAfter: goodTilBlockAfter?.toString(),
    goodTilBlockTimeBeforeOrAt,
    goodTilBlockTimeAfter,
  };
  if (parentSubaccount !== undefined) {
    orderQueryConfig.parentSubaccount = parentSubaccount;
  } else {
    orderQueryConfig.subaccountId = subaccountIds;
  }

  if (!_.isEmpty(status)) {
    // BEST_EFFORT_OPENED status is not filtered out, because it's a minor optimization,
    // is more confusing, and is not going to affect the result of the query.
    orderQueryConfig.statuses = status as OrderStatus[];
  }
  const ordering: Ordering = returnLatestOrders !== undefined
    ? returnLatestOrdersToOrdering(returnLatestOrders)
    : Ordering.DESC;
  const [
    redisOrderMap,
    { results: postgresOrders },
  ]: [
    RedisOrderMap,
    PaginationFromDatabase<OrderFromDatabase>,
  ] = await Promise.all([
    getRedisOrderMapForSubaccountIds(
      subaccountIds,
      clobPairId,
      side,
      type,
      goodTilBlockBeforeOrAt,
      goodTilBlockAfter,
      goodTilBlockTimeBeforeOrAt,
      goodTilBlockTimeAfter,
    ),
    OrderTable.findAll(
      orderQueryConfig, [], {
        ...DEFAULT_POSTGRES_OPTIONS,
        orderBy: [
          // Order by `goodTilBlock` and then order by `goodTilBlockTime`
          // This way, orders with `goodTilBlock` defined are ordered before orders with
          // `goodTilBlockTime` if order is asecnding, and after if descending
          [OrderColumns.goodTilBlock, ordering],
          [OrderColumns.goodTilBlockTime, ordering],
        ],
      },
    ),
  ]);

  const redisOrderIds: string[] = _.map(
    Object.values(redisOrderMap),
    (redisOrder: RedisOrder) => {
      return OrderTable.orderIdToUuid(redisOrder.order!.orderId!);
    },
  );
  const postgresOrderIdsToFetch: string[] = _.difference(
    redisOrderIds,
    _.map(postgresOrders, OrderColumns.id),
  );

  // Postgres is regarded as the source of truth, so for any redis orders not returned from the
  // initial postgres query, we need to fetch them from Postgres to ensure we have the most
  // accurate status. For example, if the user is querying for `status: [BEST_EFFORT_OPENED]`,
  // we need to fetch all orders from Postgres, because if the order in postgres is 'OPENED',
  // then we do not want to return this order to the user as 'BEST_EFFORT_OPENED'.
  let additionalPostgresOrders: OrderFromDatabase[] = [];
  if (!_.isEmpty(postgresOrderIdsToFetch)) {
    const { results }: PaginationFromDatabase<OrderFromDatabase> = await OrderTable.findAll({
      id: postgresOrderIdsToFetch,
    }, [], {
      ...DEFAULT_POSTGRES_OPTIONS,
    });

    additionalPostgresOrders = results;
  }

  const postgresOrderMap: PostgresOrderMap = _.keyBy(
    _.concat(postgresOrders, additionalPostgresOrders),
    OrderColumns.id,
  );

  let mergedResponses: OrderResponseObject[] = mergePostgresAndRedisOrdersToResponseObjects(
    postgresOrderMap,
    redisOrderMap,
    subaccountIdToNumber,
  );

  if (status !== undefined) {
    mergedResponses = _.filter(
      mergedResponses,
      (orderResponse: OrderResponseObject, _index: number) => {
        return status.includes(orderResponse.status);
      },
    );
  }

  return sortAndLimitResponses(
    mergedResponses,
    ordering,
    limit || config.API_LIMIT_V4,
  );

}

@Route('orders')
class OrdersController extends Controller {
  @Get('/')
  async listOrders(
    @Query() address: string,
      @Query() subaccountNumber: number,
      @Query() limit?: number,
      @Query() ticker?: string,
      @Query() side?: OrderSide,
      @Query() type?: OrderType,
      @Query() includeTypes?: OrderType[],
      @Query() excludeTypes?: OrderType[],
      @Query() status?: APIOrderStatus[],
      @Query() goodTilBlockBeforeOrAt?: number,
      @Query() goodTilBlockAfter?: number,
      @Query() goodTilBlockTimeBeforeOrAt?: IsoString,
      @Query() goodTilBlockTimeAfter?: IsoString,
      @Query() returnLatestOrders?: boolean,
  ): Promise<OrderResponseObject[]> {

    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);

    return listOrdersCommon(
      { [subaccountId]: subaccountNumber },
      undefined, // parentSubaccount
      limit,
      ticker,
      side,
      type,
      includeTypes,
      excludeTypes,
      status,
      goodTilBlockBeforeOrAt,
      goodTilBlockAfter,
      goodTilBlockTimeBeforeOrAt,
      goodTilBlockTimeAfter,
      returnLatestOrders,
    );
  }

  @Get('/parentSubaccountNumber')
  async listOrdersForParentSubaccount(
    @Query() address: string,
      @Query() parentSubaccountNumber: number,
      @Query() limit?: number,
      @Query() ticker?: string,
      @Query() side?: OrderSide,
      @Query() type?: OrderType,
      @Query() includeTypes?: OrderType[],
      @Query() excludeTypes?: OrderType[],
      @Query() status?: APIOrderStatus[],
      @Query() goodTilBlockBeforeOrAt?: number,
      @Query() goodTilBlockAfter?: number,
      @Query() goodTilBlockTimeBeforeOrAt?: IsoString,
      @Query() goodTilBlockTimeAfter?: IsoString,
      @Query() returnLatestOrders?: boolean,
  ): Promise<OrderResponseObject[]> {
    const childIdtoSubaccountNumber: Record<string, number> = {};
    getChildSubaccountNums(parentSubaccountNumber).forEach(
      (subaccountNum: number) => {
        childIdtoSubaccountNumber[SubaccountTable.uuid(address, subaccountNum)] = subaccountNum;
      },
    );

    return listOrdersCommon(
      childIdtoSubaccountNumber,
      {
        address,
        subaccountNumber: parentSubaccountNumber,
      }, // parentSubaccount
      limit,
      ticker,
      side,
      type,
      includeTypes,
      excludeTypes,
      status,
      goodTilBlockBeforeOrAt,
      goodTilBlockAfter,
      goodTilBlockTimeBeforeOrAt,
      goodTilBlockTimeAfter,
      returnLatestOrders,
    );
  }

  @Get('/:orderId')
  async getOrder(
    @Path() orderId: string,
  ): Promise<OrderResponseObject> {
    const [
      postgresOrder,
      redisOrder,
    ]: [
      OrderFromDatabase | undefined,
      RedisOrder | null,
    ] = await Promise.all([
      OrderTable.findById(orderId),
      OrdersCache.getOrder(orderId, redisReadOnlyClient),
    ]);

    // Get subaccount number and subaccountId from either Redis or Postgres
    let subaccountNumber: number | undefined;
    let subaccountId: string | undefined;
    if (redisOrder !== null) {
      subaccountNumber = redisOrder.order!.orderId!.subaccountId!.number;
      subaccountId = SubaccountTable.uuid(
        redisOrder.order!.orderId!.subaccountId!.owner,
        subaccountNumber,
      );
    } else if (postgresOrder !== undefined) {
      const subaccount = await SubaccountTable.findById(postgresOrder.subaccountId);
      if (subaccount === undefined) {
        throw new NotFoundError(`Unable to find subaccount id ${postgresOrder.subaccountId}`);
      }
      subaccountNumber = subaccount.subaccountNumber;
      subaccountId = postgresOrder.subaccountId;
    } else {
      throw new NotFoundError(`Unable to find order id ${orderId}`);
    }

    const order: OrderResponseObject | undefined = postgresAndRedisOrderToResponseObject(
      postgresOrder,
      { [subaccountId]: subaccountNumber },
      redisOrder,
    );
    if (order === undefined) {
      throw new NotFoundError(`Unable to find order id ${orderId}`);
    }

    return order;
  }
}

router.get(
  '/',
  rateLimiterMiddleware(ordersRateLimiter),
  ordersCacheControlMiddleware,
  ...CheckSubaccountSchema,
  ...CheckLimitSchema,
  ...CheckTickerOptionalQuerySchema,
  ...checkSchema({
    side: {
      in: 'query',
      isIn: {
        options: [Object.values(OrderSide)],
        errorMessage: `side must be one of ${Object.values(OrderSide)}`,
      },
      optional: true,
    },
    type: {
      in: 'query',
      isIn: {
        options: [Object.values(OrderType)],
        errorMessage: `type must be one of ${Object.values(OrderType)}`,
      },
      optional: true,
    },
    // TODO(DEC-1462): Add /active-orders endpoint fetching mainly from Redis once fully-filled
    // orders are removed from Redis. Until then, orders have to be merged with Postgres orders
    // to get the correct status.
    status: {
      in: ['query'],
      optional: true,
      customSanitizer: {
        options: sanitizeArray,
      },
      custom: {
        options: (inputArray) => validateArray(inputArray, Object.values(APIOrderStatusEnum)),
        errorMessage: `status must be one of ${Object.values(APIOrderStatusEnum)}`,
      },
    },
    goodTilBlockBeforeOrAt: {
      in: 'query',
      optional: true,
      isInt: {
        options: { gt: 0 },
      },
    },
    goodTilBlockAfter: {
      in: 'query',
      optional: true,
      isInt: {
        options: { gt: 0 },
      },
    },
    goodTilBlockTimeBeforeOrAt: {
      in: 'query',
      optional: true,
      isISO8601: true,
    },
    goodTilBlockTimeAfter: {
      in: 'query',
      optional: true,
      isISO8601: true,
    },
    returnLatestOrders: {
      in: 'query',
      isBoolean: true,
      optional: true,
    },
    includeTypes: {
      in: ['query'],
      optional: true,
      customSanitizer: {
        options: sanitizeArray,
      },
      custom: {
        options: (inputArray) => validateArray(inputArray, Object.values(OrderType)),
        errorMessage: `includeTypes must be one of ${Object.values(OrderType)}`,
      },
    },
    excludeTypes: {
      in: ['query'],
      optional: true,
      customSanitizer: {
        options: sanitizeArray,
      },
      custom: {
        options: (inputArray) => validateArray(inputArray, Object.values(OrderType)),
        errorMessage: `excludeTypes must be one of ${Object.values(OrderType)}`,
      },
    },
  }),
  query('goodTilBlock').if(query('goodTilBlockTime').exists()).isEmpty()
    .withMessage('Cannot provide both goodTilBlock and goodTilBlockTime'),
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
      limit,
      ticker,
      side,
      type,
      includeTypes,
      excludeTypes,
      status,
      goodTilBlockBeforeOrAt,
      goodTilBlockAfter,
      goodTilBlockTimeBeforeOrAt,
      goodTilBlockTimeAfter,
      returnLatestOrders,
    }: ListOrderRequest = matchedData(req) as ListOrderRequest;

    // The schema checks allow subaccountNumber to be a string, but we know it's a number here.
    const subaccountNum: number = +subaccountNumber;

    try {
      const controller: OrdersController = new OrdersController();
      const response: OrderResponseObject[] = await controller.listOrders(
        address,
        subaccountNum,
        limit,
        ticker,
        side,
        type,
        includeTypes,
        excludeTypes,
        status,
        goodTilBlockBeforeOrAt,
        goodTilBlockAfter,
        goodTilBlockTimeBeforeOrAt,
        goodTilBlockTimeAfter,
        returnLatestOrders,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'OrdersController GET /',
        'Orders error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.list_orders.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/parentSubaccountNumber',
  rateLimiterMiddleware(ordersRateLimiter),
  ordersCacheControlMiddleware,
  ...CheckParentSubaccountSchema,
  ...CheckLimitSchema,
  ...CheckTickerOptionalQuerySchema,
  ...checkSchema({
    side: {
      in: 'query',
      isIn: {
        options: [Object.values(OrderSide)],
        errorMessage: `side must be one of ${Object.values(OrderSide)}`,
      },
      optional: true,
    },
    type: {
      in: 'query',
      isIn: {
        options: [Object.values(OrderType)],
        errorMessage: `type must be one of ${Object.values(OrderType)}`,
      },
      optional: true,
    },
    // TODO(DEC-1462): Add /active-orders endpoint fetching mainly from Redis once fully-filled
    // orders are removed from Redis. Until then, orders have to be merged with Postgres orders
    // to get the correct status.
    status: {
      in: ['query'],
      optional: true,
      customSanitizer: {
        options: sanitizeArray,
      },
      custom: {
        options: (inputArray) => validateArray(inputArray, Object.values(APIOrderStatusEnum)),
        errorMessage: `status must be one of ${Object.values(APIOrderStatusEnum)}`,
      },
    },
    goodTilBlockBeforeOrAt: {
      in: 'query',
      optional: true,
      isInt: {
        options: { gt: 0 },
      },
    },
    goodTilBlockAfter: {
      in: 'query',
      optional: true,
      isInt: {
        options: { gt: 0 },
      },
    },
    goodTilBlockTimeBeforeOrAt: {
      in: 'query',
      optional: true,
      isISO8601: true,
    },
    goodTilBlockTimeAfter: {
      in: 'query',
      optional: true,
      isISO8601: true,
    },
    returnLatestOrders: {
      in: 'query',
      isBoolean: true,
      optional: true,
    },
    includeTypes: {
      in: ['query'],
      optional: true,
      customSanitizer: {
        options: sanitizeArray,
      },
      custom: {
        options: (inputArray) => validateArray(inputArray, Object.values(OrderType)),
        errorMessage: `includeTypes must be one of ${Object.values(OrderType)}`,
      },
    },
    excludeTypes: {
      in: ['query'],
      optional: true,
      customSanitizer: {
        options: sanitizeArray,
      },
      custom: {
        options: (inputArray) => validateArray(inputArray, Object.values(OrderType)),
        errorMessage: `excludeTypes must be one of ${Object.values(OrderType)}`,
      },
    },
  }),
  query('goodTilBlock').if(query('goodTilBlockTime').exists()).isEmpty()
    .withMessage('Cannot provide both goodTilBlock and goodTilBlockTime'),
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      parentSubaccountNumber,
      limit,
      ticker,
      side,
      type,
      includeTypes,
      excludeTypes,
      status,
      goodTilBlockBeforeOrAt,
      goodTilBlockAfter,
      goodTilBlockTimeBeforeOrAt,
      goodTilBlockTimeAfter,
      returnLatestOrders,
    }: ParentSubaccountListOrderRequest = matchedData(req) as ParentSubaccountListOrderRequest;

    // The schema checks allow subaccountNumber to be a string, but we know it's a number here.
    const parentSubaccountNum: number = +parentSubaccountNumber;

    try {
      const controller: OrdersController = new OrdersController();
      const response: OrderResponseObject[] = await controller.listOrdersForParentSubaccount(
        address,
        parentSubaccountNum,
        limit,
        ticker,
        side,
        type,
        includeTypes,
        excludeTypes,
        status,
        goodTilBlockBeforeOrAt,
        goodTilBlockAfter,
        goodTilBlockTimeBeforeOrAt,
        goodTilBlockTimeAfter,
        returnLatestOrders,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'OrdersController GET /parentSubaccountNumber',
        'Orders error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.list_orders_parent_subaccount.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/:orderId',
  ordersCacheControlMiddleware,
  ...checkSchema({
    orderId: {
      in: ['params'],
      isUUID: true,
    },
  }),
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      orderId,
    }: GetOrderRequest = matchedData(req) as GetOrderRequest;

    try {
      const controller: OrdersController = new OrdersController();
      const response: OrderResponseObject = await controller.getOrder(orderId);

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'OrdersController GET /:orderId',
        'Orders error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_orders.timing`,
        Date.now() - start,
      );
    }
  },
);

/* ------- ORDER HELPERS ------- */

/**
 * Gets a RedisOrderMap filtered by subaccounts, and if provided,
 * by clobPairId, side, type, goodTilBlock, and goodTilBlockTime
 * Note: When filtering by `goodTilBlock` all orders without a `goodTilBlock` will be filtered out
 * and the same with `goodTilBlockTime`. Both cannot be provided as that should be an invalid query.
 * @param subaccountIds
 * @param clobPairId
 * @param side
 * @param type
 * @param goodTilBlockBeforeOrAt
 * @param goodTilBlockTimeBeforeOrAt
 * @returns
 */
async function getRedisOrderMapForSubaccountIds(
  subaccountIds: string[],
  clobPairId?: string,
  side?: OrderSide,
  type?: OrderType,
  goodTilBlockBeforeOrAt?: number,
  goodTilBlockAfter?: number,
  goodTilBlockTimeBeforeOrAt?: IsoString,
  goodTilBlockTimeAfter?: IsoString,
): Promise<RedisOrderMap> {
  if (type !== undefined && type !== OrderType.LIMIT) {
    // TODO(DEC-1458): Add support for advanced Orders
    // We don't currently support non LIMIT orders in Redis
    return {};
  }

  const subaccountToOrderIds = await SubaccountOrderIdsCache.getOrderIdsForSubaccounts(
    subaccountIds,
    redisReadOnlyClient,
  );
  const orderIds: string[] = _.flatten(_.values(subaccountToOrderIds));

  const nullableRedisOrders: (RedisOrder | null)[] = await Promise.all(
    _.map(orderIds, (orderId: string) => OrdersCache.getOrder(orderId, redisReadOnlyClient)),
  );
  const redisOrders: RedisOrder[] = _.filter(
    nullableRedisOrders,
    (redisOrder: RedisOrder | null) => {
      if (redisOrder === null) {
        return false;
      }

      const redisClobPairId: string = redisOrder!.order!.orderId!.clobPairId.toString();
      if (clobPairId !== undefined && redisClobPairId !== clobPairId) {
        return false;
      }

      const redisSide: OrderSide = protocolTranslations.protocolOrderSideToOrderSide(
        redisOrder!.order!.side,
      );
      if (side !== undefined && side !== redisSide) {
        return false;
      }

      const redisGoodTilBlock: number | undefined = protocolTranslations
        .getGoodTilBlock(redisOrder!.order!);
      if (redisGoodTilBlock !== undefined) {
        if (goodTilBlockBeforeOrAt !== undefined && redisGoodTilBlock > goodTilBlockBeforeOrAt) {
          return false;
        }
        if (goodTilBlockAfter !== undefined && redisGoodTilBlock <= goodTilBlockAfter) {
          return false;
        }
      } else {
        // If `goodTilBlockBeforeOrAt` or `goodTilBlockAfter` is defined as a filter, filter out all
        // orders that don't have `goodTilBlock` defined
        if (goodTilBlockBeforeOrAt !== undefined || goodTilBlockAfter !== undefined) {
          return false;
        }
      }

      const redisGoodTilBlockTime: string | undefined = protocolTranslations
        .getGoodTilBlockTime(redisOrder!.order!);
      if (redisGoodTilBlockTime) {
        const redisGoodTilBlockTimeDateObj: DateTime = DateTime.fromISO(redisGoodTilBlockTime);
        if (goodTilBlockTimeBeforeOrAt !== undefined &&
          redisGoodTilBlockTimeDateObj > DateTime.fromISO(goodTilBlockTimeBeforeOrAt)
        ) {
          return false;
        }
        if (goodTilBlockTimeAfter !== undefined &&
          redisGoodTilBlockTimeDateObj <= DateTime.fromISO(goodTilBlockTimeAfter)
        ) {
          return false;
        }
      } else {
        if (goodTilBlockTimeBeforeOrAt !== undefined || goodTilBlockTimeAfter !== undefined) {
          // If `goodTilBlockTimeBeforeOrAt` or `goodTilBlockTimeAfter` is defined as a filter,
          // filter out all orders that don't have `goodTilBlockTime` defined
          return false;
        }
      }

      return true;
    },
  ) as RedisOrder[];

  return _.keyBy(redisOrders, 'id');
}

function returnLatestOrdersToOrdering(
  returnLatestOrders: boolean,
): Ordering {
  return returnLatestOrders === true ? Ordering.DESC : Ordering.ASC;
}

/**
 * Sorts the orders based on the ordering provided. If ordering is ASC, then lowest
 * goodTilBlock is first. If ordering is DESC, then highest goodTilBlock is first.
 * Then limits the number of orders to the limit provided.
 * @param orderResponses
 * @param ordering
 * @param limit
 * @returns
 */
function sortAndLimitResponses(
  orderResponses: OrderResponseObject[],
  ordering: Ordering,
  limit: number,
): OrderResponseObject[] {
  const sortedResponses: OrderResponseObject[] = orderResponses.sort(
    (a: OrderResponseObject, b: OrderResponseObject): number => (ordering === Ordering.ASC
      ? compareOrderResponses(a, b)
      : compareOrderResponses(b, a)),
  );
  return sortedResponses.slice(0, limit);
}

/**
 * Compares 2 OrderResponseObjects a and b.
 * Return:
 * - 1 if a is greater or equal to b
 * - -1 if a is less than b
 * All orders with `goodTilBlockTime` defined are considered to be greater than any orders with
 * `goodTilBlock` defined, and vice-versa all orders with `goodTilBlock` defined are considered to
 * be less than orders with `goodTilBlockTime` defined.
 * @param a
 * @param b
 * @returns
 */
function compareOrderResponses(a: OrderResponseObject, b: OrderResponseObject): number {
  // Orders with `goodTilBlock` should be ordered before orders with `goodTilBlockTime` in ascending
  // order
  if (isDefined(a.goodTilBlock) && isDefined(b.goodTilBlockTime)) {
    return -1;
  }
  if (isDefined(b.goodTilBlock) && isDefined(a.goodTilBlockTime)) {
    return 1;
  }

  if (isDefined(a.goodTilBlock) && isDefined(b.goodTilBlock)) {
    return Big(a.goodTilBlock!).lt(Big(b.goodTilBlock!)) ? -1 : 1;
  }
  if (isDefined(a.goodTilBlockTime) && isDefined(b.goodTilBlockTime)) {
    return DateTime.fromISO(a.goodTilBlockTime!) < DateTime.fromISO(b.goodTilBlockTime!) ? -1 : 1;
  }

  const errMessage: string = 'Order repsonse objects are invalid';
  logger.error({
    at: `${controllerName}#compareOrderResponses`,
    message: errMessage,
    orderA: a,
    orderB: b,
  });
  throw new Error(errMessage);
}

// eslint-disable-next-line  @typescript-eslint/no-explicit-any
function isDefined(val?: any) {
  return val !== null && val !== undefined;
}

export default router;
