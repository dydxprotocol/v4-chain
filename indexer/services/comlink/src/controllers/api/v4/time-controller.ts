import express from 'express';
import { DateTime } from 'luxon';
import { Controller, Get, Route } from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { TimeResponse } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'time-controller';

@Route('time')
class TimeController extends Controller {
  @Get('/')
  getTime(): TimeResponse {
    const time: DateTime = DateTime.utc();

    return {
      iso: time.toISO(),
      epoch: time.toSeconds(),
    };
  }
}

router.get(
  '/',
  rateLimiterMiddleware(getReqRateLimiter),
  ExportResponseCodeStats({ controllerName }),
  (_req: express.Request, res: express.Response) => {
    const controller: TimeController = new TimeController();
    return res.send(controller.getTime());
  },
);

export default router;
