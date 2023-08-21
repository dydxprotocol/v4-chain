import { logger } from '@dydxprotocol-indexer/base';
import { bigIntToBytes, bytesToBigInt, getPositionIsLong } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  Asset,
  AssetPosition,
  ClobPair,
  ClobPair_Status,
  LiquidityTier,
  MarketParam,
  MarketPrice,
  Perpetual,
} from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';
import Long from 'long';
import { DateTime } from 'luxon';

import { CURRENCY_DECIMAL_PRECISION, ONE_MILLION, QUOTE_CURRENCY_ATOMIC_RESOLUTION } from '../constants';
import { setBulkRowsForUpdate } from '../helpers/stores-helpers';
import {
  InvalidClobPairStatusError,
  LiquidityTierDoesNotExistError,
  MarketDoesNotExistError,
  PerpetualDoesNotExistError,
} from '../lib/errors';
import { protocolPriceToHuman, quantumsToHuman, quantumsToHumanFixedString } from '../lib/protocol-translations';
import * as AssetPositionTable from '../stores/asset-position-table';
import * as SubaccountTable from '../stores/subaccount-table';
import {
  AssetColumns,
  AssetCreateObject,
  BlockColumns,
  BlockCreateObject,
  FundingIndexMap,
  IsoString,
  LiquidityTiersColumns,
  LiquidityTiersCreateObject,
  MarketColumns,
  MarketCreateObject,
  MarketsMap,
  PerpetualMarketColumns,
  PerpetualMarketCreateObject,
  PerpetualMarketFromDatabase,
  PerpetualMarketStatus,
  PerpetualPositionFromDatabase,
  TransferFromDatabase,
  TransferType,
  UpdatedPerpetualPositionSubaccountKafkaObject,
} from '../types';
import genesis from './genesis.json';

export interface AssetPositionCreateObjectWithId {
  id: string,
  subaccountId: string,
  assetId: string,
  size: string,
  isLong: boolean,
}

export interface SubaccountCreateObjectWithId {
  id: string,
  address: string,
  subaccountNumber: number,
  updatedAt: IsoString,
  updatedAtHeight: string,
}

type SpecifiedClobPairStatus = Exclude<ClobPair_Status, ClobPair_Status.STATUS_UNSPECIFIED> &
Exclude<ClobPair_Status, ClobPair_Status.UNRECOGNIZED>;

const CLOB_STATUS_TO_MARKET_STATUS: Record<SpecifiedClobPairStatus, PerpetualMarketStatus> = {
  [ClobPair_Status.STATUS_ACTIVE]: PerpetualMarketStatus.ACTIVE,
  [ClobPair_Status.STATUS_CANCEL_ONLY]: PerpetualMarketStatus.CANCEL_ONLY,
  [ClobPair_Status.STATUS_PAUSED]: PerpetualMarketStatus.PAUSED,
  [ClobPair_Status.STATUS_POST_ONLY]: PerpetualMarketStatus.POST_ONLY,
};

/**
 * @description Gets the SQL to update the margin columns in `perpetual_markets` table,
 * using the `genesis.json` file from the V4 network.
 * This is only used in case the migration to liquidity tiers fails.
 *
 * @param id
 * @param initialMarginFraction
 * @param maintenanceMarginFraction
 */
export function updatePerpetualMarketMarginColumns(
  id: string,
  initialMarginFraction: string,
  maintenanceMarginFraction: string,
): string {
  return `UPDATE PERPETUAL_MARKETS
          SET "initialMarginFraction" = ${initialMarginFraction},
              "maintenanceMarginFraction" = ${maintenanceMarginFraction}
          WHERE "id" = '${id}'`;
}

export function getPerpetualMarketMarginRestoreSql(): string[] {
  const liquidityTierMap: Record<string, LiquidityTier> = _.keyBy(
    getLiquidityTiersFromGenesis(),
    (liquidityTier: LiquidityTier) => liquidityTier.id.toString(),
  );

  const marginMap: Record<string, {
    initialMarginFraction: string,
    maintenanceMarginFraction: string,
  }> = getPerpetualMarketToMarginMapping(liquidityTierMap);
  const sql: string[] = [];
  Object.keys(marginMap).forEach((id: string) => {
    sql.push(updatePerpetualMarketMarginColumns(
      id,
      marginMap[id].initialMarginFraction,
      marginMap[id].maintenanceMarginFraction,
    ));
  });
  return sql;
}

