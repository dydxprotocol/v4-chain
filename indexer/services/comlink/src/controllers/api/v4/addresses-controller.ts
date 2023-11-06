import { stats } from '@dydxprotocol-indexer/base';
import {
  AssetPositionFromDatabase,
  BlockTable,
  BlockFromDatabase,
  PerpetualPositionFromDatabase,
  PerpetualPositionStatus,
  PerpetualPositionTable,
  QueryableField,
  SubaccountFromDatabase,
  SubaccountTable,
  AssetPositionTable,
  AssetTable,
  AssetFromDatabase,
  AssetColumns,
  MarketTable,
  MarketFromDatabase,
  MarketsMap,
  MarketColumns,
  PerpetualMarketsMap,
  perpetualMarketRefresher,
  Options,
  FundingIndexUpdatesTable,
  FundingIndexMap,
} from '@dydxprotocol-indexer/postgres';
import Big from 'big.js';
import express from 'express';
import {
  matchedData,
  checkSchema,
} from 'express-validator';
import _ from 'lodash';
import {
  Route, Get, Path, Controller,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceCheck } from '../../../lib/compliance-check';
import { NotFoundError } from '../../../lib/errors';
import {
  adjustUSDCAssetPosition,
  calculateEquityAndFreeCollateral,
  filterAssetPositions,
  filterPositionsByLatestEventIdPerPerpetual,
  getFundingIndexMaps,
  getTotalUnsettledFunding,
  handleControllerError,
  getPerpetualPositionsWithUpdatedFunding,
  initializePerpetualPositionsWithFunding,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { rejectRestrictedCountries } from '../../../lib/restrict-countries';
import { CheckSubaccountSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import {
  assetPositionToResponseObject,
  perpetualPositionToResponseObject,
  subaccountToResponseObject,
} from '../../../request-helpers/request-transformer';
import {
  AddressRequest,
  PerpetualPositionsMap,
  PerpetualPositionResponseObject,
  SubaccountRequest,
  SubaccountResponseObject,
  AssetById,
  AssetPositionResponseObject,
  AssetPositionsMap,
  PerpetualPositionWithFunding,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'addresses-controller';

@Route('addresses')
class AddressesController extends Controller {
  @Get('/:address')
  public async getAddress(
    @Path() address: string,
  ): Promise<SubaccountResponseObject[]> {
    // TODO(IND-189): Use a transaction across all the DB queries
    const [subaccounts, latestBlock]:
    [SubaccountFromDatabase[], BlockFromDatabase | undefined] = await Promise.all([
      SubaccountTable.findAll(
        {
          address,
        },
        [],
      ),
      BlockTable.getLatest(),
    ]);

    if (subaccounts.length === 0 || latestBlock === undefined) {
      throw new NotFoundError(`No subaccounts found for address ${address}`);
    }

    const latestFundingIndexMap: FundingIndexMap = await FundingIndexUpdatesTable
      .findFundingIndexMap(
        latestBlock.blockHeight,
      );

    const subaccountResponses: SubaccountResponseObject[] = await Promise.all(subaccounts.map(
      async (subaccount: SubaccountFromDatabase): Promise<SubaccountResponseObject> => {
        const [
          perpetualPositions,
          assetPositions,
          assets,
          markets,
          lastUpdatedFundingIndexMap,
        ] = await Promise.all([
          getOpenPerpetualPositionsForSubaccount(
            subaccount.id,
          ),
          getAssetPositionsForSubaccount(
            subaccount.id,
          ),
          AssetTable.findAll(
            {},
            [],
          ),
          MarketTable.findAll(
            {},
            [],
          ),
          FundingIndexUpdatesTable.findFundingIndexMap(
            subaccount.updatedAtHeight,
          ),
        ]);
        const unsettledFunding: Big = getTotalUnsettledFunding(
          perpetualPositions,
          latestFundingIndexMap,
          lastUpdatedFundingIndexMap,
        );

        const updatedPerpetualPositions:
        PerpetualPositionWithFunding[] = getPerpetualPositionsWithUpdatedFunding(
          initializePerpetualPositionsWithFunding(perpetualPositions),
          latestFundingIndexMap,
          lastUpdatedFundingIndexMap,
        );

        return getSubaccountResponse(
          subaccount,
          updatedPerpetualPositions,
          assetPositions,
          assets,
          markets,
          unsettledFunding,
        );
      },
    ));

    return subaccountResponses;
  }

  @Get('/:address/subaccountNumber/:subaccountNumber')
  public async getSubaccount(
    @Path() address: string,
      @Path() subaccountNumber: number,
  ): Promise<SubaccountResponseObject> {
    // TODO(IND-189): Use a transaction across all the DB queries
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);
    const [
      subaccount,
      perpetualPositions,
      assetPositions,
      assets,
      markets,
      latestBlock,
    ]: [
      SubaccountFromDatabase | undefined,
      PerpetualPositionFromDatabase[],
      AssetPositionFromDatabase[],
      AssetFromDatabase[],
      MarketFromDatabase[],
      BlockFromDatabase | undefined,
    ] = await Promise.all([
      SubaccountTable.findById(
        subaccountId,
      ),
      getOpenPerpetualPositionsForSubaccount(
        subaccountId,
      ),
      getAssetPositionsForSubaccount(
        subaccountId,
      ),
      AssetTable.findAll(
        {},
        [],
      ),
      MarketTable.findAll(
        {},
        [],
      ),
      BlockTable.getLatest(),
    ]);

    if (subaccount === undefined || latestBlock === undefined) {
      throw new NotFoundError(
        `No subaccount found with address ${address} and subaccountNumber ${subaccountNumber}`,
      );
    }

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

    const updatedPerpetualPositions:
    PerpetualPositionWithFunding[] = getPerpetualPositionsWithUpdatedFunding(
      initializePerpetualPositionsWithFunding(perpetualPositions),
      latestFundingIndexMap,
      lastUpdatedFundingIndexMap,
    );

    const subaccountResponse: SubaccountResponseObject = await getSubaccountResponse(
      subaccount,
      updatedPerpetualPositions,
      assetPositions,
      assets,
      markets,
      unsettledFunding,
    );
    return subaccountResponse;
  }
}

router.get(
  '/:address',
  rejectRestrictedCountries,
  rateLimiterMiddleware(getReqRateLimiter),
  ...checkSchema({
    address: {
      in: ['params'],
      isString: true,
    },
  }),
  handleValidationErrors,
  complianceCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const {
      address,
    }: {
      address: string,
    } = matchedData(req) as AddressRequest;

    try {
      const controller: AddressesController = new AddressesController();
      const subaccountResponse: SubaccountResponseObject[] = await controller.getAddress(
        address,
      );

      return res.send({
        subaccounts: subaccountResponse,
      });
    } catch (error) {
      return handleControllerError(
        'AddressesController GET /:address',
        'Addresses error',
        error,
        req,
        res,
      );
    }
  },
);

router.get(
  '/:address/subaccountNumber/:subaccountNumber',
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
    }: {
      address: string,
      subaccountNumber: number,
    } = matchedData(req) as SubaccountRequest;

    try {
      const controller: AddressesController = new AddressesController();
      const subaccountResponse: SubaccountResponseObject = await controller.getSubaccount(
        address,
        subaccountNumber,
      );

      return res.send({
        subaccount: subaccountResponse,
      });
    } catch (error) {
      return handleControllerError(
        'AddressesController GET /:address/subaccountNumber/:subaccountNumber',
        'Addresses subaccount error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_addresses.timing`,
        Date.now() - start,
      );
    }
  },
);

