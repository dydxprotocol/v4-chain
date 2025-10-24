import { cacheControlMiddleware, stats } from '@dydxprotocol-indexer/base';
import {
  BlockFromDatabase,
  BlockTable,
  FundingIndexMap,
  IsoString,
  MarketColumns,
  MarketFromDatabase,
  MarketsMap,
  MarketTable,
  perpetualMarketRefresher,
  PerpetualMarketsMap,
  PerpetualPositionFromDatabase,
  PerpetualPositionStatus,
  PerpetualPositionTable,
  QueryableField,
  SubaccountFromDatabase,
  SubaccountTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import {
  checkSchema,
  matchedData,
} from 'express-validator';
import _ from 'lodash';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { defaultRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceAndGeoCheck } from '../../../lib/compliance-and-geo-check';
import { NotFoundError } from '../../../lib/errors';
import {
  getChildSubaccountNums,
  getFundingIndexMaps,
  getPerpetualPositionsWithUpdatedFunding,
  handleControllerError,
  initializePerpetualPositionsWithFunding,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import {
  CheckLimitAndCreatedBeforeOrAtSchema, CheckParentSubaccountSchema,
  CheckSubaccountSchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { perpetualPositionToResponseObject } from '../../../request-helpers/request-transformer';
import { sanitizeArray } from '../../../request-helpers/sanitizers';
import { validateArray } from '../../../request-helpers/validators';
import {
  ParentSubaccountPerpetualPositionRequest,
  PerpetualPositionRequest,
  PerpetualPositionResponse,
  PerpetualPositionWithFunding,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'perpetual-positions-controller';
const perpetualPositionsCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_PERPETUAL_POSITIONS,
);

@Route('perpetualPositions')
class PerpetualPositionsController extends Controller {
  @Get('/')
  async listPositions(
    @Query() address: string,
      @Query() subaccountNumber: number,
      @Query() status?: PerpetualPositionStatus[],
      @Query() limit?: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
  ): Promise<PerpetualPositionResponse> {
    const subaccountUuid: string = SubaccountTable.uuid(address, subaccountNumber);

    const [
      positions,
      markets,
    ]: [
      PerpetualPositionFromDatabase[],
      MarketFromDatabase[],
    ] = await Promise.all([
      PerpetualPositionTable.findAll(
        {
          subaccountId: [subaccountUuid],
          status,
          limit,
          createdBeforeOrAtHeight: createdBeforeOrAtHeight
            ? createdBeforeOrAtHeight.toString()
            : undefined,
          createdBeforeOrAt,
        },
        [QueryableField.LIMIT],
      ),
      MarketTable.findAll(
        {},
        [],
      ),
    ]);

    const openPositionsExist: boolean = positions.some(
      (position: PerpetualPositionFromDatabase) => position.status === PerpetualPositionStatus.OPEN,
    );

    let updatedPerpetualPositions:
    PerpetualPositionWithFunding[] = initializePerpetualPositionsWithFunding(positions);
    if (openPositionsExist) {
      const [
        subaccount,
        latestBlock,
      ]: [
        SubaccountFromDatabase | undefined,
        BlockFromDatabase | undefined,
      ] = await Promise.all([
        SubaccountTable.findById(
          subaccountUuid,
        ),
        BlockTable.getLatest(),
      ]);

      if (subaccount === undefined || latestBlock === undefined) {
        throw new NotFoundError(
          `Found OPEN perpetual positions but no subaccount with address ${address} ` +
          `and number ${subaccountNumber}`,
        );
      }

      const {
        lastUpdatedFundingIndexMap,
        latestFundingIndexMap,
      }: {
        lastUpdatedFundingIndexMap: FundingIndexMap,
        latestFundingIndexMap: FundingIndexMap,
      } = await getFundingIndexMaps(subaccount, latestBlock);
      updatedPerpetualPositions = getPerpetualPositionsWithUpdatedFunding(
        updatedPerpetualPositions,
        latestFundingIndexMap,
        lastUpdatedFundingIndexMap,
      );
    }

    const perpetualMarketsMap: PerpetualMarketsMap = perpetualMarketRefresher
      .getPerpetualMarketsMap();

    const marketIdToMarket: MarketsMap = _.keyBy(
      markets,
      MarketColumns.id,
    );

    return {
      positions: updatedPerpetualPositions.map((position: PerpetualPositionWithFunding) => {
        return perpetualPositionToResponseObject(
          position,
          perpetualMarketsMap,
          marketIdToMarket,
          subaccountNumber,
        );
      }),
    };
  }

  @Get('/parentSubaccountNumber')
  async listPositionsForParentSubaccount(
    @Query() address: string,
      @Query() parentSubaccountNumber: number,
      @Query() status?: PerpetualPositionStatus[],
      @Query() limit?: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
  ): Promise<PerpetualPositionResponse> {
    // Get subaccountIds for all child subaccounts of the parent subaccount
    // Create a record of subaccountId to subaccount number
    const childIdtoSubaccountNumber: Record<string, number> = {};
    getChildSubaccountNums(parentSubaccountNumber).forEach(
      (subaccountNum: number) => {
        childIdtoSubaccountNumber[SubaccountTable.uuid(address, subaccountNum)] = subaccountNum;
      },
    );
    const childSubaccountIds: string[] = Object.keys(childIdtoSubaccountNumber);

    const [
      positions,
      markets,
    ]: [
      PerpetualPositionFromDatabase[],
      MarketFromDatabase[],
    ] = await Promise.all([
      PerpetualPositionTable.findAll(
        {
          subaccountId: childSubaccountIds,
          status,
          limit,
          createdBeforeOrAtHeight: createdBeforeOrAtHeight
            ? createdBeforeOrAtHeight.toString()
            : undefined,
          createdBeforeOrAt,
        },
        [QueryableField.LIMIT],
      ),
      MarketTable.findAll(
        {},
        [],
      ),
    ]);

    const openPositionsExist: boolean = positions.some(
      (position: PerpetualPositionFromDatabase) => position.status === PerpetualPositionStatus.OPEN,
    );

    let updatedPerpetualPositions:
    PerpetualPositionWithFunding[] = initializePerpetualPositionsWithFunding(positions);

    // Only update the funding for open positions
    if (openPositionsExist) {
      const subaccountIds: string[] = updatedPerpetualPositions.map(
        (position: PerpetualPositionFromDatabase) => position.subaccountId,
      );
      const [
        subaccounts,
        latestBlock,
      ]: [
        SubaccountFromDatabase[],
        BlockFromDatabase | undefined,
      ] = await Promise.all([
        SubaccountTable.findAll(
          {
            id: subaccountIds,
          },
          [],
        ),
        BlockTable.getLatest(),
      ]);

      if (subaccounts.length === 0 || latestBlock === undefined) {
        throw new NotFoundError(
          `Found OPEN perpetual positions but no subaccounts with address ${address} ` +
            `and parent Subaccount number ${parentSubaccountNumber}`,
        );
      }

      const perpetualPositionsBySubaccount:
      { [subaccountId: string]: PerpetualPositionWithFunding[] } = _.groupBy(
        updatedPerpetualPositions,
        'subaccountId',
      );

      // For each subaccount, update all perpetual positions with the latest funding and
      // store the updated positions in updatedPerpetualPositions.
      const updatedPerpetualPositionsPromises:
      Promise<PerpetualPositionWithFunding[]>[] = subaccounts.map(
        (subaccount: SubaccountFromDatabase) => adjustPerpetualPositionsWithUpdatedFunding(
          perpetualPositionsBySubaccount[subaccount.id],
          subaccount,
          latestBlock,
        ),
      );
      updatedPerpetualPositions = _.flatten(await Promise.all(updatedPerpetualPositionsPromises));
    }

    const perpetualMarketsMap: PerpetualMarketsMap = perpetualMarketRefresher
      .getPerpetualMarketsMap();

    const marketIdToMarket: MarketsMap = _.keyBy(
      markets,
      MarketColumns.id,
    );

    return {
      positions: updatedPerpetualPositions.map((position: PerpetualPositionWithFunding) => {
        return perpetualPositionToResponseObject(
          position,
          perpetualMarketsMap,
          marketIdToMarket,
          childIdtoSubaccountNumber[position.subaccountId],
        );
      }),
    };
  }
}

async function adjustPerpetualPositionsWithUpdatedFunding(
  perpetualPositions: PerpetualPositionWithFunding[],
  subaccount: SubaccountFromDatabase,
  latestBlock: BlockFromDatabase,
): Promise<PerpetualPositionWithFunding[]> {
  const {
    lastUpdatedFundingIndexMap,
    latestFundingIndexMap,
  }: {
    lastUpdatedFundingIndexMap: FundingIndexMap,
    latestFundingIndexMap: FundingIndexMap,
  } = await getFundingIndexMaps(subaccount, latestBlock);
  return getPerpetualPositionsWithUpdatedFunding(
    perpetualPositions,
    latestFundingIndexMap,
    lastUpdatedFundingIndexMap,
  );
}

router.get(
  '/',
  rateLimiterMiddleware(defaultRateLimiter),
  perpetualPositionsCacheControlMiddleware,
  ...CheckSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtSchema,
  ...checkSchema({
    status: {
      in: ['query'],
      optional: true,
      customSanitizer: {
        options: sanitizeArray,
      },
      custom: {
        options: (inputArray) => validateArray(inputArray, Object.values(PerpetualPositionStatus)),
        errorMessage: 'status must be a valid Position Status (OPEN, etc)',
      },
    },
  }),
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
      status,
      limit,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
    }: PerpetualPositionRequest = matchedData(req) as PerpetualPositionRequest;

    // The schema checks allow subaccountNumber to be a string, but we know it's a number here.
    const subaccountNum: number = +subaccountNumber;

    try {
      const controller: PerpetualPositionsController = new PerpetualPositionsController();
      const response: PerpetualPositionResponse = await controller.listPositions(
        address,
        subaccountNum,
        status,
        limit,
        createdBeforeOrAtHeight,
        createdBeforeOrAt,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'PerpetualPositionsController GET /',
        'Perpetual positions error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_perpetual_positions.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/parentSubaccountNumber',
  rateLimiterMiddleware(defaultRateLimiter),
  perpetualPositionsCacheControlMiddleware,
  ...CheckParentSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtSchema,
  ...checkSchema({
    status: {
      in: ['query'],
      optional: true,
      customSanitizer: {
        options: sanitizeArray,
      },
      custom: {
        options: (inputArray) => validateArray(inputArray, Object.values(PerpetualPositionStatus)),
        errorMessage: 'status must be a valid Position Status (OPEN, etc)',
      },
    },
  }),
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      parentSubaccountNumber,
      status,
      limit,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
    }: ParentSubaccountPerpetualPositionRequest = matchedData(
      req,
    ) as ParentSubaccountPerpetualPositionRequest;

    // The schema checks allow subaccountNumber to be a string, but we know it's a number here.
    const parentSubaccountNum: number = +parentSubaccountNumber;

    try {
      const controller: PerpetualPositionsController = new PerpetualPositionsController();
      const response: PerpetualPositionResponse = await controller.listPositionsForParentSubaccount(
        address,
        parentSubaccountNum,
        status,
        limit,
        createdBeforeOrAtHeight,
        createdBeforeOrAt,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'PerpetualPositionsController GET /parentSubaccountNumber',
        'Perpetual positions error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_perpetual_positions_parent_subaccount.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