/**
 * @description Gets the SQL to update the liquidityTierId column in `perpetual_markets` table,
 * using the `genesis.json` file from the V4 network.
 *
 * @param marketId
 * @param liquidityTierId
 */
export function updatePerpetualMarketLiquidityTier(id: string, liquidityTierId: number): string {
  return `UPDATE PERPETUAL_MARKETS
          SET "liquidityTierId" = ${liquidityTierId}
          WHERE "id" = '${id}'`;
}

function getPerpetualMarketToLiquidityTierMapping(): Record<string, number> {
  const liquidityTierMap: Record<string, number> = {};
  const perpetuals: Perpetual[] = getPerpetualsFromGenesis();
  perpetuals.forEach((perpetual: Perpetual) => {
    liquidityTierMap[perpetual.id.toString()] = perpetual.liquidityTier;
  });
  return liquidityTierMap;
}

/**
 * Gets a map of perpetual market IDs to margin fractions.
 *
 * @param liquidityTiersMap
 */
function getPerpetualMarketToMarginMapping(
  liquidityTiersMap: Record<string, LiquidityTier>,
): Record<string, {
  initialMarginFraction: string,
  maintenanceMarginFraction: string,
}> {
  const marginMap: Record<string, {
    initialMarginFraction: string,
    maintenanceMarginFraction: string,
  }> = {};
  const perpetuals: Perpetual[] = getPerpetualsFromGenesis();
  perpetuals.forEach((perpetual: Perpetual) => {
    marginMap[perpetual.id.toString()] = {
      initialMarginFraction: ppmToString(
        liquidityTiersMap[perpetual.liquidityTier].initialMarginPpm,
      ),
      maintenanceMarginFraction: ppmToString(
        getMaintenanceMarginPpm(
          liquidityTiersMap[perpetual.liquidityTier].initialMarginPpm,
          liquidityTiersMap[perpetual.liquidityTier].maintenanceFractionPpm,
        ),
      ),
    };
  });
  return marginMap;
}

export function getPerpetualMarketLiquidityTierUpdateSql(): string[] {
  const liquidityTierMap: Record<string, number> = getPerpetualMarketToLiquidityTierMapping();
  const sql: string[] = [];
  Object.keys(liquidityTierMap).forEach((id: string) => {
    sql.push(updatePerpetualMarketLiquidityTier(id, liquidityTierMap[id]));
  });
  return sql;
}

/**
 * @description Gets the SQL to seed the `perpetual_markets` table, using the `genesis.json` file
 * from the V4 network.
 * @returns SQL statement for seeding the `perpetual_markets` table. The SQL statement will do
 * nothing if the rows in the `perpetual_markets` table already exist.
 */
export function getSeedPerpetualMarketsSql(): string {
  // Get `Perpetual` objects from the genesis app state
  const perpetuals: Perpetual[] = getPerpetualsFromGenesis();

  // Get `ClobPair` objects from the genesis app state
  const clobPairs: ClobPair[] = getClobPairsFromGenesis();

  // Get `MarketParam` objects from the genesis app state
  const marketParams: MarketParam[] = getMarketParamsFromGenesis();

  const liquidityTiers: LiquidityTier[] = getLiquidityTiersFromGenesis();

  const perpetualMarkets: PerpetualMarketCreateObject[] = getPerpetualMarketCreateObjects(
    clobPairs,
    perpetuals,
    marketParams,
    liquidityTiers,
  );
  const perpetualMarketColumns = _.keys(perpetualMarkets[0]) as PerpetualMarketColumns[];

  const marketRows: string[] = setBulkRowsForUpdate<PerpetualMarketColumns>({
    objectArray: perpetualMarkets,
    columns: perpetualMarketColumns,
    stringColumns: [
      PerpetualMarketColumns.status,
      PerpetualMarketColumns.ticker,
      PerpetualMarketColumns.marketId,
      PerpetualMarketColumns.baseAsset,
      PerpetualMarketColumns.quoteAsset,
    ],
    numericColumns: [
      PerpetualMarketColumns.atomicResolution,
      PerpetualMarketColumns.basePositionSize,
      PerpetualMarketColumns.clobPairId,
      PerpetualMarketColumns.id,
      PerpetualMarketColumns.incrementalPositionSize,
      PerpetualMarketColumns.lastPrice,
      PerpetualMarketColumns.maxPositionSize,
      PerpetualMarketColumns.minOrderBaseQuantums,
      PerpetualMarketColumns.nextFundingRate,
      PerpetualMarketColumns.openInterest,
      PerpetualMarketColumns.priceChange24H,
      PerpetualMarketColumns.quantumConversionExponent,
      PerpetualMarketColumns.stepBaseQuantums,
      PerpetualMarketColumns.subticksPerTick,
      PerpetualMarketColumns.trades24H,
      PerpetualMarketColumns.volume24H,
      PerpetualMarketColumns.liquidityTierId,
    ],
  });

  return `INSERT INTO PERPETUAL_MARKETS (${perpetualMarketColumns.map((col) => `"${col}"`).join(',')})
          VALUES ${marketRows.map((market) => `(${market})`).join(', ')}
          ON CONFLICT DO NOTHING`;
}

