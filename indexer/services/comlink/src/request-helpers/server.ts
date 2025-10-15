import bodyParser from 'body-parser';
import cors from 'cors';
import express, { Express } from 'express';
import requestId from 'express-request-id';
import responseTime from 'response-time';
import swaggerUi from 'swagger-ui-express';

import * as swaggerJson from '../../public/swagger.json';
import config from '../config';
import { logErrors } from './error-handler';
import geoHeadersMiddleware from './geo-headers-middleware';
import RequestLogger from './request-logger';
import resBodyCapture from './res-body-capture';

export default function server(
  indexV4?: express.Router,
): Express {
  const app: Express = express();

  app.use(geoHeadersMiddleware);

  app.use(responseTime({ suffix: false }));

  app.use(requestId());

  app.use(resBodyCapture);

  const corsOptions = {
    origin: config.CORS_ORIGIN,
    optionsSuccessStatus: 200,
  };

  app.use(cors(corsOptions));

  app.get('/health', (_req: express.Request, res: express.Response) => {
    res.json({ ok: true });
  });

  app.use((_req: express.Request, _res: express.Response, next: express.NextFunction) => next());

  app.use(bodyParser.json());

  app.use(RequestLogger);

  if (indexV4) {
    app.use('/v4', indexV4);
  }

  app.use('/docs', swaggerUi.serve, swaggerUi.setup(swaggerJson));

  // Log all other errors before being passed to default error handler
  app.use(logErrors);

  app.use((_req: express.Request, res: express.Response) => {
    res.status(404).json({
      error: 'Not Found',
      errorCode: 404,
    });
  });

  return app;
}
