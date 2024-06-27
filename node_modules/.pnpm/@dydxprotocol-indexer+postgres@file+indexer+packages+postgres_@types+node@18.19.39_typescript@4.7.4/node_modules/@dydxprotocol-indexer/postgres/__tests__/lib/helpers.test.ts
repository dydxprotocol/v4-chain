import { blockTimeFromIsoString } from '../../src/lib/helpers';

describe('helpers', () => {

  const expectedGoodTilBlockTimeISO: string = '2017-07-14T02:40:00.000Z';

  describe('getGoodTilBlockTimeFromIsoString', () => {
    it('gets goodTilBlockTime as ISO string for order', () => {
      expect(blockTimeFromIsoString(expectedGoodTilBlockTimeISO))
        .toEqual(1_500_000_000);
    });
  });

});
