/**
 * Utilities for parsing environment variables, with typings.
 *
 * The purpose of this module is to enforce types on all environment variables and throw an error if
 * required environment variables are missing. Services should not access `process.env` directly.
 *
 * Each service should include a config.ts file which should define a config schema, for example:
 *
 * ```
 * const configSchema = {
 *   NETWORK_ID: parseInteger(),
 *   NODE_ENV: parseString(),
 *   HISTORY_PAGE_MAX: parseNumber({ default: 100 }),
 * };
 * ```
 *
 * The config.ts file should then parse and export a config object, for example:
 *
 * ```
 * const config = parseSchema(configSchema);
 * export default config;
 * ```
 *
 * The exported config object should then have types on all its properties.
 */

import { Big } from 'big.js';
import { BigNumber } from 'bignumber.js';
import _ from 'lodash';

import { ConfigError } from './errors';
import {
  Bigable,
  BigIntable,
  BigNumberable,
  NodeEnv,
} from './types';

// A parse function takes the name of an environment variable as an argument,
// e.g. 'NODE_ENV', and parses and returns that variable from `process.env`.
type ParseFn<T> = (varName: string) => T;

// A schema base maps environment variable names to parse functions that can be
// used to parse those variables.
type SchemaBase = { [varName: string]: ParseFn<unknown> };

interface ParseOptions<T> {
  // If `default` is present, then the environment variable will be optional and will default to the
  // value of `default` when unset. In particular, `default` may be null in which case the config
  // value will be null when the environment variable is not set.
  default: T,

  // Can be specified to ensure the default value is not used when running in a certain NODE_ENV.
  requireInEnv?: NodeEnv[],
}

const NODE_ENV = process.env.NODE_ENV;

function defaultIsValid(
  options?: ParseOptions<unknown>,
): options is ParseOptions<unknown> {
  if (!options) {
    return false;
  }
  const hasDefaultValue = typeof options.default !== 'undefined';
  const requiredInEnv = (
    NODE_ENV &&
    options.requireInEnv &&
    options.requireInEnv.includes(NODE_ENV as NodeEnv)
  );
  return (hasDefaultValue && !requiredInEnv);
}

export function parseString(): ParseFn<string>;
export function parseString(options: ParseOptions<string>): ParseFn<string>;
export function parseString(options: ParseOptions<null>): ParseFn<string | null>;
export function parseString(options?: ParseOptions<string | null>): ParseFn<string | null> {
  return (varName: string) => {
    const value = process.env[varName];
    if (!value) {
      if (defaultIsValid(options)) {
        return options.default;
      }
      throw new ConfigError(`Missing required env var '${varName}' (string)`);
    }
    return value;
  };
}

export function parseBoolean(): ParseFn<boolean>;
export function parseBoolean(options: ParseOptions<boolean>): ParseFn<boolean>;
export function parseBoolean(options: ParseOptions<null>): ParseFn<boolean | null>;
export function parseBoolean(options?: ParseOptions<boolean | null>): ParseFn<boolean | null> {
  return (varName: string) => {
    const rawValue = process.env[varName];
    if (!rawValue) {
      if (defaultIsValid(options)) {
        return options.default;
      }
      throw new ConfigError(`Missing required env var '${varName}' (number)`);
    }
    if (rawValue === 'true') {
      return true;
    }
    if (rawValue === 'false') {
      return false;
    }
    throw new ConfigError(`Invalid boolean for env var '${varName}'`);
  };
}

export function parseNumber(): ParseFn<number>;
export function parseNumber(options: ParseOptions<number>): ParseFn<number>;
export function parseNumber(options: ParseOptions<null>): ParseFn<number | null>;
export function parseNumber(options?: ParseOptions<number | null>): ParseFn<number | null> {
  return (varName: string) => {
    const rawValue = process.env[varName];
    if (!rawValue) {
      if (defaultIsValid(options)) {
        return options.default;
      }
      throw new ConfigError(`Missing required env var '${varName}' (number)`);
    }
    const value = Number(rawValue);
    if (Number.isNaN(value)) {
      throw new ConfigError(`Invalid number for env var '${varName}'`);
    }
    return value;
  };
}

