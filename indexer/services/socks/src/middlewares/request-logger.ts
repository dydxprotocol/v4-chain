import { logger } from '@dydxprotocol-indexer/base';
import express from 'express';
import { ResponseWithBody } from 'src/types';

export default (
  request: express.Request,
  response: express.Response,
  next: express.NextFunction,
) => {
  response.on('finish', () => {
    const protocol: string = request.protocol;
    const host: string | undefined = request.get('host');
    const url: string = request.originalUrl;
    const fullUrl: string = `${protocol}://${host}${url}`;

    // Convert RegExpMatchArray | null into true/false (boolean).
    const isError: boolean = !!response.statusCode.toString().match(/^[^2]/);
    if (request.method !== 'GET') {
      logger.info({
        at: 'requestLogger#logRequest',
        message: {
          request: {
            url: fullUrl,
            method: request.method,
            headers: request.headers,
            query: request.query,
            body: request.body,
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
