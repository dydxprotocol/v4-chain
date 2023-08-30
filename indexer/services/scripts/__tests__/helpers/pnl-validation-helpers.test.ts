import {
  BlockTable,
  dbHelpers,
  FillTable,
  FundingIndexUpdatesTable,
  Liquidity,
  OraclePriceTable,
  OrderSide,
  OrderTable,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import {
  getPnl,
  getRealizedFunding,
  getUnsettledFunding,
} from '../../src/helpers/pnl-validation-helpers';
import Big from 'big.js';

describe('pnl-validation-helpers', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await OrderTable.create(testConstants.defaultOrder);
    await Promise.all([
      OraclePriceTable.create(testConstants.defaultOraclePrice),
      OraclePriceTable.create(testConstants.defaultOraclePrice2),
    ]);
    const blockHeights: string[] = ['3', '4', '5', '6', '7'];

    await Promise.all(blockHeights.map((height) => BlockTable.create({
      ...testConstants.defaultBlock,
      blockHeight: height,
    }),
    ));

    await Promise.all([
      FundingIndexUpdatesTable.create(testConstants.defaultFundingIndexUpdate),
      FundingIndexUpdatesTable.create({
        ...testConstants.defaultFundingIndexUpdate,
        effectiveAtHeight: '3',
        fundingIndex: '10100',
      }),
      FundingIndexUpdatesTable.create({
        ...testConstants.defaultFundingIndexUpdate,
        effectiveAtHeight: '4',
        fundingIndex: '10150',
      }),
      FundingIndexUpdatesTable.create({
        ...testConstants.defaultFundingIndexUpdate,
        effectiveAtHeight: '5',
        fundingIndex: '10200',
      }),

      FundingIndexUpdatesTable.create({
        ...testConstants.defaultFundingIndexUpdate,
        perpetualId: testConstants.defaultPerpetualMarket2.id,
        effectiveAtHeight: '3',
        fundingIndex: '100',
      }),
      FundingIndexUpdatesTable.create({
        ...testConstants.defaultFundingIndexUpdate,
        perpetualId: testConstants.defaultPerpetualMarket2.id,
        effectiveAtHeight: '4',
        fundingIndex: '150',
      }),
      FundingIndexUpdatesTable.create({
        ...testConstants.defaultFundingIndexUpdate,
        perpetualId: testConstants.defaultPerpetualMarket2.id,
        effectiveAtHeight: '5',
        fundingIndex: '200',
      }),
    ]);

    await Promise.all([
      FillTable.create(testConstants.defaultFill),
      FillTable.create({
        ...testConstants.defaultFill,
        createdAtHeight: '3',
        liquidity: Liquidity.MAKER,
        size: '3',
      }),
      FillTable.create({
        ...testConstants.defaultFill,
        eventId: testConstants.defaultTendermintEventId2,
        liquidity: Liquidity.TAKER,
        side: OrderSide.SELL,
        size: '4',
        createdAtHeight: '4',
      }),
      FillTable.create({
        ...testConstants.defaultFill,
        eventId: testConstants.defaultTendermintEventId2,
        liquidity: Liquidity.MAKER,
        side: OrderSide.SELL,
        size: '5',
        createdAtHeight: '5',
      }),
      FillTable.create({
        ...testConstants.defaultFill,
        eventId: testConstants.defaultTendermintEventId3,
        createdAtHeight: '3',
        liquidity: Liquidity.MAKER,
        size: '3',
        clobPairId: '2',
        price: '500',
      }),
      FillTable.create({
        ...testConstants.defaultFill,
        eventId: testConstants.defaultTendermintEventId3,
        liquidity: Liquidity.TAKER,
        side: OrderSide.SELL,
        size: '4',
        createdAtHeight: '4',
        clobPairId: '2',
        price: '500',
        fee: '-0.1',
      }),
    ]);
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.resetAllMocks();
  });

  // markets:
  // defaultMarket - id 0 (btc)
  // defaultMarket2 - id 1 (eth)
  //
  // oracle prices:
  // 10000 for defaultMarket at height '2'
  // 500 for defaultMarket2 at height '2'
  //
  // blocks 1-7
  //
  // funding indices:
  // btc
  // 10050 at height 2
  // 10100 at height 3
  // 10150 at height 4
  // 10200 at height 5
  // eth:
  // 100 at height 3
  // 150 at height 4
  // 200 at height 5
  //
  // fills:
  // btc, buy size 10 at price 20000, at height 2
  // btc, buy size 3 at price 20000, at height 3
  // btc, sell size 4 at price 20000, at height 4
  // btc, sell size 5 at price 20000, at height 5
  // eth, buy size 3 at price 500, at height 3
  // eth, sell size 4 at price 500, at height 4

  it('getUnrealizedFunding', async () => {
    const unrealizedFunding = await getUnsettledFunding(
      testConstants.defaultSubaccountId,
      '7',
    );
    // no unrealized funding for BTC.
    // unrealized funding for ETH: -1 * (200 - 150) = -50.
    expect(unrealizedFunding).toEqual(Big(-50));
  });

  it('getRealizedFunding', async () => {
    const realizedFunding = await getRealizedFunding(
      testConstants.defaultSubaccountId,
      '7',
    );
    // realized funding for BTC: 10*(10100-10050) + 13 * (10150-10100) + 9 * (10200 - 10150) = 1600
    // realized funding for ETH: 3*(150-100) = 150
    expect(realizedFunding).toEqual(Big(1750));
  });

  it('getPnl', async () => {
    const effectiveBeforeOrAtHeight: string = '7';
    const [
      costOfFills,
      totalValueOfOpenPositions,
      realizedFunding,
      unrealizedFunding,
      feesPaid,
    ]: [
      Big,
      Big,
      Big,
      Big,
      Big,
    ] = await Promise.all([
      FillTable.getCostOfFills(testConstants.defaultSubaccountId, effectiveBeforeOrAtHeight),
      FillTable.getTotalValueOfOpenPositions(
        testConstants.defaultSubaccountId,
        effectiveBeforeOrAtHeight,
      ),
      getRealizedFunding(testConstants.defaultSubaccountId, effectiveBeforeOrAtHeight),
      getUnsettledFunding(testConstants.defaultSubaccountId, effectiveBeforeOrAtHeight),
      FillTable.getFeesPaid(testConstants.defaultSubaccountId, effectiveBeforeOrAtHeight),
    ]);

    // cost of fills:
    // (-10-3+4+5)*20000 + (-3+4)*500 = -79500
    expect(costOfFills).toEqual(Big('-79500'));
    // total value of open positions:
    // 4 BTC * 10000 + -1 ETH * 500 = 40000 - 500 = 39500
    expect(totalValueOfOpenPositions).toEqual(Big('39500'));
    // realized funding for BTC: 10*(10100-10050) + 13 * (10150-10100) + 9 * (10200 - 10150) = 1600
    // realized funding for ETH: 3*(150-100) = 150
    expect(realizedFunding).toEqual(Big('1750'));
    // no unrealized funding for BTC.
    // unrealized funding for ETH: -1 * (200 - 150) = -50.
    expect(unrealizedFunding).toEqual(Big('-50'));
    // fees: 5*1.1-0.1 = 5.4
    expect(feesPaid).toEqual(Big('5.4'));

    const pnl: Big = await getPnl(
      testConstants.defaultSubaccountId,
      effectiveBeforeOrAtHeight,
    );
    // -79500 + 39500 - 1750 + 50 - 5.4 = -41705.4
    expect(pnl).toEqual(Big('-41705.4'));
  });
});
