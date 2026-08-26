import { IncomingMessage as IncomingMessageHttp } from 'http';

import { GeoOriginHeaders } from '@dydxprotocol-indexer/compliance';

import { IncomingMessage } from '../types';

export function getGeoOriginHeaders(req: IncomingMessage): GeoOriginHeaders {
  const geoOriginHeaders = req.headers as GeoOriginHeaders;
  return {
    'geo-origin-country': geoOriginHeaders['geo-origin-country'],
    'geo-origin-region': geoOriginHeaders['geo-origin-region'],
    'geo-origin-status': geoOriginHeaders['geo-origin-status'],
  } as GeoOriginHeaders;
}

/**
 * Derives the client address for a websocket request.
 *
 * socks sits behind Cloudflare and an AWS load balancer, so the socket's remote address is a
 * proxy, not the client. `cf-connecting-ip` is set by Cloudflare and is the most trustworthy
 * value when traffic is proxied; otherwise fall back to the left-most `x-forwarded-for` entry,
 * and finally to the socket address for direct connections (local/dev, or health checks).
 *
 * @param req HTTP request accompanying the websocket upgrade.
 * @returns the client address, or undefined if it cannot be determined.
 */
export function getClientIp(req: IncomingMessageHttp): string | undefined {
  const cfConnectingIp = req.headers['cf-connecting-ip'];
  if (typeof cfConnectingIp === 'string' && cfConnectingIp.trim().length > 0) {
    return cfConnectingIp.trim();
  }

  const forwardedFor = req.headers['x-forwarded-for'];
  const forwardedForStr = Array.isArray(forwardedFor) ? forwardedFor[0] : forwardedFor;
  if (typeof forwardedForStr === 'string') {
    const first = forwardedForStr.split(',')[0].trim();
    if (first.length > 0) {
      return first;
    }
  }

  return req.socket?.remoteAddress;
}
