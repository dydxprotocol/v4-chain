import {
  logger, setInstanceId, getInstanceId, startBugsnag, wrapBackgroundTask,
} from '@dydxprotocol-indexer/base';
import { startConsumer } from '@dydxprotocol-indexer/kafka';
import { blockHeightRefresher, perpetualMarketRefresher } from '@dydxprotocol-indexer/postgres';

import config from './config';
import {
  connect as connectToKafka,
  disconnect as disconnectFromKafka,
} from './helpers/kafka/kafka-controller';
import { Wss } from './helpers/wss';
import { MessageForwarder } from './lib/message-forwarder';
import { Subscriptions } from './lib/subscription';
import Server from './server';
import { Index } from './websocket';

let index: Index;
let messageForwarder: MessageForwarder;
let subscriptions: Subscriptions;

const shutdown = async (signal: string) => {
  logger.info({
    at: `index#${signal}`,
    message: `Received ${signal.toUpperCase()}, shutting down`,
  });

  if (index !== undefined) {
    await index.close();
  }
  if (messageForwarder !== undefined) {
    messageForwarder.stop();
  }
  if (subscriptions !== undefined) {
    subscriptions.stop();
  }
  await disconnectFromKafka();

  process.exit(0);
};

async function start(): Promise<void> {
  logger.info({
    at: 'index#start',
    message: `Starting in env ${config.NODE_ENV}`,
  });

  startBugsnag();

  logger.info({
    at: 'index#start',
    message: 'Getting instance id...',
  });

  await setInstanceId();

  logger.info({
    at: 'index#start',
    message: `Got instance id ${getInstanceId()}.`,
  });

  // Initialize PerpetualMarkets and BlockHeight cache
  await Promise.all([
    blockHeightRefresher.updateBlockHeight(),
    perpetualMarketRefresher.updatePerpetualMarkets(),
  ]);
  wrapBackgroundTask(blockHeightRefresher.start(), true, 'startUpdateBlockHeight');
  wrapBackgroundTask(perpetualMarketRefresher.start(), true, 'startUpdatePerpetualMarkets');

  logger.info({
    at: 'index#start',
    message: 'Started task loops to refresh perpetual markets, Starting websockets...',
  });

  const wss = new Wss();
  await wss.start();

  logger.info({
    at: 'index#start',
    message: 'Started websockets. Subscribing to kafka...',
  });

  await connectToKafka();

  subscriptions = new Subscriptions();
  index = new Index(wss, subscriptions);
  messageForwarder = new MessageForwarder(subscriptions, index);
  subscriptions.start(messageForwarder.forwardToClient);
  messageForwarder.start();
  await startConsumer(config.BATCH_PROCESSING_ENABLED);

  logger.info({
    at: 'index#start',
    message: 'Connected to kafka.',
  });

  startServer();

  logger.info({
    at: 'index',
    message: 'Successfully started',
  });
}

function startServer(): void {
  const app = Server();
  const port: string = config.PORT;
  app.listen(port, () => {
    logger.info({
      at: 'index#startServer',
      message: `Api server is listening on ${port}`,
    });
  });
}

wrapBackgroundTask(start(), true, 'main');

process.on('SIGINT', () => shutdown('SIGINT'));
process.on('SIGTERM', () => shutdown('SIGTERM'));
