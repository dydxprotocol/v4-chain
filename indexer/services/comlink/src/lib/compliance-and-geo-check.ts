import {
  GeoOriginHeaders,
  isRestrictedCountryHeaders,
  INDEXER_GEOBLOCKED_PAYLOAD,
  INDEXER_COMPLIANCE_BLOCKED_PAYLOAD,
  isWhitelistedAddress,
} from '@dydxprotocol-indexer/compliance';
import {
  ComplianceStatus,
  ComplianceStatusFromDatabase,
  ComplianceStatusTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';

import { AddressRequest, BlockedCode } from '../types';
import {
  create4xxResponse,
  handleInternalServerError,
} from './helpers';
import { getIpAddr, isIndexerIp } from './utils';

/**
 * Checks if the address in the request is blocked or not.
 *
 * IF the address is in the compliance_status table and has the status CLOSE_ONLY,
 * return data for the endpoint
 * ELSE IF the address has compliance_status of BLOCKED block access to the endpoint (return 403)
 * ELSE IF the origin country is restricted geography, block access to the endpoint (return 403)
 * ELSE return data for the endpoint
 * NOTE: This middleware must be used after `checkSchema` to ensure `matchData` can get the
 * address parameter from the request.
 */
export async function complianceAndGeoCheck(
  req: express.Request,
  res: express.Response,
  next: express.NextFunction,
) {
  const ipAddr: string | undefined = getIpAddr(req);

  // Don't enforce geo-blocking for internal IPs as they don't go through a proxy
  if (ipAddr !== undefined && isIndexerIp(ipAddr)) {
    return next();
  }

  const { address }: AddressRequest = matchedData(req) as AddressRequest;
  if (isWhitelistedAddress(address)) {
    return next();
  }

  if (address !== undefined) {
    try {
      const updatedStatus: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll(
        { address: [address] },
        [],
        { readReplica: true },
      );
      if (updatedStatus.length > 0) {
        if (updatedStatus[0].status === ComplianceStatus.CLOSE_ONLY ||
          updatedStatus[0].status === ComplianceStatus.FIRST_STRIKE_CLOSE_ONLY
        ) {
          return next();
        } else if (updatedStatus[0].status === ComplianceStatus.BLOCKED) {
          return create4xxResponse(
            res,
            INDEXER_COMPLIANCE_BLOCKED_PAYLOAD,
            403,
            { code: BlockedCode.COMPLIANCE_BLOCKED },
          );
        }
      }
    } catch (error) {
      return handleInternalServerError(
        'complianceAndGeoCheck',
        'complianceAndGeoCheck error',
        error,
        req,
        res,
      );
    }
  }

  if (isRestrictedCountryHeaders(req.headers as GeoOriginHeaders)) {
    return create4xxResponse(
      res,
      INDEXER_GEOBLOCKED_PAYLOAD,
      403,
      { code: BlockedCode.GEOBLOCKED },
    );
  }

  return next();
}
