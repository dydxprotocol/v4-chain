import {
  cacheControlMiddleware,
  logger,
  noCacheControlMiddleware,
  NodeEnv,
  stats,
} from '@dydxprotocol-indexer/base';
import {
  createNotification,
  NotificationDynamicFieldKey,
  NotificationType,
  sendFirebaseMessage,
} from '@dydxprotocol-indexer/notifications';
import {
  AssetFromDatabase,
  AssetPositionFromDatabase,
  AssetPositionTable,
  AssetTable,
  BlockFromDatabase,
  BlockTable,
  FirebaseNotificationTokenTable,
  FundingIndexMap,
  FundingIndexUpdatesTable,
  MarketFromDatabase,
  MarketTable,
  Options,
  perpetualMarketRefresher,
  PerpetualPositionFromDatabase,
  PerpetualPositionStatus,
  PerpetualPositionTable,
  QueryableField,
  SubaccountFromDatabase,
  SubaccountTable,
  WalletFromDatabase,
  WalletTable,
} from '@dydxprotocol-indexer/postgres';
import Big from 'big.js';
import express from 'express';
import {
  matchedData,
} from 'express-validator';
import {
  Body,
  Controller,
  Get, Path,
  Post,
  Route,
} from 'tsoa';

import { defaultRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { AccountVerificationRequiredAction, validateSignature, validateSignatureKeplr } from '../../../helpers/compliance/compliance-utils';
import { complianceAndGeoCheck } from '../../../lib/compliance-and-geo-check';
import { DatabaseError, NotFoundError } from '../../../lib/errors';
import {
  getChildSubaccountIds,
  getFundingIndexMaps,
  getSubaccountResponse,
  handleControllerError,
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
  AddressResponse,
  ParentSubaccountRequest,
  ParentSubaccountResponse,
  RegisterTokenRequest,
  SubaccountRequest,
  SubaccountResponseObject,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'addresses-controller';
const addressesCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_ADDRESSES,
);

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
      @Body() body: { token: string, language: string },
  ): Promise<void> {
    const { token, language } = body;
    const wallet = await WalletTable.findById(address);
    if (!wallet) {
      throw new NotFoundError(`No wallet found with address: ${address}`);
    }
    try {
      // Register the new token
      await FirebaseNotificationTokenTable.registerToken(
        token,
        wallet.address,
        language,
      );
    } catch (error) {
      throw new DatabaseError(`Error registering token: ${error}`);
    }
  }

  @Post('/:address/testNotification')
  public async testNotification(
    @Path() address: string,
  ): Promise<void> {
    try {
      const wallet = await WalletTable.findById(address);
      if (!wallet) {
        throw new NotFoundError(`No wallet found for address: ${address}`);
      }
      const allTokens = await FirebaseNotificationTokenTable.findAll(
        { address: wallet.address }, [],
      );
      if (allTokens.length === 0) {
        throw new NotFoundError(`No tokens found for address: ${address}`);
      }

      const notification = createNotification(NotificationType.ORDER_FILLED, {
        [NotificationDynamicFieldKey.MARKET]: 'BTC/USD',
        [NotificationDynamicFieldKey.AMOUNT]: '100',
        [NotificationDynamicFieldKey.AVERAGE_PRICE]: '1000',
      });
      await sendFirebaseMessage(allTokens, notification);
    } catch (error) {
      logger.error({
        at: 'addresses-controller#testNotification',
        message: error.message,
        error,
      });
    }
  }
}

router.get(
  '/:address',
  rateLimiterMiddleware(defaultRateLimiter),
  addressesCacheControlMiddleware,
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
  rateLimiterMiddleware(defaultRateLimiter),
  addressesCacheControlMiddleware,
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
  rateLimiterMiddleware(defaultRateLimiter),
  addressesCacheControlMiddleware,
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
  noCacheControlMiddleware,
  CheckAddressSchema,
  RegisterTokenValidationSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address, token, language = 'en', timestamp, message, signedMessage, pubKey, walletIsKeplr,
    } = matchedData(req) as RegisterTokenRequest;

    try {
      const failedValidationResponse = walletIsKeplr
        ? validateSignatureKeplr(
          res, address, message, signedMessage, pubKey,
        )
        : await validateSignature(
          res,
          AccountVerificationRequiredAction.REGISTER_TOKEN,
          address,
          timestamp,
          message,
          signedMessage,
          pubKey,
          '',
        );
      if (failedValidationResponse) {
        return failedValidationResponse;
      }

      const controller: AddressesController = new AddressesController();
      await controller.registerToken(address, { token, language });
      return res.status(200).send({});
    } catch (error) {
      return handleControllerError(
        'AddressesController POST /:address/registerToken',
        'Addresses error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.post_registerToken.timing`,
        Date.now() - start,
      );
    }
  },
);

router.post(
  '/:address/testNotification',
  rateLimiterMiddleware(defaultRateLimiter),
  noCacheControlMiddleware,
  ...CheckAddressSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    // This endpoint should only be avaliable in testnet / staging
    if (config.NODE_ENV === NodeEnv.PRODUCTION) {
      return res.status(404).send();
    }

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
