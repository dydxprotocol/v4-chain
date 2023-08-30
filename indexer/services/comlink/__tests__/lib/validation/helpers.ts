import express from 'express';

import { CheckLimitAndCreatedBeforeOrAtSchema, CheckSubaccountSchema } from '../../../src/lib/validation/schemas';
import { handleValidationErrors } from '../../../src/request-helpers/error-handler';
import Server from '../../../src/request-helpers/server';

const router: express.Router = express.Router();

router.get(
  '/check-subaccount-schema',
  ...CheckSubaccountSchema,
  handleValidationErrors,
  (req: express.Request, res: express.Response) => {
    res.sendStatus(200);
  },
);

router.get(
  '/check-limit-and-created-before-schema',
  ...CheckLimitAndCreatedBeforeOrAtSchema,
  handleValidationErrors,
  (req: express.Request, res: express.Response) => {
    res.sendStatus(200);
  },
);

export const schemaTestApp = Server(router);
