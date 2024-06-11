import {
  MarketColumns,
  MarketFromDatabase,
  MarketsMap,
  PerpetualPositionFromDatabase,
  PositionSide,
} from '../../src/types';
import { getMaintenanceMarginPpm, getUnrealizedPnl, getUnsettledFunding } from '../../src/db/helpers';
import {
  defaultFundingIndexUpdate,
  defaultPerpetualMarket,
  defaultPerpetualPosition,
  defaultPerpetualPositionId,
} from '../helpers/constants';
import * as PerpetualPositionTable from '../../src/stores/perpetual-position-table';
import * as MarketTable from '../../src/stores/market-table';
import Big from 'big.js';
import { seedData } from '../helpers/mock-generators';
import _ from 'lodash';
import { clearData } from '../../src/helpers/db-helpers';

describe('helpers', () => {

  afterEach(async () => {
    await clearData();
  });

  describe('getUnsettledFunding', () => {
    const position: PerpetualPositionFromDatabase = {
      ...defaultPerpetualPosition,
      id: defaultPerpetualPositionId,
      entryPrice: defaultPerpetualPosition.entryPrice as string,
      sumOpen: defaultPerpetualPosition.sumOpen as string,
      sumClose: defaultPerpetualPosition.sumClose as string,
    };

    it('compute unsettled funding for long position', () => {
      expect(
        getUnsettledFunding(
          position,
          {
            [defaultFundingIndexUpdate.perpetualId]: Big('12050'),
          },
          {
            [defaultFundingIndexUpdate.perpetualId]: Big('10050'),
          },
        ),
      ).toEqual(Big('-20000'));  // 10 * (10050-12050). longs pay shorts when funding index is increasing.
    });

    it('compute unsettled funding for short position', () => {
      const shortPosition: PerpetualPositionFromDatabase = {
        ...position,
        side: PositionSide.SHORT,
        size: '-10',
      };
      expect(
        getUnsettledFunding(
          shortPosition,
          {
            [defaultFundingIndexUpdate.perpetualId]: Big('12050'),
          },
          {
            [defaultFundingIndexUpdate.perpetualId]: Big('10050'),
          },
        ),
      ).toEqual(Big('20000'));  // -10 * (10050-12050). longs pay shorts when funding index is increasing.
    });

    it('compute unsettled funding for decimal position', () => {
      const shortPosition: PerpetualPositionFromDatabase = {
        ...position,
        size: '1.35',
      };
      expect(
        getUnsettledFunding(
          shortPosition,
          {
            [defaultFundingIndexUpdate.perpetualId]: Big('12050.124'),
          },
          {
            [defaultFundingIndexUpdate.perpetualId]: Big('10050'),
          },
        ),
      ).toEqual(Big('-2700.1674'));  // 1.35 * (10050-12050.124). longs pay shorts when funding index is increasing.
    });
  });

  describe('getUnrealizedPnl', () => {
    it('getUnrealizedPnl long', async () => {
      await seedData();

      const perpetualPosition: PerpetualPositionFromDatabase = await
      PerpetualPositionTable.create(defaultPerpetualPosition);

      const markets: MarketFromDatabase[] = await MarketTable.findAll({}, []);

      const marketIdToMarket: MarketsMap = _.keyBy(
        markets,
        MarketColumns.id,
      );

      const unrealizedPnl: string = getUnrealizedPnl(
        perpetualPosition,
        defaultPerpetualMarket,
        marketIdToMarket[defaultPerpetualMarket.marketId],
      );

      expect(unrealizedPnl).toEqual(Big(-50000).toFixed());
    });

    it('getUnrealizedPnl short', async () => {
      await seedData();

      const perpetualPosition: PerpetualPositionFromDatabase = await
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        side: PositionSide.SHORT,
        size: '-10',
      });

      const markets: MarketFromDatabase[] = await MarketTable.findAll({}, []);

      const marketIdToMarket: MarketsMap = _.keyBy(
        markets,
        MarketColumns.id,
      );

      const unrealizedPnl: string = getUnrealizedPnl(
        perpetualPosition,
        defaultPerpetualMarket,
        marketIdToMarket[defaultPerpetualMarket.marketId],
      );

      expect(unrealizedPnl).toEqual(Big(50000).toFixed());
    });
  });

  describe('getMaintenanceMarginPpm', () => {
    it('5% initial margin, 60% maintenance fraction', () => {
      expect(getMaintenanceMarginPpm(50_000, 600_000)).toEqual(30_000);
    });

    it('25% initial margin, 100% maintenance fraction', () => {
      expect(getMaintenanceMarginPpm(250_000, 1_000_000)).toEqual(250_000);
    });

    it('100% initial margin, 0% maintenance fraction', () => {
      expect(getMaintenanceMarginPpm(1_000_000, 0)).toEqual(0);
    });

    it('0% initial margin, 100% maintenance fraction', () => {
      expect(getMaintenanceMarginPpm(0, 1_000_000)).toEqual(0);
    });

    it('0% initial margin, 0% maintenance fraction', () => {
      expect(getMaintenanceMarginPpm(0, 0)).toEqual(0);
    });
  });
});
