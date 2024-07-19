import { IndexerPerpetualPosition, IndexerAssetPosition } from 'packages/v4-protos/build';
import {
  getPositionIsLong,
} from '../src/position-helpers';
import {
  bigIntToBytes,
} from '../src/bytes-helpers';
import Long from 'long';

const TEST_CASES: [bigint, boolean][] = [
  [BigInt('0'), false],
  [BigInt('-0'), false],
  [BigInt('1000'), true],
  [BigInt('123456'), true],
  [BigInt('-20'), false],
  [BigInt('-1'), false],
  [BigInt('-54321'), false],
];

describe('getPositionIsLong', () => {
  it('Asset position', () => {
    TEST_CASES.forEach(([i, isLong]: [bigint, boolean]) => {
      const assetPosition: IndexerAssetPosition = {
        assetId: 0,
        index: new Long(0),
        quantums: bigIntToBytes(i),
      };
      expect(getPositionIsLong(assetPosition)).toEqual(isLong);
    });
  });

  it('Perpetual position', () => {
    TEST_CASES.forEach(([i, isLong]: [bigint, boolean]) => {
      const perpetualPosition: IndexerPerpetualPosition = {
        perpetualId: 0,
        fundingIndex: bigIntToBytes(BigInt(0)),
        quantums: bigIntToBytes(i),
        fundingPayment: bigIntToBytes(BigInt(0)),
      };
      expect(getPositionIsLong(perpetualPosition)).toEqual(isLong);
    });
  });
});
