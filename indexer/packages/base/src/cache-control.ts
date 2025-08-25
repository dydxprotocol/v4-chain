import { Request, Response, NextFunction } from 'express';

export function cacheControlMiddleware(directive: string) {
  return (req: Request, res: Response, next: NextFunction) => {
    res.setHeader('Cache-Control', directive);
    next();
  };
}

export const noCacheControlMiddleware = cacheControlMiddleware('no-cache, no-store, no-transform');
