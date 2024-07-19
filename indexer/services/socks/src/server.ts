import bodyParser from 'body-parser';
import cors from 'cors';
import express, { Express } from 'express';
import expressRequestId from 'express-request-id';
import nocache from 'nocache';
import responseTime from 'response-time';

import RequestLogger from './middlewares/request-logger';
import resBodyCapture from './middlewares/res-body-capture';

export default function server(): Express {
  const app = express();

  app.use(responseTime({ suffix: false }));

  app.use(expressRequestId());

  app.use(resBodyCapture);

  const corsOptions = {
    origin: process.env.CORS_ORIGIN,
    optionsSuccessStatus: 200,
  };

  app.use(cors(corsOptions));

  app.use(nocache());

  app.get('/health', (_req: express.Request, res: express.Response) => {
    return res.status(200).json({ ok: true });
  });

  app.use((_req: express.Request, _res: express.Response, next: express.NextFunction) => next());

  app.use(bodyParser.json());

  app.use(RequestLogger);

  app.use((_req: express.Request, res: express.Response) => {
    res.status(404).json({
      error: 'Not Found',
      errorCode: 404,
    });
  });

  return app;
}
