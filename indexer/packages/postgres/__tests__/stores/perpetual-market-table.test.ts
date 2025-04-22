import {
  MarketCreateObject, PerpetualMarketFromDatabase, PerpetualMarketStatus, PerpetualMarketWithMarket,
} from '../../src/types';
import * as PerpetualMarketTable from '../../src/stores/perpetual-market-table';
import * as LiquidityTiersTable from '../../src/stores/liquidity-tiers-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import {
  defaultLiquidityTier,
  defaultLiquidityTier2,
  defaultMarket,
  defaultMarket2,
  defaultPerpetualMarket,
  invalidTicker,
} from '../helpers/constants';
import * as MarketTable from '../../src/stores/market-table';
import _ from 'lodash';

describe('PerpetualMarket store', () => {
  beforeAll(async () => {
    await migrate();
  });

  beforeEach(async () => {
    await Promise.all([
      MarketTable.create(defaultMarket),
      MarketTable.create(defaultMarket2),
    ]);
    await Promise.all([
      LiquidityTiersTable.create(defaultLiquidityTier),
      LiquidityTiersTable.create(defaultLiquidityTier2),
    ]);
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a PerpetualMarket', async () => {
    await PerpetualMarketTable.create(defaultPerpetualMarket);
  });

  it('Successfully finds all PerpetualMarkets', async () => {
    await Promise.all([
      PerpetualMarketTable.create(defaultPerpetualMarket),
      PerpetualMarketTable.create({
        ...defaultPerpetualMarket,
        id: '1',
      }),
    ]);

    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(perpetualMarkets.length).toEqual(2);
    expect(perpetualMarkets[0]).toEqual(expect.objectContaining(defaultPerpetualMarket));
    expect(perpetualMarkets[1]).toEqual(expect.objectContaining({
      ...defaultPerpetualMarket,
      id: '1',
    }));
  });

  it('Successfully finds all PerpetualMarkets joined with markets', async () => {
    await Promise.all([
      PerpetualMarketTable.create(defaultPerpetualMarket),
      PerpetualMarketTable.create({
        ...defaultPerpetualMarket,
        id: '1',
        marketId: defaultMarket2.id,
      }),
    ]);

    const perpetualWithMarkets: PerpetualMarketWithMarket[] = await PerpetualMarketTable.findAll(
      { joinWithMarkets: true },
      [],
      { readReplica: true },
    ) as PerpetualMarketWithMarket[];

    expect(perpetualWithMarkets.length).toEqual(2);
    expect(perpetualWithMarkets[0]).toEqual({
      ...defaultPerpetualMarket,
      ..._.omit(defaultMarket, 'id'),
    });
    expect(perpetualWithMarkets[1]).toEqual({
      ...defaultPerpetualMarket,
      id: '1',
      marketId: defaultMarket2.id,
      ..._.omit(defaultMarket2, 'id'),
    });
  });

  it('Successfully finds all PerpetualMarkets joined with markets, filter by market id', async () => {
    await Promise.all([
      PerpetualMarketTable.create(defaultPerpetualMarket),
      PerpetualMarketTable.create({
        ...defaultPerpetualMarket,
        id: '1',
        marketId: defaultMarket2.id,
      }),
    ]);

    const perpetualWithMarkets: PerpetualMarketWithMarket[] = await PerpetualMarketTable.findAll(
      { marketId: [1], joinWithMarkets: true },
      [],
      { readReplica: true },
    ) as PerpetualMarketWithMarket[];

    expect(perpetualWithMarkets.length).toEqual(1);
    expect(perpetualWithMarkets[0]).toEqual({
      ...defaultPerpetualMarket,
      id: '1',
      marketId: defaultMarket2.id,
      ..._.omit(defaultMarket2, 'id'),
    });
  });

  it('Successfully finds a PerpetualMarket', async () => {
    await PerpetualMarketTable.create(defaultPerpetualMarket);

    const perpetualMarket: PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable
      .findById(
        '0',
      );

    expect(perpetualMarket).toEqual(expect.objectContaining(defaultPerpetualMarket));
  });

  it('Successfully finds a PerpetualMarket by market id', async () => {
    await PerpetualMarketTable.create(defaultPerpetualMarket);

    const perpetualMarket: PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable
      .findByMarketId(
        defaultMarket.id,
      );

    expect(perpetualMarket).toEqual(expect.objectContaining(defaultPerpetualMarket));
  });

  it('Successfully finds a PerpetualMarket by clob pair id', async () => {
    await PerpetualMarketTable.create(defaultPerpetualMarket);

    const perpetualMarket: PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable
      .findByClobPairId(
        '1',
      );

    expect(perpetualMarket).toEqual(expect.objectContaining(defaultPerpetualMarket));
  });

  it('Unable finds a PerpetualMarket', async () => {
    const perpetualMarket: PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable
      .findById(
        '0',
      );

    expect(perpetualMarket).toEqual(undefined);
  });

  it('Unable finds a PerpetualMarket by market id', async () => {
    const perpetualMarket: PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable
      .findByMarketId(
        defaultMarket.id,
      );

    expect(perpetualMarket).toEqual(undefined);
  });

  it.each([
    ['market with ticker exists', defaultPerpetualMarket.ticker, defaultPerpetualMarket],
    ['market with ticker does not exist', invalidTicker, undefined],
  ])('Finds a PerpetualMarket by ticker: %s', async (
    _name: string,
    ticker: string,
    expectedPerpetualMarket?: Object,
  ) => {
    await PerpetualMarketTable.create(defaultPerpetualMarket);

    const perpetualMarket: PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable
      .findByTicker(
        ticker,
      );

    if (expectedPerpetualMarket !== undefined) {
      expect(perpetualMarket).toEqual(expect.objectContaining(defaultPerpetualMarket));
    } else {
      expect(perpetualMarket).toEqual(expectedPerpetualMarket);
    }
  });

  it('Successfully updates a perpetual market', async () => {
    await PerpetualMarketTable.create(defaultPerpetualMarket);

    const perpetualMarket: PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable
      .update({
        id: defaultPerpetualMarket.id,
        trades24H: 100,
      });

    expect(perpetualMarket).toEqual(expect.objectContaining({
      ...defaultPerpetualMarket,
      trades24H: 100,
    }));
  });

  it('Successfully winds down a perpetual market', async () => {
    await PerpetualMarketTable.create(defaultPerpetualMarket);

    const perpetualMarket: PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable
      .update({
        id: defaultPerpetualMarket.id,
        status: PerpetualMarketStatus.FINAL_SETTLEMENT,
      });

    expect(perpetualMarket).toEqual(expect.objectContaining({
      ...defaultPerpetualMarket,
      status: PerpetualMarketStatus.FINAL_SETTLEMENT,
    }));
  });

  it('Successfully updates a perpetual market by market id', async () => {
    const market: MarketCreateObject = {
      id: 5,
      pair: 'DYDX-USD',
      exponent: -5,
      minPriceChangePpm: 50,
    };

    await MarketTable.create(market);
    await PerpetualMarketTable.create({
      ...defaultPerpetualMarket,
      marketId: 5,
    });

    const perpetualMarket: PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable
      .updateByMarketId({
        marketId: 5,
        trades24H: 100,
      });

    expect(perpetualMarket).toEqual(expect.objectContaining({
      ...defaultPerpetualMarket,
      marketId: 5,
      trades24H: 100,
    }));
  });
});