/**
 * @description Gets the SQL to seed the `markets` table, using the `genesis.json` file
 * from the V4 network.
 * @returns SQL statement for seeding the `markets` table. The SQL statement will do
 * nothing if the rows in the `markets` table already exist.
 */
export function getSeedMarketsSql(): string {
  // Get `MarketParam` and `MarketPrice` objects from the genesis app state
  const marketParams: MarketParam[] = getMarketParamsFromGenesis();
  const marketPrices: MarketPrice[] = getMarketPricesFromGenesis();

  const marketCreateObjects: MarketCreateObject[] = [];
  marketParams.forEach((marketParam, index) => {
    marketCreateObjects.push(getMarketCreateObject(marketParam, marketPrices[index]));
  });

  const marketColumns = _.keys(marketCreateObjects[0]) as MarketColumns[];

  const marketRows: string[] = setBulkRowsForUpdate<MarketColumns>({
    objectArray: marketCreateObjects,
    columns: marketColumns,
    stringColumns: [
      MarketColumns.pair,
      MarketColumns.oraclePrice,
    ],
    numericColumns: [
      MarketColumns.id,
      MarketColumns.exponent,
      MarketColumns.minPriceChangePpm,
    ],
  });

  return `INSERT INTO MARKETS (${marketColumns.map((col) => `"${col}"`).join(',')})
          VALUES ${marketRows.map((market) => `(${market})`).join(', ')}
          ON CONFLICT DO NOTHING`;
}

/**
 * @description Gets the SQL to seed the `assets` table, using the `genesis.json` file
 * from the V4 network.
 * @returns SQL statement for seeding the `assets` table. The SQL statement will do
 * nothing if the rows in the `assets` table already exist.
 */
export function getSeedAssetsSql(): string {
  // Get `Asset` objects from the genesis app state
  const assets: Asset[] = getAssetsFromGenesis();

  const assetCreateObjects: AssetCreateObject[] = assets.map((asset: Asset) => {
    return getAssetCreateObject(asset);
  });

  const assetColumns = _.keys(assetCreateObjects[0]) as AssetColumns[];

  const assetRows: string[] = setBulkRowsForUpdate<AssetColumns>({
    objectArray: assetCreateObjects,
    columns: assetColumns,
    stringColumns: [
      AssetColumns.symbol,
      AssetColumns.id,
    ],
    numericColumns: [
      AssetColumns.atomicResolution,
      AssetColumns.marketId,
    ],
    booleanColumns: [
      AssetColumns.hasMarket,
    ],
  });

  return `INSERT INTO ASSETS (${assetColumns.map((col) => `"${col}"`).join(',')})
          VALUES ${assetRows.map((asset) => `(${asset})`).join(', ')}
          ON CONFLICT DO NOTHING`;
}

/**
 * @description Gets the SQL to seed the `liquidity_tiers` table, using the `genesis.json` file
 * from the V4 network.
 * @returns SQL statement for seeding the `liquidity_tiers` table. The SQL statement will do
 * nothing if the rows in the `liquidity_tiers` table already exist.
 */
export function getSeedLiquidityTiersSql(): string {
  // Get `LiquidityTier` objects from the genesis app state
  const liquidityTiers: LiquidityTier[] = getLiquidityTiersFromGenesis();

  const liquidityTierCreateObjects:
  LiquidityTiersCreateObject[] = liquidityTiers.map((liquidityTier: LiquidityTier) => {
    return getLiquidityTiersCreateObject(liquidityTier);
  });

  const liquidityTierColumns = _.keys(liquidityTierCreateObjects[0]) as LiquidityTiersColumns[];

  const liquidityTierRows: string[] = setBulkRowsForUpdate<LiquidityTiersColumns>({
    objectArray: liquidityTierCreateObjects,
    columns: liquidityTierColumns,
    stringColumns: [
      LiquidityTiersColumns.name,
      LiquidityTiersColumns.basePositionNotional,
    ],
    numericColumns: [
      LiquidityTiersColumns.id,
      LiquidityTiersColumns.initialMarginPpm,
      LiquidityTiersColumns.maintenanceFractionPpm,
    ],
  });

  return `INSERT INTO LIQUIDITY_TIERS (${liquidityTierColumns.map((col) => `"${col}"`).join(',')})
          VALUES ${liquidityTierRows.map((liquidityTier) => `(${liquidityTier})`).join(', ')}
          ON CONFLICT DO NOTHING`;
}

