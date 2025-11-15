import { getFixedRepresentation } from '../helpers/helpers';
import {
  PerpetualPositionFromDatabase,
  AssetPositionFromDatabase,
  PerpetualPositionTable,
  testConstants,
  testMocks,
  PerpetualMarketsMap,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  PerpetualMarketColumns,
  dbHelpers,
  MarketFromDatabase,
  MarketTable,
  MarketsMap,
  MarketColumns,
  FundingIndexUpdatesTable,
  SubaccountTable,
  SubaccountFromDatabase,
  BlockFromDatabase,
  BlockTable,
  FundingIndexMap,
  FundingIndexUpdatesCreateObject,
  USDC_SYMBOL,
  PositionSide,
  helpers,
  PerpetualPositionStatus,
  LiquidityTiersFromDatabase,
  LiquidityTiersTable,
  liquidityTierRefresher,
  PnlTicksFromDatabase,
  PnlTicksTable,
  AssetFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import {
  adjustUSDCAssetPosition,
  calculateEquityAndFreeCollateral,
  filterAssetPositions,
  filterPositionsByLatestEventIdPerPerpetual,
  getFundingIndexMaps,
  getMarginFraction,
  getSignedNotionalAndRisk,
  getTotalUnsettledFunding,
  getPerpetualPositionsWithUpdatedFunding,
  initializePerpetualPositionsWithFunding,
  getChildSubaccountNums,
  aggregateHourlyPnlTicks,
  getSubaccountResponse,
} from '../../src/lib/helpers';
import _ from 'lodash';
import Big from 'big.js';
import {
  defaultLiquidityTier,
  defaultMarket,
  defaultPerpetualMarket,
  defaultPerpetualMarket2,
  defaultTendermintEventId,
  defaultTendermintEventId2,
  defaultTendermintEventId3,
} from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';
import {
  AggregatedPnlTick, AssetPositionsMap, PerpetualPositionWithFunding, SubaccountResponseObject,
} from '../../src/types';
import { ZERO, ZERO_USDC_POSITION } from '../../src/lib/constants';
import { DateTime } from 'luxon';

describe('helpers', () => {
  afterEach(async () => {
    await dbHelpers.clearData();
  });

  const zeroSizeAssetPosition: AssetPositionFromDatabase = {
    ...testConstants.defaultAssetPosition,
    id: '',
    size: '0.0',
  };

  const defaultAssetPosition2: AssetPositionFromDatabase = {
    ...testConstants.defaultAssetPosition,
    id: '5',
    size: '99269.783787',
  };

  it('getFixedRepresentation', () => {
    const fixedRep: string = getFixedRepresentation(150125);
    expect(fixedRep).toEqual('150125');
  });

  it('filterAssetPositions with 0 size', () => {
    const assetPositions: AssetPositionFromDatabase[] = [
      zeroSizeAssetPosition,
      defaultAssetPosition2,
    ];
    const filteredPositions: AssetPositionFromDatabase[] = filterAssetPositions(assetPositions);
    expect(filteredPositions).toHaveLength(1);
    expect(filteredPositions[0]).toEqual(
      expect.objectContaining({
        ...defaultAssetPosition2,
      }),
    );
  });

  it('calculateEquityAndFreeCollateral', async () => {
    await testMocks.seedData();
    await liquidityTierRefresher.updateLiquidityTiers();

    const perpetualPosition: PerpetualPositionFromDatabase = await
    PerpetualPositionTable.create(testConstants.defaultPerpetualPosition);

    const usdcPositionSize: string = '175000';

    const [perpetualMarkets, markets]: [PerpetualMarketFromDatabase[],
      MarketFromDatabase[]] = await Promise.all([
      PerpetualMarketTable.findAll({}, []),
      MarketTable.findAll({}, []),
    ]);

    const perpetualIdToMarket: PerpetualMarketsMap = _.keyBy(
      perpetualMarkets,
      PerpetualMarketColumns.id,
    );
    const marketIdToMarket: MarketsMap = _.keyBy(
      markets,
      MarketColumns.id,
    );

    const {
      equity,
      freeCollateral,
    }: {
      equity: string,
      freeCollateral: string,
    } = calculateEquityAndFreeCollateral(
      [perpetualPosition],
      perpetualIdToMarket,
      marketIdToMarket,
      usdcPositionSize,
    );

    expect(equity).toEqual('325000');
    expect(freeCollateral).toEqual('317500');
  });

  it('calculateEquityAndFreeCollateral with SHORT position', async () => {
    await testMocks.seedData();
    await liquidityTierRefresher.updateLiquidityTiers();

    const perpetualPosition: PerpetualPositionFromDatabase = await
    PerpetualPositionTable.create({
      ...testConstants.defaultPerpetualPosition,
      side: PositionSide.SHORT,
      size: '-10',
    });

    const usdcPositionSize: string = '175000';

    const [perpetualMarkets, markets]: [PerpetualMarketFromDatabase[],
      MarketFromDatabase[]] = await Promise.all([
      PerpetualMarketTable.findAll({}, []),
      MarketTable.findAll({}, []),
    ]);

    const perpetualIdToMarket: PerpetualMarketsMap = _.keyBy(
      perpetualMarkets,
      PerpetualMarketColumns.id,
    );
    const marketIdToMarket: MarketsMap = _.keyBy(
      markets,
      MarketColumns.id,
    );

    const {
      equity,
      freeCollateral,
    }: {
      equity: string,
      freeCollateral: string,
    } = calculateEquityAndFreeCollateral(
      [perpetualPosition],
      perpetualIdToMarket,
      marketIdToMarket,
      usdcPositionSize,
    );

    expect(equity).toEqual('25000');
    expect(freeCollateral).toEqual('17500');
  });

  it('filterPositionsByLatestEventIdPerPerpetual', async () => {
    await testMocks.seedData();

    const perpetualPosition: PerpetualPositionFromDatabase = await
    PerpetualPositionTable.create({
      ...testConstants.defaultPerpetualPosition,
      lastEventId: defaultTendermintEventId,
      openEventId: defaultTendermintEventId,
    });

    const perpetualPosition2: PerpetualPositionFromDatabase = await
    PerpetualPositionTable.create({
      ...testConstants.defaultPerpetualPosition,
      perpetualId: defaultPerpetualMarket2.id,
      lastEventId: defaultTendermintEventId2,
      openEventId: defaultTendermintEventId2,
    });

    const perpetualPosition3: PerpetualPositionFromDatabase = await
    PerpetualPositionTable.create({
      ...testConstants.defaultPerpetualPosition,
      lastEventId: defaultTendermintEventId3,
      openEventId: defaultTendermintEventId3,
    });

    const filteredPerpetualPositions: PerpetualPositionFromDatabase[
    ] = filterPositionsByLatestEventIdPerPerpetual(
      initializePerpetualPositionsWithFunding([
        perpetualPosition,
        perpetualPosition2,
        perpetualPosition3,
      ]),
    );

    expect(filteredPerpetualPositions).toHaveLength(2);
    expect(filteredPerpetualPositions[0]).toEqual(
      expect.objectContaining({
        ...perpetualPosition3,
      }),
    );
    expect(filteredPerpetualPositions[1]).toEqual(
      expect.objectContaining({
        ...perpetualPosition2,
      }),
    );
  });

  it('maintenance fraction', async () => {
    const liquidityTierFromDatabase: LiquidityTiersFromDatabase = await
    LiquidityTiersTable.create(defaultLiquidityTier);
    expect(
      getMarginFraction(
        { liquidityTier: liquidityTierFromDatabase, initial: true },
      ),
    ).toEqual(Big('0.05'));
    expect(
      getMarginFraction(
        { liquidityTier: liquidityTierFromDatabase, initial: false },
      ),
    ).toEqual(Big('0.03'));
  });

  it('getSignedNotionalAndRisk', async () => {
    await LiquidityTiersTable.create(defaultLiquidityTier);
    await liquidityTierRefresher.updateLiquidityTiers();
    const perpetualMarketFromDatabase: PerpetualMarketFromDatabase = {
      ...defaultPerpetualMarket,
      id: '1',
    };
    const market: MarketFromDatabase = {
      ...defaultMarket,
      oraclePrice: '10000',
    };
    const bigSize: Big = Big('20');
    expect(
      getSignedNotionalAndRisk(
        { perpetualMarket: perpetualMarketFromDatabase, market, size: bigSize },
      ),
    ).toEqual(
      {
        signedNotional: Big('200000'),
        individualRisk: {
          initial: Big('10000'),
          maintenance: Big('6000'),
        },
      },
    );
  });

  describe('getFundingIndexMaps', () => {
    const fundingIndexUpdate3: FundingIndexUpdatesCreateObject = {
      ...testConstants.defaultFundingIndexUpdate,
      fundingIndex: '500',
      effectiveAtHeight: '3',
      eventId: testConstants.defaultTendermintEventId2,
    };

    it('returns FundingIndexMap', async () => {
      await testMocks.seedData();
      await BlockTable.create({
        ...testConstants.defaultBlock,
        blockHeight: '3',
      });
      await Promise.all([
        FundingIndexUpdatesTable.create(testConstants.defaultFundingIndexUpdate),
        FundingIndexUpdatesTable.create(fundingIndexUpdate3),
      ]);

      const [
        subaccount,
        latestBlock,
      ]: [
        SubaccountFromDatabase | undefined,
        BlockFromDatabase,
      ] = await Promise.all([
        SubaccountTable.findById(testConstants.defaultSubaccountId),
        BlockTable.getLatest(),
      ]);

      const {
        lastUpdatedFundingIndexMap,
        latestFundingIndexMap,
      }: {
        lastUpdatedFundingIndexMap: FundingIndexMap,
        latestFundingIndexMap: FundingIndexMap,
      } = await getFundingIndexMaps(
        subaccount!,
        latestBlock!,
      );

      expect(
        lastUpdatedFundingIndexMap[testConstants.defaultFundingIndexUpdate.perpetualId]
          .toString(),
      ).toEqual(testConstants.defaultFundingIndexUpdate.fundingIndex);
      expect(
        lastUpdatedFundingIndexMap[testConstants.defaultPerpetualMarket2.id]
          .toString(),
      ).toEqual(ZERO.toString());
      expect(
        lastUpdatedFundingIndexMap[testConstants.defaultPerpetualMarket3.id]
          .toString(),
      ).toEqual(ZERO.toString());
      expect(latestFundingIndexMap[fundingIndexUpdate3.perpetualId].toString())
        .toEqual(fundingIndexUpdate3.fundingIndex);
      expect(latestFundingIndexMap[testConstants.defaultPerpetualMarket2.id].toString())
        .toEqual(ZERO.toString());
      expect(latestFundingIndexMap[testConstants.defaultPerpetualMarket3.id].toString())
        .toEqual(ZERO.toString());
    });
  });

  describe('getTotalUnsettledFunding', () => {
    it('gets unsettled funding', async () => {
      await testMocks.seedData();

      const perpetualPosition: PerpetualPositionFromDatabase = await
      PerpetualPositionTable.create({
        ...testConstants.defaultPerpetualPosition,
        lastEventId: defaultTendermintEventId,
        openEventId: defaultTendermintEventId,
      });
      const perpetualPosition2: PerpetualPositionFromDatabase = await
      PerpetualPositionTable.create({
        ...testConstants.defaultPerpetualPosition,
        perpetualId: defaultPerpetualMarket2.id,
        lastEventId: defaultTendermintEventId2,
        openEventId: defaultTendermintEventId2,
      });

      const lastUpdatedFundingIndexMap: FundingIndexMap = {
        [perpetualPosition.perpetualId]: Big('100'),
        [perpetualPosition2.perpetualId]: Big('1000'),
      };
      const latestFundingIndexMap: FundingIndexMap = {
        [perpetualPosition.perpetualId]: Big('200'),
        [perpetualPosition2.perpetualId]: Big('2000'),
      };

      const unsettledFunding: Big = getTotalUnsettledFunding(
        [perpetualPosition, perpetualPosition2],
        latestFundingIndexMap,
        lastUpdatedFundingIndexMap,
      );

      expect(unsettledFunding).toEqual(
        Big(perpetualPosition.size).times('-100').plus(
          Big(perpetualPosition2.size).times('-1000'),
        ),
      );
    });
  });

  describe('adjustUSDCAssetPosition', () => {
    it.each([
      ['long', PositionSide.LONG, '1300', '1300'],
      ['short', PositionSide.SHORT, '700', '-700'],
    ])('adjusts USDC position size in returned map, size: [%s]', (
      _name: string,
      side: PositionSide,
      expectedPositionSize: string,
      expectedAdjustedPositionSize: string,
    ) => {
      const assetPositions: AssetPositionsMap = {
        [USDC_SYMBOL]: {
          ...ZERO_USDC_POSITION,
          side,
          size: '1000',
        },
        BTC: {
          symbol: 'BTC',
          side: PositionSide.LONG,
          assetId: '0',
          size: '1',
          subaccountNumber: 0,
        },
      };
      const unsettledFunding: Big = Big('300');

      const {
        assetPositionsMap,
        adjustedUSDCAssetPositionSize,
      }: {
        assetPositionsMap: AssetPositionsMap,
        adjustedUSDCAssetPositionSize: string,
      } = adjustUSDCAssetPosition(assetPositions, unsettledFunding);

      // Original asset positions object should be unchanged
      expect(assetPositions).toEqual({
        [USDC_SYMBOL]: {
          ...ZERO_USDC_POSITION,
          side,
          size: '1000',
        },
        BTC: {
          symbol: 'BTC',
          side: PositionSide.LONG,
          assetId: '0',
          size: '1',
          subaccountNumber: 0,
        },
      });
      expect(assetPositionsMap).toEqual({
        [USDC_SYMBOL]: {
          ...ZERO_USDC_POSITION,
          side,
          size: expectedPositionSize,
        },
        BTC: {
          symbol: 'BTC',
          side: PositionSide.LONG,
          assetId: '0',
          size: '1',
          subaccountNumber: 0,
        },
      });
      expect(adjustedUSDCAssetPositionSize).toEqual(expectedAdjustedPositionSize);
    });

    it.each([
      ['long', 'short', PositionSide.LONG, PositionSide.LONG, '300', '500', '800', '800'],
      ['short', 'long', PositionSide.SHORT, PositionSide.SHORT, '300', '-500', '800', '-800'],
    ])('flips USDC position side, original side [%s], flipped side [%s]', (
      _name: string,
      _secondName: string,
      side: PositionSide,
      expectedSide: PositionSide,
      positionSize: string,
      unsettledFunding: string,
      expectedPositionSize: string,
      expectedAdjustedPositionSize: string,
    ) => {
      const assetPositions: AssetPositionsMap = {
        [USDC_SYMBOL]: {
          ...ZERO_USDC_POSITION,
          side,
          size: positionSize,
        },
        BTC: {
          symbol: 'BTC',
          side: PositionSide.LONG,
          assetId: '0',
          size: '1',
          subaccountNumber: 0,
        },
      };

      const {
        assetPositionsMap,
        adjustedUSDCAssetPositionSize,
      }: {
        assetPositionsMap: AssetPositionsMap,
        adjustedUSDCAssetPositionSize: string,
      } = adjustUSDCAssetPosition(assetPositions, Big(unsettledFunding));

      // Original asset positions object should be unchanged
      expect(assetPositions).toEqual({
        [USDC_SYMBOL]: {
          ...ZERO_USDC_POSITION,
          side,
          size: positionSize,
        },
        BTC: {
          symbol: 'BTC',
          side: PositionSide.LONG,
          assetId: '0',
          size: '1',
          subaccountNumber: 0,
        },
      });
      expect(assetPositionsMap).toEqual({
        [USDC_SYMBOL]: {
          ...ZERO_USDC_POSITION,
          side: expectedSide,
          size: expectedPositionSize,
        },
        BTC: {
          symbol: 'BTC',
          side: PositionSide.LONG,
          assetId: '0',
          size: '1',
          subaccountNumber: 0,
        },
      });
      expect(adjustedUSDCAssetPositionSize).toEqual(expectedAdjustedPositionSize);
    });

    it.each([
      ['long', '300', PositionSide.LONG],
      ['short', '-300', PositionSide.SHORT],
    ])('adjusts USDC position when USDC position doesn\'t exist, side [%s]', (
      _name: string,
      funding: string,
      expectedSide: PositionSide,
    ) => {
      const assetPositions: AssetPositionsMap = {
        BTC: {
          symbol: 'BTC',
          side: PositionSide.LONG,
          assetId: '0',
          size: '1',
          subaccountNumber: 0,
        },
      };

      const {
        assetPositionsMap,
        adjustedUSDCAssetPositionSize,
      }: {
        assetPositionsMap: AssetPositionsMap,
        adjustedUSDCAssetPositionSize: string,
      } = adjustUSDCAssetPosition(assetPositions, Big(funding));

      // Original asset positions object should be unchanged
      expect(assetPositions).toEqual({
        BTC: {
          symbol: 'BTC',
          side: PositionSide.LONG,
          assetId: '0',
          size: '1',
          subaccountNumber: 0,
        },
      });
      expect(assetPositionsMap).toEqual({
        [USDC_SYMBOL]: {
          ...ZERO_USDC_POSITION,
          side: expectedSide,
          size: Big(funding).abs().toString(),
        },
        BTC: {
          symbol: 'BTC',
          side: PositionSide.LONG,
          assetId: '0',
          size: '1',
          subaccountNumber: 0,
        },
      });
      expect(adjustedUSDCAssetPositionSize).toEqual(funding);
    });

    it.each([
      ['long', PositionSide.LONG, '300', '-300'],
      ['short', PositionSide.SHORT, '300', '300'],
    ])('removes USDC position when resulting USDC position size is 0, side [%s]', (
      _name: string,
      side: PositionSide,
      positionSize: string,
      unsettledFunding: string,
    ) => {
      const assetPositions: AssetPositionsMap = {
        [USDC_SYMBOL]: {
          ...ZERO_USDC_POSITION,
          side,
          size: positionSize,
        },
        BTC: {
          symbol: 'BTC',
          side: PositionSide.LONG,
          assetId: '0',
          size: '1',
          subaccountNumber: 0,
        },
      };

      const {
        assetPositionsMap,
        adjustedUSDCAssetPositionSize,
      }: {
        assetPositionsMap: AssetPositionsMap,
        adjustedUSDCAssetPositionSize: string,
      } = adjustUSDCAssetPosition(assetPositions, Big(unsettledFunding));

      // Original asset positions object should be unchanged
      expect(assetPositions).toEqual({
        [USDC_SYMBOL]: {
          ...ZERO_USDC_POSITION,
          side,
          size: positionSize,
        },
        BTC: {
          symbol: 'BTC',
          side: PositionSide.LONG,
          assetId: '0',
          size: '1',
          subaccountNumber: 0,
        },
      });
      expect(assetPositionsMap).toEqual({
        BTC: {
          symbol: 'BTC',
          side: PositionSide.LONG,
          assetId: '0',
          size: '1',
          subaccountNumber: 0,
        },
      });
      expect(adjustedUSDCAssetPositionSize).toEqual(ZERO.toString());
    });
  });

  describe('getPerpetualPositionsWithUpdatedFunding', () => {
    let perpetualPosition: PerpetualPositionFromDatabase;
    let perpetualPosition2: PerpetualPositionFromDatabase;
    let lastUpdatedFundingIndexMap: FundingIndexMap;
    let latestFundingIndexMap: FundingIndexMap;

    beforeEach(async () => {
      await testMocks.seedData();

      perpetualPosition = await
      PerpetualPositionTable.create({
        ...testConstants.defaultPerpetualPosition,
        lastEventId: defaultTendermintEventId,
        openEventId: defaultTendermintEventId,
      });
      perpetualPosition2 = await
      PerpetualPositionTable.create({
        ...testConstants.defaultPerpetualPosition,
        perpetualId: defaultPerpetualMarket2.id,
        lastEventId: defaultTendermintEventId2,
        openEventId: defaultTendermintEventId2,
      });

      lastUpdatedFundingIndexMap = {
        [perpetualPosition.perpetualId]: Big('100'),
        [perpetualPosition2.perpetualId]: Big('1000'),
      };
      latestFundingIndexMap = {
        [perpetualPosition.perpetualId]: Big('200'),
        [perpetualPosition2.perpetualId]: Big('2000'),
      };
    });

    it('updates OPEN perpetual positions', () => {
      const updatedPerpetualPositions:
      PerpetualPositionWithFunding[] = getPerpetualPositionsWithUpdatedFunding(
        initializePerpetualPositionsWithFunding([perpetualPosition, perpetualPosition2]),
        latestFundingIndexMap,
        lastUpdatedFundingIndexMap,
      );

      expect(updatedPerpetualPositions[0].unsettledFunding).toEqual(
        helpers.getUnsettledFunding(
          perpetualPosition,
          latestFundingIndexMap,
          lastUpdatedFundingIndexMap,
        ).toFixed(),
      );

      expect(updatedPerpetualPositions[1].unsettledFunding).toEqual(
        helpers.getUnsettledFunding(
          perpetualPosition2,
          latestFundingIndexMap,
          lastUpdatedFundingIndexMap,
        ).toFixed(),
      );
    });

    it.each([
      [PerpetualPositionStatus.CLOSED],
      [PerpetualPositionStatus.LIQUIDATED],
    ])('does not modify positions with status %s', (
      status: PerpetualPositionStatus,
    ) => {
      const positionWithStatus: PerpetualPositionWithFunding = {
        ...perpetualPosition,
        status,
        unsettledFunding: '0',
      };

      const updatedPerpetualPositions:
      PerpetualPositionWithFunding[] = getPerpetualPositionsWithUpdatedFunding(
        [positionWithStatus],
        latestFundingIndexMap,
        lastUpdatedFundingIndexMap,
      );

      expect(updatedPerpetualPositions[0].unsettledFunding)
        .toEqual('0');
    });
  });

  describe('getChildSubaccountNums', () => {
    it('Gets a list of all possible child subaccount numbers for a parent subaccount 0', () => {
      const childSubaccounts = getChildSubaccountNums(0);
      expect(childSubaccounts.length).toEqual(1000);
      expect(childSubaccounts[0]).toEqual(0);
      expect(childSubaccounts[1]).toEqual(128);
      expect(childSubaccounts[999]).toEqual(128 * 999);
    });
    it('Gets a list of all possible child subaccount numbers for a parent subaccount 127', () => {
      const childSubaccounts = getChildSubaccountNums(127);
      expect(childSubaccounts.length).toEqual(1000);
      expect(childSubaccounts[0]).toEqual(127);
      expect(childSubaccounts[1]).toEqual(128 + 127);
      expect(childSubaccounts[999]).toEqual(128 * 999 + 127);
    });
  });

  describe('getChildSubaccountNums', () => {
    it('Throws an error if the parent subaccount number is greater than or equal to the maximum parent subaccount number', () => {
      expect(() => getChildSubaccountNums(128)).toThrowError('Parent subaccount number must be less than 128');
    });
  });

  describe('getSubaccountResponse', () => {
    it('gets subaccount response with adjusted perpetual positions', () => {
      // Helper function does not care about ids.
      const id: string = 'mock-id';
      const perpetualPositions: PerpetualPositionFromDatabase[] = [{
        ...testConstants.defaultPerpetualPosition,
        id,
        entryPrice: '20000',
        sumOpen: '10',
        sumClose: '0',
      }];
      const assetPositions: AssetPositionFromDatabase[] = [{
        ...testConstants.defaultAssetPosition,
        id,
      }];
      const lastUpdatedFundingIndexMap: FundingIndexMap = {
        0: Big('10000'),
        1: Big('0'),
        2: Big('0'),
        3: Big('0'),
        4: Big('0'),
      };
      const latestUpdatedFundingIndexMap: FundingIndexMap = {
        0: Big('10050'),
        1: Big('0'),
        2: Big('0'),
        3: Big('0'),
        4: Big('0'),
      };
      const assets: AssetFromDatabase[] = [{
        ...testConstants.defaultAsset,
        id: '0',
      }];
      const markets: MarketFromDatabase[] = [
        testConstants.defaultMarket,
      ];
      const subaccount: SubaccountFromDatabase = {
        ...testConstants.defaultSubaccount,
        id,
      };
      const perpetualMarketsMap: PerpetualMarketsMap = {
        0: {
          ...testConstants.defaultPerpetualMarket,
        },
      };

      const response: SubaccountResponseObject = getSubaccountResponse(
        subaccount,
        perpetualPositions,
        assetPositions,
        assets,
        markets,
        perpetualMarketsMap,
        '3',
        latestUpdatedFundingIndexMap,
        lastUpdatedFundingIndexMap,
      );

      expect(response).toEqual({
        address: testConstants.defaultAddress,
        subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        equity: getFixedRepresentation(159500),
        freeCollateral: getFixedRepresentation(152000),
        marginEnabled: true,
        updatedAtHeight: testConstants.defaultSubaccount.updatedAtHeight,
        latestProcessedBlockHeight: '3',
        openPerpetualPositions: {
          [testConstants.defaultPerpetualMarket.ticker]: {
            market: testConstants.defaultPerpetualMarket.ticker,
            size: testConstants.defaultPerpetualPosition.size,
            side: testConstants.defaultPerpetualPosition.side,
            entryPrice: getFixedRepresentation(
              testConstants.defaultPerpetualPosition.entryPrice!,
            ),
            maxSize: testConstants.defaultPerpetualPosition.maxSize,
            // 200000 + 10*(10000-10050)=199500
            netFunding: getFixedRepresentation('199500'),
            realizedPnl: getFixedRepresentation('100'),
            // size * (index-entry) = 10*(15000-20000) = -50000
            unrealizedPnl: getFixedRepresentation(-50000),
            status: testConstants.defaultPerpetualPosition.status,
            sumOpen: testConstants.defaultPerpetualPosition.sumOpen,
            sumClose: testConstants.defaultPerpetualPosition.sumClose,
            createdAt: testConstants.defaultPerpetualPosition.createdAt,
            createdAtHeight: testConstants.defaultPerpetualPosition.createdAtHeight,
            exitPrice: undefined,
            closedAt: undefined,
            subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
          },
        },
        assetPositions: {
          [testConstants.defaultAsset.symbol]: {
            symbol: testConstants.defaultAsset.symbol,
            size: '9500',
            side: PositionSide.LONG,
            assetId: testConstants.defaultAssetPosition.assetId,
            subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
          },
        },
      });
    });
  });

  describe('aggregateHourlyPnlTicks', () => {
    it('aggregates single pnl tick', () => {
      const pnlTick: PnlTicksFromDatabase = {
        ...testConstants.defaultPnlTick,
        id: PnlTicksTable.uuid(
          testConstants.defaultPnlTick.subaccountId,
          testConstants.defaultPnlTick.createdAt,
        ),
      };

      const aggregatedPnlTicks: AggregatedPnlTick[] = aggregateHourlyPnlTicks([pnlTick]);
      expect(
        aggregatedPnlTicks,
      ).toEqual(
        [expect.objectContaining(
          {
            pnlTick: expect.objectContaining(testConstants.defaultPnlTick),
            numTicks: 1,
          },
        )],
      );
    });

    it('aggregates multiple pnl ticks same height and de-dupes ticks', () => {
      const pnlTick: PnlTicksFromDatabase = {
        ...testConstants.defaultPnlTick,
        id: PnlTicksTable.uuid(
          testConstants.defaultPnlTick.subaccountId,
          testConstants.defaultPnlTick.createdAt,
        ),
      };
      const pnlTick2: PnlTicksFromDatabase = {
        ...testConstants.defaultPnlTick,
        subaccountId: testConstants.defaultSubaccountId2,
        id: PnlTicksTable.uuid(
          testConstants.defaultSubaccountId2,
          testConstants.defaultPnlTick.createdAt,
        ),
      };
      const blockHeight2: string = '80';
      const blockTime2: string = DateTime.fromISO(pnlTick.createdAt).plus({ hour: 1 }).toISO();
      const pnlTick3: PnlTicksFromDatabase = {
        ...testConstants.defaultPnlTick,
        id: PnlTicksTable.uuid(
          testConstants.defaultPnlTick.subaccountId,
          blockTime2,
        ),
        blockHeight: blockHeight2,
        blockTime: blockTime2,
        createdAt: blockTime2,
      };
      const blockHeight3: string = '81';
      const blockTime3: string = DateTime.fromISO(pnlTick.createdAt).plus({ minute: 61 }).toISO();
      const pnlTick4: PnlTicksFromDatabase = {
        ...testConstants.defaultPnlTick,
        subaccountId: testConstants.defaultSubaccountId2,
        id: PnlTicksTable.uuid(
          testConstants.defaultSubaccountId2,
          blockTime3,
        ),
        equity: '1',
        totalPnl: '2',
        netTransfers: '3',
        blockHeight: blockHeight3,
        blockTime: blockTime3,
        createdAt: blockTime3,
      };
      const blockHeight4: string = '82';
      const blockTime4: string = DateTime.fromISO(pnlTick.createdAt).startOf('hour').plus({ minute: 63 }).toISO();
      // should be de-duped
      const pnlTick5: PnlTicksFromDatabase = {
        ...testConstants.defaultPnlTick,
        subaccountId: testConstants.defaultSubaccountId2,
        id: PnlTicksTable.uuid(
          testConstants.defaultSubaccountId2,
          blockTime4,
        ),
        equity: '1',
        totalPnl: '2',
        netTransfers: '3',
        blockHeight: blockHeight4,
        blockTime: blockTime4,
        createdAt: blockTime4,
      };

      const aggregatedPnlTicks: AggregatedPnlTick[] = aggregateHourlyPnlTicks(
        [pnlTick, pnlTick2, pnlTick3, pnlTick4, pnlTick5],
      );
      expect(aggregatedPnlTicks).toEqual(
        expect.arrayContaining([
          // Combined pnl tick at initial hour
          expect.objectContaining({
            pnlTick: expect.objectContaining({
              equity: (parseFloat(testConstants.defaultPnlTick.equity) +
              parseFloat(pnlTick2.equity)).toString(),
              totalPnl: (parseFloat(testConstants.defaultPnlTick.totalPnl) +
                  parseFloat(pnlTick2.totalPnl)).toString(),
              netTransfers: (parseFloat(testConstants.defaultPnlTick.netTransfers) +
                  parseFloat(pnlTick2.netTransfers)).toString(),
            }),
            numTicks: 2,
          }),
          // Combined pnl tick at initial hour + 1 hour and initial hour + 1 hour, 1 minute
          expect.objectContaining({
            pnlTick: expect.objectContaining({
              equity: (parseFloat(pnlTick3.equity) +
              parseFloat(pnlTick4.equity)).toString(),
              totalPnl: (parseFloat(pnlTick3.totalPnl) +
                  parseFloat(pnlTick4.totalPnl)).toString(),
              netTransfers: (parseFloat(pnlTick3.netTransfers) +
                  parseFloat(pnlTick4.netTransfers)).toString(),
            }),
            numTicks: 2,
          }),
        ]),
      );
    });
  });
});
