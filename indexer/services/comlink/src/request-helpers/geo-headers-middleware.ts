import express from 'express';

export default function geoHeadersMiddleware(
  req: express.Request,
  res: express.Response,
  next: express.NextFunction,
): void {
  const country = req.headers['geo-origin-country'];
  const exposedHeaders = [];

  if (country !== undefined) {
    exposedHeaders.push('geo-origin-country');
    res.set('Geo-Origin-Country', country);
  }

  const region = req.headers['geo-origin-region'];
  if (region !== undefined) {
    exposedHeaders.push('geo-origin-region');
    res.set('Geo-Origin-Region', region);
  }

  const status = req.headers['geo-origin-status'];
  if (status !== undefined) {
    exposedHeaders.push('geo-origin-status');
    res.set('Geo-Origin-Status', status);
  }

  if (exposedHeaders.length > 0) {
    res.header('Access-Control-Expose-Headers', exposedHeaders.join(', '));
  }

  next();
}
