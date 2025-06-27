import { isValidLanguageCode } from '@dydxprotocol-indexer/notifications';
import {
  perpetualMarketRefresher,
  MAX_PARENT_SUBACCOUNTS,
  CHILD_SUBACCOUNT_MULTIPLIER,
} from '@dydxprotocol-indexer/postgres';
import { decode } from 'bech32';
import { body, checkSchema, ParamSchema } from 'express-validator';

import config from '../../config';

export const CheckSubaccountSchema = checkSchema({
  address: {
    in: ['params', 'query'],
    isString: true,
    custom: {
      options: isValidAddress,
    },
    errorMessage: 'address must be a valid dydx address',
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
    custom: {
      options: isValidAddress,
    },
    errorMessage: 'address must be a valid dydx address',
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
    custom: {
      options: isValidAddress,
    },
    errorMessage: 'address must be a valid address',
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

const paginationSchemaRecord: Record<string, ParamSchema> = {
  page: {
    in: ['query'],
    optional: true,
    isInt: {
      options: { gt: 0 },
    },
    errorMessage: 'page must be a non-negative integer',
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

const transferBetweenSchemaRecord: Record<string, ParamSchema> = {
  ...createdBeforeOrAtSchemaRecord,
  sourceAddress: {
    in: ['params', 'query'],
    isString: true,
  },
  sourceSubaccountNumber: {
    in: ['params', 'query'],
    isInt: {
      options: { gt: -1, lt: MAX_PARENT_SUBACCOUNTS * CHILD_SUBACCOUNT_MULTIPLIER + 1 },
    },
    errorMessage: 'subaccountNumber must be a non-negative integer less than 128001',
  },
  recipientAddress: {
    in: ['params', 'query'],
    isString: true,
  },
  recipientSubaccountNumber: {
    in: ['params', 'query'],
    isInt: {
      options: { gt: -1, lt: MAX_PARENT_SUBACCOUNTS * CHILD_SUBACCOUNT_MULTIPLIER + 1 },
    },
    errorMessage: 'subaccountNumber must be a non-negative integer less than 128001',
  },
};

export const CheckLimitSchema = checkSchema(limitSchemaRecord);

export const CheckPaginationSchema = checkSchema(paginationSchemaRecord);

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

const checkZeroPaymentsOptionalParamSchema: ParamSchema = {
  in: 'query',
  optional: true,
  isBoolean: true,
};

export const CheckZeroPaymentsOptionalParamSchema = checkSchema({
  showZeroPayments: checkZeroPaymentsOptionalParamSchema,
});

export const CheckTickerParamSchema = checkSchema({
  ticker: checkTickerParamSchema,
});

export const CheckTickerOptionalQuerySchema = checkSchema({
  ticker: checkTickerOptionalQuerySchema,
});

export const CheckHistoricalBlockTradingRewardsSchema = checkSchema({
  ...checkAddressSchemaRecord,
  ...limitSchemaRecord,
  startingBeforeOrAt: {
    in: ['query'],
    optional: true,
    isISO8601: true,
  },
  startingBeforeOrAtHeight: {
    in: ['query'],
    optional: true,
    isInt: {
      options: { gt: -1 },
    },
    errorMessage: 'startingBeforeOrAtHeight must be a non-negative integer',
  },
});

export const CheckTransferBetweenSchema = checkSchema(transferBetweenSchemaRecord);

export const RegisterTokenValidationSchema = [
  body('token')
    .exists().withMessage('Token is required')
    .isString()
    .withMessage('Token must be a string')
    .notEmpty()
    .withMessage('Token cannot be empty'),
  body('timestamp')
    .exists().withMessage('timestamp is required')
    .isNumeric()
    .withMessage('timestamp must be a number')
    .notEmpty()
    .withMessage('timestamp cannot be empty'),
  body('message')
    .exists().withMessage('message is required')
    .isString()
    .withMessage('message must be a string')
    .notEmpty()
    .withMessage('message cannot be empty'),
  body('signedMessage')
    .exists().withMessage('signedMessage is required')
    .isString()
    .withMessage('signedMessage must be a string')
    .notEmpty()
    .withMessage('signedMessage cannot be empty'),
  body('pubKey')
    .exists().withMessage('pubKey is required')
    .isString()
    .withMessage('pubKey must be a string')
    .notEmpty()
    .withMessage('pubKey cannot be empty'),
  body('walletIsKeplr')
    .exists().withMessage('walletIsKeplr is required')
    .isBoolean()
    .withMessage('walletIsKeplr must be a boolean')
    .notEmpty()
    .withMessage('walletIsKeplr cannot be empty'),
  body('language')
    .optional()
    .isString()
    .withMessage('Language must be a string')
    .custom((value: string) => {
      if (!isValidLanguageCode(value)) {
        throw new Error('Invalid language code');
      }
      return true;
    }),
];

function verifyIsBech32(address: string): Error | undefined {
  try {
    decode(address);
  } catch (error) {
    return error;
  }

  return undefined;
}

export function isValidDydxAddress(address: string): boolean {
  // An address is valid if it starts with `dydx1` and is Bech32 format.
  return address.startsWith('dydx1') && (verifyIsBech32(address) === undefined);
}

export function isValidAddress(address: string): boolean {
  // Address is valid if its under 90 characters and alphanumeric
  return address.length <= 90 && /^[a-zA-Z0-9]*$/.test(address);
}
