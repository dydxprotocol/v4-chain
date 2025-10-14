import { isValidLanguageCode } from '@dydxprotocol-indexer/notifications';
import {
  perpetualMarketRefresher,
  MAX_PARENT_SUBACCOUNTS,
  CHILD_SUBACCOUNT_MULTIPLIER,
} from '@dydxprotocol-indexer/postgres';
import { decode } from 'bech32';
import { body, checkSchema, ParamSchema } from 'express-validator';

import config from '../../config';
import { SigninMethod } from '../../types';

const addressSchema = {
  isString: true as const,
  custom: {
    options: isValidAddress,
  },
  errorMessage: 'address must be a valid dydx address',
};

export const CheckSubaccountSchema = checkSchema({
  address: {
    in: ['params', 'query'],
    ...addressSchema,
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
    ...addressSchema,
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

const checkDailyOptionalParamSchema: ParamSchema = {
  in: 'query',
  optional: true,
  isBoolean: true,
  toBoolean: true,
};

const checkBridgeSchema: Record<string, ParamSchema> = {
  // Validate the event object structure
  event: {
    in: 'body',
    isObject: true,
    errorMessage: 'Event must be an object',
  },
  // for solana
  'event.transaction': {
    in: 'body',
    optional: true,
    isArray: true,
    errorMessage: 'Event.transaction must be an array',
  },
  'event.transaction.*.meta': {
    in: 'body',
    optional: true,
    isArray: true,
    errorMessage: 'Event.transaction.transaction must be an array',
  },
  'event.transaction.*.meta.*.post_token_balances': {
    in: 'body',
    optional: true,
    isArray: true,
    errorMessage: 'Event.transaction.meta.post_token_balances must be an array',
  },
  'event.transaction.*.meta.*.post_token_balances.*.amount': {
    in: 'body',
    optional: true,
    isString: true,
    errorMessage: 'Event.transaction.meta.post_token_balances.amount must be a string',
  },
  // for evm
  'event.activity': {
    in: 'body',
    optional: true,
    isArray: true,
    errorMessage: 'Event.activity must be an array',
  },
  'event.activity.*.fromAddress': {
    in: 'body',
    optional: true,
    isString: true,
    errorMessage: 'Activity fromAddress must be a string',
  },
  'event.activity.*.toAddress': {
    in: 'body',
    isString: true,
    optional: true,
    errorMessage: 'Activity toAddress must be a string',
  },
  'event.activity.*.asset': {
    in: 'body',
    isString: true,
    optional: true,
    errorMessage: 'Activity asset must be a string',
  },
  'event.activity.*.value': {
    in: 'body',
    optional: true,
    errorMessage: 'Activity value must be a number',
  },
  'event.network': {
    in: 'body',
    isString: true,
    optional: true,
    errorMessage: 'Event network must be a string',
  },
  // Webhook metadata
  id: {
    in: 'body',
    isString: true,
    optional: true,
  },
  type: {
    in: 'body',
    isString: true,
    optional: true,
  },
  webhookId: {
    in: 'body',
    isString: true,
    optional: true,
  },
};

// Validation schemas
const signInSchema: Record<string, ParamSchema> = {
  signinMethod: {
    in: ['body'],
    isIn: {
      options: [[SigninMethod.SOCIAL, SigninMethod.PASSKEY, SigninMethod.EMAIL]],
    },
    errorMessage: `Must be one of: ${SigninMethod.SOCIAL}, ${SigninMethod.PASSKEY}, ${SigninMethod.EMAIL}`,
  },
  userEmail: {
    in: ['body'],
    optional: true,
    isEmail: true,
    errorMessage: 'Must be a valid email address',
    custom: {
      options: (value: string, { req }) => {
        // Require email for EMAIL signin method
        if (req.body.signinMethod === SigninMethod.EMAIL && !value) {
          throw new Error('userEmail is required for email signin');
        }
        return true;
      },
    },
  },
  magicLink: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Magic link must be a string',
  },
  targetPublicKey: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Target public key must be a string',
    custom: {
      options: (value: string, { req }) => {
        // Require targetPublicKey for EMAIL and SOCIAL signin methods
        const signinMethod = req.body.signinMethod;
        if ((signinMethod === SigninMethod.EMAIL || signinMethod === SigninMethod.SOCIAL) &&
          !value) {
          throw new Error('targetPublicKey is required for email and social signin');
        }
        return true;
      },
    },
  },
  // Passkey params
  challenge: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Challenge must be a string',
  },
  attestation: {
    in: ['body'],
    optional: true,
    isObject: true,
    errorMessage: 'Attestation must be an object',
  },
  provider: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Provider must be a string',
  },
  oidcToken: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'OIDC token must be a string',
  },
};

const uploadDydxAddressSchema: Record<string, ParamSchema> = {
  dydxAddress: {
    in: ['body'],
    isString: true,
    errorMessage: 'dydxAddress must be a string',
  },
  signature: {
    in: ['body'],
    isString: true,
    errorMessage: 'signature must be a string',
  },
};

const appleLoginRedirectSchema: Record<string, ParamSchema> = {
  state: {
    in: ['query'],
    isString: true,
    notEmpty: true,
    errorMessage: 'state (public key) is required and must be a non-empty string',
  },
  code: {
    in: ['query'],
    isString: true,
    notEmpty: true,
    errorMessage: 'code (authorization code) is required and must be a non-empty string',
  },
};

const getDepositAddressSchema: Record<string, ParamSchema> = {
  dydxAddress: {
    in: ['params'],
    ...addressSchema,
  },
};

export const CheckSignInSchema = checkSchema(signInSchema);

export const CheckUploadDydxAddressSchema = checkSchema(uploadDydxAddressSchema);

export const CheckAppleLoginRedirectSchema = checkSchema(appleLoginRedirectSchema);

export const CheckGetDepositAddressSchema = checkSchema(getDepositAddressSchema);

export const CheckBridgeSchema = checkSchema(checkBridgeSchema);

export const CheckZeroPaymentsOptionalParamSchema = checkSchema({
  showZeroPayments: checkZeroPaymentsOptionalParamSchema,
});

export const CheckDailyOptionalSchema = checkSchema({
  daily: checkDailyOptionalParamSchema,
});

export const CheckTickerParamSchema = checkSchema({
  ticker: checkTickerParamSchema,
});

export const CheckTickerOptionalQuerySchema = checkSchema({
  ticker: checkTickerOptionalQuerySchema,
});

export const CheckMarketOptionalQuerySchema = checkSchema({
  market: checkTickerOptionalQuerySchema,
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

export const UpdateReferralCodeSchema = (withTimestamp: boolean = true) => checkSchema({
  address: {
    in: ['body'],
    ...addressSchema,
  },
  newCode: {
    in: ['body'],
    isString: true,
    errorMessage: 'newCode must be a valid string',
    custom: {
      options: validateReferralCode,
    },
  },
  signedMessage: {
    in: ['body'],
    isString: true,
    errorMessage: 'signedMessage must be a valid string',
  },
  pubKey: {
    in: ['body'],
    isString: true,
    errorMessage: 'pubKey must be a valid string',
  },
  ...(withTimestamp ? {
    timestamp: {
      in: ['body'],
      isInt: true,
      errorMessage: 'timestamp must be a valid integer',
    },
  } : {}),
});

function validateReferralCode(code: string): boolean {
  if (code.length < 3 || code.length > 32 || !/^[a-zA-Z0-9]*$/.test(code)) {
    return false;
  }
  return true;
}

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
