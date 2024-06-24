import {logger} from '@dydxprotocol-indexer/base';

export default function print(): void {
  logger.info({
    at: 'blank-worker#print',
    message: 'Hello, World!',
  });
}
