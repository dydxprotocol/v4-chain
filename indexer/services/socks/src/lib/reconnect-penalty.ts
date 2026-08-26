import { getInstanceId, logger, stats } from '@dydxprotocol-indexer/base';

import config from '../config';

// Upper bound on how many expired entries a single penalty application will sweep. Keeps the
// sweep constant-time regardless of how large the map has grown.
const EVICTION_SWEEP_LIMIT: number = 8;

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
  //
  // A Map rather than an object so that iteration order is insertion order. Every penalty uses
  // the same duration, so insertion order is also expiry order, which makes both the expiry
  // sweep and the capacity eviction O(1) amortized from the front. That matters because
  // eviction runs precisely when the map is full -- i.e. while under attack.
  private penaltyExpiryMs: Map<string, number>;
  // Connection id -> client address, so a connection can be penalised by the code that drops it
  // without having to thread request headers through the subscription layer.
  private clientIpByConnectionId: Map<string, string>;

  constructor() {
    this.penaltyExpiryMs = new Map();
    this.clientIpByConnectionId = new Map();
  }

  /**
   * Associates a connection with the client address it originated from.
   */
  public trackConnection(connectionId: string, clientIp?: string): void {
    if (clientIp === undefined) {
      return;
    }
    this.clientIpByConnectionId.set(connectionId, clientIp);
  }

  /**
   * Stops tracking a closed connection. Any penalty on the client address is intentionally left
   * in place, since it must outlive the connection that earned it.
   */
  public removeConnection(connectionId: string): void {
    this.clientIpByConnectionId.delete(connectionId);
  }

  /**
   * Places the client behind `connectionId` in the penalty box.
   */
  public penalizeConnection(connectionId: string, reason: string): void {
    const clientIp: string | undefined = this.clientIpByConnectionId.get(connectionId);
    if (clientIp === undefined) {
      return;
    }
    this.penalizeClient(clientIp, reason);
  }

  public penalizeClient(clientIp: string, reason: string): void {
    if (!config.RECONNECT_PENALTY_ENABLED) {
      return;
    }

    // Re-insert rather than update, so a re-offending client moves to the back and the
    // insertion-order-equals-expiry-order invariant holds.
    this.penaltyExpiryMs.delete(clientIp);
    this.evictExpired();
    this.evictToCapacity();

    this.penaltyExpiryMs.set(clientIp, Date.now() + config.RECONNECT_PENALTY_MS);

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

    const expiryMs: number | undefined = this.penaltyExpiryMs.get(clientIp);
    if (expiryMs === undefined) {
      return 0;
    }

    const remainingMs: number = expiryMs - Date.now();
    if (remainingMs <= 0) {
      this.penaltyExpiryMs.delete(clientIp);
      return 0;
    }

    return remainingMs;
  }

  /**
   * Drops entries that have already expired, from the front and bounded per call, so the sweep
   * never scales with the size of the map.
   */
  private evictExpired(): void {
    const now: number = Date.now();
    let swept: number = 0;

    for (const [clientIp, expiryMs] of this.penaltyExpiryMs) {
      if (expiryMs > now || swept >= EVICTION_SWEEP_LIMIT) {
        break;
      }
      this.penaltyExpiryMs.delete(clientIp);
      swept += 1;
    }
  }

  /**
   * Makes room for one more entry by dropping the oldest, which under a uniform penalty duration
   * is also the soonest to expire.
   */
  private evictToCapacity(): void {
    while (this.penaltyExpiryMs.size >= config.RECONNECT_PENALTY_MAX_TRACKED_CLIENTS) {
      const oldest: string | undefined = this.penaltyExpiryMs.keys().next().value;
      if (oldest === undefined) {
        return;
      }
      this.penaltyExpiryMs.delete(oldest);
    }
  }

  /**
   * Clears all tracked state. Used to reset between tests.
   */
  public clear(): void {
    this.penaltyExpiryMs = new Map();
    this.clientIpByConnectionId = new Map();
  }
}

// Shared instance. The subscription layer writes to it when it drops a connection, and the
// websocket server reads from it when verifying an upgrade.
export const reconnectPenalty: ReconnectPenalty = new ReconnectPenalty();
