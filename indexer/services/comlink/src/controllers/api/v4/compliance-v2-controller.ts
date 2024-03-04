import { Secp256k1, sha256, ExtendedSecp256k1Signature } from '@cosmjs/crypto';
import { logger, stats, TooManyRequestsError } from '@dydxprotocol-indexer/base';
import {
  ComplianceReason,
  ComplianceStatus,
  ComplianceStatusFromDatabase,
  ComplianceStatusTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';
import { DateTime } from 'luxon';
import {
  Controller, Get, Path, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceProvider } from '../../../helpers/compliance/compliance-clients';
import { DYDX_ADDRESS_PREFIX, GEOBLOCK_REQUEST_TTL_SECONDS } from '../../../lib/constants';
import { create4xxResponse, handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { getIpAddr } from '../../../lib/utils';
import { CheckAddressSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { ComplianceRequest, ComplianceV2Response, SetComplianceStatusRequest } from '../../../types';
import { ComplianceControllerHelper } from './compliance-controller';

const router: express.Router = express.Router();
const controllerName: string = 'compliance-v2-controller';

export enum ComplianceAction {
  ONBOARD = 'ONBOARD',
  CONNECT = 'CONNECT',
}

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
      if (restricted) {
        const complianceStatus: ComplianceStatusFromDatabase[] = await
        ComplianceStatusTable.findAll(
          { address: [address] },
          [],
        );
        let complianceStatusFromDatabase: ComplianceStatusFromDatabase | undefined;
        if (complianceStatus.length === 0) {
          complianceStatusFromDatabase = await ComplianceStatusTable.upsert({
            address,
            status: ComplianceStatus.BLOCKED,
            reason: ComplianceReason.COMPLIANCE_PROVIDER,
            updatedAt: DateTime.utc().toISO(),
          });
        } else {
          complianceStatusFromDatabase = await ComplianceStatusTable.update({
            address,
            status: ComplianceStatus.CLOSE_ONLY,
            reason: ComplianceReason.COMPLIANCE_PROVIDER,
            updatedAt: DateTime.utc().toISO(),
          });
        }
        return {
          status: complianceStatusFromDatabase!.status,
          reason: complianceStatusFromDatabase!.reason,
        };
      } else {
        return {
          status: ComplianceStatus.COMPLIANT,
        };
      }
    }
  }
}

router.get(
  '/screen/:address',
  rateLimiterMiddleware(getReqRateLimiter),
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
      signedMessage: Uint8Array,
      pubkey: Uint8Array,
      timestamp: number,  // UNIX timestamp in seconds
    } = req.body;

    try {
      if (!address.startsWith(DYDX_ADDRESS_PREFIX)) {
        return create4xxResponse(
          res,
          `Address ${address} is not a valid dYdX V4 address`,
        );
      }

      // Verify the timestamp is within GEOBLOCK_REQUEST_TTL_SECONDS seconds of the current time
      const now = DateTime.now().toSeconds();
      if (Math.abs(now - timestamp) > GEOBLOCK_REQUEST_TTL_SECONDS) {
        return create4xxResponse(
          res,
          `Timestamp is not within the valid range of ${GEOBLOCK_REQUEST_TTL_SECONDS} seconds`,
        );
      }

      // Prepare the message for verification
      const messageToSign: string = `${message}${action}${currentStatus || ''}${timestamp}`;
      const messageHash: Uint8Array = sha256(Buffer.from(messageToSign));
      const signature: ExtendedSecp256k1Signature = ExtendedSecp256k1Signature
        .fromFixedLength(signedMessage);

      // Verify the signature
      const isValidSignature: boolean = await
      Secp256k1.verifySignature(signature, messageHash, pubkey);
      if (!isValidSignature) {
        return create4xxResponse(
          res,
          'Signature verification failed',
        );
      }

      // TODO(OTE-141): Implement logic.
      const response = {
        status: ComplianceStatus.COMPLIANT,
      };

      return res.send(response);
    } catch (error) {
      logger.error({
        at: 'ComplianceV2Controller POST /geoblock',
        message,
        error,
        params: JSON.stringify(req.params),
        query: JSON.stringify(req.query),
      });
      return create4xxResponse(
        res,
        error.message,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.geo_block.timing`,
        Date.now() - start,
      );
    }
  },
);

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
