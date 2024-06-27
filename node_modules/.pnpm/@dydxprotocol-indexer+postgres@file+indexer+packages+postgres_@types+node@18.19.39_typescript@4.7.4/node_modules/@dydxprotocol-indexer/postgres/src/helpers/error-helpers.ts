import { logger } from '@dydxprotocol-indexer/base';

import { ValidationError } from '../lib/errors';

export function logAndThrowValidationError(message: string): void {
  logger.error({
    at: 'errorHelpers#logAndThrowValidationError',
    message,
  });
  throw new ValidationError(message);
}
