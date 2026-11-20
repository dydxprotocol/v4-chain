import { isIPv6 } from 'net';

import express from 'express';

import config from '../config';

const indexerIps: string[] = config.INDEXER_INTERNAL_IPS.split(',').map((l) => l.toLowerCase());

export function isIndexerIp(ipAddress: string): boolean {
  return indexerIps.includes(ipAddress);
}

/**
 * Rate-limit key for an IP: IPv4 as is, IPv6 as its /64 (or /48). A single IPv6 host usually holds
 * a whole /64 and one site often a /48, so per-address keys give them unlimited buckets.
 */
export function ipRateLimitKey(ipAddress: string, ipv6PrefixBits: 48 | 64 = 64): string {
  if (!isIPv6(ipAddress)) {
    return ipAddress;
  }
  const mappedV4: RegExpMatchArray | null = ipAddress.match(/^::ffff:(\d+\.\d+\.\d+\.\d+)$/i);
  if (mappedV4) {
    return mappedV4[1];
  }
  const [head, tail] = ipAddress.split('::');
  const headGroups: string[] = head ? head.split(':') : [];
  const tailGroups: string[] = tail ? tail.split(':') : [];
  const groups: string[] = tail === undefined
    ? headGroups
    : [
      ...headGroups,
      ...Array(8 - headGroups.length - tailGroups.length).fill('0'),
      ...tailGroups,
    ];
  const values: number[] = groups.map((g) => parseInt(g, 16));
  // IPv4-mapped in hex form (::ffff:102:304): key it as the IPv4 address, not as ::/64.
  if (values.length === 8 && values.slice(0, 5).every((v) => v === 0) && values[5] === 0xffff) {
    return [values[6], values[7]].flatMap((v) => [Math.floor(v / 256), v % 256]).join('.');
  }
  const groupCount: number = ipv6PrefixBits / 16;
  return `${values.slice(0, groupCount).map((v) => v.toString(16)).join(':')}::/${ipv6PrefixBits}`;
}

/**
 * Checks a client-supplied magic link against the allowed origin and path prefix.
 */
export function isAllowedMagicLink(magicLink: string, allowed: string): boolean {
  let link: URL;
  let allowedUrl: URL;
  try {
    link = new URL(magicLink);
    allowedUrl = new URL(allowed);
  } catch {
    return false;
  }
  // Turnkey emails the raw string, so require it to be in canonical form.
  return link.href === magicLink &&
    link.username === '' &&
    link.password === '' &&
    link.origin === allowedUrl.origin &&
    link.pathname.startsWith(allowedUrl.pathname);
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
