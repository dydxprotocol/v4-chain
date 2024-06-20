import { parentPort } from 'worker_threads';

import { getMessagesToForward } from '../../helpers/from-kafka-helpers';

parentPort?.on('message', (data) => {
  const { topic, message } = data;
  const messagesToForward = getMessagesToForward(topic, message);
  parentPort?.postMessage(messagesToForward);
});
