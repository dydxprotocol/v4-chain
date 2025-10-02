import { GeoOriginStatus } from '@dydxprotocol-indexer/compliance';
import express from 'express';
import request from 'supertest';

import geoHeadersMiddleware from '../../src/request-helpers/geo-headers-middleware';

describe('geoOriginHeadersMiddleware', () => {
  let app: express.Application;

  beforeEach(() => {
    app = express();
    app.use(geoHeadersMiddleware);

    app.get('/test', (req: express.Request, res: express.Response) => {
      res.json({ success: true });
    });
  });

  it('reflect all geo headers when all are present', async () => {
    const response = await request(app)
      .get('/test')
      .set('Geo-Origin-Country', 'US')
      .set('Geo-Origin-Region', 'CA')
      .set('Geo-Origin-Status', GeoOriginStatus.OK);

    expect(response.headers['geo-origin-country']).toBe('US');
    expect(response.headers['geo-origin-region']).toBe('CA');
    expect(response.headers['geo-origin-status']).toBe(GeoOriginStatus.OK);
  });

  it('reflect only Geo-Origin-Country when only country header is present', async () => {
    const response = await request(app)
      .get('/test')
      .set('Geo-Origin-Country', 'FR');

    expect(response.headers['geo-origin-country']).toBe('FR');
    expect(response.headers['geo-origin-region']).toBeUndefined();
    expect(response.headers['geo-origin-status']).toBeUndefined();
  });

  it('reflect only Geo-Origin-Region when only region header is present', async () => {
    const response = await request(app)
      .get('/test')
      .set('Geo-Origin-Region', 'NY');

    expect(response.headers['geo-origin-country']).toBeUndefined();
    expect(response.headers['geo-origin-region']).toBe('NY');
    expect(response.headers['geo-origin-status']).toBeUndefined();
  });

  it('reflect only Geo-Origin-Status when only status header is present', async () => {
    const response = await request(app)
      .get('/test')
      .set('Geo-Origin-Status', GeoOriginStatus.RESTRICTED);

    expect(response.headers['geo-origin-country']).toBeUndefined();
    expect(response.headers['geo-origin-region']).toBeUndefined();
    expect(response.headers['geo-origin-status']).toBe(GeoOriginStatus.RESTRICTED);
  });

  it('should not set any headers when no geo headers are present', async () => {
    const response = await request(app)
      .get('/test');

    expect(response.headers['geo-origin-country']).toBeUndefined();
    expect(response.headers['geo-origin-region']).toBeUndefined();
    expect(response.headers['geo-origin-status']).toBeUndefined();
  });

  it('reflect empty string when headers are set but empty', async () => {
    const response = await request(app)
      .get('/test')
      .set('Geo-Origin-Country', '')
      .set('Geo-Origin-Region', '')
      .set('Geo-Origin-Status', '');

    expect(response.headers['geo-origin-country']).toBe('');
    expect(response.headers['geo-origin-region']).toBe('');
    expect(response.headers['geo-origin-status']).toBe('');
  });

  it('handle mixed empty and non-empty headers', async () => {
    const response = await request(app)
      .get('/test')
      .set('Geo-Origin-Country', 'GB')
      .set('Geo-Origin-Region', '')
      .set('Geo-Origin-Status', GeoOriginStatus.RESTRICTED);

    expect(response.headers['geo-origin-country']).toBe('GB');
    expect(response.headers['geo-origin-region']).toBe('');
    expect(response.headers['geo-origin-status']).toBe(GeoOriginStatus.RESTRICTED);
  });

  it('handle case-sensitive header names correctly', async () => {
    const response = await request(app)
      .get('/test')
      .set('geo-origin-country', 'DE') // lowercase
      .set('GEO-ORIGIN-REGION', 'BE') // uppercase
      .set('Geo-Origin-Status', GeoOriginStatus.OK); // mixed case

    // Express normalizes header names, so these should all work
    expect(response.headers['geo-origin-country']).toBe('DE');
    expect(response.headers['geo-origin-region']).toBe('BE');
    expect(response.headers['geo-origin-status']).toBe(GeoOriginStatus.OK);
  });
});
