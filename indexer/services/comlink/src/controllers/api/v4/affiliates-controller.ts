import { stats } from '@dydxprotocol-indexer/base';
import {
  WalletTable,
  AffiliateReferredUsersTable,
  SubaccountTable,
  SubaccountUsernamesTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { NotFoundError, UnexpectedServerError } from '../../../lib/errors';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
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
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'affiliates-controller';

// TODO(OTE-731): replace api stubs with real logic
@Route('affiliates')
class AffiliatesController extends Controller {
  @Get('/metadata')
  async getMetadata(
    @Query() address: string,
  ): Promise<AffiliateMetadataResponse> {
    const [walletRow, referredUserRows, subaccountRows] = await Promise.all([
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
    if (subaccountRows.length === 0) {
      // error logging will be performed by handleInternalServerError
      throw new UnexpectedServerError(`Subaccount 0 not found for address ${address}`);
    }
    const subaccountId = subaccountRows[0].id;

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
    @Query() referralCode: string, // eslint-disable-line @typescript-eslint/no-unused-vars
  ): Promise<AffiliateAddressResponse> {
    // simulate a delay
    await new Promise((resolve) => setTimeout(resolve, 100));
    return {
      address: 'some_address',
    };
  }

  @Get('/snapshot')
  async getSnapshot(
    @Query() offset?: number,
      @Query() limit?: number,
      @Query() sortByReferredFees?: boolean,
  ): Promise<AffiliateSnapshotResponse> {
    const finalOffset = offset ?? 0;
    const finalLimit = limit ?? 1000;
    // eslint-disable-next-line
    const finalSortByReferredFees = sortByReferredFees ?? false;

    // simulate a delay
    await new Promise((resolve) => setTimeout(resolve, 100));

    const snapshot: AffiliateSnapshotResponseObject = {
      affiliateAddress: 'some_address',
      affiliateReferralCode: 'TempCode123',
      affiliateEarnings: 100,
      affiliateReferredTrades: 1000,
      affiliateTotalReferredFees: 100,
      affiliateReferredUsers: 10,
      affiliateReferredNetProtocolEarnings: 1000,
      affiliateReferredTotalVolume: 1000000,
    };

    const affiliateSnapshots: AffiliateSnapshotResponseObject[] = [];
    for (let i = 0; i < finalLimit; i++) {
      affiliateSnapshots.push(snapshot);
    }

    const response: AffiliateSnapshotResponse = {
      affiliateList: affiliateSnapshots,
      total: finalLimit,
      currentOffset: finalOffset,
    };

    return response;
  }

  @Get('/total_volume')
  public async getTotalVolume(
    @Query() address: string, // eslint-disable-line @typescript-eslint/no-unused-vars
  ): Promise<AffiliateTotalVolumeResponse> {
    // simulate a delay
    await new Promise((resolve) => setTimeout(resolve, 100));
    return {
      totalVolume: 111.1,
    };
  }
}

router.get(
  '/metadata',
  rateLimiterMiddleware(getReqRateLimiter),
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
  rateLimiterMiddleware(getReqRateLimiter),
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

router.get(
  '/snapshot',
  rateLimiterMiddleware(getReqRateLimiter),
  ...checkSchema({
    offset: {
      in: ['query'],
      isInt: true,
      toInt: true,
      optional: true,
      errorMessage: 'offset must be a valid integer',
    },
    limit: {
      in: ['query'],
      isInt: true,
      toInt: true,
      optional: true,
      errorMessage: 'limit must be a valid integer',
    },
    sortByReferredFees: {
      in: ['query'],
      isBoolean: true,
      toBoolean: true,
      optional: true,
      errorMessage: 'sortByReferredFees must be a boolean',
    },
  }),
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      offset,
      limit,
      sortByReferredFees,
    }: AffiliateSnapshotRequest = matchedData(req) as AffiliateSnapshotRequest;

    try {
      const controller: AffiliatesController = new AffiliatesController();
      const response: AffiliateSnapshotResponse = await controller.getSnapshot(
        offset,
        limit,
        sortByReferredFees,
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
  rateLimiterMiddleware(getReqRateLimiter),
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
