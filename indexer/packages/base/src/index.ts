import { addTransportsToLogger } from './add-transports-to-logger';

export { baseConfigSchema, baseConfigSecrets } from './config';
export * from './config-util';
export * from './errors';
export { default as logger } from './logger';
export * from './types';
export * from './sanitization';
export { default as stats } from './stats';
export * from './tasks';
export * from './axios';
export * from './constants';
export * from './bugsnag';
export * from './stats-util';
export * from './date-helpers';
export * from './instance-id';
export * from './az-id';
export { cacheControlMiddleware, noCacheControlMiddleware } from './cache-control';

// Do this outside logger.ts to avoid a dependency cycle with logger transports that may trigger
// additional logging.
addTransportsToLogger();
