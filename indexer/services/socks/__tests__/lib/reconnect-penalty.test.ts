import { ReconnectPenalty } from '../../src/lib/reconnect-penalty';
import config from '../../src/config';

describe('ReconnectPenalty', () => {
  let penalty: ReconnectPenalty;

  const connectionId: string = 'connectionId';
  const clientIp: string = '203.0.113.7';
  const reason: string = 'subscription-limit-abuse';

  const defaultPenaltyMs: number = config.RECONNECT_PENALTY_MS;
  const defaultEnabled: boolean = config.RECONNECT_PENALTY_ENABLED;
  const defaultMaxTracked: number = config.RECONNECT_PENALTY_MAX_TRACKED_CLIENTS;

  // Fake timers are unusable under this jest/node combination, and this class only depends on
  // Date.now, so drive the clock directly.
  let nowMs: number = 1_700_000_000_000;

  function advanceTime(ms: number): void {
    nowMs += ms;
  }

  beforeEach(() => {
    nowMs = 1_700_000_000_000;
    jest.spyOn(Date, 'now').mockImplementation(() => nowMs);
    config.RECONNECT_PENALTY_ENABLED = true;
    config.RECONNECT_PENALTY_MS = defaultPenaltyMs;
    config.RECONNECT_PENALTY_MAX_TRACKED_CLIENTS = defaultMaxTracked;
    penalty = new ReconnectPenalty();
  });

  afterEach(() => {
    jest.restoreAllMocks();
    config.RECONNECT_PENALTY_ENABLED = defaultEnabled;
    config.RECONNECT_PENALTY_MS = defaultPenaltyMs;
    config.RECONNECT_PENALTY_MAX_TRACKED_CLIENTS = defaultMaxTracked;
  });

  it('reports no penalty for an untracked client', () => {
    expect(penalty.getRemainingPenaltyMs(clientIp)).toEqual(0);
  });

  it('reports no penalty for an undefined client address', () => {
    expect(penalty.getRemainingPenaltyMs(undefined)).toEqual(0);
  });

  it('penalizes the client behind a tracked connection', () => {
    penalty.trackConnection(connectionId, clientIp);
    penalty.penalizeConnection(connectionId, reason);

    expect(penalty.getRemainingPenaltyMs(clientIp)).toEqual(config.RECONNECT_PENALTY_MS);
  });

  it('does not penalize other clients', () => {
    penalty.trackConnection(connectionId, clientIp);
    penalty.penalizeConnection(connectionId, reason);

    expect(penalty.getRemainingPenaltyMs('198.51.100.1')).toEqual(0);
  });

  it('is a no-op when the connection has no known client address', () => {
    penalty.penalizeConnection('unknownConnectionId', reason);

    expect(penalty.getRemainingPenaltyMs(clientIp)).toEqual(0);
  });

  it('counts the penalty down and expires it', () => {
    penalty.trackConnection(connectionId, clientIp);
    penalty.penalizeConnection(connectionId, reason);

    advanceTime(config.RECONNECT_PENALTY_MS - 1000);
    expect(penalty.getRemainingPenaltyMs(clientIp)).toEqual(1000);

    advanceTime(1000);
    expect(penalty.getRemainingPenaltyMs(clientIp)).toEqual(0);
  });

  it('outlives the connection that earned it', () => {
    penalty.trackConnection(connectionId, clientIp);
    penalty.penalizeConnection(connectionId, reason);
    penalty.removeConnection(connectionId);

    expect(penalty.getRemainingPenaltyMs(clientIp)).toEqual(config.RECONNECT_PENALTY_MS);
  });

  it('stops mapping a removed connection to its client address', () => {
    penalty.trackConnection(connectionId, clientIp);
    penalty.removeConnection(connectionId);
    penalty.penalizeConnection(connectionId, reason);

    expect(penalty.getRemainingPenaltyMs(clientIp)).toEqual(0);
  });

  it('extends an existing penalty on a repeat offence', () => {
    penalty.trackConnection(connectionId, clientIp);
    penalty.penalizeConnection(connectionId, reason);

    advanceTime(config.RECONNECT_PENALTY_MS - 1000);
    penalty.penalizeConnection(connectionId, reason);

    expect(penalty.getRemainingPenaltyMs(clientIp)).toEqual(config.RECONNECT_PENALTY_MS);
  });

  it('does nothing when disabled', () => {
    config.RECONNECT_PENALTY_ENABLED = false;
    penalty.trackConnection(connectionId, clientIp);
    penalty.penalizeConnection(connectionId, reason);

    expect(penalty.getRemainingPenaltyMs(clientIp)).toEqual(0);
  });

  it('bounds the number of tracked clients', () => {
    config.RECONNECT_PENALTY_MAX_TRACKED_CLIENTS = 4;

    for (let i = 0; i < 20; i += 1) {
      penalty.penalizeClient(`203.0.113.${i}`, reason);
      // Stagger so that eviction of the soonest-expiring entry is well defined.
      advanceTime(1);
    }

    let tracked: number = 0;
    for (let i = 0; i < 20; i += 1) {
      if (penalty.getRemainingPenaltyMs(`203.0.113.${i}`) > 0) {
        tracked += 1;
      }
    }

    expect(tracked).toBeLessThanOrEqual(config.RECONNECT_PENALTY_MAX_TRACKED_CLIENTS);
    // The most recent offender must still be penalised.
    expect(penalty.getRemainingPenaltyMs('203.0.113.19')).toBeGreaterThan(0);
  });

  it('clears all state', () => {
    penalty.trackConnection(connectionId, clientIp);
    penalty.penalizeConnection(connectionId, reason);
    penalty.clear();

    expect(penalty.getRemainingPenaltyMs(clientIp)).toEqual(0);
  });
});
