/**
 * Add tranports to the Winston logger.
 *
 * Do this outside logger.ts to avoid a dependency cycle with logger transports that may trigger
 * additional logging (see wrapBackgroundTask() in PagerDutyTransport).
 */

import winston from 'winston';
import { TransportStreamOptions } from 'winston-transport';

import config from './config';
import logger from './logger';
import {
  BugsnagTransport,
  StackTransport,
} from './logger-transports';
import { safeJsonStringify } from './sanitization';

/**
 * A Winston formatter to pretty-print logs in development and test.
 */
const alignedWithColorsAndTime = winston.format.combine(
  winston.format.colorize(),
  winston.format.timestamp(),
  winston.format.printf((info) => {
    const {
      level,
      ...args
    } = info;
    const ts = new Date().toISOString().slice(0, 19).replace('T', ' ');
    const argsString = safeJsonStringify(args, 2);
    return `${ts} [${level}]: ${argsString}`;
  }),
);

export function addTransportsToLogger(): void {
  // Send errors to Bugsnag (won't actually send unless BUGSNAG_KEY is set to a valid key).
  logger.add(
    new BugsnagTransport({
      level: 'error',

      // Disable since the Bugsnag client already reports and logs unhandled errors.
      handleExceptions: false,
      handleRejections: false,
    } as TransportStreamOptions),
  );

  // Send stack traces of any errors to the console.
  if (config.isTest() || config.isDevelopment()) {
    logger.add(
      new StackTransport({
        level: 'error',

        // Disable since the Bugsnag client already reports and logs unhandled errors.
        handleExceptions: false,
        handleRejections: false,
      } as TransportStreamOptions),
    );
  }

  // Send all logs to the console.
  if (!config.isTest() || config.ENABLE_LOGS_IN_TEST) {
    logger.add(
      new winston.transports.Console({
        level: config.LOG_LEVEL,
        format: (config.isTest() || config.isDevelopment()) ? alignedWithColorsAndTime : undefined,

        // Disable in development and test since the output is too verbose and mostly redundant with
        // the stack logged by the Bugsnag client.
        handleExceptions: config.isProduction() || config.isStaging(),
        handleRejections: config.isProduction() || config.isStaging(),
      } as TransportStreamOptions),
    );
  }
}
