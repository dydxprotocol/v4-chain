import {
  AssetCreateObject,
  LiquidityTiersCreateObject,
  MarketColumns,
  MarketFromDatabase,
  MarketsMap,
  PerpetualPositionFromDatabase,
  PositionSide,
  SubaccountFromDatabase,
} from '../../src/types';
import {
  AssetPositionCreateObjectWithId,
  getAssetCreateObject,
  getAssetPositionCreateObject,
  getLiquidityTiersCreateObject,
  getMaintenanceMarginPpm,
  getUnrealizedPnl,
  getUnsettledFunding,
  SubaccountCreateObjectWithId,
} from '../../src/db/helpers';
import {
  Asset, AssetPosition, LiquidityTier, MarketParam, MarketPrice,
} from '@dydxprotocol-indexer/v4-protos';
import { bigIntToBytes } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  expectAsset,
  expectAssetPositionCreateObject,
  expectLiquidityTier,
  expectMarketParamAndPrice,
  expectSubaccount,
} from './helpers';
import {
  createdDateTime,
  createdHeight,
  defaultAddress,
  defaultFundingIndexUpdate,
  defaultPerpetualMarket,
  defaultPerpetualPosition,
  defaultPerpetualPositionId,
  defaultSubaccount,
  defaultSubaccountId,
} from '../helpers/constants';
import * as SubaccountTable from '../../src/stores/subaccount-table';
import * as PerpetualPositionTable from '../../src/stores/perpetual-position-table';
import * as MarketTable from '../../src/stores/market-table';
import { CURRENCY_DECIMAL_PRECISION, USDC_DENOM, USDC_SYMBOL } from '../../src';
import Long from 'long';
import Big from 'big.js';
import { seedData } from '../helpers/mock-generators';
import _ from 'lodash';
import { clearData } from '../../src/helpers/db-helpers';

describe('helpers', () => {

  afterEach(async () => {
    await clearData();
  });

  describe('getAssetCreateObjects', () => {
    const defaultAsset: Asset = {
      id: 0,
      symbol: USDC_SYMBOL,
      denom: USDC_DENOM,
      denomExponent: -6,
      atomicResolution: -2,
      hasMarket: false,
      marketId: 0,
      longInterest: Long.fromValue(1000),
    };

    it('create asset object from asset proto', () => {
      const assetCreateObject: AssetCreateObject = getAssetCreateObject(defaultAsset);
      expectAsset(assetCreateObject, defaultAsset);
    });
  });

  describe('getLiquidityTiersCreateObjects', () => {
    const defaultLiquidityTier: LiquidityTier = {
      id: 0,
      name: 'tier1',
      initialMarginPpm: 50000,
      maintenanceFractionPpm: 30000,
      basePositionNotional: Long.fromValue(10000),
      impactNotional: Long.fromValue(500),
    };

    it('create LiquidityTiers object from LiquidityTiers proto', () => {
      const liquidityTiersCreateObject:
      LiquidityTiersCreateObject = getLiquidityTiersCreateObject(defaultLiquidityTier);
      expectLiquidityTier(liquidityTiersCreateObject, defaultLiquidityTier);
    });
  });

  describe('getAssetPositionCreateObjects', () => {
    const atomicResolution: number = 2;
    const defaultAssetPosition: AssetPosition = {
      assetId: 0,
      quantums: bigIntToBytes(BigInt(1000)),
      index: Long.fromValue(0),
    };

    it('create asset position object from AssetPosition proto', () => {
      const assetPositionCreateObject:
      AssetPositionCreateObjectWithId = getAssetPositionCreateObject(
        defaultSubaccountId, defaultAssetPosition, atomicResolution);
      expectAssetPositionCreateObject(
        assetPositionCreateObject, defaultAssetPosition, atomicResolution);
    });
  });

  describe('expectSubaccount', () => {
    const subaccount: SubaccountFromDatabase = {
      id: defaultSubaccountId,
      address: defaultAddress,
      subaccountNumber: 0,
      updatedAt: createdDateTime.toISO(),
      updatedAtHeight: createdHeight,
    };
    const subaccountCreateObject: SubaccountCreateObjectWithId = {
      id: SubaccountTable.uuid(defaultSubaccount.address, defaultSubaccount.subaccountNumber),
      ...defaultSubaccount,
    };

    it('expect subaccount', () => {
      expectSubaccount(subaccount, subaccountCreateObject);
    });
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
      ).toEqual(Big('20000'));  // 10 * (12050-10050)
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
      ).toEqual(Big('-20000'));  // -10 * (12050-10050)
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
      ).toEqual(Big('2700.1674'));  // 1.35 * (12050.124-10050)
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
        perpetualPosition, defaultPerpetualMarket, marketIdToMarket,
      );

      expect(unrealizedPnl).toEqual(Big(-50000).toFixed(CURRENCY_DECIMAL_PRECISION));
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
        perpetualPosition, defaultPerpetualMarket, marketIdToMarket,
      );

      expect(unrealizedPnl).toEqual(Big(50000).toFixed(CURRENCY_DECIMAL_PRECISION));
    });
  });

  describe('expectMarketParamAndPrice', () => {

    it('expect market', () => {
      const marketFromDb: MarketFromDatabase = {
        id: 0,
        pair: 'BTC-USD',
        exponent: -5,
        minPriceChangePpm: 50,
        oraclePrice: '50000',
      };
      const marketParam: MarketParam = {
        id: 0,
        pair: 'BTC-USD',
        exponent: -5,
        exchangeConfigJson: '{exchanges:[{"exchangeName":"Binance","ticker":"BTCUSDT"},{"exchangeName":"BinanceUS","ticker":"BTCUSD"}]}',
        minExchanges: 1,
        minPriceChangePpm: 50,
      };
      const marketPrice: MarketPrice = {
        id: 0,
        exponent: -5,
        price: Long.fromValue(5_000_000_000),
      };
      expectMarketParamAndPrice(marketFromDb, marketParam, marketPrice);
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
