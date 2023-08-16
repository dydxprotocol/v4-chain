/* eslint-disable no-console */

import util from 'util';

import bugsnag from '@bugsnag/js';
import Transport from 'winston-transport';

import config from './config';
import { safeJsonStringify } from './sanitization';
import { InfoObject } from './types';

/**
 * A Winston transport to log the stack traces of errors during development.
 */
export class StackTransport extends Transport {
  log(info: InfoObject, callback: () => void) {
    setImmediate(() => {
      if (config.NODE_ENV === 'test') {
        return;
      }
      if (info.error) {
        const { error } = info;
        // Note that util.inspect() is able to handle circular objects.
        const errorString = error && (error.stack || util.inspect(error));
        console.error(errorString);
      }
    });
    if (callback) {
      callback();
    }
  }
}

/**
 * A Winston transport which reports an error to Bugsnag.
 */
export class BugsnagTransport extends Transport {
  log(info: InfoObject, callback: () => void) {
    setImmediate(() => {
      if (bugsnag.isStarted() && config.SEND_BUGSNAG_ERRORS) {
        // Add service name to be able to distinguish which service the error came from
        const bugsnagInfo: InfoObject = {
          ...info,
          serviceName: config.SERVICE_NAME,
        };
        bugsnag.notify(
          new Error(
            safeJsonStringify(
              (info.error && (info.error.stack || util.inspect(info.error))) ||
              info.message ||
              {},
            ),
          ),
          (event) => {
            // eslint-disable-next-line no-param-reassign
            event.severity = 'error';
            event.addMetadata('info', bugsnagInfo);
            // eslint-disable-next-line no-param-reassign
            event.groupingHash = `${info.message}  ${info.at} ${config.SERVICE_NAME}`;
          },
        );
      }
    });
    if (callback) {
      callback();
    }
  }
}
