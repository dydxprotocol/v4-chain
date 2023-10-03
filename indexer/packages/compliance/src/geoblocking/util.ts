import config from '../config';

const RESTRICTED_COUNTRY_CODES: Set<string> = new Set(config.RESTRICTED_COUNTRIES.split(','));

export function isRestrictedCountry(ipCountry: string): boolean {
  return RESTRICTED_COUNTRY_CODES.has(ipCountry);
}
