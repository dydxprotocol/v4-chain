import {
  perpetualMarketRefresher,
  MAX_PARENT_SUBACCOUNTS,
  CHILD_SUBACCOUNT_MULTIPLIER,
} from '@dydxprotocol-indexer/postgres';
import { checkSchema, ParamSchema } from 'express-validator';

import config from '../../config';

export const CheckSubaccountSchema = checkSchema({
  address: {
    in: ['params', 'query'],
    isString: true,
  },
  subaccountNumber: {
    in: ['params', 'query'],
    isInt: {
      options: { gt: -1, lt: MAX_PARENT_SUBACCOUNTS * CHILD_SUBACCOUNT_MULTIPLIER + 1 },
    },
    errorMessage: 'subaccountNumber must be a non-negative integer less than 128001',
  },
});

export const CheckParentSubaccountSchema = checkSchema({
  address: {
    in: ['params', 'query'],
    isString: true,
  },
  parentSubaccountNumber: {
    in: ['params', 'query'],
    isInt: {
      options: { gt: -1, lt: MAX_PARENT_SUBACCOUNTS },
    },
    errorMessage: 'parentSubaccountNumber must be a non-negative integer less than 128',
  },
});

export const checkAddressSchemaRecord: Record<string, ParamSchema> = {
  address: {
    in: ['params'],
    isString: true,
  },
};

export const CheckAddressSchema = checkSchema(checkAddressSchemaRecord);

const limitSchemaRecord: Record<string, ParamSchema> = {
  limit: {
    in: ['query'],
    errorMessage: 'limit must be a positive integer that is not greater than max: ' +
      `${config.API_LIMIT_V4}`,
    customSanitizer: {
      options: (value?: number | string): number => {
        return value !== undefined ? +value : config.API_LIMIT_V4;
      },
    },
    custom: {
      options: (value: number) => {
        // Custom validator to ensure the value is a positive integer
        if (value <= 0 || value > config.API_LIMIT_V4 || !Number.isInteger(value)) {
          throw new Error(`limit must be a positive integer that is not greater than max: ${config.API_LIMIT_V4}`);
        }
        return true;
      },
    },
  },
};

const createdBeforeOrAtSchemaRecord: Record<string, ParamSchema> = {
  createdBeforeOrAtHeight: {
    in: ['query'],
    optional: true,
    isInt: {
      options: { gt: -1 },
    },
    errorMessage: 'createdBeforeOrAtHeight must be a non-negative integer',
  },
  createdBeforeOrAt: {
    in: ['query'],
    optional: true,
    isISO8601: true,
  },
};

const effectiveBeforeOrAtSchemaRecord: Record<string, ParamSchema> = {
  effectiveBeforeOrAtHeight: {
    in: ['query'],
    optional: true,
    isInt: {
      options: { gt: -1 },
    },
    errorMessage: 'effectiveBeforeOrAtHeight must be a non-negative integer',
  },
  effectiveBeforeOrAt: {
    in: ['query'],
    optional: true,
    isISO8601: true,
  },
};

const createdOnOrAfterSchemaRecord: Record<string, ParamSchema> = {
  createdOnOrAfterHeight: {
    in: ['query'],
    optional: true,
    isInt: {
      options: { gt: -1 },
    },
    errorMessage: 'createdOnOrAfterHeight must be a non-negative integer',
  },
  createdOnOrAfter: {
    in: ['query'],
    optional: true,
    isISO8601: true,
  },
};

export const CheckLimitSchema = checkSchema(limitSchemaRecord);

export const CheckLimitAndCreatedBeforeOrAtSchema = checkSchema({
  ...limitSchemaRecord,
  ...createdBeforeOrAtSchemaRecord,
});

export const CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema = checkSchema({
  ...limitSchemaRecord,
  ...createdBeforeOrAtSchemaRecord,
  ...createdOnOrAfterSchemaRecord,
});

export const CheckEffectiveBeforeOrAtSchema = checkSchema({
  ...effectiveBeforeOrAtSchemaRecord,
});

const checkTickerParamSchema: ParamSchema = {
  in: 'params',
  isString: true,
  custom: {
    options: perpetualMarketRefresher.isValidPerpetualMarketTicker,
    errorMessage: 'ticker must be a valid ticker (BTC-USD, etc)',
  },
};

const checkTickerOptionalQuerySchema: ParamSchema = {
  ...checkTickerParamSchema,
  in: 'query',
  optional: true,
};

export const CheckTickerParamSchema = checkSchema({
  ticker: checkTickerParamSchema,
});

export const CheckTickerOptionalQuerySchema = checkSchema({
  ticker: checkTickerOptionalQuerySchema,
});

export const CheckLimitAndYieldParamsSchema = checkSchema({
  ...limitSchemaRecord,
  createdAtOrBeforeHeight: {
    in: ['query'],
    optional: true,
    isString: true,
    errorMessage: 'createdOnOrAfterHeight must be a string',
  },
});
