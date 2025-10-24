import {
  logger,
  stats,
  TooManyRequestsError,
} from '@dydxprotocol-indexer/base';
import {
  ComplianceClientResponse,
  INDEXER_COMPLIANCE_BLOCKED_PAYLOAD,
  NOT_IN_BLOCKCHAIN_RISK_SCORE,
} from '@dydxprotocol-indexer/compliance';
import { ComplianceDataCreateObject, ComplianceDataFromDatabase, ComplianceTable } from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import _ from 'lodash';
import { DateTime } from 'luxon';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import {
  defaultRateLimiter,
  screenProviderGlobalLimiter,
  screenProviderLimiter,
} from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceProvider } from '../../../helpers/compliance/compliance-clients';
import { create4xxResponse, handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { getIpAddr } from '../../../lib/utils';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { ComplianceRequest, ComplianceResponse } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'compliance-controller';
const UNCACHED_QUERY_POINTS: number = 1;
const GLOBAL_RATE_LIMIT_KEY: string = 'screenQueryProviderGlobal';

@Route('screen')
export class ComplianceControllerHelper extends Controller {
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

    if (complianceData !== undefined) {
      stats.increment(
        `${config.SERVICE_NAME}.${controllerName}.compliance_data_cache_hit`,
        { provider: complianceProvider.provider },
      );
    }

    // Immediately return for blocked addresses, do not refresh
    if (complianceData?.blocked) {
      return {
        restricted: true,
        reason: INDEXER_COMPLIANCE_BLOCKED_PAYLOAD,
      };
    }

    if (complianceData === undefined || DateTime.fromISO(complianceData.updatedAt) < ageThreshold) {
      await checkRateLimit(this.ipAddress);

      if (complianceData === undefined) {
        stats.increment(
          `${config.SERVICE_NAME}.${controllerName}.compliance_data_cache_miss`,
          { provider: complianceProvider.provider },
        );
      } else {
        stats.increment(
          `${config.SERVICE_NAME}.${controllerName}.refresh_compliance_data_cache`,
          { provider: complianceProvider.provider },
        );
      }

      const response:
      ComplianceClientResponse = await complianceProvider.client.getComplianceResponse(
        address,
      );
      // Don't upsert invalid addresses (address causing ellitic error) to compliance table.
      // When the elliptic request fails with 404, getComplianceResponse returns
      // riskScore=NOT_IN_BLOCKCHAIN_RISK_SCORE
      if (response.riskScore === undefined ||
        Number(response.riskScore) === NOT_IN_BLOCKCHAIN_RISK_SCORE) {
        return {
          restricted: false,
          reason: undefined,
        };
      }

      complianceData = await ComplianceTable.upsert({
        ..._.omitBy(response, _.isUndefined) as ComplianceDataCreateObject,
        provider: complianceProvider.provider,
        updatedAt: DateTime.utc().toISO(),
      });
    }

    return {
      restricted: complianceData.blocked,
      reason: complianceData.blocked ? INDEXER_COMPLIANCE_BLOCKED_PAYLOAD : undefined,
    };
  }
}

router.get(
  '/',
  rateLimiterMiddleware(defaultRateLimiter),
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

      const controller: ComplianceControllerHelper = new ComplianceControllerHelper(ipAddress);
      const response: ComplianceResponse = await controller.screen(address);

      return res.send(response);
    } catch (error) {
      if (error instanceof TooManyRequestsError) {
        stats.increment(
          `${config.SERVICE_NAME}.${controllerName}.compliance_screen_rate_limited_attempts`,
          { provider: complianceProvider.provider },
        );
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
        req,
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
