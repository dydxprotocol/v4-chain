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
  // Send stack traces of any errors to the console.
  if (config.isTest() || config.isDevelopment()) {
    logger.add(
      new StackTransport({
        level: 'error',
        handleExceptions: true,
        handleRejections: true,
      } as TransportStreamOptions),
    );
  }

  // Send all logs to the console.
  if (!config.isTest() || config.ENABLE_LOGS_IN_TEST) {
    logger.add(
      new winston.transports.Console({
        level: config.LOG_LEVEL,
        format: (config.isTest() || config.isDevelopment()) ? alignedWithColorsAndTime : undefined,

        handleExceptions: config.isProduction() || config.isStaging(),
        handleRejections: config.isProduction() || config.isStaging(),
      } as TransportStreamOptions),
    );
  }
}
