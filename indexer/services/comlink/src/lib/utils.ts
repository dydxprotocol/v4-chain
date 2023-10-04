import config from '../config';

const indexerIps: string[] = config.INDEXER_INTERNAL_IPS.split(',').map((l) => l.toLowerCase());
const RESTRICTED_COUNTRY_CODES: Set<string> = new Set(config.RESTRICTED_COUNTRIES.split(','));

export function isIndexerIp(ipAddress: string): boolean {
  return indexerIps.includes(ipAddress);
}

export function isRestrictedCountry(ipCountry: string): boolean {
  return RESTRICTED_COUNTRY_CODES.has(ipCountry);
}
