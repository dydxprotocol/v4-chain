import { IndexerOrderId } from '@klyraprotocol-indexer/v4-protos';
import { getOrderIdHash, isLongTermOrder, isStatefulOrder } from '../src/order-helpers';
import { ORDER_FLAG_CONDITIONAL, ORDER_FLAG_LONG_TERM, ORDER_FLAG_SHORT_TERM } from '../src';

describe('getOrderIdHash', () => {
  it('hashes an order id correctly', () => {
    const orderId: IndexerOrderId = {
      subaccountId: {
        owner: 'klyra199tqg4wdlnu4qjlxchpd7seg454937hju8xa57',
        number: 0,
      },
      clientId: 0,
      orderFlags: 0,
      clobPairId: 0,
    };
    const hash: Buffer = getOrderIdHash(orderId);

    const expectedHash: Buffer = Buffer.from([
      0x82, 0xba, 0x4e, 0xbd, 0x7b, 0x36, 0x58, 0x3c, 0x30, 0x37, 0x5b, 0x88, 0xb5, 0x9a, 0x8f,
      0x34, 0x9b, 0x7a, 0x4e, 0xfa, 0x6e, 0xe3, 0x67, 0x65, 0x3c, 0xdf, 0x50, 0x40, 0x7e, 0xa3,
      0x26, 0x27,
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