/**
 * @description Gets the SQL to seed the `blocks` table.
 *
 * This needs to be run before the seeding of any tables that have blockHeight foreign key.
 *
 * @returns SQL statement for seeding the `blocks` table. The SQL statement will do
 * nothing if the rows in the `blocks` table already exist.
 */
export function getSeedBlocksSql() {
  const blockCreateObjects:
  BlockCreateObject[] = [{
    // Setting block height to -1, so the genesis.json initialized structs from the full node will
    // be block 0
    blockHeight: '-1',
    time: DateTime.utc().toISO(),
  }];

  const blockColumns = _.keys(blockCreateObjects[0]) as BlockColumns[];

  const blockRows: string[] = setBulkRowsForUpdate<BlockColumns>({
    objectArray: blockCreateObjects,
    columns: blockColumns,
    timestampColumns: [
      BlockColumns.time,
    ],
    numericColumns: [
      BlockColumns.blockHeight,
    ],
  });

  return `INSERT INTO BLOCKS (${blockColumns.map((col) => `"${col}"`).join(',')})
          VALUES ${blockRows.map((block) => `(${block})`).join(', ')}
          ON CONFLICT DO NOTHING`;
}

/**
 * @description Gets `Perpetual` objects from genesis state.
 * @returns
 */
export function getPerpetualsFromGenesis(): Perpetual[] {
  // Get `Perpetual` objects from the genesis app state
  const perpetuals: Perpetual[] = genesis.app_state.perpetuals.perpetuals.map(
    (genesisPerpetual, index: number): Perpetual => {
      return {
        ...genesisPerpetual,
        marketId: genesisPerpetual.market_id,
        atomicResolution: genesisPerpetual.atomic_resolution,
        defaultFundingPpm: genesisPerpetual.default_funding_ppm,
        liquidityTier: genesisPerpetual.liquidity_tier,
        // Reference https://github.com/dydxprotocol/v4/blob/main/x/perpetuals/keeper/perpetual.go#L34
        // Id for each perpetual is the number of perpetuals, which is equal to the index of the
        // perpetual in the array.
        id: index,
        // Reference https://github.com/dydxprotocol/v4/blob/main/x/perpetuals/keeper/perpetual.go#L43
        // Currently both fields are set to 0 on V4 when a Perpetual is created.
        fundingIndex: bigIntToBytes(BigInt(0)),
        openInterest: Long.fromValue(0),
      };
    },
  );

  return perpetuals;
}

/**
 * @description Gets `Asset` objects from genesis state.
 * @returns
 */
export function getAssetsFromGenesis(): Asset[] {
  // Get `Asset` objects from the genesis app state
  const assets: Asset[] = genesis.app_state.assets.assets.map(
    (asset): Asset => {
      return {
        atomicResolution: asset.atomic_resolution,
        symbol: asset.symbol,
        denom: asset.denom,
        denomExponent: parseInt(asset.denom_exponent, 10),
        hasMarket: asset.has_market,
        id: asset.id,
        marketId: asset.market_id,
        longInterest: Long.fromValue(asset.long_interest),
      };
    },
  );

  return assets;
}

/**
 * @description Gets `SubaccountCreateObject` objects from genesis state.
 *
 * @returns a list of SubaccountCreateObjects
 */
export function getSubaccountCreateObjectsFromGenesis(): SubaccountCreateObjectWithId[] {
  const subaccountCreateObjects: SubaccountCreateObjectWithId[] = [];
  // Get SubaccountCreateObjects from the genesis app state
  genesis.app_state.subaccounts.subaccounts
    .forEach(
      (subaccount) => {
        subaccountCreateObjects.push({
          id: SubaccountTable.uuid(subaccount.id.owner, subaccount.id.number),
          address: subaccount.id.owner,
          subaccountNumber: subaccount.id.number,
          updatedAt: DateTime.utc().toISO(),
          updatedAtHeight: '1',
        });
      });

  return subaccountCreateObjects;
}

