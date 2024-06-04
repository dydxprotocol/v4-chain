import { logger, stats, TooManyRequestsError } from '@dydxprotocol-indexer/base';
import { ComplianceClientResponse, INDEXER_COMPLIANCE_BLOCKED_PAYLOAD } from '@dydxprotocol-indexer/compliance';
import { ComplianceDataFromDatabase, ComplianceTable } from '@dydxprotocol-indexer/postgres';
import { DateTime } from 'luxon';
import {
  Controller, Get, Query,
} from 'tsoa';

import {
  screenProviderLimiter,
  screenProviderGlobalLimiter,
} from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceProvider } from '../../../helpers/compliance/compliance-clients';
import {ComplianceResponse } from '../../../types';

const controllerName: string = 'compliance-controller';
const UNCACHED_QUERY_POINTS: number = 1;
const GLOBAL_RATE_LIMIT_KEY: string = 'screenQueryProviderGlobal';


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
      complianceData = await ComplianceTable.upsert({
        ...response,
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