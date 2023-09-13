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
import { ComplianceRequest, ComplianceResponse } from '../../../types';
import { handleValidationErrors } from '../../../request-helpers/error-handler';

const router: express.Router = express.Router();
const controllerName: string = 'compliance-controller';

@Route('screen')
class ComplianceController extends Controller {
  @Get('/')
  screen(
    @Query() address: string,
  ): ComplianceResponse {
    // TODO(IND-372): Add logic to either use cached data or query provider
    // Dummy logic for front-end testing, returns true if the address ends in a letter between
    // 'a' and 'm'
    if (
      address.charCodeAt(address.length - 1) > 'a'.charCodeAt(0) &&
      address.charCodeAt(address.length - 1) < 'm'.charCodeAt(0)
    ) {
      return {
        restricted: true,
      };
    }

    return {
      restricted: false,
    };
  }
}

router.get(
  '/',
  // TODO(IND-372): Add custom rate-limiter around un-cached requests / global rate-limit
  rateLimiterMiddleware(getReqRateLimiter),
  ...checkSchema({
    address: {
      in: ['query'],
      isString: true,
    },
  }),
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  (req: express.Request, res: express.Response) => {
    const start: number = Date.now();

    const {
      address,
    }: {
      address: string,
    } = matchedData(req) as ComplianceRequest;

    try {
      const controller: ComplianceController = new ComplianceController();
      const response: ComplianceResponse = controller.screen(address);

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'ComplianceController GET /',
        'Compliance error',
        error,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.compliance_screen.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
