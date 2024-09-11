import { stats } from '@dydxprotocol-indexer/base';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import {
  AffiliateAddressRequest,
  AffiliateReferralCodeRequest,
  AffiliateReferralCodeResponse,
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
    @Query() address: string, // eslint-disable-line @typescript-eslint/no-unused-vars
  ): Promise<AffiliateReferralCodeResponse> {
    // simulate a delay
    await new Promise((resolve) => setTimeout(resolve, 100));
    return {
      referralCode: 'TempCode123',
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
    }: AffiliateReferralCodeRequest = matchedData(req) as AffiliateReferralCodeRequest;

    try {
      const controller: AffiliatesController = new AffiliatesController();
      const response: AffiliateReferralCodeResponse = await controller.getMetadata(address);
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
