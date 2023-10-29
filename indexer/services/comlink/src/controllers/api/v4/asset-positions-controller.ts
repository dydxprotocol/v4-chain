import { stats } from '@dydxprotocol-indexer/base';
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

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceCheck } from '../../../lib/compliance-check';
import {
  adjustUSDCAssetPosition,
  filterAssetPositions,
  getFundingIndexMaps,
  getTotalUnsettledFunding,
  handleControllerError,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { rejectRestrictedCountries } from '../../../lib/restrict-countries';
import {
  CheckSubaccountSchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { assetPositionToResponseObject } from '../../../request-helpers/request-transformer';
import {
  AssetById,
  AssetPositionRequest,
  AssetPositionResponse,
  AssetPositionResponseObject,
  AssetPositionsMap,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'asset-positions-controller';

@Route('assetPositions')
class AddressesController extends Controller {
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
      BlockFromDatabase | undefined,
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

    const sortedAssetPositions:
    AssetPositionFromDatabase[] = filterAssetPositions(assetPositions);

    const idToAsset: AssetById = _.keyBy(
      assets,
      AssetColumns.id,
    );

    let assetPositionsMap: AssetPositionsMap = _.chain(sortedAssetPositions)
      .map(
        (position: AssetPositionFromDatabase) => assetPositionToResponseObject(position, idToAsset),
      ).keyBy(
        (positionResponse: AssetPositionResponseObject) => positionResponse.symbol,
      ).value();

    // If a subaccount, latest block, and perpetual positions exist, calculate the unsettled funding
    // for positions and adjust the returned USDC position
    if (subaccount !== undefined && latestBlock !== undefined && perpetualPositions.length > 0) {
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

    return {
      positions: Object.values(assetPositionsMap),
    };
  }
}

router.get(
  '/',
  rejectRestrictedCountries,
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckSubaccountSchema,
  handleValidationErrors,
  complianceCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
    }: AssetPositionRequest = matchedData(req) as AssetPositionRequest;

    try {
      const controller: AddressesController = new AddressesController();
      const response: AssetPositionResponse = await controller.getAssetPositions(
        address,
        subaccountNumber,
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

export default router;
