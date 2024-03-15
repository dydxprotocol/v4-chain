import { stats } from '@dydxprotocol-indexer/base';
import express from 'express';

import config from '../config';

export interface ResponseStatsOptions {
  controllerName: string,
}

export default (options?: ResponseStatsOptions) => (
  request: express.Request,
  response: express.Response,
  next: express.NextFunction,
) => {
  response.on('finish', () => {
    stats.increment(`${config.SERVICE_NAME}.${options?.controllerName}.response_status_code.${response.statusCode}`,
      1, { path: request.route.path, method: request.method });
    stats.increment(`${config.SERVICE_NAME}.${options?.controllerName}.response_status_code`,
      1, { path: request.route.path, method: request.method });
    stats.increment(`${config.SERVICE_NAME}.response_status_code.${response.statusCode}`,
      1, { path: request.route.path, method: request.method });
  });

  return next();
};
