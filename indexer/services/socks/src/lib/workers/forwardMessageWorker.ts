import {parentPort} from 'worker_threads';
import {MessageForwarder} from '../message-forwarder';

parentPort?.on('message', async (data) => {
  const { subscriptions, index, messageToForward } = data;
  await MessageForwarder.getInstance(subscriptions, index).forwardMessage(messageToForward);
  parentPort?.postMessage({ status: 'done' });
});
