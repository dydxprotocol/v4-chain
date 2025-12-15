import {
  logger,
  stats,
  TooManyRequestsError,
} from '@dydxprotocol-indexer/base';
import {
  GeoOriginHeaders,
  isRestrictedCountryHeaders,
  isWhitelistedAddress,
  INDEXER_GEOBLOCKED_PAYLOAD,
} from '@dydxprotocol-indexer/compliance';
import {
  ComplianceReason,
  ComplianceStatus,
  ComplianceStatusFromDatabase,
  ComplianceStatusTable,
  WalletFromDatabase,
  WalletTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';
import { DateTime } from 'luxon';
import {
  Controller, Get, Path, Route,
} from 'tsoa';

import { defaultRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceProvider } from '../../../helpers/compliance/compliance-clients';
import {
  ComplianceAction, validateSignature, validateSignatureKeplr,
} from '../../../helpers/compliance/compliance-utils';
import { DYDX_ADDRESS_PREFIX } from '../../../lib/constants';
import { create4xxResponse, handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { getIpAddr } from '../../../lib/utils';
import { CheckAddressSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import {
  BlockedCode,
  ComplianceRequest,
  ComplianceV2Response,
  SetComplianceStatusRequest,
} from '../../../types';
import { ComplianceControllerHelper } from './compliance-controller';

const router: express.Router = express.Router();
const controllerName: string = 'compliance-v2-controller';

@Route('compliance')
class ComplianceV2Controller extends Controller {
  private ipAddress: string;

  constructor(ipAddress: string) {
    super();
    this.ipAddress = ipAddress;
  }

  @Get('/screen/:address')
  async screen(
    @Path() address: string,
  ): Promise<ComplianceV2Response> {
    const controller: ComplianceControllerHelper = new ComplianceControllerHelper(this.ipAddress);
    const {
      restricted,
    }: {
      restricted: boolean,
    } = await controller.screen(address);
    if (!address.startsWith(DYDX_ADDRESS_PREFIX)) {
      if (restricted) {
        return {
          status: ComplianceStatus.BLOCKED,
          reason: ComplianceReason.COMPLIANCE_PROVIDER,
        };
      } else {
        return {
          status: ComplianceStatus.COMPLIANT,
        };
      }
    } else {
      const complianceStatus: ComplianceStatusFromDatabase[] = await
      ComplianceStatusTable.findAll(
        { address: [address] },
        [],
      );
      if (restricted) {
        let complianceStatusFromDatabase: ComplianceStatusFromDatabase | undefined;
        const updatedAt: string = DateTime.utc().toISO();
        if (complianceStatus.length === 0) {
          complianceStatusFromDatabase = await ComplianceStatusTable.upsert({
            address,
            status: ComplianceStatus.BLOCKED,
            reason: ComplianceReason.COMPLIANCE_PROVIDER,
            updatedAt,
          });
        } else if (
          complianceStatus[0].status !== ComplianceStatus.CLOSE_ONLY &&
          complianceStatus[0].status !== ComplianceStatus.BLOCKED
        ) {
          complianceStatusFromDatabase = await ComplianceStatusTable.update({
            address,
            status: ComplianceStatus.CLOSE_ONLY,
            reason: ComplianceReason.COMPLIANCE_PROVIDER,
            updatedAt,
          });
        } else {
          complianceStatusFromDatabase = complianceStatus[0];
        }
        return {
          status: complianceStatusFromDatabase!.status,
          reason: complianceStatusFromDatabase!.reason,
          updatedAt: complianceStatusFromDatabase!.updatedAt,
        };
      } else {
        if (complianceStatus.length === 0) {
          return {
            status: ComplianceStatus.COMPLIANT,
          };
        } else {
          return {
            status: complianceStatus[0].status,
            reason: complianceStatus[0].reason,
            updatedAt: complianceStatus[0].updatedAt,
          };
        }
      }
    }
  }
}

router.get(
  '/screen/:address',
  rateLimiterMiddleware(defaultRateLimiter),
  ...CheckAddressSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();

    const {
      address,
    }: {
      address: string,
    } = matchedData(req) as ComplianceRequest;
    if (isWhitelistedAddress(address)) {
      return res.send({
        status: ComplianceStatus.COMPLIANT,
      });
    }

    try {
      // Rate limiter middleware ensures the ip address can be found from the request
      const ipAddress: string = getIpAddr(req)!;

      const controller: ComplianceV2Controller = new ComplianceV2Controller(ipAddress);
      const response: ComplianceV2Response = await controller.screen(address);

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
        'ComplianceV2Controller GET /screen/:address',
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

router.post(
  '/geoblock',
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();

    const {
      address,
      message,
      currentStatus,
      action,
      signedMessage,
      pubkey,
      timestamp,
    }: {
      address: string,
      message: string,
      currentStatus?: string,
      action: ComplianceAction,
      signedMessage: string, // base64 encoded
      pubkey: string, // base64 encoded
      timestamp: number,  // UNIX timestamp in seconds
    } = req.body;

    try {
      const failedValidationResponse = await validateSignature(
        res, action, address, timestamp, message, signedMessage, pubkey, currentStatus,
      );
      if (failedValidationResponse) {
        return failedValidationResponse;
      }
      return await checkCompliance(req, res, address, action, false);
    } catch (error) {
      return handleError(error, 'geoblock', message, req, res);
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.geo_block.timing`,
        Date.now() - start,
      );
    }
  },
);

router.post(
  '/geoblock-keplr',
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();

    const {
      address,
      message,
      action,
      signedMessage,
      pubkey,
    }: {
      address: string,
      message: string,
      action: ComplianceAction,
      signedMessage: string, // base64 encoded
      pubkey: string, // base64 encoded
    } = req.body;

    try {
      const failedValidationResponse = validateSignatureKeplr(
        res, address, message, signedMessage, pubkey,
      );
      if (failedValidationResponse) {
        return failedValidationResponse;
      }
      return await checkCompliance(req, res, address, action, true);
    } catch (error) {
      return handleError(error, 'geoblock-keplr', message, req, res);
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.geo_block_keplr.timing`,
        Date.now() - start,
      );
    }
  },
);

async function checkCompliance(
  req: express.Request,
  res: express.Response,
  address: string,
  action: ComplianceAction,
  forKeplr: boolean,
): Promise<express.Response> {
  if (isWhitelistedAddress(address)) {
    return res.send({
      status: ComplianceStatus.COMPLIANT,
      updatedAt: DateTime.utc().toISO(),
    });
  }

  if (isRestrictedCountryHeaders(req.headers as GeoOriginHeaders)) {
    return create4xxResponse(
      res,
      INDEXER_GEOBLOCKED_PAYLOAD,
      403,
      { code: BlockedCode.GEOBLOCKED },
    );
  }

  const [
    complianceStatus,
    wallet,
  ]: [
    ComplianceStatusFromDatabase[],
    WalletFromDatabase | undefined,
  ] = await Promise.all([
    ComplianceStatusTable.findAll(
      { address: [address] },
      [],
    ),
    WalletTable.findById(address),
  ]);

  const updatedAt: string = DateTime.utc().toISO();
  const complianceStatusFromDatabase:
  ComplianceStatusFromDatabase | undefined = await upsertComplianceStatus(
    req,
    action,
    address,
    wallet,
    complianceStatus,
    updatedAt,
  );
  if (complianceStatus.length === 0 ||
    complianceStatus[0] !== complianceStatusFromDatabase) {
    if (complianceStatusFromDatabase !== undefined &&
      complianceStatusFromDatabase.status !== ComplianceStatus.COMPLIANT
    ) {
      stats.increment(
        `${config.SERVICE_NAME}.${controllerName}.geo_block${forKeplr ? '_keplr' : ''}.compliance_status_changed.count`,
        {
          newStatus: complianceStatusFromDatabase!.status,
        },
      );
    }
  }

  const response = {
    status: complianceStatusFromDatabase!.status,
    reason: complianceStatusFromDatabase!.reason,
    updatedAt: complianceStatusFromDatabase!.updatedAt,
  };

  return res.send(response);
}

function handleError(
  error: Error, endpointName: string, message: string, req: express.Request, res: express.Response,
): express.Response {
  logger.error({
    at: `ComplianceV2Controller POST /${endpointName}`,
    message,
    error,
    params: JSON.stringify(req.params),
    query: JSON.stringify(req.query),
    body: JSON.stringify(req.body),
  });
  return create4xxResponse(
    res,
    error.message,
  );
}

/**
 * If the address doesn't exist in the compliance table:
 * - if the request is from a restricted country
 *  - the request must be blocked upstream
 * - else if the request is from a non-restricted country:
 *  - set the status to COMPLIANT
 *  - return compliant status
 * If the address does exist in the compliance table:
 * - return the existing status
 */
// eslint-disable-next-line @typescript-eslint/require-await
async function upsertComplianceStatus(
  req: express.Request,
  action: ComplianceAction,
  address: string,
  wallet: WalletFromDatabase | undefined,
  complianceStatus: ComplianceStatusFromDatabase[],
  updatedAt: string,
): Promise<ComplianceStatusFromDatabase | undefined> {

  if (complianceStatus.length === 0) {
    return ComplianceStatusTable.upsert({
      address,
      status: ComplianceStatus.COMPLIANT,
      updatedAt,
    });
  }

  return complianceStatus[0];
}

if (config.EXPOSE_SET_COMPLIANCE_ENDPOINT) {
  router.post(
    '/setStatus',
    handleValidationErrors,
    ExportResponseCodeStats({ controllerName }),
    async (req: express.Request, res: express.Response) => {
      const start: number = Date.now();

      const {
        address,
        status,
        reason,
      }: {
        address: string,
        status: ComplianceStatus,
        reason?: ComplianceReason,
      } = req.body as SetComplianceStatusRequest;

      try {
        if (!address.startsWith(DYDX_ADDRESS_PREFIX)) {
          return create4xxResponse(
            res,
            `Address ${address} is not a dydx address`,
          );
        }
        const complianceStatus: ComplianceStatusFromDatabase = await ComplianceStatusTable.upsert({
          address,
          status,
          reason,
          updatedAt: DateTime.utc().toISO(),
        });
        const response: ComplianceV2Response = {
          status: complianceStatus.status,
          reason: complianceStatus.reason,
        };

        return res.send(response);
      } catch (error) {
        return handleControllerError(
          'ComplianceV2Controller POST /setStatus',
          'Compliance error',
          error,
          req,
          res,
        );
      } finally {
        stats.timing(
          `${config.SERVICE_NAME}.${controllerName}.set_compliance.timing`,
          Date.now() - start,
        );
      }
    },
  );
}

export default router;