/**
 * @description Gets `AssetPosition` objects from genesis state.
 * @returns a map from subaccountId to its asset positions
 */
export function getAssetPositionsFromGenesis(): _.Dictionary<AssetPosition[]> {
  const assetPositionMapping: _.Dictionary<AssetPosition[]> = {};
  // Get `AssetPosition` objects from the genesis app state
  genesis.app_state.subaccounts.subaccounts
    .forEach(
      (subaccount) => {
        const subaccountNumber = subaccount.id.number;
        const subaccountOwner = subaccount.id.owner;
        const subaccountId = SubaccountTable.uuid(
          subaccountOwner,
          subaccountNumber,
        );
        const assetPositions: AssetPosition[] = [];
        _.forEach(subaccount.asset_positions, (assetPosition) => {
          assetPositions.push({
            assetId: assetPosition.asset_id,
            index: Long.fromValue(assetPosition.index),
            quantums: bigIntToBytes(BigInt(assetPosition.quantums)),
          });
        });
        assetPositionMapping[subaccountId] = assetPositions;
      },
    );

  return assetPositionMapping;
}

export function getLiquidityTiersCreateObject(liquidityTier: LiquidityTier):
  LiquidityTiersCreateObject {
  return {
    id: liquidityTier.id,
    name: liquidityTier.name,
    initialMarginPpm: liquidityTier.initialMarginPpm.toString(),
    maintenanceFractionPpm: liquidityTier.maintenanceFractionPpm.toString(),
    basePositionNotional: quantumsToHuman(
      liquidityTier.basePositionNotional.toString(),
      QUOTE_CURRENCY_ATOMIC_RESOLUTION,
    ).toFixed(6),
  };
}

export function getAssetCreateObject(asset: Asset): AssetCreateObject {
  return {
    id: BigInt(asset.id).toString(),
    symbol: asset.symbol,
    atomicResolution: asset.atomicResolution,
    hasMarket: asset.hasMarket,
    marketId: asset.marketId,
  };
}

export function getAssetPositionCreateObject(
  subaccountId: string,
  assetPosition: AssetPosition,
  atomicResolution: number,
): AssetPositionCreateObjectWithId {
  return {
    id: AssetPositionTable.uuid(subaccountId, assetPosition.assetId.toString()),
    subaccountId,
    assetId: assetPosition.assetId.toString(),
    size: quantumsToHumanFixedString(bytesToBigInt(assetPosition.quantums).toString(),
      atomicResolution),
    isLong: getPositionIsLong(assetPosition),
  };
}

/**
 * Get `ClobPairs` from genesis.
 * @returns
 */
export function getClobPairsFromGenesis(): ClobPair[] {
  const clobPairs: ClobPair[] = genesis.app_state.clob.clob_pairs.map(
    (genesisClobpair): ClobPair => {
      return {
        ...genesisClobpair,
        subticksPerTick: genesisClobpair.subticks_per_tick,
        quantumConversionExponent: genesisClobpair.quantum_conversion_exponent,
        // Since the field in the proto is a uint64, this corresponds to a `BigInt` and not `number`
        stepBaseQuantums: Long.fromValue(genesisClobpair.step_base_quantums),
        minOrderBaseQuantums: Long.fromValue(genesisClobpair.min_order_base_quantums),
        perpetualClobMetadata: {
          perpetualId: genesisClobpair.perpetual_clob_metadata.perpetual_id,
        },
        // No status is set for the clob pairs in V4, so default to `ACTIVE` for now
        // TODO(DEC-600): Update this when clob pair statuses are fleshed out alongside governance
        status: ClobPair_Status.STATUS_ACTIVE,
      };
    },
  );

  return clobPairs;
}

export function getMarketParamsFromGenesis(): MarketParam[] {
  const markets: MarketParam[] = genesis.app_state.prices.market_params.map(
    (genesisMarketParam, index: number): MarketParam => {
      return {
        ...genesisMarketParam,
        minExchanges: genesisMarketParam.min_exchanges,
        minPriceChangePpm: genesisMarketParam.min_price_change_ppm,
        exchangeConfigJson: '',
        id: index,
      };
    },
  );

  return markets;
}

export function getMarketPricesFromGenesis(): MarketPrice[] {
  const marketPrices: MarketPrice[] = genesis.app_state.prices.market_prices.map(
    (genesisMarketPrice, index: number): MarketPrice => {
      return {
        ...genesisMarketPrice,
        id: index,
        exponent: genesisMarketPrice.exponent,
        price: Long.fromNumber(genesisMarketPrice.price),
      };
    },
  );
  return marketPrices;
}

