import express from 'express';

export default function geoOriginHeadersMiddleware(
  req: express.Request,
  res: express.Response,
  next: express.NextFunction,
): void {
  const geoCountry = req.headers['geo-origin-country'];
  if (geoCountry !== undefined) {
    res.set('Geo-Origin-Country', geoCountry);
  }

  const geoRegion = req.headers['geo-origin-region'];
  if (geoRegion !== undefined) {
    res.set('Geo-Origin-Region', geoRegion);
  }

  const geoStatus = req.headers['geo-origin-status'];
  if (geoStatus !== undefined) {
    res.set('Geo-Origin-Status', geoStatus);
  }

  next();
}
