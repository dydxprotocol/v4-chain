import { Big } from 'big.js';
import { BigNumber } from 'bignumber.js';

export type BigNumberable = BigNumber | string | number;

export type Bigable = Big | string | number;

export type BigIntable = bigint | string | number;

export enum NodeEnv {
  DEVELOPMENT = 'development',
  PRODUCTION = 'production',
  STAGING = 'staging',
  TEST = 'test',
}

export enum BugsnagReleaseStage {
  DEVELOPMENT = 'development',
  PRODUCTION = 'production',
  STAGING = 'staging',
}

export interface PagerDutyInfo {
  message: {},
  id?: string,
}

// Enforce type constraints on the objects passed into Winston logging functions.
export interface InfoObject extends PagerDutyInfo {
  [key: string]: unknown,
  // Note: If message were missing, the info object would get wrapped as { message: infoObject },
  // which is not what we want since it can prevent errors from being reported as expected.
  message: {},
  at: string,
  // Require `error` to be the right type.
  error?: Error,
}