/**
 * Gets LiquidityTiers from geneis.
 * @returns
 */
export function getLiquidityTiersFromGenesis(): LiquidityTier[] {
  const liquidityTiers: LiquidityTier[] = genesis.app_state.perpetuals.liquidity_tiers.map(
    (genesisLiquidityTier, index: number): LiquidityTier => {
      return {
        ...genesisLiquidityTier,
        basePositionNotional: Long.fromValue(genesisLiquidityTier.base_position_notional),
        initialMarginPpm: genesisLiquidityTier.initial_margin_ppm,
        maintenanceFractionPpm: genesisLiquidityTier.maintenance_fraction_ppm,
        impactNotional: Long.fromValue(genesisLiquidityTier.impact_notional),
        id: index,
      };
    },
  );
  return liquidityTiers;
}

/**
 * @description Given the initial `ClobPair` and `Perpetual` objects, generate a list of
 * `PerpetualMarketCreateObjects`.
 * @param clobPairs Initial `ClobPair` objects.
 * @param perpetuals Initial `Perpetual` objects.
 * @param marketParams Initial `MarketParam` objects.
 * @returns `PerpetualMarketCreateObjects` corresponding to each `ClobPair`, `Perpetual` and
 * `MarketParam`.
 * @throws `InvalidClobPairError` or `PerpetualDoesNotExist` error if the passed in `ClobPair` and
 * `Perpetual` objects are invalid.
 */
export function getPerpetualMarketCreateObjects(
  clobPairs: ClobPair[],
  perpetuals: Perpetual[],
  marketParams: MarketParam[],
  liquidityTiers: LiquidityTier[],
): PerpetualMarketCreateObject[] {
  const perpetualsById: { [id: number]: Perpetual } = _.keyBy(perpetuals, 'id');
  const marketParamsById: { [id: number]: MarketParam } = _.keyBy(marketParams, 'id');
  const liquidityTiersById: { [id: number]: LiquidityTier } = _.keyBy(liquidityTiers, 'id');

  // For each ClobPair, either create a `PerpetualMarketCreateObject` or return undefined if it's
  // not a perpetual ClobPair.
  const perpetualMarketsOrUndefined: (PerpetualMarketCreateObject | undefined)[] = clobPairs.map(
    (clobPair: ClobPair): PerpetualMarketCreateObject | undefined => {
      const perpetualForClobPair: Perpetual | undefined = getPerpetualForClobPair(
        clobPair, perpetualsById,
      );

      // If the ClobPair doesn't contain a reference to a perpetual, `undefined` is returned by
      // `getPerpetualForClobPair`
      if (perpetualForClobPair === undefined) {
        // Skip any ClobPairs that do not contain a reference to a perpetual
        return undefined;
      }

      const marketForPerpetual: MarketParam = getMarketParamForPerpetual(
        perpetualForClobPair, marketParamsById,
      );

      if (!(perpetualForClobPair.liquidityTier in liquidityTiersById)) {
        throw new LiquidityTierDoesNotExistError(
          perpetualForClobPair.liquidityTier,
          perpetualForClobPair.id,
        );
      }

      return getPerpetualMarketCreateObject(
        clobPair,
        perpetualForClobPair,
        marketForPerpetual,
      );
    },
  );

  // Filter out any undefined objects in the list.
  const perpetualMarkets: PerpetualMarketCreateObject[] = perpetualMarketsOrUndefined.filter(
    (
      perpetualMarket: PerpetualMarketCreateObject | undefined,
    ): perpetualMarket is PerpetualMarketCreateObject => {
      return perpetualMarket !== undefined;
    },
  );

  return perpetualMarkets;
}

/**
 * @description Get the `Perpetual` object referenced by a `ClobPair` object.
 * @param clobPair `ClobPair` object to get a `Perpetual` object for.
 * @param perpetualsById Map of `Perpetual` objects to their id.
 * @returns `Perpetual` referenced by the `ClobPair` or `undefined` if the `ClobPair` does not
 * reference a `Perpetual`.
 * @throws `PerpetualDoesNotExistError` if the `Perpetual` referenced by a `ClobPair` does not exist
 * in the passed in map of `Perpetual` objects.
 */
