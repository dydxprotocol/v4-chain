import { logger, stats, TooManyRequestsError } from '@dydxprotocol-indexer/base';
import { ComplianceClientResponse } from '@dydxprotocol-indexer/compliance';
import { ComplianceDataFromDatabase, ComplianceTable } from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import { DateTime } from 'luxon';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import {
  getReqRateLimiter,
  screenProviderLimiter,
  screenProviderGlobalLimiter,
} from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceProvider } from '../../../helpers/compliance/compliance-clients';
import { create4xxResponse, handleControllerError } from '../../../lib/helpers';
import { getIpAddr, rateLimiterMiddleware } from '../../../lib/rate-limit';
import { rejectRestrictedCountries } from '../../../lib/restrict-countries';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { ComplianceRequest, ComplianceResponse } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'compliance-controller';
const UNCACHED_QUERY_POINTS: number = 1;
const GLOBAL_RATE_LIMIT_KEY: string = 'screenQueryProviderGlobal';

@Route('screen')
class ComplianceController extends Controller {
  private ipAddress: string;

  constructor(ipAddress: string) {
    super();
    this.ipAddress = ipAddress;
  }

  @Get('/')
  async screen(
    @Query() address: string,
  ): Promise<ComplianceResponse> {
    const ageThreshold: DateTime = DateTime.utc().minus({
      seconds: config.MAX_AGE_SCREENED_ADDRESS_COMPLIANCE_DATA_SECONDS,
    });

    let complianceData:
    ComplianceDataFromDatabase | undefined = await ComplianceTable.findByAddressAndProvider(
      address,
      complianceProvider.provider,
    );

    if (complianceData === undefined || DateTime.fromISO(complianceData.updatedAt) < ageThreshold) {
      await checkRateLimit(this.ipAddress);
      // TODO(IND-369): Use Ellptic client
      const response:
      ComplianceClientResponse = await complianceProvider.client.getComplianceResponse(
        address,
      );
      complianceData = await ComplianceTable.upsert({
        ...response,
        provider: complianceProvider.provider,
        updatedAt: DateTime.utc().toISO(),
      });
    }

    return {
      restricted: complianceData.blocked,
    };
  }
}

router.get(
  '/',
  rejectRestrictedCountries,
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
      // Rate limiter middleware ensures the ip address can be found from the request
      const ipAddress: string = getIpAddr(req)!;

      const controller: ComplianceController = new ComplianceController(ipAddress);
      const response: ComplianceResponse = await controller.screen(address);

      return res.send(response);
    } catch (error) {
      if (error instanceof TooManyRequestsError) {
        return create4xxResponse(
          res,
          'Too many requests',
          429,
        );
      }
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

async function checkRateLimit(ipAddress: string) {
  try {
    await Promise.all([
      screenProviderLimiter.consume(ipAddress, UNCACHED_QUERY_POINTS),
      screenProviderGlobalLimiter.consume(GLOBAL_RATE_LIMIT_KEY, UNCACHED_QUERY_POINTS),
    ]);
  } catch (reject) {
    if (reject instanceof Error) {
      logger.error({
        at: 'rate-limit',
        message: 'redis error when checking rate limit',
        reject,
      });
    } else {
      throw new TooManyRequestsError('Too many requests');
    }
  }
}

export default router;
