import { stats } from '@dydxprotocol-indexer/base';
import {
  Ordering,
  YieldParamsFromDatabase,
  YieldParamsTable,
  YieldParamsColumns,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';
import _ from 'lodash';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceCheck } from '../../../lib/compliance-check';
import { NotFoundError } from '../../../lib/errors';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { rejectRestrictedCountries } from '../../../lib/restrict-countries';
import {
  CheckLimitAndYieldParamsSchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import {
  yieldParamsToResponseObject,
} from '../../../request-helpers/request-transformer';
import {
  YieldParamsResponse,
  YieldParamsResponseObject,
  YieldParamsRequest,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'yield-params-controller';

@Route('yieldParams')
class YieldParamsController extends Controller {
  @Get('/')
  async getYieldParams(
      @Query() createdBeforeOrAtHeight?: string,
  ): Promise<YieldParamsResponse> {

    // [YBCP-30]: Add cache for yield params
    const query = createdBeforeOrAtHeight !== undefined 
      ? { createdBeforeOrAtHeight: createdBeforeOrAtHeight } 
      : {};
    const allYieldParams: YieldParamsFromDatabase[] | undefined = await YieldParamsTable.findAll(
        query, 
        [], {
        orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
    });

    if (allYieldParams === undefined) {
      throw new NotFoundError(
        `No yield params found before or at ${createdBeforeOrAtHeight}`,
      );
    }

    if (allYieldParams.length === 0) {
      return { allYieldParams: [] };
    }

    const resultParams: YieldParamsResponse = {
      allYieldParams: allYieldParams.map((yieldParams: YieldParamsFromDatabase) => {
        return yieldParamsToResponseObject(yieldParams);
      }),
    }

    return resultParams;
  }


  @Get('/latestYieldParams')
  async getLatestYieldParams(): Promise<YieldParamsResponse> {
    // [YBCP-30]: Add cache for yield params
    const yieldParams: YieldParamsFromDatabase | undefined = await YieldParamsTable.getLatest()

    if (yieldParams === undefined) {
      throw new NotFoundError(
        `No lates yield params found`,
      );
    }

    return {
        allYieldParams: [yieldParamsToResponseObject(yieldParams)],
    };
  }
}

router.get(
  '/',
  rejectRestrictedCountries,
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckLimitAndYieldParamsSchema,
  handleValidationErrors,
  complianceCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const matchedDataObject = matchedData(req);
    const yieldParamsGetRequest: YieldParamsRequest = {
      createdBeforeOrAtHeight: matchedDataObject.createdAtOrBeforeHeight,
    }

    try {
        const controllers: YieldParamsController = new YieldParamsController();
        const response: YieldParamsResponse = await controllers.getYieldParams(
          yieldParamsGetRequest.createdBeforeOrAtHeight,
        );
        return res.send(response);
    } catch (error) {
        return handleControllerError(
            'YieldParamsController GET /',
            'YieldParams error',
            error,
            req,
            res
        );
    } finally {
        stats.timing(
            `${config.SERVICE_NAME}.${controllerName}.get_yield_params.timing`,
            Date.now() - start,
        );
    }
  },
);

router.get(
  '/latestYieldParams',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckLimitAndYieldParamsSchema,
  handleValidationErrors,
  complianceCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    matchedData(req);

    try {
        const controller: YieldParamsController = new YieldParamsController();
        const response: YieldParamsResponse = await controller.getLatestYieldParams();
        return res.send(response);
    } catch (error) {
        return handleControllerError(
            'YieldParamsController GET /latestYieldParams',
            'YieldParams error',
            error,
            req,
            res,
        );
    } finally {
        stats.timing(
            `${config.SERVICE_NAME}.${controllerName}.get_latest_yield_params.timing`,
            Date.now() - start,
        );
    }
  },
);

export default router;
