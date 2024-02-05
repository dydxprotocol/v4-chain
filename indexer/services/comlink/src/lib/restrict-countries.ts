import {
  CountryHeaders,
  isRestrictedCountryHeaders,
  INDEXER_GEOBLOCKED_PAYLOAD,
  INDEXER_COMPLIANCE_BLOCKED_PAYLOAD,
} from '@dydxprotocol-indexer/compliance';
import {
  ComplianceStatus,
  ComplianceStatusFromDatabase,
  ComplianceStatusTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';

import { AddressRequest, BlockedCode } from '../types';
import { create4xxResponse } from './helpers';
import { getIpAddr, isIndexerIp } from './utils';

/**
 * Return an error code for users that access the API from a restricted country
 */
export async function rejectRestrictedCountries(
  req: express.Request,
  res: express.Response,
  next: express.NextFunction,
) {
  const ipAddr: string | undefined = getIpAddr(req);

  // Don't enforce geo-blocking for internal IPs as they don't go through a proxy
  if (ipAddr !== undefined && isIndexerIp(ipAddr)) {
    return next();
  }

  const {
    address,
  }: {
    address: string,
  } = matchedData(req) as AddressRequest;
  console.log('address', address);
  const updatedStatus: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll(
    { address: [address] },
    [],
    { readReplica: true },
  );
  if (updatedStatus.length > 0) {
    if (updatedStatus[0].status === ComplianceStatus.CLOSE_ONLY) {
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

  if (isRestrictedCountryHeaders(req.headers as CountryHeaders)) {
    return create4xxResponse(
      res,
      INDEXER_GEOBLOCKED_PAYLOAD,
      403,
      { code: BlockedCode.GEOBLOCKED },
    );
  }

  return next();
}
