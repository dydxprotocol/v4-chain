import { INDEXER_COMPLIANCE_BLOCKED_PAYLOAD } from '@dydxprotocol-indexer/compliance';
import { ComplianceDataFromDatabase, ComplianceTable } from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';

import { AddressRequest, BlockedCode } from '../types';
import { create4xxResponse } from './helpers';

/**
 * Checks if the address in the request is blocked or not.
 * Returns 403 if the address is blocked, otherwise if there is no address in the request, or the
 * address does not exist in the database or is not blocked, continue onto the next middleware.
 * NOTE: This middleware must be used after `checkSchema` to ensure `matchData` can get the
 * address parameter from the request.
 */
export async function complianceCheck(
  req: express.Request,
  res: express.Response,
  next: express.NextFunction,
) {
  // Check for the address parameter in either query params or path params
  const { address }: AddressRequest = matchedData(req) as AddressRequest;
  if (address === undefined) {
    return next();
  }

  // Search for any compliance data indicating the address is blocked
  const complianceData: ComplianceDataFromDatabase[] = await ComplianceTable.findAll(
    { address: [address], blocked: true },
    [],
  );
  if (complianceData.length > 0) {
    return create4xxResponse(
      res,
      INDEXER_COMPLIANCE_BLOCKED_PAYLOAD,
      403,
      { code: BlockedCode.COMPLIANCE_BLOCKED },
    );
  }

  return next();
}
