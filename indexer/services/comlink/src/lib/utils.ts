import express from 'express';

import config from '../config';

const indexerIps: string[] = config.INDEXER_INTERNAL_IPS.split(',').map((l) => l.toLowerCase());

export function isIndexerIp(ipAddress: string): boolean {
  return indexerIps.includes(ipAddress);
}

export function getIpAddr(req: express.Request): string | undefined {
  const {
    'cf-connecting-ip': cloudflareIP,
    'x-forwarded-for': loadBalancerHeader,
  } = req.headers as {
    'cf-connecting-ip'?: string,
    'x-forwarded-for'?: string,
  };

  // get ip address
  const loadBalancerIPs: string[] | undefined = loadBalancerHeader?.replace(/\s/g, '').split(',');
  const firstLoadBalancerIP: string | undefined = loadBalancerIPs?.[0];

  return cloudflareIP || firstLoadBalancerIP;
}
