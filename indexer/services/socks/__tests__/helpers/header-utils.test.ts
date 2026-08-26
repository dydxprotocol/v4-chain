import { IncomingMessage as IncomingMessageHttp } from 'http';
import { Socket } from 'net';

import { getClientIp } from '../../src/helpers/header-utils';

describe('getClientIp', () => {
  function requestWithHeaders(
    headers: { [key: string]: string | string[] },
    remoteAddress?: string,
  ): IncomingMessageHttp {
    const socket: Socket = new Socket();
    if (remoteAddress !== undefined) {
      Object.defineProperty(socket, 'remoteAddress', { value: remoteAddress });
    }
    const req: IncomingMessageHttp = new IncomingMessageHttp(socket);
    Object.entries(headers).forEach(([key, value]) => {
      req.headers[key] = value;
    });
    return req;
  }

  it('prefers cf-connecting-ip, which is the only header carrying the real client past Cloudflare', () => {
    const req = requestWithHeaders({
      'cf-connecting-ip': '203.0.113.7',
      'x-forwarded-for': '198.51.100.1, 172.68.0.1',
    }, '10.0.1.5');

    expect(getClientIp(req)).toEqual('203.0.113.7');
  });

  it('falls back to the left-most x-forwarded-for entry', () => {
    const req = requestWithHeaders({
      'x-forwarded-for': '198.51.100.1, 172.68.0.1',
    }, '10.0.1.5');

    expect(getClientIp(req)).toEqual('198.51.100.1');
  });

  it('trims whitespace around the forwarded address', () => {
    const req = requestWithHeaders({ 'x-forwarded-for': '  198.51.100.1  ,172.68.0.1' });

    expect(getClientIp(req)).toEqual('198.51.100.1');
  });

  it('handles a repeated x-forwarded-for header', () => {
    const req = requestWithHeaders({ 'x-forwarded-for': ['198.51.100.1', '172.68.0.1'] });

    expect(getClientIp(req)).toEqual('198.51.100.1');
  });

  it('falls back to the socket address when no proxy headers are present', () => {
    const req = requestWithHeaders({}, '10.0.1.5');

    expect(getClientIp(req)).toEqual('10.0.1.5');
  });

  it('ignores an empty cf-connecting-ip', () => {
    const req = requestWithHeaders({
      'cf-connecting-ip': '   ',
      'x-forwarded-for': '198.51.100.1',
    });

    expect(getClientIp(req)).toEqual('198.51.100.1');
  });

  it('ignores an empty x-forwarded-for', () => {
    const req = requestWithHeaders({ 'x-forwarded-for': '' }, '10.0.1.5');

    expect(getClientIp(req)).toEqual('10.0.1.5');
  });
});
