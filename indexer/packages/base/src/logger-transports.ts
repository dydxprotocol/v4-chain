/* eslint-disable no-console */

import util from 'util';

import Transport from 'winston-transport';

import config from './config';
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
