import { stats } from '@dydxprotocol-indexer/base';
import {
  PerpetualPositionStatus,
  SubaccountTable,
  PerpetualPositionFromDatabase,
  PerpetualPositionTable,
  IsoString,
  PerpetualMarketsMap,
  QueryableField,
  MarketFromDatabase,
  MarketTable,
  MarketsMap,
  MarketColumns,
  perpetualMarketRefresher,
  SubaccountFromDatabase,
  BlockFromDatabase,
  BlockTable,
  FundingIndexMap,
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

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceCheck } from '../../../lib/compliance-check';
import { NotFoundError } from '../../../lib/errors';
import {
  getFundingIndexMaps,
  handleControllerError,
  getPerpetualPositionsWithUpdatedFunding,
  initializePerpetualPositionsWithFunding,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { rejectRestrictedCountries } from '../../../lib/restrict-countries';
import {
  CheckLimitAndCreatedBeforeOrAtSchema,
  CheckSubaccountSchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { perpetualPositionToResponseObject } from '../../../request-helpers/request-transformer';
import { sanitizeArray } from '../../../request-helpers/sanitizers';
import { validateArray } from '../../../request-helpers/validators';
import { PerpetualPositionRequest, PerpetualPositionResponse, PerpetualPositionWithFunding } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'perpetual-positions-controller';

@Route('perpetualPositions')
class PerpetualPositionsController extends Controller {
  @Get('/')
  async listPositions(
    @Query() address: string,
      @Query() subaccountNumber: number,
      @Query() status: PerpetualPositionStatus[],
      @Query() limit: number,
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
        return perpetualPositionToResponseObject(position, perpetualMarketsMap, marketIdToMarket);
      }),
    };
  }
}

router.get(
  '/',
  rejectRestrictedCountries,
  rateLimiterMiddleware(getReqRateLimiter),
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
  complianceCheck,
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

    try {
      const controller: PerpetualPositionsController = new PerpetualPositionsController();
      const response: PerpetualPositionResponse = await controller.listPositions(
        address,
        subaccountNumber,
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

export default router;
