import util from 'util';

import bugsnag from '@bugsnag/js';
import { NodeConfig } from '@bugsnag/node';
import bugsnagExpress from '@bugsnag/plugin-express';

import config from './config';
import { BugsnagReleaseStage } from './types';

/**
 * Configuration options for the Bugsnag client.
 *
 * The Bugsnag client handles unhandled exceptions and unhandled promise rejections by default.
 * This occurs separately from BugsnagTransport.
 *
 * Bugsnag's default onUncaughtException option calls process.exit(1) when an unhandled exception
 * occurs. It is recommended to do the same for unhandled promise rejections, so we replicate the
 * behavior of onUncaughtException here in onUnhandledRejection.
 */
const OPTIONS: NodeConfig = {
  apiKey: config.BUGSNAG_KEY,
  logger: config.isDevelopment() || config.isTest() ? {
    debug: () => null,
    info: logBugsnag,
    warn: logBugsnag,
    error: logBugsnag,
  } : null,
  releaseStage: config.BUGSNAG_RELEASE_STAGE || BugsnagReleaseStage.DEVELOPMENT,
  plugins: [bugsnagExpress],
  onUnhandledRejection: (error, _event, logger) => {
    const errorString = error && error.stack ? error.stack : util.inspect(error);
    logger.error(`Unhandled rejection…\n${errorString}`);
    process.exit(1);
  },
};

/**
 * Initialize the Bugsnag client.
 */
export function startBugsnag(): void {
  bugsnag.start(OPTIONS);
}

/**
 * Bugsnag logging in development and test.
 *
 * This is used, for example, when an unhandled exception or unhandled promise rejection occurs.
 */
function logBugsnag(info: {}) {
  // eslint-disable-next-line no-console
  console.log('Bugsnag:', info);
}
