import { logger, stats, cacheControlMiddleware } from '@dydxprotocol-indexer/base';
import {
  AddressUsername,
  WalletTable,
  AffiliateInfoTable,
  AffiliateReferredUsersTable,
  SubaccountTable,
  SubaccountUsernamesTable,
  AffiliateInfoFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import {
  Body,
  Controller, Get, Query, Route,
  Post,
} from 'tsoa';

import { defaultRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { AccountVerificationRequiredAction, validateSignature, validateSignatureKeplr } from '../../../helpers/compliance/compliance-utils';
import { InvalidParamError, NotFoundError, UnexpectedServerError } from '../../../lib/errors';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { UpdateReferralCodeSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import {
  AffiliateAddressRequest,
  AffiliateMetadataRequest,
  AffiliateMetadataResponse,
  AffiliateAddressResponse,
  AffiliateSnapshotResponse,
  AffiliateSnapshotResponseObject,
  AffiliateSnapshotRequest,
  AffiliateTotalVolumeResponse,
  AffiliateTotalVolumeRequest,
  CreateReferralCodeResponse,
  CreateReferralCodeRequest,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'affiliates-controller';
const affiliatesCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_AFFILIATES,
);
const affiliatesMetadataCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_AFFILIATES_METADATA,
);

// TODO(OTE-731): replace api stubs with real logic
@Route('affiliates')
class AffiliatesController extends Controller {
  @Get('/metadata')
  async getMetadata(
    @Query() address: string,
  ): Promise<AffiliateMetadataResponse> {
    const [walletRow, referredUserRows, subaccountZeroRows] = await Promise.all([
      WalletTable.findById(address),
      AffiliateReferredUsersTable.findByAffiliateAddress(address),
      SubaccountTable.findAll(
        {
          address,
          subaccountNumber: 0,
        },
        [],
      ),
    ]);

    // Check that the address exists
    if (!walletRow) {
      throw new NotFoundError(`Wallet with address ${address} not found`);
    }

    // Check if the address is an affiliate (has referred users)
    const isVolumeEligible = Number(walletRow.totalVolume) >= config.VOLUME_ELIGIBILITY_THRESHOLD;
    const isAffiliate = referredUserRows !== undefined ? referredUserRows.length > 0 : false;

    // No need to check subaccountRows.length > 1 as subaccountNumber is unique for an address
    if (subaccountZeroRows.length === 0) {
      // error logging will be performed by handleInternalServerError
      throw new UnexpectedServerError(`Subaccount 0 not found for address ${address}`);
    } else if (subaccountZeroRows.length > 1) {
      logger.error({
        at: 'affiliates-controller#snapshot',
        message: `More than 1 username exist for address: ${address}`,
        subaccountZeroRows,
      });
    }
    const subaccountId = subaccountZeroRows[0].id;

    // Get subaccount0 username, which is the referral code
    const usernameRows = await SubaccountUsernamesTable.findAll(
      {
        subaccountId: [subaccountId],
      },
      [],
    );
    // No need to check usernameRows.length > 1 as subAccountId is unique (foreign key constraint)
    // This error can happen if a user calls this endpoint before subaccount-username-generator
    // has generated the username
    if (usernameRows.length === 0) {
      stats.increment(`${config.SERVICE_NAME}.${controllerName}.get_metadata.subaccount_username_not_found`);
      throw new UnexpectedServerError(`Username not found for subaccount ${subaccountId}`);
    }
    const referralCode = usernameRows[0].username;

    return {
      referralCode,
      isVolumeEligible,
      isAffiliate,
    };
  }

  @Get('/address')
  async getAddress(
    @Query() referralCode: string,
  ): Promise<AffiliateAddressResponse> {
    const usernameRow = await SubaccountUsernamesTable.findByUsername(referralCode);
    if (!usernameRow) {
      throw new NotFoundError(`Referral code ${referralCode} does not exist`);
    }
    const subAccountId = usernameRow.subaccountId;

    const subaccountRow = await SubaccountTable.findById(subAccountId);
    // subaccountRow should never be undefined because of foreign key constraint between subaccounts
    // and subaccount_usernames tables
    if (!subaccountRow) {
      throw new UnexpectedServerError(`Subaccount ${subAccountId} not found`);
    }
    const address = subaccountRow.address;

    return {
      address,
    };
  }

  @Post('/referralCode')
  async updateCode(
    @Body() body: {
      address: string,
      newCode: string,
    },
  ): Promise<CreateReferralCodeResponse> {
    const {
      address,
      newCode,
    }: {
      address: string,
      newCode: string,
    } = body;

    // Check if the referral code already exists.
    // There is a unique constraint but doing this allows us to have a better error message.
    const existingUsernameRow = await SubaccountUsernamesTable.findByUsername(
      newCode,
    );
    if (existingUsernameRow) {
      throw new InvalidParamError('Referral code already exists');
    }

    const subAccount = await SubaccountTable.findAll(
      {
        address,
        subaccountNumber: 0,
      },
      [],
    );
    // There is a code-level restriction, but it is possible to
    // have more than one subaccount for an address
    // It is also possible for there to be no username, if the task to create it has not run yet
    if (subAccount.length !== 1) {
      throw new InvalidParamError(
        'Referral code update not available yet - please try again later',
      );
    }

    const subaccountId = subAccount[0].id;

    try {
      // there are assumptions here that
      // 1. There is only one entry per subaccountId
      // 2. There is already an entry for the subaccountId
      await SubaccountUsernamesTable.update({
        username: newCode,
        subaccountId,
      });
    } catch (error) {
      throw new UnexpectedServerError('Failed to update referral code - please try again later');
    }

    return {
      referralCode: newCode,
    };
  }

  @Get('/snapshot')
  async getSnapshot(
    @Query() addressFilter?: string[],
      @Query() offset?: number,
      @Query() limit?: number,
      @Query() sortByAffiliateEarning?: boolean,
  ): Promise<AffiliateSnapshotResponse> {
    const finalAddressFilter: string[] = addressFilter ?? [];
    const finalOffset: number = offset ?? 0;
    const finalLimit: number = limit ?? 1000;
    const finalsortByAffiliateEarning: boolean = sortByAffiliateEarning ?? false;

    const infos: AffiliateInfoFromDatabase[] = await AffiliateInfoTable
      .paginatedFindWithAddressFilter(
        finalAddressFilter,
        finalOffset,
        finalLimit,
        finalsortByAffiliateEarning,
      );

    // Get referral codes
    const addressUsernames:
    AddressUsername[] = await SubaccountUsernamesTable.findByAddress(
      infos.map((info) => info.address),
    );
    const addressUsernameMap: Record<string, string> = {};
    addressUsernames.forEach((addressUsername) => {
      addressUsernameMap[addressUsername.address] = addressUsername.username;
    });
    if (addressUsernames.length !== infos.length) {
      const addressesNotFound: string = infos
        .map((info) => info.address)
        .filter((address) => !(address in addressUsernameMap))
        .join(', ');

      logger.warning({
        at: 'affiliates-controller#snapshot',
        message: `Could not find referral code for the following addresses: ${addressesNotFound}`,
      });
    }

    const affiliateSnapshots: AffiliateSnapshotResponseObject[] = infos.map((info) => ({
      affiliateAddress: info.address,
      affiliateReferralCode:
        info.address in addressUsernameMap ? addressUsernameMap[info.address] : '',
      affiliateEarnings: Number(info.affiliateEarnings),
      affiliateReferredTrades: Number(info.referredMakerTrades) + Number(info.referredTakerTrades),
      affiliateTotalReferredFees: Number(info.totalReferredMakerFees) +
      Number(info.totalReferredTakerFees) +
      Number(info.totalReferredMakerRebates),
      affiliateReferredUsers: Number(info.totalReferredUsers),
      affiliateReferredNetProtocolEarnings: Number(info.totalReferredMakerFees) +
      Number(info.totalReferredTakerFees) +
      Number(info.totalReferredMakerRebates) -
      Number(info.affiliateEarnings),
      affiliateReferredTotalVolume: Number(info.referredTotalVolume),
      affiliateReferredMakerFees: Number(info.totalReferredMakerFees),
      affiliateReferredTakerFees: Number(info.totalReferredTakerFees),
      affiliateReferredMakerRebates: Number(info.totalReferredMakerRebates),
    }));

    const response: AffiliateSnapshotResponse = {
      affiliateList: affiliateSnapshots,
      currentOffset: finalOffset,
      total: affiliateSnapshots.length,
    };

    return response;
  }

  @Get('/total_volume')
  public async getTotalVolume(
    @Query() address: string,
  ): Promise<AffiliateTotalVolumeResponse> {
    // Check that the address exists
    const walletRow = await WalletTable.findById(address);
    if (!walletRow) {
      throw new NotFoundError(`Wallet with address ${address} not found`);
    }

    return {
      totalVolume: Number(walletRow.totalVolume),
    };
  }
}

router.get(
  '/metadata',
  rateLimiterMiddleware(defaultRateLimiter),
  affiliatesMetadataCacheControlMiddleware,
  ...checkSchema({
    address: {
      in: ['query'],
      isString: true,
      errorMessage: 'address must be a valid string',
    },
  }),
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
    }: AffiliateMetadataRequest = matchedData(req) as AffiliateMetadataRequest;

    try {
      const controller: AffiliatesController = new AffiliatesController();
      const response: AffiliateMetadataResponse = await controller.getMetadata(address);
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'AffiliatesController GET /metadata',
        'Affiliates referral code error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_metadata.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/address',
  rateLimiterMiddleware(defaultRateLimiter),
  affiliatesCacheControlMiddleware,
  ...checkSchema({
    referralCode: {
      in: ['query'],
      isString: true,
      errorMessage: 'referralCode must be a valid string',
    },
  }),
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      referralCode,
    }: AffiliateAddressRequest = matchedData(req) as AffiliateAddressRequest;

    try {
      const controller: AffiliatesController = new AffiliatesController();
      const response: AffiliateAddressResponse = await controller.getAddress(referralCode);
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'AffiliatesController GET /address',
        'Affiliates address error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_address.timing`,
        Date.now() - start,
      );
    }
  },
);

