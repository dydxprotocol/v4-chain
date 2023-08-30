import { logger, safeJsonStringify } from '@dydxprotocol-indexer/base';
import express from 'express';

import config from '../config';
import { ResponseWithBody } from '../types';

export default (
  request: express.Request,
  response: express.Response,
  next: express.NextFunction,
) => {
  response.on('finish', () => {
    const { protocol } : { protocol: string } = request;
    const host: string | undefined = request.get('host');
    const url: string = request.originalUrl;
    const fullUrl: string = `${protocol}://${host}${url}`;

    const isError: RegExpMatchArray | null = response.statusCode.toString().match(/^[^2]/);
    // Don't log GET requests unless configured to
    const shouldLogMethod: boolean = request.method !== 'GET' || config.LOG_GETS;
    if (shouldLogMethod || response.statusCode !== 200) {
      logger.info({
        at: 'requestLogger#logRequest',
        message: {
          request: {
            url: fullUrl,
            method: request.method,
            headers: request.headers,
            query: request.query,
            body: safeJsonStringify(request.body),
          },
          response: {
            statusCode: response.statusCode,
            errorBody: isError && (response as ResponseWithBody).body,
            statusMessage: response.statusMessage,
            headers: response.getHeaders(),
          },
        },
      });
    }
  });

  return next();
};
