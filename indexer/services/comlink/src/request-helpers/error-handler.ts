import { logger } from '@dydxprotocol-indexer/base';
import express from 'express';
import { validationResult, matchedData } from 'express-validator';
import { isEqual, isObject, forEach } from 'lodash';

export function handleUnexpectedFieldErrors(
  req: express.Request,
  res: express.Response,
  next: express.NextFunction,
) {
  const errors: { msg: string }[] = [];

  /**
   * @description Recursively find keys in `object` that are missing or not equal in `base`.
   */
  function difference(object: {}, base: Record<string, {}>, parentPath: string = '') {
    forEach(object, (value, key) => {
      if (!isEqual(value, base[key])) {
        if (isObject(value) && isObject(base[key])) {
          difference(value, base[key] as Record<string, {}>, `${parentPath}${key}.`);
        } else {
          errors.push({ msg: `Unexpected field: '${parentPath}${key}'` });
        }
      }
    });
  }

  const body = matchedData(req);
  difference(req.body, body);
  if (errors.length === undefined || errors.length === 0) {
    return next();
  }

  return res.status(400).json({ errors });
}

export function handleValidationErrors(
  req: express.Request,
  res: express.Response,
  next: express.NextFunction,
) {
  const errors = validationResult(req);
  if (!errors.isEmpty()) {
    return res.status(400).json({ errors: errors.array() });
  }
  return next();
}

export function logErrors(
  error: Error,
  req: express.Request,
  res: express.Response,
  next: express.NextFunction,
) {
  logger.error({
    at: 'error-handerl#logErrors',
    message: `Encountered error: ${error}`,
    error,
  });
  next(error);
}
