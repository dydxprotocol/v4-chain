import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  APIOrderStatus,
  APIOrderStatusEnum,
  DEFAULT_POSTGRES_OPTIONS,
  IsoString,
  OrderColumns,
  OrderFromDatabase,
  Ordering,
  OrderSide,
  OrderTable,
  OrderType,
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

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { redisClient } from '../../../helpers/redis/redis-controller';
import { complianceCheck } from '../../../lib/compliance-check';
import { NotFoundError } from '../../../lib/errors';
import {
  handleControllerError,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { rejectRestrictedCountries } from '../../../lib/restrict-countries';
import {
  CheckLimitSchema,
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
  GetOrderRequest, ListOrderRequest, OrderResponseObject, PostgresOrderMap, RedisOrderMap,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'orders-controller';

@Route('orders')
class OrdersController extends Controller {
  @Get('/')
  async listOrders(
    @Query() address: string,
      @Query() subaccountNumber: number,
      @Query() limit: number,
      @Query() ticker?: string,
      @Query() side?: OrderSide,
      @Query() type?: OrderType,
      @Query() status?: APIOrderStatus[],
      @Query() goodTilBlockBeforeOrAt?: number,
      @Query() goodTilBlockTimeBeforeOrAt?: IsoString,
      @Query() returnLatestOrders?: boolean,
  ): Promise<OrderResponseObject[]> {
    let clobPairId: string | undefined;
    if (ticker !== undefined) {
      clobPairId = perpetualMarketRefresher.getClobPairIdFromTicker(ticker);
    }

    const ordering: Ordering = returnLatestOrders !== undefined
      ? returnLatestOrdersToOrdering(returnLatestOrders)
      : Ordering.DESC;
    const [
      redisOrderMap,
      postgresOrders,
    ]: [
      RedisOrderMap,
      OrderFromDatabase[],
    ] = await Promise.all([
      getRedisOrderMapForSubaccountId(
        SubaccountTable.uuid(address, subaccountNumber),
        clobPairId,
        side,
        type,
        goodTilBlockBeforeOrAt,
        goodTilBlockTimeBeforeOrAt,
      ),
      OrderTable.findAll(
        {
          subaccountId: [SubaccountTable.uuid(address, subaccountNumber)],
          limit,
          clobPairId,
          side,
          type,
          goodTilBlockBeforeOrAt: goodTilBlockBeforeOrAt?.toString(),
          goodTilBlockTimeBeforeOrAt,
        }, [], {
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

    const postgresOrderMap: PostgresOrderMap = _.keyBy(postgresOrders, OrderColumns.id);

    let mergedResponses: OrderResponseObject[] = mergePostgresAndRedisOrdersToResponseObjects(
      postgresOrderMap,
      redisOrderMap,
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
      limit,
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
      OrdersCache.getOrder(orderId, redisClient),
    ]);

    const order: OrderResponseObject | undefined = postgresAndRedisOrderToResponseObject(
      postgresOrder,
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
  rejectRestrictedCountries,
  rateLimiterMiddleware(getReqRateLimiter),
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
    goodTilBlockTimeBeforeOrAt: {
      in: 'query',
      optional: true,
      isISO8601: true,
    },
    returnLatestOrders: {
      in: 'query',
      isBoolean: true,
      optional: true,
    },
  }),
  query('goodTilBlock').if(query('goodTilBlockTime').exists()).isEmpty()
    .withMessage('Cannot provide both goodTilBlock and goodTilBlockTime'),
  handleValidationErrors,
  complianceCheck,
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
      status,
      goodTilBlockBeforeOrAt,
      goodTilBlockTimeBeforeOrAt,
      returnLatestOrders,
    }: ListOrderRequest = matchedData(req) as ListOrderRequest;

    try {
      const controller: OrdersController = new OrdersController();
      const response: OrderResponseObject[] = await controller.listOrders(
        address,
        subaccountNumber,
        limit,
        ticker,
        side,
        type,
        status,
        goodTilBlockBeforeOrAt,
        goodTilBlockTimeBeforeOrAt,
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
  '/:orderId',
  rejectRestrictedCountries,
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
 * Gets a RedisOrderMap filtered by subaccount, and if provided,
 * by clobPairId, side, type, goodTilBlock, and goodTilBlockTime
 * Note: When filtering by `goodTilBlock` all orders without a `goodTilBlock` will be filtered out
 * and the same with `goodTilBlockTime`. Both cannot be provided as that should be an invalid query.
 * @param subaccountId
 * @param clobPairId
 * @param side
 * @param type
 * @param goodTilBlockBeforeOrAt
 * @param goodTilBlockTimeBeforeOrAt
 * @returns
 */
async function getRedisOrderMapForSubaccountId(
  subaccountId: string,
  clobPairId?: string,
  side?: OrderSide,
  type?: OrderType,
  goodTilBlockBeforeOrAt?: number,
  goodTilBlockTimeBeforeOrAt?: IsoString,
): Promise<RedisOrderMap> {
  if (type !== undefined && type !== OrderType.LIMIT) {
    // TODO(DEC-1458): Add support for advanced Orders
    // We don't currently support non LIMIT orders in Redis
    return {};
  }

  const subaccountOrderIds: string[] = await SubaccountOrderIdsCache.getOrderIdsForSubaccount(
    subaccountId,
    redisClient,
  );

  const nullableRedisOrders: (RedisOrder | null)[] = await Promise.all(
    _.map(subaccountOrderIds, (orderId: string) => OrdersCache.getOrder(orderId, redisClient)),
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
      } else {
        // If `goodTilBlockBeforeOrAt` is defined as a filter, filter out all orders that don't have
        // `goodTilBlock` defined
        if (goodTilBlockBeforeOrAt !== undefined) {
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
      } else {
        if (goodTilBlockTimeBeforeOrAt !== undefined) {
          // If `goodTilBlockTimeBeforeOrAt` is defined as a filter, filter out all orders that
          // don't have `goodTilBlockTime` defined
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
