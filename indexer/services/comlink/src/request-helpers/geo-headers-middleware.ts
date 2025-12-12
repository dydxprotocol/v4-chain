import express from 'express';

export default function geoHeadersMiddleware(
  req: express.Request,
  res: express.Response,
  next: express.NextFunction,
): void {
  const country = req.headers['geo-origin-country'];
  if (country !== undefined) {
    res.header('Access-Control-Expose-Headers', 'geo-origin-country');
    res.set('Geo-Origin-Country', country);
  }

  const region = req.headers['geo-origin-region'];
  if (region !== undefined) {
    res.header('Access-Control-Expose-Headers', 'geo-origin-region');
    res.set('Geo-Origin-Region', region);
  }

  const status = req.headers['geo-origin-status'];
  if (status !== undefined) {
    res.header('Access-Control-Expose-Headers', 'geo-origin-status');
    res.set('Geo-Origin-Status', status);
  }

  next();
}
