import { logger, stats } from '@dydxprotocol-indexer/base';
import { redis } from '@dydxprotocol-indexer/redis';
import { RateLimiterRedis } from 'rate-limiter-flexible';
import { createClient, RedisClient } from 'redis';

import {
  defaultRateLimiter,
  ratelimitRedis,
  resetRateLimitStoreThrottles,
  screenProviderGlobalLimiter,
} from '../../src/caches/rate-limiters';
import config from '../../src/config';
import { TurnkeyError } from '../../src/lib/errors';
import { clearStoreErrorLogTimes, rateLimiterMiddleware } from '../../src/lib/rate-limit';
import { reserveNewSuborgSlot } from '../../src/lib/turnkey-helpers';

// These tests drive the real rate-limit client against the test Redis, stalling, demoting or
// locking that Redis through a second connection. They need DEBUG, REPLICAOF and CONFIG SET, which
// the compose image (redis:5.0.6-alpine) allows; Redis 7 refuses DEBUG by default.
describe('rate-limit store', () => {
  const client: any = ratelimitRedis.client;
  let admin: RedisClient;
  let req: any;
  let res: any;
  let next: jest.Mock;

  const send = (name: string, args: string[] = []): Promise<unknown> => new Promise(
    (resolve, reject) => admin.send_command(name, args, (e, r) => (e ? reject(e) : resolve(r))),
  );
  const sleep = (ms: number): Promise<void> => new Promise((r) => setTimeout(r, ms));
  const waitFor = async (check: () => boolean, ms: number): Promise<boolean> => {
    const deadline: number = Date.now() + ms;
    while (Date.now() < deadline) {
      if (check()) {
        return true;
      }
      await sleep(20);
    }
    return check();
  };

  beforeAll(async () => {
    config.RATE_LIMIT_ENABLED = true;
    admin = createClient({ url: config.RATE_LIMIT_REDIS_URL });
    admin.on('error', () => {});
    await ratelimitRedis.connect();
  });

  beforeEach(() => {
    req = { headers: { 'cf-connecting-ip': '203.0.113.20' }, get: jest.fn().mockReturnThis() };
    res = {
      status: jest.fn().mockReturnThis(),
      json: jest.fn().mockReturnThis(),
      set: jest.fn().mockReturnThis(),
    };
    next = jest.fn();
  });

  afterEach(async () => {
    jest.restoreAllMocks();
    clearStoreErrorLogTimes();
    resetRateLimitStoreThrottles();
    expect(await waitFor(() => client.ready, 5000)).toBe(true);
    await redis.deleteAllAsync(client);
  });

  afterAll(() => {
    config.RATE_LIMIT_ENABLED = false;
    admin.quit();
  });

  describe('stalled store', () => {
    const timeoutMs: number = config.RATE_LIMIT_STORE_TIMEOUT_MS;
    let stalled: Promise<unknown>;

    beforeEach(async () => {
      config.RATE_LIMIT_STORE_TIMEOUT_MS = 200;
      // The client stays ready while the server does not answer, as on a half-open connection.
      stalled = send('debug', ['sleep', '2']);
      await sleep(50);
    });

    afterEach(async () => {
      await stalled;
      config.RATE_LIMIT_STORE_TIMEOUT_MS = timeoutMs;
    });

    it('fails open within the timeout and does not run the request again later', async () => {
      const increment = jest.spyOn(stats, 'increment');
      const start: number = Date.now();
      await rateLimiterMiddleware(defaultRateLimiter)(req, res, next);
      expect(Date.now() - start).toBeLessThan(1000);
      expect(next).toHaveBeenCalledTimes(1);
      expect(increment).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.rate_limit_store_error`, 1,
        { key_prefix: defaultRateLimiter.keyPrefix },
      );
      await stalled;
      await sleep(50);
      expect(next).toHaveBeenCalledTimes(1);
    });

    it('refuses an opt-in route within the timeout', async () => {
      const start: number = Date.now();
      await rateLimiterMiddleware(defaultRateLimiter, undefined, undefined, true)(req, res, next);
      expect(Date.now() - start).toBeLessThan(1000);
      expect(res.status).toHaveBeenCalledWith(503);
      expect(res.json).toHaveBeenCalledWith({
        errors: [{ code: 'RATE_LIMIT_STORE_UNAVAILABLE', msg: 'Service temporarily unavailable' }],
      });
      await stalled;
      await sleep(50);
      expect(next).not.toHaveBeenCalled();
    });

    it('refuses a new sub-org within the timeout', async () => {
      const start: number = Date.now();
      await expect(reserveNewSuborgSlot()).rejects.toThrow(TurnkeyError);
      expect(Date.now() - start).toBeLessThan(1000);
    });

    it('drops the connection on a timeout, so later calls fail at once instead of queueing', async () => {
      await expect(defaultRateLimiter.consume('203.0.113.22')).rejects.toThrow('did not answer');
      expect(await waitFor(() => !client.ready, 500)).toBe(true);
      expect(client.command_queue.length).toBe(0);
      const start: number = Date.now();
      await expect(defaultRateLimiter.consume('203.0.113.22')).rejects.toThrow();
      expect(Date.now() - start).toBeLessThan(100);
    });

    it('drops again, within the floor, when the reconnect\'s ready check gets no answer', async () => {
      let ends: number = 0;
      const onEnd = (): void => { ends += 1; };
      client.on('end', onEnd);
      try {
        await expect(defaultRateLimiter.consume('203.0.113.24')).rejects.toThrow('did not answer');
        // Reconnected, with the ready check waiting on the stalled server.
        expect(await waitFor(() => client.connected && !client.ready, 1000)).toBe(true);
        await expect(defaultRateLimiter.consume('203.0.113.24')).rejects.toThrow('did not answer');
        expect(await waitFor(() => ends === 2, 300)).toBe(true);
      } finally {
        client.removeListener('end', onEnd);
      }
    });

    it('answers screening from memory within the timeout, and counts it', async () => {
      const increment = jest.spyOn(stats, 'increment');
      const start: number = Date.now();
      await expect(screenProviderGlobalLimiter.consume('stalled')).resolves.toBeDefined();
      expect(Date.now() - start).toBeLessThan(1000);
      expect(increment).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.rate_limit_fallback_used`, 1,
        { key_prefix: screenProviderGlobalLimiter.keyPrefix },
      );
    });
  });

  it('reconnects when the server turns into a replica (READONLY)', async () => {
    await send('replicaof', ['127.0.0.1', '1']);
    try {
      await expect(defaultRateLimiter.consume('203.0.113.21')).rejects.toThrow('READONLY');
      // Reconnected to the replica, whose master link is down, so it is not ready.
      expect(await waitFor(() => !client.ready, 2000)).toBe(true);
    } finally {
      await send('replicaof', ['no', 'one']);
    }
    expect(await waitFor(() => client.ready, 5000)).toBe(true);
  });

  it('drops the connection at most once per 5 s on READONLY', async () => {
    jest.spyOn(RateLimiterRedis.prototype, 'consume').mockRejectedValue(
      new Error("ERR Error running script: -READONLY You can't write against a read only replica."),
    );
    const destroy = jest.spyOn(client.stream, 'destroy').mockImplementation(() => client.stream);
    const start: number = Date.now();
    const now = jest.spyOn(Date, 'now').mockReturnValue(start);

    await expect(defaultRateLimiter.consume('203.0.113.23')).rejects.toThrow('READONLY');
    now.mockReturnValue(start + 4_999);
    await expect(defaultRateLimiter.consume('203.0.113.23')).rejects.toThrow('READONLY');
    expect(destroy).toHaveBeenCalledTimes(1);
    now.mockReturnValue(start + 5_000);
    await expect(defaultRateLimiter.consume('203.0.113.23')).rejects.toThrow('READONLY');
    expect(destroy).toHaveBeenCalledTimes(2);
  });

  it('reconnects after a failed ready check once the cause is gone', async () => {
    const retryMs: number = config.RATE_LIMIT_STORE_READY_RETRY_MS;
    config.RATE_LIMIT_STORE_READY_RETRY_MS = 100;
    await send('config', ['set', 'requirepass', 'pw']);
    try {
      client.stream.destroy();
      expect(await waitFor(() => !client.ready, 1000)).toBe(true);
      // Every reconnect now fails its ready check with NOAUTH.
      await sleep(1500);
      expect(client.ready).toBe(false);
    } finally {
      await send('auth', ['pw']);
      await send('config', ['set', 'requirepass', '']);
      config.RATE_LIMIT_STORE_READY_RETRY_MS = retryMs;
    }
    expect(await waitFor(() => client.ready, 5000)).toBe(true);
  });

  it('logs reconnect attempts once per interval and counts every one', async () => {
    const error = jest.spyOn(logger, 'error');
    const distribution = jest.spyOn(stats, 'distribution');
    const start: number = Date.now() + 1_000_000_000;
    const now = jest.spyOn(Date, 'now').mockReturnValue(start);
    const options: any = { attempt: 50, error: new Error('ECONNREFUSED') };

    expect(client.options.retry_strategy(options)).toBe(config.REDIS_RECONNECT_TIMEOUT_MS);
    now.mockReturnValue(start + 9_999);
    expect(client.options.retry_strategy(options)).toBe(config.REDIS_RECONNECT_TIMEOUT_MS);
    expect(error).toHaveBeenCalledTimes(1);
    expect(distribution.mock.calls.filter(
      ([stat]) => stat === 'redis.connection_retries',
    )).toHaveLength(2);
    now.mockReturnValue(start + 10_000);
    client.options.retry_strategy(options);
    expect(error).toHaveBeenCalledTimes(2);

    // Once the client is ready again, the next outage's first attempt is logged at once.
    now.mockRestore();
    client.stream.destroy();
    expect(await waitFor(() => !client.ready, 1000)).toBe(true);
    expect(await waitFor(() => client.ready, 5000)).toBe(true);
    jest.spyOn(Date, 'now').mockReturnValue(start + 10_001);
    client.options.retry_strategy(options);
    expect(error).toHaveBeenCalledTimes(3);
  });

  it('does not time out a reply that arrived while the event loop was blocked', async () => {
    const timeoutMs: number = config.RATE_LIMIT_STORE_TIMEOUT_MS;
    config.RATE_LIMIT_STORE_TIMEOUT_MS = 200;
    let ends: number = 0;
    const onEnd = (): void => { ends += 1; };
    client.on('end', onEnd);
    try {
      // From the check phase, the next loop turn runs timers before it reads the socket; from the
      // timers phase it would not, and the test would pass without the fix.
      await new Promise((r) => setImmediate(r));
      const call: Promise<unknown> = defaultRateLimiter.consume('203.0.113.25');
      const end: number = Date.now() + 400;
      while (Date.now() < end) {
        // Busy-wait: Redis answers while comlink cannot read the reply.
      }
      await expect(call).resolves.toBeDefined();
      await sleep(50);
      expect(ends).toBe(0);
      expect(client.ready).toBe(true);
    } finally {
      client.removeListener('end', onEnd);
      config.RATE_LIMIT_STORE_TIMEOUT_MS = timeoutMs;
    }
  });

  // Last: without the fix the client ends for good, which would break the tests after it.
  it('keeps reconnecting after more than an hour of retries', async () => {
    client.retry_totaltime = 3_700_000;
    client.stream.destroy();
    expect(await waitFor(() => !client.ready, 1000)).toBe(true);
    expect(await waitFor(() => client.ready, 5000)).toBe(true);
  });
});
