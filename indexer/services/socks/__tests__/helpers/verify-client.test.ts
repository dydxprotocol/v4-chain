import { IncomingMessage as IncomingMessageHttp, OutgoingHttpHeaders } from 'http';
import { Socket } from 'net';

import config from '../../src/config';
import { verifyClientNotPenalized } from '../../src/helpers/wss';
import {
  RATE_LIMIT_REASON_HEADER,
  RATE_LIMIT_REASON_SUBSCRIPTION_LIMIT_ABUSE,
} from '../../src/lib/constants';
import { reconnectPenalty } from '../../src/lib/reconnect-penalty';

describe('verifyClientNotPenalized', () => {
  const clientIp: string = '203.0.113.7';
  const defaultPenaltyMs: number = config.RECONNECT_PENALTY_MS;
  const defaultTrust: boolean = config.TRUST_FORWARDED_HEADERS;

  function infoForIp(ip?: string): {
    origin: string, secure: boolean, req: IncomingMessageHttp,
  } {
    const req: IncomingMessageHttp = new IncomingMessageHttp(new Socket());
    if (ip !== undefined) {
      req.headers['cf-connecting-ip'] = ip;
    }
    return { origin: 'https://dydx.trade', secure: true, req };
  }

  let nowMs: number = 1_700_000_000_000;

  function advanceTime(ms: number): void {
    nowMs += ms;
  }

  beforeEach(() => {
    nowMs = 1_700_000_000_000;
    jest.spyOn(Date, 'now').mockImplementation(() => nowMs);
    // These cases describe a deployment whose ingress is locked to the trusted proxy chain, so
    // the forwarded client address can be believed.
    config.TRUST_FORWARDED_HEADERS = true;
    config.RECONNECT_PENALTY_ENABLED = true;
    config.RECONNECT_PENALTY_MS = defaultPenaltyMs;
    reconnectPenalty.clear();
  });

  afterEach(() => {
    jest.restoreAllMocks();
    reconnectPenalty.clear();
    config.RECONNECT_PENALTY_MS = defaultPenaltyMs;
    config.TRUST_FORWARDED_HEADERS = defaultTrust;
  });

  it('accepts the upgrade when forwarded headers are not trusted', () => {
    config.TRUST_FORWARDED_HEADERS = false;
    reconnectPenalty.penalizeClient(clientIp, RATE_LIMIT_REASON_SUBSCRIPTION_LIMIT_ABUSE);
    const callback: jest.Mock = jest.fn();

    // The address cannot be established safely, so the cooldown must not be applied to a
    // caller-supplied value.
    verifyClientNotPenalized(infoForIp(clientIp), callback);

    expect(callback).toHaveBeenCalledWith(true);
  });

  it('accepts the upgrade for a client with no penalty', () => {
    const callback: jest.Mock = jest.fn();

    verifyClientNotPenalized(infoForIp(clientIp), callback);

    expect(callback).toHaveBeenCalledWith(true);
  });

  it('rejects the upgrade with HTTP 429 while the client is penalised', () => {
    reconnectPenalty.penalizeClient(clientIp, RATE_LIMIT_REASON_SUBSCRIPTION_LIMIT_ABUSE);
    const callback: jest.Mock = jest.fn();

    verifyClientNotPenalized(infoForIp(clientIp), callback);

    expect(callback).toHaveBeenCalledTimes(1);
    const [result, code, statusMessage, headers] = callback.mock.calls[0] as [
      boolean, number, string, OutgoingHttpHeaders,
    ];
    expect(result).toEqual(false);
    // The 429 is the signal Cloudflare can see; a websocket close frame is not.
    expect(code).toEqual(429);
    expect(statusMessage).toEqual('Too Many Requests');
    expect(headers['Retry-After']).toEqual(String(config.RECONNECT_PENALTY_MS / 1000));
    expect(headers[RATE_LIMIT_REASON_HEADER])
      .toEqual(RATE_LIMIT_REASON_SUBSCRIPTION_LIMIT_ABUSE);
  });

  it('rounds Retry-After up to the next whole second', () => {
    config.RECONNECT_PENALTY_MS = 1500;
    reconnectPenalty.penalizeClient(clientIp, RATE_LIMIT_REASON_SUBSCRIPTION_LIMIT_ABUSE);
    const callback: jest.Mock = jest.fn();

    verifyClientNotPenalized(infoForIp(clientIp), callback);

    expect(callback.mock.calls[0][3]['Retry-After']).toEqual('2');
  });

  it('accepts the upgrade again once the penalty expires', () => {
    reconnectPenalty.penalizeClient(clientIp, RATE_LIMIT_REASON_SUBSCRIPTION_LIMIT_ABUSE);
    advanceTime(config.RECONNECT_PENALTY_MS);
    const callback: jest.Mock = jest.fn();

    verifyClientNotPenalized(infoForIp(clientIp), callback);

    expect(callback).toHaveBeenCalledWith(true);
  });

  it('does not penalise a different client', () => {
    reconnectPenalty.penalizeClient(clientIp, RATE_LIMIT_REASON_SUBSCRIPTION_LIMIT_ABUSE);
    const callback: jest.Mock = jest.fn();

    verifyClientNotPenalized(infoForIp('198.51.100.1'), callback);

    expect(callback).toHaveBeenCalledWith(true);
  });

  it('accepts the upgrade when the client address cannot be determined', () => {
    reconnectPenalty.penalizeClient(clientIp, RATE_LIMIT_REASON_SUBSCRIPTION_LIMIT_ABUSE);
    const callback: jest.Mock = jest.fn();

    verifyClientNotPenalized(infoForIp(undefined), callback);

    expect(callback).toHaveBeenCalledWith(true);
  });
});
