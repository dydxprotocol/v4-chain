import { DateTime } from 'luxon';

import { IsoString } from '../types';
import vaultAddresses from './vault-addresses.json';

// Gets `goodTilBlockTime` from an ISO string.
export function blockTimeFromIsoString(isoString: IsoString): number {
  const dateTime: DateTime = DateTime.fromISO(isoString, { zone: 'utc' });
  return Math.floor(dateTime.toMillis() / 1000);
}

// TODO (ENG-65): future-proof pre-generated vault addresses.
// Set of the addresses of vaults that quote on clob pairs 0 to 999.
export const VAULTS_CLOB_0_TO_999: Set<string> = new Set(vaultAddresses);
// Comma-separated list of the addresses of vaults that quote on clob pairs 0 to 999.
export const VAULTS_CLOB_0_TO_999_STR_CONCAT: string = vaultAddresses.map((address) => `'${address}'`).join(',');
