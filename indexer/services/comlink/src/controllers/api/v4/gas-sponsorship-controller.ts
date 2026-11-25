import { stats } from '@dydxprotocol-indexer/base';
import express from 'express';
import { matchedData } from 'express-validator';
import {
  Body, Controller, Post, Route,
} from 'tsoa';

import { gasSponsorshipRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceAndGeoCheck } from '../../../lib/compliance-and-geo-check';
import {
  GasSponsorshipError,
  GasSponsorshipErrorCode,
  GasSponsorshipGrant,
  getGasSponsorshipService,
} from '../../../lib/gas-sponsorship';
import { create4xxResponse, handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { CheckGasSponsorshipSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { GasSponsorshipRequestBody } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'gas-sponsorship-controller';

function respondWithSponsorshipError(
  error: GasSponsorshipError,
  res: express.Response,
): express.Response {
  // Not only 4xx: the helper builds the error body this API uses for any status.
  return create4xxResponse(res, error.message, error.httpStatus, { code: error.code });
}

/**
 * Answers before anything else runs while sponsorship is switched off, so that a switched-off
 * endpoint does no work at all.
 */
function requireGasSponsorshipEnabled(
  _req: express.Request,
  res: express.Response,
  next: express.NextFunction,
): express.Response | void {
  if (!config.GAS_SPONSORSHIP_ENABLED) {
    return respondWithSponsorshipError(
      new GasSponsorshipError(
        GasSponsorshipErrorCode.DISABLED,
        'Gas sponsorship is not available',
      ),
      res,
    );
  }
  return next();
}

@Route('gasSponsorship')
class GasSponsorshipController extends Controller {
  /**
   * Requests a fee allowance on Injective for an address that holds USDC to deposit to dYdX but
   * no INJ to pay for gas.
   *
   * The response names the granter. To use the allowance, the address sends its own transfer
   * with that granter set as the fee granter, and a fee of at most gasLimit * gasPrice.
   *
   * The allowance can be used once its status is CONFIRMED. If the status is PENDING, repeat the
   * request: it returns the same allowance, and never grants a second one.
   */
  @Post('/')
  async requestSponsorship(
    @Body() body: GasSponsorshipRequestBody,
  ): Promise<GasSponsorshipGrant> {
    return getGasSponsorshipService().sponsor({
      granteeAddress: body.injectiveAddress,
      dydxAddress: body.address,
    });
  }
}

router.post(
  '/',
  requireGasSponsorshipEnabled,
  rateLimiterMiddleware(gasSponsorshipRateLimiter, undefined, undefined, true),
  ...CheckGasSponsorshipSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    try {
      const body: GasSponsorshipRequestBody = matchedData(req) as GasSponsorshipRequestBody;
      const controller: GasSponsorshipController = new GasSponsorshipController();
      const grant: GasSponsorshipGrant = await controller.requestSponsorship(body);
      return res.send(grant);
    } catch (error) {
      if (error instanceof GasSponsorshipError) {
        stats.increment(
          `${config.SERVICE_NAME}.${controllerName}.post_gas_sponsorship.refused`,
          { code: error.code },
        );
        return respondWithSponsorshipError(error, res);
      }
      return handleControllerError(
        'GasSponsorshipController POST /',
        'Gas sponsorship error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.post_gas_sponsorship.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