export function parseInteger(): ParseFn<number>;
export function parseInteger(options: ParseOptions<number>): ParseFn<number>;
export function parseInteger(options: ParseOptions<null>): ParseFn<number | null>;
export function parseInteger(options?: ParseOptions<number | null>): ParseFn<number | null> {
  return (varName: string) => {
    const rawValue = process.env[varName];
    if (!rawValue) {
      if (defaultIsValid(options)) {
        if (options.default !== null && !Number.isInteger(options.default)) {
          throw new ConfigError(`Expected integer default value for env var '${varName}'`);
        }
        return options.default;
      }
      throw new ConfigError(`Missing required env var '${varName}' (integer)`);
    }
    const value = Number(rawValue);
    if (!Number.isInteger(value)) {
      throw new ConfigError(`Invalid integer for env var '${varName}'`);
    }
    return value;
  };
}

export function parseBigInt(): ParseFn<bigint>;
export function parseBigInt(options: ParseOptions<BigIntable>): ParseFn<bigint>;
export function parseBigInt(options: ParseOptions<null>): ParseFn<bigint | null>;
export function parseBigInt(options?: ParseOptions<BigIntable | null>): ParseFn<bigint | null> {
  return (varName: string) => {
    const rawValue = process.env[varName];
    if (!rawValue) {
      if (defaultIsValid(options)) {
        try {
          return (options.default === null)
            ? null
            : BigInt(options.default);
        } catch (e) {
          throw new ConfigError(`Expected integer default value for env var '${varName}'`);
        }
      }
      throw new ConfigError(`Missing required env var '${varName}' (BigInt)`);
    }
    try {
      return BigInt(rawValue);
    } catch (e) {
      throw new ConfigError(`Invalid BigInt for env var '${varName}'`);
    }
  };
}

export function parseBN(): ParseFn<BigNumber>;
export function parseBN(options: ParseOptions<BigNumberable>): ParseFn<BigNumber>;
export function parseBN(options: ParseOptions<null>): ParseFn<BigNumber | null>;
export function parseBN(options?: ParseOptions<BigNumberable | null>): ParseFn<BigNumber | null> {
  return (varName: string) => {
    const rawValue = process.env[varName];
    if (!rawValue) {
      if (defaultIsValid(options)) {
        return options.default === null ? null : new BigNumber(options.default);
      }
      throw new ConfigError(`Missing required env var '${varName}' (BigNumber)`);
    }
    const value = new BigNumber(rawValue);
    if (value.isNaN()) {
      throw new ConfigError(`Invalid BigNumber for env var '${varName}'`);
    }
    return value;
  };
}

export function parseBig(): ParseFn<Big>;
export function parseBig(options: ParseOptions<Bigable>): ParseFn<Big>;
export function parseBig(options: ParseOptions<null>): ParseFn<Big | null>;
export function parseBig(options?: ParseOptions<Bigable | null>): ParseFn<Big | null> {
  return (varName: string) => {
    const rawValue = process.env[varName];
    if (!rawValue) {
      if (defaultIsValid(options)) {
        return options.default === null ? null : new Big(options.default);
      }
      throw new ConfigError(`Missing required env var '${varName}' (Big)`);
    }
    try {
      return new Big(rawValue);
    } catch (error) {
      throw new ConfigError(`Invalid Big for env var '${varName}'`);
    }
  };
}

/**
 * Process the schema and parse environment variables.
 *
 * Use type inference to presserve type information including which values may or may not be null.
 */
export function parseSchema<T extends SchemaBase>(
  schema: T,
  { prefix }: { prefix?: string } = {},
): {
  [K in keyof T]: T[K] extends ParseFn<infer U> ? U : never;
} & {
  isDevelopment: () => boolean,
  isStaging: () => boolean,
  isProduction: () => boolean,
  isTest: () => boolean,
} {
  const config = _.mapValues(schema, (parseFn: ParseFn<T>, varName: string) => {
    const fullVarName = prefix ? `${prefix}_${varName}` : varName;
    return parseFn(fullVarName);
  }) as { [K in keyof T]: T[K] extends ParseFn<infer U> ? U : never };

  // Include helper functions.
  return {
    ...config,
    isDevelopment: () => NODE_ENV === NodeEnv.DEVELOPMENT,
    isStaging: () => NODE_ENV === NodeEnv.STAGING,
    isProduction: () => NODE_ENV === NodeEnv.PRODUCTION,
    isTest: () => NODE_ENV === NodeEnv.TEST,
  };
}