function getPerpetualForClobPair(
  clobPair: ClobPair,
  perpetualsById: { [id: number]: Perpetual },
): Perpetual | undefined {
  if (clobPair.perpetualClobMetadata !== undefined) {
    const perpetualIdForClobPair: number = clobPair.perpetualClobMetadata.perpetualId;

    if (!(perpetualIdForClobPair in perpetualsById)) {
      throw new PerpetualDoesNotExistError(perpetualIdForClobPair, clobPair.id);
    }

    return perpetualsById[perpetualIdForClobPair];
  } else {
    // TODO(DEC-689) Add support for spot markets with assets metadata
    return undefined;
  }
}

/**
 * @description Get the `MarketParam` object referenced by a `Perpetual` object.
 * @param perpetual `Perpetual` object to get a `MarketParam` object for.
 * @param marketParamsById Map of market ids to `MarketParam` objects.
 * @returns `MarketParam` referenced by the given `Perpetual`.
 */
function getMarketParamForPerpetual(
  perpetual: Perpetual,
  marketParamsById: { [id: number]: MarketParam },
): MarketParam {
  if (!(perpetual.marketId in marketParamsById)) {
    throw new MarketDoesNotExistError(perpetual.marketId, perpetual.id);
  }

  return marketParamsById[perpetual.marketId];
}

/**
 * @description Given a clob pair and it's corresponding perpetual, generate the `PerpetualMarket`
 * to create.
 * @param clobPair `ClobPair` for the `PerpetualMarket`.
 * @param perpetual `Perpetual` for the `PerpetualMarket`.
 * @param marketParam: `MarketParam` for the `Perpetual`.
 * @param liquidityTier `LiquidityTier` for the perpetual.
 * @returns `PerpetualMarketCreateObject` corresponding to the `Perpetual`, ClobPair` and
 * `MarketParam` passed in.
 * Note: This function assumes the passed in `Perpetual`, `ClobPair` and `MarketParam` match, and
 * does no validation for this.
 */
function getPerpetualMarketCreateObject(
  clobPair: ClobPair,
  perpetual: Perpetual,
  marketParam: MarketParam,
): PerpetualMarketCreateObject {
  return {
    id: perpetual.id.toString(),
    clobPairId: clobPair.id.toString(),
    ticker: perpetual.ticker,
    marketId: marketParam.id,
    status: clobStatusToMarketStatus(clobPair.status),
    // TODO(DEC-744): Remove base asset, quote asset.
    baseAsset: '',
    quoteAsset: '',
    // TODO(DEC-745): Initialized as 0, will be updated by roundtable task to valid values.
    lastPrice: '0',
    priceChange24H: '0',
    trades24H: 0,
    volume24H: '0',
    // TODO(DEC-746): Add funding index update events and logic to indexer.
    nextFundingRate: '0',
    // TODO(DEC-744): Remove base, incremental and maxPositionSize if not available in V4.
    basePositionSize: '0',
    incrementalPositionSize: '0',
    maxPositionSize: '0',
    openInterest: quantumsToHuman(
      perpetual.openInterest.toString(),
      perpetual.atomicResolution,
    ).toFixed(6),
    quantumConversionExponent: clobPair.quantumConversionExponent,
    atomicResolution: perpetual.atomicResolution,
    subticksPerTick: clobPair.subticksPerTick,
    minOrderBaseQuantums: Number(clobPair.minOrderBaseQuantums),
    stepBaseQuantums: Number(clobPair.stepBaseQuantums),
    liquidityTierId: perpetual.liquidityTier,
  };
}

/**
 * @description Given the initial `MarketParam` and `MarketPrice` objects, generate a
 * `MarketCreateObject`.
 * @param marketParam Initial `MarketParam` object.
 * @param marketPrice Initial `MarketPrice` object.
 * @returns `CreateMarketObject` corresponding to each `MarketParam`, `MarketPrice` pair passed in.
 * Note: This function assumes the passed in `MarketParam` and `MarketPrice` objects match, and
 * does no validation for this.
 */
function getMarketCreateObject(
  marketParam: MarketParam,
  marketPrice: MarketPrice,
): MarketCreateObject {
  return {
    id: marketParam.id,
    pair: marketParam.pair,
    exponent: marketParam.exponent,
    minPriceChangePpm: marketParam.minPriceChangePpm,
    oraclePrice: protocolPriceToHuman(marketPrice.price.toString(), marketPrice.exponent),
  };
}

/**
 * Converts a parts-per-million value to the string representation of the number. 1 ppm, or
 * parts-per-million is equal to 10^-6 (0.000001).
 * @param ppm Parts-per-million value.
 * @returns String representation of the parts-per-million value as a floating point number.
 */
