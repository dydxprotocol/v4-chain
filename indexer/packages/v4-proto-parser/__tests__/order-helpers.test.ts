import { IndexerOrderId } from '@dydxprotocol-indexer/v4-protos';
import { getOrderIdHash, isLongTermOrder, isStatefulOrder } from '../src/order-helpers';
import { ORDER_FLAG_CONDITIONAL, ORDER_FLAG_LONG_TERM, ORDER_FLAG_SHORT_TERM } from '../src';

describe('getOrderIdHash', () => {
  // Test cases match test cases in V4
  // https://github.com/dydxprotocol/v4/blob/311411a3ce92230d4866a7c4abb1422fbc4ef3b9/indexer/off_chain_updates/off_chain_updates_test.go#L278-L299
  it('hashes an order id correctly', () => {
    const orderId: IndexerOrderId = {
      subaccountId: {
        owner: 'dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4',
        number: 0,
      },
      clientId: 0,
      orderFlags: 0,
      clobPairId: 0,
    };
    const hash: Buffer = getOrderIdHash(orderId);

    const expectedHash: Buffer = Buffer.from([
      0x5c, 0x2e, 0xc7, 0xcb, 0xfc, 0xf2, 0x63, 0xf6, 0xcb, 0x37, 0x44, 0x12, 0x60, 0xcc, 0x8b,
      0x71, 0xdd, 0x28, 0xb7, 0xfc, 0x8f, 0x00, 0xff, 0x00, 0xc1, 0x39, 0x16, 0x45, 0xce, 0x53,
      0x21, 0x95,
    ]);
    expect(hash).toEqual(expectedHash);
  });

  it('hashes empty order id correctly', () => {
    const orderId: IndexerOrderId = {
      subaccountId: {
        owner: '',
        number: 0,
      },
      clientId: 0,
      orderFlags: 0,
      clobPairId: 0,
    };
    const hash: Buffer = getOrderIdHash(orderId);

    const expectedHash: Buffer = Buffer.from([
      0x10, 0x2b, 0x51, 0xb9, 0x76, 0x5a, 0x56, 0xa3, 0xe8, 0x99, 0xf7, 0xcf, 0x0e, 0xe3, 0x8e,
      0x52, 0x51, 0xf9, 0xc5, 0x03, 0xb3, 0x57, 0xb3, 0x30, 0xa4, 0x91, 0x83, 0xeb, 0x7b, 0x15,
      0x56, 0x04,
    ]);
    expect(hash).toEqual(expectedHash);
  });
});

describe('isStatefulOrder', () => {
  it.each([
    [ORDER_FLAG_SHORT_TERM.toString(), 'string', false],
    ['4', 'string', false],
    [ORDER_FLAG_CONDITIONAL.toString(), 'string', true],
    [ORDER_FLAG_LONG_TERM.toString(), 'string', true],
    [ORDER_FLAG_SHORT_TERM, 'number', false],
    [3, 'number', false],
    [ORDER_FLAG_CONDITIONAL, 'number', true],
    [ORDER_FLAG_LONG_TERM, 'number', true],
    ['abc', 'string', false],
  ])('Checks if flag %s with type %s is a stateful order', (
    flag: number | string,
    _type: string,
    isStateful: boolean,
  ) => {
    expect(isStatefulOrder(flag)).toEqual(isStateful);
  });
});

describe('isLongTermOrder', () => {
  it.each([
    [ORDER_FLAG_SHORT_TERM.toString(), 'string', false],
    ['4', 'string', false],
    [ORDER_FLAG_CONDITIONAL.toString(), 'string', false],
    [ORDER_FLAG_LONG_TERM.toString(), 'string', true],
    [ORDER_FLAG_SHORT_TERM, 'number', false],
    [3, 'number', false],
    [ORDER_FLAG_CONDITIONAL, 'number', false],
    [ORDER_FLAG_LONG_TERM, 'number', true],
    ['abc', 'string', false],
  ])('Checks if flag %s with type %s is a long term order', (
    flag: number | string,
    _type: string,
    isLongTerm: boolean,
  ) => {
    expect(isLongTermOrder(flag)).toEqual(isLongTerm);
  });
});
