import {
  ICancelOrder,
  IPlaceOrder,
  OrderFlags,
  Transfer,
} from '../../src/types';
import { MAX_SUBACCOUNT_NUMBER, MAX_UINT_32 } from '../../src/lib/constants';
import { UserError } from '../../src/lib/errors';
import {
  isValidAddress,
  validateCancelOrderMessage,
  validatePlaceOrderMessage,
  validateTransferMessage,
} from '../../src/lib/validation';
import {
  TEST_ADDRESS, defaultCancelOrder, defaultOrder, defaultTransfer,
} from '../helpers/constants';
import Long from 'long';

const MAX_UINT_32_PLUS_1: number = MAX_UINT_32 + 1;
const MAX_SUBACCOUNT_NUMBER_PLUS_1: number = MAX_SUBACCOUNT_NUMBER + 1;

describe('Validations', () => {
  it.each([
    ['valid', defaultOrder, undefined],
    [
      '0 value clientId',
      { ...defaultOrder, clientId: 0 },
      undefined,
    ],
    [
      'underflow clientId',
      { ...defaultOrder, clientId: -1 },
      new UserError(`clientId: ${-1} is not a valid uint32`),
    ],
    [
      'overflow clientId',
      { ...defaultOrder, clientId: MAX_UINT_32_PLUS_1 },
      new UserError(`clientId: ${MAX_UINT_32_PLUS_1} is not a valid uint32`),
    ],
    [
      'underflow quantums',
      { ...defaultOrder, quantums: Long.NEG_ONE },
      new UserError(`quantums: ${-1} cannot be <= 0`),
    ],
    [
      'underflow goodTilBlock',
      { ...defaultOrder, goodTilBlock: -1 },
      new UserError(`goodTilBlock: ${-1} is not a valid uint32 or is 0`),
    ],
    [
      'overflow goodTilBlock',
      { ...defaultOrder, goodTilBlock: MAX_UINT_32_PLUS_1 },
      new UserError(`goodTilBlock: ${MAX_UINT_32_PLUS_1} is not a valid uint32 or is 0`),
    ],
    [
      '0 goodTilBlock',
      { ...defaultOrder, goodTilBlock: 0 },
      new UserError(`goodTilBlock: ${0} is not a valid uint32 or is 0`),
    ],
    [
      'underflow subticks',
      { ...defaultOrder, subticks: Long.NEG_ONE },
      new UserError(`subticks: ${-1} cannot be <= 0`),
    ],
  ])('Validate order: %s', (_name: string, order: IPlaceOrder, expectedError: UserError | undefined) => {
    const validationError: UserError | void = validatePlaceOrderMessage(0, order);
    expect(validationError).toEqual(expectedError);
  });

  it.each([
    ['valid', defaultCancelOrder, undefined],
    [
      '0 value clientId',
      { ...defaultCancelOrder, clientId: 0 },
      undefined,
    ],
    [
      'underflow clientId',
      { ...defaultCancelOrder, clientId: -1 },
      new UserError(`clientId: ${-1} is not a valid uint32`),
    ],
    [
      'overflow clientId',
      { ...defaultCancelOrder, clientId: MAX_UINT_32_PLUS_1 },
      new UserError(`clientId: ${MAX_UINT_32_PLUS_1} is not a valid uint32`),
    ],
    [
      'underflow goodTilBlock',
      { ...defaultCancelOrder, goodTilBlock: -1 },
      new UserError(`goodTilBlock: ${-1} is not a valid uint32 or is 0`),
    ],
    [
      'overflow goodTilBlock',
      { ...defaultCancelOrder, goodTilBlock: MAX_UINT_32_PLUS_1 },
      new UserError(`goodTilBlock: ${MAX_UINT_32_PLUS_1} is not a valid uint32 or is 0`),
    ],
    [
      '0 goodTilBlock',
      { ...defaultCancelOrder, goodTilBlock: 0 },
      new UserError(`goodTilBlock: ${0} is not a valid uint32 or is 0`),
    ],
    [
      'contains GoodTilBlockTime',
      { ...defaultCancelOrder, goodTilBlockTime: 10 },
      new UserError('goodTilBlockTime is 10, but should not be set for non-stateful orders'),
    ],
    [
      'stateful order - valid',
      {
        ...defaultCancelOrder,
        orderFlags: OrderFlags.LONG_TERM,
        goodTilBlock: undefined,
        goodTilBlockTime: 10,
      },
      undefined,
    ],
    [
      'stateful order - undefined goodTilBlockTime',
      { ...defaultCancelOrder, orderFlags: OrderFlags.LONG_TERM },
      new UserError(`goodTilBlockTime: ${undefined} is not a valid uint32 or is 0`),
    ],
    [
      'stateful order - zero goodTilBlockTime',
      { ...defaultCancelOrder, orderFlags: OrderFlags.LONG_TERM, goodTilBlockTime: 0 },
      new UserError(`goodTilBlockTime: ${0} is not a valid uint32 or is 0`),
    ],
    [
      'stateful order - underflow goodTilBlockTime',
      { ...defaultCancelOrder, orderFlags: OrderFlags.LONG_TERM, goodTilBlockTime: -1 },
      new UserError(`goodTilBlockTime: ${-1} is not a valid uint32 or is 0`),
    ],
    [
      'stateful order - overflow goodTilBlockTime',
      {
        ...defaultCancelOrder,
        orderFlags: OrderFlags.LONG_TERM,
        goodTilBlockTime: MAX_UINT_32_PLUS_1,
      },
      new UserError(`goodTilBlockTime: ${MAX_UINT_32_PLUS_1} is not a valid uint32 or is 0`),
    ],
    [
      'stateful order - has GoodTilBlock',
      {
        ...defaultCancelOrder,
        orderFlags: OrderFlags.LONG_TERM,
        goodTilBlock: 10,
        goodTilBlockTime: 10,
      },
      new UserError('goodTilBlock is 10, but should not be set for stateful orders'),
    ],
  ])('Validate cancel order: %s', (_name: string, order: ICancelOrder, expectedError: UserError | undefined) => {
    const validationError: UserError | void = validateCancelOrderMessage(0, order);
    expect(validationError).toEqual(expectedError);
  });

  it.each([
    ['valid', defaultTransfer, undefined],
    [
      'underflow senderSubaccountNumber',
      { ...defaultTransfer, sender: { owner: TEST_ADDRESS, number: -1 } },
      new UserError(`senderSubaccountNumber: ${-1} cannot be < 0 or > ${MAX_SUBACCOUNT_NUMBER}`),
    ],
    [
      'exceeds max subaccount number - senderSubaccountNumber',
      { ...defaultTransfer, sender: { owner: TEST_ADDRESS, number: MAX_SUBACCOUNT_NUMBER_PLUS_1 } },
      new UserError(
        `senderSubaccountNumber: ${MAX_SUBACCOUNT_NUMBER_PLUS_1} cannot be < 0 or > ${MAX_SUBACCOUNT_NUMBER}`,
      ),
    ],
    [
      '0 senderSubaccountNumber',
      { ...defaultTransfer, sender: { owner: TEST_ADDRESS, number: 0 } },
      undefined,
    ],
    [
      'underflow recipientSubaccountNumber',
      { ...defaultTransfer, recipient: { owner: TEST_ADDRESS, number: -1 } },
      new UserError(`recipientSubaccountNumber: ${-1} cannot be < 0 or > ${MAX_SUBACCOUNT_NUMBER}`),
    ],
    [
      'exceeds max subaccount number - recipient.subaccountNumber',
      {
        ...defaultTransfer,
        recipient: { owner: TEST_ADDRESS, number: MAX_SUBACCOUNT_NUMBER_PLUS_1 },
      },
      new UserError(
        `recipientSubaccountNumber: ${MAX_SUBACCOUNT_NUMBER_PLUS_1} cannot be < 0 or > ${MAX_SUBACCOUNT_NUMBER}`,
      ),
    ],
    [
      '0 recipientSubaccountNumber',
      { ...defaultTransfer, recipient: { owner: TEST_ADDRESS, number: 0 } },
      undefined,
    ],
    [
      'non-zero asset id',
      { ...defaultTransfer, assetId: 1 },
      new UserError(`asset id: ${1} not supported`),
    ],
    [
      '0 amount',
      { ...defaultTransfer, amount: Long.ZERO },
      new UserError(`amount: ${0} cannot be <= 0`),
    ],
    [
      'too short recipientAddress',
      {
        ...defaultTransfer,
        recipient: {
          owner: 'dydx14063jves4u9zhm7eja5ltf3t8zspxd92qnk23',
          number: 0,
        },
      },
      new UserError('Error: Invalid checksum for dydx14063jves4u9zhm7eja5ltf3t8zspxd92qnk23'),
    ],
    [
      'invalid recipientAddress',
      {
        ...defaultTransfer,
        recipient: {
          owner: 'fakeAddress1234',
          number: 0,
        },
      },
      new UserError('Error: Mixed-case string fakeAddress1234'),
    ],
  ])(
    'Validate transfer: %s',
    (_name: string, transfer: Transfer, expectedError: UserError | undefined) => {
      const validationError: UserError | void = validateTransferMessage(transfer);
      expect(validationError).toEqual(expectedError);
    });

  it.each([
    ['valid', 'dydx17xpfvakm2amg962yls6f84z3kell8c5leqdyt2', true],
    ['invalid: does not start with dydx1', 'dydx27xpfvakm2amg962yls6f84z3kell8c5leqdyt2', false],
    ['invalid: too short', 'dydx17xpfvakm2amg962yls6f84z3kell8c5leqdyt', false],
    ['invalid: too long', 'dydx17xpfvakm2amg962yls6f84z3kell8c5leqdyt2s', false],
  ])('Validate address: %s', (_name: string, address: string, expectedResult: boolean) => {
    const validationResult: boolean = isValidAddress(address);
    expect(validationResult).toEqual(expectedResult);
  });
});
