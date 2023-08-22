import { logger } from '@dydxprotocol-indexer/base';
import { bigIntToBytes, bytesToBigInt, getPositionIsLong } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  Asset,
  AssetPosition,
  ClobPairStatus,
  LiquidityTier,
  MarketParam,
  MarketPrice,
  Perpetual,
  PerpetualMarketCreateEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';
import Long from 'long';
import { DateTime } from 'luxon';

import { CURRENCY_DECIMAL_PRECISION, ONE_MILLION, QUOTE_CURRENCY_ATOMIC_RESOLUTION } from '../constants';
import { setBulkRowsForUpdate } from '../helpers/stores-helpers';
import { InvalidClobPairStatusError } from '../lib/errors';
import { protocolPriceToHuman, quantumsToHuman, quantumsToHumanFixedString } from '../lib/protocol-translations';
import * as AssetPositionTable from '../stores/asset-position-table';
import * as SubaccountTable from '../stores/subaccount-table';
import {
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

type SpecifiedClobPairStatus =
  Exclude<ClobPairStatus, ClobPairStatus.CLOB_PAIR_STATUS_UNSPECIFIED> &
  Exclude<ClobPairStatus, ClobPairStatus.UNRECOGNIZED>;

const CLOB_STATUS_TO_MARKET_STATUS: Record<SpecifiedClobPairStatus, PerpetualMarketStatus> = {
  [ClobPairStatus.CLOB_PAIR_STATUS_ACTIVE]: PerpetualMarketStatus.ACTIVE,
  [ClobPairStatus.CLOB_PAIR_STATUS_CANCEL_ONLY]: PerpetualMarketStatus.CANCEL_ONLY,
  [ClobPairStatus.CLOB_PAIR_STATUS_PAUSED]: PerpetualMarketStatus.PAUSED,
  [ClobPairStatus.CLOB_PAIR_STATUS_POST_ONLY]: PerpetualMarketStatus.POST_ONLY,
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
 * @description Given a PerpetualMarketCreateEventV1 event, generate the `PerpetualMarket`
 * to create.
 */
export function getPerpetualMarketCreateObject(
  perpetualMarketCreateEventV1: PerpetualMarketCreateEventV1,
): PerpetualMarketCreateObject {
  return {
    id: perpetualMarketCreateEventV1.id.toString(),
    clobPairId: perpetualMarketCreateEventV1.clobPairId.toString(),
    ticker: perpetualMarketCreateEventV1.ticker,
    marketId: perpetualMarketCreateEventV1.marketId,
    status: clobStatusToMarketStatus(perpetualMarketCreateEventV1.status),
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
    openInterest: '0',
    quantumConversionExponent: perpetualMarketCreateEventV1.quantumConversionExponent,
    atomicResolution: perpetualMarketCreateEventV1.atomicResolution,
    subticksPerTick: perpetualMarketCreateEventV1.subticksPerTick,
    minOrderBaseQuantums: Number(perpetualMarketCreateEventV1.minOrderBaseQuantums),
    stepBaseQuantums: Number(perpetualMarketCreateEventV1.stepBaseQuantums),
    liquidityTierId: perpetualMarketCreateEventV1.liquidityTier,
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

function clobStatusToMarketStatus(clobPairStatus: ClobPairStatus): PerpetualMarketStatus {
  if (
    clobPairStatus !== ClobPairStatus.CLOB_PAIR_STATUS_UNSPECIFIED &&
    clobPairStatus !== ClobPairStatus.UNRECOGNIZED &&
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
