import express from 'express';

import { create4xxResponse } from './helpers';

/**
 * Parses a comma-separated allowlist string into a normalized Set of origins.
 *
 * Origins are trimmed and compared case-insensitively (per RFC 6454 the scheme
 * and host portions of an origin are case-insensitive, and browsers normalize
 * the `Origin` header to lowercase in practice).
 */
function parseAllowedOrigins(allowedOrigins: string): Set<string> {
  return new Set(
    allowedOrigins
      .split(',')
      .map((o) => o.trim().toLowerCase())
      .filter((o) => o.length > 0),
  );
}

/**
 * Express middleware that restricts requests to a configured allowlist of
 * `Origin` header values.
 *
 * Behavior:
 *   - If `allowedOrigins` is empty (no list configured), the middleware is a
 *     no-op and lets all requests through. This preserves backwards
 *     compatibility for deployments that have not opted in.
 *   - If the request carries an `Origin` header that is not in the allowlist,
 *     the request is rejected with HTTP 403.
 *   - If the request has no `Origin` header (e.g. native mobile apps,
 *     server-to-server callers, curl), the middleware lets the request
 *     through. Browser-driven cross-origin requests are still gated by the
 *     global CORS middleware; this allowlist is an additional defense in depth
 *     specifically for clients that DO advertise an origin.
 *
 * @param allowedOrigins Comma-separated list of allowed origins
 *   (e.g. `"https://dydx.trade,https://staging.dydx.trade"`).
 */
export function originAllowlistMiddleware(allowedOrigins: string) {
  const allowed: Set<string> = parseAllowedOrigins(allowedOrigins);

  return (
    req: express.Request,
    res: express.Response,
    next: express.NextFunction,
  ): void | express.Response => {
    if (allowed.size === 0) {
      return next();
    }

    // Signal to caches that this response depends on the Origin header.
    res.vary('Origin');

    const origin: string | undefined = req.get('Origin');
    if (origin !== undefined && !allowed.has(origin.toLowerCase())) {
      return create4xxResponse(res, 'Origin not allowed', 403);
    }

    return next();
  };
}
