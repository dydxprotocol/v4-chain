import { StatsD } from 'hot-shots';

import config from './config';
import logger from './logger';

const stats = new StatsD({
  // The host to send stats to
  // default: localhost
  host: config.STATSD_HOST,

  // The port to send stats to
  // default: 8125
  port: config.STATSD_PORT,

  // Expose this StatsD instance globally?
  // default: false
  globalize: false,

  // Cache the initial dns lookup to host
  // default: false
  /* cacheDns: false, */

  // Create mock instance, sending no stats to the server. Allows data to be read from mockBuffer
  // default: false
  mock: config.isTest(),

  // Tags that will be added to every metric. Can be either an object or list of tags
  // default: {}
  globalTags: [
    `NODE_ENV:${config.NODE_ENV}`,
    `SERVICE_NAME:${config.SERVICE_NAME}`,
  ],
  // If larger than 0, metrics will be buffered and only sent data length greater than the size
  // default: 0
  /* maxBufferSize: 0, */

  // If buffering is in use, this is the time in ms to always flush any buffered metrics
  // default: 1000
  /* bufferFlushInterval: 1000, */

  // Use Telegrafs StatsD line protocol, which is slightly different than the rest
  // default: false
  /* telegraf: false, */

  // Sends only a sample of data to StatsD for all StatsD methods.
  // Can be overriden at the method level.
  // default: 1
  /* sampleRate: 1, */

  // A function with one argument. It is called to handle various errors
  // default: none, errors are thrown/logger to console
  /* errorHandler: */

  // Use the default interface on a Linux system. Useful when running in containers
  /* useDefaultRoute: */

  // Use tcp option for TCP protocol. Defaults to UDP otherwise
  /* protocol: */
});

stats.socket.on('error', (error) => {
  logger.error({
    at: 'stats#onError',
    message: error.message,
    error,
  });
});

export default stats;
