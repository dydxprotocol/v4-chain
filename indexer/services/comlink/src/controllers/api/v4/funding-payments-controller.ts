import { cacheControlMiddleware, stats } from '@dydxprotocol-indexer/base';
import {
  FundingPaymentsFromDatabase,
  FundingPaymentsQueryConfig,
  FundingPaymentsTable,
  IsoString,
  Ordering,
  QueryableField,
  SubaccountTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { fundingRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceAndGeoCheck } from '../../../lib/compliance-and-geo-check';
import {
  getChildSubaccountNums,
  handleControllerError,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import {
  CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  CheckPaginationSchema,
  CheckParentSubaccountSchema,
  CheckSubaccountSchema,
  CheckTickerOptionalQuerySchema,
  CheckZeroPaymentsOptionalParamSchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { fundingPaymentsToResponseObject } from '../../../request-helpers/request-transformer';
import {
  FundingPaymentResponse,
  FundingPaymentResponseObject,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'funding-payments-controller';
const fundingPaymentsCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_FUNDING,
);

@Route('fundingPayments')
export class FundingPaymentController extends Controller {
  @Get('/')
  async getFundingPayments(
    @Query() address: string,
      @Query() subaccountNumber: number,
      @Query() limit?: number,
      @Query() ticker?: string,
      @Query() afterOrAt?: IsoString,
      @Query() page?: number,
      @Query() zeroPayments?: boolean,
  ): Promise<FundingPaymentResponse> {
    const subaccountId: string = SubaccountTable.uuid(
      address,
      subaccountNumber,
    );

    const queryConfig: FundingPaymentsQueryConfig = {
      subaccountId: [subaccountId],
      ticker,
      createdOnOrAfter: afterOrAt,
      limit,
      page,
      zeroPayments,
    };

    const {
      results: fundingPayments,
      limit: pageSize,
      offset,
      total,
    } = await FundingPaymentsTable.findAll(
      queryConfig,
      [QueryableField.LIMIT],
      page !== undefined
        ? { orderBy: [['createdAtHeight', Ordering.DESC]] }
        : undefined,
    );

    return {
      fundingPayments: fundingPayments.map(
        (
          fundingPayment: FundingPaymentsFromDatabase,
        ): FundingPaymentResponseObject => {
          return fundingPaymentsToResponseObject(
            fundingPayment,
            subaccountNumber,
          );
        },
      ),
      pageSize,
      totalResults: total,
      offset,
    };
  }

  @Get('/parentSubaccount')
  // Note: This is expected to be used for FE only, where `parentSubaccount -> childSubaccount`
  // mapping is relevant. API traders should use `fundingPayments/` instead.
  async getFundingPaymentsForParentSubaccount(
    @Query() address: string,
      @Query() parentSubaccountNumber: number,
      @Query() limit?: number,
      @Query() afterOrAt?: IsoString,
      @Query() page?: number,
      @Query() zeroPayments?: boolean,
  ): Promise<FundingPaymentResponse> {
    const childIdtoSubaccountNumber: Record<string, number> = {};
    getChildSubaccountNums(parentSubaccountNumber).forEach(
      (subaccountNum: number) => {
        childIdtoSubaccountNumber[
          SubaccountTable.uuid(address, subaccountNum)
        ] = subaccountNum;
      },
    );

    const queryConfig: FundingPaymentsQueryConfig = {
      parentSubaccount: {
        address,
        subaccountNumber: parentSubaccountNumber,
      },
      createdOnOrAfter: afterOrAt,
      limit,
      page,
      zeroPayments,
    };

    const {
      results: fundingPayments,
      limit: pageSize,
      offset,
      total,
    } = await FundingPaymentsTable.findAll(
      queryConfig,
      [QueryableField.LIMIT],
      page !== undefined
        ? { orderBy: [['createdAtHeight', Ordering.DESC]] }
        : undefined,
    );

    return {
      fundingPayments: fundingPayments.map(
        (
          fundingPayment: FundingPaymentsFromDatabase,
        ): FundingPaymentResponseObject => fundingPaymentsToResponseObject(
          fundingPayment,
          childIdtoSubaccountNumber[fundingPayment.subaccountId],
        ),
      ),
      pageSize,
      totalResults: total,
      offset,
    };
  }
}

router.get(
  '/',
  rateLimiterMiddleware(fundingRateLimiter),
  fundingPaymentsCacheControlMiddleware,
  ...CheckSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  ...CheckPaginationSchema,
  ...CheckTickerOptionalQuerySchema,
  ...CheckZeroPaymentsOptionalParamSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address, subaccountNumber, limit, ticker, createdOnOrAfter, page, showZeroPayments,
    } = matchedData(req) as {
      address: string,
      subaccountNumber: number,
      limit?: number,
      ticker?: string,
      createdOnOrAfter?: IsoString,
      page?: number,
      showZeroPayments?: boolean,
    };

    try {
      const controller: FundingPaymentController = new FundingPaymentController();
      const response: FundingPaymentResponse = await controller.getFundingPayments(
        address,
        subaccountNumber,
        limit,
        ticker,
        createdOnOrAfter,
        page,
        showZeroPayments,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'FundingPaymentsController GET /',
        'Funding payments error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_funding_payments.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/parentSubaccount',
  rateLimiterMiddleware(fundingRateLimiter),
  fundingPaymentsCacheControlMiddleware,
  ...CheckParentSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  ...CheckPaginationSchema,
  ...CheckTickerOptionalQuerySchema,
  ...CheckZeroPaymentsOptionalParamSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address, parentSubaccountNumber, limit, page, createdOnOrAfter, showZeroPayments,
    } = matchedData(req) as {
      address: string,
      parentSubaccountNumber: number,
      limit?: number,
      createdOnOrAfter?: IsoString,
      page?: number,
      showZeroPayments?: boolean,
    };

    const parentSubaccountNum: number = +parentSubaccountNumber;

    try {
      const ctrl: FundingPaymentController = new FundingPaymentController();
      const response: FundingPaymentResponse = await ctrl.getFundingPaymentsForParentSubaccount(
        address,
        parentSubaccountNum,
        limit,
        createdOnOrAfter,
        page,
        showZeroPayments,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'FundingPaymentsController GET /parentSubaccount',
        'Funding payments error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_funding_payments_for_parent_subaccount.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
