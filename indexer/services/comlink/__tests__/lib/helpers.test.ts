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
  TDAI_SYMBOL,
  PositionSide,
  helpers,
  PerpetualPositionStatus,
  LiquidityTiersFromDatabase,
  LiquidityTiersTable,
  liquidityTierRefresher,
} from '@dydxprotocol-indexer/postgres';
import {
  adjustTDAIAssetPosition,
  calculateEquityAndFreeCollateral,
  filterAssetPositions,
  filterPositionsByLatestEventIdPerPerpetual,
  getFundingIndexMaps,
  getMarginFraction,
  getSignedNotionalAndRisk,
  getTotalUnsettledFunding,
  getPerpetualPositionsWithUpdatedFunding,
  getChildSubaccountNums,
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
import { ZERO, ZERO_TDAI_POSITION } from '../../src/lib/constants';

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

    const tdaiPositionSize: string = '175000';

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
      tdaiPositionSize,
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

    const tdaiPositionSize: string = '175000';

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
      tdaiPositionSize,
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
      spotPrice: '10000',
      pnlPrice: '10000',
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

  describe('adjustTDAIAssetPosition', () => {
    it.each([
      ['long', PositionSide.LONG, '1300', '1300'],
      ['short', PositionSide.SHORT, '700', '-700'],
    ])('adjusts TDAI position size in returned map, size: [%s]', (
      _name: string,
      side: PositionSide,
      expectedPositionSize: string,
      expectedAdjustedPositionSize: string,
    ) => {
      const assetPositions: AssetPositionsMap = {
        [TDAI_SYMBOL]: {
          ...ZERO_TDAI_POSITION,
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
        adjustedTDAIAssetPositionSize,
      }: {
        assetPositionsMap: AssetPositionsMap,
        adjustedTDAIAssetPositionSize: string
      } = adjustTDAIAssetPosition(assetPositions, unsettledFunding);

      // Original asset positions object should be unchanged
      expect(assetPositions).toEqual({
        [TDAI_SYMBOL]: {
          ...ZERO_TDAI_POSITION,
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
        [TDAI_SYMBOL]: {
          ...ZERO_TDAI_POSITION,
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
      expect(adjustedTDAIAssetPositionSize).toEqual(expectedAdjustedPositionSize);
    });

    it.each([
      ['long', 'short', PositionSide.LONG, PositionSide.LONG, '300', '500', '800', '800'],
      ['short', 'long', PositionSide.SHORT, PositionSide.SHORT, '300', '-500', '800', '-800'],
    ])('flips TDAI position side, original side [%s], flipped side [%s]', (
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
        [TDAI_SYMBOL]: {
          ...ZERO_TDAI_POSITION,
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
        adjustedTDAIAssetPositionSize,
      }: {
        assetPositionsMap: AssetPositionsMap,
        adjustedTDAIAssetPositionSize: string
      } = adjustTDAIAssetPosition(assetPositions, Big(unsettledFunding));

      // Original asset positions object should be unchanged
      expect(assetPositions).toEqual({
        [TDAI_SYMBOL]: {
          ...ZERO_TDAI_POSITION,
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
        [TDAI_SYMBOL]: {
          ...ZERO_TDAI_POSITION,
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
      expect(adjustedTDAIAssetPositionSize).toEqual(expectedAdjustedPositionSize);
    });

    it.each([
      ['long', '300', PositionSide.LONG],
      ['short', '-300', PositionSide.SHORT],
    ])('adjusts TDAI position when TDAI position doesn\'t exist, side [%s]', (
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
        adjustedTDAIAssetPositionSize,
      }: {
        assetPositionsMap: AssetPositionsMap,
        adjustedTDAIAssetPositionSize: string
      } = adjustTDAIAssetPosition(assetPositions, Big(funding));

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
        [TDAI_SYMBOL]: {
          ...ZERO_TDAI_POSITION,
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
      expect(adjustedTDAIAssetPositionSize).toEqual(funding);
    });

    it.each([
      ['long', PositionSide.LONG, '300', '-300'],
      ['short', PositionSide.SHORT, '300', '300'],
    ])('removes TDAI position when resulting TDAI position size is 0, side [%s]', (
      _name: string,
      side: PositionSide,
      positionSize: string,
      unsettledFunding: string,
    ) => {
      const assetPositions: AssetPositionsMap = {
        [TDAI_SYMBOL]: {
          ...ZERO_TDAI_POSITION,
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
        adjustedTDAIAssetPositionSize,
      }: {
        assetPositionsMap: AssetPositionsMap,
        adjustedTDAIAssetPositionSize: string
      } = adjustTDAIAssetPosition(assetPositions, Big(unsettledFunding));

      // Original asset positions object should be unchanged
      expect(assetPositions).toEqual({
        [TDAI_SYMBOL]: {
          ...ZERO_TDAI_POSITION,
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
      expect(adjustedTDAIAssetPositionSize).toEqual(ZERO.toString());
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
});
