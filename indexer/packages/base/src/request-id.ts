import { randomUUID } from 'node:crypto';

import {
  NextFunction,
  Request,
  RequestHandler,
  Response,
} from 'express';

declare module 'express-serve-static-core' {
  interface Request {
    id: string,
  }
}

declare module 'http' {
  interface IncomingHttpHeaders {
    'X-Request-Id'?: string,
  }
}

const headerKey: 'X-Request-Id' = 'X-Request-Id';

export function requestId(): RequestHandler {
  return function requestIdHandler(req: Request, res: Response, next: NextFunction) {
    // get ignores case
    req.id = req.get(headerKey) || randomUUID();
    res.setHeader(headerKey, req.id);
    next();
  };
}