router.post(
  '/referralCode',
  ...UpdateReferralCodeSchema(true),
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();

    const {
      address,
      newCode,
      signedMessage,
      pubKey,
      timestamp,
    }: CreateReferralCodeRequest = req.body;

    try {
      const failedValidationResponse = await validateSignature(
        res,
        AccountVerificationRequiredAction.UPDATE_CODE,
        address,
        timestamp,
        newCode,
        signedMessage,
        pubKey,
      );
      if (failedValidationResponse) {
        return failedValidationResponse;
      }

      const controller: AffiliatesController = new AffiliatesController();
      const response: CreateReferralCodeResponse = await controller.updateCode({
        address,
        newCode,
      });
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'AffiliatesController POST /referralCode',
        'Affiliates referral code error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.create_code.timing`,
        Date.now() - start,
      );
    }
  },
);

// Keplr wallet does uses a completely different signature format
// so we need to have a separate endpoint for it
router.post(
  '/referralCode-keplr',
  ...UpdateReferralCodeSchema(false),
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();

    const {
      address,
      newCode,
      signedMessage,
      pubKey,
    }: CreateReferralCodeRequest = req.body;

    try {
      const failedValidationResponse = await validateSignatureKeplr(
        res,
        address,
        newCode,
        signedMessage,
        pubKey,
      );
      if (failedValidationResponse) {
        return failedValidationResponse;
      }

      const controller: AffiliatesController = new AffiliatesController();
      const response: CreateReferralCodeResponse = await controller.updateCode({
        address,
        newCode,
      });
      return res.send(response);

    } catch (error) {
      return handleControllerError(
        'AffiliatesController POST /referralCode-keplr',
        'Affiliates referral code error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.create_code.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/snapshot',
  rateLimiterMiddleware(defaultRateLimiter),
  affiliatesCacheControlMiddleware,
  ...checkSchema({
    addressFilter: {
      in: ['query'],
      optional: true,
      customSanitizer: {
        options: (value) => {
          // Split the comma-separated string into an array
          return typeof value === 'string' ? value.split(',') : value;
        },
      },
      custom: {
        options: (values) => {
          return Array.isArray(values) &&
            values.length > 0 &&
            values.every((val) => typeof val === 'string');
        },
      },
      errorMessage: 'addressFilter must be a non-empy array of comma separated strings',
    },
    offset: {
      in: ['query'],
      optional: true,
      isInt: {
        options: { min: 0 },
      },
      toInt: true,
      errorMessage: 'offset must be a valid integer',
    },
    limit: {
      in: ['query'],
      optional: true,
      isInt: {
        options: { min: 1 },
      },
      toInt: true,
      errorMessage: 'limit must be a valid integer',
    },
    sortByAffiliateEarning: {
      in: ['query'],
      isBoolean: true,
      toBoolean: true,
      optional: true,
      errorMessage: 'sortByAffiliateEarning must be a boolean',
    },
  }),
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      addressFilter,
      offset,
      limit,
      sortByAffiliateEarning,
    }: AffiliateSnapshotRequest = matchedData(req) as AffiliateSnapshotRequest;

    try {
      const controller: AffiliatesController = new AffiliatesController();
      const response: AffiliateSnapshotResponse = await controller.getSnapshot(
        addressFilter,
        offset,
        limit,
        sortByAffiliateEarning,
      );
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'AffiliatesController GET /snapshot',
        'Affiliates snapshot error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_snapshot.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/total_volume',
  rateLimiterMiddleware(defaultRateLimiter),
  affiliatesCacheControlMiddleware,
  ...checkSchema({
    address: {
      in: ['query'],
      isString: true,
      errorMessage: 'address must be a valid string',
    },
  }),
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
    }: AffiliateTotalVolumeRequest = matchedData(req) as AffiliateTotalVolumeRequest;

    try {
      const controller: AffiliatesController = new AffiliatesController();
      const response: AffiliateTotalVolumeResponse = await controller.getTotalVolume(address);
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'AffiliateTotalVolumeResponse GET /total_volume',
        'Affiliate total volume error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_total_volume.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
