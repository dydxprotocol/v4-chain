import { stats } from '@dydxprotocol-indexer/base';
import {
  IsoString,
  Ordering,
  FundingPaymentsFromDatabase,
  SubaccountTable,
  QueryableField,
  FundingPaymentsQueryConfig,
  FundingPaymentsTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceAndGeoCheck } from '../../../lib/compliance-and-geo-check';
import {
  handleControllerError,
  getChildSubaccountNums,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import {
  CheckLimitSchema,
  CheckParentSubaccountSchema,
  CheckSubaccountSchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { fundingPaymentsToResponseObject } from '../../../request-helpers/request-transformer';
import {
  FundingPaymentResponseObject,
  FundingPaymentResponse,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'funding-payments-controller';

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
        ? { orderBy: [['createdAt', Ordering.DESC]] }
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
      subaccountId: Object.keys(childIdtoSubaccountNumber),
      createdOnOrAfter: afterOrAt,
      limit,
      page,
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
        ? { orderBy: [['createdAt', Ordering.DESC]] }
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
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckSubaccountSchema,
  ...CheckLimitSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address, subaccountNumber, limit, ticker, afterOrAt, page,
    } = matchedData(req) as {
      address: string,
      subaccountNumber: number,
      limit?: number,
      ticker?: string,
      afterOrAt?: IsoString,
      page?: number,
    };

    try {
      const controller: FundingPaymentController = new FundingPaymentController();
      const response: FundingPaymentResponse = await controller.getFundingPayments(
        address,
        subaccountNumber,
        limit,
        ticker,
        afterOrAt,
        page,
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
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckParentSubaccountSchema,
  ...CheckLimitSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address, parentSubaccountNumber, limit, afterOrAt, page,
    } = matchedData(req) as {
      address: string,
      parentSubaccountNumber: number,
      limit?: number,
      afterOrAt?: IsoString,
      page?: number,
    };

    const parentSubaccountNum: number = +parentSubaccountNumber;

    try {
      const ctrl: FundingPaymentController = new FundingPaymentController();
      const response: FundingPaymentResponse = await ctrl.getFundingPaymentsForParentSubaccount(
        address,
        parentSubaccountNum,
        limit,
        afterOrAt,
        page,
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