export function ppmToString(ppm: number): string {
  return Big(ppm).div(1_000_000).toFixed(6);
}

/**
 * Calculates maintenance margin based on initial margin and maintenance fraction.
 * maintenance margin = initial margin * maintenance fraction
 * @param initialMarginPpm Initial margin in parts-per-million.
 * @param maintenanceFractionPpm Maintenance fraction in parts-per-million.
 * @returns Maintenance margin in parts-per-million.
 */
export function getMaintenanceMarginPpm(
  initialMarginPpm: number,
  maintenanceFractionPpm: number,
): number {
  return Big(initialMarginPpm).times(maintenanceFractionPpm).div(ONE_MILLION).toNumber();
}

function clobStatusToMarketStatus(clobPairStatus: ClobPair_Status): PerpetualMarketStatus {
  if (
    clobPairStatus !== ClobPair_Status.STATUS_UNSPECIFIED &&
    clobPairStatus !== ClobPair_Status.UNRECOGNIZED &&
    clobPairStatus in CLOB_STATUS_TO_MARKET_STATUS
  ) {
    return CLOB_STATUS_TO_MARKET_STATUS[clobPairStatus];
  } else {
    throw new InvalidClobPairStatusError(clobPairStatus);
  }
}

/**
 * Computes the unsettled funding for a position.
 *
 * To compute the net USDC balance for a subaccount, sum the result of this function for all
 * open perpetual positions, and subtract the sum from the latest USDC asset position for
 * this subaccount.
 *
 * @param position
 * @param latestFundingIndex
 * @param lastUpdateFundingIndex
 */
export function getUnsettledFunding(
  position: PerpetualPositionFromDatabase,
  latestFundingIndexMap: FundingIndexMap,
  lastUpdateFundingIndexMap: FundingIndexMap,
): Big {
  return Big(position.size).times(
    latestFundingIndexMap[position.perpetualId].minus(
      lastUpdateFundingIndexMap[position.perpetualId],
    ),
  );
}

/**
 * Get unrealized pnl for a perpetual position. If the perpetual market is not found in the
 * markets map or the oracle price is not found in the market, return 0.
 *
 * @param position Perpetual position object from the database, or the updated
 * perpetual position subaccountKafkaObject.
 * @param perpetualMarketsMap Map of perpetual ids to perpetual market objects from the database.
 * @returns Unrealized pnl of the position.
 */
export function getUnrealizedPnl(
  position: PerpetualPositionFromDatabase | UpdatedPerpetualPositionSubaccountKafkaObject,
  perpetualMarket: PerpetualMarketFromDatabase,
  marketsMap: MarketsMap,
): string {
  if (marketsMap[perpetualMarket.marketId] === undefined) {
    logger.error({
      at: 'getUnrealizedPnl',
      message: 'Market is undefined',
      marketId: perpetualMarket.marketId,
    });
    return Big(0).toFixed(CURRENCY_DECIMAL_PRECISION);
  }
  if (marketsMap[perpetualMarket.marketId].oraclePrice === undefined) {
    logger.error({
      at: 'getUnrealizedPnl',
      message: 'Oracle price is undefined for market',
      marketId: perpetualMarket.marketId,
    });
    return Big(0).toFixed(CURRENCY_DECIMAL_PRECISION);
  }
  return (
    Big(position.size).times(
      Big(marketsMap[perpetualMarket.marketId].oraclePrice!).minus(position.entryPrice),
    )
  ).toFixed(CURRENCY_DECIMAL_PRECISION);
}

/**
 * Gets the transfer type for a subaccount.
 *
 * If sender/recipient are both subaccounts, then it is a transfer_in/transfer_out.
 * If sender/recipient are wallet addresses, then it is a deposit/withdrawal.
 *
 * @param transfer
 * @param subaccountId
 */
export function getTransferType(
  transfer: TransferFromDatabase,
  subaccountId: string,
): TransferType {
  if (transfer.senderSubaccountId === subaccountId) {
    if (transfer.recipientSubaccountId) {
      return TransferType.TRANSFER_OUT;
    } else {
      return TransferType.WITHDRAWAL;
    }
  } else if (transfer.recipientSubaccountId === subaccountId) {
    if (transfer.senderSubaccountId) {
      return TransferType.TRANSFER_IN;
    } else {
      return TransferType.DEPOSIT;
    }
  }
  throw new Error(`Transfer ${transfer.id} does not involve subaccount ${subaccountId}`);
}
