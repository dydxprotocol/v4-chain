import {
  ExtendedSecp256k1Signature, Secp256k1, ripemd160, sha256,
} from '@cosmjs/crypto';
import { toBech32 } from '@cosmjs/encoding';
import { logger, stats, TooManyRequestsError } from '@dydxprotocol-indexer/base';
import { CountryHeaders, isRestrictedCountryHeaders } from '@dydxprotocol-indexer/compliance';
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
import _ from 'lodash';
import { DateTime } from 'luxon';
import {
  Controller, Get, Path, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceProvider } from '../../../helpers/compliance/compliance-clients';
import { getGeoComplianceReason } from '../../../helpers/compliance/compliance-utils';
import { DYDX_ADDRESS_PREFIX, GEOBLOCK_REQUEST_TTL_SECONDS } from '../../../lib/constants';
import { create4xxResponse, handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { getIpAddr } from '../../../lib/utils';
import { CheckAddressSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import {
  ComplianceRequest, ComplianceV2Response, SetComplianceStatusRequest,
} from '../../../types';
import { ComplianceControllerHelper } from './compliance-controller';

const router: express.Router = express.Router();
const controllerName: string = 'compliance-v2-controller';

export enum ComplianceAction {
  CONNECT = 'CONNECT',
  VALID_SURVEY = 'VALID_SURVEY',
  INVALID_SURVEY = 'INVALID_SURVEY',
}

const COMPLIANCE_PROGRESSION: Partial<Record<ComplianceStatus, ComplianceStatus>> = {
  [ComplianceStatus.COMPLIANT]: ComplianceStatus.FIRST_STRIKE_CLOSE_ONLY,
  [ComplianceStatus.FIRST_STRIKE]: ComplianceStatus.CLOSE_ONLY,
};

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
        const updatedAt: string = DateTime.utc().toISO();
        if (complianceStatus.length === 0) {
          complianceStatusFromDatabase = await ComplianceStatusTable.upsert({
            address,
            status: ComplianceStatus.BLOCKED,
            reason: ComplianceReason.COMPLIANCE_PROVIDER,
            updatedAt,
          });
        } else {
          complianceStatusFromDatabase = await ComplianceStatusTable.update({
            address,
            status: ComplianceStatus.CLOSE_ONLY,
            reason: ComplianceReason.COMPLIANCE_PROVIDER,
            updatedAt,
          });
        }
        return {
          status: complianceStatusFromDatabase!.status,
          reason: complianceStatusFromDatabase!.reason,
          updatedAt,
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
      signedMessage: string, // base64 encoded
      pubkey: string, // base64 encoded
      timestamp: number,  // UNIX timestamp in seconds
    } = req.body;

    try {
      if (!address.startsWith(DYDX_ADDRESS_PREFIX)) {
        return create4xxResponse(
          res,
          `Address ${address} is not a valid dYdX V4 address`,
        );
      }

      const pubkeyArray: Uint8Array = new Uint8Array(Buffer.from(pubkey, 'base64'));
      if (address !== generateAddress(pubkeyArray)) {
        return create4xxResponse(
          res,
          `Address ${address} does not correspond to the pubkey provided ${pubkey}`,
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
      const messageToSign: string = `${message}:${action}"${currentStatus || ''}:${timestamp}`;
      const messageHash: Uint8Array = sha256(Buffer.from(messageToSign));
      const signedMessageArray: Uint8Array = new Uint8Array(Buffer.from(signedMessage, 'base64'));
      const signature: ExtendedSecp256k1Signature = ExtendedSecp256k1Signature
        .fromFixedLength(signedMessageArray);

      // Verify the signature
      const isValidSignature: boolean = await Secp256k1.verifySignature(
        signature,
        messageHash,
        pubkeyArray,
      );
      if (!isValidSignature) {
        return create4xxResponse(
          res,
          'Signature verification failed',
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
      ComplianceStatusFromDatabase | undefined = await upsertComplicanceStatus(
        req,
        action,
        address,
        wallet,
        complianceStatus,
        updatedAt,
      );

      const response = {
        status: complianceStatusFromDatabase!.status,
        reason: complianceStatusFromDatabase!.reason,
        updatedAt,
      };

      return res.send(response);
    } catch (error) {
      logger.error({
        at: 'ComplianceV2Controller POST /geoblock',
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
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.geo_block.timing`,
        Date.now() - start,
      );
    }
  },
);

function generateAddress(pubkeyArray: Uint8Array): string {
  return toBech32('dydx', ripemd160(sha256(pubkeyArray)));
}

/**
 * If the address doesn't exist in the compliance table:
 * - if the request is from a restricted country:
 *  - if the action is CONNECT and no wallet, set the status to BLOCKED
 *  - if the action is CONNECT and wallet exists, set the status to FIRST_STRIKE_CLOSE_ONLY
 * - else if the request is from a non-restricted country:
 *  - set the status to COMPLIANT
 *
 * if the address is COMPLIANT:
 * - the ONLY action should be CONNECT. VALID_SURVEY/INVALID_SURVEY are no-ops.
 * - if the request is from a restricted country:
 *  - set the status to FIRST_STRIKE_CLOSE_ONLY
 *
 * if the address is FIRST_STRIKE_CLOSE_ONLY:
 * - the ONLY actions should be VALID_SURVEY/INVALID_SURVEY/CONNECT. CONNECT
 * are no-ops.
 * - if the action is VALID_SURVEY:
 *   - set the status to FIRST_STRIKE
 * - if the action is INVALID_SURVEY:
 *   - set the status to CLOSE_ONLY
 *
 * if the address is FIRST_STRIKE:
 * - the ONLY action should be CONNECT. VALID_SURVEY/INVALID_SURVEY are no-ops.
 * - if the request is from a restricted country:
 *  - set the status to CLOSE_ONLY
 */
async function upsertComplicanceStatus(
  req: express.Request,
  action: ComplianceAction,
  address: string,
  wallet: WalletFromDatabase | undefined,
  complianceStatus: ComplianceStatusFromDatabase[],
  updatedAt: string,
): Promise<ComplianceStatusFromDatabase | undefined> {
  if (complianceStatus.length === 0) {
    if (!isRestrictedCountryHeaders(req.headers as CountryHeaders)) {
      return ComplianceStatusTable.upsert({
        address,
        status: ComplianceStatus.COMPLIANT,
        updatedAt,
      });
    }

    // If address is restricted and is not onboarded then block
    if (_.isUndefined(wallet)) {
      return ComplianceStatusTable.upsert({
        address,
        status: ComplianceStatus.BLOCKED,
        reason: getGeoComplianceReason(req.headers as CountryHeaders)!,
        updatedAt,
      });
    }

    return ComplianceStatusTable.upsert({
      address,
      status: ComplianceStatus.FIRST_STRIKE_CLOSE_ONLY,
      reason: getGeoComplianceReason(req.headers as CountryHeaders)!,
      updatedAt,
    });
  }

  if (
    complianceStatus[0].status === ComplianceStatus.FIRST_STRIKE ||
    complianceStatus[0].status === ComplianceStatus.COMPLIANT
  ) {
    if (
      isRestrictedCountryHeaders(req.headers as CountryHeaders) &&
      action === ComplianceAction.CONNECT
    ) {
      return ComplianceStatusTable.update({
        address,
        status: COMPLIANCE_PROGRESSION[complianceStatus[0].status],
        reason: getGeoComplianceReason(req.headers as CountryHeaders)!,
        updatedAt,
      });
    }
  } else if (complianceStatus[0].status === ComplianceStatus.FIRST_STRIKE_CLOSE_ONLY) {
    if (action === ComplianceAction.VALID_SURVEY) {
      return ComplianceStatusTable.update({
        address,
        status: ComplianceStatus.FIRST_STRIKE,
        updatedAt,
      });
    } else if (action === ComplianceAction.INVALID_SURVEY) {
      return ComplianceStatusTable.update({
        address,
        status: ComplianceStatus.CLOSE_ONLY,
        updatedAt,
      });
    }
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
