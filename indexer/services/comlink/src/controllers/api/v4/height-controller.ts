import { stats } from '@dydxprotocol-indexer/base';
import { BlockTable, BlockFromDatabase } from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { Controller, Get, Route } from 'tsoa';

import { defaultRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { NotFoundError } from '../../../lib/errors';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { HeightResponse } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'height-controller';

@Route('height')
class HeightController extends Controller {
  @Get('/')
  async getHeight(): Promise<HeightResponse> {
    try {
      const latestBlock: BlockFromDatabase = await BlockTable.getLatest();
      return {
        height: latestBlock.blockHeight,
        time: latestBlock.time,
      };
    } catch {
      throw new NotFoundError('No blocks found');
    }
  }
}

router.get(
  '/',
  rateLimiterMiddleware(defaultRateLimiter),
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    try {
      const controller: HeightController = new HeightController();
      const response: HeightResponse = await controller.getHeight();

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'HeightController GET /',
        'Height error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_latest_block_height.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
