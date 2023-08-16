import config from '../config';

const indexerIps: string[] = config.INDEXER_INTERNAL_IPS.split(',').map((l) => l.toLowerCase());

export function isIndexerIp(ipAddress: string) {
  return indexerIps.includes(ipAddress);
}
