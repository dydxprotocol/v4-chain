import { MarketFromDatabase } from '../../src/types';
import * as MarketTable from '../../src/stores/market-table';
import {
  clearData,
  migrate,
  teardown,
} from '../../src/helpers/db-helpers';
import { UniqueViolationError } from 'objection';
import { defaultMarket, defaultMarket2 } from '../helpers/constants';
import Transaction from '../../src/helpers/transaction';

describe('Market store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
    jest.resetAllMocks();
  });

  afterAll(async () => {
    await teardown();
    jest.clearAllMocks();
  });

  it('Successfully creates a market', async () => {
    await MarketTable.create(defaultMarket);
  });

  it('Fails to create second market with the same ID', async () => {
    try {
      await Promise.all([
        MarketTable.create(defaultMarket),
        MarketTable.create(defaultMarket),
      ]);
    } catch (e) {
      expect(e).toBeInstanceOf(UniqueViolationError);
    }
  });

  it('Successfully finds all Markets', async () => {
    await Promise.all([
      MarketTable.create(defaultMarket),
      MarketTable.create(defaultMarket2),
    ]);

    const markets: MarketFromDatabase[] = await MarketTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(markets.length).toEqual(2);
    expect(markets[0]).toEqual(expect.objectContaining(defaultMarket));
    expect(markets[1]).toEqual(expect.objectContaining(defaultMarket2));
  });

  it('Successfully finds market with pair', async () => {
    await Promise.all([
      MarketTable.create(defaultMarket),
      MarketTable.create(defaultMarket2),
    ]);

    const market: MarketFromDatabase = await MarketTable.findByPair(
      defaultMarket.pair,
      {},
    ) as MarketFromDatabase;

    expect(market).toEqual(expect.objectContaining(defaultMarket));
  });

  it('Successfully finds a market', async () => {
    await MarketTable.create(defaultMarket);

    const market: MarketFromDatabase | undefined = await MarketTable.findById(
      defaultMarket.id,
    );

    expect(market).toEqual(expect.objectContaining(defaultMarket));
  });

  it('Unable to find a market', async () => {
    const market: MarketFromDatabase | undefined = await MarketTable.findById(
      defaultMarket.id,
    );
    expect(market).toEqual(undefined);
  });

  it('Successfully updates a market', async () => {
    await MarketTable.create(defaultMarket);

    const market: MarketFromDatabase | undefined = await MarketTable.update({
      id: defaultMarket.id,
      minPriceChangePpm: 100,
    });

    expect(market).toEqual(expect.objectContaining({
      ...defaultMarket,
      minPriceChangePpm: 100,
    }));
  });

  it('Successfully updates a market created in the same transaction', async () => {
    const txId: number = await Transaction.start();
    await MarketTable.create(defaultMarket, { txId });
    const market: MarketFromDatabase | undefined = await MarketTable.update(
      {
        id: defaultMarket.id,
        minPriceChangePpm: 100,
      },
      {
        txId,
      },
    );
    expect(market).toEqual(expect.objectContaining({
      ...defaultMarket,
      minPriceChangePpm: 100,
    }));
    await Transaction.commit(txId);
  });

  it('Fails to update market to have same pair as existing market', async () => {
    try {
      await MarketTable.create(defaultMarket);
      await MarketTable.create({
        id: 1, pair: 'ETH-USD', exponent: -5, minPriceChangePpm: 100,
      });
      await MarketTable.update({
        id: defaultMarket.id,
        pair: 'ETH-USD',
      });
    } catch (e) {
      expect(e).toBeInstanceOf(UniqueViolationError);
    }
  });
});
