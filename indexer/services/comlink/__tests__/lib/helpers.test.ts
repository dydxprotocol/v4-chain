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
} from '@dydxprotocol-indexer/postgres';
import {
  adjustUSDCAssetPosition,
  calculateEquityAndFreeCollateral,
  filterAssetPositions,
  filterPositionsByLatestEventIdPerPerpetual,
  getFundingIndexMaps,
  getAdjustedMarginFraction,
  getSignedNotionalAndRisk,
  getTotalUnsettledFunding,
  getPerpetualPositionsWithUpdatedFunding,
  initializePerpetualPositionsWithFunding,
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
import { AssetPositionsMap, PerpetualPositionWithFunding } from '../../src/types';
import { ZERO, ZERO_USDC_POSITION } from '../../src/lib/constants';

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
    ] = await filterPositionsByLatestEventIdPerPerpetual(
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

  it.each([
    ['less than base', 20, 0.05],
    ['base', 1_000_000, 0.05],
    ['greater than base', 4_000_000, 0.1],
    ['max', 400_000_000, 1],
    ['greater than max', 4_000_000_000, 1],
  ])('getAdjustedInitialMarginFraction: %s', async (
    _name: string,
    notionalValue: number,
    expectedResult: number,
  ) => {
    const liquidityTierFromDatabase: LiquidityTiersFromDatabase = await
    LiquidityTiersTable.create(defaultLiquidityTier);
    const positionNotional: Big = Big(notionalValue);
    expect(
      getAdjustedMarginFraction(
        { liquidityTier: liquidityTierFromDatabase, positionNotional, initial: true },
      ),
    ).toEqual(Big(expectedResult));
  });

  it.each([
    ['less than base', 20, 0.03],
    ['base', 1_000_000, 0.03],
    ['greater than base', 4_000_000, 0.06],
    ['greater than max', 4_000_000_000, 1],
  ])('getAdjustedMaintenanceMarginFraction: %s', async (
    _name: string,
    notionalValue: number,
    expectedResult: number,
  ) => {
    const liquidityTierFromDatabase: LiquidityTiersFromDatabase = await
    LiquidityTiersTable.create(defaultLiquidityTier);
    const positionNotional: Big = Big(notionalValue);
    expect(
      getAdjustedMarginFraction(
        { liquidityTier: liquidityTierFromDatabase, positionNotional, initial: false },
      ),
    ).toEqual(Big(expectedResult));
  });

  it.each([
    ['less than base', 20, 200_000, 10_000, 6_000],
    ['base', 100, 1_000_000, 50_000, 30_000],
    ['greater than base', 400, 4_000_000, 400_000, 240_000],
    ['max', 400_000, 4_000_000_000, 4_000_000_000, 4_000_000_000],
    ['less than base SHORT', -20, -200_000, 10_000, 6_000],
    ['base SHORT', -100, -1_000_000, 50_000, 30_000],
    ['greater than base SHORT', -400, -4_000_000, 400_000, 240_000],
    ['max SHORT', -400_000, -4_000_000_000, 4_000_000_000, 4_000_000_000],
  ])('getSignedNotionalAndRisk: %s', async (
    _name: string,
    size: number,
    signedNotional: number,
    initial: number,
    maintenance: number,
  ) => {
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
    const bigSize: Big = Big(size);
    expect(
      getSignedNotionalAndRisk(
        { perpetualMarket: perpetualMarketFromDatabase, market, size: bigSize },
      ),
    ).toEqual(
      {
        signedNotional: Big(signedNotional),
        individualRisk: {
          initial: Big(initial),
          maintenance: Big(maintenance),
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
        BlockFromDatabase | undefined,
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

      expect(Object.keys(lastUpdatedFundingIndexMap)).toHaveLength(3);
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
      expect(Object.keys(latestFundingIndexMap)).toHaveLength(3);
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
        },
      };
      const unsettledFunding: Big = Big('300');

      const {
        assetPositionsMap,
        adjustedUSDCAssetPositionSize,
      }: {
        assetPositionsMap: AssetPositionsMap,
        adjustedUSDCAssetPositionSize: string
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
        },
      };

      const {
        assetPositionsMap,
        adjustedUSDCAssetPositionSize,
      }: {
        assetPositionsMap: AssetPositionsMap,
        adjustedUSDCAssetPositionSize: string
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
        },
      };

      const {
        assetPositionsMap,
        adjustedUSDCAssetPositionSize,
      }: {
        assetPositionsMap: AssetPositionsMap,
        adjustedUSDCAssetPositionSize: string
      } = adjustUSDCAssetPosition(assetPositions, Big(funding));

      // Original asset positions object should be unchanged
      expect(assetPositions).toEqual({
        BTC: {
          symbol: 'BTC',
          side: PositionSide.LONG,
          assetId: '0',
          size: '1',
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
        },
      };

      const {
        assetPositionsMap,
        adjustedUSDCAssetPositionSize,
      }: {
        assetPositionsMap: AssetPositionsMap,
        adjustedUSDCAssetPositionSize: string
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
        },
      });
      expect(assetPositionsMap).toEqual({
        BTC: {
          symbol: 'BTC',
          side: PositionSide.LONG,
          assetId: '0',
          size: '1',
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
});
