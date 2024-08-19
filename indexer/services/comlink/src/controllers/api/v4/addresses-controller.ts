import { stats } from '@dydxprotocol-indexer/base';
import {
  createNotification, NotificationType, NotificationDynamicFieldKey, sendFirebaseMessage,
} from '@dydxprotocol-indexer/notifications';
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
  MarketTable,
  MarketFromDatabase,
  Options,
  FundingIndexUpdatesTable,
  FundingIndexMap,
  WalletTable,
  WalletFromDatabase,
<<<<<<< HEAD
  perpetualMarketRefresher,
=======
  TokenTable,
>>>>>>> ce9ce9ab (Register token with postgres from comlink)
} from '@dydxprotocol-indexer/postgres';
import Big from 'big.js';
import express from 'express';
import {
  matchedData,
} from 'express-validator';
import {
  Route, Get, Path, Controller,
  Post,
  Body,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceAndGeoCheck } from '../../../lib/compliance-and-geo-check';
import { BadRequestError, DatabaseError, NotFoundError } from '../../../lib/errors';
import {
  getFundingIndexMaps,
  handleControllerError,
  getChildSubaccountIds,
  getSubaccountResponse,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import {
  CheckAddressSchema,
  CheckParentSubaccountSchema,
  CheckSubaccountSchema,
  RegisterTokenValidationSchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import {
  AddressRequest,
  SubaccountRequest,
  SubaccountResponseObject,
  AddressResponse,
  ParentSubaccountResponse,
  ParentSubaccountRequest,
  RegisterTokenRequest,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'addresses-controller';

@Route('addresses')
class AddressesController extends Controller {
  @Get('/:address')
  public async getAddress(
    @Path() address: string,
  ): Promise<AddressResponse> {
    // TODO(IND-189): Use a transaction across all the DB queries
    const [subaccounts, latestBlock, wallet]: [
      SubaccountFromDatabase[],
      BlockFromDatabase,
      WalletFromDatabase | undefined,
    ] = await Promise.all([
      SubaccountTable.findAll(
        {
          address,
        },
        [],
      ),
      BlockTable.getLatest(),
      WalletTable.findById(address),
    ]);

    if (subaccounts.length === 0) {
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

        return getSubaccountResponse(
          subaccount,
          perpetualPositions,
          assetPositions,
          assets,
          markets,
          perpetualMarketRefresher.getPerpetualMarketsMap(),
          latestBlock.blockHeight,
          latestFundingIndexMap,
          lastUpdatedFundingIndexMap,
        );
      },
    ));

    return {
      subaccounts: subaccountResponses,
      totalTradingRewards: wallet !== undefined ? wallet.totalTradingRewards : '0',
    };
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
      BlockFromDatabase,
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

    if (subaccount === undefined) {
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

    const subaccountResponse: SubaccountResponseObject = getSubaccountResponse(
      subaccount,
      perpetualPositions,
      assetPositions,
      assets,
      markets,
      perpetualMarketRefresher.getPerpetualMarketsMap(),
      latestBlock.blockHeight,
      latestFundingIndexMap,
      lastUpdatedFundingIndexMap,
    );
    return subaccountResponse;
  }

  @Get('/:address/parentSubaccountNumber/:parentSubaccountNumber')
  public async getParentSubaccount(
    @Path() address: string,
      @Path() parentSubaccountNumber: number,
  ): Promise<ParentSubaccountResponse> {

    const childSubaccountIds: string[] = getChildSubaccountIds(address, parentSubaccountNumber);

    // TODO(IND-189): Use a transaction across all the DB queries
    const [subaccounts, latestBlock]: [
      SubaccountFromDatabase[],
      BlockFromDatabase,
    ] = await Promise.all([
      SubaccountTable.findAll(
        {
          id: childSubaccountIds,
          address,
        },
        [],
      ),
      BlockTable.getLatest(),
    ]);

    if (subaccounts.length === 0) {
      throw new NotFoundError(`No subaccounts found for address ${address} and parentSubaccountNumber ${parentSubaccountNumber}`);
    }

    const latestFundingIndexMap: FundingIndexMap = await FundingIndexUpdatesTable
      .findFundingIndexMap(
        latestBlock.blockHeight,
      );

    const [assets, markets]: [AssetFromDatabase[], MarketFromDatabase[]] = await Promise.all([
      AssetTable.findAll(
        {},
        [],
      ),
      MarketTable.findAll(
        {},
        [],
      ),
    ]);
    const subaccountResponses: SubaccountResponseObject[] = await Promise.all(subaccounts.map(
      async (subaccount: SubaccountFromDatabase): Promise<SubaccountResponseObject> => {
        const [
          perpetualPositions,
          assetPositions,
          lastUpdatedFundingIndexMap,
        ] = await Promise.all([
          getOpenPerpetualPositionsForSubaccount(
            subaccount.id,
          ),
          getAssetPositionsForSubaccount(
            subaccount.id,
          ),
          FundingIndexUpdatesTable.findFundingIndexMap(
            subaccount.updatedAtHeight,
          ),
        ]);

        return getSubaccountResponse(
          subaccount,
          perpetualPositions,
          assetPositions,
          assets,
          markets,
          perpetualMarketRefresher.getPerpetualMarketsMap(),
          latestBlock.blockHeight,
          latestFundingIndexMap,
          lastUpdatedFundingIndexMap,
        );
      },
    ));

    return {
      address,
      parentSubaccountNumber,
      equity: subaccountResponses.reduce(
        (acc: Big, subaccount: SubaccountResponseObject): Big => acc.plus(subaccount.equity),
        Big(0),
      ).toString(),
      freeCollateral: subaccountResponses.reduce(
        // eslint-disable-next-line max-len
        (acc: Big, subaccount: SubaccountResponseObject): Big => acc.plus(subaccount.freeCollateral),
        Big(0),
      ).toString(),
      childSubaccounts: subaccountResponses,
    };
  }

  @Post('/:address/registerToken')
  public async registerToken(
    @Path() address: string,
      @Body() body: { token: string },
  ): Promise<void> {
    const { token } = body;
    if (!token) {
      throw new BadRequestError('Invalid Token in request');
    }

    const foundAddress = await WalletTable.findById(address);
    if (!foundAddress) {
      throw new NotFoundError(`No address found with address: ${address}`);
    }

    try {
      await TokenTable.registerToken(
        token,
        address,
      );
    } catch (error) {
      throw new DatabaseError(`Error registering token: ${error}`);
    }
  }

  @Post('/:address/testNotification')
  public async testNotification(
    @Path() address: string,
  ): Promise<void> {
    const wallet = await WalletTable.findById(address);
    if (!wallet) {
      throw new NotFoundError(`No wallet found for address: ${address}`);
    }

    try {
      const notification = createNotification(NotificationType.ORDER_FILLED, {
        [NotificationDynamicFieldKey.MARKET]: 'BTC/USD',
        [NotificationDynamicFieldKey.AMOUNT]: '100',
        [NotificationDynamicFieldKey.AVERAGE_PRICE]: '1000',
      });
      await sendFirebaseMessage(wallet.address, notification);
    } catch (error) {
      throw new Error('Failed to send test notification');
    }
  }
}

router.get(
  '/:address',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckAddressSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
    }: {
      address: string,
    } = matchedData(req) as AddressRequest;

    try {
      const controller: AddressesController = new AddressesController();
      const addressResponse: AddressResponse = await controller.getAddress(
        address,
      );

      return res.send(addressResponse);
    } catch (error) {
      return handleControllerError(
        'AddressesController GET /:address',
        'Addresses error',
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

router.get(
  '/:address/subaccountNumber/:subaccountNumber',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckSubaccountSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
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
        `${config.SERVICE_NAME}.${controllerName}.get_subaccount.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/:address/parentSubaccountNumber/:parentSubaccountNumber',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckParentSubaccountSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      parentSubaccountNumber,
    }: {
      address: string,
      parentSubaccountNumber: number,
    } = matchedData(req) as ParentSubaccountRequest;

    // The schema checks allow subaccountNumber to be a string, but we know it's a number here.
    const parentSubaccountNum = +parentSubaccountNumber;

    try {
      const controller: AddressesController = new AddressesController();
      const subaccountResponse: ParentSubaccountResponse = await controller.getParentSubaccount(
        address,
        parentSubaccountNum,
      );

      return res.send({
        subaccount: subaccountResponse,
      });
    } catch (error) {
      return handleControllerError(
        'AddressesController GET /:address/parentSubaccountNumber/:parentSubaccountNumber',
        'Addresses subaccount error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_parentSubaccount.timing`,
        Date.now() - start,
      );
    }
  },
);

router.post(
  '/:address/registerToken',
  CheckAddressSchema,
  RegisterTokenValidationSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const { address, token } = matchedData(req) as RegisterTokenRequest;

    try {
      const controller: AddressesController = new AddressesController();
      await controller.registerToken(address, { token });
      return res.status(200).send({});
    } catch (error) {
      return handleControllerError(
        'AddressesController POST /:address/registerToken',
        'Addresses error',
        error,
        req,
        res,
      );
    }
  },
);

router.post(
  '/:address/testNotification',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckAddressSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const { address } = matchedData(req) as AddressRequest;

    try {
      const controller: AddressesController = new AddressesController();
      await controller.testNotification(address);
      return res.status(200).send({ message: 'Test notification sent successfully' });
    } catch (error) {
      return handleControllerError(
        'AddressesController POST /:address/testNotification',
        'Test notification error',
        error,
        req,
        res,
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
  latestBlockHeight: string,
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
        subaccount.subaccountNumber,
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
        subaccount.subaccountNumber,
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
    latestBlockHeight,
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
