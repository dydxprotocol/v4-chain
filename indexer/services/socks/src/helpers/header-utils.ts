import { IncomingMessage as IncomingMessageHttp } from 'http';

import { GeoOriginHeaders } from '@dydxprotocol-indexer/compliance';

import config from '../config';
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
 * `cf-connecting-ip` and `x-forwarded-for` are set by the caller, and a load balancer forwards
 * whatever it was given. They are believable here because origin ingress accepts connections
 * from Cloudflare proxies alone, so the caller is always Cloudflare.
 * `TRUST_FORWARDED_HEADERS` exists to withdraw that trust if the origin ever becomes directly
 * reachable, since a caller could then name an arbitrary victim and have that address penalised.
 *
 * When a forwarded header is present but not trusted, this returns undefined rather than the
 * socket address. Behind a proxy the socket address is the proxy itself, shared by every client,
 * so keying a penalty on it would refuse everyone. Returning undefined makes the address-based
 * cooldown inert, which is the safe failure mode; the per-connection drop is unaffected.
 *
 * @param req HTTP request accompanying the websocket upgrade.
 * @returns the client address, or undefined if it cannot be established safely.
 */
export function getClientIp(req: IncomingMessageHttp): string | undefined {
  const cfConnectingIp = req.headers['cf-connecting-ip'];
  const forwardedFor = req.headers['x-forwarded-for'];
  const hasForwardedHeaders: boolean = cfConnectingIp !== undefined || forwardedFor !== undefined;

  if (!hasForwardedHeaders) {
    // No proxy in the path, so the peer is the client.
    return req.socket?.remoteAddress;
  }

  if (!config.TRUST_FORWARDED_HEADERS) {
    return undefined;
  }

  if (typeof cfConnectingIp === 'string' && cfConnectingIp.trim().length > 0) {
    return cfConnectingIp.trim();
  }

  const forwardedForStr = Array.isArray(forwardedFor) ? forwardedFor[0] : forwardedFor;
  if (typeof forwardedForStr === 'string') {
    const first = forwardedForStr.split(',')[0].trim();
    if (first.length > 0) {
      return first;
    }
  }

  return req.socket?.remoteAddress;
}
