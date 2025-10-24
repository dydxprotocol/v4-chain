import { stats, cacheControlMiddleware } from '@dydxprotocol-indexer/base';
import {
  AssetColumns,
  AssetFromDatabase,
  AssetPositionFromDatabase,
  AssetPositionTable,
  AssetTable,
  BlockFromDatabase,
  BlockTable,
  FundingIndexMap,
  PerpetualPositionFromDatabase,
  PerpetualPositionStatus,
  PerpetualPositionTable,
  QueryableField,
  SubaccountFromDatabase,
  SubaccountTable,
} from '@dydxprotocol-indexer/postgres';
import Big from 'big.js';
import express from 'express';
import { matchedData } from 'express-validator';
import _ from 'lodash';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { defaultRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceAndGeoCheck } from '../../../lib/compliance-and-geo-check';
import { NotFoundError } from '../../../lib/errors';
import {
  adjustUSDCAssetPosition,
  filterAssetPositions,
  getChildSubaccountIds,
  getFundingIndexMaps,
  getTotalUnsettledFunding,
  handleControllerError,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { CheckParentSubaccountSchema, CheckSubaccountSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import {
  assetPositionToResponseObject,
} from '../../../request-helpers/request-transformer';
import {
  AssetById,
  AssetPositionRequest,
  AssetPositionResponse,
  AssetPositionResponseObject,
  AssetPositionsMap,
  ParentSubaccountAssetPositionRequest,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'asset-positions-controller';
const assetPositionsCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_ASSET_POSITIONS,
);

@Route('assetPositions')
class AssetPositionsController extends Controller {
  @Get('/')
  async getAssetPositions(
    @Query() address: string,
      @Query() subaccountNumber: number,
  ): Promise<AssetPositionResponse> {
    const subaccountUuid: string = SubaccountTable.uuid(address, subaccountNumber);

    // TODO(IND-189): Use a transaction across all the DB queries
    const [
      subaccount,
      assetPositions,
      perpetualPositions,
      // TODO(DEC-656): Change to a cache in Redis or local instead of querying DB.
      assets,
      latestBlock,
    ]: [
      SubaccountFromDatabase | undefined,
      AssetPositionFromDatabase[],
      PerpetualPositionFromDatabase[],
      AssetFromDatabase[],
      BlockFromDatabase,
    ] = await Promise.all([
      SubaccountTable.findById(
        subaccountUuid,
      ),
      AssetPositionTable.findAll(
        {
          subaccountId: [subaccountUuid],
        },
        [QueryableField.SUBACCOUNT_ID],
      ),
      PerpetualPositionTable.findAll(
        {
          subaccountId: [subaccountUuid],
          status: [PerpetualPositionStatus.OPEN],
        },
        [QueryableField.SUBACCOUNT_ID],
      ),
      AssetTable.findAll(
        {},
        [],
      ),
      BlockTable.getLatest(),
    ]);

    if (subaccount === undefined) {
      throw new NotFoundError(
        `No subaccount found with address ${address} and subaccountNumber ${subaccountNumber}`,
      );
    }

    const sortedAssetPositions:
    AssetPositionFromDatabase[] = filterAssetPositions(assetPositions);

    const idToAsset: AssetById = _.keyBy(
      assets,
      AssetColumns.id,
    );

    const assetPositionsMap: AssetPositionsMap = await adjustAssetPositionsWithFunding(
      sortedAssetPositions,
      perpetualPositions,
      idToAsset,
      subaccount,
      latestBlock,
    );

    return {
      positions: Object.values(assetPositionsMap),
    };
  }

  @Get('/parentSubaccountNumber')
  async getAssetPositionsForParentSubaccount(
    @Query() address: string,
      @Query() parentSubaccountNumber: number,
  ): Promise<AssetPositionResponse> {

    const childSubaccountUuids: string[] = getChildSubaccountIds(address, parentSubaccountNumber);

    // TODO(IND-189): Use a transaction across all the DB queries
    const [
      subaccounts,
      assetPositions,
      perpetualPositions,
      // TODO(DEC-656): Change to a cache in Redis or local instead of querying DB.
      assets,
      latestBlock,
    ]: [
      SubaccountFromDatabase[],
      AssetPositionFromDatabase[],
      PerpetualPositionFromDatabase[],
      AssetFromDatabase[],
      BlockFromDatabase,
    ] = await Promise.all([
      SubaccountTable.findAll(
        {
          id: childSubaccountUuids,
        },
        [QueryableField.ID],
      ),
      AssetPositionTable.findAll(
        {
          subaccountId: childSubaccountUuids,
        },
        [QueryableField.SUBACCOUNT_ID],
      ),
      PerpetualPositionTable.findAll(
        {
          subaccountId: childSubaccountUuids,
          status: [PerpetualPositionStatus.OPEN],
        },
        [QueryableField.SUBACCOUNT_ID],
      ),
      AssetTable.findAll(
        {},
        [],
      ),
      BlockTable.getLatest(),
    ]);

    const sortedAssetPositions:
    AssetPositionFromDatabase[] = filterAssetPositions(assetPositions);

    const idToAsset: AssetById = _.keyBy(
      assets,
      AssetColumns.id,
    );

    const assetPositionsBySubaccount:
    { [subaccountId: string]: AssetPositionFromDatabase[] } = _.groupBy(
      sortedAssetPositions,
      'subaccountId',
    );

    const perpetualPositionsBySubaccount:
    { [subaccountId: string]: PerpetualPositionFromDatabase[] } = _.groupBy(
      perpetualPositions,
      'subaccountId',
    );

    // For each subaccount, adjust the asset positions with the unsettled funding and return the
    // asset positions per subaccount
    const assetPositionsPromises = subaccounts.map(async (subaccount) => {
      const adjustedAssetPositionsMap: AssetPositionsMap = await adjustAssetPositionsWithFunding(
        assetPositionsBySubaccount[subaccount.id] || [],
        perpetualPositionsBySubaccount[subaccount.id] || [],
        idToAsset,
        subaccount,
        latestBlock,
      );
      return Object.values(adjustedAssetPositionsMap);
    });

    const assetPositionsResponse: AssetPositionResponseObject[] = (
      await Promise.all(assetPositionsPromises)
    ).flat();

    return {
      positions: assetPositionsResponse,
    };
  }
}

/**
 * Helper function to adjust the asset positions with the unsettled funding
 * per subaccount
 * @param assetPositions pulled from DB
 * @param perpetualPositions pulled from DB
 * @param idToAsset mapping of assetId to asset
 * @param subaccount subaccount for which the asset positions are being adjusted
 * @param latestBlock
 * @returns AssetPositionsMap
 */
async function adjustAssetPositionsWithFunding(
  assetPositions: AssetPositionFromDatabase[],
  perpetualPositions: PerpetualPositionFromDatabase[],
  idToAsset: AssetById,
  subaccount: SubaccountFromDatabase,
  latestBlock: BlockFromDatabase,
): Promise<AssetPositionsMap> {
  let assetPositionsMap: AssetPositionsMap = _.chain(assetPositions)
    .map(
      (position: AssetPositionFromDatabase) => assetPositionToResponseObject(
        position,
        idToAsset,
        subaccount.subaccountNumber),
    ).keyBy(
      (positionResponse: AssetPositionResponseObject) => positionResponse.symbol,
    ).value();

  // If the latest block, and perpetual positions exist, calculate the unsettled funding
  // for positions and adjust the returned USDC position
  if (perpetualPositions.length > 0) {
    const {
      lastUpdatedFundingIndexMap,
      latestFundingIndexMap,
    }: {
      lastUpdatedFundingIndexMap: FundingIndexMap,
      latestFundingIndexMap: FundingIndexMap,
    } = await getFundingIndexMaps(subaccount, latestBlock);
    const unsettledFunding: Big = getTotalUnsettledFunding(
      perpetualPositions,
      latestFundingIndexMap,
      lastUpdatedFundingIndexMap,
    );

    // Adjust the USDC asset position
    const {
      assetPositionsMap: adjustedAssetPositionsMap,
    }: {
      assetPositionsMap: AssetPositionsMap,
      adjustedUSDCAssetPositionSize: string,
    } = adjustUSDCAssetPosition(assetPositionsMap, unsettledFunding);
    assetPositionsMap = adjustedAssetPositionsMap;
  }

  return assetPositionsMap;
}

router.get(
  '/',
  rateLimiterMiddleware(defaultRateLimiter),
  assetPositionsCacheControlMiddleware,
  ...CheckSubaccountSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
    }: AssetPositionRequest = matchedData(req) as AssetPositionRequest;

    // The schema checks allow subaccountNumber to be a string, but we know it's a number here.
    const subaccountNum : number = +subaccountNumber;

    try {
      const controller: AssetPositionsController = new AssetPositionsController();
      const response: AssetPositionResponse = await controller.getAssetPositions(
        address,
        subaccountNum,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'AssetPositionsController GET /',
        'Asset positions error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_asset_positions.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/parentSubaccountNumber',
  rateLimiterMiddleware(defaultRateLimiter),
  assetPositionsCacheControlMiddleware,
  ...CheckParentSubaccountSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      parentSubaccountNumber,
    }: ParentSubaccountAssetPositionRequest = matchedData(
      req,
    ) as ParentSubaccountAssetPositionRequest;

    // The schema checks allow subaccountNumber to be a string, but we know it's a number here.
    const parentSubaccountNum: number = +parentSubaccountNumber;

    try {
      const controller: AssetPositionsController = new AssetPositionsController();
      const response: AssetPositionResponse = await controller.getAssetPositionsForParentSubaccount(
        address,
        parentSubaccountNum,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'AssetPositionsController GET /parentSubaccountNumber',
        'Asset positions error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_asset_positions_parent_subaccount.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
