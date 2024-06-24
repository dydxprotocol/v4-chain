import { parentPort } from 'worker_threads';

import getMessagesToForward from './from-kafka-helpers';

parentPort?.on('message', (data) => {
  const { topic, message, clobPairIdToTickerMap } = data;
  const messagesToForward = getMessagesToForward({ topic, message, clobPairIdToTickerMap });
  parentPort?.postMessage(messagesToForward);
});
