/**
 * Logging configuration using Winston.
 */

import winston from 'winston';

import config from './config';
import { InfoObject } from './types';

// Fix types. The methods available depend on the levels used. We're using syslog levels, so these
// methods don't actually exist on our logger object.
type UnusedLevels = 'warn' | 'help' | 'data' | 'prompt' | 'http' | 'verbose' | 'input' | 'silly';

// Enforce type constraints on the objects passed into Winston logging functions.
interface LeveledLogMethod {
  (infoObject: InfoObject): winston.Logger,
}
// Exclude the functions whose type we want to change from the base definition. This seems to be
// enough (and the only way I've found) to trick TypeScript into accepting the modified LoggerExport
// as a valid extension of the base winston.Logger type.
type SyslogLevels = 'emerg' | 'alert' | 'crit' | 'error' | 'warning' | 'notice' | 'info' | 'debug';
export interface LoggerExport extends Omit<winston.Logger, UnusedLevels | SyslogLevels> {
  emerg: LeveledLogMethod,
  alert: LeveledLogMethod,
  crit: LeveledLogMethod,
  error: LeveledLogMethod,
  warning: LeveledLogMethod,
  notice: LeveledLogMethod,
  info: LeveledLogMethod,
  debug: LeveledLogMethod,
}

const logger: LoggerExport = winston.createLogger({
  levels: winston.config.syslog.levels,
  level: config.LOG_LEVEL,
  format: winston.format.combine(
    winston.format((info) => {
      return {
        ...info,           // info contains some symbols that are lost when the object is cloned.
        error: info.error,
      };
    })(),
    winston.format.json(),
  ),

  // Don't have Winston exit on uncaught errors. Bugsnag will exit when it is done handling them.
  exitOnError: false,
});

export default logger;
