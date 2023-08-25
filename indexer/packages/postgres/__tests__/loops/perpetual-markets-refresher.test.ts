import * as PerpetualMarketTable from '../../src/stores/perpetual-market-table';
import {
  getClobPairIdFromTicker,
  getClobPairIdToPerpetualMarket,
  getPerpetualMarketsMap,
  getPerpetualMarketFromClobPairId,
  getPerpetualMarketFromTicker,
  getPerpetualMarketTicker,
  getTickerToPerpetualMarketForTest,
  isValidPerpetualMarketTicker,
  updatePerpetualMarkets,
  clear,
  addPerpetualMarket,
  getPerpetualMarketFromId,
} from '../../src/loops/perpetual-market-refresher';
import { PerpetualMarketColumns, PerpetualMarketFromDatabase } from '../../src/types';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { seedData } from '../helpers/mock-generators';
import _ from 'lodash';
import { defaultPerpetualMarket } from '../helpers/constants';

describe('perpetual_markets_refresher', () => {
  let perpetualMarkets: PerpetualMarketFromDatabase[];
  const invalidTicker: string = 'INVALID-INVALID';
  const invalidClobPairId: string = '4125';

  const newId: string = '3';
  const newTicker: string = 'NEW-TICKER';
  const newClobPairId: string = '3';

  beforeAll(async () => {
    await migrate();
    await seedData();
  });

  beforeEach(async () => {
    await updatePerpetualMarkets();
    perpetualMarkets = await PerpetualMarketTable.findAll(
      {},
      [],
      { readReplica: true },
    );
  });

  afterEach(async () => {
    await clear();
  });

  afterAll(async () => {
    await clearData();
    await teardown();
  });

  describe('updatePerpetualMarkets', () => {
    it('updates in-memory mapping of perpetual markets', () => {
      const clobPairIdToPerpetualMarket: Record<
        string,
        PerpetualMarketFromDatabase> = getClobPairIdToPerpetualMarket();
      const tickerToPerpetualMarket: Record<
        string,
        PerpetualMarketFromDatabase> = getTickerToPerpetualMarketForTest();
      const idToPerpetualMarket: Record<
        string,
        PerpetualMarketFromDatabase> = getPerpetualMarketsMap();

      perpetualMarkets.forEach(
        (perpetualMarket: PerpetualMarketFromDatabase) => {
          expect(clobPairIdToPerpetualMarket[perpetualMarket.clobPairId]).toEqual(perpetualMarket);
          expect(tickerToPerpetualMarket[perpetualMarket.ticker]).toEqual(perpetualMarket);
          expect(idToPerpetualMarket[perpetualMarket.id]).toEqual(perpetualMarket);
        },
      );

      Object.keys(clobPairIdToPerpetualMarket).forEach(
        (clobPairId: string) => {
          expect(_.map(perpetualMarkets, PerpetualMarketColumns.clobPairId))
            .toContain(clobPairId);
        },
      );

      Object.keys(tickerToPerpetualMarket).forEach(
        (ticker: string) => {
          expect(_.map(perpetualMarkets, PerpetualMarketColumns.ticker))
            .toContain(ticker);
        },
      );

      Object.keys(idToPerpetualMarket).forEach(
        (id: string) => {
          expect(_.map(perpetualMarkets, PerpetualMarketColumns.id))
            .toContain(id);
        },
      );
    });
  });

  describe('isValidPerpetualMarketTicker', () => {
    it('returns true for valid ticker', () => {
      expect(isValidPerpetualMarketTicker(perpetualMarkets[0].ticker)).toEqual(true);
    });

    it('returns false for invalid ticker', () => {
      expect(isValidPerpetualMarketTicker(invalidTicker)).toEqual(false);
    });
  });

  describe('getPerpetualMarketTicker', () => {
    it('gets ticker for clob pair id', () => {
      expect(getPerpetualMarketTicker(perpetualMarkets[0].clobPairId)).toEqual(
        perpetualMarkets[0].ticker,
      );
    });

    it('returns undefined for invalid clob pair id', () => {
      expect(getPerpetualMarketTicker(invalidClobPairId)).toBeUndefined();
    });
  });

  describe('getClobPairIdFromTicker', () => {
    it('gets clob pair id for ticker', () => {
      expect(getClobPairIdFromTicker(perpetualMarkets[0].ticker)).toEqual(
        perpetualMarkets[0].clobPairId,
      );
    });

    it('returns undefined for invalid ticker', () => {
      expect(getClobPairIdFromTicker(invalidTicker)).toBeUndefined();
    });
  });

  describe('getPerpetualMarketFromTicker', () => {
    it('gets perpetual market for ticker', () => {
      expect(getPerpetualMarketFromTicker(perpetualMarkets[0].ticker)).toEqual(
        perpetualMarkets[0],
      );
    });

    it('returns undefined for invalid ticker', () => {
      expect(getPerpetualMarketFromTicker(invalidTicker)).toBeUndefined();
    });
  });

  describe('getPerpetualMarketFromClobPairId', () => {
    it('gets perpetual market for clob pair id', () => {
      expect(getPerpetualMarketFromClobPairId(perpetualMarkets[0].clobPairId)).toEqual(
        perpetualMarkets[0],
      );
    });

    it('returns undefined for invalid clob pair id', () => {
      expect(getPerpetualMarketFromClobPairId(invalidClobPairId)).toBeUndefined();
    });
  });

  describe('addPerpetualMarket', () => {
    it.each([
      [
        'id',
        {
          ...defaultPerpetualMarket,
          clobPairId: newClobPairId,
          ticker: newTicker,
        },
        `Perpetual market with id ${defaultPerpetualMarket.id} already exists`,
      ],
      [
        'clobPairId',
        {
          ...defaultPerpetualMarket,
          id: newId,
          ticker: newTicker,
        },
        `Perpetual market with clob pair id ${defaultPerpetualMarket.clobPairId} already exists`,
      ],
      [
        'ticker',
        {
          ...defaultPerpetualMarket,
          id: newId,
          clobPairId: newClobPairId,
        },
        `Perpetual market with ticker ${defaultPerpetualMarket.ticker} already exists`,
      ],
    ])(
      'fails to add perpetual market with duplicate %s',
      (
        _ignore: string,
        perpetualMarket: PerpetualMarketFromDatabase,
        errorMessage: string,
      ) => {
        expect(() => {
          addPerpetualMarket(perpetualMarket);
        }).toThrow(errorMessage);
      },
    );

    it('successfully adds perpetual market', () => {
      addPerpetualMarket({
        ...defaultPerpetualMarket,
        id: newId,
        clobPairId: newClobPairId,
        ticker: newTicker,
      });

      expect(getPerpetualMarketFromId(newId)).not.toBeUndefined();
    });
  });
});
