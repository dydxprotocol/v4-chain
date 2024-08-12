// Gets `goodTilBlockTime` from an ISO string.
import { DateTime } from 'luxon';

import { IsoString } from '../types';
import vaultAddresses from './vault-addresses.json';

export function blockTimeFromIsoString(isoString: IsoString): number {
  const dateTime: DateTime = DateTime.fromISO(isoString, { zone: 'utc' });
  return Math.floor(dateTime.toMillis() / 1000);
}

export function getVaultAddresses(): string[] {
  return vaultAddresses;
}
