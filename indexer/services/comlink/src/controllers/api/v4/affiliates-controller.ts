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
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { AffiliateReferralCodeRequest, AffiliateReferralCodeResponse } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'affiliates-controller';

@Route('affiliates')
class AffiliatesController extends Controller {
  @Get('/referral_code')
  async getReferralCode(
    @Query() address: string, // eslint-disable-line @typescript-eslint/no-unused-vars
  ): Promise<AffiliateReferralCodeResponse> {
    // TODO: OTE-731 replace apit stubs with real logic
    // simulate a delay
    await new Promise((resolve) => setTimeout(resolve, 100));
    return {
      referralCode: 'TempCode123',
    };
  }
}

router.get(
  '/referral_code',
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
    }: AffiliateReferralCodeRequest = matchedData(req) as AffiliateReferralCodeRequest;

    try {
      const controller: AffiliatesController = new AffiliatesController();
      const response: AffiliateReferralCodeResponse = await controller.getReferralCode(address);
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'AffiliatesController GET /referral_code',
        'Affiliates referral code error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_referral_code.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
