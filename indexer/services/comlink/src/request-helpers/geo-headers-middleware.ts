import express from 'express';

export default function geoHeadersMiddleware(
  req: express.Request,
  res: express.Response,
  next: express.NextFunction,
): void {
  const country = req.headers['geo-origin-country'];
  if (country !== undefined) {
    res.set('Geo-Origin-Country', country);
  }

  const region = req.headers['geo-origin-region'];
  if (region !== undefined) {
    res.set('Geo-Origin-Region', region);
  }

  const status = req.headers['geo-origin-status'];
  if (status !== undefined) {
    res.set('Geo-Origin-Status', status);
  }

  next();
}
