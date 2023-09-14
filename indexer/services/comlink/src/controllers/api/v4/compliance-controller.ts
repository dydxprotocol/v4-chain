import { stats } from '@dydxprotocol-indexer/base';
import { ComplianceClientResponse } from '@dydxprotocol-indexer/compliance';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { placeHolderProvider } from '../../../helpers/compliance/compliance-clients';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { ComplianceRequest, ComplianceResponse } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'compliance-controller';

@Route('screen')
class ComplianceController extends Controller {
  @Get('/')
  async screen(
    @Query() address: string,
  ): Promise<ComplianceResponse> {
    // TODO(IND-372): Add logic to either use cached data or query provider
    // TODO(IND-369): Use Ellptic client
    const response:
    ComplianceClientResponse = await placeHolderProvider.client.getComplianceResponse(
      address,
    );

    return {
      restricted: response.blocked,
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
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();

    const {
      address,
    }: {
      address: string,
    } = matchedData(req) as ComplianceRequest;

    try {
      const controller: ComplianceController = new ComplianceController();
      const response: ComplianceResponse = await controller.screen(address);

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
