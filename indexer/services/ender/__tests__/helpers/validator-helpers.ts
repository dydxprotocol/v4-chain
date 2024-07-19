import { logger } from '@dydxprotocol-indexer/base';

import { defaultHeight } from './constants';

export function expectLoggedParseMessageError(
  className: string,
  message: string,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  additionalParameters: any,
): void {
  expect(logger.error).toHaveBeenCalledWith({
    at: `${className}#logAndThrowParseMessageError`,
    message,
    blockHeight: defaultHeight,
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    ...additionalParameters,
  });
}

export function expectDidntLogError(): void {
  expect(logger.error).not.toHaveBeenCalled();
}
