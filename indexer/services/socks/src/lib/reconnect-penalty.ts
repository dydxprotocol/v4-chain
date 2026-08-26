import { getInstanceId, logger, stats } from '@dydxprotocol-indexer/base';

import config from '../config';

/**
 * Tracks clients that were disconnected for abusing the subscription limit, so that their
 * reconnection attempts can be refused at the websocket handshake.
 *
 * This exists because a websocket close frame is invisible to the edge: Cloudflare proxies the
 * tunnelled connection but does not inspect its frames, so it cannot rate limit on "was dropped".
 * The handshake is HTTP, so refusing it with a 429 gives Cloudflare something it can count and
 * escalate on, and gives the client an unambiguous, documented reason for the refusal.
 *
 * Penalties are held in memory and are per-socks-instance. That is deliberate: the goal is to
 * stop a single instance from being monopolised, and to surface a durable signal to the edge,
 * which is where cross-instance enforcement belongs.
 */
export class ReconnectPenalty {
  // Client address -> epoch ms at which the penalty expires.
  private penaltyExpiryMs: { [clientIp: string]: number };
  // Connection id -> client address, so a connection can be penalised by the code that drops it
  // without having to thread request headers through the subscription layer.
  private clientIpByConnectionId: { [connectionId: string]: string };

  constructor() {
    this.penaltyExpiryMs = {};
    this.clientIpByConnectionId = {};
  }

  /**
   * Associates a connection with the client address it originated from.
   */
  public trackConnection(connectionId: string, clientIp?: string): void {
    if (clientIp === undefined) {
      return;
    }
    this.clientIpByConnectionId[connectionId] = clientIp;
  }

  /**
   * Stops tracking a closed connection. Any penalty on the client address is intentionally left
   * in place, since it must outlive the connection that earned it.
   */
  public removeConnection(connectionId: string): void {
    delete this.clientIpByConnectionId[connectionId];
  }

  /**
   * Places the client behind `connectionId` in the penalty box.
   */
  public penalizeConnection(connectionId: string, reason: string): void {
    const clientIp: string | undefined = this.clientIpByConnectionId[connectionId];
    if (clientIp === undefined) {
      return;
    }
    this.penalizeClient(clientIp, reason);
  }

  public penalizeClient(clientIp: string, reason: string): void {
    if (!config.RECONNECT_PENALTY_ENABLED) {
      return;
    }

    this.evictIfFull();

    this.penaltyExpiryMs[clientIp] = Date.now() + config.RECONNECT_PENALTY_MS;

    logger.info({
      at: 'reconnect-penalty#penalizeClient',
      message: 'Client placed in reconnect penalty box',
      reason,
      penaltyMs: config.RECONNECT_PENALTY_MS,
    });

    stats.increment(
      `${config.SERVICE_NAME}.reconnect_penalty_applied`,
      1,
      undefined,
      {
        reason,
        instance: getInstanceId(),
      },
    );
  }

  /**
   * @returns the number of milliseconds remaining on the client's penalty, or 0 if the client is
   * not penalised.
   */
  public getRemainingPenaltyMs(clientIp?: string): number {
    if (!config.RECONNECT_PENALTY_ENABLED || clientIp === undefined) {
      return 0;
    }

    const expiryMs: number | undefined = this.penaltyExpiryMs[clientIp];
    if (expiryMs === undefined) {
      return 0;
    }

    const remainingMs: number = expiryMs - Date.now();
    if (remainingMs <= 0) {
      delete this.penaltyExpiryMs[clientIp];
      return 0;
    }

    return remainingMs;
  }

  /**
   * Drops expired entries, and if the map is still at capacity, the entry expiring soonest. This
   * bounds memory so a client cycling through addresses cannot grow the map without limit.
   */
  private evictIfFull(): void {
    if (Object.keys(this.penaltyExpiryMs).length < config.RECONNECT_PENALTY_MAX_TRACKED_CLIENTS) {
      return;
    }

    const now: number = Date.now();
    Object.entries(this.penaltyExpiryMs).forEach(([clientIp, expiryMs]: [string, number]) => {
      if (expiryMs <= now) {
        delete this.penaltyExpiryMs[clientIp];
      }
    });

    const entries: [string, number][] = Object.entries(this.penaltyExpiryMs);
    if (entries.length >= config.RECONNECT_PENALTY_MAX_TRACKED_CLIENTS && entries.length > 0) {
      const soonest: [string, number] = entries.reduce(
        (a: [string, number], b: [string, number]) => (a[1] <= b[1] ? a : b),
      );
      delete this.penaltyExpiryMs[soonest[0]];
    }
  }

  /**
   * Clears all tracked state. Used to reset between tests.
   */
  public clear(): void {
    this.penaltyExpiryMs = {};
    this.clientIpByConnectionId = {};
  }
}

// Shared instance. The subscription layer writes to it when it drops a connection, and the
// websocket server reads from it when verifying an upgrade.
export const reconnectPenalty: ReconnectPenalty = new ReconnectPenalty();
