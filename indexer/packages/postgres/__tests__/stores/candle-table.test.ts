import { CandleFromDatabase, CandleUpdateObject } from '../../src/types';
import * as CandleTable from '../../src/stores/candle-table';
import {
  clearData,
  teardown,
} from '../../src/helpers/db-helpers';
import { defaultCandle, defaultCandleId, defaultPerpetualMarket2 } from '../helpers/constants';

describe('CandleTable', () => {
  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a Candle', async () => {
    await CandleTable.create(defaultCandle);
  });

  it('Successfully finds all Candles', async () => {
    await Promise.all([
      CandleTable.create(defaultCandle),
      CandleTable.create({
        ...defaultCandle,
        ticker: defaultPerpetualMarket2.ticker,
      }),
    ]);

    const candles: CandleFromDatabase[] = await CandleTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(candles.length).toEqual(2);
    expect(candles[0]).toEqual(expect.objectContaining({
      ...defaultCandle,
      ticker: defaultPerpetualMarket2.ticker,
    }));
    expect(candles[1]).toEqual(expect.objectContaining(defaultCandle));
  });

  it('Successfully finds a Candle', async () => {
    await CandleTable.create(defaultCandle);

    const candle: CandleFromDatabase | undefined = await CandleTable.findById(
      defaultCandleId,
    );

    expect(candle).toEqual(expect.objectContaining(defaultCandle));
  });

  it('Fails to finds a nonexistent Candle', async () => {
    const candle: CandleFromDatabase | undefined = await CandleTable.findById(
      defaultCandleId,
    );
    expect(candle).toEqual(undefined);
  });

  it('Successfully updates a candle', async () => {
    await CandleTable.create(defaultCandle);

    const updatedCandle: CandleUpdateObject = {
      id: defaultCandleId,
      open: '100',
      orderbookMidPriceClose: '200',
      orderbookMidPriceOpen: '300',
    };

    await CandleTable.update(updatedCandle);

    const candle: CandleFromDatabase | undefined = await CandleTable.findById(
      defaultCandleId,
    );

    expect(candle).toEqual(expect.objectContaining(updatedCandle));
  });
});