/**
 * Gets subaccount response objects given the subaccount, perpetual positions and perpetual markets
 * @param subaccount Subaccount to get response for, from the database
 * @param positions List of perpetual positions held by the subaccount, from the database
 * @param markets List of perpetual markets, from the database
 * @param assetPositions List of asset positions held by the subaccount, from the database
 * @param assets List of assets from the database
 * @param unsettledFunding Total unsettled funding across all open perpetual positions for the
 *                         subaccount
 * @returns Response object for the subaccount
 */
async function getSubaccountResponse(
  subaccount: SubaccountFromDatabase,
  perpetualPositions: PerpetualPositionWithFunding[],
  assetPositions: AssetPositionFromDatabase[],
  assets: AssetFromDatabase[],
  markets: MarketFromDatabase[],
  unsettledFunding: Big,
): Promise<SubaccountResponseObject> {
  const perpetualMarketsMap: PerpetualMarketsMap = perpetualMarketRefresher
    .getPerpetualMarketsMap();
  const marketIdToMarket: MarketsMap = _.keyBy(
    markets,
    MarketColumns.id,
  );

  const filteredPerpetualPositions: PerpetualPositionWithFunding[
  ] = await filterPositionsByLatestEventIdPerPerpetual(perpetualPositions);

  const perpetualPositionResponses:
  PerpetualPositionResponseObject[] = filteredPerpetualPositions.map(
    (perpetualPosition: PerpetualPositionWithFunding): PerpetualPositionResponseObject => {
      return perpetualPositionToResponseObject(
        perpetualPosition,
        perpetualMarketsMap,
        marketIdToMarket,
      );
    },
  );

  const perpetualPositionsMap: PerpetualPositionsMap = _.keyBy(
    perpetualPositionResponses,
    'market',
  );

  const assetIdToAsset: AssetById = _.keyBy(
    assets,
    AssetColumns.id,
  );

  const sortedAssetPositions:
  AssetPositionFromDatabase[] = filterAssetPositions(assetPositions);

  const assetPositionResponses: AssetPositionResponseObject[] = sortedAssetPositions.map(
    (assetPosition: AssetPositionFromDatabase): AssetPositionResponseObject => {
      return assetPositionToResponseObject(
        assetPosition,
        assetIdToAsset,
      );
    },
  );

  const assetPositionsMap: AssetPositionsMap = _.keyBy(
    assetPositionResponses,
    'symbol',
  );
  const {
    assetPositionsMap: adjustedAssetPositionsMap,
    adjustedUSDCAssetPositionSize,
  }: {
    assetPositionsMap: AssetPositionsMap,
    adjustedUSDCAssetPositionSize: string,
  } = adjustUSDCAssetPosition(assetPositionsMap, unsettledFunding);

  const {
    equity,
    freeCollateral,
  }: {
    equity: string,
    freeCollateral: string,
  } = calculateEquityAndFreeCollateral(
    filteredPerpetualPositions,
    perpetualMarketsMap,
    marketIdToMarket,
    adjustedUSDCAssetPositionSize,
  );

  return subaccountToResponseObject({
    subaccount,
    equity,
    freeCollateral,
    openPerpetualPositions: perpetualPositionsMap,
    assetPositions: adjustedAssetPositionsMap,
  });
}

// eslint-disable-next-line  @typescript-eslint/require-await
async function getOpenPerpetualPositionsForSubaccount(
  subaccountId: string,
  options: Options = {},
): Promise<PerpetualPositionFromDatabase[]> {
  // Don't await the promise, since that will start the DB query. Knex (the database library used)
  // will only start executing a DB query once `then` is called, and not when the promise is
  // instantiated.
  return PerpetualPositionTable.findAll(
    {
      subaccountId: [subaccountId],
      status: [PerpetualPositionStatus.OPEN],
    },
    [QueryableField.SUBACCOUNT_ID, QueryableField.STATUS],
    options,
  );
}

// eslint-disable-next-line  @typescript-eslint/require-await
async function getAssetPositionsForSubaccount(
  subaccountId: string,
  options: Options = {},
): Promise<AssetPositionFromDatabase[]> {
  // Don't await the promise, since that will start the DB query. Knex (the database library used)
  // will only start executing a DB query once `then` is called, and not when the promise is
  // instantiated.
  return AssetPositionTable.findAll(
    {
      subaccountId: [subaccountId],
    },
    [QueryableField.SUBACCOUNT_ID],
    options,
  );
}

export default router;
