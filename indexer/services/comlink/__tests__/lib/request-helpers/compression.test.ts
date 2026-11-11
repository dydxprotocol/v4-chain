import { logger } from '@dydxprotocol-indexer/base';
import express from 'express';
import http from 'http';
import { AddressInfo } from 'net';
import zlib from 'zlib';

import config, { configSchema } from '../../../src/config';
import server from '../../../src/request-helpers/server';

// Well above the 1 KB compression threshold.
const bigBody = {
  items: Array.from({ length: 200 }, (_, i) => ({
    id: i,
    ticker: 'BTC-USD',
    price: '65000.12345',
    size: '0.0100',
  })),
};
const smallBody = { ok: true };

interface RawResponse {
  statusCode: number,
  headers: http.IncomingHttpHeaders,
  raw: Buffer,
}

// Plain http client so the raw bytes on the wire can be inspected (supertest auto-decompresses).
async function rawGet(
  app: express.Express,
  path: string,
  headers: http.OutgoingHttpHeaders = {},
): Promise<RawResponse> {
  const listener: http.Server = app.listen(0);
  try {
    const { port } = listener.address() as AddressInfo;
    return await new Promise<RawResponse>((resolve, reject) => {
      http.get({ port, path, headers }, (res: http.IncomingMessage) => {
        const chunks: Buffer[] = [];
        res.on('data', (chunk: Buffer) => chunks.push(chunk));
        res.on('end', () => resolve({
          statusCode: res.statusCode!,
          headers: res.headers,
          raw: Buffer.concat(chunks),
        }));
        res.on('error', reject);
      }).on('error', reject);
    });
  } finally {
    listener.close();
  }
}

function buildApp(): express.Express {
  const router: express.Router = express.Router();
  router.get('/big', (_req: express.Request, res: express.Response) => {
    res.json(bigBody);
  });
  router.get('/small', (_req: express.Request, res: express.Response) => {
    res.json(smallBody);
  });
  router.get('/err', (_req: express.Request, res: express.Response) => {
    res.status(400).json(bigBody);
  });
  return server(router);
}

describe('comlink response compression', () => {
  const originalCompressionEnabled: boolean = config.COMPRESSION_ENABLED;

  beforeEach(() => {
    // Don't depend on the environment's rollback setting; tests that need it off set it below.
    config.COMPRESSION_ENABLED = true;
  });

  afterEach(() => {
    config.COMPRESSION_ENABLED = originalCompressionEnabled;
    jest.restoreAllMocks();
  });

  it('is enabled by default when COMPRESSION_ENABLED is unset', () => {
    const originalEnv: string | undefined = process.env.COMPRESSION_ENABLED;
    delete process.env.COMPRESSION_ENABLED;
    try {
      expect(configSchema.COMPRESSION_ENABLED('COMPRESSION_ENABLED')).toBe(true);
    } finally {
      if (originalEnv !== undefined) {
        process.env.COMPRESSION_ENABLED = originalEnv;
      }
    }
  });

  it('gzips large JSON responses when the client accepts gzip', async () => {
    const res: RawResponse = await rawGet(buildApp(), '/v4/big', { 'Accept-Encoding': 'gzip' });

    expect(res.statusCode).toBe(200);
    expect(res.headers['content-encoding']).toBe('gzip');
    expect(res.headers.vary).toMatch(/Accept-Encoding/i);
    expect(res.headers['content-length']).toBeUndefined();

    const uncompressed: Buffer = zlib.gunzipSync(res.raw);
    expect(JSON.parse(uncompressed.toString('utf8'))).toEqual(bigBody);
    expect(res.raw.length).toBeLessThan(uncompressed.length / 4);
  });

  it('keeps CORS headers on compressed responses', async () => {
    const res: RawResponse = await rawGet(buildApp(), '/v4/big', { 'Accept-Encoding': 'gzip' });

    expect(res.headers['content-encoding']).toBe('gzip');
    expect(res.headers['access-control-allow-origin']).toBe(config.CORS_ORIGIN);
  });

  it.each([
    ['identity', { 'Accept-Encoding': 'identity' }],
    ['no Accept-Encoding header', {}],
  ])('returns plain JSON for %s', async (_name: string, headers: http.OutgoingHttpHeaders) => {
    const res: RawResponse = await rawGet(buildApp(), '/v4/big', headers);

    expect(res.statusCode).toBe(200);
    expect(res.headers['content-encoding']).toBeUndefined();
    expect(JSON.parse(res.raw.toString('utf8'))).toEqual(bigBody);
  });

  it('does not compress responses below the size threshold', async () => {
    const res: RawResponse = await rawGet(buildApp(), '/v4/small', { 'Accept-Encoding': 'gzip' });

    expect(res.headers['content-encoding']).toBeUndefined();
    expect(JSON.parse(res.raw.toString('utf8'))).toEqual(smallBody);
  });

  it('does not compress /health', async () => {
    const res: RawResponse = await rawGet(buildApp(), '/health', { 'Accept-Encoding': 'gzip' });

    expect(res.headers['content-encoding']).toBeUndefined();
    expect(JSON.parse(res.raw.toString('utf8'))).toEqual({ ok: true });
  });

  it('logs the uncompressed error body when the response is gzipped', async () => {
    const debugSpy: jest.SpyInstance = jest.spyOn(logger, 'debug');

    const res: RawResponse = await rawGet(buildApp(), '/v4/err', { 'Accept-Encoding': 'gzip' });
    expect(res.statusCode).toBe(400);
    expect(res.headers['content-encoding']).toBe('gzip');

    // RequestLogger logs on the response 'finish' event.
    await new Promise((resolve) => setImmediate(resolve));

    const requestLog = debugSpy.mock.calls
      .map((call) => call[0])
      .find((log) => log.at === 'requestLogger#logRequest');
    expect(requestLog).toBeDefined();
    expect(requestLog.message.response.statusCode).toBe(400);
    expect(JSON.parse(requestLog.message.response.errorBody)).toEqual(bigBody);
  });

  it('does not compress when COMPRESSION_ENABLED is false', async () => {
    config.COMPRESSION_ENABLED = false;

    const res: RawResponse = await rawGet(buildApp(), '/v4/big', { 'Accept-Encoding': 'gzip' });

    expect(res.statusCode).toBe(200);
    expect(res.headers['content-encoding']).toBeUndefined();
    expect(JSON.parse(res.raw.toString('utf8'))).toEqual(bigBody);
  });
});
